package cache_invalidation

import (
	"mapexIam/src/modules/cache_invalidation/application/ports"
	"mapexIam/src/modules/cache_invalidation/application/services"
	"mapexIam/src/modules/cache_invalidation/interfaces/message/consumers"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the cache invalidation service in the DIG container.
// Service is registered as CacheInvalidationServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with CacheInvalidationServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide CacheInvalidationService (returns CacheInvalidationServicePort interface)
	c.Provide(services.New)

	logger.Info("[MODULE:CacheInvalidation] Services registered (with CacheInvalidationServicePort)")
}

// InitListeners initializes NATS event listeners for cache invalidation.
// Uses Phase 4 (InitListeners) to ensure all repositories and services from
// other modules are already registered in the DI container.
func InitListeners() {
	c := container.GetContainer()

	err := c.Invoke(func(
		natsBus *natsModel.Bus,
		service ports.CacheInvalidationServicePort,
	) {
		consumers.NewCacheInvalidationConsumer(natsBus, service)

		logger.Info("[MODULE:CacheInvalidation] Cache invalidation consumer started successfully")
	})

	if err != nil {
		logger.Error(err, "[MODULE:CacheInvalidation] Failed to invoke cache invalidation consumer initialization")
	}
}
