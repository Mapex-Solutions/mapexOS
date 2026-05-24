# Journey: Slack trigger

## What this journey proves

The Slack trigger of the triggers service fires through the live
pipeline and POSTs its webhook payload to a real HTTP listener
(same shape as a Slack incoming webhook would receive).

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger POSTs webhook. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetry → gateway → js-executor → router → trigger POSTs webhook. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-slack
```

## Requirements

- Live stack with default ports.
- Free port `11010` on the host (shared HTTP sink).
