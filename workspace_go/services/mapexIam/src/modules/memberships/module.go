package memberships

import (
	"log"

	"github.com/gofiber/fiber/v2"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	ports "mapexIam/src/modules/memberships/application/ports"
	services "mapexIam/src/modules/memberships/application/services"
	repositories "mapexIam/src/modules/memberships/domain/repositories"
	collection "mapexIam/src/modules/memberships/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/memberships/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the memberships repositories in the DIG container.
// Registers repository as interface (Hexagonal Architecture).
func InitRepositories() {
	c := container.GetContainer()

	// Provide MembershipRepository (returns interface, MongoDB implementation)
	c.Provide(func(m *mongoManager.MongoManager) repositories.MembershipRepository {
		return collection.New(m)
	})

	logger.Info("[MODULE:Memberships] Repositories registered")
}

// InitServices registers the memberships services in the DIG container.
// Service is registered as MembershipServicePort interface (Hexagonal Architecture).
// Uses dependency injection pattern with MembershipServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Provide MembershipService (returns MembershipServicePort interface)
	// DIG automatically injects MembershipServiceDependenciesInjection struct with all dependencies
	// Cache invalidation is now handled by shared authorization_cache module
	c.Provide(services.New)

	logger.Info("[MODULE:Memberships] Services registered (with MembershipServicePort)")
}

// InitInterfaces registers the memberships HTTP routes
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.MembershipServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Register /api/v1/memberships routes
		routesV1 := app.Group(
			"/api/v1/memberships",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(routesV1, service)

		// Register /api/v1/me routes with authentication
		meRoutes := app.Group(
			"/api/v1/me",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterMeRoutes(meRoutes, service)

		logger.Info("[MODULE:Memberships] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke memberships module interfaces: %v", err)
	}
}
