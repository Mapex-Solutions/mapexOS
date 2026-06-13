package services

import (
	"workflow/src/modules/archiver/application/di"
	archiverTypes "workflow/src/modules/archiver/application/types"
	"workflow/src/modules/archiver/domain/repositories"
	runtimePorts "workflow/src/modules/runtime/application/ports"
)

type ArchiverService struct {
	deps di.ArchiverServiceDependenciesInjection
}

// classifiedState carries the per-bucket payload that ProcessStateBatch produces
// during classification and consumes during the write phase.
type classifiedState struct {
	createdStubs       []repositories.LightweightExecution
	waitingUpdates     []repositories.WaitingUpdate
	resumedIDs         []string
	terminalExecutions []*runtimePorts.WorkflowExecution
	terminalKVKeys     []string
	refs               []archiverTypes.MsgRef
}
