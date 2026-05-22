# Bounded Context: Scripts

**Service:** js-workflow-executor
**Module path:** `src/modules/scripts/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

Orchestrates workflow code node execution end-to-end. For each incoming request, it resolves the script source from TieredCache, attaches cached V8 bytecode when available, dispatches to the engine's Piscina worker, stores newly-produced bytecode for future cold starts, and publishes a `WorkflowScriptCallback` to the `WORKFLOW-RESUME` subject supplied by the runtime. Also owns cache-invalidation orchestration (node-level and workflow-level, L0+L1 only) triggered by FANOUT broadcasts from the Go workflow service.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Code node | A workflow DAG node whose behavior is user JavaScript producing `{ output, statePatch }` | A generic IoT script (js-executor domain) |
| Script source | The raw JS string for a code node, cached at `{orgId}/{workflowId}/scripts/{nodeId}` | V8 bytecode artifact |
| Callback | `WorkflowScriptCallback` published to the runtime-supplied subject to resume the DAG | NATS ack/nack |
| Execution token | Opaque token echoed from input to callback for runtime validation (optional, backward-compatible) | NATS msg-id |
| State patch | Partial state object the script returns; merged into the workflow instance state by the runtime | Full state replacement |
| Invalidate nodes / workflow | L0+L1 cache drop for script source + bytecode; L2 (MinIO) owned by Go service | Hard delete of the definition |

## Published Events (outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| WorkflowScriptCallback | `callbackSubject` from input (runtime-supplied, typically `mapexos.workflow.resume.callback.{instanceId}`) (inferred) | `WorkflowScriptCallback` `{ instanceId, nodeId, executionToken?, status, output?, statePatch?, error? }` | `workflow` Go service (`WORKFLOW-RESUME`) |

Not yet mirrored in `@mapexos/schemas` / `workspace_go/packages/contracts/` (inferred gap).

## Consumed Events (inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | — |

Module has no own NATS consumers; it is driven in-process by the `events` module consumers.

## Driving Ports (inbound — who calls this module)

- `WorkflowScriptServicePort.execute` called by the `WORKFLOW-JS-CODE` queue consumer (events module).
- `WorkflowScriptServicePort.invalidateNodes` / `invalidateWorkflow` called by the `DefinitionInvalidate` FANOUT consumer (events module).
- Registered under `WORKFLOW_SCRIPT_SERVICE_TOKEN` in the tsyringe container.

## Driven Ports (outbound — what this module requires)

- `ScriptSourceCachePort` — `TieredScriptSourceCacheAdapter` over `TieredCacheClient('ScriptSourceCache')` (L0 RAM + L1 Redis + L2 MinIO fallback).
- `BytecodeCachePort` (from engine) — `TieredBytecodeCache` over `TieredCacheClient('BytecodeCache')` + `MinIOClient('MinIOWorkflowsClient')`.
- `ScriptEngineServicePort` (from engine) — `runWorkflowScript` dispatch.
- `CallbackPublisherPort` — `NatsCallbackPublisher` backed by `NatsBus` (`@mapexos/infrastructure`).
- `Logger` from `@mapexos/microservices`.

## Invariants and Business Rules

- Script source and cached bytecode MUST be fetched in parallel (`Promise.all`) to minimize latency.
- Missing script source produces a `SCRIPT_NOT_FOUND` error callback and returns without dispatching to the engine.
- Newly-produced bytecode (`result.newBytecode`) is stored fire-and-forget; failures log a `[SERVICE:WorkflowScript]` warning but do not fail the execution.
- A callback MUST be published for every execute call, success or error; publish failures are logged but do not throw (consumer ACK semantics rely on this).
- Idempotency is external: publish-side Runtime MsgId dedup + consumer-side CAS on NATS KV (not enforced by this module).
- Cache-key format is authoritative here: source `{orgId}/{workflowId}/scripts/{nodeId}`; bytecode via `BytecodeCachePort.buildCacheKey`.
- Invalidation only clears L0+L1; L2 (MinIO) is managed exclusively by the Go workflow service.
- Per-node `timeout` is expressed in seconds in `WorkflowScriptInput` and converted to `timeoutMs` when dispatched to the engine.

## Known Cross-Context Interactions

- Called by `@modules/events` consumers for both execution and invalidation flows.
- Calls `@modules/engine` for V8 dispatch and bytecode cache operations (shares `BytecodeCachePort`).
- Publishes callbacks consumed by the Go `workflow` service (`WORKFLOW-RESUME` stream) to resume the DAG at the code node.
- Shares the L2 bytecode/MinIO bucket lifecycle with the Go workflow service: Go writes/deletes the source of truth, this module only invalidates local tiers on FANOUT.
- Distinct from `js-executor`'s `scripts` module: that one runs IoT-event scripts with a different input/output shape and publishes to `mapexos.js.script.result`.
