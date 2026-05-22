# Architecture

## Design

The Workflow service follows **Hexagonal Architecture (Ports & Adapters)** with strict layer
separation. Each module has its own domain, application, infrastructure, and interface layers.
Dependencies flow inward: interfaces → application → domain. Cross-module communication happens
through ports (interfaces) wired via DIG dependency injection.

## 3-Entity Model

| Entity | Module | Collection | Storage |
|---|---|---|---|
| **Definition** | definitions | `workflow_definitions` | TieredCache L0/L1 → MongoDB |
| **Instance** | instances | `workflow_instances` | TieredCache L0/L1 → MongoDB |
| **Execution** | runtime | `workflow_executions` | NATS KV (hot) + MongoDB (archive) |

```
Definition 1 ──→ N Instance 1 ──→ N Execution
```

- **Definition** = reusable template (DAG, states schema, externalInputs schema)
- **Instance** = config to run (which definition, filled externalInputs, orgId)
- **Execution** = mutable run record (status, state values, DAG progress, executionPath)

## Project Structure

```
src/
├── main.go
├── bootstrap/                          # Infrastructure wiring
│   ├── config.go                       # Environment configuration
│   ├── nats.go                         # NATS Core + KV + Bus
│   ├── mongo.go                        # MongoDB manager
│   ├── redis.go                        # Redis Shared (auth only)
│   ├── cache.go                        # TieredCache (definitions + plugins + instances) + MinIO
│   ├── fiber.go                        # HTTP server
│   ├── health.go                       # Health endpoint
│   ├── metrics.go                      # Prometheus registry
│   ├── middlewares.go                  # Permission + Coverage
│   └── shutdown.go                    # Graceful shutdown hooks
│
├── shared/configuration/
│   ├── application/config.go           # DefaultConfiguration (all env vars)
│   └── modules/config.go              # Module initialization order
│
└── modules/
    ├── app/module.go                   # App initialization orchestrator
    ├── definitions/                    # Workflow definition CRUD
    ├── engine/                         # Condition evaluation + value resolution
    ├── instances/                      # Instance config CRUD + TieredCache
    ├── runtime/                        # DAG execution engine + execution API
    ├── plugins/                        # Plugin marketplace manifests
    ├── credentials/                    # Envelope encryption + CRUD
    ├── archiver/                       # State events → MongoDB persistence
    └── reconciler/                     # Timer sweep for async timeouts
```

## Module Responsibilities

### definitions
CRUD operations for workflow definitions. Stores definition JSON in MinIO (source of truth) with
TieredCache (L0 RAM + L1 Disk) for fast reads. Publishes cache invalidation events via NATS Fanout
on updates/deletes. Manages code node scripts (write/delete) in MinIO on create/update/delete.

### engine
Pure computation module with no I/O. Provides condition evaluation (comparison, datetime, string,
group operators) and value resolution (field paths, literals, expressions). Used by Runtime
executors during DAG traversal.

### instances
CRUD for workflow instance configs. Provides HTTP endpoints at `/api/v1/workflow_instances` for
creating, reading, updating, and deleting instance configurations. Uses TieredCache (L0 RAM + L1
Disk) with MongoDB fallback for fast reads during execution.

### runtime
Core execution engine. Receives trigger and resume messages via NATS consumers. Loads instance
config (TieredCache) and definition (TieredCache), creates WorkflowExecution in NATS KV, traverses
DAG with per-step checkpointing. Publishes state lifecycle events to WORKFLOW-STATE stream.

Includes 17 node executors organized in three categories:
- **Inline** (7): start, end, condition, switch, set_state, log, goto
- **Async** (5): delay, wait_signal, code, subworkflow, trigger_event
- **Control** (5): fanout, merge, sequence, loop, wait_for

### archiver
Consumes WORKFLOW-STATE stream events in batches. Writes to `workflow_executions` collection:
- `created` → lightweight InsertOne stub (~200B) for immediate frontend listing
- `waiting` → updates timerExpiresAt and status
- `resumed` → clears timerExpiresAt
- `completed`/`failed`/`cancelled` → reads full state from NATS KV, upserts complete document
  (~5-25KB), deletes KV entry

Also handles short timer fast-path: timers expiring within 1 minute get a direct NATS scheduled
message instead of waiting for the Reconciler sweep.

### reconciler
Sweeps MongoDB for executions with expired timers (timerExpiresAt <= now). Validates via NATS KV
that the execution is still waiting, then publishes resume messages to WORKFLOW-RESUME.

## Module Initialization Order

```
1. definitions    (no dependencies)
2. plugins        (no dependencies)
3. credentials    (no dependencies)
4. engine         (no dependencies, pure computation)
5. instances      (depends on TieredCache "instances")
6. runtime        (depends on engine + definitions + instances)
7. archiver       (depends on KV + MongoDB)
8. reconciler     (depends on KV + publisher)
```

## State Persistence Strategy

```
                ┌──────────────────────────────────┐
                │           TieredCache             │
                │   L0 RAM → L1 Disk → MongoDB      │
                │   For: Definitions + Instances     │
                └──────────────┬───────────────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                       │
  Runtime reads           Runtime reads            Runtime reads
  definition              instance config          execution state
  (TieredCache)           (TieredCache)            (NATS KV)
        │                      │                       │
        │                      │              ┌────────┴────────┐
        │                      │              │    NATS KV       │
        │                      │              │ Key: exec:{uuid} │
        │                      │              │ ~1-5KB per exec  │
        │                      │              └────────┬────────┘
        │                      │                       │
        │                      │                 WORKFLOW-STATE
        │                      │                 stream events
        │                      │                       │
        │                      │              ┌────────┴────────┐
        │                      │              │    Archiver      │
        │                      │              └────────┬────────┘
        │                      │                       │
        │                      │              ┌────────┴────────┐
        │                      │              │    MongoDB       │
        │                      │              │  (archive)       │
        │                      │              └─────────────────┘
        │                      │              Execution writes:
        │                      │              1. created   → stub (~200B)
        │                      │              2. waiting   → update timer
        │                      │              3. resumed   → clear timer
        │                      │              4. terminal  → full upsert (~5-25KB)
```

## Request Flow: Execution Command (Entry Point)

See [WORKFLOW-EXECUTION Stream](workflow-execution-stream.md) for full documentation.

1. Message arrives on `WORKFLOW-EXECUTION` stream with `{ mode, event, metadata, data }`
2. Execution Consumer dispatches by `mode`:
   - `newInstance` → forwards to WORKFLOW-TRIGGER with `data.instanceId`
   - `signal` → reads KV, validates, publishes to WORKFLOW-RESUME
   - `signalOrStart` → tries signal, falls back to newInstance

## Request Flow: Trigger Workflow

1. Message arrives on `WORKFLOW-TRIGGER` stream (from Execution Consumer or subworkflow)
2. RuntimeService loads instance config from TieredCache (`instanceId`)
3. RuntimeService loads definition from TieredCache (`instance.definitionId`)
4. GraphBuilder constructs DAG with adjacency lists + parsed configs
5. Generates execution UUID (user-provided or GUIDv4)
6. New WorkflowExecution created with initial state
7. Execution saved to NATS KV (`exec:{uuid}`)
8. State event `workflow.state.created` published → Archiver inserts stub
9. `executeInline` starts DAG traversal from start node

## Request Flow: Execute Inline

1. Look up current node in graph
2. Get executor for node type from registry
3. Execute node → receive result (OutputHandles, StatePatch, WaitRequest, LogEntries)
4. Apply StatePatch + NodeOutput to execution
5. **KV Put** (per-step checkpoint)
6. If WaitRequest → suspend (publish timer/callback if needed) → return
7. Resolve next nodes from OutputHandles + adjacency
8. If multiple next nodes → fanout (goroutines per branch)
9. Loop until end node, async suspension, or MaxInlineSteps reached

## Request Flow: Resume Execution

1. Message arrives on `WORKFLOW-RESUME` stream (`{ workflowUUID, nodeId, data }`)
2. RuntimeService loads execution from NATS KV
3. Loads instance config (TieredCache) + definition (TieredCache)
4. Validates resume matches current WaitRequest
5. Applies resume data (output, statePatch, signalData)
6. Clears WaitRequest, advances to next node
7. Continues with `executeInline` from the resumed node
