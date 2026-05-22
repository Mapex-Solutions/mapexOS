# Bounded Context: Trigger Execution Events

**Service:** triggers
**Module path:** `src/modules/events/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

Owns the runtime execution path for triggers and workflow plugin actions inside the `triggers` service. Consumes execution requests from NATS, fetches the trigger configuration (via the `triggers` module port), resolves `{{payload.x}}` placeholders, selects the correct executor adapter by type, performs the outbound side effect (HTTP, MQTT, NATS, RabbitMQ, WebSocket, Email, Slack, Teams), and publishes an audit event to the `events` service. Also handles fully resolved workflow plugin pipelines (hooks + operation) and publishes a resume callback back to the `workflow` service.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Executor | Infrastructure adapter that performs a single protocol call (HTTP/MQTT/etc.) for one trigger type | Workflow executor (runtime) in `workflow` service |
| ExecutorRegistry | Factory mapping `triggerType` to its `TriggerExecutor` adapter | Router's matching registry |
| TriggerExecuteEvent | Payload from Router asking this module to fire a stored trigger by `triggerId` | Generic NATS event envelope |
| PluginExecutionEvent | Payload from Workflow carrying a fully resolved action pipeline (no DB fetch needed) | Stored `Trigger` entity |
| Resume callback | Reply message published to the workflow-supplied `callbackSubject` to unblock a paused workflow node | Router-side ACK |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| TriggerExecuted (audit) | `mapexos.events.trigger` | `contracts/services/events/events/TriggerEventDTO` | events service (ClickHouse sink) |
| WorkflowResume | dynamic `callbackSubject` from request | inline `{instanceId, nodeId, status, output?, error?, executionToken?}` | workflow service (runtime module) |
| Outbound side-effects | external (HTTP/MQTT/NATS/SMTP/webhook) | protocol-specific, driven by `TriggerConfig` | external systems |

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| TriggerExecute | `mapexos.trigger.router.execute` (stream `TRIGGERS`) | `contracts/services/triggers/triggers/TriggerExecuteEvent` | router service |
| WorkflowExecute (trigger or plugin mode) | `mapexos.trigger.workflow.execute` (stream `TRIGGERS`) | inline `{mode, orgId, pathKey, workflowId, instanceId, nodeId, callbackSubject, executionToken, data}` | workflow service |

## Driving Ports (what can call this module)

- NATS durable consumer on `mapexos.trigger.router.execute` (`ProcessTriggerExecutionBatch`, V2 batch with retry + DLQ).
- NATS durable consumer on `mapexos.trigger.workflow.execute` (`ProcessWorkflowExecutionBatch`, routes by `mode=trigger|plugin`).
- Legacy single-message entry point `ProcessTriggerExecution` (V1, kept for compatibility).

## Driven Ports (what this module requires)

- `ports.ExecutorRegistry` — resolves `triggerType` → `TriggerExecutor` (HTTP, MQTT, NATS, RabbitMQ, WebSocket, Email, Slack, Teams).
- `triggers.ports.TriggerServicePort` — cache-aside fetch of stored `Trigger` config by ID.
- `natsModel.CorePublisher` — fire-and-forget publish of audit events and resume callbacks, with per-batch `FlushConnection`.
- `bootstrap.TriggerMetrics` — Prometheus counters/histograms for batch size, cache hit/miss, executor duration, DLQ outcomes.

## Invariants and Business Rules

- A disabled trigger (`Enabled=false` or nil) MUST be ACKed without side effects; never retried.
- Invalid JSON, missing executor, failed placeholder resolution, or failed `TriggerConfigToMap` are permanent errors → `msg.Reject(reason)` → straight to DLQ.
- Executor runtime errors (HTTP 5xx, network, SMTP) are transient → `msg.Nack(err)` with retry/backoff until `maxAttempts`, then DLQ.
- Batch processing runs in three phases: parallel execute (bounded by `trigger_executor_workers`, default 50) → single `FlushConnection` → sequential ACK/Nack/Reject.
- Audit event `msgId` is `{eventTrackerId}-triggerlog` to enable JetStream deduplication end-to-end.
- Workflow `plugin` mode MUST publish a resume to `callbackSubject`; missing subject is silently skipped (fire-and-forget workflows).

## Known Cross-Context Interactions

- Pulls trigger definitions from the sibling **triggers** module via its application port (cache-aside over Redis + MongoDB).
- Consumes execution requests produced by the **router** service (match → execute fan-out).
- Publishes audit rows consumed by the **events** service for ClickHouse persistence.
- Replies to the **workflow** service runtime for paused plugin nodes via dynamic `callbackSubject`.
- Outbound executors talk to arbitrary external systems — this is the ONLY module in the triggers service that performs external I/O as a side effect.
