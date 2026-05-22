# Events Service

## Summary
Stores and serves event data for MapexOS. Consumes multiple NATS JetStream streams
(raw, router, trigger, business rule, JS executor, DLQ, processed) in configurable
batches, persists them into ClickHouse via bulk inserts, and exposes query APIs for
analytics and UI. Supports EVA (Entity-Value-Attribute) dynamic field mapping for
processed events, and enforces per-organization retention policies.

## Responsibilities
- Persist event streams from NATS to ClickHouse (batch fetch + bulk insert)
- Expose query APIs for event datasets (raw, router, trigger, businessrule, jsexec, store)
- Manage retention policies (per org, per event type)
- Resolve EVA dynamic fields using tiered template cache (L0 RAM, L1 disk, L2 S3)
- Store failed events via DLQ consumer

## Architecture
Modular Go service with two modules: `events` (persistence + query APIs) and `retention`
(policy management).

## Documentation
Deep‑dive documentation (architecture, endpoints, configuration, observability, tests, benchmarks):
- [docs/index.md](docs/index.md)

## How to run
```bash
# build
go build -o bin/events src/main.go

# run
./bin/events
```

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
