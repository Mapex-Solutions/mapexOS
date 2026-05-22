# Fase 1 — Smoke do trigger Slack pela conectividade

## O que este teste prova

Trigger Slack dispara pela cadeia healthmonitor → router → triggers e
faz POST do webhook para um listener HTTP.

Webhooks Slack são POST HTTP puros, então esta fase reusa o mesmo
sink HTTP in-process que a fase do trigger HTTP usa. Duas transições
reais de saúde → dois POSTs de webhook recebidos.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_slack/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-slack
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Verifique com `./run-tests.sh check`.
- Porta `11010` livre no host (sink HTTP compartilhado — sobrescreva `SAGA_TRIGGER_SINK_BIND_ADDR`).
- Quando o serviço triggers roda em Docker, defina `SAGA_TRIGGER_SINK_URL=http://host.docker.internal:11010`.
- Usuário admin seed provisionado (`admin@mapex.local`).
