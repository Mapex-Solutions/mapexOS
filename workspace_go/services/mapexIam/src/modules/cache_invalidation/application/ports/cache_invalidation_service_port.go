package ports

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// CacheInvalidationServicePort defines the contract for handling cache invalidation events.
//
// The consumer (interface layer) receives the raw NATS message and delegates to `HandleEvent`.
// The service is responsible for routing by event type and invalidating the appropriate caches
// (auth cache, coverage cache) through its driven repositories.
type CacheInvalidationServicePort interface {
	// HandleEvent routes a cache invalidation message to the appropriate handler by event type,
	// performs the invalidations, and Acks/Nacks the NATS message accordingly.
	HandleEvent(msg *natsModel.Message)
}
