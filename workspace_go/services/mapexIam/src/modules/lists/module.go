package lists

import (
	"log"

	"github.com/gofiber/fiber/v2"

	ports "mapexIam/src/modules/lists/application/ports"
	services "mapexIam/src/modules/lists/application/services"
	repositories "mapexIam/src/modules/lists/domain/repositories"
	collection "mapexIam/src/modules/lists/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/lists/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the lists repositories in the DIG container.
// Registers repository as interface (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide ListRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.ListRepository {
		return collection.New(m)
	})

	logger.Info("[MODULE:Lists] Repositories registered")
}

// InitServices registers the lists services in the DIG container.
// Service is registered as ListServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with ListServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide ListService (returns ListServicePort interface)
	// DIG automatically injects ListServiceDependenciesInjection struct with all dependencies
	c.Provide(services.New)

	logger.Info("[MODULE:Lists] Services registered (returns ListServicePort)")
}

// InitInterfaces registers the lists HTTP routes.
// Handlers receive ListServicePort interface (not concrete implementation).
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.ListServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Create a new group for the lists routes
		routesV1 := app.Group(
			"/api/v1/lists",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Lists] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke lists module interfaces: %v", err)
	}
}
