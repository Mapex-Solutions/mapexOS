package di

import (
	"triggers/src/modules/triggers/application/ports"
	"triggers/src/modules/triggers/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	"go.uber.org/dig"
)

// TriggerServiceDependenciesInjection aggregates all dependencies required
// by the TriggerService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - TriggerRepository: Repository for trigger persistence operations
//   - CacheRepository: Repository for caching trigger data (GetOrSetEx pattern)
//   - AppCache: Service-private cache for counter and general cache operations
//   - CacheKeyBuilder: Port that produces Redis cache keys for the module
//
// The dig.In tag enables automatic dependency injection by the dig container.
type TriggerServiceDependenciesInjection struct {
	dig.In
	TriggerRepository repositories.TriggerRepository
	CacheRepository   repositories.CacheRepository
	AppCache          common.AppCache
	CacheKeyBuilder   ports.TriggerCacheKeyBuilderPort
}
