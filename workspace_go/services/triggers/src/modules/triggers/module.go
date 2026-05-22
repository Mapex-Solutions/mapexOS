package triggers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"triggers/src/modules/triggers/application/ports"
	service "triggers/src/modules/triggers/application/services"
	redisCache "triggers/src/modules/triggers/infrastructure/cache/redis"
	collection "triggers/src/modules/triggers/infrastructure/persistence/mongo"
	routes "triggers/src/modules/triggers/interfaces/http/routes"

	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
)

// InitRepositories registers the triggers repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)
	c.Provide(redisCache.NewKeyBuilder)
	logger.Info("[MODULE:Triggers] Repositories registered")
}

// InitServices registers the triggers services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Triggers] Services registered")
}

// InitInterfaces registers the triggers HTTP routes
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.TriggerServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := configuration.GetIntValue("ctx_timeout")

		// Create route group for triggers API
		routesV1 := app.Group(
			"/api/v1/triggers",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Triggers] Routes registered")

	}); err != nil {
		log.Fatalf("failed to invoke triggers module interfaces: %v", err)
	}
}
