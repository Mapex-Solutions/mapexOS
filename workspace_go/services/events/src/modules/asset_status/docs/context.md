# Bounded Context: Asset Status

**Service:** events
**Module path:** `src/modules/asset_status/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose
Persists asset connectivity transitions (offline/online) as an immutable history in ClickHouse (`asset_status_history`) and serves cursor-paginated queries over it. The module is a sink for the healthmonitor's transition events and a read API for the UI's asset connectivity timeline. It does not decide when an asset is offline — that is the assets service's responsibility.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| AssetStatusEvent | One ClickHouse row representing a single offline/online transition for an asset | Live health state in assets service (Redis `HealthRepository`) — ephemeral, not a row |
| EventType | Always `"offline"` or `"online"` for this module | `event_type` in the main `events` table (telemetry/alarm/command) |
| MissCount | How many consecutive heartbeat scans missed before the offline decision fired (offline transitions only) | Retry count on the NATS consumer |
| ThresholdMinutes | Offline threshold in effect at transition time, carried from `HealthMonitorConfig` | A general-purpose grace period |
| EventId | UUID generated once per transition, shared with the matching `route.execute` message (inferred) | `event_tracker_id` used by the main processed-events pipeline |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — (this module does not publish; it is a terminal consumer) |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Asset connectivity transition | `mapexos.events.asset_status_save` (stream `EVENTS-ASSET-STATUS`) | FLAT JSON mapped to `domain/entities.AssetStatusEvent` (inferred — emitted by assets healthmonitor alert publisher) | assets (healthmonitor) |

## Driving Ports (inbound)
- NATS consumer `{service_name}-asset-status-save` on stream `EVENTS-ASSET-STATUS`, subject `mapexos.events.asset_status_save`, queue group `{service_name}-EVENTS-ASSET-STATUS-GROUP`. Uses `DefaultRetryPolicy` and DLQ with `eventType=asset-status`.
- HTTP `GET /api/v1/events/connectivity_history` (org-wide list) — perm `events.asset_status.list`.
- HTTP `GET /api/v1/events/assets/:assetUUID/connectivity_history` (asset-scoped) — perm `events.asset_status.list`. Both routes use `InjectRequestContext` for org filtering and accept `from`/`to`/`eventType` query filters with cursor pagination.

## Driven Ports (outbound)
- `domain/repositories.AssetStatusRepository` — implemented by `infrastructure/persistence/clickhouse`. Methods: `BulkInsert` (batch write on consumer path), `FindByCursor` (time-cursor-ordered query for HTTP).
- ClickHouse driver (`driver.Conn`) injected at module init; no direct DB access from the application layer.

## Invariants and Business Rules
- Every consumed message is individually Ack/Nack/Reject'd by the service — `ProcessAssetStatusBatch` always returns `nil` so the NATS lib does not re-handle the batch.
- Invalid payloads are Rejected (→ DLQ) immediately; insertion failures are Nack'd for retry with backoff (inferred from the shared batch contract in this service).
- `BulkInsert` with an empty slice is a documented no-op, so the service need not pre-filter empty batches.
- TTL on `asset_status_history` is governed by the retention module (platform-level policy, default 7 days, range 1–90); this module never issues DDL.
- `LastSeenAt` and `MissCount` are only meaningful on `offline` transitions.

## Known Cross-Context Interactions
- Receives transition events from **assets/healthmonitor** (source of truth for offline/online decisions).
- TTL for the backing ClickHouse table is applied by **retention** via `ApplyAssetStatusHistoryTTL` (platform policy `asset_status_history`).
- Consumed by the frontend asset connectivity timeline (UI) via the two HTTP endpoints.
