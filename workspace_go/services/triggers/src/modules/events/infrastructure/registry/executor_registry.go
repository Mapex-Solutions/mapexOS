package registry

import (
	"triggers/src/modules/events/application/ports"
	"triggers/src/modules/events/infrastructure/communications/email"
	"triggers/src/modules/events/infrastructure/communications/slack"
	"triggers/src/modules/events/infrastructure/communications/teams"
	"triggers/src/modules/events/infrastructure/technical/http"
	"triggers/src/modules/events/infrastructure/technical/mqtt"
	natsexec "triggers/src/modules/events/infrastructure/technical/nats"
	"triggers/src/modules/events/infrastructure/technical/rabbitmq"
	wsexec "triggers/src/modules/events/infrastructure/technical/websocket"
)

// NewExecutorRegistry creates and initializes the executor registry.
//
// This is the infrastructure factory that knows about all concrete executor implementations.
// It will be registered in the DI container and injected into the EventService.
//
// Following the Factory pattern, the registry:
// - Instantiates all available executors
// - Registers them by trigger type
// - Returns the registry interface (port), not the concrete type
//
// Returns:
//   - ports.ExecutorRegistry: The registry interface implementation
func NewExecutorRegistry() ports.ExecutorRegistry {
	registry := &executorRegistry{
		executors: make(map[string]ports.TriggerExecutor),
	}

	// Register all available executor adapters
	// Each executor implements the ports.TriggerExecutor interface

	// Technical executors
	registry.register(http.NewHTTPExecutor())
	registry.register(mqtt.NewMQTTExecutor())
	registry.register(rabbitmq.NewRabbitMQExecutor())
	registry.register(natsexec.NewNATSExecutor())
	registry.register(wsexec.NewWebSocketExecutor())

	// Communication executors
	registry.register(email.NewEmailExecutor())
	registry.register(teams.NewTeamsExecutor())
	registry.register(slack.NewSlackExecutor())

	return registry
}

// register adds an executor to the internal registry map.
//
// This is a private helper method used during initialization.
//
// Parameters:
//   - executor: The executor adapter to register
func (r *executorRegistry) register(executor ports.TriggerExecutor) {
	r.executors[executor.GetType()] = executor
}

// GetExecutor retrieves an executor by trigger type.
//
// This implements the ports.ExecutorRegistry interface method.
//
// Parameters:
//   - triggerType: The type of trigger (http, email, teams, slack, mqtt, rabbitmq, nats, websocket)
//
// Returns:
//   - ports.TriggerExecutor: The executor adapter for the trigger type
//   - bool: True if executor exists, false otherwise
func (r *executorRegistry) GetExecutor(triggerType string) (ports.TriggerExecutor, bool) {
	executor, exists := r.executors[triggerType]
	return executor, exists
}
