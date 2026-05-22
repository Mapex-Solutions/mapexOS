# Router Service

Routes incoming asset events to downstream processors based on RouteGroup conditional match rules, with fan-out to multiple NATS destinations.

## Service Info

| Field | Value |
|---|---|
| Port | `5003` |
| DB | `dev-router` |
| Config | `src/shared/configuration/application/config.go` |
| Module Config | `src/shared/configuration/modules/config.go` |
| Entry Point | `src/main.go` |
| Build | `go build -o bin/router src/main.go` |

## Modules

| Module | Path | Responsibility |
|---|---|---|
| app | `src/modules/app/` | Module orchestrator — initializes all modules in 3 phases (repositories, services, interfaces) following the order defined in module config |
| routegroups | `src/modules/routegroups/` | CRUD management of RouteGroup entities via REST API (external JWT-auth + internal API-key MS-to-MS endpoints). Persists to MongoDB, exposes cache handler for RouteGroup lookups used by the events module |
| events | `src/modules/events/` | Core routing pipeline — consumes `route.execute` from NATS JetStream (WorkQueue pattern), resolves asset context via TieredCache (L0 RAM, L1 Disk, L2 MinIO, Fallback HTTP), evaluates match rules using a domain MatchEvaluator (operators: eq/neq/gt/gte/lt/lte/in/nin, policies: all/any), publishes matched events to downstream subjects (RuleEngine, Triggers, Event Store, Lakehouse, Notifications), and emits routing history for audit. Also runs a FANOUT consumer for asset cache invalidation across all instances |

## Key Decisions

- **Three-phase batch processing**: events are processed in parallel (Phase 1: worker pool sized at `NumCPU*2`), then flushed (Phase 2: single `FlushConnection`), then ACK/Nack/Reject sequentially (Phase 3) — avoids goroutine explosion and ensures fire-and-forget publishes are batched
- **TieredCache for asset resolution**: L0 (RAM) -> L1 (Disk) -> L2 (MinIO) -> Fallback (HTTP to Assets service) — avoids direct DB queries on the hot path
- **FANOUT cache invalidation**: separate NATS consumer broadcasts asset invalidation to all Router instances so each local L0+L1 cache stays consistent
- **JetStream dedup via MsgId**: each publish uses `{eventTrackerId}-{routerIdx}` as MsgId to prevent duplicate downstream processing on retries
- **Domain MatchEvaluator is pure logic**: no I/O dependencies, supports dot-notation field paths for nested event data, fully unit-testable
- **DI via DIG container**: all modules register dependencies through a centralized DIG container following Hexagonal Architecture (services depend on port interfaces, not concrete implementations)
- **Reject vs Nack distinction**: validation errors (malformed JSON, missing fields) are Rejected to DLQ; processing failures (cache miss, publish error) are Nacked for retry

## Docs

Full documentation at `docs/`. See `docs/index.md` for the complete map.
