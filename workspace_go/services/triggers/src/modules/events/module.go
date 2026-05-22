package events

import (
	service "triggers/src/modules/events/application/services"
	registry "triggers/src/modules/events/infrastructure/registry"
	consumers "triggers/src/modules/events/interfaces/message/consumers"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the events services in the DIG container.
//
// Following Hexagonal Architecture:
// Register infrastructure adapters FIRST (ExecutorRegistry)
// Then register application services (EventService)
//
// The DI container will automatically inject ExecutorRegistry into EventService
// via the EventServiceDependenciesInjection struct.
func InitServices() {
	c := container.GetContainer()

	// Register infrastructure layer: ExecutorRegistry adapter
	// This provides concrete implementations of trigger executors (HTTP, Email, Slack, Teams, etc.)
	c.Provide(registry.NewExecutorRegistry)

	// Register application layer: EventService
	// Depends on ExecutorRegistry (port interface) injected by DI
	c.Provide(service.New)

	logger.Info("[MODULE:Events] Services registered (ExecutorRegistry + EventService)")
}

// InitInterfaces registers the events consumers (NATS)
func InitInterfaces() {
	c := container.GetContainer()

	// Start trigger execution consumer (from Router Service)
	if err := c.Invoke(consumers.NewTriggerExecuteConsumer); err != nil {
		logger.Error(err, "[CONSUMER:TriggerExecute] Failed to start trigger execution consumer")
	}

	// Start plugin execution consumer (from Workflow Service)
	if err := c.Invoke(consumers.NewPluginExecuteConsumer); err != nil {
		logger.Error(err, "[CONSUMER:PluginExecute] Failed to start plugin execution consumer")
	}

	logger.Info("[MODULE:Events] Consumers registered")
}
