# Bounded Context: Runtime

**Service:** workflow
**Module path:** `src/modules/runtime/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

The DAG execution engine. Owns the `WorkflowExecution` root aggregate (hot state in NATS KV `exec.{uuid}`), the DAG walker, the executor registry (inline + control + async), and the dispatch fabric to external services. Consumes three NATS streams (`WORKFLOW-EXECUTION`, `WORKFLOW-RESUME`, `WORKFLOW-SCHEDULE`) and produces lifecycle events on `WORKFLOW-STATE` plus execution dispatches to `WORKFLOW-JS-CODE`, `WORKFLOW-SCHEDULE`, and the Triggers Service. Handles suspension/resume for async nodes, timer-based scheduling via NATS JetStream scheduled messages, retry-with-backoff, subworkflow parent↔child callbacks, and fanout concurrency via CAS on KV. Never writes to MongoDB — the Archiver owns persistence.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| WorkflowExecution | Root aggregate: one run of an instance; hot JSON in NATS KV, upserted to Mongo only on terminal by Archiver | WorkflowInstance (config) or WorkflowDefinition (template) |
| Walker | In-process DAG traversal loop that advances nodes until a suspension point or terminal state | Consumer (which just dispatches into `HandleExecution`/`HandleResume`) |
| Suspension | Node transition to `waiting` state with an external dispatch (code/subworkflow/trigger/plugin) awaiting callback | Completion (walker continues normally) |
| Resume | Inbound message restarting a suspended node — timer/signal/callback/re-enqueue variants | Retry (handled internally via retry timer + CAS) |
| NATS Schedule | JetStream-scheduled future message on `WORKFLOW-SCHEDULE` used for timer-based resumes (delay/retry/timeout) | Internal Go timer (not used — survives pod restart) |
| Execution token + Msg-Id | Per-dispatch idempotency pair: token validates callback authenticity, Msg-Id dedups external redelivery | Workflow UUID (the execution's identity) |
| CAS | Compare-And-Swap on KV used for concurrent fanout callback merges (up to `MaxCASRetries=5` before Nack) | Transaction (there is none — NATS-native CAS only) |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| StateEvent (created/waiting/resumed/completed/failed/cancelled) | `mapexos.workflow.state.{status}` (stream `WORKFLOW-STATE`) | `shared/types.StateEvent` | archiver |
| ResumeMessage (re-enqueue) | `mapexos.workflow.resume.reenqueue.{instanceId}` (stream `WORKFLOW-RESUME`) | `shared/types.ResumeMessage` | runtime (self) |
| ResumeMessage (signal/timer/callback) | `mapexos.workflow.resume.{signal|timer|callback}.{instanceId}` | `shared/types.ResumeMessage` | runtime (self) |
| CodeExecutionRequest | `mapexos.workflow.js.code` (stream `WORKFLOW-JS-CODE`) | `interfaces/message.CodeExecutionRequest` | js-workflow-executor |
| SubworkflowExecution | `mapexos.workflow.execution.subworkflow.{instanceId}` (stream `WORKFLOW-EXECUTION`, `mode=subworkflow`) | execution command payload | runtime (self) |
| WorkflowTriggerRequest | `mapexos.trigger.workflow.execute` | `interfaces/message.WorkflowTriggerRequest` (mode `trigger` or `plugin`) | triggers service |
| NATS Schedule (timer) | `mapexos.workflow.schedule.{wfUUID}.{nodeID}` (stream `WORKFLOW-SCHEDULE`) | resume payload | runtime's ScheduleFire consumer |

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| ExecutionCommand (newInstance / signal / signalOrStart / subworkflow) | `mapexos.workflow.execution.>` (stream `WORKFLOW-EXECUTION`) | command payload (`mode` in body) | http_gateway, router, runtime (subworkflow self-dispatch) |
| ResumeMessage | `mapexos.workflow.resume.>` (stream `WORKFLOW-RESUME`) | `shared/types.ResumeMessage` | runtime (self), js-workflow-executor (callbacks), triggers service (callbacks) |
| ScheduleFired | `mapexos.workflow.schedule.fired` (stream `WORKFLOW-SCHEDULE`) | scheduled resume body | NATS JetStream (schedule fire) |

## Driving Ports (what can call this module)

- NATS consumers: `WORKFLOW-EXECUTION` (3-mode dispatch), `WORKFLOW-RESUME`, `WORKFLOW-SCHEDULE` (schedule-fire → re-publish to `WORKFLOW-RESUME`).
- Cross-module Go API: `RuntimeServicePort.ExecuteByInstanceID(...)` called in-process by the `instances` HTTP execute endpoint.

## Driven Ports (what this module requires)

- `ExecutionStateRepository` (NATS KV `WORKFLOW-INSTANCES`, key `exec.{uuid}`).
- `RuntimePublisherPort` (NATS) — state events, resume, schedule publish/purge, code/subworkflow/trigger dispatch.
- `DefinitionLoaderPort` (TieredCache-wrapped `DefinitionRepository` from the definitions module).
- `ConditionEvaluatorPort` + `ValueResolverPort` (from engine module).
- `VaultPort` — HTTP to mapexVault for credential decryption at dispatch time.
- `PluginRepo` — manifest lookup for plugin-backed nodes (from plugins module).
- Metrics registry (checkpoint duration, execution counters, dispatch counters).

## Invariants and Business Rules

- Runtime NEVER writes to MongoDB. The Archiver is the only writer.
- Dispatch-before-checkpoint: `suspendExecution` MUST dispatch the external call first; if dispatch fails, no checkpoint is written and NATS redelivers. Msg-Id dedups on redelivery.
- Terminal transitions (`failExecution` / `completeExecution`) MUST call `PurgeAllSchedules(wfUUID)` before completing to cancel pending timers.
- Subworkflow child completion publishes `PublishCallbackResume` to the parent's `CallbackSubject` with `ParentExecutionToken` for callback authenticity.
- Fanout mode `firstCompleted`: upon first terminal branch, remaining waiting nodes are cancelled (path status flipped to `cancelled`) and schedules purged.
- Concurrent fanout merges use CAS on KV with up to `MaxCASRetries=5` attempts; beyond that the message is Nacked for NATS redelivery.
- Timer-based nodes (`core/delay`, `retry`) DO NOT dispatch externally — the NATS Schedule itself carries the resume payload.
- Retry timer preserves `__retryAttempt` across suspension cycles on the node state.

## Known Cross-Context Interactions

- Publishes `WORKFLOW-STATE` consumed by **archiver** (the only MongoDB writer for executions).
- Dispatches JS code to **js-workflow-executor** on `WORKFLOW-JS-CODE`; callbacks arrive back on `WORKFLOW-RESUME`.
- Dispatches plugin/trigger actions to the **triggers service** on `mapexos.trigger.workflow.execute`; callbacks arrive back on `WORKFLOW-RESUME`.
- Reads definitions via **definitions** module (TieredCache) and plugin manifests via **plugins** module (TieredCache).
- Uses **engine** module for all condition evaluation and field-value resolution.
- Decrypts credentials via the **mapexVault** microservice.
