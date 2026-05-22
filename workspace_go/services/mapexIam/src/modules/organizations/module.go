package organizations

import (
	"log"

	"github.com/gofiber/fiber/v2"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	ports "mapexIam/src/modules/organizations/application/ports"
	services "mapexIam/src/modules/organizations/application/services"
	repositories "mapexIam/src/modules/organizations/domain/repositories"
	collection "mapexIam/src/modules/organizations/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/organizations/interfaces/http/routes"
	internalRoutes "mapexIam/src/modules/organizations/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the organizations repositories in the DIG container.
// Registers repository as interface (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide OrganizationRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.OrganizationRepository {
		return collection.New(m)
	})

	logger.Info("[MODULE:Organizations] Repositories registered")
}

// InitServices registers the organizations services in the DIG container.
// Service is registered as OrganizationServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with OrganizationServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide OrganizationService (returns OrganizationServicePort interface)
	// DIG automatically injects OrganizationServiceDependenciesInjection struct with all dependencies
	c.Provide(services.New)

	logger.Info("[MODULE:Organizations] Services registered (with OrganizationServicePort)")
}

// InitInterfaces registers the organizations HTTP routes.
// Handlers receive OrganizationServicePort interface (not concrete implementation).
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.OrganizationServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Set routes and versions
		routesV1 := app.Group(
			"/api/v1/organizations",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(routesV1, service)

		// Register internal routes (for inter-service communication)
		internalRoutes.RegisterInternalRoutes(app, service)

		logger.Info("[MODULE:Organizations] Routes registered (public + internal)")

	}); err != nil {
		log.Fatalf("failed to invoke organizations module interfaces: %v", err)
	}
}
