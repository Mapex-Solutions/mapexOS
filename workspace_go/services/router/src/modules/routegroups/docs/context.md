# Bounded Context: RouteGroups (Router)

**Service:** router
**Module path:** `src/modules/routegroups/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-26

## Purpose

This bounded context owns the lifecycle and retrieval of `RouteGroup` aggregates — the configuration artifact that tells the router how to fan an event out to one or more downstream kinds (`save_event`, `lake_house`, `notification`, `trigger`, `workflow`), each with its own optional `MatchConfig` (`all`/`any` policy over `MatchRule` conditions). The module exposes a REST CRUD surface for UI/admin use, an internal MS-to-MS bulk lookup endpoint, and a service port used by the sibling `events` module to resolve a group entity by ID at dispatch time. Persistence is MongoDB (`routegroups` collection); a Redis cache-aside layer keeps hot route groups and an organization-scoped counter fresh.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| RouteGroup | Aggregate: named, versioned, enabled-flaggable collection of `Router` entries with multi-tenant scoping (`orgId` + `pathKey`) | Route (a single `Router` entry inside the group) |
| Router (entry) | Single dispatch descriptor inside a RouteGroup: `kind` + optional `Match` + kind-specific data (`LakeHouse`/`Notification`/`Trigger`/`SaveEvent`/`Workflow`) | `router` the microservice |
| isSystem | `true` → MAPEX-wide global template, visible to every tenant, no `orgId`/`pathKey` | `isTemplate` |
| isTemplate | `true` → shared template usable by Vendor/Customer ancestors, scoped by `pathKey` | `isSystem` |
| MatchConfig | `{policy: "all"\|"any", rules: []MatchRule}` attached to a Router to gate dispatch | Event source filter (owned by `events` module) |
| Coverage / RequestContext | Middleware-injected data (`OrgContext`, `OrgContextData.PathKey`) used for hierarchical org filtering on list/count | Auth token claims |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — |

This module does not publish NATS events. Mutations are persisted to MongoDB and reflected into Redis; cross-replica propagation is handled implicitly via the shared cache (inferred — no NATS invalidation is emitted from here).

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | — |

No NATS consumers — this module is HTTP-driven and in-process (called by the `events` module via its service port).

## Driving Ports (what can call this module)

- HTTP routes under `/api/v1/route_groups` (Auth middleware + per-verb permissions from `permissions/router`): `GET /`, `GET /counter`, `POST /`, `GET /:routeGroupId`, `PATCH /:routeGroupId`, `DELETE /:routeGroupId`.
- GET `/api/v1/route-groups` accepts an optional `kinds` query parameter (repeatable; valid values: `lake_house|notification|trigger|save_event|workflow`). When present, the result set is restricted to RouteGroups whose every router.kind is in the requested set AND that have at least one router. Empty/omitted = no filter (backward-compatible).
- Internal HTTP routes under `/api/internal/v1/routegroups` (API-Key middleware): `GET /?ids=a,b,c&projection=...` for bulk MS-to-MS lookup.
- In-process: the `events` module calls `RouteGroupServicePort.GetRouteGroupEntityById` during route dispatch.

## Driven Ports (what this module requires)

- `repositories.RouteGroupRepository` → Mongo adapter (`infrastructure/persistence/mongo`, collection `routegroups`).
- `repositories.CacheRepository` → Redis (`common.Cache` + `common.CacheGetOrSetEx`) for per-entity cache-aside.
- `common.AppCache` → service-private Redis used for the org-scoped counter cache.
- `*bootstrap.RouterMetrics` (Prometheus histograms/counters for operations, cache hit/miss, list result size).
- Middlewares: `contextInjector`, `auth`, `apiKey`, `coverage.InjectRequestContext`, `permission.RequirePermission`, `requestValidation`.

## Invariants and Business Rules

- Multi-tenant scoping is mandatory on writes:
  - `isSystem=true` → `orgId` and `pathKey` forced to nil.
  - `isTemplate=true` → only Vendor/Customer `pathKey`s are allowed (`orgfilter.ValidateTemplateCreation`).
  - otherwise → `orgId` + `pathKey` are copied from the `RequestContext` (`orgfilter.ValidateOrgContextForNonSystem`).
- `GET /:id` uses cache-aside with TTL `RouteGroupCacheTTL = 60m` under key `ROUTE_GROUP:{id}`; cache hit/miss is metered.
- `POST` and `DELETE` invalidate the org-scoped counter cache key `counter:route_groups:{orgId}` (TTL 6h).
- `List` enforces hierarchical org visibility via `orgfilter.BuildOrgFilter` + optional system/template ancestors; a missing or deliberate `isSystem=false` / `isTemplate=false` filter excludes those branches.
- `GetRouteGroupsByIds` silently skips not-found IDs — partial results are valid for MS-to-MS callers.
- Services always convert entity ↔ DTO at the boundary via `mapper.EntityToDto` / `mapper.DtoToEntity`; entities never leak through HTTP. The in-process `GetRouteGroupEntityById` is the single exception for sibling-module performance.

## Known Cross-Context Interactions

- Consumed in-process by the sibling **events** module via `RouteGroupServicePort.GetRouteGroupEntityById` to resolve router definitions during dispatch (cache-aside benefits both).
- Exposes bulk retrieval to other microservices through `/api/internal/v1/routegroups` (API-Key authenticated) — callers are expected to be MS peers (inferred; exact consumers not enumerated here).
- Relies on shared coverage/permission/auth middlewares from `packages/microservices` and on `orgfilter` for the tenant hierarchy rules that also govern assets and other tenant-scoped aggregates.
- Does not emit cross-service contracts — all payloads are local DTOs; no reciprocity with `workspace_js/packages/schemas/` is required from this module today (inferred).
