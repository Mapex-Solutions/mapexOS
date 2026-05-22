# Triggers Service

Outbound execution layer of MapexOS — consumes trigger execution events from NATS, resolves dynamic placeholders, and dispatches actions through a registry of executors (HTTP, MQTT, RabbitMQ, NATS, WebSocket, Email, Teams, Slack).

## Service Info

| Field | Value |
|---|---|
| Port | `5006` |
| DB | `dev-triggers` |
| Config | `src/shared/configuration/application/configMap.go` |
| Module Config | `src/shared/configuration/modules/config.go` |
| Entry Point | `src/main.go` |
| Build | `go build -o bin/triggers src/main.go` |

## Modules

| Module | Path | Responsibility |
|---|---|---|
| app | `src/modules/app/` | Orchestrates module initialization in 3 phases (repositories, services, interfaces) following the order defined in module config |
| triggers | `src/modules/triggers/` | CRUD API for trigger configurations. Provides REST endpoints at `/api/v1/triggers` with auth middleware. Persists trigger definitions to MongoDB (`triggers` collection) and exposes a cache handler for Redis-backed low-latency reads |
| events | `src/modules/events/` | Execution pipeline. Consumes `trigger.*.execute` from NATS JetStream (stream `TRIGGERS`), resolves placeholders in trigger configs using event payload data, dispatches to the correct executor via `ExecutorRegistry`, and publishes results to `events.trigger`. Executors: HTTP, MQTT, RabbitMQ, NATS, WebSocket, Email, Teams, Slack |

## Key Decisions

- **Executor Registry pattern** — all executors implement a `TriggerExecutor` port interface and are registered by type string in a factory registry (`infrastructure/registry/`). The `EventService` receives the registry via DI and selects the executor at runtime by `triggerType`.
- **Hexagonal Architecture** — strict separation: `domain/` (entities + repository interfaces), `application/` (services + ports + DTOs), `infrastructure/` (persistence + executors), `interfaces/` (HTTP handlers + NATS consumers).
- **DI via DIG container** — services never import ports directly; everything is wired through the DIG container with constructor injection.
- **Redis cache for trigger definitions** — trigger configs are cached in Redis for low-latency execution; MongoDB is the source of truth.
- **NATS JetStream consumer** — trigger execution events are consumed from the `TRIGGERS` stream on subject `trigger.*.execute` with configurable batch size and fetch timeout.
- **Graceful shutdown** — uses a shutdown manager with 15-second timeout, registered in bootstrap.

## Docs

Full documentation at `docs/`. See `docs/index.md` for the complete map.
