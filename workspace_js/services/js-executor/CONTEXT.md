# JS-Executor Service

High-performance Node.js microservice that executes user-defined JavaScript in isolated V8 worker threads. Consumes NATS events with IoT payloads and processes them through a decode → validate → transform pipeline, then publishes processed results downstream.

## Service Info

| Field | Value |
|---|---|
| Port | `8000` |
| Config | `src/shared/configuration/application/configMap.ts` |
| Module Config | `src/shared/configuration/modules/config.ts` |
| Entry Point | `src/main.ts` |
| Build | `npm run build` → `dist/` |
| Runtime | Node.js ≥24, isolated-vm, Piscina |

## Modules

| Module | Path | Responsibility |
|---|---|---|
| engine | `src/modules/engine/` | V8 script execution via Piscina worker pool. Each worker owns a V8 Isolate (32MB limit). Per-worker script compilation cache, context recycling every 10k events, OOM recovery. Two pools: single-event (HTTP) and batch (NATS). |
| scripts | `src/modules/scripts/` | Orchestrates script execution pipeline. Fetches Asset/Template from TieredCache, builds ScriptSet (decode, validation, transform), dispatches to engine, publishes results to NATS. HTTP endpoints for script testing. |
| events | `src/modules/events/` | 4 NATS consumers: JsExecuteConsumer (queue, `js.execute.*`), MqttDataConsumer (queue, `dt.>`), AssetInvalidateConsumer (fanout), TemplateInvalidateConsumer (fanout). |
| app | `src/modules/app/` | Root HTTP routes (minimal). |

## Communication

| Direction | Protocol | Subject/Endpoint | Purpose |
|---|---|---|---|
| IN | NATS Queue | `js.execute.*` | Execute requests from HTTP Gateway |
| IN | NATS Queue | `dt.>` | MQTT device telemetry |
| IN | NATS Fanout | `fanout.asset.invalidate` | Asset cache invalidation |
| IN | NATS Fanout | `fanout.template.invalidate` | Template cache invalidation |
| OUT | NATS | downstream subjects | Processed event payloads |
| IN | HTTP | `POST /api/v1/scripts/test` | Script testing (JWT auth) |
| IN | HTTP | `GET /api/v1/scripts/templates/:id/*` | Template scripts (JWT auth) |
| IN | HTTP | `GET /internal/templates/:orgId/:id/*` | Internal access (API key) |
| OUT | HTTP | Assets Service (`5002`) | Fallback on cache miss |

## Key Decisions

- **Isolated V8 execution**: Each worker thread owns a V8 Isolate (isolated-vm) with 32MB heap limit and 10s timeout for safe user code execution.
- **TieredCache (L0/L1/L2)**: L0=RAM (~50µs), L1=Disk (~500µs), L2=MinIO (source of truth). NATS fanout invalidation. HTTP fallback to Assets Service.
- **Batch processing**: NATS consumers use batch workers with direct NATS publish (fire-and-forget + single flush) for throughput.
- **Auto-tuning**: Single `CPU_LIMIT` knob derives worker count, batch size, and consumer concurrency.
- **Bytecode caching**: Compiled V8 bytecode stored in L1/L2 for cross-pod reuse.
- **DI via tsyringe**: Services depend on port interfaces, concrete implementations injected via container.

## Docs

Full documentation at `docs/`. See `docs/index.md` for the complete map.
