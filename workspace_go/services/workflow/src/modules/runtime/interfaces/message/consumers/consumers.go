package consumers

import (
	"workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/interfaces/message/consumers/schedule_fire"
	"workflow/src/modules/runtime/interfaces/message/consumers/workflow_execution"
	"workflow/src/modules/runtime/interfaces/message/consumers/workflow_resume"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel file - exports all consumer initialization functions
 *
 * Following Hexagonal Architecture:
 * - Consumers (Interface Layer) only receive messages and call service
 * - Service (Application Layer) handles all business logic and message lifecycle
 */

// NewWorkflowResumeConsumer creates the workflow resume consumer
func NewWorkflowResumeConsumer(bus *natsModel.Bus, service ports.RuntimeServicePort) *natsModel.Consumer {
	return workflow_resume.NewConsumer(bus, service)
}

// NewWorkflowExecutionConsumer creates the workflow execution consumer (3-mode dispatch)
func NewWorkflowExecutionConsumer(bus *natsModel.Bus, service ports.RuntimeServicePort) *natsModel.Consumer {
	return workflow_execution.NewConsumer(bus, service)
}

// NewScheduleFireConsumer creates the schedule fire consumer.
// Consumes fired NATS Schedule messages from WORKFLOW-SCHEDULE stream
// and re-publishes to WORKFLOW-RESUME for HandleResume processing.
func NewScheduleFireConsumer(bus *natsModel.Bus, service ports.RuntimeServicePort) *natsModel.Consumer {
	return schedule_fire.NewConsumer(bus, service)
}
