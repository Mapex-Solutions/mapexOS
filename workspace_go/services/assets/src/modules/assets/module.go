package assets

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assets/application/di"
	"assets/src/modules/assets/application/ports"
	service "assets/src/modules/assets/application/services"
	redisCache "assets/src/modules/assets/infrastructure/cache/redis"
	routerProvider "assets/src/modules/assets/infrastructure/httpclient/router"
	natsProvider "assets/src/modules/assets/infrastructure/messaging/nats"
	minioProvider "assets/src/modules/assets/infrastructure/storage/minio"
	collection "assets/src/modules/assets/infrastructure/persistence/mongo"
	consumers "assets/src/modules/assets/interfaces/message/consumers"
	routes "assets/src/modules/assets/interfaces/http/routes"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// InitRepositories registers the assets repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)
	c.Provide(routerProvider.NewRouteGroupPort)         // Register RouteGroupPort for Router service communication
	c.Provide(minioProvider.NewAssetStoragePort)        // Register AssetStoragePort for object storage operations
	c.Provide(redisCache.NewCacheKeyBuilderAdapter)     // Register CacheKeyBuilderPort for Redis key construction
	// L2 writes retry publisher — feeds the durable fallback stream
	// when synchronous MinIO writes fail; the in-module consumer
	// drains it back against current Mongo state.
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) ports.L2WritesPublisherPort {
		return natsProvider.NewL2WritesPublisherAdapter(params.Bus)
	})
	logger.Info("[MODULE:Assets] Repositories registered")
}

// InitServices registers the assets services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Assets] Services registered")
}

// InitInterfaces registers the assets routes (HTTP)
//
// Routes registered:
//   - External routes (/api/v1/assets): JWT auth for frontend/API clients
//   - Internal routes (/internal/assets): API Key auth for service-to-service communication
//
// Internal routes are used by TieredCache consumers (Router, JS-Executor) as a fallback
// when L2 (MinIO) cache miss occurs. They fetch from MongoDB and repopulate L2.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.AssetServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := configuration.GetIntValue("ctx_timeout")

		// External routes (JWT auth)
		routesV1 := app.Group(
			"/api/v1/assets",

			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		// Register the external routes
		routes.RegisterRoutes(routesV1, service)

		// Internal routes (API Key auth) - for TieredCache fallback
		// Used by: Router, JS-Executor when L2 (MinIO) cache miss occurs
		internalApiKey, _ := configuration.GetStringValue("internal_api_key")
		internalRoutes := app.Group(
			"/internal/assets",

			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)

		// Register the internal routes
		routes.RegisterInternalRoutes(internalRoutes, service)

		// Internal asset-auth routes (API Key auth) — broker plugin L3
		// fallback. Mounted on a sibling group so the broker reads
		// /internal/asset-auth/:uuid without colliding with the read-model
		// endpoint at /internal/assets/:uuid.
		assetAuthInternalRoutes := app.Group(
			"/internal/asset-auth",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		routes.RegisterAssetAuthInternalRoutes(assetAuthInternalRoutes, service)

	}); err != nil {
		log.Fatalf("failed to invoke assets module: %v", err)
	}

	logger.Info("[MODULE:Assets] Routes registered")

	// Start the L2 sync fallback consumer. Drains MAPEXOS-L2-WRITES
	// retry messages and reconciles against current Mongo state when
	// the synchronous MinIO write fails on the happy path.
	if err := c.Invoke(func(params di.AssetConsumerDependenciesInjection) {
		consumers.NewAssetL2SyncConsumer(params.CoreBus, params.AssetService)
	}); err != nil {
		logger.Error(err, "[CONSUMER:AssetL2Sync] Failed to start L2 sync consumer")
	}

	logger.Info("[MODULE:Assets] Consumers registered (NATS Core)")
}
