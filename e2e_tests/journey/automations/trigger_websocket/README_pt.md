# Journey: Trigger WebSocket

## O que esta journey prova

Trigger WebSocket conecta + escreve em um servidor WS in-process via
pipeline ao vivo; validado via oracle events_trigger.

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger disca para o sink WS. |
| `phase2_event_pipeline` *(planejado)* | POST telemetria → gateway → js-executor → router → trigger disca para o sink WS. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-websocket
```

## Requisitos

- Stack viva nas portas default.
- Porta `11026` livre no host (sobrescreva `SAGA_TRIGGER_WS_BIND_ADDR`).
