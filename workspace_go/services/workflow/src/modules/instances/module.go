package instances

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/instances/application/ports"
	service "workflow/src/modules/instances/application/services"
	cacheLoader "workflow/src/modules/instances/infrastructure/cache"
	"workflow/src/modules/instances/domain/repositories"
	collection "workflow/src/modules/instances/infrastructure/persistence/mongo"
	routes "workflow/src/modules/instances/interfaces/http/routes"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// InitRepositories registers the instances repositories in the DIG container.
func InitRepositories() {
	c := container.GetContainer()

	// MongoDB repository
	c.Provide(collection.New)

	// InstanceLoader wraps TieredCache (L0→L1) + MongoDB fallback
	c.Provide(func(params struct {
		container.In
		Cache common.TieredCache              `name:"instances"`
		Repo  repositories.InstanceRepository
	}) ports.InstanceLoaderPort {
		return cacheLoader.New(params.Cache, params.Repo)
	})

	logger.Info("[MODULE:Instances] Repositories registered")
}

// InitServices registers the instances services in the DIG container.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Instances] Services registered")
}

// InitInterfaces registers the instances routes (HTTP).
func InitInterfaces() {
	c := container.GetContainer()

	// HTTP Routes
	if err := c.Invoke(func(app *fiber.App, service ports.InstancesServicePort) {
		ctxTimeout, _ := configuration.GetIntValue("ctx_timeout")

		routesV1 := app.Group(
			"/api/v1/workflow_instances",

			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)

	}); err != nil {
		log.Fatalf("failed to invoke instances HTTP routes: %v", err)
	}

	logger.Info("[MODULE:Instances] Interfaces registered (routes)")
}
