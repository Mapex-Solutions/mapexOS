package collection

import (
	"context"

	"workflow/src/modules/plugins/domain/entities"
	"workflow/src/modules/plugins/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates a repository for PluginManifest entities.
func New(m *manager.MongoManager) repositories.PluginManifestRepository {
	mdl := model.New[entities.PluginManifest](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

// Create inserts a new PluginManifest entity into MongoDB.
func (r *repository) Create(ctx context.Context, entity *entities.PluginManifest) (*entities.PluginManifest, error) {
	return r.model.CreateOne(ctx, entity)
}

// FindById retrieves a PluginManifest by its MongoDB ObjectId.
func (r *repository) FindById(ctx context.Context, id *string) (*entities.PluginManifest, error) {
	return r.model.FindByID(ctx, *id)
}

// FindByPluginId retrieves a PluginManifest by its unique pluginId field.
func (r *repository) FindByPluginId(ctx context.Context, pluginId string) (*entities.PluginManifest, error) {
	query := model.Map{"pluginId": pluginId}
	return r.model.FindOne(ctx, &query, nil)
}

// FindByIdAndUpdate updates a PluginManifest by ID using $set operator.
func (r *repository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.PluginManifest, error) {
	_id, err := model.ToObjectID(*id)
	if err != nil {
		return nil, err
	}

	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	return r.model.FindOneAndUpdate(ctx, &query, &update, &options)
}

// DeleteById removes a PluginManifest by its MongoDB ObjectId.
func (r *repository) DeleteById(ctx context.Context, id *string) error {
	_id, err := model.ToObjectID(*id)
	if err != nil {
		return err
	}

	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

// FindWithFilters retrieves a paginated list of plugin manifests.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.PluginManifest], error) {
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

// CountDocuments returns the count of plugin manifests matching the filters.
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	return r.model.DIRECT().CountDocuments(ctx, filters)
}
