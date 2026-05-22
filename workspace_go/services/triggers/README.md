# Triggers Service

## Summary
Outbound execution layer of MapexOS. Consumes trigger execution events from NATS JetStream
(produced by RuleEngine and Router), resolves dynamic placeholders from the event payload,
and dispatches actions through a registry of executors (HTTP, MQTT, RabbitMQ, NATS, WebSocket,
Email, Teams, Slack). Publishes execution results to `events.trigger` for auditing and analytics.

## Responsibilities
- Execute trigger actions consumed from NATS (`trigger.*.execute`)
- Resolve placeholders in trigger configs using event payload data
- Manage trigger configuration and templates via REST CRUD API
- Cache trigger definitions in Redis for low-latency execution
- Publish execution results to `events.trigger`

## Architecture
Go service organized into two modules: `triggers` (configuration management) and `events`
(execution pipeline). Executors are selected by `triggerType` via a factory registry.

## Documentation
Deep‑dive documentation (architecture, endpoints, configuration, observability, tests, benchmarks):
- [docs/index.md](docs/index.md)

## How to run
```bash
# build
go build -o bin/triggers src/main.go

# run
./bin/triggers
```

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
