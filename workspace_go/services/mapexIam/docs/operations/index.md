# Operations

## Start
```bash
./bin/mapexos
```

## Smoke checks
- `GET /health`
- `POST /api/v1/auth/login`

## Dependencies
| Dependency | Required | Purpose |
|---|---|---|
| MongoDB | yes | IAM data persistence (users, orgs, roles, groups, memberships, lists) |
| Redis (app) | yes | Application cache |
| Redis (shared) | yes | Authorization and coverage caches |
| NATS | yes | Cache invalidation events + list name sync |

> **Note**: The service automatically initializes its modules in the correct dependency order. See [Architecture](../architecture/index.md) for details.

## Benchmarks
See [Benchmarks](../benchmarks/index.md) for load-testing methodology, scripts, and results.
