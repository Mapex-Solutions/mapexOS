# Context: Automations

## What this context covers

Saga journeys for the **automations** domain — trigger CRUD + end-to-end
execution against protocol-specific receivers.

| Journey | Receiver / Oracle |
|---|---|
| [`trigger_http`](./trigger_http/)         | in-process HTTP sink |
| [`trigger_email`](./trigger_email/)       | in-process SMTP sink (`emersion/go-smtp`) |
| [`trigger_websocket`](./trigger_websocket/) | in-process WS sink (`gorilla/websocket`) + events_trigger |
| [`trigger_slack`](./trigger_slack/)       | in-process HTTP sink (Slack webhook = POST HTTP) |
| [`trigger_teams`](./trigger_teams/)       | in-process HTTP sink (Teams webhook = POST HTTP) |
| [`trigger_mqtt`](./trigger_mqtt/)         | events_trigger oracle (broker publish path) |
| [`trigger_nats`](./trigger_nats/)         | events_trigger oracle (NATS publish path) |
| [`trigger_rabbitmq`](./trigger_rabbitmq/) | events_trigger oracle (AMQP publish path) |

Each `trigger_<type>` follows the same layout:

```
trigger_<type>/
├── README.md / README_pt.md
├── phase1_connectivity/  ─ healthmonitor force online/offline → trigger
│   ├── journey.go        (one-line // comment per saga item)
│   ├── journey_test.go
│   ├── README.md
│   └── README_pt.md
└── phase2_event_pipeline/  (planned — POST telemetry → gateway → js-executor → router → trigger)
```

## How to run all

```bash
cd e2e_tests
./run-tests.sh saga triggers
```

## Requirements

- Default stack ports (`./run-tests.sh check` confirms).
- Free ports on the host: `11010` (HTTP sink), `11025` (SMTP sink), `11026` (WS sink).
- Brokers reachable from the triggers service: MQTT (`:1883`), NATS (`:4222`), RabbitMQ (`:5672`).
