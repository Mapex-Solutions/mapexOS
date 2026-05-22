# Journey: RabbitMQ trigger

## What this journey proves

The RabbitMQ trigger of the triggers service publishes via the live
pipeline; verified via the events_trigger oracle.

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger publishes to queue. |
| `phase2_event_pipeline` *(planned)* | POST telemetry → gateway → js-executor → router → trigger publishes to queue. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-rabbitmq
```

## Requirements

- Live stack with default ports.
- RabbitMQ broker reachable by the triggers service.
