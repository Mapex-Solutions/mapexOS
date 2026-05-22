# Operations

## Start
```bash
./bin/events
```

## Smoke checks
- `GET /metrics`
- `GET /health`
- `GET /api/v1/events/raw` with valid auth

## Dependencies
| Dependency | Required | Purpose |
|---|---|---|
| ClickHouse | yes | Event persistence (bulk insert + queries) |
| NATS | yes | Event stream consumption (7 consumers) |
| Redis (app) | yes | Application cache |
| Redis (shared) | yes | Shared cache (permissions, coverage) |
| MongoDB | yes | Retention policy persistence |
| MinIO | yes | Templates cache L2 (EVA field resolution) |

## Module Initialization Order
`retention` must initialize before `events` — the Events module requires retention policies to resolve TTL.

## Benchmarks
See [Benchmarks](../benchmarks/index.md) for load-testing methodology, scripts, and results.
