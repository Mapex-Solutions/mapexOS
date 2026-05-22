package collection

import (
	"events/src/modules/retention/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// repository is the MongoDB-backed implementation of the RetentionRepository
// port. It wraps model.Model[RetentionPolicy] for CRUD and pagination.
type repository struct {
	model *model.Model[entities.RetentionPolicy]
}
