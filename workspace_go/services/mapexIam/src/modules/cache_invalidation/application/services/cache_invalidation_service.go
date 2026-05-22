package services

import (
	"mapexIam/src/modules/cache_invalidation/application/di"
	"mapexIam/src/modules/cache_invalidation/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

var _ ports.CacheInvalidationServicePort = (*CacheInvalidationService)(nil)

// New creates a new CacheInvalidationService. Returns the port interface (Hexagonal Architecture).
func New(deps di.CacheInvalidationServiceDependenciesInjection) ports.CacheInvalidationServicePort {
	return &CacheInvalidationService{deps: deps}
}

// HandleEvent routes a cache invalidation NATS message to the appropriate handler by event type,
// performs the invalidations, and Acks/Nacks the NATS message accordingly. It is the single
// public entry point called by the consumer layer.
func (s *CacheInvalidationService) HandleEvent(msg *natsModel.Message) {
	if err := s.handleCacheInvalidationEvent(msg.Data); err != nil {
		msg.Nack(err)
		return
	}
	msg.Ack()
}
