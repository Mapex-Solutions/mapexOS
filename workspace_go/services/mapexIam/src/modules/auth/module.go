package auth

import (
	"log"

	"github.com/gofiber/fiber/v2"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	ports "mapexIam/src/modules/auth/application/ports"
	services "mapexIam/src/modules/auth/application/services"
	repositories "mapexIam/src/modules/auth/domain/repositories"
	cacheImpl "mapexIam/src/modules/auth/infrastructure/cache/redis"
	lockImpl "mapexIam/src/modules/auth/infrastructure/lock/redis"
	collection "mapexIam/src/modules/auth/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/auth/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeyMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the auth repositories in the DIG container.
// Registers repositories as interfaces (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide AuthRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.AuthRepository {
		return collection.New(m)
	})

	// Provide SessionRepository (returns interface, cache-based implementation)
	c.Provide(func(cache common.AppCache) repositories.SessionRepository {
		return cacheImpl.New(cache)
	})

	// Provide LockManagerPort adapter (wraps redisLock.LockManager driver so
	// the application DI layer depends on a port, not the concrete driver).
	c.Provide(lockImpl.NewLockManagerAdapter)

	// Provide AuthorizationCacheRepository (returns interface, uses DI struct pattern)
	c.Provide(cacheImpl.NewAuthorizationCacheRepository)

	// Provide CoverageCacheRepository (returns interface, uses DI struct pattern)
	c.Provide(cacheImpl.NewCoverageCacheRepository)

	logger.Info("[MODULE:Auth] Repositories registered (MongoDB + Session + AuthCache + CoverageCache)")
}

// InitServices registers the auth services in the DIG container.
// Registers services as ports/interfaces for dependency injection following Hexagonal Architecture.
func InitServices() {
	c := container.GetContainer()

	// Register AuthService (returns AuthServicePort)
	c.Provide(services.New)

	logger.Info("[MODULE:Auth] Services registered (AuthService)")
}

// InitInterfaces registers the auth HTTP routes.
// Handlers receive AuthServicePort interface (not concrete implementation).
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(
		app *fiber.App,
		service ports.AuthServicePort,
		authCacheRepo repositories.AuthorizationCacheRepository,
		coverageCacheRepo repositories.CoverageCacheRepository,
	) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Public routes
		routesV1 := app.Group("/auth", ctxInjector.ContextInjector(ctxTimeout))
		routes.RegisterRoutes(routesV1, service)

		// Internal routes (protected with API Key)
		internalApiKey, err := config.GetStringValue("internal_api_key")
		if err != nil {
			log.Fatalf("internal_api_key not configured: %v", err)
		}

		internalRoutesV1 := app.Group(
			"/internal/auth",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeyMw.ApiKeyAuthMiddleware(internalApiKey),
		)
		routes.RegisterInternalRoutes(internalRoutesV1, authCacheRepo, coverageCacheRepo)

		logger.Info("[MODULE:Auth] Routes registered (public + internal)")

	}); err != nil {
		log.Fatalf("failed to invoke auth module interfaces: %v", err)
	}
}
