# Bounded Context: DataSources

**Service:** http_gateway
**Module path:** `src/modules/datasources/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose

Owns the lifecycle and configuration of `DataSource` records — the webhook/API ingestion endpoints that external systems use to push IoT events into MapexOS. A DataSource stores the protocol (http/mqtt), mode (pull/push), authentication strategy (apiKey, jwt, oauth2, ip_whitelist, none), rate-limit and working-hours policies, and the `AssetBind` rule that associates incoming payloads with a target Asset. This module provides CRUD over those records and keeps a Redis cache so the hot path in the `events` module can resolve a DataSource by ID on every request without hitting MongoDB.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| DataSource | Configured ingestion endpoint (auth, bind, limits) persisted in MongoDB `data_sources` collection | `router` RouteGroup — that matches events to rules post-ingestion |
| AssetBind | Rule that maps an incoming payload to an Asset (type + UUID field path) | Asset entity itself, which lives in the `assets` service |
| Auth | Per-DataSource auth config (apiKey / jwt / oauth2 / ip_whitelist / none) | Platform user/JWT auth on the admin HTTP routes |
| PathKey | Hierarchical org path used for range queries under an org subtree (inferred) | MongoDB ObjectId of the org itself |
| Mode | pull or push, declares whether the gateway fetches or is called | HTTP verb on the admin route |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — |

This module does not publish NATS events. State changes are persisted to MongoDB and cached in Redis only.

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | — |

No NATS consumers. All traffic enters via the admin HTTP API.

## Driving Ports (what can call this module)

- HTTP admin routes under `/api/v1/data_sources` (GET list, POST create, GET/:id, PATCH/:id, DELETE/:id), protected by platform `AuthMiddleware` and per-route `RequirePermission` (`DatasourceList/Create/Read/Update/Delete`).
- `DataSourceServicePort` consumed in-process by the `events` module (middleware calls `GetDataSourceById` for every incoming webhook).

## Driven Ports (what this module requires)

- `DataSourceRepository` — MongoDB adapter over the `data_sources` collection.
- `CacheRepository` — Redis-backed `Cache` + `CacheGetOrSetEx` (AppCache) used to memoize DataSource reads with a 24h TTL under key prefix `DATA_SOURCE` (inferred: used on reads via `GetOrSet`).
- Coverage middleware `InjectRequestContext` for org-scoped listing with `PathKey` hierarchical filtering.

## Invariants and Business Rules

- Every DataSource is multi-tenant: `orgId` and `pathKey` are assigned from `RequestContext` on create and are the primary filter on list.
- `Auth.Type` is one of `apiKey | jwt | oauth2 | ip_whitelist | none`; the matching sub-struct (`APIKey`, `JWT`, `OAuth2`, `IPWhitelist`) must be present for the event-ingestion middleware to authorize traffic.
- `Mode` is `pull | push` and `Protocol` is `http | mqtt` (inferred from entity tags; no runtime enforcement shown).
- Cache TTL for a DataSource record is 24h because configuration changes are rare; updates SHOULD invalidate the cache (inferred — not verified in service code shown).
- List results are projected and paginated through the shared `PaginatedResult[T]` model and `orgfilter.BuildOrgFilter` for `includeChildren` PathKey range queries.

## Known Cross-Context Interactions

- The `events` module (same service) resolves a DataSource on every webhook via `DataSourceServicePort.GetDataSourceById` inside `CustomAuthMiddleware`, then applies the configured auth strategy before publishing to NATS.
- `AssetBind.Data.AssetId` / `UUIDField` is forwarded to `js-executor` in the event payload so it can look up the target Asset in the Assets cache (source of truth for asset metadata lives in the `assets` service).
- Admin routes depend on the shared `coverage` middleware and `permissions/http_gateway` catalog; changes to those contracts are cross-service.
