# Timeout Configuration for ALL Async Workflow Nodes

> Status: IMPLEMENTED
> Date: 2026-03-27

## Context

Currently only `core/delay` calculates `expiresAt` and uses the Reconciler for timeout. All other async nodes (`wait_signal`, `wait_for`, `trigger_event`, `code`, `subworkflow`, plugin nodes) can wait **forever** if a callback/signal never arrives. This is a critical reliability gap.

**Goal:** Every async node MUST have a timeout with `expiresAt`. When it expires, the node either fails the execution (`TIMEOUT_EXCEEDED`) or routes to a `timeout` output handle (if `enableOutput: true`).

---

## Current Flow (how it works today)

### Step-by-step: Executor → NodeState → Suspend → DB → Reconciler

1. **DAG walker** calls `executeStep()` for each node
2. `executeStep` builds `NodeExecutionContext` (state, config, graph) and calls `executor.Execute(ctx, execCtx)`
3. **Executor** returns `NodeExecutionResult` with `NodeState` map containing `waitType` + optional `expiresAt`
4. DAG walker detects `waitType` → calls `suspendExecution(execution, nodeID, nodeType, nodeState)`
5. `suspendExecution` extracts `expiresAt` from NodeState → sets `execution.TimerExpiresAt`
6. **Checkpoint** → writes full execution to **NATS KV**
7. **Publishes** waiting state event → Archiver writes stub to **MongoDB** (with `timerExpiresAt`)
8. **Reconciler** sweeps MongoDB every 1min: `timerExpiresAt <= now + SweepInterval`
9. For each expired timer → publishes **resume message** to NATS

### Key: No extra DB writes
- Inline steps (sync): only NATS KV checkpoint
- Suspend (async): KV + MongoDB stub (ALREADY happens) — now with `timerExpiresAt` filled instead of nil
- Terminal: MongoDB full + ClickHouse

Adding `expiresAt` to all async executors adds **zero extra DB writes**.

---

## Architecture Decision

### Timeout lives at NODE level, NOT inside config

```json
{
  "id": "n_telegram_message_1",
  "type": "telegram/message",
  "label": "Telegram Message",
  "position": { "x": 195, "y": 270 },
  "config": { ... },
  "timeout": {
    "duration": 30,
    "unit": "seconds",
    "enableOutput": false
  },
  "parentNodeId": ""
}
```

### TimeoutConfig struct (Go)

```go
type TimeoutConfig struct {
    Duration     int    `json:"duration" bson:"duration"`
    Unit         string `json:"unit" bson:"unit"`
    EnableOutput bool   `json:"enableOutput" bson:"enableOutput"`
}
```

### Resolution priority

1. **Node instance** (`node.timeout`) — user configured in editor
2. **Plugin manifest default** (`nodeType.timeout`) — plugin author default
3. **Hardcoded platform default** — per waitType

### Default timeouts by waitType

| waitType | Default Duration | Rationale |
|---|---|---|
| `callback` (trigger_event, plugin, code) | 30 seconds | External HTTP calls should be fast |
| `callback` (subworkflow) | 1 hour | Child workflows can be long |
| `signal` (wait_signal) | 24 hours | Signals may arrive much later |
| `condition` (wait_for) | 24 hours | Polling conditions may take time |
| `timer` (delay) | N/A — IS the timeout | Duration is the timer itself |

### enableOutput behavior

| `enableOutput` | Timeout expires → |
|---|---|
| `false` (default) | `failExecution()` with `TIMEOUT_EXCEEDED` — workflow dies |
| `true` | Resume with `outputHandle: "timeout"` — workflow continues |

---

## Current State of Each Async Executor

| Node | waitType | Has timeout field? | Sets expiresAt? |
|---|---|---|---|
| `core/delay` | timer | Duration + unit (IS the timer) | **YES** ✅ |
| `core/code` | callback | `timeout: int` (ms) | **NO** ❌ |
| `core/subworkflow` | callback | `timeout: TimeoutConfig` | **NO** ❌ |
| `core/wait_signal` | signal | `timeout: string` + `maxTimeoutCycles` | **NO** ❌ |
| `core/wait_for` | condition | `timeout: string` + `maxTimeoutCycles` | **NO** ❌ |
| `core/trigger_event` | callback | **NONE** | **NO** ❌ |
| Plugin nodes | callback | **NONE** | **NO** ❌ |

---

## Implementation Phases

### Phase 1 — Go Entity + Config Parsing
- Add `EnableOutput` to `TimeoutConfig`
- Add `Timeout *TimeoutConfig` to `WorkflowNode` entity + DTO
- Pass timeout to `NodeExecutionContext`
- Parse `node.Timeout` in config_parsing

### Phase 2 — Go Executors
- Each async executor calculates `expiresAt` and includes in NodeState
- Shared helper: `calculateExpiresAt(timeout *TimeoutConfig, defaultDuration time.Duration) time.Time`

### Phase 3 — Go Reconciler + Resume
- Include `enableOutput` in resume message
- HandleResume: check enableOutput → route via "timeout" handle or failExecution

### Phase 4 — Plugin Manifest DSL
- Add `timeout` to nodeType in manifests
- Update plugin entity + Zod schema

### Phase 5 — Frontend Schemas + Interfaces
- Add `timeout` to WorkflowNode + PluginNodeType interfaces
- Add Zod schema field

### Phase 6 — Frontend Core Plugins
- Add timeout defaults to every async node type

### Phase 7 — Frontend UI
- Shared TimeoutConfig component (duration + unit + enableOutput)
- Render in NodeConfigPanel for async nodes

### Phase 8 — Frontend Dynamic Handles
- If enableOutput → inject "timeout" output handle

### Phase 9 — i18n
- Timeout section labels (en-US + pt-BR)

### Phase 10 — Tests
- Each executor test: verify expiresAt in NodeState

---

## Open Questions (to discuss before implementing)

1. How should `NodeExecutionContext` receive the timeout? New field `Timeout *TimeoutConfig`?
2. For wait_signal/wait_for: keep `maxTimeoutCycles` or replace with standard timeout?
3. Plugin manifest: should timeout be per-nodeType or per-operation?
4. Reconciler resume for timeout: should it use a different subject or same resume flow?

---

## Files to Modify

**Go Backend (17+ files):**
- `runtime/domain/entities/node_configs.go`
- `runtime/domain/entities/execution_context.go`
- `definitions/domain/entities/workflow_definition.go`
- `contracts/services/workflow/definitions/dto.go`
- `runtime/domain/services/config_parsing.go`
- `executors/async/wait_signal.go`
- `executors/control/wait_for.go`
- `executors/async/trigger_event.go`
- `executors/async/code.go`
- `executors/async/subworkflow.go`
- `executors/async/plugin.go`
- NEW: `executors/async/timeout_helper.go`
- `reconciler/application/services/reconciler_service.go`
- `runtime/application/services/runtime_service.go`
- `shared/types/resume_message.go`
- `plugins/domain/entities/plugin_manifest.go`

**JS Schemas (2 files):**
- `packages/schemas/src/workflows/schemas/definitions/definitions.schema.ts`
- `packages/schemas/src/workflows/schemas/plugins/plugins.schema.ts`

**Frontend (10+ files):**
- `components/workflow/interfaces/workflowNode.interface.ts`
- `components/workflow/interfaces/workflowPlugin.interface.ts`
- `components/workflow/constants/corePlugins.constant.ts`
- NEW: `components/workflow/TimeoutConfig/TimeoutConfig.vue`
- `NodeConfigPanel/NodeConfigPanel.vue`
- `useWorkflowEditorState.ts`
- `buildDefaultConfig.ts`
- `resolveNodeHandles.ts`
- i18n files

**Plugin Manifests (7 files):**
- telegram, slack, discord, teams, email, mqtt, http-request
