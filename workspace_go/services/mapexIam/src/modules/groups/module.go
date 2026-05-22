package groups

import (
	"log"

	"github.com/gofiber/fiber/v2"

	ports "mapexIam/src/modules/groups/application/ports"
	services "mapexIam/src/modules/groups/application/services"
	repositories "mapexIam/src/modules/groups/domain/repositories"
	counterCache "mapexIam/src/modules/groups/infrastructure/cache/redis"
	collection "mapexIam/src/modules/groups/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/groups/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the groups repositories in the DIG container.
// Registers repository as interface (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide GroupRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.GroupRepository {
		return collection.New(m)
	})

	// Provide GroupMemberRepository (junction table for scalable member management)
	c.Provide(func(m *mongoManager.MongoManager) repositories.GroupMemberRepository {
		return collection.NewGroupMemberRepository(m)
	})

	logger.Info("[MODULE:Groups] Repositories registered (Group + GroupMember)")
}

// InitServices registers the groups services in the DIG container.
// Service is registered as GroupServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with GroupServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide CounterCachePort adapter (infrastructure implementation injected
	// into GroupService as a port, preserving application → infrastructure
	// dependency direction).
	c.Provide(counterCache.NewCounterCacheAdapter)

	// Provide GroupQueryService (returns GroupQueryServicePort interface)
	// Query-only service for cross-domain consumption - no circular dependencies
	c.Provide(services.NewGroupQueryService)

	// Provide GroupService (returns GroupServicePort interface)
	// DIG automatically injects GroupServiceDependenciesInjection struct with all dependencies
	c.Provide(services.New)

	logger.Info("[MODULE:Groups] Services registered (with GroupServicePort + GroupQueryServicePort)")
}

// InitInterfaces registers the groups HTTP routes.
// Handlers receive GroupServicePort interface (not concrete implementation).
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.GroupServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Set routes and versions
		routesV1 := app.Group(
			"/api/v1/groups",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(routesV1, service)

		logger.Info("[MODULE:Groups] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke groups module interfaces: %v", err)
	}
}
