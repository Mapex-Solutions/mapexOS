# Journey: Trigger Teams

## O que esta journey prova

O trigger Teams do serviço triggers dispara pelo pipeline ao vivo e
faz POST do payload MessageCard do webhook para um listener HTTP real
(mesmo formato que um webhook do Teams receberia).

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger faz POST do webhook. |
| `phase2_event_pipeline` *(planejado)* | POST telemetria → gateway → js-executor → router → trigger faz POST do webhook. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-teams
```

## Requisitos

- Stack viva nas portas default.
- Porta `11010` livre no host (sink HTTP compartilhado).
