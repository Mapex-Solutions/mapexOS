# Journey: WebSocket trigger

## What this journey proves

WebSocket trigger connects + writes to an in-process WS server via
the live pipeline; verified via the events_trigger oracle.

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger dials WS sink. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetry → gateway → js-executor → router → trigger dials WS sink. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-websocket
```

## Requirements

- Live stack with default ports.
- Free port `11026` on the host (override `SAGA_TRIGGER_WS_BIND_ADDR`).
