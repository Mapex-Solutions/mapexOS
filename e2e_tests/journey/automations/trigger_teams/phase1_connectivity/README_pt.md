# Fase 1 — Smoke do trigger Teams pela conectividade

## O que este teste prova

Trigger Teams dispara pela cadeia healthmonitor → router → triggers e
faz POST do MessageCard para um listener HTTP.

Webhooks Teams são POST HTTP, então esta fase reusa o mesmo sink HTTP
compartilhado. Duas transições reais de saúde → dois POSTs recebidos.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_teams/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-teams
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`.
- Porta `11010` livre no host (sink HTTP compartilhado — sobrescreva `SAGA_TRIGGER_SINK_BIND_ADDR`).
- Quando o serviço triggers roda em Docker, defina `SAGA_TRIGGER_SINK_URL=http://host.docker.internal:11010`.
- Usuário admin seed provisionado (`admin@mapex.local`).
