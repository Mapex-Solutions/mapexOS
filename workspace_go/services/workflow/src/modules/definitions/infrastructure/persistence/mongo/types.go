package collection

import (
	"workflow/src/modules/definitions/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type repository struct {
	model *model.Model[entities.WorkflowDefinition]
}
