package di

import (
	"router/src/bootstrap"
	"router/src/modules/events/application/ports"
	routegroupPorts "router/src/modules/routegroups/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// EventServiceDependenciesInjection aggregates all dependencies required
// by the EventService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - AssetCache: TieredCache for asset read models (L0=RAM, L1=Disk, L2=MinIO)
//   - RouteGroupService: For fetching RouteGroup details (already has cache)
//   - NatsBus: Publisher interface for NATS messaging (Hexagonal Architecture)
//   - Metrics: Prometheus metrics for observability
//
// Asset routing info is now fetched via TieredCache instead of HTTP calls.
// Cache invalidation is handled via NATS FANOUT consumer.
//
// The dig.In tag enables automatic dependency injection by the dig container.
type EventServiceDependenciesInjection struct {
	dig.In

	// AssetCache provides TieredCache for asset read models
	// L0 (RAM): Hot cache for frequently accessed assets
	// L1 (Disk): Persistent cache across restarts
	// L2 (MinIO): Source of truth (AssetReadModel JSON)
	AssetCache common.TieredCache `name:"assets"`

	// TemplateCache provides cached access to template metadata for event enrichment
	// Resolves templateName + templateDescription from AssetTemplateID
	TemplateCache ports.TemplateCachePort

	// RouteGroupService provides access to route group details
	RouteGroupService routegroupPorts.RouteGroupServicePort

	// NatsBus provides NATS core publishing (fire-and-forget + batch flush)
	NatsBus natsModel.CorePublisher

	// Metrics provides Prometheus metrics for service observability
	Metrics *bootstrap.RouterMetrics
}
