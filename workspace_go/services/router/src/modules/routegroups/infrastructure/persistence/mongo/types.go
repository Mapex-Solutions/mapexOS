package collection

import (
	"router/src/modules/routegroups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type repository struct {
	model *model.Model[entities.RouteGroup]
}
