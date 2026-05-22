package roles

import (
	"log"

	"github.com/gofiber/fiber/v2"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	ports "mapexIam/src/modules/roles/application/ports"
	services "mapexIam/src/modules/roles/application/services"
	repositories "mapexIam/src/modules/roles/domain/repositories"
	collection "mapexIam/src/modules/roles/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/roles/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the roles repositories in the DIG container.
// Registers repository as interface (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide RoleRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.RoleRepository {
		return collection.New(m)
	})

	logger.Info("[MODULE:Roles] Repositories registered")
}

// InitServices registers the roles services in the DIG container.
// Service is registered as RoleServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with RoleServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide RoleService (returns RoleServicePort interface)
	// DIG automatically injects RoleServiceDependenciesInjection struct with all dependencies
	// Cache invalidation is now handled by shared authorization_cache module
	c.Provide(services.New)

	logger.Info("[MODULE:Roles] Services registered (with RoleServicePort)")
}

// InitInterfaces registers the roles HTTP routes
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.RoleServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Set routes and versions
		routesV1 := app.Group(
			"/api/v1/roles",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(routesV1, service)

		logger.Info("[MODULE:Roles] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke roles module interfaces: %v", err)
	}
}
