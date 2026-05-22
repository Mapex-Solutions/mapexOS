# Router Service

## Summary
Routes incoming asset events to downstream processors based on RouteGroup matching rules.
Consumes events from NATS JetStream, resolves asset context via TieredCache, evaluates
conditional match rules, and publishes routed events to multiple downstream subjects
(RuleEngine, Triggers, Event Store, Lakehouse, Notifications).

## Responsibilities
- Evaluate RouteGroup match rules for incoming events (policy `all`/`any`, operators `eq`/`neq`/`gt`/`gte`/`lt`/`lte`/`in`/`nin`)
- Publish routing results to downstream NATS subjects (fan-out to multiple destinations)
- Manage RouteGroup configuration via REST API (CRUD + internal MS-to-MS endpoint)
- Resolve asset context through TieredCache (L0 RAM, L1 Disk, L2 MinIO, Fallback HTTP)
- Handle cache invalidation via NATS FANOUT consumer
- Emit routing history for audit and UI visualization

## Architecture
Modular Go service with two modules: `routegroups` (configuration management + cache) and
`events` (routing pipeline + NATS consumers).

## Documentation
Deep-dive documentation (architecture, endpoints, configuration, observability, tests, benchmarks):
- [docs/index.md](docs/index.md)

## How to run
```bash
# build
go build -o bin/router src/main.go

# run
./bin/router
```

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
