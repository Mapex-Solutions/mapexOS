package collection

import (
	"mapexIam/src/modules/groups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type repository struct {
	model *model.Model[entities.Group]
}

type groupMemberRepository struct {
	model *model.Model[entities.GroupMember]
}
