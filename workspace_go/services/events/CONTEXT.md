# Events Service

Stores and serves event data for MapexOS. Consumes multiple NATS JetStream streams in configurable batches, persists them into ClickHouse via bulk inserts, and exposes query APIs for analytics and UI.

## Service Info

| Field | Value |
|---|---|
| Port | `5004` |
| DB | `mapexos` (ClickHouse) + `events` (MongoDB, for retention policies) |
| Config | `src/shared/configuration/application/config.go` |
| Module Config | `src/shared/configuration/modules/config.go` |
| Entry Point | `src/main.go` |
| Build | `go build -o bin/events src/main.go` |

## Modules

| Module | Path | Responsibility |
|---|---|---|
| retention | `src/modules/retention/` | Manages per-organization retention policies (CRUD + defaults). Stores policies in MongoDB, caches TTL lookups in Redis. Listens to `org_created` NATS events to auto-create 8 default retention policies for new orgs. Exposes `/api/v1/retention` HTTP routes. **Must init before `events`** because EventService depends on RetentionServicePort for TTL resolution. |
| events | `src/modules/events/` | Consumes 7 NATS streams (raw, save/store, router, businessrule, trigger, jsexec, DLQ) via batch fetch, parses/validates messages in parallel (bounded worker pool), and bulk-inserts into ClickHouse. Resolves EVA (Entity-Value-Attribute) dynamic field mapping for processed events using a tiered template cache (L0 RAM, L1 disk, L2 S3). Exposes `/api/v1/events` HTTP routes with cursor-based pagination for all event types. |
| app | `src/modules/app/` | Orchestrates module initialization in 3 phases: (1) Repositories, (2) Services, (3) Interfaces (HTTP routes + NATS consumers). Reads module config to determine init order. |

## Key Decisions

- **ClickHouse as primary event store** -- all event types (raw, processed, router, businessrule, trigger, jsexec, DLQ) stored in separate ClickHouse tables optimized for time-series analytics
- **NATS batch fetch + ClickHouse bulk insert** -- batch size configurable via `nats_batch_size` (default 10K), matching ClickHouse's recommended 10K-100K rows per insert for optimal throughput
- **Parallel parse with bounded worker pool** -- `processBatchParallel[T]` uses `NumCPU()*2` goroutines, each writing to its own result slot (no contention), then a single bulk insert for all valid entities
- **EVA (Entity-Value-Attribute) dynamic fields** -- processed events resolve field mappings from asset templates via a tiered cache (L0 RAM/L1 disk/L2 S3), mapping dynamic data into typed ClickHouse MAP columns (`eva_number`, `eva_string`, `eva_bool`, `eva_date`)
- **Retention MUST init before events** -- EventService depends on RetentionServicePort to resolve TTL per organization per table
- **DLQ consumer never retries** -- ACKs all messages even on failure to prevent infinite redelivery loops
- **Graceful shutdown** -- P0 drains HTTP, P5 closes ClickHouse/MongoDB/Redis/NATS concurrently
- **DI via DIG container** -- services never import ports directly; wired through DI container following hexagonal architecture

## Docs

Full documentation at `docs/`. See `docs/index.md` for the complete map.
