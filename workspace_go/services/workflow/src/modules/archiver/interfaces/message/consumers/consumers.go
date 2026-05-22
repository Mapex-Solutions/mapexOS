package consumers

import (
	"workflow/src/modules/archiver/application/ports"
	"workflow/src/modules/archiver/interfaces/message/consumers/workflow_state"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// NewWorkflowStateConsumer starts the WORKFLOW-STATE batch consumer for archiving.
func NewWorkflowStateConsumer(bus *natsModel.Bus, service ports.ArchiverServicePort) *natsModel.Consumer {
	return workflow_state.NewConsumer(bus, service)
}
