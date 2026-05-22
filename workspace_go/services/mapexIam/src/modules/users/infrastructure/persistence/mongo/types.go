package collection

import (
	"mapexIam/src/modules/users/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type repository struct {
	model *model.Model[entities.User]
}
