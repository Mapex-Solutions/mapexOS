package users

import (
	"log"

	"github.com/gofiber/fiber/v2"

	ports "mapexIam/src/modules/users/application/ports"
	services "mapexIam/src/modules/users/application/services"
	repositories "mapexIam/src/modules/users/domain/repositories"
	counterCache "mapexIam/src/modules/users/infrastructure/cache/redis"
	collection "mapexIam/src/modules/users/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/users/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the users repositories in the DIG container.
// Registers repository as interface (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide UserRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.UserRepository {
		return collection.New(m)
	})

	logger.Info("[MODULE:Users] Repositories registered (MongoDB)")
}

// InitServices registers the users services in the DIG container.
// Service is registered as UserServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with UserServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide CounterCachePort adapter (infrastructure implementation injected
	// into UserService as a port, preserving application → infrastructure
	// dependency direction).
	c.Provide(counterCache.NewCounterCacheAdapter)

	// Provide UserService (returns UserServicePort interface)
	// DIG automatically injects UserServiceDependenciesInjection struct with all dependencies
	c.Provide(services.New)

	logger.Info("[MODULE:Users] Services registered (returns UserServicePort)")
}

// InitInterfaces registers the users HTTP routes.
// Handlers receive UserServicePort interface (not concrete implementation).
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.UserServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Create a new group for the users routes
		routesV1 := app.Group(
			"/api/v1/users",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Users] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke users module interfaces: %v", err)
	}
}
