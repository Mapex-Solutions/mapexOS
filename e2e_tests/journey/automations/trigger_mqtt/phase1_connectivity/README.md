# Phase 1 — Connectivity-driven MQTT trigger smoke

## What this test proves

MQTT trigger fires through the healthmonitor → router → triggers
chain and publishes to an MQTT broker. Validation uses the
events_trigger oracle: the triggers service marks success=true after
the broker publish call returns, so a non-empty success count in
`/api/v1/events/trigger` is equivalent to a real subscriber observing
the message. Two real health transitions → two successful executions.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_mqtt/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-mqtt
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- MQTT broker reachable from the triggers service at `tcp://localhost:1883` (override `broker` in `SagaMqttTrigger` for other targets).
- Seed admin user provisioned (`admin@mapex.local`).
