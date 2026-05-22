# Bounded Context: Triggers (CRUD + Catalog)

**Service:** triggers
**Module path:** `src/modules/triggers/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose

Owns the lifecycle of `Trigger` entities — the reusable, stored definitions that describe an outbound side effect (HTTP call, MQTT/NATS/RabbitMQ publish, WebSocket push, Email/Slack/Teams message). Exposes the authoritative REST API for creating, listing, reading, updating, and deleting triggers, plus a cache-aside read path consumed at runtime by the sibling `events` module. Enforces multi-tenant isolation (org + pathKey) and the Template Resources pattern (system / template / local visibility).

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Trigger | Stored definition of an outbound action: type + category + resolved config schema | Router "match rule" that fires it |
| TriggerConfig | Discriminated union of protocol configs (`Http`, `Mqtt`, `Nats`, `Rabbitmq`, `Websocket`, `Email`, `Teams`, `Slack`); exactly one is populated per entity | Workflow plugin action (transient, not stored here) |
| Category | High-level grouping — `technical` vs `communication` | Trigger type (finer: `http`, `slack`, etc.) |
| isSystem / isTemplate | Template Resources flags: system = global MAPEX library; template = vendor/customer inheritable; both false = org-local | `enabled` (runtime on/off switch) |
| Counter | Cached per-org count of triggers, served by the `/counter` endpoint | Prometheus metric |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| (none) | — | — | — |

This module is CRUD-only; it does not publish domain events. Runtime fan-out is the responsibility of the `events` module + router service.

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| (none) | — | — | — |

No NATS consumers. All inbound traffic is HTTP.

## Driving Ports (what can call this module)

- REST API under `/api/v1/triggers` (Fiber router group with `ContextInjector` + `AuthMiddleware`):
  - `GET /` — paginated, filtered list (coverage-aware org filtering).
  - `GET /counter` — cached count per org.
  - `POST /` — create.
  - `GET /:id` — read by id.
  - `PATCH /:id` — partial update.
  - `DELETE /:id` — delete.
- In-process port `TriggerServicePort.GetTriggerById(ctx, id, *CacheMetrics)` used by the `events` module at execution time.

## Driven Ports (what this module requires)

- `domain.repositories.TriggerRepository` — MongoDB collection access (create/find/update/delete/count/paginate).
- `domain.repositories.CacheRepository` — Redis-backed `common.Cache` + `CacheGetOrSetEx` for the `TRIGGER:{id}` cache-aside read path.
- `common.AppCache` — service-private cache used for the per-org trigger counter (`counter:triggers:{orgId}`, TTL 6h).
- HTTP middlewares from `packages/microservices`: `requestValidation`, `permissionMw.RequirePermission`, `coverageMw.InjectRequestContext`.

## Invariants and Business Rules

- Every non-system trigger MUST carry both `orgId` and `pathKey`; `isSystem=true` triggers MUST have neither (global catalog).
- `TriggerType` MUST match exactly one populated branch of `TriggerConfig`; validation rejects mismatches at the DTO layer (Zod contract mirror).
- Reads go through `TRIGGER:{triggerId}` cache with 60-minute TTL; writes (`Update`, `Delete`) MUST invalidate the cache entry.
- List queries are coverage-filtered by `RequestContext` (hierarchical pathKey expansion from the coverage cache) — users only see triggers their org tree can access.
- Counter cache is invalidated on every create/delete for the affected `orgId`.
- Permission gates: `TriggerList`, `TriggerRead`, `TriggerCreate`, `TriggerUpdate`, `TriggerDelete` from `packages/permissions/triggers`.

## Known Cross-Context Interactions

- Read by the sibling **events** module at execution time via `TriggerServicePort.GetTriggerById` (cache-aside, reports `CacheMetrics` for Prometheus).
- DTOs are aliased from `packages/contracts/services/triggers/triggers` (single source of truth, mirrored in `workspace_js/packages/schemas`).
- Consumed by the frontend (`workspace_js/apps/mapexOS`) via the generated API wrapper in `@mapexos/apis` for the Triggers management UI.
- Indirectly referenced by the **router** service: router match rules point at `triggerId`s managed here.
