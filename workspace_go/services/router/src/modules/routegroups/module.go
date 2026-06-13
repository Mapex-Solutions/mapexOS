package routegroups

import (
	"fmt"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"router/src/modules/routegroups/application/ports"
	service "router/src/modules/routegroups/application/services"
	collection "router/src/modules/routegroups/infrastructure/persistence/mongo"
	routes "router/src/modules/routegroups/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
)

// InitRepositories registers the routegroups repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)
	logger.Info("[MODULE:Routegroups] Repositories registered")
}

// InitServices registers the routegroups services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Routegroups] Services registered")
}

// InitInterfaces registers the routegroups HTTP routes
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *web.App, service ports.RouteGroupServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Create route group
		routesV1 := app.Group(
			"/api/v1/route_groups",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Routegroups] Routes registered")

		// Internal routes (API Key authentication for MS-to-MS communication)
		apiKey, _ := config.GetStringValue("internal_api_key")
		internalRoutesV1 := app.Group(
			"/internal/route_groups",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(apiKey),
		)

		routes.RegisterInternalRoutes(internalRoutesV1, service)
		logger.Info("[MODULE:Routegroups] Internal routes registered")

	}); err != nil {
		logger.Panic(fmt.Sprintf("[MODULE:Routegroups] failed to invoke module: %v", err))
	}
}
