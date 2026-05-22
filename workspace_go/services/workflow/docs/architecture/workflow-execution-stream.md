# WORKFLOW-EXECUTION Stream

## Overview

The `WORKFLOW-EXECUTION` stream is the **single entry point** for all execution commands in the workflow service. The Router service (and other producers) publish messages to this stream, and the Execution Consumer in the runtime module processes them.

Routing is done by the `mode` field in the payload — **never by subject**. The subject is always a fixed pattern.

## Stream Configuration

| Field | Value |
|---|---|
| Stream | `WORKFLOW-EXECUTION` |
| Subject | `workflow.execution.>` |
| Storage | File |
| Retention | Work (removed after ACK) |

## Consumer

| Field | Value |
|---|---|
| Durable | `{serviceName}-workflow-execution` |
| Queue Group | `{serviceName}-WORKFLOW-EXECUTION-GROUP` |
| Handler | `RuntimeServicePort.HandleExecution()` |
| Processing | Per-message (MessageHandlerV2) |
| Retry | Default exponential backoff |
| DLQ | Yes (ServiceType: "workflow", EventType: "workflow-execution") |

## Message Structure

```json
{
  "mode": "newInstance | signal | signalOrStart",
  "event": {},
  "metadata": {},
  "data": {}
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `mode` | string | Yes | Dispatch mode: `newInstance`, `signal`, `signalOrStart` |
| `event` | object | No | Event payload from the Router (original sensor/trigger data) |
| `metadata` | object | No | Router metadata (routerId, matchRuleId, timestamp, etc.) |
| `data` | object | Yes | Mode-specific configuration (schema defined per mode below) |

## Modes

### mode: `newInstance`

Creates a new execution from an existing instance config.

**When:** Router detects an event that matches a rule pointing to an instance.

**data schema:**

| Field | Type | Required | Description |
|---|---|---|---|
| `instanceId` | string | Yes | MongoDB ObjectID of the instance config (created by UI) |
| `workflowUUID` | string | No | Custom UUID for the execution. If absent, runtime generates GUIDv4 |

**Example:**

```json
{
  "mode": "newInstance",
  "event": { "temperatura": 42, "sensorId": "sensor-001" },
  "metadata": { "routerId": "rule-abc", "receivedAt": "2026-03-20T10:00:00Z" },
  "data": {
    "instanceId": "683f1a2b3c4d5e6f7a8b9c0d",
    "workflowUUID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
}
```

**Runtime flow:**
1. Load instance config from TieredCache (`data.instanceId`)
2. Load definition from TieredCache (`instance.definitionId`)
3. Generate execution UUID (`data.workflowUUID` or GUIDv4)
4. Create `WorkflowExecution` in NATS KV (`exec:{uuid}`)
5. Publish `workflow.state.created` → Archiver inserts stub to MongoDB
6. Run DAG from `__start__`

---

### mode: `signal`

Delivers a signal to a waiting execution. Used when an external event needs to resume a suspended workflow (e.g., user approval, webhook callback).

**When:** Router detects an event that should signal a specific execution.

**data schema:**

| Field | Type | Required | Description |
|---|---|---|---|
| `workflowUUID` | string | Yes | UUID of the execution to signal |
| `signalName` | string | Yes | Name of the signal (must match a `wait_signal` node's configured signalName) |
| `signalData` | object | No | Data payload to inject into the execution state |

**Example:**

```json
{
  "mode": "signal",
  "event": {},
  "metadata": {},
  "data": {
    "workflowUUID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "signalName": "approval",
    "signalData": { "approved": true, "approvedBy": "john@example.com" }
  }
}
```

**Runtime flow:**
1. Read execution from NATS KV (`exec:{data.workflowUUID}`)
2. Validate execution status is `waiting`
3. Find node in `nodeStates` waiting for `data.signalName`
4. Publish resume message to `WORKFLOW-RESUME` with signal data
5. Runtime `HandleResume` picks up and continues DAG from the waiting node

---

### mode: `signalOrStart`

Tries to deliver a signal first. If the execution is not found or not waiting, falls back to creating a new execution (like `newInstance`).

**When:** Router doesn't know if an execution is already running. Common for event-driven workflows that may or may not have started yet.

**data schema:**

| Field | Type | Required | Description |
|---|---|---|---|
| `instanceId` | string | Yes | Instance config ID (used for fallback `newInstance`) |
| `workflowUUID` | string | Yes | UUID to try signal on, or to use for new execution |
| `signalName` | string | Yes | Signal name to try delivering |
| `signalData` | object | No | Data payload for the signal |

**Example:**

```json
{
  "mode": "signalOrStart",
  "event": { "temperatura": 42 },
  "metadata": { "routerId": "rule-xyz" },
  "data": {
    "instanceId": "683f1a2b3c4d5e6f7a8b9c0d",
    "workflowUUID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "signalName": "sensorData",
    "signalData": { "value": 42 }
  }
}
```

**Runtime flow:**
1. Try `signal` flow:
   - Read execution from KV (`exec:{data.workflowUUID}`)
   - If found + waiting for `data.signalName` → deliver signal → done
2. If signal fails (not found, not waiting, wrong signal):
   - Fall back to `newInstance` flow using `data.instanceId`

---

## Producers

| Producer | When | Typical mode |
|---|---|---|
| **Router Service** | Event matches a routing rule | `newInstance`, `signal`, `signalOrStart` |
| **UI / API** | User clicks "Execute" | `newInstance` |
| **External Webhooks** | Callback arrives | `signal` |

## Related Streams

| Stream | Relationship |
|---|---|
| `WORKFLOW-TRIGGER` | Execution Consumer forwards `newInstance` data here (internal) |
| `WORKFLOW-RESUME` | Signal delivery publishes resume messages here |
| `WORKFLOW-STATE` | Runtime publishes lifecycle events here (created, waiting, etc.) |
