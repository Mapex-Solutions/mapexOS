# Module e2e: triggers / triggers

## Escopo

Blocos de construção de saga para o módulo `triggers` do serviço
triggers (`workspace_go/services/triggers/src/modules/triggers/`).
Este pacote exercita o contrato CRUD público de triggers nos oito
tipos suportados pela plataforma — HTTP, Email (SMTP), WebSocket,
Slack, Teams, MQTT, NATS e RabbitMQ — compondo um step de create por
tipo junto com um sink, broker ou container in-process que as
journeys de saga sobem para receber a saída do trigger. Não há
funções `Test*` aqui: a suite vive em `journey/automations/` e
carrega estes steps, builders de payload e asserts para dirigir o
serviço triggers ponta-a-ponta.

## Endpoints exercitados

- `POST /api/v1/triggers` — registra um trigger em qualquer um dos
  oito tipos suportados. Cada step `Create*Trigger` envia uma captura
  literal da UI e parametriza apenas o nome e o destino (URL de
  endpoint, host/porta do broker, URL do server, etc.) para que o
  trigger dispare contra o sink gerenciado pela saga.
- `DELETE /api/v1/triggers/{id}` — chamado pelo Compensate de todo
  step de create; idempotente contra `404`.

## Fixtures

O pacote não tem fixtures JSON; os corpos de payload são literais Go
capturados da UI e parametrizados por execução pelos builders em
`payloads/`.

| Arquivo                              | Propósito                                                                                            |
|--------------------------------------|------------------------------------------------------------------------------------------------------|
| `payloads/saga_simple_trigger.go`    | Trigger HTTP genérico; endpoint reescrito para `constants.TriggerSinkURL`.                            |
| `payloads/saga_email_trigger.go`     | Trigger de email apontando para o SMTP sink in-process.                                               |
| `payloads/saga_websocket_trigger.go` | Trigger WebSocket apontando para o server de upgrade WS in-process.                                   |
| `payloads/saga_slack_trigger.go`     | Trigger Slack; webhook URL reescrita para o HTTP sink de modo que o post seja capturado localmente.   |
| `payloads/saga_teams_trigger.go`     | Trigger Microsoft Teams; webhook URL reescrita para o HTTP sink.                                      |
| `payloads/saga_mqtt_trigger.go`      | Trigger MQTT apontando para o broker mochi-mqtt embarcado que a saga inicia.                          |
| `payloads/saga_nats_trigger.go`      | Trigger NATS apontando para o servidor NATS embarcado que a saga inicia.                              |
| `payloads/saga_rabbitmq_trigger.go`  | Trigger RabbitMQ apontando para o container RabbitMQ efêmero que a saga sobe (testcontainers).        |

A pasta `steps/` espelha um `Create*Trigger` por payload mais os
steps de ciclo de vida do sink/broker/container
(`start_test_sink.go`, `start_smtp_sink.go`, `start_websocket_sink.go`,
`start_mqtt_broker.go`, `start_nats_server.go`,
`start_rabbitmq_container.go`) e `compensate_helpers.go` para a
lógica de teardown compartilhada.

A pasta `asserts/` traz `assert_sink_hit.go` (HTTP, WebSocket,
contadores genéricos) e `assert_smtp_received.go` (validação de
subject/to/from/body sem depender de um parser MX externo).

## Como rodar

Este pacote não tem arquivos `*_test.go`; rodar `go test
./services/triggers/triggers` vai imprimir `[no test files]`. Os
steps são consumidos pelas journeys de automation em
`journey/automations/trigger_*/`:

```bash
cd e2e_tests

# Todas as oito journeys de trigger
go test -tags=saga ./journey/automations/...

# Um tipo de trigger ponta-a-ponta
go test -tags=saga ./journey/automations/trigger_http/...
go test -tags=saga ./journey/automations/trigger_email/...
```

## Resultado em caso de PASS

Quando as journeys consumidoras passam, este pacote provou que o
módulo triggers respeita `POST /api/v1/triggers` (e o `DELETE` do
Compensate) para todos os tipos suportados, e que a configuração de
trigger resolvida realmente dispara contra um destino vivo —
observado por um sink, broker ou container in-process em vez do
round-trip pelo serviço de events.

## Requisitos

- `triggers` disponível na porta `5006` (override via
  `TRIGGERS_URL`).
- `mapexos` disponível na porta `5000` para gerar o token de admin.
- Portas `11010` (HTTP sink), `11025` (SMTP sink), `11026` (WebSocket
  sink) livres no host; os steps de broker/server/container ocupam
  portas atribuídas pelo SO.
- Docker disponível para o step de testcontainers do RabbitMQ.

## Notas

- Cada tipo de trigger tem seu próprio namespace de bag keys (ver
  `steps/keys.go`), então uma journey pode compor vários triggers na
  mesma execução sem colisão.
- Os triggers HTTP, Slack e Teams compartilham o mesmo HTTP sink
  genérico — o assert apenas compara o corpo capturado contra o
  envelope de webhook esperado para o tipo.
- Os servidores NATS / mochi-mqtt embarcados e o container RabbitMQ
  são derrubados pelo Compensate da saga; execuções com falha podem
  deixar containers para trás que aparecem em `docker ps`.
- Oito tipos de trigger são exercitados inline pelos steps de create;
  novos tipos adicionados ao módulo exigem um novo par
  `create_<tipo>_trigger.go` + `saga_<tipo>_trigger.go`, não uma
  alteração nos arquivos existentes.
