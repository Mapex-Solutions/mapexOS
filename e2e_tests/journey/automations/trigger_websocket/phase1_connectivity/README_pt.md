# Fase 1 — Smoke do trigger WebSocket pela conectividade

## O que este teste prova

Trigger WebSocket dispara pela cadeia healthmonitor → router →
triggers e conecta em um servidor WS in-process, escrevendo um
frame.

A fase sobe um sink WS local em `WsSinkBindAddr` (`/ws`), cria um
trigger WebSocket apontando para ele, liga dois route groups
(`online` / `offline`) ao trigger, e então guia um asset por duas
transições reais de saúde. O oracle events_trigger confirma que o
trigger publicou com sucesso em cada transição.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_websocket/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-websocket
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- Porta `11026` livre no host (sobrescreva `SAGA_TRIGGER_WS_BIND_ADDR`).
- Quando o serviço triggers roda em Docker, defina `SAGA_TRIGGER_WS_URL=ws://host.docker.internal:11026/ws`.
- Usuário admin seed provisionado (`admin@mapex.local`).
