# Bounded Context: Engine

**Service:** js-workflow-executor
**Module path:** `src/modules/engine/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose

Owns the V8 sandboxed execution substrate for workflow code nodes. Wraps a Piscina worker-thread pool where each worker hosts its own V8 Isolate with per-worker script/bytecode caching, and exposes a single `runWorkflowScript` dispatch entry point to the rest of the service. Also owns the `BytecodeCachePort` contract used to persist compiled V8 bytecode across cold starts. Pure execution substrate — no knowledge of NATS, MinIO wiring details, or workflow business semantics.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Isolate | A V8 Isolate owned by a single Piscina worker thread | The `NodeVM`/`vm2` sandboxes used by js-executor (different service, different engine) |
| Worker | A Piscina-managed Node worker thread hosting one Isolate | A workflow "worker/executor" in the Go workflow service |
| Cache key | `{orgId}/{workflowId}/scripts/{nodeId}` for source, `{orgId}/{definitionId}/bytecode/{nodeId}` for bytecode | Generic TieredCache keys |
| Bytecode | Cached V8 compiled bytecode for a code node — definition-scoped, reused across instances | Script source string |
| OOM | Transient Isolate out-of-memory — isolate is disposed and the event NACK'd for retry | A permanent script error |
| Context recycle | Periodic Isolate/context reset every N events to prevent slow leaks | Full worker replacement |

## Published Events (outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — |

Engine publishes no NATS events directly; results flow back in-process as `PiscinaWorkerOutput` and are forwarded by the `scripts` module callback publisher.

## Consumed Events (inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | — |

Engine consumes no NATS events directly; it is invoked in-process by `WorkflowScriptService`.

## Driving Ports (inbound — who calls this module)

- In-process `ScriptEngineServicePort` consumed by `@modules/scripts` `WorkflowScriptService.execute` for every code node execution.
- `initialize` / `shutdown` lifecycle called by `engine/module.ts` during bootstrap Phase 4 and graceful shutdown.
- `getPoolStats` read by the metrics/health bootstrap (via `WORKFLOW_EXECUTOR_METRICS`) to export Piscina gauges.

## Driven Ports (outbound — what this module requires)

- `Piscina` worker pool (infrastructure, owned by this module) — thread pool + structured-clone transport.
- `BytecodeCachePort` implementations (see `scripts` module DI) — TieredCache (L0 RAM + L1 Redis) over MinIO (L2).
- `Logger` (`@mapexos/microservices`) for `[SERVICE:ScriptEngine]` logs.
- Optional `ScriptEngineMetrics` + `PiscinaPoolMetrics` Prom counters/histograms/gauges.
- `ConfigModule` values: `isolate_memory_limit_mb`, `worker_script_timeout_ms`, `context_recycle_interval`, worker count (via `resolvePiscinaWorkers`).

## Invariants and Business Rules

- Each worker owns exactly one V8 Isolate; `minThreads = maxThreads` for a stable pool size.
- Structured-clone boundary: `PiscinaWorkerInput` / `PiscinaWorkerOutput` MUST be plain-serializable (no functions, no class instances, no circular refs).
- `runWorkflowScript` auto-initializes if not yet initialized — lazy-safe.
- V8 OOM MUST surface as `OOMError` and trigger consumer NACK for retry; all other script errors are returned as `{ success: false, error }` and ACK'd.
- Bytecode is definition-scoped (`{orgId}/{definitionId}/bytecode/{nodeId}`); invalidation clears only L0+L1 — L2 (MinIO) is managed by the Go workflow service.
- Per-script `timeoutMs` from node config overrides the worker default; absent value falls back to `worker_script_timeout_ms`.
- Context recycle every `contextRecycleInterval` events mitigates long-running Isolate leaks.

## Known Cross-Context Interactions

- Invoked by `@modules/scripts` (`WorkflowScriptService`) — the only in-process caller.
- Reads `BytecodeCachePort` (`TieredBytecodeCache`) backed by Redis + MinIO; L2 writes/deletes on MinIO are coordinated with the Go `workflow` service, which owns the source-of-truth bytecode blobs.
- Consumes `ConfigModule` and metrics from the service bootstrap; registered under `SCRIPT_ENGINE_SERVICE_TOKEN`.
- Distinct from the `js-executor` service engine: that one runs V8 for IoT event scripts; this one executes workflow DAG code nodes.
