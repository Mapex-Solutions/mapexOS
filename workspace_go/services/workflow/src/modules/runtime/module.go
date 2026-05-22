package runtime

import (
	"log"

	defRepos "workflow/src/modules/definitions/domain/repositories"
	"workflow/src/modules/runtime/application/ports"
	service "workflow/src/modules/runtime/application/services"
	defCache "workflow/src/modules/runtime/infrastructure/cache"
	natsMessaging "workflow/src/modules/runtime/infrastructure/messaging/nats"
	natsPersistence "workflow/src/modules/runtime/infrastructure/persistence/nats"
	"workflow/src/bootstrap"
	"workflow/src/modules/runtime/interfaces/message/consumers"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the runtime repositories in the DIG container.
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(natsPersistence.NewExecutionStateRepository)
	logger.Info("[MODULE:Runtime] Repositories registered")
}

// InitServices registers the runtime services in the DIG container.
func InitServices() {
	c := container.GetContainer()
	c.Provide(natsMessaging.NewRuntimePublisher)

	// DefinitionLoader: TieredCache (name:"definitions") + MongoDB DefinitionRepository
	c.Provide(func(params struct {
		container.In
		Cache common.TieredCache `name:"definitions"`
		Repo  defRepos.DefinitionRepository
	}) ports.DefinitionLoaderPort {
		return defCache.New(params.Cache, params.Repo)
	})

	c.Provide(service.New)
	logger.Info("[MODULE:Runtime] Services registered")
}

// InitInterfaces registers the runtime consumers (NATS).
// Consumer references are registered in the ConsumerRegistry (from DIG) for graceful shutdown.
func InitInterfaces() {
	c := container.GetContainer()

	// NATS Consumers
	if err := c.Invoke(func(params struct {
		container.In
		Bus      *natsModel.Bus `name:"core"`
		Service  ports.RuntimeServicePort
		Registry *bootstrap.ConsumerRegistry
	}) {
		params.Registry.Register(consumers.NewWorkflowResumeConsumer(params.Bus, params.Service))
		params.Registry.Register(consumers.NewWorkflowExecutionConsumer(params.Bus, params.Service))
		params.Registry.Register(consumers.NewScheduleFireConsumer(params.Bus, params.Service))
	}); err != nil {
		log.Fatalf("failed to invoke runtime consumers: %v", err)
	}

	logger.Info("[MODULE:Runtime] Interfaces registered (consumers)")
}
