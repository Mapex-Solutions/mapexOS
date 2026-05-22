# Bounded Context: Assets

**Service:** assets
**Module path:** `src/modules/assets/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-11

## Purpose
Owns the lifecycle of IoT Assets — the physical/logical devices that emit events in MapexOS. MongoDB is the source of truth for asset configuration (protocol, MQTT credentials, current X.509 cert metadata, route groups, health-monitor config, geo). The module exposes CRUD over HTTP for the frontend and a single internal read-model endpoint (`GET /internal/assets/:assetUUID`) that downstream consumers — Router, JS-Executor, Events, the mapex-mqtt-broker plugin — hit as the L3 fallback of their TieredCaches. It publishes a denormalized `AssetReadModel` to MinIO (L2) carrying everything every consumer needs, including `Protocol.Mqtt.PasswordHash` and `CurrentCert.Serial` for the broker plugin's local CONNECT decisions. Cache invalidation rides NATS FANOUT.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Asset | A single IoT device entity with protocol (HTTP/MQTT), credentials, org scoping, optional health-monitor config, and the active cert metadata when MQTT cert auth is used | `AssetTemplate` (classification/scripts) and `AssetReadModel` (denormalized L2 payload) |
| AssetUUID | Business identifier used across the platform (NATS, MQTT, cache keys) | MongoDB `_id` which is internal |
| PathKey | Hierarchical org path used for multi-tenant range queries | `OrgID` which is the single owning org |
| RouteGroupIds | References to Router service groups the asset's events are routed to | `HealthMonitor.OfflineRouteGroupIds/OnlineRouteGroupIds` which target health transitions |
| AssetReadModel | Denormalized payload stored at `{orgId}/{assetUUID}.json` in MinIO bucket `mapex-assets` — single payload that carries the broker plugin's `Protocol.Mqtt.PasswordHash` + `CurrentCert.Serial` as well as the public asset state read by Router / JS-Executor / Events | The Mongo `Asset` entity (source of truth) |
| CurrentCert | The asset's currently-active MQTT device cert metadata (serial, fingerprint, expiry). Empty for password-only assets or after revoke. | `mqttRevokedCertificates` (30-day TTL audit collection owned by `mqttcerts`) |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| Asset cache invalidate | `mapexos.fanout.asset.invalidate` (stream `FANOUT`) | JSON with `assetUUID`, `orgId` (inferred from `publishAssetInvalidate`) | Router, JS-Executor, Events, mapex-mqtt-broker plugin |
| Asset read model write | MinIO bucket key `mapex-assets/{orgId}/{assetUUID}.json` (object storage, not NATS) | `AssetReadModel` from `contracts/services/assets/assets` | Router, JS-Executor, Events, mapex-mqtt-broker plugin via TieredCache L2 |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | None (this module has no NATS consumers) |

## Driving Ports (inbound — who calls this module)
- HTTP `/api/v1/assets` (JWT auth): `GET /`, `POST /`, `GET /counter`, `GET /:assetId`, `PATCH /:assetId`, `DELETE /:assetId` — frontend/API clients
- HTTP `GET /internal/assets/:assetUUID` (API-Key auth): the L3 read-model fallback for every consumer's TieredCache (Router, JS-Executor, Events, mapex-mqtt-broker plugin). On hit, the handler also repopulates the L2 MinIO entry inline so the next reader hits the cache.
- Go port `AssetServicePort` consumed in-process by `mqttcerts` (for `GetByUUID` on issue/revoke) and `healthmonitor` (for asset enrichment on heartbeat/presence handling)

## Driven Ports (outbound — what this module requires)
- `repositories.AssetRepository` — MongoDB persistence (`infrastructure/persistence/mongo`)
- `ports.AssetStoragePort` — MinIO writer for the asset read model (`infrastructure/storage/minio`); writes the single `AssetReadModel` carrying `PasswordHash` + `CurrentCert`
- `ports.RouteGroupPort` — HTTP client to Router service for name lookup and router-kind validation (`infrastructure/httpclient/router`)
- `common.AppCache` — Redis DB 0 for the per-org asset counter cache
- `natsModel.Fanout` (name `core`) — publishes invalidation messages
- `assettemplatePort.AssetTemplateRepository` — cross-module dependency to enrich responses with classification
- `healthPorts.HealthRepository` — read-only enrichment with live `healthStatus`/`lastSeenAt`; list path uses `GetLastSeenBatch` + `IsAlertedBatch` (one Redis round-trip per page, not per asset)
- `healthPorts.HealthLifecyclePort` — clears Redis health state on `HealthMonitor.Enabled` true→false transition and on delete

## Invariants and Business Rules
- `OrgID` + `PathKey` MUST come from `RequestContext` on create; asset is always scoped to one org
- `HealthMonitor.Enabled=true` requires at least one `OfflineRouteGroupIds` or `OnlineRouteGroupIds` entry, and every referenced group must contain only `trigger` or `workflow` router kinds (validated via `RouteGroupPort.GetRouterKindsByIds` — fails fast with 422 before persistence)
- Create/Update/Delete MUST invalidate the counter cache key `counter:assets:{orgId}` and fan out `mapexos.fanout.asset.invalidate`
- The MinIO `AssetReadModel` MUST be written before fanout invalidation publishes — consumers re-fetching after the invalidate must see fresh data
- `AssetReadModel.Protocol.Mqtt.PasswordHash` and `AssetReadModel.CurrentCert` MUST mirror the entity's current state on every write — the broker plugin makes CONNECT decisions off these fields and has no fallback path
- `AssetUUID` is the business identity used across NATS/MQTT/cache — Mongo `_id` is internal only
- `healthStatus` and `healthStatusChangedAt` MUST be written atomically on every transition — the repository method `UpdateHealthStatusWithChangedAt` issues a single Mongo `$set` so downstream readers can never observe a divergent pair
- DELETE is HARD — drops the asset entity along with its `currentCert` subdoc and the L2 `{orgId}/{assetUUID}.json`. `mqttRevokedCertificates` rows for the asset are cleaned up by `mqttcerts` on the same delete (LGPD: no orphan cert metadata after asset removal).

## Known Cross-Context Interactions
- Router service: synchronous HTTP for route-group name/kind resolution; async FANOUT for cache invalidation
- The asset create/edit wizard's Health step (HealthMonitoringSection) requests `GET /api/v1/route-groups?kinds=trigger&kinds=workflow` so users can only select RouteGroups acceptable to `validateHealthMonitorConfig` (mirrors `HealthStatusAllowedRouterKinds`). Step 3 (regular RouteGroupIds) remains unfiltered.
- JS-Executor, Events: consume the MinIO read model and the FANOUT invalidation via their TieredCache.
- mapex-mqtt-broker plugin (out-of-process): consumes the SAME MinIO read model via its own TieredCache (L1 Pebble → L2 MinIO → L3 GET /internal/assets/:uuid) and projects out `PasswordHash` + `CurrentCert.Serial` to decide every CONNECT locally — no HTTP auth callout, no separate bucket.
- HealthMonitor module (same service): writes live `healthStatus` that this module reads at response time. On Asset Update with `HealthMonitor.Enabled` transitioning true→false (and on Asset Delete), the service calls `HealthLifecyclePort.ClearAssetState` to remove all Redis health state for the asset AND resets Mongo `healthStatus` to `unknown` with `healthStatusChangedAt` updated.
- AssetTemplates module (same service): repository read for classification enrichment in responses
- MqttCerts module (same service): writes `Asset.CurrentCert` on issue and clears it on revoke, then fans out `mapexos.fanout.asset.invalidate` so every consumer (broker plugin included) picks up the new cert state via the standard TieredCache flow.
