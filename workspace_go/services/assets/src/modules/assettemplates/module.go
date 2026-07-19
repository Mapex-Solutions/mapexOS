package assettemplates

import (
	"log"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	assetsPorts "assets/src/modules/assets/application/ports"
	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/application/ports"
	service "assets/src/modules/assettemplates/application/services"
	adapters "assets/src/modules/assettemplates/infrastructure/adapters"
	redisCache "assets/src/modules/assettemplates/infrastructure/cache/redis"
	listsclient "assets/src/modules/assettemplates/infrastructure/httpclient/listsclient"
	marketplaceclient "assets/src/modules/assettemplates/infrastructure/httpclient/marketplaceclient"
	natsAdapter "assets/src/modules/assettemplates/infrastructure/messaging/nats"
	collection "assets/src/modules/assettemplates/infrastructure/persistence/mongo"
	minioProvider "assets/src/modules/assettemplates/infrastructure/storage/minio"
	routes "assets/src/modules/assettemplates/interfaces/http/routes"
	consumers "assets/src/modules/assettemplates/interfaces/message/consumers"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
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
	c.Provide(collection.NewFieldVocabulary)              // Register FieldVocabularyRepository for the authoring UI vocabulary
	c.Provide(minioProvider.NewTemplateStoragePort)       // Register TemplateStoragePort for script storage
	c.Provide(redisCache.NewCacheKeyBuilderAdapter)       // Register CacheKeyBuilderPort for Redis key construction
	c.Provide(collection.NewMigrationPlanRepository)      // Register MigrationPlanRepository for template migration plans
	c.Provide(collection.NewMigrationExecutionRepository) // Register MigrationExecutionRepository for per-asset executions
	logger.Info("[MODULE:AssetTemplates] Repositories registered")
}

// InitServices registers the assettemplates services in the DIG container
func InitServices() {
	c := container.GetContainer()

	// Migration start-timer scheduler over the core NATS ScheduleManager.
	c.Provide(func(params struct {
		container.In
		SM natsModel.ScheduleManager `name:"core"`
	}) ports.MigrationSchedulerPort {
		return natsAdapter.NewMigrationScheduler(params.SM)
	})

	// Template switcher — rebinds an asset to a target template via the assets service.
	c.Provide(func(svc assetsPorts.AssetServicePort) ports.TemplateSwitcherPort {
		return adapters.NewTemplateSwitcherAdapter(svc)
	})

	// Asset usage — counts assets referencing a template (marketplace uninstall guard) via the assets service.
	c.Provide(func(svc assetsPorts.AssetServicePort) ports.AssetUsagePort {
		return adapters.NewAssetUsageAdapter(svc)
	})

	// Marketplace catalog client (fetch bundle) + mapexIam lists client (resolve
	// org-scoped classification) for the install flow.
	c.Provide(marketplaceclient.NewMarketplaceClient)
	c.Provide(listsclient.NewListsClient)

	c.Provide(service.New)
	logger.Info("[MODULE:AssetTemplates] Services registered")
}

// InitInterfaces registers the assettemplates routes (HTTP) and consumers (NATS Core)
func InitInterfaces() {
	c := container.GetContainer()

	// Register HTTP routes
	if err := c.Invoke(func(app *web.App, service ports.AssetTemplateServicePort) {

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

	// Start migration timers consumer for fired plan start-timer messages.
	if err := c.Invoke(func(params di.AssetTemplateConsumerDependenciesInjection) {
		consumers.NewMigrationTimersConsumer(params.CoreBus, params.AssetTemplateService)
	}); err != nil {
		logger.Error(err, "[CONSUMER:MigrationTimers] Failed to start migration timers consumer")
	}

	logger.Info("[MODULE:AssetTemplates] Consumers registered (NATS Core)")
}
