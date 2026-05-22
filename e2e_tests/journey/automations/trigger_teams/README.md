# Journey: Teams trigger

## What this journey proves

The Teams trigger of the triggers service fires through the live
pipeline and POSTs its webhook MessageCard payload to a real HTTP
listener (same shape as a Teams incoming webhook would receive).

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger POSTs webhook. |
| `phase2_event_pipeline` *(planned)* | POST telemetry → gateway → js-executor → router → trigger POSTs webhook. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-teams
```

## Requirements

- Live stack with default ports.
- Free port `11010` on the host (shared HTTP sink).
