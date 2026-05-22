package collection

import (
	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// repository is the MongoDB-backed adapter implementing repositories.AssetTemplateRepository.
// It wraps the generic *model.Model[entities.Assettemplate] used by all repository
// methods declared in assettemplate_repository.go.
type repository struct {
	model *model.Model[entities.Assettemplate]
}
