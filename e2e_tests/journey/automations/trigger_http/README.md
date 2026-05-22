# Journey: HTTP trigger

## What this journey proves

The HTTP trigger of the triggers service fires through the live
pipeline and POSTs to a real HTTP listener.

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger POSTs to the sink. |
| `phase2_event_pipeline` *(planned)* | POST telemetry to gateway → js-executor → router → trigger POSTs to the sink. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-http
```

## Requirements

- Live stack with default ports.
- Free port `11010` on the host (override `SAGA_TRIGGER_SINK_BIND_ADDR`).
