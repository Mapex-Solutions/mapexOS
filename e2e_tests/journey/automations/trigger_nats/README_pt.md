# Journey: Trigger NATS

## O que esta journey prova

O trigger NATS do serviço triggers publica em um subject NATS via
pipeline ao vivo; validado via oracle events_trigger.

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger publica no subject. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetria → gateway → js-executor → router → trigger publica no subject. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-nats
```

## Requisitos

- Stack viva nas portas default.
- NATS alcançável pelo serviço triggers.
