# Workflow Dispatcher

Simulates a Router event by publishing directly to the `WORKFLOW-EXECUTION` stream. Useful for testing the workflow service without needing the full pipeline (HTTP Gateway → JS Executor → Router).

## Setup

```bash
cd scripts/dispatcher
npm install nats
```

## Usage

```bash
# Send default payload with auto-generated workflowUUID
node index.js

# Send with custom workflowUUID
node index.js --uuid my-custom-uuid

# Send a different payload file
node index.js --file custom-payload.json
```

## What it does

1. Loads `payload.json` (or custom file)
2. Generates a unique `workflowUUID`, `eventTrackerId`, and `executionId`
3. Connects to NATS with service credentials
4. Publishes to subject `workflow.execution.router` (WORKFLOW-EXECUTION stream)
5. The workflow service HandleExecution consumer picks it up

## Payload format

The payload follows the `WorkflowExecutionMessage` contract:

```json
{
  "mode": "newInstance",
  "data": {
    "instanceId": "...",
    "workflowUUID": "auto-generated"
  },
  "event": { ... },
  "orgId": "...",
  "pathKey": "..."
}
```

## Modes

| Mode | Description | Required data fields |
|------|-------------|---------------------|
| `newInstance` | Create a new workflow execution | `instanceId`, `workflowUUID` (auto) |
| `signal` | Send signal to a waiting execution | `workflowUUID`, `signalName`, `signalData` |
| `signalOrStart` | Try signal first, fallback to newInstance | `instanceId`, `workflowUUID`, `signalName` |

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_USER` | `service` | NATS username |
| `NATS_PASSWORD` | `service_secret` | NATS password |
