# Fase 2: Pipeline de eventos (trigger MQTT)

## O que este teste prova

Trigger MQTT dispara pelo pipeline completo de telemetria: um POST
no gateway é transformado pelo js-executor com o script do template
do asset, roteado pelo router de trigger, e o executor do trigger
MQTT publica em um broker mochi-mqtt in-process em porta efêmera.

## Resultado em caso de PASS

- Telemetria postada no gateway vira evento no stream do asset.
- O js-executor executa o script do template e emite um payload
  padronizado.
- O router de trigger do route group seleciona o trigger MQTT.
- `events_trigger` registra ao menos uma execução bem-sucedida.
- O último request data resolvido do trigger contém o prefixo de
  tópico escopado pela saga `mapex-saga/mqtt/`.
- Teardown do asset propaga em cascata pelo Compensate.

## Resultado em caso de FAIL

- Gateway rejeita o POST (API key da data source errada / não
  provisionada).
- Script do template ausente ou o js-executor não consegue rodá-lo
  (sem payload padronizado, o router não vê nada).
- Route group não ligado ao trigger, ou o `RouteGroupIds` do asset
  desviou do RG do trigger configurado.
- Porta do broker MQTT não livre, ou configuração MQTT do trigger
  (host / porta / tópico) desviou do broker gerenciado pela saga.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -v ./journey/automations/trigger_mqtt/phase2_event_pipeline
```

O ponto de entrada é `TestJourney` em `journey_test.go`.

## Composição

`Run` compõe a fase0 (bootstrap de IAM em
`journey/iot/mqtt_broker_auth/phase0_iam_bootstrap`) antes dos
`Items()` desta fase. A fase1 não é incluída — a fase2 sobe seu
próprio broker MQTT, trigger MQTT, route groups, data source,
template e asset de conectividade, e então faz o POST de telemetria
que dispara o trigger.
