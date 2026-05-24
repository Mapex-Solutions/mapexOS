# Journey: Trigger MQTT

## O que esta journey prova

O trigger MQTT do serviço triggers publica em um broker MQTT via
pipeline ao vivo; validado via oracle events_trigger.

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger publica no broker. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetria → gateway → js-executor → router → trigger publica no broker. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-mqtt
```

## Requisitos

- Stack viva nas portas default.
- Broker MQTT alcançável pelo serviço triggers.
