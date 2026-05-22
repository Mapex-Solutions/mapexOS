package archiver

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/archiver/application/ports"
	service "workflow/src/modules/archiver/application/services"
	collection "workflow/src/modules/archiver/infrastructure/persistence/mongo/collection"
	mongoMgrAdapter "workflow/src/modules/archiver/infrastructure/persistence/mongo/manager"
	"workflow/src/bootstrap"
	"workflow/src/modules/archiver/interfaces/message/consumers"
	routes "workflow/src/modules/archiver/interfaces/http/routes"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the archiver repositories in the DIG container.
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)
	c.Provide(mongoMgrAdapter.New)
	logger.Info("[MODULE:Archiver] Repositories registered")
}

// InitServices registers the archiver services in the DIG container.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Archiver] Services registered")
}

// InitInterfaces registers the archiver consumers (NATS) and HTTP routes.
// Consumer references are registered in the ConsumerRegistry (from DIG) for graceful shutdown.
func InitInterfaces() {
	c := container.GetContainer()

	// HTTP Routes
	if err := c.Invoke(func(app *fiber.App, service ports.ArchiverServicePort) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		routesV1 := app.Group(
			"/api/v1/workflow_executions",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Archiver] HTTP routes registered")
	}); err != nil {
		log.Fatalf("failed to invoke archiver HTTP routes: %v", err)
	}

	// NATS Consumers
	if err := c.Invoke(func(params struct {
		container.In
		Bus      *natsModel.Bus             `name:"core"`
		Service  ports.ArchiverServicePort
		Registry *bootstrap.ConsumerRegistry
	}) {
		params.Registry.Register(consumers.NewWorkflowStateConsumer(params.Bus, params.Service))
	}); err != nil {
		log.Fatalf("failed to invoke archiver consumers: %v", err)
	}

	logger.Info("[MODULE:Archiver] Interfaces registered (routes + consumers)")
}
