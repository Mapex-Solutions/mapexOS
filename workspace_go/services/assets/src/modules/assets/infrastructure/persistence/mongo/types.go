package collection

import (
	"assets/src/modules/assets/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// repository is the MongoDB-backed adapter implementing repositories.AssetRepository.
// It wraps the generic *model.Model[entities.Asset] used by all repository methods
// declared in asset_repository.go.
type repository struct {
	model *model.Model[entities.Asset]
}
