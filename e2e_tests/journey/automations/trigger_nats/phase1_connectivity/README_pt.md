# Fase 1 — Smoke do trigger NATS pela conectividade

## O que este teste prova

Trigger NATS dispara pela cadeia healthmonitor → router → triggers e
publica em um subject NATS. Verificação via oracle events_trigger:
success=true é setado depois que a publish do NATS retorna. Duas
transições reais de saúde → duas execuções de sucesso.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_nats/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-nats
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- NATS alcançável pelo serviço triggers em `nats://localhost:4222` (sobrescreva `server` em `SagaNatsTrigger` para outros destinos).
- Usuário admin seed provisionado (`admin@mapex.local`).
