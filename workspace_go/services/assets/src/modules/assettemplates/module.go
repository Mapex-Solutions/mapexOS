package assettemplates

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/application/ports"
	service "assets/src/modules/assettemplates/application/services"
	redisCache "assets/src/modules/assettemplates/infrastructure/cache/redis"
	minioProvider "assets/src/modules/assettemplates/infrastructure/storage/minio"
	collection "assets/src/modules/assettemplates/infrastructure/persistence/mongo"
	consumers "assets/src/modules/assettemplates/interfaces/message/consumers"
	routes "assets/src/modules/assettemplates/interfaces/http/routes"

	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
)

// InitRepositories registers the assettemplates repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)
	c.Provide(minioProvider.NewTemplateStoragePort)     // Register TemplateStoragePort for script storage
	c.Provide(redisCache.NewCacheKeyBuilderAdapter)     // Register CacheKeyBuilderPort for Redis key construction
	logger.Info("[MODULE:AssetTemplates] Repositories registered")
}

// InitServices registers the assettemplates services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:AssetTemplates] Services registered")
}

// InitInterfaces registers the assettemplates routes (HTTP) and consumers (NATS Core)
func InitInterfaces() {
	c := container.GetContainer()

	// Register HTTP routes
	if err := c.Invoke(func(app *fiber.App, service ports.AssetTemplateServicePort) {

		// Set default timeot for this router
		ctxTimeout, _ := configuration.GetIntValue("ctx_timeout")

		// External routes (JWT auth)
		routesV1 := app.Group(
			"/api/v1/asset_templates",

			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		// Register the routes and handlers
		routes.RegisterRoutes(routesV1, service)

		// Internal routes (API Key auth) - for TieredCache fallback
		internalApiKey, _ := configuration.GetStringValue("internal_api_key")
		internalRoutes := app.Group(
			"/internal/templates",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		routes.RegisterInternalRoutes(internalRoutes, service)

	}); err != nil {
		log.Fatalf("failed to invoke assettemplates module: %v", err)
	}

	logger.Info("[MODULE:AssetTemplates] Routes registered")

	// Start list name updated consumer for denormalized name synchronization
	// Uses NATS Core connection (port 4222) for JetStream streams
	if err := c.Invoke(func(params di.AssetTemplateConsumerDependenciesInjection) {
		consumers.NewListNameUpdatedConsumer(params.CoreBus, params.AssetTemplateService)
	}); err != nil {
		logger.Error(err, "[CONSUMER:ListNameUpdated] Failed to start list name updated consumer")
	}

	logger.Info("[MODULE:AssetTemplates] Consumers registered (NATS Core)")
}
