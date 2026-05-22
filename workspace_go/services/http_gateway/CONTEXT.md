# HTTP Gateway Service

High-performance HTTP gateway for MapexOS event ingestion — receives webhooks from external systems, authenticates per Data Source config, and publishes normalized events to NATS JetStream.

## Service Info

| Field | Value |
|---|---|
| Port | `5001` |
| DB | `dev-http_gateway` |
| Config | `src/shared/configuration/application/config.go` |
| Module Config | `src/shared/configuration/modules/config.go` |
| Entry Point | `src/main.go` |
| Build | `go run http_gateway/src/main.go` |

## Modules

| Module | Path | Responsibility |
|---|---|---|
| app | `src/modules/app/` | Module orchestrator — initializes all business modules in 3 phases (repositories, services, interfaces) following the order defined in `modules/config.go` |
| datasources | `src/modules/datasources/` | CRUD management of Data Source entities (name, auth, rate limit, asset binding, working hours). Exposes REST API at `/api/v1/data_sources`. Uses Redis cache-aside pattern for fast lookups. Auth strategies per Data Source: `apiKey`, `jwt`, `oauth2`, `ip_whitelist`, `none`. Multi-tenant with hierarchical org filtering via PathKey |
| events | `src/modules/events/` | Webhook event ingestion at `POST /api/v1/events?ds={id}`. Applies per-Data-Source authentication via custom middleware (delegates to datasources module for config lookup). On success, publishes event to `processor.js.execute` NATS subject for JS-Executor pipeline. On auth failure, publishes security event to `events.raw` for monitoring. No own repository — stateless processing only |

## Key Decisions

- **Cache-aside pattern** for Data Source lookups — Redis cache with TTL to minimize MongoDB hits on the hot event ingestion path
- **Per-Data-Source authentication** — each Data Source defines its own auth strategy (JWT, OAuth2, API Key, IP whitelist, none), enforced by a custom middleware on the events endpoint
- **Fire-and-forget auth failure events** — authentication failures are published asynchronously to `events.raw` NATS subject for security monitoring without blocking the response
- **Minimal payload to JS-Executor** — only `orgId` and `assetBind` are sent; downstream services fetch metadata from Asset cache to avoid stale data
- **Hexagonal Architecture** — services depend on port interfaces, DI via dig container, strict separation between domain/application/infrastructure/interfaces layers
- **Middleware ordering** — Validation runs before Permission (fail-fast pattern) on datasources routes; events route uses Validation then CustomAuth then Handler
- **Graceful shutdown** — uses `shutdown.Manager` with 15s timeout for clean teardown of Fiber, MongoDB, Redis, and NATS connections

## Docs

Full documentation at `docs/`. See `docs/index.md` for the complete map.
