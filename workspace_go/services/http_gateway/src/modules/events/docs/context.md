# Bounded Context: Events

**Service:** http_gateway
**Module path:** `src/modules/events/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

Ingestion edge for HTTP-delivered IoT events. Accepts webhook/API payloads on `POST /api/v1/events?ds={dataSourceId}`, resolves the target `DataSource` from the `datasources` module, enforces its per-DataSource auth policy (apiKey/jwt/oauth2/ip_whitelist/none), stamps the payload with a tracker UUID, and hands it off to the downstream pipeline through NATS. Success publishes to `js-executor`; auth failures publish to the raw event stream for security monitoring. This module owns no domain state — it is a thin, metered adapter between HTTP and NATS.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Event | Arbitrary `map[string]any` body posted by an external system for one DataSource | Domain event in DDD — here it is a raw ingestion payload |
| eventTrackerId | UUID minted at ingestion to correlate the event across http_gateway → js-executor → router → events (storage) | `threadId` (= DataSource ID) used for grouping by source |
| Auth failure event | Raw event published with `success=false` and error message for security audit, even when the request is rejected | Admin auth errors on `/api/v1/data_sources` admin routes |
| sourceType | Marker in the js-executor payload (`"http"`) telling the executor which gateway produced the event | DataSource `protocol` field (http/mqtt) |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| ProcessorExecute | `mapexos.processor.js.execute` | `contracts/services/http_gateway/events.ProcessorExecutePayload` | js-executor |
| RawEvent (auth failure) | `mapexos.events.raw` | `contracts/services/events/events.RawEventDTO` with `Success=false` | events service (ClickHouse `events_raw`) |
| AssetHeartbeatV1 (TKT-2026-0034) | `mapexos.asset.heartbeat.{orgId}` (constant `hmContract.SubjectAssetHeartbeat` + `.{orgId}`) | `{orgId, assetUUID, pathKey, ts}` (matches js-executor's implicit publish shape so the consumer is origin-agnostic) | assets/healthmonitor consumer |

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | — |

No NATS consumers. All input arrives over HTTP.

## Driving Ports (what can call this module)

- HTTP route `POST /api/v1/events` with query DTO `EvenIdentificationDto` (`ds` = DataSource ID), guarded by `CustomAuthMiddleware`.
- HTTP route `POST /api/v1/heartbeat` (TKT-2026-0036 reformulation) — body `{ "assetUUID": "<v>" }` (required), returns `204 No Content` on success. Same query DTO + auth chain as `/events` (resolves DataSource via `ds` query param). Used by explicit-mode HTTP-protocol assets (`HealthMonitorConfig.heartbeatMode='explicit'` + `protocol.type='http'`) to send heartbeats. orgId and pathKey come from the resolved DataSource — never from the body — so a compromised body cannot spoof a different tenant. The legacy `AssetBind.Type='fixedAssetId'` constraint was removed; any DataSource shape works. Errors: `404` if dataSource is missing or has no orgId; `422` if `assetUUID` is empty; `403` if DataSource is disabled; `500` on publish failure.
- `EventServicePort` is also called from the auth middleware itself (to publish `PublishAuthFailure` when validation rejects a request).

## Driven Ports (what this module requires)

- `NatsBus` (`packages/infrastructure/nats`) — publishes to `mapexos.processor.js.execute` and `mapexos.events.raw`.
- `dsPorts.DataSourceServicePort` (cross-module, in-process) — loads the DataSource config in the auth middleware.
- `bootstrap.HttpGatewayMetrics` — Prometheus counters/histograms for `EventsProcessed`, `EventsPublished`, `EventAuthTotal/Duration/Failures`, `EventPayloadSize`, `EventProcessingDuration`, plus `HeartbeatsTotal{status}` and `HeartbeatDuration` (TKT-2026-0034).

## Invariants and Business Rules

- Every successful request publishes exactly once to `mapexos.processor.js.execute`; publish failure surfaces as `INTERNAL_SERVER_ERROR` and increments error counters — no retries in this layer.
- Auth failures MUST fire-and-forget a `RawEventDTO{Success:false}` to `mapexos.events.raw` (goroutine, best-effort) before returning `401`.
- Payload sent to js-executor is minimized to `orgId` + `assetBind`; `name`, `description`, `pathKey` are deliberately omitted because js-executor reads them from the Asset cache (source of truth).
- `eventTrackerId` is a fresh UUID v4 per request and is the only cross-pipeline correlation id generated here.
- Unsupported `Auth.Type` values short-circuit with `401` and are still audited via `PublishAuthFailure`.
- Request body is parsed twice on the failure path (once lazily in `parseEventBody`, once in the handler) — acceptable because the failure path does not reach the handler.
- **Disabled DataSource gate (TKT-2026-0036)**: `CustomAuthMiddleware` rejects any request whose resolved DataSource has `Enabled=false` (or nil) with `403` BEFORE running the auth switch. Reuses the existing `EventAuthTotal`/`EventAuthDuration`/`EventAuthFailures` metrics with the new label value `type="disabled"`. `PublishAuthFailure` still fires for the security audit.
- **Heartbeat body shape (TKT-2026-0036)**: orgId and pathKey are derived from the resolved DataSource (server-side, `c.Locals`). The body provides only `assetUUID`. The published payload to `mapexos.asset.heartbeat.{orgId}` matches the shape js-executor publishes implicitly so the assets/healthmonitor consumer is origin-agnostic.

## Known Cross-Context Interactions

- Depends on `datasources` (same service) for every request to resolve auth config and asset-bind.
- Feeds the `js-executor` (workspace_js) pipeline on `mapexos.processor.js.execute`. The Go-side contract is declared in `packages/contracts/services/http_gateway/events` (`ProcessorExecutePayload`, `SubjectProcessorJSExecute`). TS-side reciprocity in `workspace_js/packages/schemas/src/services/http_gateway/events/` is pending (tracked as follow-up).
- Feeds the `events` service `events_raw` ClickHouse storage via `mapexos.events.raw` using `contracts/services/events/events.RawEventDTO` (authoritative cross-service contract).
- Metrics are scraped by the platform Prometheus stack; label cardinality is bounded by `authType`, `status`, and fixed subject strings.
