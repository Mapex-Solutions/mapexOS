# HTTP Gateway Documentation

## Overview
The HTTP Gateway is the HTTP entrypoint of MapexOS. It receives external webhook events, authenticates them according to each Data Source configuration, and publishes normalized events to NATS JetStream for downstream processing (JS‑Executor, Events, RuleEngine, etc.). It also exposes a CRUD API to manage Data Sources, including authentication strategy, asset binding, and rate limit settings.

## Responsibilities
- Ingest events via `POST /api/v1/events`
- Authenticate each event using the Data Source configuration (JWT, OAuth2, API Key, IP whitelist, or none)
- Publish events to NATS JetStream (`processor.js.execute`)
- Publish auth failure security events to NATS (`events.raw`) for monitoring
- Manage Data Sources via REST API (create/read/update/delete)
- Cache Data Source lookups in Redis (cache-aside pattern)

## Non‑Responsibilities
- Business rule evaluation (RuleEngine)
- Script execution (JS‑Executor)
- Event persistence (Events service)

## Primary Data Flow
1. External system sends webhook to `POST /api/v1/events?ds=<dataSourceId>`
2. Gateway resolves the Data Source (cache-aside in Redis)
3. Gateway validates auth based on the Data Source auth type
   - On auth failure: publishes a security event to `events.raw` (fire-and-forget) and returns 401
4. Gateway publishes event to NATS JetStream (`processor.js.execute`)

## Docs Map
- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)
