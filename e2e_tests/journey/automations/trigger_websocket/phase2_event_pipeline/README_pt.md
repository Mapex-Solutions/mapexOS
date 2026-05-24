# Fase 2: Pipeline de eventos (trigger WebSocket)

## O que este teste prova

Trigger WebSocket dispara pelo pipeline completo de telemetria: um
POST no gateway é transformado pelo js-executor com o script do
template do asset, roteado pelo router de trigger, e o executor do
trigger WebSocket abre conexão com um sink WS in-process e envia a
mensagem.

## Resultado em caso de PASS

- Telemetria postada no gateway vira evento no stream do asset.
- O js-executor executa o script do template e emite um payload
  padronizado.
- O router de trigger do route group seleciona o trigger WebSocket.
- `events_trigger` registra ao menos uma execução bem-sucedida.
- O último request data resolvido do trigger contém o marcador da
  saga `"saga":"websocket"`.
- Teardown do asset propaga em cascata pelo Compensate.

## Resultado em caso de FAIL

- Gateway rejeita o POST (API key da data source errada / não
  provisionada).
- Script do template ausente ou o js-executor não consegue rodá-lo
  (sem payload padronizado, o router não vê nada).
- Route group não ligado ao trigger, ou o `RouteGroupIds` do asset
  desviou do RG do trigger configurado.
- Sink WS não associado, ou configuração WebSocket do trigger (URL /
  template de mensagem) desviou do sink gerenciado pela saga.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -v ./journey/automations/trigger_websocket/phase2_event_pipeline
```

O ponto de entrada é `TestJourney` em `journey_test.go`.

## Composição

`Run` compõe a fase0 (bootstrap de IAM em
`journey/iot/mqtt_broker_auth/phase0_iam_bootstrap`) antes dos
`Items()` desta fase. A fase1 não é incluída — a fase2 sobe seu
próprio sink WS, trigger WebSocket, route groups, data source,
template e asset de conectividade, e então faz o POST de telemetria
que dispara o trigger.
