package collection

import (
	"workflow/src/modules/archiver/domain/repositories"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"go.uber.org/dig"
)

type repository struct {
	model *model.Model[runtimePorts.WorkflowExecution]
}

// RepositoryOut provides the ArchiveRepository from the concrete type.
type RepositoryOut struct {
	dig.Out
	ArchiveRepo repositories.ArchiveRepository
}
