package collection

import (
	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// fieldVocabularyRepository is the MongoDB-backed adapter implementing
// repositories.FieldVocabularyRepository. It wraps the generic
// *model.Model[entities.FieldVocabulary] bound to the field_vocabulary
// collection.
type fieldVocabularyRepository struct {
	model *model.Model[entities.FieldVocabulary]
}
