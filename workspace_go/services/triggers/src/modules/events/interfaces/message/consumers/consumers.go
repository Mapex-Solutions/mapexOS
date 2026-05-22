package consumers

import (
	"triggers/src/modules/events/application/ports"
	"triggers/src/modules/events/interfaces/message/consumers/plugin_execute"
	"triggers/src/modules/events/interfaces/message/consumers/trigger_execute"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel file - exports all consumer initialization functions
 *
 * Following Hexagonal Architecture:
 * - Consumers (Interface Layer) only receive messages and call service
 * - Service (Application Layer) handles all business logic and message lifecycle
 */

// NewTriggerExecuteConsumer creates the trigger execution consumer (from Router Service)
func NewTriggerExecuteConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	trigger_execute.NewConsumer(bus, eventService)
}

// NewPluginExecuteConsumer creates the plugin execution consumer (from Workflow Service)
func NewPluginExecuteConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	plugin_execute.NewConsumer(bus, eventService)
}
