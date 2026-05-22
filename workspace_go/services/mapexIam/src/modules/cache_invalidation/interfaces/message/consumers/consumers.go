package consumers

import (
	"mapexIam/src/modules/cache_invalidation/application/ports"
	"mapexIam/src/modules/cache_invalidation/interfaces/message/consumers/cache_invalidation"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel file - exports all consumer initialization functions
 *
 * Following Hexagonal Architecture:
 * - Consumers (Interface Layer) only wire transport → service port
 * - Service (Application Layer) handles all business logic
 */

// NewCacheInvalidationConsumer wires the centralized cache invalidation consumer.
// Returns the underlying NATS consumer (or nil if startup failed — error is logged inside).
func NewCacheInvalidationConsumer(bus *natsModel.Bus, service ports.CacheInvalidationServicePort) *natsModel.Consumer {
	return cache_invalidation.NewConsumer(bus, service)
}
