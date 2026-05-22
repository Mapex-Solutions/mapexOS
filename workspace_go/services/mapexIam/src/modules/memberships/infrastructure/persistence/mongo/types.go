package collection

import (
	"mapexIam/src/modules/memberships/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type repository struct {
	model *model.Model[entities.Membership]
}
