package consumers

import (
	"router/src/modules/events/application/ports"
	"router/src/modules/events/interfaces/message/consumers/asset_invalidate"
	"router/src/modules/events/interfaces/message/consumers/route_execute"
	"router/src/modules/events/interfaces/message/consumers/template_invalidate"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel file - exports all consumer initialization functions
 *
 * Following Hexagonal Architecture:
 * - Consumers (Interface Layer) only receive messages and call service
 * - Service (Application Layer) handles all business logic and message lifecycle
 *
 * Consumer Patterns:
 * - WorkQueue: Durable, load-balanced (route_execute) - uses StartConsumer
 * - FANOUT: Ephemeral, broadcast (asset_invalidate) - uses SubscribeFanout
 */

// NewRouteExecuteConsumer creates the route execution consumer (WorkQueue pattern)
func NewRouteExecuteConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	route_execute.NewConsumer(bus, eventService)
}

// NewAssetInvalidateConsumer creates the asset cache invalidation consumer (FANOUT pattern)
func NewAssetInvalidateConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	asset_invalidate.NewConsumer(bus, eventService)
}

// NewTemplateInvalidateConsumer creates the template cache invalidation consumer (FANOUT pattern)
func NewTemplateInvalidateConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	template_invalidate.NewConsumer(bus, eventService)
}
