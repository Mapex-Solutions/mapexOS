# HTTP Gateway Service

## Summary
High-performance HTTP gateway for MapexOS event ingestion. Receives webhooks from external
systems, authenticates each request according to the Data Source configuration (JWT, OAuth2,
API Key, IP whitelist, or none), and publishes normalized events to NATS JetStream for
downstream processing. Data Source configurations are cached in Redis (cache-aside pattern)
to minimize database lookups on the hot path.

## Responsibilities
- Ingest events via `POST /api/v1/events`
- Manage Data Sources (auth, rate limit, asset binding) via REST CRUD API
- Publish events to NATS (`processor.js.execute`)
- Publish auth failure security events to NATS (`events.raw`) for monitoring
- Cache Data Source lookups in Redis (cache-aside)

## Architecture
Modular Go service with two modules: `datasources` (configuration management + cache) and
`events` (ingestion + authentication).

## Documentation
Deep‑dive documentation (architecture, endpoints, configuration, observability, tests, benchmarks):
- [docs/index.md](docs/index.md)

## How to run
```bash
# build
go build -o bin/http_gateway src/main.go

# run
./bin/http_gateway
```

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
