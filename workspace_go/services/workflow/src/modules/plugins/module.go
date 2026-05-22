package plugins

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/plugins/application/ports"
	service "workflow/src/modules/plugins/application/services"
	"workflow/src/modules/plugins/domain/repositories"
	cacheLoader "workflow/src/modules/plugins/infrastructure/cache"
	collection "workflow/src/modules/plugins/infrastructure/persistence/mongo"
	pluginFanout "workflow/src/modules/plugins/interfaces/message/consumers/plugin_fanout"
	routes "workflow/src/modules/plugins/interfaces/http/routes"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// InitRepositories registers the plugins repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()

	// MongoDB repository
	c.Provide(collection.New)

	// PluginLoader wraps TieredCache (L0→L1) + MongoDB fallback
	c.Provide(func(params struct {
		container.In
		Cache common.TieredCache                    `name:"plugins"`
		Repo  repositories.PluginManifestRepository
	}) ports.PluginLoaderPort {
		return cacheLoader.New(params.Cache, params.Repo)
	})

	logger.Info("[MODULE:Plugins] Repositories registered")
}

// InitServices registers the plugins services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Plugins] Services registered")
}

// InitInterfaces registers the plugins routes (HTTP) and NATS fanout subscription
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(params struct {
		container.In
		App     *fiber.App
		Service ports.PluginServicePort
		NatsBus natsModel.Fanout `name:"core"`
	}) {
		// Set default timeout for this router
		ctxTimeout, _ := configuration.GetIntValue("ctx_timeout")

		// External routes (JWT auth)
		routesV1 := params.App.Group(
			"/api/v1/plugins",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, params.Service)

		// Subscribe to FANOUT for cross-pod cache invalidation
		pluginFanout.NewConsumer(params.NatsBus, params.Service)

	}); err != nil {
		log.Fatalf("failed to invoke plugins module: %v", err)
	}

	logger.Info("[MODULE:Plugins] Routes registered + FANOUT subscription active")
}
