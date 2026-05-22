package di

import (
	"triggers/src/bootstrap"
	"triggers/src/modules/events/application/ports"
	triggerPorts "triggers/src/modules/triggers/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// EventServiceDependenciesInjection aggregates all dependencies required
// by the EventService.
//
// Following Hexagonal Architecture and Dependency Inversion Principle:
// - Application depends on PORTS (interfaces), never on concrete implementations
// - Infrastructure provides the concrete adapters via DI container
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - TriggerService: For fetching trigger configuration by ID (uses cache)
//   - ExecutorRegistry: For getting the appropriate executor based on trigger type
//   - NatsBus: CorePublisher for fire-and-forget NATS publishing with batch flush
//   - Metrics: Prometheus metrics for observability
//
// The dig.In tag enables automatic dependency injection by the dig container.
type EventServiceDependenciesInjection struct {
	dig.In

	// TriggerService provides access to trigger details (with cache-aside pattern)
	// This is a cross-domain service dependency (application port)
	TriggerService triggerPorts.TriggerServicePort

	// ExecutorRegistry provides access to trigger execution adapters
	// This is an infrastructure dependency injected via the application port interface
	ExecutorRegistry ports.ExecutorRegistry

	// NatsBus for publishing trigger execution events to events service (fire-and-forget)
	NatsBus natsModel.CorePublisher

	// Metrics provides Prometheus observability for the triggers service
	Metrics *bootstrap.TriggerMetrics
}
