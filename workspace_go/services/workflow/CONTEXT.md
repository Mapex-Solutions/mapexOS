# Workflow Service

DAG execution engine of MapexOS. Designs workflows as directed acyclic graphs, triggers them via events or API, executes with state management, crash recovery, and async suspension.

## Service Info

| Field | Value |
|---|---|
| Port | `5007` |
| DB | `dev-workflow` |
| Config | `src/shared/configuration/application/config.go` |
| Module Config | `src/shared/configuration/modules/config.go` |
| Entry Point | `src/main.go` |
| Build | `go run workflow/src/main.go` |

## 3-Entity Architecture

| Entity | Module | Collection | Storage | Responsibility |
|---|---|---|---|---|
| **Definition** | `definitions` | `workflow_definitions` | TieredCache L0/L1 → MongoDB | Blueprint (DAG, states schema, externalInputs schema) |
| **Instance** | `instances` | `workflow_instances` | TieredCache L0/L1 → MongoDB | Config to run (definitionId, filled externalInputs) |
| **Execution** | `runtime` | `workflow_executions` | NATS KV (hot) + MongoDB (archive) | Mutable execution record (status, state, DAG progress) |

### Cardinality

```
Definition 1 ──→ N Instance 1 ──→ N Execution
```

- One definition can have many instances (different inputs)
- One instance can trigger many executions (each run)
- Each execution has a UUID (`workflowUUID`): user-provided OR auto-generated GUIDv4

## Modules

| Module | Path | Responsibility |
|---|---|---|
| definitions | `src/modules/definitions/` | Workflow definition CRUD, validation, MinIO storage, TieredCache |
| engine | `src/modules/engine/` | Condition evaluators, value resolver, 22 operators |
| instances | `src/modules/instances/` | Instance config CRUD, TieredCache |
| runtime | `src/modules/runtime/` | DAG execution, NATS consumers (trigger, resume), 17 node executors |
| plugins | `src/modules/plugins/` | Plugin marketplace manifests, TieredCache, FANOUT invalidation |
| credentials | `src/modules/credentials/` | Envelope encryption (AES-256-GCM), CRUD, test endpoint, loadOptions proxy |
| archiver | `src/modules/archiver/` | NATS consumer → MongoDB BulkWrite for execution state events |
| reconciler | `src/modules/reconciler/` | Timer sweep for async timeouts (delay, wait_signal, wait_for) |
| app | `src/modules/app/` | App-level bootstrap, health checks |

## API Routes

| Route | Module | Purpose |
|---|---|---|
| `/api/v1/workflow_definitions` | definitions | Definition CRUD |
| `/api/v1/workflow_instances` | instances | Instance config CRUD |
| `/api/v1/workflow_executions` | runtime | List/get/cancel/signal executions |

## State Strategy

- **NATS KV**: Hot execution state during DAG traversal — per-step checkpointing for crash recovery. Key: `exec:{uuid}`
- **MongoDB**: Permanent archive — lightweight stub on creation, full upsert on terminal
- **TieredCache (L0 RAM + L1 Disk)**: For WorkflowDefinitions, PluginManifests, AND WorkflowInstances

## Primary Data Flow

```
Event/API → WORKFLOW-TRIGGER stream
  → Runtime loads instance config (TieredCache) + definition (TieredCache)
  → Creates WorkflowExecution in KV (exec:{uuid})
  → Publishes workflow.state.created → Archiver inserts stub to MongoDB
  → executeInline traverses DAG node-by-node with per-step KV checkpoint
  → Inline nodes (condition, setState, log, switch) execute synchronously
  → Async nodes (delay, waitSignal, code, subworkflow) suspend with WaitRequest
  → Resume via WORKFLOW-RESUME stream (callback, timer, signal)
  → Terminal state → Archiver persists full execution to MongoDB, deletes KV
  → Reconciler manages timers, publishes resume on expiration
```

## Plugin System

- Plugin DSL v1.1 — spec at `docs/plugins/DSL_JSON.md`
- Manifest root: `loadOptions` map for dynamic dropdowns, `nodeTypes` for action nodes
- Properties are pure UI+Data — HTTP mapping lives in `operations.*.request.body`
- loadOptions proxy: `POST /api/v1/credentials/:id/load_options/:resourceKey`
- JS transform via goja (pure Go ES5.1, sandboxed, 10s timeout) in `packages/utils/jsrunner/`
- Credential encryption via `packages/utils/envelope/` (AES-256-GCM, Master Key + per-record DEK)

## Separation of Concerns

| Concern | Owner | Rule |
|---|---|---|
| DAG + inline processing | Workflow Service | ONLY orchestration |
| HTTP execution | Triggers Service | All external HTTP calls at runtime |
| Script execution | JS Executor | V8 sandbox, scales independently |
| Logging/observability | Events Service | — |
| Design-time HTTP (credential test, loadOptions) | Workflow Service | Exception: direct HTTP allowed |

## Key Decisions

- **Redis NOT used** — NATS KV replaces it (native CAS via revision)
- **Workflow service NEVER makes HTTP requests at runtime** — all via NATS pipeline
- **DLQ**: centralized `MAPEXOS-DLQ` stream with `DLQPolicy` struct
- **Inbound nodes**: `behavior: "inbound"` on nodeTypes replaces `triggerTypes` (design approved, not implemented)
- **NOT copying n8n** — study competitors, adapt to our microservice architecture

## Docs

Full documentation at `docs/`. See `docs/index.md` for the complete map.

| Section | Path | Content |
|---|---|---|
| Architecture | `docs/architecture/` | Module structure, DI, layers |
| Endpoints | `docs/endpoints/` | HTTP API reference |
| Configuration | `docs/configuration/` | Config keys, env vars |
| Operations | `docs/operations/` | Build, run, deploy |
| Observability | `docs/observability/` | Metrics, logging |
| Tests | `docs/tests/` | Test strategy, fixtures |
| Benchmarks | `docs/benchmarks/` | Performance baselines |
| Plugins | `docs/plugins/` | DSL JSON spec |
