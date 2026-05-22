package collection

import (
	"workflow/src/modules/instances/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type repository struct {
	model *model.Model[entities.WorkflowInstance]
}
