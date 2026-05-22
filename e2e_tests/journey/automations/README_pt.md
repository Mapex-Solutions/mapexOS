# Contexto: Automations

## O que este contexto cobre

Saga journeys do domínio **automations** — CRUD de triggers + execução
ponta a ponta contra receivers específicos do protocolo.

| Journey | Status | Receiver / Oracle |
|---|---|---|
| [`trigger_http`](./trigger_http/)         | ✅ phase1 passa | sink HTTP in-process |
| [`trigger_email`](./trigger_email/)       | ✅ phase1 passa | sink SMTP in-process (`emersion/go-smtp`) |
| [`trigger_websocket`](./trigger_websocket/) | ✅ phase1 passa | sink WS in-process (`gorilla/websocket`) + events_trigger |
| [`trigger_slack`](./trigger_slack/)       | ⏳ skip — bug backend | (reusaria sink HTTP) |
| [`trigger_teams`](./trigger_teams/)       | ⏳ skip — bug backend | (reusaria sink HTTP) |
| [`trigger_mqtt`](./trigger_mqtt/)         | ⏳ skip — auth do broker | oracle events_trigger |
| [`trigger_nats`](./trigger_nats/)         | ⏳ skip — auth do broker | oracle events_trigger |
| [`trigger_rabbitmq`](./trigger_rabbitmq/) | ⏳ skip — precisa broker de teste | oracle events_trigger |

Cada `trigger_<type>` segue o mesmo layout:

```
trigger_<type>/
├── README.md / README_pt.md
├── phase1_connectivity/  ─ force online/offline via healthmonitor → trigger
│   ├── journey.go        (um // comentário por item da saga)
│   ├── journey_test.go
│   ├── README.md
│   └── README_pt.md
└── phase2_event_pipeline/  (planejado — POST telemetria → gateway → js-executor → router → trigger)
```

## Bugs descobertos por estas journeys

| Descoberta | Onde | Status |
|---|---|---|
| **Slack executor espera config flat** (`config["webhookUrl"]`) mas recebe a forma union `config["slack"]`; também lê `config["text"]` em vez de `config["slack"]["message"]`. | `workspace_go/services/triggers/src/modules/events/infrastructure/communications/slack/slack_executor.go` | journey trigger_slack bloqueada |
| **Teams executor espera config flat** (mesma divergência). | `workspace_go/.../communications/teams/teams_executor.go` | journey trigger_teams bloqueada |
| **Email executor sempre usa `smtp.PlainAuth`** mesmo com credenciais vazias, o que falha contra servidores SMTP que não anunciam AUTH. A journey trigger_email passa porque o sink `go-smtp` da saga implementa explicitamente `AuthSession`. | `workspace_go/.../communications/email/email_executor.go:112` | journey passa; bug latente em prod contra SMTPs sem AUTH |
| **MQTT executor constrói `protocol://broker:port`** assumindo `broker` é hostname puro. UI / docs do contrato deveriam deixar isso explícito. | `workspace_go/.../technical/mqtt/mqtt_executor.go` | confirmado via payload da saga |
| **Broker MQTT da plataforma exige auth por asset** sem caminho para creds de trigger. | (arquitetura) | bloqueia smoke trigger_mqtt |
| **NATS da plataforma exige auth** sem creds expostas no config do trigger. | (arquitetura) | bloqueia smoke trigger_nats |

## Como rodar tudo

```bash
cd e2e_tests
./run-tests.sh saga triggers
```

## Requisitos

- Portas default da stack (`./run-tests.sh check` confirma).
- Portas livres no host: `11010` (sink HTTP), `11025` (sink SMTP), `11026` (sink WS).
- Para as journeys atualmente skipped: veja cada README de fase para a peça faltante.
