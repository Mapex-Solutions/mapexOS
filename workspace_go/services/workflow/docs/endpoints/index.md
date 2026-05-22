# Endpoints

## HTTP API

### Workflow Definitions

**Base:** `/api/v1/workflows`

| Method | Path | Permission | Description |
|--------|------|-----------|-------------|
| GET | `/` | `workflows.list` | List definitions (paginated + filtered) |
| POST | `/` | `workflows.create` | Create definition |
| GET | `/:workflowId` | `workflows.read` | Get definition by ID |
| PATCH | `/:workflowId` | `workflows.update` | Update definition |
| DELETE | `/:workflowId` | `workflows.delete` | Delete definition |

**Auth:** JWT or OAuth2 (configurable via `AUTH_STRATEGY`)

**Middleware stack:** Auth → Coverage → Permission

**Query parameters (GET /):**
- `page` (int): Page number
- `limit` (int): Items per page
- Standard filters by orgId, enabled, name

### Workflow Instances

**Base:** `/api/v1/workflow-instances`

| Method | Path | Permission | Description |
|--------|------|-----------|-------------|
| GET | `/` | `workflows.instances.list` | List instances (paginated + filtered) |
| GET | `/:instanceId` | `workflows.instances.read` | Get instance by ID (hybrid KV/MongoDB) |
| POST | `/:instanceId/cancel` | `workflows.instances.cancel` | Cancel a running instance |
| POST | `/:instanceId/signal` | `workflows.instances.signal` | Send signal to waiting instance |

**GetInstanceById logic:**
- Terminal status (completed/failed/cancelled) → returns from MongoDB (full archive)
- Non-terminal status (running/waiting) → returns from NATS KV (fresh state)

### Observability

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | none | Health check (MongoDB, Redis, NATS, MinIO) |
| GET | `/metrics` | none | Prometheus metrics |

## NATS Streams

### Inbound (consumed by this service)

| Stream | Subject | Consumer Type | Module | Purpose |
|--------|---------|--------------|--------|---------|
| `WORKFLOW-TRIGGER` | `workflow.trigger.>` | Single message | Runtime | Start new workflow instances |
| `WORKFLOW-RESUME` | `workflow.resume.>` | Single message | Runtime | Resume suspended instances |
| `WORKFLOW-STATE` | `workflow.state.>` | Batch (500 msgs) | Archiver | Persist state lifecycle to MongoDB |
| `WORKFLOW-RECONCILER` | `workflow.reconciler.>` | Single message | Reconciler | Timer registrations |
| `WORKFLOW-SIGNAL` | `workflow.signal.>` | Single message | Instances | Route external signals to waiting instances |

### Outbound (published by this service)

| Subject | Published by | Purpose |
|---------|-------------|---------|
| `workflow.state.created` | RuntimeService | Notify Archiver of new instance |
| `workflow.state.waiting` | RuntimeService | Notify Archiver of suspension (includes timer info) |
| `workflow.state.resumed` | RuntimeService | Notify Archiver to clear timer |
| `workflow.state.completed` | RuntimeService | Notify Archiver of completion |
| `workflow.state.failed` | RuntimeService | Notify Archiver of failure |
| `workflow.state.cancelled` | InstancesService | Notify Archiver of cancellation |
| `workflow.resume.reenqueue.{id}` | RuntimeService | Re-enqueue after MaxInlineSteps exceeded |
| `workflow.resume.callback.{id}` | RuntimeService | Callback subject for code/subworkflow async |
| `workflow.resume.timer.{id}` | ReconcilerService | Resume expired timer |
| `workflow.resume.timer.{id}` | ArchiverService | Short timer fast-path (< 1min) |
| `workflow.js.code` | RuntimeService | Dispatch code execution to JS Executor |
| `workflow.trigger.subworkflow.{id}` | RuntimeService | Trigger child workflow |
| `workflow.trigger.event` | RuntimeService | Dispatch trigger event to Trigger Service |
| `fanout.workflow.definition.invalidate` | DefinitionService | Cache invalidation across instances |

### Outbound Streams (created by this service, consumed externally)

| Stream | Subject | Consumed by | Purpose |
|--------|---------|------------|---------|
| `WORKFLOW-JS-CODE` | `workflow.js.code` | JS Workflow Executor | Execute JavaScript code nodes |
| `WORKFLOW-LOGS` | `workflow.logs.>` | Archiver (future) | Per-step execution logs |

### NATS KV

| Bucket | Key Format | Purpose |
|--------|-----------|---------|
| `WORKFLOW-INSTANCES` | Instance ID | Hot state during execution (~1-5KB per instance) |
