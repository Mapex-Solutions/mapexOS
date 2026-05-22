# Tests

## Run
```bash
go test ./... -count=1
```

## Prerequisites
- MongoDB running locally (for integration tests)
- Redis running locally (for cache handler tests)
- NATS JetStream running locally (for integration tests)

Unit tests (no external dependencies) can be run standalone.

## Test File Locations
| File | Module | Description |
|---|---|---|
| `src/modules/events/application/services/placeholder_resolver_test.go` | events | Placeholder resolution logic |
| `src/modules/events/application/services/event_service_test.go` | events | Event service unit tests |
| `src/modules/events/application/services/event_service_integration_test.go` | events | Integration tests (requires infra) |
| `src/modules/events/infrastructure/communications/email/email_executor_test.go` | events | Email executor |
| `src/modules/events/infrastructure/communications/slack/slack_executor_test.go` | events | Slack executor |
| `src/modules/events/infrastructure/communications/teams/teams_executor_test.go` | events | Teams executor |
| `src/modules/events/infrastructure/registry/executor_registry_test.go` | events | Executor factory registry |
| `src/modules/events/infrastructure/technical/http/http_executor_test.go` | events | HTTP executor |
| `src/modules/events/infrastructure/technical/mqtt/mqtt_executor_test.go` | events | MQTT executor |
| `src/modules/events/infrastructure/technical/nats/nats_executor_test.go` | events | NATS executor |
| `src/modules/events/infrastructure/technical/rabbitmq/rabbitmq_executor_test.go` | events | RabbitMQ executor |
| `src/modules/events/infrastructure/technical/websocket/websocket_executor_test.go` | events | WebSocket executor |
| `src/modules/triggers/application/handlers/cache.handler_test.go` | triggers | Cache handler |
| `src/modules/triggers/application/services/trigger_service_test.go` | triggers | Trigger service unit tests |
