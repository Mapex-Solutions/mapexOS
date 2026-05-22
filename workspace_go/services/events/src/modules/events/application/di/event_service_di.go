package di

import (
	"events/src/bootstrap"
	"events/src/modules/events/application/ports"
	"events/src/modules/events/domain/repositories"
	retentionPorts "events/src/modules/retention/application/ports"

	"go.uber.org/dig"
)

// EventServiceDependenciesInjection aggregates all dependencies required
// by the EventService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - EventRepo: For storing events in ClickHouse
//   - RetentionService: For getting retention days per organization
//   - TemplateCache: For EVA field resolution from AssetTemplate.DynamicFields
//
// The dig.In tag enables automatic dependency injection by the dig container.
type EventServiceDependenciesInjection struct {
	dig.In

	// EventRepo provides event storage operations in ClickHouse
	EventRepo repositories.EventRepository

	// RetentionService provides retention days for each organization
	RetentionService retentionPorts.RetentionServicePort

	// TemplateCache provides cached access to AssetTemplate for EVA field resolution
	TemplateCache ports.TemplateCachePort

	// Metrics provides Prometheus observability for the events service
	Metrics *bootstrap.EventsMetrics
}
