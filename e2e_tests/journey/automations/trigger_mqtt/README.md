# Journey: MQTT trigger

## What this journey proves

The MQTT trigger of the triggers service publishes to an MQTT broker
via the live pipeline; verified via the events_trigger oracle.

| Phase | Firing path |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Healthmonitor force-offline / force-online → trigger publishes to broker. |
| `phase2_event_pipeline` *(planned)* | POST telemetry → gateway → js-executor → router → trigger publishes to broker. |

## How to run

```bash
cd e2e_tests
./run-tests.sh saga trigger-mqtt
```

## Requirements

- Live stack with default ports.
- MQTT broker reachable by the triggers service.
