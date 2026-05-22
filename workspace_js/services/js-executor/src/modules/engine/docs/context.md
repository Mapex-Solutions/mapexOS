# Bounded Context: Engine

**Service:** js-executor
**Module path:** `src/modules/engine/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Pure script-execution engine for IoT event payloads. Owns a Piscina worker-thread pool; each worker holds its own V8 Isolate (via `isolated-vm`) and an in-memory compiled-script cache keyed by template. Runs the fixed three-phase pipeline decode → validation → transform, returning the final standardized payload or a structured error. Has no knowledge of NATS, HTTP, assets, or publishing — only receives scripts + raw payload and returns the result.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
| --- | --- | --- |
| ScriptSet | Triple of decode/validation/transform source strings to run per event | `AssetScripts` in scripts module (same shape, domain-specific) |
| Pipeline | The ordered decode → validation → transform execution over one event | Router pipeline (downstream) |
| Isolate | A V8 Isolate owned by a single Piscina worker thread (`isolated-vm`) | Node worker thread itself — one-to-one per worker |
| Context | `ivm.Context` created per event inside an Isolate; holds `payload` global | V8 realm / JS execution context in general |
| Bytecode cache | L0/L1/L2 store of compiled V8 cachedData keyed by `{orgId}/{templateId}/{scriptType}` | Per-worker `scriptCache` Map (in-memory only) |
| Template ID | Key for per-worker compiled-script reuse across events sharing the same template | Asset ID (different identifier) |
| OOM | V8 isolate disposed because it exceeded `memoryLimitMb`; raised as `OOMError` for NACK/retry | Node-level OOM (process crash) |
| Context recycle | Periodic isolate disposal every `contextRecycleInterval` events to bound memory | Context release (per-event cleanup) |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
| --- | --- | --- | --- |
| (none) | — | — | — |

The engine is purely synchronous/RPC-style — it never publishes to NATS. Results are returned in-process to the `scripts` module.

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
| --- | --- | --- | --- |
| (none) | — | — | — |

The engine has no NATS subscriptions. It is invoked in-process via `ScriptEngineServicePort`.

## Driving Ports (inbound — who calls this module)
- `ScriptEngineServicePort` (`application/ports/script_engine_service_port.ts`)
  - `initialize()` / `shutdown()` — lifecycle of the Piscina pool, called from `module.ts` `initListeners`.
  - `runScriptPipeline(rawPayload, userScripts, cacheContext?)` — single-event path; used by `scripts` module for HTTP test endpoints and sample_payload generation.
  - `runBatch(events)` — N-parallel dispatch to the pool; used by the `scripts` batch handler for MQTT/HTTP NATS consumers.
  - `getPoolStats()` — Piscina stats for Prometheus / health endpoints.

## Driven Ports (outbound — what this module requires)
- `BytecodeCachePort` (`application/ports/bytecode_cache_port.ts`) — get/set/invalidate compiled V8 bytecode across L0 (RAM), L1 (disk), L2 (MinIO). Implemented by `infrastructure/cache/TieredBytecodeCache` backed by `@mapexos/infrastructure` `TieredCacheClient` + `MinIOClient`. *(Port is declared and adapter exists, but `ScriptEngineService` currently keeps the reference only for DI wiring — the Piscina worker itself uses a per-worker `scriptCache` Map and does not consult the port directly.)* (inferred)
- `Logger` (`@mapexos/microservices`) — structured logging in the application service only; workers log to stderr via exceptions.
- `Piscina` runtime (`piscina`) — managed thread pool created from `PiscinaOptions { workers, workerPath }` and `PiscinaWorkerConfig { memoryLimitMb, timeoutMs, contextRecycleInterval, mapexValidatorCode }`.
- Prometheus metrics (optional `ScriptEngineMetrics`, `PiscinaPoolMetrics`) — injected from the bootstrap `JsExecutorMetrics`.

## Invariants and Business Rules
- `initialize()` is idempotent; `runScriptPipeline` / `runBatch` auto-initialize on first call.
- Each worker uses `minThreads = maxThreads = workers` (stable pool, no autoscale). `workers` resolves from `CPU_LIMIT - 1` via `resolvePiscinaWorkers`.
- Every user script is wrapped in an IIFE and MUST assign `result`; missing `result` is surfaced as the user-friendly message "The script must define a variable called 'result' with the return value".
- Scripts are compiled once per `(worker, templateId)` pair and cached in the worker's in-memory `scriptCache` Map; cache is cleared on OOM or recycle.
- Validation scripts are preceded by running the static `mapexValidatorCode` setup script in the same context (injects `MapexValidator`).
- If the Isolate is disposed mid-run, the worker returns `isOOM: true` and the application service converts it to `OOMError` so upstream NACKs for retry; plain script failures return `success: false` with `failedAt` set to the failing phase.
- Per-event: a fresh `ivm.Context` is created, `payload` is injected, and the context is released in `finally`; the Isolate itself is recycled every `contextRecycleInterval` events.
- Script timeouts are enforced at `timeoutMs` per `script.runSync` call.
- The engine never mutates the input payload; each phase's output becomes the next phase's `payload` global.

## Known Cross-Context Interactions
- Scripts module (`src/modules/scripts`): sole in-process consumer. Calls `runScriptPipeline` (HTTP/test paths) and `runBatch` (NATS batch paths). Provides the `ScriptSet` and `BytecodeCacheContext`.
- Shared constants (`@shared/constants`): supplies `resolvePiscinaWorkers`, metrics token, and tuning knobs for `memoryLimitMb`, `timeoutMs`, `contextRecycleInterval`.
- `@mapexos/infrastructure`: `TieredCacheClient` ("BytecodeCache") and `MinIOClient` ("MinIOBytecodeClient") back the bytecode port.
- `@mapexos/schemas`: supplies the `StandardizedPayload` shape that the final transform output is typed as (validated by the scripts module, not here).
- Worker file is self-contained: imports only `isolated-vm` and `worker_threads` and is resolved via `path.resolve(__dirname, 'infrastructure/worker/piscina-worker.{ts,js}')` — no DI, no path aliases cross the thread boundary.
