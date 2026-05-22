# MapexOS API Service

## Summary
Core IAM and tenant management API for MapexOS. Manages authentication (JWT/OAuth2),
multi-level organization hierarchy (Vendor → Customer → …), flexible RBAC (custom roles,
groups, memberships), lookup lists, and onboarding workflows. Publishes cache invalidation
events via NATS JetStream so downstream services refresh authorization and coverage caches.
10 modules initialized in dependency order.

## Responsibilities
- Authenticate users (JWT HS256 / OAuth2 JWKS RS256) and manage sessions
- CRUD for users, organizations, roles, groups, memberships, and lists
- Build and serve authorization and coverage caches for internal services
- Publish cache invalidation events via NATS JetStream (`MAPEXOS_CACHE_INVALIDATION`)
- Orchestrate onboarding workflows (user + membership creation)

## Architecture
Modular Go service with 10 coordinated modules managing IAM, organizations, roles, groups,
memberships, and onboarding workflows.

## Documentation
Deep‑dive documentation (architecture, endpoints, configuration, observability, tests, benchmarks):
- [docs/index.md](docs/index.md)

## How to run
```bash
# build
go build -o bin/mapexos src/main.go

# run
./bin/mapexos
```

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
