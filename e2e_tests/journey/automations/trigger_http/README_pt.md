# Journey: Trigger HTTP

## O que esta journey prova

O trigger HTTP do serviço triggers dispara pelo pipeline ao vivo e
faz POST em um listener HTTP real.

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger faz POST no sink. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetria no gateway → js-executor → router → trigger faz POST no sink. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-http
```

## Requisitos

- Stack viva nas portas default.
- Porta `11010` livre no host (sobrescreva `SAGA_TRIGGER_SINK_BIND_ADDR`).
