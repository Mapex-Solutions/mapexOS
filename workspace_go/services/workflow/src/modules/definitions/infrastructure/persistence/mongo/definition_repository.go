package collection

import (
	"context"

	"workflow/src/modules/definitions/domain/entities"
	"workflow/src/modules/definitions/domain/repositories"
	"workflow/src/modules/definitions/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

var _ repositories.DefinitionRepository = (*repository)(nil)

// New creates and returns a repository for the WorkflowDefinition entity.
func New(m *manager.MongoManager) repositories.DefinitionRepository {
	mdl := model.New[entities.WorkflowDefinition](m.GetDatabase(), constants.CollectionName, model.Config{
		Indexes: constants.Indexes,
	})
	return &repository{model: mdl}
}

// Create inserts a new WorkflowDefinition entity into the repository.
func (r *repository) Create(ctx context.Context, def *entities.WorkflowDefinition) (*entities.WorkflowDefinition, error) {
	created, err := r.model.CreateOne(ctx, def)
	return created, err
}

// FindById retrieves a WorkflowDefinition by its MongoDB ID.
func (r *repository) FindById(ctx context.Context, id *string) (*entities.WorkflowDefinition, error) {
	retData, err := r.model.FindByID(ctx, *id)
	return retData, err
}

// FindByIdAndUpdate updates a WorkflowDefinition and returns the updated document.
func (r *repository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.WorkflowDefinition, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a WorkflowDefinition by its MongoDB ID.
func (r *repository) DeleteById(ctx context.Context, id *string) error {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of WorkflowDefinition entities.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.WorkflowDefinition], error) {
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}

	result, err := r.model.FindByOffset(ctx, filters, pagination, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CountDocuments counts documents matching the provided filters.
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}
