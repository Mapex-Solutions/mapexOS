package authorization_cache

import (
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"mapexIam/src/modules/authorization_cache/domain/repositories"
	cacheImpl "mapexIam/src/modules/authorization_cache/infrastructure/cache/redis"
)

// InitRepositories registers the authorization cache repository in the DIG container.
// This is a SHARED module - used by multiple modules (memberships, roles, groups, organizations).
//
// Architecture Pattern: Dependency Injection
//   - Registers repository as interface (Hexagonal Architecture)
//   - Returns repositories.AuthCacheRepository (port/interface)
//   - Implementation is hidden (infrastructure detail)
//
// Usage in other modules:
//   import authCacheRepos "mapexIam/src/modules/authorization_cache/domain/repositories"
//
//   type XXXServiceDependenciesInjection struct {
//       dig.In
//       AuthCacheRepo authCacheRepos.AuthCacheRepository
//   }
func InitRepositories() {
	c := container.GetContainer()

	// Provide AuthCacheRepository (returns interface, Redis implementation)
	// IMPORTANT: Uses SharedCache (DB 5) for cross-service authorization data
	c.Provide(func(sharedCache common.SharedCache) repositories.AuthCacheRepository {
		return cacheImpl.New(sharedCache)
	})

	logger.Info("[MODULE:AuthorizationCache] Repositories registered")
}
