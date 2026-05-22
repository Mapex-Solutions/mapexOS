# Fase 1 — Smoke do trigger MQTT pela conectividade

## O que este teste prova

Trigger MQTT dispara pela cadeia healthmonitor → router → triggers e
publica em um broker MQTT. Verificação via oracle events_trigger: o
serviço triggers marca success=true depois que a publish do broker
retorna, então um count de sucessos em `/api/v1/events/trigger`
equivale a um subscriber real observando a mensagem. Duas transições
reais de saúde → duas execuções de sucesso.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_mqtt/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-mqtt
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- Broker MQTT alcançável pelo serviço triggers em `tcp://localhost:1883` (sobrescreva `broker` em `SagaMqttTrigger` para outros destinos).
- Usuário admin seed provisionado (`admin@mapex.local`).
