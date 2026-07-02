package collection

import (
	"context"

	"assets/src/modules/assettemplates/domain/entities"
	"assets/src/modules/assettemplates/domain/repositories"
	"assets/src/modules/assettemplates/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// NewFieldVocabulary creates and returns the field vocabulary repository.
// It binds a generic Model to the field_vocabulary collection.
func NewFieldVocabulary(m *manager.MongoManager) repositories.FieldVocabularyRepository {
	mdl := model.New[entities.FieldVocabulary](m.GetDatabase(), constants.FieldVocabularyCollectionName, model.Config{
		Indexes: constants.FieldVocabularyIndexes,
	})
	return &fieldVocabularyRepository{model: mdl}
}

// ListFieldVocabulary returns every enabled vocabulary field visible to the
// caller. Visibility is the union of the central platform standard
// (isSystem=true) and the caller's own org-scoped entries, mirroring the
// asset-template list filter: orgfilter.BuildOrgFilter(...) OR isSystem=true.
func (r *fieldVocabularyRepository) ListFieldVocabulary(ctx context.Context, reqContext *reqCtx.RequestContext) ([]entities.FieldVocabulary, error) {
	orConditions := []model.Map{}
	if orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: reqContext}); len(orgFilter) > 0 {
		orConditions = append(orConditions, orgFilter)
	}
	orConditions = append(orConditions, model.Map{"isSystem": true})

	filters := model.Map{"enabled": true}
	if len(orConditions) > 0 {
		filters["$or"] = orConditions
	}

	cursor, err := r.model.DIRECT().Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]entities.FieldVocabulary, 0)
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Compile-time check to ensure the adapter implements the repository port.
var _ repositories.FieldVocabularyRepository = (*fieldVocabularyRepository)(nil)
