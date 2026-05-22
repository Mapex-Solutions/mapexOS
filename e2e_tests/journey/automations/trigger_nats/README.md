# Journey: NATS trigger

## What this journey proves

The NATS trigger of the triggers service publishes to a NATS subject
via the live pipeline; verified via the events_trigger oracle.

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger publishes to subject. |
| `phase2_event_pipeline` *(planned)* | POST telemetry → gateway → js-executor → router → trigger publishes to subject. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-nats
```

## Requirements

- Live stack with default ports.
- NATS server reachable by the triggers service.
