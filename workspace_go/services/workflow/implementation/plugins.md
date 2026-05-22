# Plugin System — Architecture & Implementation Plan

## Table of Contents

1. [Design Philosophy — MapexOS Is Not Another n8n](#1-design-philosophy--mapexos-is-not-another-n8n)
2. [Service Responsibility Map](#2-service-responsibility-map)
3. [Current Architecture — What Already Exists](#3-current-architecture--what-already-exists)
4. [Competitive Research — n8n Study](#4-competitive-research--n8n-study)
5. [Scale Analysis — The 500 Plugins Problem](#5-scale-analysis--the-500-plugins-problem)
6. [Pipeline Architecture — How Plugin Execution Works](#6-pipeline-architecture--how-plugin-execution-works)
7. [Integration Manifest Schema](#7-integration-manifest-schema)
8. [Pipeline Execution — Step-by-Step](#8-pipeline-execution--step-by-step)
9. [Backend Implementation](#9-backend-implementation)
10. [Frontend Integration](#10-frontend-integration)
11. [Credential System](#11-credential-system)
12. [Webhook & Trigger Architecture](#12-webhook--trigger-architecture) → [`plugin-webhook-trigger.md`](./plugin-webhook-trigger.md)
13. ~~System-Managed Resources~~ (merged into Section 12)
14. [Storage & Caching Strategy](#14-storage--caching-strategy)
15. [What Needs To Be Built](#15-what-needs-to-be-built)
16. [Migration Path](#16-migration-path)
17. [File Reference](#17-file-reference)

---

## 1. Design Philosophy — MapexOS Is Not Another n8n

> **CRITICAL RULE**: We study competitors (n8n, Zapier, Temporal) to understand patterns and learn
> from their decisions. We do NOT copy their architecture. MapexOS has its own microservice
> topology with strict separation of concerns. Each service has ONE responsibility and scales
> INDEPENDENTLY. This is what makes MapexOS enterprise-grade.

### 1.1 Why We Don't Copy n8n

n8n is a **monolith** — one Node.js process handles webhook ingestion, DAG execution, HTTP
requests, script execution, credential storage, and logging. This works for small teams but
fails catastrophically at enterprise scale:

- Cannot scale HTTP execution independently from DAG traversal
- Cannot scale script execution independently from HTTP execution
- Single failure domain — one bad HTTP call blocks the entire engine
- No crash recovery for mid-HTTP-call failures
- No backpressure — overloaded HTTP calls starve the DAG walker

### 1.2 MapexOS Core Principle: Separation of Concerns

Every piece of work has exactly ONE service responsible for it. No service does another
service's job. Every inter-service boundary is NATS (async, durable, scalable).

```
┌──────────────────────────────────────────────────────────────┐
│                     WORKFLOW SERVICE (Go)                      │
│                                                                │
│  DOES:                          NEVER DOES:                    │
│  ✅ DAG traversal (inline)      ❌ HTTP requests               │
│  ✅ Pipeline orchestration       ❌ Script execution            │
│  ✅ Suspend / Resume / KV        ❌ AI processing               │
│  ✅ Dispatch async work          ❌ External I/O of any kind    │
│                                                                │
│  The workflow service is LIGHTWEIGHT.                          │
│  It decides WHAT to do, not HOW to do it.                     │
└───────────────┬──────────────────────┬────────────────────────┘
                │ NATS                 │ NATS
                ▼                      ▼
┌───────────────────────┐   ┌──────────────────────────────────┐
│  TRIGGERS SERVICE     │   │  JS WORKFLOW EXECUTOR            │
│  (Go)                 │   │  (Node.js)                       │
│                       │   │                                  │
│  DOES:                │   │  DOES:                           │
│  ✅ ALL HTTP requests │   │  ✅ ALL script execution          │
│  ✅ Template resolving│   │  ✅ preScript / postScript hooks  │
│  ✅ Retry / backoff   │   │  ✅ Code node execution           │
│  ✅ Circuit breaking  │   │  ✅ V8 sandbox (32MB, 10s)        │
│                       │   │                                  │
│  Scales: 100 pods     │   │  Scales: 50 pods                 │
│  if HTTP is bottleneck│   │  if CPU is bottleneck            │
└───────────────────────┘   └──────────────────────────────────┘
                │ NATS                 │ NATS
                ▼                      ▼
┌───────────────────────┐   ┌──────────────────────────────────┐
│  EVENTS SERVICE       │   │  AI SERVICE (future)             │
│                       │   │                                  │
│  DOES:                │   │  DOES:                           │
│  ✅ Step logging      │   │  ✅ LLM calls                    │
│  ✅ Audit trail       │   │  ✅ Embeddings                   │
│  ✅ Observability     │   │  ✅ AI pipeline processing        │
└───────────────────────┘   └──────────────────────────────────┘
```

### 1.3 The Key Insight: Pipelines + KV Checkpoints

When a plugin node executes, it becomes a **pipeline** — a sequence of steps where each step:
1. Is dispatched to the appropriate service via NATS
2. Executes asynchronously in that service
3. Returns result via NATS callback to WORKFLOW-RESUME
4. Workflow saves result to **NATS KV** (checkpoint)
5. Workflow dispatches the next step (or continues DAG if done)

This gives us:
- **Crash recovery** — if any service dies mid-step, the workflow resumes from last KV checkpoint
- **Independent scaling** — each service scales based on its own bottleneck
- **Full observability** — every step generates state events for the Archiver
- **Accumulated state** — each step can read results from all previous steps

---

## 2. Service Responsibility Map

| Service | Language | Responsibility | Scales Based On | Never Does |
|---------|----------|---------------|-----------------|------------|
| **Workflow** | Go | DAG traversal, inline processing, pipeline orchestration, KV checkpoint | Concurrent workflows | HTTP, scripts, AI |
| **Triggers** | Go | ALL HTTP execution, template resolution, retry, circuit breaker | HTTP throughput | DAG logic, scripts |
| **JS Executor** | Node.js | ALL script execution (V8 sandbox), preScript/postScript, code nodes | CPU for V8 | HTTP, DAG logic |
| **Events** | Go | Logging, audit trail, step observability, archiving | Write throughput | Execution logic |
| **AI** (future) | TBD | LLM calls, embeddings, AI pipeline | GPU / API limits | HTTP, scripts |

### 2.1 Inter-Service Communication

ALL communication between services is via **NATS JetStream**. No direct HTTP between services.
No shared databases between services. No synchronous RPC.

```
Workflow → NATS → Triggers    (HTTP execution request)
Triggers → NATS → Workflow    (HTTP result callback)
Workflow → NATS → JS Executor (script execution request)
JS Exec  → NATS → Workflow    (script result callback)
Workflow → NATS → Events      (state change events)
Workflow → NATS → AI Service  (AI execution request, future)
AI Svc   → NATS → Workflow    (AI result callback, future)
```

---

## 3. Current Architecture — What Already Exists

### 3.1 Backend: ExecutorRegistry (Go)

The runtime already uses a registry pattern for dispatching node execution:

```
src/modules/runtime/domain/
├── entities/
│   └── executor_port.go          # NodeExecutor interface
└── executors/
    ├── executor.go               # ExecutorRegistry (map[string]NodeExecutor)
    ├── registry_builder.go       # BuildRegistry — registers all 17 executors
    ├── inline/                   # 7 inline executors
    ├── async/                    # 5 async executors
    └── control/                  # 5 control executors
```

**NodeExecutor interface:**

```go
type NodeExecutor interface {
    Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error)
    NodeType() string
}
```

**ExecutorRegistry:**

```go
type ExecutorRegistry struct {
    executors map[string]NodeExecutor
}

func (r *ExecutorRegistry) Register(executor NodeExecutor) {
    r.executors[executor.NodeType()] = executor
}

func (r *ExecutorRegistry) Get(nodeType string) (NodeExecutor, error) {
    executor, ok := r.executors[nodeType]
    if !ok {
        return nil, fmt.Errorf("%w: %s", ErrExecutorNotFound, nodeType)
    }
    return executor, nil
}
```

Currently, 17 node executors are hardcoded at boot via `BuildRegistry()`. Every integration plugin
must be a compiled Go struct implementing `NodeExecutor`. This doesn't scale.

### 3.2 Frontend: Plugin Registry (Pinia Store)

```
workspace_js/
├── packages/
│   ├── workflow-sdk/src/interfaces/
│   │   └── workflowPlugin.interface.ts    # WorkflowPlugin, PluginNodeType, etc.
│   └── workflow-plugin-core/
│       ├── src/constants/corePlugins.constant.ts  # 6 core plugins, 18 node types
│       ├── src/nodes/                             # Canvas + config components
│       ├── src/validators/                        # Node config validators
│       └── src/i18n/                              # en-US, pt-BR translations
└── apps/mapexOS/src/
    └── stores/pluginRegistry/
        ├── index.ts   # Pinia store definition
        ├── state.ts   # plugins: Map, nodeTypeMap: Map
        ├── actions.ts # registerPlugin, unregisterPlugin, getNodeType
        └── getters.ts # catalog (grouped by category), counts
```

**WorkflowPlugin interface:**

```typescript
interface WorkflowPlugin {
  id: string;
  name: string;
  version: string;
  category: PluginCategory;
  icon: string;
  nodeTypes: PluginNodeType[];
  onActivate?: (context: PluginActivationContext) => void;
  onDeactivate?: () => void;
}
```

**PluginNodeType interface (key fields):**

```typescript
interface PluginNodeType {
  type: string;                                    // "core/delay", "slack/send_message"
  label: string;
  icon: string;
  color: string;
  description: string;
  inputs: HandleDefinition[];
  outputs: HandleDefinition[];
  configSchema: Record<string, unknown>;

  // Form rendering — THE KEY TO SCALING
  properties?: NodePropertyDefinition[];           // Declarative auto-form
  configComponent?: Component;                     // Custom Vue component (optional)

  canvasComponent?: Component;                     // Defaults to GenericWorkflowNode
  validate?: (config) => ValidationResult;
  defaults?: Record<string, unknown>;
  resolveOutputs?: HandleResolver;
  resolveInputs?: HandleResolver;
  outputHints?: Array<{ path: string; description: string }>;
}
```

**NodePropertyDefinition — auto-form system:**

```typescript
interface NodePropertyDefinition {
  name: string;              // config key
  displayName: string;       // form label
  type: 'string' | 'number' | 'boolean' | 'options' | 'json';
  default: unknown;
  hint?: string;
  required?: boolean;
  options?: { label: string; value: string | number }[];
  displayOptions?: { show?: Record<string, unknown[]> };
}
```

When `properties[]` is defined and no `configComponent` is provided, the editor renders a
`DynamicNodeForm` that auto-generates the form fields from the property definitions.

### 3.3 Registration Flow (Current)

```
Boot → corePlugins.constant.ts → bootWorkflowPlugins() → registerPlugin() per plugin
                                                              ↓
                                                  pluginRegistry store
                                                  ├── plugins Map
                                                  └── nodeTypeMap Map
```

All core plugins are currently hardcoded in the `workflow-plugin-core` package. Integration
plugins (Slack, Google, etc.) don't exist yet.

### 3.4 JS Workflow Executor (Existing)

The JS Workflow Executor service already provides:

- **Piscina V8 worker pool** with isolated-vm (32MB heap, 10s timeout per script)
- **NATS consumer** for code node execution (`WORKFLOW-JS-CODE` stream)
- **TieredCache** for script source (L0 RAM → L1 Disk → L2 MinIO)
- **Bytecode cache** for compiled V8 scripts (fast re-execution)
- **Sandboxed globals**: `event`, `state`, `variables`, `nodes` — NO filesystem, NO network
- **Callback pattern**: publishes result to `workflow.resume.callback.{instanceId}`

This is the **same infrastructure** we reuse for plugin preScript/postScript hooks.

### 3.5 Triggers Service (Existing)

The Triggers Service already provides:

- **PlaceholderResolver**: `{{path.to.field}}` template syntax with recursive map resolution
- **HTTPExecutor**: HTTP request building (method, url, headers, body), status validation
- **ExecutorRegistry**: factory pattern for different executor types
- **Batch processing**: parallel execute → flush → ACK/NACK
- **DLQ pattern**: with retry policies and dead letter queue

This is the **same infrastructure** we reuse for plugin HTTP execution.

---

## 4. Competitive Research — n8n Study

> **WARNING**: This section documents what we LEARNED from n8n, not what we COPY.
> n8n's architecture is a monolith where one process does everything. MapexOS is a
> microservice platform where each concern has its own service. The patterns are interesting
> but the architecture is fundamentally different.

### 4.1 n8n Architecture

n8n uses an **npm package per integration** model:

```
packages/
└── nodes-base/src/nodes/          # 400+ integrations
    ├── Slack/
    │   ├── Slack.node.ts          # INodeType implementation
    │   ├── SlackDescription.ts    # properties[] definitions
    │   └── v2/                    # versioned nodes
    ├── Google/
    │   ├── Gmail.node.ts
    │   ├── GoogleSheets.node.ts
    │   └── ...
    └── ...
```

| Aspect | n8n Approach | MapexOS Approach |
|--------|-------------|-----------------|
| **Execution** | In-process, same Node.js runtime | Service-per-concern, NATS between |
| **HTTP** | In the engine process | Triggers Service (dedicated) |
| **Scripts** | In the engine process, no sandbox | JS Executor (V8 isolate, 32MB) |
| **Scaling** | Vertical only (bigger machine) | Horizontal per service |
| **Properties** | Declarative `properties[]` → auto forms | Same — JSON manifest `properties[]` |
| **Auth** | Separate `credentials` system | Same concept, different storage |
| **Crash recovery** | None during HTTP call | KV checkpoint between every step |
| **Backpressure** | None | NATS JetStream native |

### 4.2 What We Learned (Patterns Worth Adopting)

| Pattern | n8n Source | MapexOS Adaptation |
|---------|-----------|-------------------|
| Declarative `properties[]` | `INodeProperties` auto-form | Already implemented in DynamicNodeForm |
| Credential as separate entity | `ICredentialType` | Section 11 — Credential System |
| Template syntax in requests | `={{value}}` | Our `{{config.field}}` syntax |
| Webhook lifecycle | `webhookMethods: check/create/delete` | Saga Workflow automation (Section 12) |
| OAuth2 inheritance | `extends: ['oAuth2Api']` | Credential type inheritance |
| Dynamic dropdowns | `loadOptionsDependsOn` | Backend proxy endpoint + `dependsOn` field |
| fixedCollection | Nested structured objects in forms | `NodePropertyDefinition` type extension |

### 4.3 What We Explicitly Reject

| n8n Pattern | Why We Reject It |
|-------------|-----------------|
| **In-process HTTP** | Blocks DAG walker. Can't scale HTTP independently. No crash recovery. |
| **In-process scripting** | No sandbox. Malicious code = full process access. |
| **Monolithic bundle** | 400+ nodes = ~250MB. Can't hot-reload. Can't isolate per tenant. |
| **8 built-in postReceive types** | Over-engineering. Our V8 postScript handles all transforms. |
| **Webhook per workflow node** | Bypasses DataSource pipeline. Loses rate limit, auth, multi-protocol. |
| **Per-property routing** | `routing.send.type: 'body'` is clever but fragmented. Our `execution.request.body` with templates is clearer. |

### 4.4 n8n Plugin Study — Complexity Distribution

From studying 400+ n8n nodes and 11 packages from skriptfabrik:

| Plugin | n8n LOC | MapexOS Pipeline | Why |
|--------|---------|-----------------|-----|
| Telegram | ~1,500 | HTTP only (1 step) | Pure REST, token in URL |
| Slack | ~2,000 | HTTP only (1 step) | Pure REST, OAuth bearer |
| SendGrid | ~800 | HTTP only (1 step) | Pure REST, API key header |
| GitHub | ~3,000 | HTTP only OR HTTP+postScript (1-2 steps) | Most ops 1 step, pagination needs postScript |
| Google Sheets | ~4,000 | preScript+HTTP+postScript (3 steps) | OAuth refresh + response mapping |
| Salesforce | ~6,000 | preScript+HTTP+postScript (3 steps) | ~40 loadOptions, SOQL builder, cascading dropdowns |
| Stripe | ~3,500 | HTTP+postScript (2 steps) | Cursor pagination, nested objects |
| AWS S3 | ~2,500 | preScript+HTTP+postScript (3 steps) | HMAC-SHA256 signing, multipart |
| Postgres | N/A | Native (Go executor via NATS) | Binary protocol, connection pooling |

**Key finding**: ~90% of plugins need only 1 step (HTTP only). ~8% need 2-3 steps (with JS hooks).
~2% need native executors (binary protocols, DB connections).

### 4.5 Patterns Discovered in Deep Analysis

**From Stripe**: `fixedCollection` for nested objects (address, metadata), cursor pagination via
`has_more` + `starting_after`, webhook signature with timestamp validation (replay attack prevention).

**From Salesforce**: Cascading dropdowns via `loadOptionsDependsOn`, SOQL query builder with
`conditionsUi`, ~40 `loadOptions` methods for dynamic field loading.

**From AWS S3**: HMAC-SHA256 signing (SigV4) — perfect preScript use case. Multipart chunked
upload with MD5 integrity per part.

**From skriptfabrik (11 packages)**: Optimistic locking retry (Kaufland), JWT service account
auth (Google Enhanced), multiple auth types per node (Sentry), `notice` property type for
informational messages.

---

## 5. Scale Analysis — The 500 Plugins Problem

### 5.1 Why Per-Plugin Microservices Don't Work

If every integration (Google, Slack, GitHub, ...) is a separate container:

- 500 containers × ~50MB RAM = **25 GB of RAM** just idle
- 500 separate CI/CD pipelines
- Operational nightmare

### 5.2 How MapexOS Solves It: JSON Manifests + Generic Services

One JSON manifest per plugin (~5KB), stored in MongoDB:

- 500 plugins × 5KB = **2.5 MB** total
- TieredCache (L0 RAM + L1 Disk): ~10μs per lookup
- Compared to 25 GB RAM for 500 containers: **10,000x more efficient**

The manifest defines WHAT to do (UI form + execution plan). The existing services
(Triggers, JS Executor) do the HOW. No new services per plugin. No new containers.

---

## 6. Pipeline Architecture — How Plugin Execution Works

### 6.1 Core Concept

A plugin node execution is a **pipeline** of steps. The workflow service orchestrates the
pipeline by dispatching each step to the appropriate service and checkpointing between steps.

```
Plugin Node → IntegrationExecutor builds pipeline → DAG suspends
    │
    │  Pipeline = ordered list of steps:
    │  [preScript?, http, postScript?]
    │
    │  Each step is a full async cycle:
    │  ┌──────────────────────────────────────────────┐
    │  │ 1. Workflow dispatches step → NATS            │
    │  │ 2. Service executes (Triggers OR JS Executor) │
    │  │ 3. Service publishes callback → WORKFLOW-RESUME│
    │  │ 4. Workflow receives callback                  │
    │  │ 5. Workflow saves result → KV checkpoint       │
    │  │ 6. If more steps → goto 1                      │
    │  │ 7. If done → continue DAG                      │
    │  └──────────────────────────────────────────────┘
```

### 6.2 Pipeline Variants

Not every plugin needs all steps. The pipeline is built from the manifest:

| Plugin Type | Pipeline Steps | Example |
|-------------|---------------|---------|
| Simple HTTP (90%) | `[http]` | Telegram, Slack, SendGrid |
| HTTP + response transform | `[http, postScript]` | Stripe, GitHub pagination |
| Auth computation + HTTP | `[preScript, http]` | AWS SigV4 signing |
| Full pipeline | `[preScript, http, postScript]` | Salesforce, Google Sheets |

### 6.3 The Full Pipeline Flow (3 Steps Example)

```
Workflow DAG → IntegrationExecutor
    │
    │ Builds pipeline: [preScript, http, postScript]
    │ Saves pipeline in NodeState → KV checkpoint
    │ DAG suspends (waitType: "plugin_pipeline")
    │
    │ ===== Step 0: preScript =====
    │
    ├─→ NATS dispatch → JS Executor
    │                        │
    │                        │ V8 runs preScript
    │                        │ (e.g., compute HMAC signature)
    │                        │
    │                        └─→ NATS → WORKFLOW-RESUME
    │                                        │
    │              Workflow RESUME ← KV checkpoint
    │              (saves preScript result in pipeline.results[0])
    │              (pipeline.currentStep = 1)
    │
    │ ===== Step 1: http =====
    │
    ├─→ NATS dispatch → Triggers Service
    │                        │
    │                        │ Resolve templates (using results from step 0)
    │                        │ Make HTTP request
    │                        │ Retry on 429/5xx
    │                        │
    │                        └─→ NATS → WORKFLOW-RESUME
    │                                        │
    │              Workflow RESUME ← KV checkpoint
    │              (saves HTTP response in pipeline.results[1])
    │              (pipeline.currentStep = 2)
    │
    │ ===== Step 2: postScript =====
    │
    ├─→ NATS dispatch → JS Executor
    │                        │
    │                        │ V8 runs postScript
    │                        │ (receives HTTP response from step 1)
    │                        │ (e.g., reshape data, extract fields)
    │                        │
    │                        └─→ NATS → WORKFLOW-RESUME
    │                                        │
    │              Workflow RESUME ← KV checkpoint
    │              (saves postScript result as final output)
    │              (pipeline.currentStep = DONE)
    │
    └─→ Continue DAG → next node
```

### 6.4 Why Every Step Returns to Workflow (KV Checkpoint)

This is **critical**. Each step does NOT forward directly to the next service. It ALWAYS
returns to the workflow service first. Why:

1. **Crash recovery** — If Triggers Service crashes during HTTP, the workflow knows which
   step it was on. On restart, it re-dispatches step 1 (not step 0 again).

2. **State accumulation** — Each step's result is saved in KV. The next step can read
   results from ALL previous steps, not just the immediately previous one.

3. **Observability** — Each resume generates `workflow.state.resumed` → Archiver logs it.
   Full step-by-step timeline in MongoDB.

4. **Single orchestrator** — The workflow service is the ONLY entity that knows the
   pipeline state. No service-to-service coupling. Triggers doesn't know about JS Executor.
   JS Executor doesn't know about Triggers.

### 6.5 Pipeline NodeState in KV

```json
{
  "waitType": "plugin_pipeline",
  "pipeline": {
    "steps": [
      { "service": "js", "action": "preScript", "scriptRef": "scripts/node-abc/preScript" },
      { "service": "triggers", "action": "http", "request": { "method": "POST", "url": "..." } },
      { "service": "js", "action": "postScript", "scriptRef": "scripts/node-abc/postScript" }
    ],
    "currentStep": 0,
    "results": {}
  },
  "credentials": { "ref": "cred-xyz", "resolved": { "botToken": "***", "baseUrl": "..." } },
  "config": { "channel": "#general", "text": "Hello" },
  "context": {
    "state": { "userId": "user-123" },
    "event": { "type": "user.created" },
    "nodes": { "prev-node": { "result": 42 } }
  }
}
```

As each step completes, `results` accumulates:

```json
{
  "results": {
    "0": { "signature": "hmac-sha256-abc...", "timestamp": "2026-03-12T..." },
    "1": { "status": 200, "data": { "ok": true, "result": { "message_id": 123 } } },
    "2": { "output": { "messageId": 123 }, "statePatch": { "lastSent": "2026-03-12T..." } }
  }
}
```

---

## 7. Integration Manifest Schema

### 7.1 Complete Manifest Structure

```json
{
  "_id": "ObjectId",
  "pluginId": "slack",
  "name": "Slack",
  "version": "1.0.0",
  "category": "integrations",
  "icon": "chat",
  "color": "purple-7",
  "description": "Send messages and interact with Slack",

  "credential": {
    "type": "slackOAuth2",
    "required": true
  },

  "nodeTypes": [
    {
      "type": "slack/send_message",
      "label": "Send Message",
      "icon": "send",
      "description": "Send a message to a Slack channel",
      "inputs": [{ "id": "in", "label": "In", "position": "top" }],
      "outputs": [
        { "id": "success", "label": "Success", "position": "bottom", "color": "#4caf50" },
        { "id": "error", "label": "Error", "position": "bottom", "color": "#ef5350" }
      ],
      "properties": [
        {
          "name": "channel",
          "displayName": "Channel",
          "type": "string",
          "default": "",
          "required": true,
          "hint": "Slack channel name or ID (e.g., #general)"
        },
        {
          "name": "text",
          "displayName": "Message",
          "type": "string",
          "default": "",
          "required": true,
          "hint": "Message text — supports {{state.fieldName}} templates"
        },
        {
          "name": "unfurl_links",
          "displayName": "Unfurl Links",
          "type": "boolean",
          "default": true
        }
      ],
      "defaults": {
        "channel": "",
        "text": "",
        "unfurl_links": true
      },
      "outputHints": [
        { "path": "ts", "description": "Message timestamp (unique ID)" },
        { "path": "channel", "description": "Channel where message was sent" }
      ],
      "execution": {
        "request": {
          "method": "POST",
          "url": "https://slack.com/api/chat.postMessage",
          "headers": {
            "Content-Type": "application/json"
          },
          "body": {
            "channel": "{{config.channel}}",
            "text": "{{config.text}}",
            "unfurl_links": "{{config.unfurl_links}}"
          }
        },
        "response": {
          "successPath": "ok",
          "dataPath": "message",
          "errorPath": "error"
        },
        "retry": {
          "maxAttempts": 3,
          "backoffMs": 1000,
          "retryOnStatus": [429, 500, 502, 503]
        }
      }
    }
  ],

  "enabled": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### 7.2 Execution Block — Without Hooks (1 Step Pipeline)

When no hooks are defined, the pipeline has a single step: HTTP. The Triggers Service
handles everything (template resolution, HTTP call, response parsing, retry).

```json
{
  "execution": {
    "request": {
      "method": "POST|GET|PUT|PATCH|DELETE",
      "url": "https://api.example.com/v1/resource",
      "headers": { "Content-Type": "application/json" },
      "query": { "page": "{{config.page}}" },
      "body": { "field": "{{config.value}}" }
    },
    "response": {
      "successPath": "ok",
      "dataPath": "data",
      "errorPath": "error"
    },
    "retry": {
      "maxAttempts": 3,
      "backoffMs": 1000,
      "retryOnStatus": [429, 500, 502, 503]
    }
  }
}
```

### 7.3 Execution Block — With Hooks (2-3 Step Pipeline)

When hooks are defined, each hook becomes a pipeline step dispatched to JS Executor.

```json
{
  "execution": {
    "hooks": {
      "preScript": {
        "script": "// Compute HMAC signature\nconst sig = computeHMAC(auth.secret, JSON.stringify(request.body));\nrequest.headers['X-Signature'] = sig;\nreturn request;",
        "timeout": 5000
      },
      "postScript": {
        "script": "// Reshape response\nconst items = response.data.items.map(i => ({ id: i.id, name: i.name }));\nreturn { output: { items, total: items.length } };",
        "timeout": 10000
      }
    },
    "request": {
      "method": "POST",
      "url": "https://api.example.com/v1/resource",
      "headers": { "Content-Type": "application/json" },
      "body": { "channel": "{{config.channel}}" }
    },
    "response": {
      "successPath": "ok",
      "dataPath": "data",
      "errorPath": "error"
    }
  }
}
```

**Pipeline built from this manifest:**
```
steps = [
  { service: "js",       action: "preScript",  script: "..." },
  { service: "triggers", action: "http",        request: { method, url, headers, body } },
  { service: "js",       action: "postScript",  script: "..." }
]
```

### 7.4 Template Resolution

All `{{...}}` expressions are resolved at execution time by the Triggers Service (for HTTP
steps) or by the JS Executor (available as context variables in scripts):

| Prefix | Source | Example |
|--------|--------|---------|
| `config.*` | Node configuration (user inputs) | `{{config.channel}}` |
| `state.*` | Workflow state variables | `{{state.userId}}` |
| `event.*` | Trigger event payload | `{{event.data.name}}` |
| `credentials.*` | Decrypted credential values | `{{credentials.botToken}}` |
| `nodes.*` | Previous node outputs | `{{nodes.nodeId.field}}` |
| `pipeline.*` | Results from previous pipeline steps | `{{pipeline.0.signature}}` |

### 7.5 NodePropertyDefinition Extensions

For richer integration forms, extend `NodePropertyDefinition` with new types:

```typescript
type NodePropertyType =
  | 'string' | 'number' | 'boolean' | 'options' | 'json'   // existing
  | 'fieldValue'       // FieldValueSelector (event/state/literal/node_output)
  | 'credential'       // Credential picker (select saved credentials)
  | 'code'             // Monaco editor (for hook scripts)
  | 'collection'       // Array of key-value pairs
  | 'multiOptions'     // Multi-select dropdown
  | 'fixedCollection'  // Nested structured objects (address, metadata)
  | 'notice';          // Informational text (not an input field)
```

**fixedCollection** — for nested object forms:
```json
{
  "name": "address",
  "displayName": "Address",
  "type": "fixedCollection",
  "options": [
    {
      "displayName": "Details",
      "name": "details",
      "values": [
        { "name": "city", "displayName": "City", "type": "string" },
        { "name": "country", "displayName": "Country", "type": "string" }
      ]
    }
  ]
}
```

**Dynamic dropdowns** — with backend proxy for option loading:

`loadOptions` are defined as a **map at the manifest root level**. Properties reference them by key:

```json
// Manifest root level — loadOptions map
"loadOptions": {
  "getRecordTypes": {
    "request": { "method": "GET", "path": "/api/record-types?resource={{dependsOn.resource}}" },
    "dataPath": "data",
    "valuePath": "id",
    "labelPath": "name"
  }
}

// In a property — references the key
{
  "name": "recordType",
  "displayName": "Record Type",
  "type": "options",
  "typeOptions": {
    "loadOptions": "getRecordTypes",
    "loadOptionsDependsOn": ["resource"]
  }
}
```

Dynamic dropdown requests go through a **backend proxy endpoint**
(`POST /api/v1/credentials/:credentialId/load_options/:resourceKey`) to avoid exposing credentials to the browser.
The backend decrypts the credential, resolves template variables, makes the HTTP call, and returns `[{label, value}]`.

**Two extraction modes:**
- **Simple** (80%): `dataPath` + `valuePath` + `labelPath` extract directly from the HTTP response.
- **Transform** (20%): An optional inline JS ES5 script receives the response as `data` and returns `[{label, value}]`. Used for filtering, concatenation, flattening nested structures, etc. Executed server-side via goja (pure Go ES5.1 runtime) in a sandboxed environment.

---

## 8. Pipeline Execution — Step-by-Step

### 8.1 IntegrationExecutor — Pipeline Builder

The single executor that handles all plugin nodes. It reads the manifest, builds the pipeline,
and suspends the DAG.

```go
type IntegrationExecutor struct {
    manifestIndex map[string]*PluginNodeManifest  // RAM index: nodeType → manifest
    credStore     CredentialStorePort             // Encrypted credential access
}

func (e *IntegrationExecutor) NodeType() string {
    return "integration" // fallback executor in registry
}

func (e *IntegrationExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    nodeType := execCtx.NodeType // e.g., "slack/send_message"

    // 1. O(1) lookup in RAM index
    manifest, ok := e.manifestIndex[nodeType]
    if !ok {
        return nil, fmt.Errorf("no manifest found for node type: %s", nodeType)
    }

    // 2. Resolve credential
    creds, _ := e.credStore.Get(ctx, manifest.PluginID, execCtx.OrgID)

    // 3. Build pipeline from manifest
    steps := buildPipeline(manifest.Execution)

    // 4. Suspend DAG with pipeline in NodeState
    return &entities.NodeExecutionResult{
        OutputHandles: []string{},  // DAG suspends (no output handles = suspend)
        NodeState: map[string]interface{}{
            "waitType": "plugin_pipeline",
            "pipeline": map[string]interface{}{
                "steps":       steps,
                "currentStep": 0,
                "results":     map[string]interface{}{},
            },
            "credentials": creds,
            "config":      execCtx.ParsedConfig,
            "context": map[string]interface{}{
                "state":  execCtx.State,
                "event":  execCtx.EventPayload,
                "nodes":  execCtx.NodeOutputs,
            },
        },
    }, nil
}

// buildPipeline constructs the ordered step list from the manifest execution block.
func buildPipeline(exec *ExecutionBlock) []PipelineStep {
    var steps []PipelineStep

    if exec.Hooks != nil && exec.Hooks.PreScript != nil {
        steps = append(steps, PipelineStep{
            Service: "js",
            Action:  "preScript",
            Script:  exec.Hooks.PreScript.Script,
            Timeout: exec.Hooks.PreScript.Timeout,
        })
    }

    steps = append(steps, PipelineStep{
        Service: "triggers",
        Action:  "http",
        Request: exec.Request,
        Response: exec.Response,
        Retry:   exec.Retry,
    })

    if exec.Hooks != nil && exec.Hooks.PostScript != nil {
        steps = append(steps, PipelineStep{
            Service: "js",
            Action:  "postScript",
            Script:  exec.Hooks.PostScript.Script,
            Timeout: exec.Hooks.PostScript.Timeout,
        })
    }

    return steps
}
```

### 8.2 Pipeline Dispatch in lifecycle.go

New case in `dispatchByNodeType()`:

```go
func (s *RuntimeService) dispatchByNodeType(
    instance *entities.WorkflowInstance,
    nodeID string,
    nodeType string,
    nodeState map[string]interface{},
) error {
    waitType := mapget.String(nodeState, "waitType")

    switch waitType {
    // ... existing cases: "callback", "timer", "signal" ...

    case "plugin_pipeline":
        return s.dispatchPipelineStep(instance, nodeID, nodeState)

    default:
        return nil
    }
}
```

### 8.3 Pipeline Step Dispatch

```go
func (s *RuntimeService) dispatchPipelineStep(
    instance *entities.WorkflowInstance,
    nodeID string,
    nodeState map[string]interface{},
) error {
    pipeline := extractPipeline(nodeState)
    step := pipeline.Steps[pipeline.CurrentStep]

    // Build context with accumulated results from previous steps
    ctx := buildStepContext(pipeline, nodeState)

    switch step.Service {
    case "js":
        return s.deps.RuntimePublisher.DispatchPluginJS(
            instance, nodeID, step, ctx,
        )
    case "triggers":
        return s.deps.RuntimePublisher.DispatchPluginHTTP(
            instance, nodeID, step, ctx,
        )
    case "ai":
        return s.deps.RuntimePublisher.DispatchPluginAI(
            instance, nodeID, step, ctx,
        )
    default:
        return fmt.Errorf("unknown pipeline service: %s", step.Service)
    }
}
```

### 8.4 Pipeline Resume Handler

When a step callback arrives, the workflow decides: next step or continue DAG?

```go
func (s *RuntimeService) handlePluginPipelineResume(
    instance *entities.WorkflowInstance,
    nodeID string,
    callbackResult map[string]interface{},
) (string, error) {
    nodeState := instance.NodeStates[nodeID]
    pipeline := extractPipeline(nodeState)

    // Save result of current step
    stepIdx := fmt.Sprintf("%d", pipeline.CurrentStep)
    pipeline.Results[stepIdx] = callbackResult
    pipeline.CurrentStep++

    // Update NodeState in instance
    updatePipelineInNodeState(nodeState, pipeline)

    // KV checkpoint — ALWAYS between steps
    if err := s.checkpoint(instance); err != nil {
        return "", err
    }

    // Pipeline complete?
    if pipeline.CurrentStep >= len(pipeline.Steps) {
        // Extract final output from last step result
        lastResult := pipeline.Results[fmt.Sprintf("%d", len(pipeline.Steps)-1)]
        output := extractOutput(lastResult)
        statePatch := extractStatePatch(lastResult)

        // Apply state patch to instance
        if statePatch != nil {
            for k, v := range statePatch {
                instance.State[k] = v
            }
        }

        // Save node output
        instance.NodeOutputs[nodeID] = output

        // Continue DAG from this node
        return nodeID, nil
    }

    // More steps → dispatch next
    instance.Status = entities.StatusWaiting
    return "", s.dispatchPipelineStep(instance, nodeID, nodeState)
}
```

### 8.5 NATS Streams & Subjects

| Stream | Subject | Direction | Consumer | Location |
|--------|---------|-----------|----------|----------|
| `WORKFLOW-PLUGIN-HTTP` | `workflow.plugin.http` | Workflow → Triggers | `triggers-plugin-http` | **Triggers Service** |
| `WORKFLOW-JS-CODE` | `workflow.plugin.js` | Workflow → JS Exec | `js-executor-plugin` | **JS Executor** |
| `WORKFLOW-RESUME` | `workflow.resume.callback.{instanceId}` | Any → Workflow | `workflow-resume` (existing) | Workflow Service |
| `WORKFLOW-STATE` | `workflow.state.*` | Workflow → Archiver | Existing | Events/Archiver |

**Key**: `WORKFLOW-PLUGIN-HTTP` is consumed by the **Triggers Service**, NOT by the workflow
service. The workflow service NEVER makes HTTP calls.

### 8.6 Plugin HTTP Consumer (in Triggers Service)

New NATS consumer added to the Triggers Service that handles generic HTTP execution
for plugin nodes. Reuses the existing PlaceholderResolver and HTTPExecutor patterns.

```go
// In Triggers Service — NOT in Workflow Service
type PluginHTTPConsumer struct {
    httpExecutor     *HTTPExecutor         // Existing, battle-tested
    templateEngine   *PlaceholderResolver  // Existing, reused
}

func (c *PluginHTTPConsumer) Handle(ctx context.Context, req *PluginHTTPRequest) {
    // 1. Build template context from pipeline results + credentials
    tplCtx := buildTemplateContext(req)

    // 2. Resolve all {{...}} templates in request
    resolvedReq := c.templateEngine.ResolveMap(req.Step.Request, tplCtx)

    // 3. Build and execute HTTP request (same as existing trigger execution)
    httpReq := buildHTTPRequest(resolvedReq)
    resp, err := c.httpExecutor.ExecuteWithRetry(ctx, httpReq, req.Step.Retry)

    // 4. Parse response using manifest response config
    output, success := parseResponse(resp, req.Step.Response)

    // 5. Publish callback to WORKFLOW-RESUME
    publishCallback(req.CallbackSubject, &CallbackResult{
        InstanceID: req.InstanceID,
        NodeID:     req.NodeID,
        Status:     statusFromSuccess(success),
        Data:       output,
    })
}
```

### 8.7 Plugin JS Consumer (in JS Executor)

New NATS consumer in the JS Workflow Executor that handles preScript/postScript hooks.
Reuses the existing Piscina V8 worker pool.

```typescript
// In JS Workflow Executor — NOT in Workflow Service
async function handlePluginScript(input: PluginScriptRequest): Promise<void> {
    const { instanceId, nodeId, callbackSubject, step, context } = input;

    try {
        // Run script in V8 sandbox (same Piscina pool as code nodes)
        const result = await scriptEngine.runScript({
            script: step.script,
            context: {
                // Previous step results available as `pipeline`
                pipeline: context.pipelineResults,
                // Same globals as code node
                config: context.config,
                state: context.state,
                event: context.event,
                credentials: context.credentials,
                nodes: context.nodes,
                // For postScript: HTTP response from previous step
                response: context.pipelineResults?.[String(context.currentStep - 1)],
                // For preScript: raw request template to modify
                request: step.action === 'preScript' ? context.request : undefined,
            },
            timeout: step.timeout || 10000,
        });

        await publishCallback(callbackSubject, {
            instanceId,
            nodeId,
            status: 'success',
            data: result,
        });

    } catch (error) {
        await publishCallback(callbackSubject, {
            instanceId,
            nodeId,
            status: 'error',
            data: { code: 'PLUGIN_SCRIPT_ERROR', message: error.message },
        });
    }
}
```

### 8.8 preScript Hook — Context & Purpose

**Purpose**: Modify the HTTP request BEFORE template resolution and execution. Can compute
dynamic auth (HMAC, JWT), add conditional fields, build complex request bodies.

```typescript
// Available in V8 sandbox:
request     // Raw request template from manifest
config      // User-configured node values
state       // Workflow state variables
event       // Trigger event payload
credentials // Decrypted credential values
nodes       // Previous node outputs
pipeline    // Results from previous pipeline steps (empty for first step)

// Script example: compute HMAC-SHA256 signature (AWS SigV4 style)
const timestamp = new Date().toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z';
const payload = JSON.stringify(request.body);
const signature = hmacSHA256(credentials.secretKey, timestamp + payload);

request.headers['X-Signature'] = signature;
request.headers['X-Timestamp'] = timestamp;
return request;

// Output: modified request object → used by HTTP step
```

### 8.9 postScript Hook — Context & Purpose

**Purpose**: Transform the HTTP response BEFORE it becomes the node output. Can reshape
data, extract nested fields, update workflow state.

```typescript
// Available in V8 sandbox:
response    // HTTP response from previous step: { status, data, headers }
config      // User-configured node values
state       // Workflow state variables
event       // Trigger event payload
credentials // Decrypted credential values
nodes       // Previous node outputs
pipeline    // Results from ALL previous pipeline steps

// Script example: extract and reshape Slack messages
const items = response.data.messages.map(m => ({
    id: m.ts,
    text: m.text,
    author: m.user
}));

return {
    output: { items, count: items.length, hasMore: response.data.has_more },
    statePatch: { lastCursor: response.data.response_metadata?.next_cursor }
};

// output → saved as nodeOutputs[nodeId]
// statePatch → merged into workflow instance state
```

---

## 9. Backend Implementation

### 9.1 New Module: `plugins` (in Workflow Service)

```
src/modules/plugins/
├── domain/
│   └── entities/
│       ├── manifest.go           # IntegrationManifest struct
│       ├── credential.go         # Credential struct
│       └── pipeline.go           # PipelineStep, Pipeline structs
├── application/
│   ├── ports/
│   │   ├── manifest_repository.go
│   │   └── credential_store.go
│   ├── services/
│   │   └── plugin_service.go     # CRUD for manifests
│   └── di/
│       └── container.go
├── infrastructure/
│   ├── repositories/
│   │   └── manifest_mongo.go     # MongoDB repository
│   └── cache/
│       └── manifest_loader.go    # TieredCache loader (like definitions)
└── interfaces/
    └── http/
        ├── routes/
        │   └── plugin_routes.go
        └── handlers/
            ├── plugin_handler.go     # CRUD endpoints
            └── options_handler.go    # Dynamic dropdown proxy
```

### 9.2 ExecutorRegistry Enhancement

Add a fallback for unknown node types (integration plugins):

```go
func (r *ExecutorRegistry) Get(nodeType string) (NodeExecutor, error) {
    // 1. Exact match (core executors)
    if executor, ok := r.executors[nodeType]; ok {
        return executor, nil
    }

    // 2. Fallback to IntegrationExecutor for plugin nodes
    if r.fallback != nil {
        return r.fallback, nil
    }

    return nil, fmt.Errorf("%w: %s", ErrExecutorNotFound, nodeType)
}
```

### 9.3 Template Engine

Resolves `{{...}}` expressions in manifest strings. Used by the Triggers Service
when executing HTTP steps.

```go
type TemplateEngine struct{}

func (t *TemplateEngine) Resolve(template string, ctx *TemplateContext) (string, error) {
    // ctx contains: config, state, event, credentials, nodes, pipeline
    // Regex: \{\{([^}]+)\}\}
    // Split prefix: config.channel → ctx.Config["channel"]
    // Deep path: event.data.user.name → mapget.Get(ctx.Event, "data.user.name")
    // Pipeline: pipeline.0.signature → ctx.Pipeline["0"]["signature"]
}

type TemplateContext struct {
    Config      map[string]interface{}
    State       map[string]interface{}
    Event       map[string]interface{}
    Credentials map[string]interface{}
    Nodes       map[string]interface{}
    Pipeline    map[string]interface{}  // Results from previous steps
}
```

### 9.4 Credential Store (Implemented)

Credentials are stored as **first-class entities** in the `credentials` collection with
**envelope encryption** (AES-256-GCM). Each credential gets a unique DEK (Data Encryption Key)
encrypted by a Master Key (`CREDENTIAL_MASTER_KEY` env var).

**Entity**: `modules/credentials/domain/entities/credential.go`
```go
type Credential struct {
    ID             model.ObjectId  `bson:"_id,omitempty"`
    Name           string          `bson:"name"`
    PluginId       string          `bson:"pluginId"`
    CredentialType string          `bson:"credentialType"`
    OrgId          *model.ObjectId `bson:"orgId,omitempty"`
    EncryptedDEK   []byte          `bson:"encryptedDEK"` // Never exposed via API (json:"-")
    DEKNonce       []byte          `bson:"dekNonce"`
    EncryptedData  []byte          `bson:"encryptedData"`
    DataNonce      []byte          `bson:"dataNonce"`
    Created        time.Time       `bson:"created"`
    Updated        time.Time       `bson:"updated"`
}
```

**Service port**: `modules/credentials/application/ports/credential_service_port.go`
```go
type CredentialServicePort interface {
    CreateCredential(ctx, requestContext, dto) (*CredentialResponse, error)
    GetCredentialById(ctx, id) (*CredentialResponse, error)
    UpdateCredentialById(ctx, id, dto) (*CredentialResponse, error)
    DeleteCredentialById(ctx, id) (map[string]bool, error)
    GetCredentials(ctx, requestContext, query) (*PaginatedResult[CredentialResponse], error)
    TestCredential(ctx, id) (map[string]bool, error)
    GetCredentialSchema(ctx, pluginId) (*CredentialSchemaResponse, error)
    LoadOptions(ctx, credentialId, resourceKey, dependsOn) ([]LoadOptionsItem, error)
    DecryptCredential(ctx, id) (map[string]interface{}, error) // Internal: runtime only
}
```

**Encryption package**: `packages/utils/envelope/` — reusable AES-256-GCM envelope encryption.

**HTTP routes** (`/api/v1/credentials`):
- `GET /` — List credentials (paginated, multi-tenant)
- `POST /` — Create credential (encrypt + persist)
- `GET /schema/:pluginId` — Get credential schema from plugin manifest (UI form rendering)
- `GET /:id` — Get credential by ID (never returns encrypted data)
- `PATCH /:id` — Update credential (re-encrypts if data changed)
- `DELETE /:id` — Delete credential
- `POST /:id/test` — Test credential against plugin's test endpoint
- `POST /:id/load_options/:resourceKey` — Proxy loadOptions request (decrypt → HTTP → extract/transform)

### 9.5 Dynamic Dropdown Proxy Endpoint (Implemented)

For `loadOptions` properties, the backend proxies the HTTP request to avoid exposing
credentials in the browser. The endpoint lives in the **credentials module** because it
requires credential decryption.

```
POST /api/v1/credentials/:credentialId/load_options/:resourceKey
Body (optional): { "dependsOn": { "baseId": "appXXX", "sheetName": "Sheet1" } }
Response: [{ "label": "General Chat", "value": "-100123" }, ...]
```

**Implementation** (`modules/credentials/application/services/credential_service.go`):
```go
func (s *CredentialService) LoadOptions(ctx, credentialId, resourceKey string, dependsOn map[string]string) ([]LoadOptionsItem, error) {
    // 1. Decrypt credential → get plaintext data (e.g., {"botToken": "123:ABC"})
    cred := s.deps.CredentialRepo.FindById(ctx, &credentialId)
    credData := s.decryptEntity(cred)

    // 2. Read manifest by credential.pluginId → find loadOptions[resourceKey]
    manifest := s.deps.PluginRepo.FindByPluginId(ctx, cred.PluginId)
    loadOptDef := manifest.LoadOptions[resourceKey]

    // 3. Resolve templates in request.path
    //    {{credentials.botToken}} → from decrypted data
    //    {{dependsOn.baseId}}     → from request body
    path := resolveTemplates(loadOptDef.Request.Path, credData, dependsOn)

    // 4. HTTP call to baseUrl + resolvedPath
    resp := httpCall(manifest.BaseUrl + path, loadOptDef.Request.Method)

    // 5a. If transform exists → run JS ES5 via goja (pure Go, sandboxed, 5s timeout)
    if loadOptDef.Transform != "" {
        return jsrunner.RunTransform(ctx, loadOptDef.Transform, resp)
    }

    // 5b. Simple mode → extract via dataPath + valuePath + labelPath
    return extractLoadOptions(resp, loadOptDef.DataPath, loadOptDef.ValuePath, loadOptDef.LabelPath)
}
```

**JS Transform** (`packages/utils/jsrunner/`) — for complex transformations (filtering, concatenation, flattening):
```js
// Example: OpenAI models — filter and sort
function transform(data) {
  return data.data
    .filter(function(m) { return m.id.indexOf("gpt") === 0; })
    .map(function(m) { return { label: m.id, value: m.id }; });
}
```
Runs in goja (pure Go ES5.1 runtime, no CGO). Sandboxed: no FS, no network, 5s timeout.

---

## 10. Frontend Integration

### 10.1 What Already Works (No Changes Needed)

| Feature | How It Works | Location |
|---------|-------------|----------|
| Plugin registration | `pluginRegistry.registerPlugin()` | `stores/pluginRegistry/actions.ts` |
| Catalog grouping | `catalog` getter groups by category | `stores/pluginRegistry/getters.ts` |
| Canvas rendering | `GenericWorkflowNode` resolves from registry | `workflow-plugin-core/nodes/_shared/` |
| Auto-form | `DynamicNodeForm` renders from `properties[]` | `createEditWorkflowPage/components/` |
| Handle resolution | `resolveNodeHandles()` in SDK | `workflow-sdk/` |
| Config panel | Falls back to DynamicNodeForm when no configComponent | `NodeConfigPanel.vue` |

### 10.2 What Needs To Be Added

1. **Plugin loader service** — fetch manifests from API on editor boot:

```typescript
async function loadIntegrationPlugins(registerFn: (plugin: WorkflowPlugin) => void) {
    const manifests = await pluginApi.listEnabled();
    for (const manifest of manifests) {
        const plugin = manifestToPlugin(manifest);
        registerFn(plugin);
    }
}
```

2. **Manifest-to-Plugin converter** — maps JSON manifest to `WorkflowPlugin`:

```typescript
function manifestToPlugin(manifest: IntegrationManifest): WorkflowPlugin {
    return {
        id: manifest.pluginId,
        name: manifest.name,
        version: manifest.version,
        category: manifest.category as PluginCategory,
        icon: manifest.icon,
        nodeTypes: manifest.nodeTypes.map(nt => ({
            type: nt.type,
            label: nt.label,
            icon: nt.icon || manifest.icon,
            color: nt.color || manifest.color,
            description: nt.description,
            inputs: nt.inputs,
            outputs: nt.outputs,
            configSchema: {},
            properties: nt.properties,          // Direct 1:1 mapping
            defaults: nt.defaults,
            outputHints: nt.outputHints,
            // No canvasComponent → GenericWorkflowNode used
            // No configComponent → DynamicNodeForm used
        })),
    };
}
```

3. **Plugin management page** — admin UI for enabling/disabling/configuring integration plugins
4. **Credential management UI** — secure forms for entering API keys, OAuth flows
5. **DynamicNodeForm extensions** — fixedCollection, notice, dynamic dropdown support

---

## 11. Credential System

> **DECISION**: Credentials are **first-class entities** stored separately from plugin manifests.
> This replaced the previous `auth.inject` pattern.

### 11.1 Why Credentials Are Separate

- Same credential (e.g., Slack OAuth token) shared across multiple workflows
- User configures credential ONCE, all workflows share it
- Credential defines HOW it injects into HTTP requests (headers, query, body)
- Credential has a test endpoint for validation
- OAuth2 credentials support inheritance (base type + per-provider overrides)

### 11.2 Credential Type Schema

```json
{
  "_id": "ObjectId",
  "name": "telegramApi",
  "displayName": "Telegram API",
  "pluginId": "telegram",

  "properties": [
    {
      "name": "botToken",
      "displayName": "Bot Token",
      "type": "string",
      "secret": true,
      "required": true
    },
    {
      "name": "baseUrl",
      "displayName": "Base URL",
      "type": "string",
      "default": "https://api.telegram.org"
    }
  ],

  "authenticate": {
    "type": "generic",
    "properties": {
      "headers": {},
      "qs": {},
      "body": {}
    }
  },

  "test": {
    "request": {
      "baseURL": "={{credentials.baseUrl}}/bot{{credentials.botToken}}",
      "url": "/getMe",
      "method": "GET"
    }
  }
}
```

### 11.3 Four Injection Targets

The `authenticate.properties` object supports 4 injection targets:

| Target | HTTP Mapping | Use Case |
|--------|-------------|----------|
| `headers` | HTTP request headers | Bearer tokens, API keys, custom auth |
| `qs` | URL query string parameters | API keys in URL |
| `body` | Request body fields | SOAP auth, form-based auth |
| `auth` | HTTP Basic Authentication | Legacy APIs |

```json
{
  "authenticate": {
    "type": "generic",
    "properties": {
      "headers": {
        "Authorization": "=Bearer {{credentials.accessToken}}"
      }
    }
  }
}
```

### 11.4 Real Examples

**Telegram** (token in URL path — no header injection needed):
```json
{
  "name": "telegramApi",
  "properties": [
    { "name": "botToken", "type": "string", "secret": true, "required": true },
    { "name": "baseUrl", "type": "string", "default": "https://api.telegram.org" }
  ],
  "authenticate": { "type": "generic", "properties": {} },
  "test": {
    "request": {
      "baseURL": "={{credentials.baseUrl}}/bot{{credentials.botToken}}",
      "url": "/getMe"
    }
  }
}
```

**Slack** (Bearer header):
```json
{
  "name": "slackApi",
  "properties": [
    { "name": "accessToken", "type": "string", "secret": true, "required": true }
  ],
  "authenticate": {
    "type": "generic",
    "properties": {
      "headers": { "Authorization": "=Bearer {{credentials.accessToken}}" }
    }
  },
  "test": {
    "request": { "baseURL": "https://slack.com", "url": "/api/users.profile.get" },
    "rules": [
      { "type": "responseBody", "key": "error", "value": "invalid_auth", "message": "Invalid token" }
    ]
  }
}
```

**MOCO** (multi-tenant with subdomain):
```json
{
  "name": "mocoApi",
  "properties": [
    { "name": "subDomain", "displayName": "Sub-Domain", "type": "string", "required": true },
    { "name": "apiKey", "displayName": "API Key", "type": "string", "secret": true, "required": true }
  ],
  "authenticate": {
    "type": "generic",
    "properties": {
      "headers": { "Authorization": "=Token token={{credentials.apiKey}}" }
    }
  },
  "test": {
    "request": {
      "baseURL": "=https://{{credentials.subDomain}}.mocoapp.com/api/v1",
      "url": "/session"
    }
  }
}
```

### 11.5 OAuth2 Inheritance

OAuth2 credentials extend a base type, only overriding URLs and scopes:

```json
{
  "name": "slackOAuth2",
  "displayName": "Slack OAuth2",
  "extends": "oAuth2Api",
  "properties": [
    { "name": "authUrl", "type": "hidden", "default": "https://slack.com/oauth/v2/authorize" },
    { "name": "accessTokenUrl", "type": "hidden", "default": "https://slack.com/api/oauth.v2.access" },
    { "name": "scope", "type": "hidden", "default": "chat:write channels:read" },
    { "name": "authentication", "type": "hidden", "default": "body" }
  ]
}
```

The base `oAuth2Api` credential type provides: Client ID, Client Secret fields,
grant type selection, token refresh logic, token storage.

### 11.6 Dual Auth Per Node

A node can support multiple credential types (e.g., API key AND OAuth2):

```json
{
  "type": "github/create_issue",
  "credentials": [
    {
      "type": "githubApi",
      "required": true,
      "displayOptions": { "show": { "authentication": ["accessToken"] } }
    },
    {
      "type": "githubOAuth2",
      "required": true,
      "displayOptions": { "show": { "authentication": ["oAuth2"] } }
    }
  ],
  "properties": [
    {
      "name": "authentication",
      "displayName": "Authentication",
      "type": "options",
      "options": [
        { "label": "Access Token", "value": "accessToken" },
        { "label": "OAuth2", "value": "oAuth2" }
      ],
      "default": "accessToken"
    }
  ]
}
```

### 11.7 Credential Instance Storage

User-configured credential values in `plugins_credentials` collection:

```json
{
  "_id": "ObjectId",
  "credentialType": "telegramApi",
  "name": "My Telegram Bot",
  "orgId": "ObjectId",
  "data": "<AES-256-GCM encrypted JSON>",
  "createdAt": "...",
  "updatedAt": "..."
}
```

### 11.8 How Credentials Flow in the Pipeline

1. IntegrationExecutor looks up which credential type the manifest requires
2. Gets the user's saved credential instance for that type + org
3. Decrypts the values
4. Stores decrypted values in pipeline NodeState under `credentials`
5. Triggers Service uses `{{credentials.*}}` in template resolution
6. JS Executor receives `credentials` as context variable in V8 sandbox

---

## 12. Webhook & Trigger Architecture

> **DECISION (ATUALIZADO)**: O setup de webhook é GUIADO (wizard dentro do trigger node),
> não escondido via Saga. O user vê e entende cada entidade (DataSource, RouteGroup) mas
> o sistema pré-preenche defaults inteligentes. Zero magia oculta.
>
> **Documento completo**: [`plugin-webhook-trigger.md`](./plugin-webhook-trigger.md)

### 12.1 Resumo da Arquitetura

Plugin webhooks usam a pipeline existente: DataSource → JS Executor → Router → Workflow.
Novo `kind: "workflow"` no Router com dois modos: `newInstance` e `runningInstance`.

### 12.2 Dois Ciclos de Vida Independentes

- **Credential controla DataSource** (1:1): cria/remove DataSource quando credential é criada/removida
- **Workflow controla RouteGroup** (1:1): cria/remove RouteGroup quando workflow é ativado/desativado

### 12.3 O Wizard (3 steps dentro do trigger node)

1. **Data Source**: seleciona existente ou cria novo (linked à credential)
2. **Routing**: events, mode (new/running instance), correlation, filters
3. **Review**: mostra tudo que será criado, user confirma

### 12.4 Signal Delivery (Temporal-inspired)

Para `runningInstance` mode, o Workflow Service usa correlation index no NATS KV
para resolver instância ativa e entregar signal. Se não encontra, aplica `onMiss`
policy (create, drop, queue).

### 12.5 Entidades Visíveis com `managedBy`

Entidades criadas via wizard são VISÍVEIS (não escondidas). Usam `managedBy` ref
em vez de `createdBySystem: true`. User pode ver, editar (com warning), deletar
(com confirmação).

```go
type ManagedByRef struct {
    WorkflowID string `bson:"workflowId" json:"workflowId"`
    NodeID     string `bson:"nodeId"     json:"nodeId"`
    PluginID   string `bson:"pluginId"   json:"pluginId"`
}
```

### 12.6 Workflow Service NEVER Makes HTTP Requests

Regra absoluta: o workflow service só orquestra o DAG. Todo I/O externo via NATS:

| Action | Who Executes | Via |
|--------|-------------|-----|
| Plugin HTTP execution | **Triggers Service** | NATS: `WORKFLOW-PLUGIN-HTTP` |
| Plugin JS hooks | **JS Workflow Executor** | NATS: `WORKFLOW-JS-CODE` |
| Webhook registration | **Triggers Service** | NATS (async) |
| Credential test | **Triggers Service** | NATS (async) |

> Para a documentação completa com fluxos end-to-end, cenários multi-workflow,
> manifest triggerTypes schema, signal delivery implementation e escala:
> **[`plugin-webhook-trigger.md`](./plugin-webhook-trigger.md)**

---

## 14. Storage & Caching Strategy

### 14.1 MongoDB Collections

```
plugins_manifests         # Integration manifest documents
├── _id: ObjectId
├── pluginId: string (unique index)
├── name, version, category, icon, color, description
├── credential: { type, required }
├── nodeTypes: NodeTypeManifest[]
├── enabled: boolean
├── createdAt, updatedAt
└── Index: { pluginId: 1 } unique, { enabled: 1 }

plugins_credential_types  # Credential type definitions
├── _id: ObjectId
├── name: string (unique index)
├── displayName, pluginId
├── properties: CredentialProperty[]
├── authenticate: AuthenticateConfig
├── test: TestConfig
├── extends: string (optional)
└── Index: { name: 1 } unique, { pluginId: 1 }

plugins_credentials       # Encrypted credential instances per org
├── _id: ObjectId
├── credentialType: string
├── name: string
├── orgId: ObjectId
├── data: binary (AES-256-GCM encrypted JSON)
├── createdAt, updatedAt
└── Index: { credentialType: 1, orgId: 1 }
```

### 14.2 TieredCache

Integration manifests follow the same TieredCache pattern as workflow definitions:

```
L0 (RAM)  — All enabled manifests loaded at boot (~2.5MB for 500 plugins)
L1 (Disk) — Fallback for cold starts
L2 (MinIO) — Source of truth for manifests (optional, MongoDB can be SoT)
```

**Cache invalidation**: NATS fanout on manifest CRUD (same as definition cache invalidation).

---

## 15. What Needs To Be Built

### 15.1 Workflow Service (Go)

| Component | Effort | Priority |
|-----------|--------|----------|
| `plugins` module (CRUD + MongoDB repo) | Medium | P0 |
| `IntegrationManifest` entity + `Pipeline` structs | Small | P0 |
| `IntegrationExecutor` (builds pipeline, suspends DAG) | Medium | P0 |
| `manifestIndex` (RAM map, boot from MongoDB) | Small | P0 |
| `ExecutorRegistry` fallback match | Small | P0 |
| Pipeline resume handler in RuntimeService | Medium | P0 |
| `dispatchPipelineStep()` routing logic | Small | P0 |
| `DispatchPluginHTTP` / `DispatchPluginJS` in RuntimePublisher | Small | P0 |
| `lifecycle.go` — new `plugin_pipeline` case | Small | P0 |
| `ManifestLoader` (TieredCache) | Small | P1 |
| `CredentialStore` (encrypted MongoDB) | Medium | P1 |
| Dynamic dropdown proxy endpoint | Small | P1 |
| Manifest validation (JSON Schema) | Small | P1 |

### 15.2 Triggers Service (Go)

| Component | Effort | Priority |
|-----------|--------|----------|
| `PluginHTTPConsumer` (NATS consumer for `WORKFLOW-PLUGIN-HTTP`) | Medium | P0 |
| Template resolution with `{{...}}` and pipeline context | Small | P0 |
| Response parsing (successPath, dataPath, errorPath) | Small | P0 |
| Retry logic (429/5xx backoff, configurable per manifest) | Medium | P1 |
| Credential test endpoint consumer | Small | P1 |
| Webhook register/unregister consumer | Small | P1 |

### 15.3 JS Workflow Executor (Node.js)

| Component | Effort | Priority |
|-----------|--------|----------|
| `PluginScriptConsumer` (NATS consumer for plugin JS steps) | Medium | P0 |
| Hook runner with pipeline context (preScript/postScript) | Small | P0 |
| Template resolver for JS side | Small | P1 |

### 15.4 Frontend (TypeScript/Vue)

| Component | Effort | Priority |
|-----------|--------|----------|
| Plugin API wrapper (`@mapexos/apis`) | Small | P0 |
| Plugin Zod schemas (`@mapexos/schemas`) | Small | P0 |
| `manifestToPlugin()` converter | Small | P0 |
| Plugin loader (fetch on editor boot) | Small | P0 |
| Plugin management admin page | Large | P1 |
| Credential management UI | Medium | P1 |
| `DynamicNodeForm` extensions (fixedCollection, notice, dynamic dropdown) | Medium | P1 |

### 15.5 Existing — No Changes Needed

| Component | Why |
|-----------|-----|
| `pluginRegistry` store | Already supports dynamic plugin registration |
| `GenericWorkflowNode` | Works for any node type |
| `DynamicNodeForm` | Renders from `properties[]` — perfect for manifests |
| `BaseWorkflowNode` | Handles all visual rendering |
| Core 17 executors | Untouched — only IntegrationExecutor added |
| DAG execution engine | Already dispatches by `node.Type` — just needs registry fallback |
| NATS KV checkpoint | Already used by all async nodes — pipeline reuses it |
| Archiver | Already consumes `workflow.state.*` events — pipeline steps generate them |
| Reconciler | Already handles timeouts — pipeline steps have timeouts |

---

## 16. Migration Path

### Phase 1: Foundation (MVP) — Single-Step Pipelines

1. Create `plugins` module with MongoDB CRUD
2. Build `IntegrationExecutor` (pipeline builder, suspends DAG)
3. Build `manifestIndex` (boot from MongoDB, RAM map)
4. Add fallback to `ExecutorRegistry`
5. Build pipeline resume handler in RuntimeService
6. Build `PluginHTTPConsumer` in **Triggers Service** (generic HTTP execution)
7. Build `TemplateEngine` with credential + pipeline context
8. Add `plugin_pipeline` dispatch case in `lifecycle.go`
9. Create NATS stream: `WORKFLOW-PLUGIN-HTTP`
10. Frontend: plugin API + loader + `manifestToPlugin()`
11. Create 3 sample manifests: Telegram, Slack, generic HTTP

**Result**: Single-step pipelines (JSON manifest → Triggers Service HTTP) work end-to-end.

### Phase 2: Multi-Step Pipelines (preScript/postScript)

1. Build `PluginScriptConsumer` in **JS Workflow Executor**
2. Hook runner with pipeline context (preScript/postScript)
3. Extend pipeline resume handler to chain steps with KV checkpoint
4. Credential store with encryption
5. Create sample manifest with hooks: Salesforce (preScript + HTTP + postScript)

**Result**: Multi-step pipelines with JS hooks work end-to-end.

### Phase 3: Polish

1. Plugin management admin page
2. Credential management UI + OAuth2 flow support
3. DynamicNodeForm extensions (fixedCollection, notice, dynamic dropdown)
4. Dynamic dropdown proxy endpoint
5. Cache invalidation via NATS fanout
6. Retry + circuit breaker in PluginHTTPConsumer

**Result**: Production-ready plugin system with credential management.

### Phase 4: Advanced

1. Saga workflows for webhook provisioning
2. System-managed resources (`createdBySystem` flag)
3. Plugin marketplace concept (browse + install)
4. Plugin versioning + migration
5. Per-tenant plugin enable/disable
6. Rate limiting per integration
7. AI Service pipeline step (future)

---

## 17. File Reference

### Workflow Service (Go) — Existing

| File | Purpose |
|------|---------|
| `runtime/domain/entities/executor_port.go` | `NodeExecutor` interface |
| `runtime/domain/executors/executor.go` | `ExecutorRegistry` — add fallback |
| `runtime/domain/executors/registry_builder.go` | `BuildRegistry` — register IntegrationExecutor |
| `runtime/application/services/runtime_service.go` | DAG execution loop + pipeline resume handler |
| `runtime/application/services/lifecycle.go` | `dispatchByNodeType()` — add `plugin_pipeline` |
| `runtime/application/services/dag_step.go` | `executeStep()` — executor dispatch per node |
| `runtime/domain/services/config_parsing.go` | `parseNodeConfig()` — add default case for plugin config |
| `runtime/infrastructure/cache/definition_loader.go` | TieredCache pattern to reuse for manifests |
| `runtime/infrastructure/messaging/nats/runtime_publisher.go` | Add `DispatchPluginHTTP()`, `DispatchPluginJS()` |

### Triggers Service (Go) — Existing Patterns to Reuse

| File | Purpose |
|------|---------|
| `events/application/handlers/placeholder_resolver.go` | `ResolvePlaceholders()` — reuse in PluginHTTPConsumer |
| `events/infrastructure/technical/http/http_executor.go` | HTTP request building — reuse in PluginHTTPConsumer |
| `events/infrastructure/registry/executor_registry.go` | ExecutorRegistry factory — reference pattern |
| `events/application/services/event_service.go` | Batch processing pattern — reuse for plugin consumer |

### JS Workflow Executor (Node.js) — Existing

| File | Purpose |
|------|---------|
| `modules/events/interfaces/message/consumers/WorkflowCodeConsumer/` | Reference for PluginScriptConsumer |
| `modules/scripts/application/services/workflow-script.service.ts` | Script execution — reuse Piscina dispatch |
| `modules/engine/application/services/script-engine.service.ts` | V8 worker pool — runs hook scripts |
| `modules/engine/infrastructure/worker/piscina-worker.ts` | Worker thread — V8 isolate execution |

### Frontend (TypeScript)

| File | Purpose |
|------|---------|
| `workflow-sdk/src/interfaces/workflowPlugin.interface.ts` | `WorkflowPlugin`, `PluginNodeType`, `NodePropertyDefinition` |
| `workflow-plugin-core/src/constants/corePlugins.constant.ts` | 6 core plugins with 18 node types |
| `workflow-plugin-core/src/nodes/_shared/GenericWorkflowNode.vue` | Default canvas component |
| `stores/pluginRegistry/actions.ts` | `registerPlugin()`, `unregisterPlugin()`, `getNodeType()` |
| `stores/pluginRegistry/getters.ts` | `catalog` getter — groups by category |
| `createEditWorkflowPage/components/DynamicNodeForm/` | Auto-form from `properties[]` |
| `createEditWorkflowPage/components/NodeConfigPanel.vue` | Config panel — DynamicNodeForm fallback |

### n8n Study Repos (Reference Only)

| Repo | Location | Purpose |
|------|----------|---------|
| awesome-n8n | `/home/thiago/Documents/Projects/MAPEX/General_Tests/n8n_study/awesome-n8n/` | Community node index, download rankings |
| n8n-nodes-starter | `/home/thiago/Documents/Projects/MAPEX/General_Tests/n8n_study/n8n-nodes-starter/` | Official template, INodeType, credentials |
| n8n-nodes (skriptfabrik) | `/home/thiago/Documents/Projects/MAPEX/General_Tests/n8n_study/n8n-nodes/` | 11 packages: Channable, Clockify, Google Enhanced, etc. |
