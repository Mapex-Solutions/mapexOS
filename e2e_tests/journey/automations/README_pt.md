# Contexto: Automations

## O que este contexto cobre

Saga journeys do domínio **automations** — CRUD de triggers + execução
ponta a ponta contra receivers específicos do protocolo.

| Journey | Receiver / Oracle |
|---|---|
| [`trigger_http`](./trigger_http/)         | sink HTTP in-process |
| [`trigger_email`](./trigger_email/)       | sink SMTP in-process (`emersion/go-smtp`) |
| [`trigger_websocket`](./trigger_websocket/) | sink WS in-process (`gorilla/websocket`) + events_trigger |
| [`trigger_slack`](./trigger_slack/)       | sink HTTP in-process (Slack webhook = POST HTTP) |
| [`trigger_teams`](./trigger_teams/)       | sink HTTP in-process (Teams webhook = POST HTTP) |
| [`trigger_mqtt`](./trigger_mqtt/)         | oracle events_trigger (caminho de publish no broker) |
| [`trigger_nats`](./trigger_nats/)         | oracle events_trigger (caminho de publish no NATS) |
| [`trigger_rabbitmq`](./trigger_rabbitmq/) | oracle events_trigger (caminho de publish AMQP) |

Cada `trigger_<type>` segue o mesmo layout:

```
trigger_<type>/
├── README.md / README_pt.md
├── phase1_connectivity/  ─ healthmonitor force online/offline → trigger
│   ├── journey.go        (um // comentário por item da saga)
│   ├── journey_test.go
│   ├── README.md
│   └── README_pt.md
└── phase2_event_pipeline/  POST telemetria → gateway → js-executor → router → trigger
    ├── journey.go
    ├── journey_test.go
    ├── README.md
    └── README_pt.md
```

## Como rodar tudo

```bash
cd e2e_tests
./run-tests.sh saga triggers
```

## Requisitos

- Portas default da stack (`./run-tests.sh check` confirma).
- Portas livres no host: `11010` (sink HTTP), `11025` (sink SMTP), `11026` (sink WS).
- Brokers acessíveis pelo serviço triggers: MQTT (`:1883`), NATS (`:4222`), RabbitMQ (`:5672`).
