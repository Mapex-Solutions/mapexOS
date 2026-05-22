# Workflow Service Documentation

## Overview

The Workflow service is the **DAG execution engine** of MapexOS. It allows users to design
workflows as directed acyclic graphs (nodes + edges), trigger them via events or API calls,
and execute them with full state management, crash recovery, and async suspension support.

## 3-Entity Architecture

The service uses three distinct entities with clear ownership:

| Entity | What it is | Storage |
|---|---|---|
| **Definition** | Reusable blueprint (DAG, states schema, externalInputs schema) | TieredCache → MongoDB |
| **Instance** | Configuration to run (which definition, filled externalInputs) | TieredCache → MongoDB |
| **Execution** | Mutable execution record (status, state, DAG progress) | NATS KV + MongoDB |

```
Definition 1 ──→ N Instance 1 ──→ N Execution
```

The service follows a two-tier state strategy for executions: **NATS KV** stores the hot state
during execution (per-step checkpointing for crash recovery), while **MongoDB** stores the
permanent archive (lightweight stub on creation, updates on waiting/resumed, full upsert on
terminal). Both definitions and instances are cached via **TieredCache** (L0 RAM + L1 Disk)
for fast access during execution.

## Responsibilities

- CRUD workflow definitions with versioning and MinIO storage
- CRUD workflow instance configs with TieredCache
- Trigger workflow executions via NATS JetStream
- Execute DAG traversal with 17 node executors
- Per-step state checkpointing in NATS KV
- Archive execution lifecycle events to MongoDB via Archiver module
- Timer sweep for async timeouts via Reconciler module
- Dispatch async requests to external services (JS code executor, subworkflows, triggers)

## Non-Responsibilities

- Event ingestion from external systems (HTTP Gateway)
- JavaScript code execution (JS Workflow Executor service)
- Business rule evaluation (RuleEngine)
- Event routing (Router)

## Primary Data Flow

1. External event or API call publishes a trigger message to `WORKFLOW-TRIGGER` stream
2. Runtime consumer loads instance config + definition (TieredCache), creates execution in NATS KV
3. `executeInline` traverses the DAG node-by-node with per-step KV checkpointing
4. Inline nodes (condition, set_state, log, switch) execute synchronously within the cycle
5. Async nodes (delay, wait_signal, code, subworkflow) suspend execution with a WaitRequest
6. On suspension, RuntimeService publishes `workflow.state.waiting` (with timer info for Archiver)
7. Suspended executions are resumed via `WORKFLOW-RESUME` stream (callback, timer, or signal)
8. On resume, RuntimeService publishes `workflow.state.resumed`, continues execution
9. On completion/failure, RuntimeService publishes `workflow.state.completed` or `workflow.state.failed`
10. Archiver consumes all state events, persists to MongoDB, deletes KV entry on terminal events
11. Reconciler manages timers for delay/wait_signal/wait_for nodes, publishes resume on expiration

## Docs Map

- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)
- [Plugins](plugins/index.md)
