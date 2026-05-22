# Operations

## Start
```bash
./bin/http_gateway
```

## Smoke checks
- `GET /health` — should return `{"status": "healthy", ...}` with HTTP 200
- `GET /metrics` — should return Prometheus text format
- `POST /api/v1/events?ds=<validId>` with a valid Data Source and appropriate auth

## Dependencies
| Dependency | Required | Purpose |
|---|---|---|
| MongoDB | yes | Data Source persistence |
| Redis (app) | yes | Data Source cache (cache-aside) |
| Redis (shared) | yes | Permission and coverage cache |
| NATS | yes | Event publishing (`processor.js.execute`, `events.raw`) |
| MapexOS API | yes | Permission and coverage resolution (`MAPEXOS_URL`) |

## Benchmarks
See [Benchmarks](../benchmarks/index.md) for methodology and results.
