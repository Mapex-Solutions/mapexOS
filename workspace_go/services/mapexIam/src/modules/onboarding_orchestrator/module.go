package onboarding_orchestrator

import (
	"log"

	"github.com/gofiber/fiber/v2"

	ports "mapexIam/src/modules/onboarding_orchestrator/application/ports"
	service "mapexIam/src/modules/onboarding_orchestrator/application/services"
	mongoAdapter "mapexIam/src/modules/onboarding_orchestrator/infrastructure/persistence/mongo"
	routes "mapexIam/src/modules/onboarding_orchestrator/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the onboarding application services in the DIG container
// This module is an Application Service that orchestrates User and Membership services
// to provide atomic user onboarding with organization/role assignments.
func InitServices() {
	c := container.GetContainer()

	// Provide MongoManagerPort adapter (wraps mongoManager.MongoManager driver
	// so the application DI layer depends on a port, not the concrete driver).
	c.Provide(mongoAdapter.NewMongoManagerAdapter)

	c.Provide(service.New)
	logger.Info("[MODULE:Onboarding] Services registered")
}

// InitInterfaces registers the onboarding HTTP routes
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.UserOnboardingServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Create route group for onboarding endpoints
		routesV1 := app.Group(
			"/api/v1/onboarding",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		// Register onboarding routes
		routes.RegisterRoutes(routesV1, service)

		logger.Info("[MODULE:Onboarding] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke onboarding module interfaces: %v", err)
	}
}
