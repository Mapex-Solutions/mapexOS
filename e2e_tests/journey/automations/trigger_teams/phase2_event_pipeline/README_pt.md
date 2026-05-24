# Fase 2: Pipeline de eventos (trigger Teams)

## O que este teste prova

Trigger Teams dispara pelo pipeline completo de telemetria: um POST
no gateway é transformado pelo js-executor com o script do template
do asset, roteado pelo router de trigger, e o executor do trigger
Teams entrega o webhook (webhooks Teams são POST HTTP) para um sink
HTTP in-process que substitui o endpoint de webhook do Teams.

## Resultado em caso de PASS

- Telemetria postada no gateway vira evento no stream do asset.
- O js-executor executa o script do template e emite um payload
  padronizado.
- O router de trigger do route group seleciona o trigger Teams.
- `events_trigger` registra ao menos uma execução bem-sucedida.
- O sink HTTP in-process recebe o POST do webhook Teams.
- Teardown do asset propaga em cascata pelo Compensate.

## Resultado em caso de FAIL

- Gateway rejeita o POST (API key da data source errada / não
  provisionada).
- Script do template ausente ou o js-executor não consegue rodá-lo
  (sem payload padronizado, o router não vê nada).
- Route group não ligado ao trigger, ou o `RouteGroupIds` do asset
  desviou do RG do trigger configurado.
- `webhookUrl` do trigger Teams desviou do sink HTTP gerenciado pela
  saga, ou o sink não é alcançável a partir do serviço triggers.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -v ./journey/automations/trigger_teams/phase2_event_pipeline
```

O ponto de entrada é `TestJourney` em `journey_test.go`.

## Composição

`Run` compõe a fase0 (bootstrap de IAM em
`journey/iot/mqtt_broker_auth/phase0_iam_bootstrap`) antes dos
`Items()` desta fase. A fase1 não é incluída — a fase2 sobe seu
próprio sink HTTP, trigger Teams, route groups, data source,
template e asset de conectividade, e então faz o POST de telemetria
que dispara o trigger.
