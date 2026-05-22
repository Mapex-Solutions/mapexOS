package collection

import (
	"context"

	"workflow/src/modules/instances/domain/entities"
	"workflow/src/modules/instances/domain/repositories"
	"workflow/src/modules/instances/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// Compile-time check
var _ repositories.InstanceRepository = (*repository)(nil)

// New creates and returns a repository for the WorkflowInstance entity.
func New(m *manager.MongoManager) repositories.InstanceRepository {
	mdl := model.New[entities.WorkflowInstance](m.GetDatabase(), constants.CollectionName, model.Config{
		Indexes: constants.Indexes,
	})
	return &repository{model: mdl}
}

// Create inserts a new WorkflowInstance entity into the repository.
func (r *repository) Create(ctx context.Context, instance *entities.WorkflowInstance) (*entities.WorkflowInstance, error) {
	created, err := r.model.CreateOne(ctx, instance)
	return created, err
}

// FindById retrieves a WorkflowInstance by its MongoDB ID.
func (r *repository) FindById(ctx context.Context, id *string) (*entities.WorkflowInstance, error) {
	retData, err := r.model.FindByID(ctx, *id)
	return retData, err
}

// FindByIdAndUpdate updates a WorkflowInstance and returns the updated document.
func (r *repository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.WorkflowInstance, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a WorkflowInstance by its MongoDB ID.
func (r *repository) DeleteById(ctx context.Context, id *string) error {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of WorkflowInstance entities.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.WorkflowInstance], error) {
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

// CountDocuments counts instances matching the provided filters.
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}
