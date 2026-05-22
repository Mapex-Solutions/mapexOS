package definitions

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/definitions/application/ports"
	service "workflow/src/modules/definitions/application/services"
	minioProvider "workflow/src/modules/definitions/infrastructure/storage/minio"
	collection "workflow/src/modules/definitions/infrastructure/persistence/mongo"
	routes "workflow/src/modules/definitions/interfaces/http/routes"

	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// InitRepositories registers the definitions repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)
	c.Provide(minioProvider.NewDefinitionStoragePort)
	logger.Info("[MODULE:Definitions] Repositories registered")
}

// InitServices registers the definitions services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Definitions] Services registered")
}

// InitInterfaces registers the definitions routes (HTTP)
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.DefinitionServicePort) {
		// Set default timeout for this router
		ctxTimeout, _ := configuration.GetIntValue("ctx_timeout")

		// External routes (JWT auth)
		routesV1 := app.Group(
			"/api/v1/workflow_definitions",

			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		// Register the external routes
		routes.RegisterRoutes(routesV1, service)

		// Internal routes (API Key auth) — for TieredCache fallback
		// Used by: js-workflow-executor when L2 (MinIO) cache miss occurs
		internalApiKey, _ := configuration.GetStringValue("internal_api_key")
		internalRoutes := app.Group(
			"/internal/workflow-scripts",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)

		routes.RegisterInternalRoutes(internalRoutes, service)

	}); err != nil {
		log.Fatalf("failed to invoke definitions module: %v", err)
	}

	logger.Info("[MODULE:Definitions] Routes registered")
}
