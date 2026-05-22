# Bounded Context: Events

**Service:** js-workflow-executor
**Module path:** `src/modules/events/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

Owns the inbound NATS boundary of the service. Wires two JetStream consumers that drive the workflow code execution pipeline: a load-balanced queue consumer on `WORKFLOW-JS-CODE` for per-node execution requests, and a broadcast FANOUT consumer on `mapexos.fanout.workflow.definition.invalidate` for definition cache invalidation. It decodes NATS envelopes, sets multi-tenant DLQ context, and delegates to `WorkflowScriptServicePort`. Contains no business logic.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Queue consumer | Load-balanced durable: exactly one pod processes each message | FANOUT subscriber (broadcast) |
| FANOUT consumer | Broadcast subscriber: every pod receives every message | Queue-group consumer |
| Invalidation | L0+L1 cache drop for a definition or its code nodes | L2 (MinIO) deletion — owned by the Go workflow service |
| Callback | Execution result published to `WORKFLOW-RESUME` (performed by `scripts` module, not here) | NATS ack/nack |
| Granular invalidation | Payload carries `nodeIds[]` — only those nodes' caches are cleared | Workflow-level fallback (TTL expiry) |

## Published Events (outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — |

Events module only consumes; downstream callbacks are published by `scripts/infrastructure/adapters/nats_callback_publisher.ts` (see `scripts` context).

## Consumed Events (inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| WorkflowCodeRequest | `mapexos.workflow.js.code` (stream `WORKFLOW-JS-CODE`, durable `js-workflow-executor-code`) | `WorkflowScriptInput` (scripts/application/ports) (inferred — not yet in `@mapexos/schemas`) | `workflow` Go service (runtime DAG dispatcher) |
| DefinitionInvalidate | `mapexos.fanout.workflow.definition.invalidate` (stream `FANOUT`) | `DefinitionInvalidatePayload` `{ orgId, definitionId, nodeIds? }` (inferred — not yet in `@mapexos/schemas`) | `workflow` Go service on create/update/delete of a definition |

## Driving Ports (inbound — who calls this module)

- NATS `WORKFLOW-JS-CODE` queue consumer (batch handler via `NatsBus.startConsumer` with `batchMessageHandlerV2`).
- NATS FANOUT subscriber on `mapexos.fanout.workflow.definition.invalidate` via `NatsBus.subscribeFanout`.
- `initListeners()` is invoked by the service bootstrap in Phase 4.

## Driven Ports (outbound — what this module requires)

- `WorkflowScriptServicePort` (`@modules/scripts/application/ports`) — `execute`, `invalidateNodes`, `invalidateWorkflow`.
- `NatsBus` from `@mapexos/infrastructure` — `startConsumer`, `subscribeFanout`, `ensureFanoutStream`.
- `Logger` + `ConfigModule` from `@mapexos/microservices`.
- Optional Prom metrics injected from bootstrap: `executionDuration`, `executionsTotal`, `batchSize`.
- Shared constants: `SERVICE_NAME`, `SERVICE_TYPE`, `DEFAULT_RETRY_POLICY`, `resolveConsumerConfig` (`@shared/constants`).

## Invariants and Business Rules

- Queue consumer MUST set `msg.orgId` and `msg.pathKey` from payload before handling — required by DLQ for multi-tenant filtering.
- Messages missing `orgId`/`workflowId`/`nodeId`/`instanceId` are ACK'd with a warning (poison-pill protection).
- Script-execution errors ACK the message (callback already published with `status:'error'`); only `OOMError` triggers NACK for retry with backoff.
- FANOUT stream MUST be ensured before subscribing (`maxAge 5m`, `maxMsgs 10000`, subjects `mapexos.fanout.>`).
- FANOUT handler is best-effort: failures are logged and swallowed — no NACK semantics on fanout broadcast.
- Granular invalidation path (with `nodeIds`) is preferred; empty/missing `nodeIds` falls back to workflow-wide TTL expiry.
- Consumer files contain wiring only — no business logic, no type/const declarations leak outside the designated `interfaces/`/`constants/` files.

## Known Cross-Context Interactions

- Upstream publisher: Go `workflow` service emits both the `WORKFLOW-JS-CODE` execution requests and the `FANOUT` invalidation broadcasts.
- Downstream: invokes `@modules/scripts` `WorkflowScriptService` which in turn publishes execution callbacks to the subject provided in `input.callbackSubject` (typically `WORKFLOW-RESUME`).
- Cross-service contract reciprocity for `WorkflowScriptInput` and `DefinitionInvalidatePayload` is expected in `workspace_js/packages/schemas/` and `workspace_go/packages/contracts/` (inferred — currently defined locally).
