# Fase 1 — Smoke do trigger RabbitMQ pela conectividade

## O que este teste prova

Trigger RabbitMQ dispara pela cadeia healthmonitor → router →
triggers e publica em uma fila. Verificação via oracle
events_trigger: success=true é setado depois que a publish do AMQP
retorna. Duas transições reais de saúde → duas execuções de sucesso.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_rabbitmq/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-rabbitmq
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- Broker RabbitMQ alcançável pelo serviço triggers em `amqp://guest:guest@localhost:5672/` (sobrescreva `host` / `port` / `username` / `password` em `SagaRabbitmqTrigger` para outros destinos).
- Usuário admin seed provisionado (`admin@mapex.local`).
