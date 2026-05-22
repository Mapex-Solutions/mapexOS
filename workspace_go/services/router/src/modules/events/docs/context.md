# Bounded Context: Events (Router)

**Service:** router
**Module path:** `src/modules/events/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

This bounded context is the core routing pipeline of the `router` service. It consumes incoming route execution messages (one per asset event or health-status transition), resolves the target asset from a TieredCache, selects the applicable route groups based on `eventSource` (`assetEvent` → `asset.RouteGroupIds`, `healthStatus` → `asset.HealthMonitor.{Offline,Online}RouteGroupIds`), evaluates per-router match rules, enriches the payload with asset and template metadata, and fans the result out to kind-specific downstream subjects (save_event, lake_house, notification, trigger, workflow). It also owns the FANOUT consumers that keep the local asset and template caches coherent with upstream Assets changes.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Route Execution | A single inbound `mapexos.route.execute` message carrying `orgId`, `assetUUID`, `event`, `eventTrackerId`, `eventSource` | Trigger execution (downstream, owned by `triggers` service) |
| eventSource | Discriminator that selects which route-group list on the asset is used (`assetEvent` \| `healthStatus`) | HTTP-Gateway "source" (ingestion channel) |
| Router (kind) | One entry inside a RouteGroup: `save_event`, `lake_house`, `notification`, `trigger`, `workflow` | `router` the microservice |
| Match | Per-router rule set evaluated against the event payload with `all`/`any` policy | Mongo filter on route groups |
| Enriched Event | Generic base (asset context + tracking IDs) + kind-specific fields added by `addRouterData` | Raw sensor event payload |
| Tier (L0/L1/L2) | TieredCache layers — RAM / Disk / MinIO (source of truth) for `AssetReadModel` | Route group cache (Redis, owned by `routegroups`) |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| SaveEvent dispatch | `mapexos.events.save` | enriched event (map) with `source=asset`, `threadId` | events |
| LakeHouse dispatch | `mapexos.events.lake_house` | enriched event + `lakeHouseId`, `metadata` | events / lakehouse sink (inferred) |
| Notification dispatch | `mapexos.events.notification` | enriched event + `notificationId`, `metadata` | events / notification consumer (inferred) |
| Trigger dispatch | `mapexos.trigger.router.execute` | enriched event + `triggerId`, `source=router`, `payload` | triggers |
| Workflow dispatch | `mapexos.workflow.execution.router` | enriched event + `mode`, `data`, `metadata` | workflow |
| Router history | `mapexos.events.router` | `eventTypes.RouterHistoryEvent` (per-router results + conditions) | events |

MsgId convention for JetStream dedup: `{eventTrackerId}-{routerIdx}` for dispatches, `{eventTrackerId}-history` for history.

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Route execute | `mapexos.route.execute` (stream `ROUTE-GROUPS`, WorkQueue + retry/DLQ) | `{orgId, assetUUID, event, eventTrackerId, eventSource, pathKey?}` (inferred contract) | http_gateway, assets (healthmonitor) |
| Asset invalidate | `mapexos.fanout.asset.invalidate` (stream `FANOUT`) | `eventTypes.AssetInvalidateEvent` = `{orgId, assetUUID}` | assets |
| Template invalidate | `mapexos.fanout.template.invalidate` (stream `FANOUT`) | `{orgId, templateId}` (inferred from ParseInvalidateEvent) | assets |

## Driving Ports (what can call this module)

- NATS WorkQueue consumer on `mapexos.route.execute` (`ROUTE-GROUPS` stream, batched, retry/DLQ via `DefaultRetryPolicy`).
- NATS FANOUT subscription on `mapexos.fanout.asset.invalidate` (one delivery per replica, ephemeral).
- NATS FANOUT subscription on `mapexos.fanout.template.invalidate` (one delivery per replica, ephemeral).
- No HTTP surface — this module is purely message-driven.

## Driven Ports (what this module requires)

- `ports.EventServicePort` (self, exposed to consumers).
- `ports.MatchEvaluatorPort` → `domainServices.MatchEvaluator` (pure, no I/O).
- `routegroupPorts.RouteGroupServicePort` (cross-module, for `GetRouteGroupEntityById`).
- `common.TieredCache` named `"assets"` (L0 RAM / L1 Disk / L2 MinIO `AssetReadModel`).
- `templateCache.TemplateCache` (wraps TieredCache `"templates"` for enrichment).
- `natsModel.CorePublisher` (fire-and-forget `PublishCore` + batched `FlushConnection`).
- `*bootstrap.RouterMetrics` (Prometheus counters/histograms for cache tiers, dispatches, match evaluations).

## Invariants and Business Rules

- Route IDs are NEVER trusted from the payload — the asset in `TieredCache("assets")` is the single source of truth for `RouteGroupIds` / `HealthMonitor.{Offline,Online}RouteGroupIds`.
- `eventSource=healthStatus` routing is restricted to router kinds in `HealthStatusAllowedRouterKinds = {trigger, workflow}`; other kinds are silently skipped at debug level.
- Missing `orgId`, `assetUUID`, `event`, or invalid JSON → `msg.Reject(reason)` (DLQ, no retry). Processing failures after validation → `msg.Nack(err)` (retry with backoff).
- `ProcessEventBatch` follows a strict 3-phase pipeline: parallel route processing → single `FlushConnection` → sequential ACK/Nack/Reject. Publishes are buffered until the flush.
- FANOUT cache invalidations clear only L0 + L1 for key `{orgId}/{assetUUID}` (or `{orgId}/{templateId}`); L2 is authoritative and refilled lazily on next read.
- Every dispatched event carries `eventTrackerId` for end-to-end tracing and a per-dispatch `executionId` UUID.

## Known Cross-Context Interactions

- Consumes `AssetReadModel` projections produced by the **assets** service (via MinIO L2 + FANOUT invalidations).
- Consumes `AssetTemplate` projections produced by the **assets** service (same mechanism) for event enrichment (`templateName`, `templateDescription`).
- Delegates RouteGroup lookups to the sibling **routegroups** module (same service) via its port — only `application/ports` types are imported, never `domain/entities`.
- Publishes to **events** (save / history / lake_house / notification — inferred for the latter two), **triggers**, and **workflow** — all via NATS subjects declared in `application/constants/subjects.go`.
- Receives `route.execute` primarily from **http_gateway** and from **assets** (healthmonitor) on health-status transitions (inferred from `eventSource` discriminator values).
- Note: events Go does NOT cache assets — only AssetTemplate. The FANOUT pipeline split is: router consumes both asset+template invalidate; events Go consumes only template.invalidate; js-executor consumes both.
