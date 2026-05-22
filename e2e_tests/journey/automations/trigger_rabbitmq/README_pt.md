# Journey: Trigger RabbitMQ

## O que esta journey prova

O trigger RabbitMQ do serviço triggers publica via pipeline ao vivo;
validado via oracle events_trigger.

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger publica na fila. |
| `phase2_event_pipeline` *(planejado)* | POST telemetria → gateway → js-executor → router → trigger publica na fila. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-rabbitmq
```

## Requisitos

- Stack viva nas portas default.
- Broker RabbitMQ alcançável pelo serviço triggers.
