# Tests

## How to Run
```bash
go test ./... -count=1
```

Run with verbose output:
```bash
go test ./... -count=1 -v
```

Run a specific package:
```bash
go test ./src/modules/events/application/services/... -count=1 -v
go test ./src/modules/events/domain/services/... -count=1 -v
go test ./src/modules/routegroups/application/services/... -count=1 -v
go test ./src/modules/routegroups/interfaces/http/handlers/... -count=1 -v
```

## Prerequisites
- Go 1.25+ installed
- No external services required (tests use mocks for MongoDB, Redis, NATS, and TieredCache)

## Test Coverage
Test files exist for:
- `src/modules/events/application/services/event_service_test.go` -- event processing pipeline
- `src/modules/events/domain/services/match_evaluator_test.go` -- match rule evaluation logic
- `src/modules/routegroups/application/services/routegroup_service_test.go` -- RouteGroup CRUD service
- `src/modules/routegroups/interfaces/http/handlers/routegroup_handler_test.go` -- external HTTP handlers
- `src/modules/routegroups/interfaces/http/handlers/routegroup_internal_handler_test.go` -- internal HTTP handlers
- `src/modules/events/interfaces/message/consumers/asset_invalidate/handler_test.go` -- asset invalidation consumer

## Mocks
Mocks are located alongside each module:
- `src/modules/events/mocks/` -- match evaluator mock
- `src/modules/routegroups/mocks/` -- RouteGroup repository, cache repository, service mocks
- `src/shared/mocks/` -- app cache, NATS bus, TieredCache mocks
