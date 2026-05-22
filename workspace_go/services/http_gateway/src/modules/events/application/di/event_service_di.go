package di

import (
	"http_gateway/src/bootstrap"
	"http_gateway/src/modules/events/application/ports"

	"go.uber.org/dig"
)

// EventServiceDependenciesInjection aggregates all dependencies required
// by the EventService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - NatsBus: NATS publish surface (Publish + PublishCore) typed as the
//     module-local port so unit tests can substitute a mock without
//     standing up a real broker. The concrete *natsModel.Bus satisfies it.
//   - Metrics: Service-specific Prometheus metrics for instrumentation.
//
// The dig.In tag enables automatic dependency injection by the dig container.
type EventServiceDependenciesInjection struct {
	dig.In
	NatsBus ports.EventBusPort
	Metrics *bootstrap.HttpGatewayMetrics
}
