package consumers

import (
	"events/src/modules/events/application/ports"
	"events/src/modules/events/interfaces/message/consumers/events_businessrule"
	"events/src/modules/events/interfaces/message/consumers/events_dlq"
	"events/src/modules/events/interfaces/message/consumers/events_jsexec"
	"events/src/modules/events/interfaces/message/consumers/events_raw"
	"events/src/modules/events/interfaces/message/consumers/events_router"
	"events/src/modules/events/interfaces/message/consumers/events_save"
	"events/src/modules/events/interfaces/message/consumers/events_trigger"
	"events/src/modules/events/interfaces/message/consumers/events_workflow"
	"events/src/modules/events/interfaces/message/consumers/template_invalidate"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel file - exports all consumer initialization functions
 *
 * Following Hexagonal Architecture:
 * - Consumers (Interface Layer) only receive messages and call service
 * - Service (Application Layer) handles all business logic and message lifecycle
 */

// NewEventsSaveConsumer creates the events save consumer
func NewEventsSaveConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_save.NewConsumer(bus, eventService)
}

// NewEventsJsExecConsumer creates the JS executor events consumer
func NewEventsJsExecConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_jsexec.NewConsumer(bus, eventService)
}

// NewEventsRawConsumer creates the raw events consumer
func NewEventsRawConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_raw.NewConsumer(bus, eventService)
}

// NewEventsDLQConsumer creates the DLQ events consumer
func NewEventsDLQConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_dlq.NewConsumer(bus, eventService)
}

// NewEventsRouterConsumer creates the router execution history consumer
func NewEventsRouterConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_router.NewConsumer(bus, eventService)
}

// NewEventsBusinessRuleConsumer creates the business rule execution history consumer
func NewEventsBusinessRuleConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_businessrule.NewConsumer(bus, eventService)
}

// NewEventsTriggerConsumer creates the trigger execution history consumer
func NewEventsTriggerConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_trigger.NewConsumer(bus, eventService)
}

// NewEventsWorkflowConsumer creates the workflow execution history consumer
func NewEventsWorkflowConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	events_workflow.NewConsumer(bus, eventService)
}

// NewTemplateInvalidateConsumer subscribes the events service to the FANOUT
// template-invalidation subject so its TieredCache (L0+L1) stays consistent
// when assets edits/renames an AssetTemplate.
func NewTemplateInvalidateConsumer(bus *natsModel.Bus, service ports.EventServicePort) {
	template_invalidate.NewConsumer(bus, service)
}
