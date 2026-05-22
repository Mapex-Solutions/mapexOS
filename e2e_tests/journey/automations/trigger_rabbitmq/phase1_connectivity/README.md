# Phase 1 — Connectivity-driven RabbitMQ trigger smoke

## What this test proves

RabbitMQ trigger fires through the healthmonitor → router → triggers
chain and publishes to a queue. Validation uses the events_trigger
oracle: success=true is set after the AMQP publish returns. Two real
health transitions → two successful executions.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_rabbitmq/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-rabbitmq
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- RabbitMQ broker reachable from the triggers service at `amqp://guest:guest@localhost:5672/` (override `host` / `port` / `username` / `password` in `SagaRabbitmqTrigger` for other targets).
- Seed admin user provisioned (`admin@mapex.local`).
