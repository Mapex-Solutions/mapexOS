# Bounded Context: Retention

**Service:** events
**Module path:** `src/modules/retention/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns the retention-policy catalog (MongoDB) that tells the rest of the platform how many days each ClickHouse table should keep data for. Provides the per-org, per-table lookup that `events` stamps on each row as `retentionDays`, seeds platform-level policies (asset_status_history) at startup, bootstraps 8 defaults per new organization, and applies platform-level DDL (TTL) on `asset_status_history` when operators edit the UI slider.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| RetentionPolicy | Mongo document keyed by `{orgId, type}` (unique compound), holding `retentionDays`, `enabled`, `name`, `pathKey` | ClickHouse `TTL` clause — policy drives TTL but is not TTL |
| Type | ClickHouse logical table name the policy targets (e.g., `events`, `eventsRaw`, `asset_status_history`) | Mongo collection name |
| Platform policy | Policy with no `orgId` — global default (currently only `asset_status_history`) | Org-scoped policy (8 defaults per org) |
| RetentionDays | uint16 day count returned by `GetRetentionDays(orgId, tableName)` with fallback to `DefaultRetentionDays = 1` | TTL seconds in ClickHouse (days × 86400) |
| CacheKeyPrefix | `RETENTION_POLICY:{orgId}:{type}`, TTL 24h | NATS subject |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — (no NATS publishes; side effect is a ClickHouse `ALTER TABLE ... MODIFY TTL` on `asset_status_history`) |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Organization created | `mapexos.events.organization.created` (stream `MAPEXOS`) | `contracts/services/mapexIam/organizations.OrganizationCreatedEvent` | mapexos |

## Driving Ports (inbound)
- NATS consumer `{service_name}-retention-org-created` on stream `MAPEXOS`, subject `mapexos.events.organization.created`, queue group `{service_name}-RETENTION-ORG-CREATED-GROUP`. Uses `DefaultRetryPolicy` and DLQ `eventType=organization.created`. Invokes `CreateDefaultPolicies(orgId, pathKey)`.
- HTTP under `/api/v1/retention` (auth + context injected):
  - `GET /` list (paginated/filtered) — perm `retention.list`.
  - `GET /:retentionPolicyId` — perm `retention.read`.
  - `PUT /` upsert by `{orgId, type}` — perm `retention.update`.
  - `DELETE /:retentionPolicyId` — perm `retention.update`.
- Internal (in-process) programmatic port: `RetentionServicePort.GetRetentionDays` called by the events module at write time.

## Driven Ports (outbound)
- `domain/repositories.RetentionRepository` (MongoDB) — implemented by `infrastructure/persistence/mongo`. Supports `Create`, `FindById`, `FindByOrgIdAndType`, `Upsert` (key `{orgId, type}`), `DeleteById`, `FindWithFilters`.
- `domain/repositories.CacheRepository` — Redis client provided as the DIG-named `"app"` Redis instance; used for cache-aside reads (`GetOrSet` with `RetentionCacheTTL = 24h`).
- ClickHouse driver — used only by `ApplyAssetStatusHistoryTTL` to issue `ALTER TABLE ... MODIFY TTL`.

## Invariants and Business Rules
- Unique per org+type: `{orgId, type}` is the upsert key. Platform policies use `orgId = nil`.
- `GetRetentionDays` is cache-aside with a 24h TTL and always returns a value — fallback is `DefaultRetentionDays = 1` day if no policy exists.
- `CreateDefaultPolicies` seeds the 8 org-level retention types when a new organization is created; idempotent via upsert semantics.
- `SeedPlatformPolicies` runs at module boot and inserts the `asset_status_history` row (7-day default, 1–90 range); re-running is safe.
- `asset_status_history` is platform-scoped (no orgId) and bound to `[AssetStatusHistoryMinDays, AssetStatusHistoryMaxDays] = [1, 90]`; the upsert path calls `ApplyAssetStatusHistoryTTL` to push the change into ClickHouse DDL.
- On org-created consumer failures: `Nack` for retry; malformed / empty `orgId|pathKey` events are ACK'd and skipped (logged).

## Known Cross-Context Interactions
- Consumes org lifecycle from **mapexos** (organization.created).
- Called in-process by **events** (this service) to stamp `RetentionDays` per row at write time.
- Governs ClickHouse TTL for the **asset_status** module's `asset_status_history` table via platform policy + `ALTER TABLE MODIFY TTL`.
- Exposed to the frontend's retention settings UI (slider for `asset_status_history`, CRUD for per-org policies).
