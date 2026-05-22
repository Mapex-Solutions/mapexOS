package collection

import (
	"context"

	"http_gateway/src/modules/datasources/domain/entities"
	"http_gateway/src/modules/datasources/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

var _ repositories.DataSourceRepository = (*repository)(nil)

// New creates and returns a repository for the DataSource entity.
// It accepts a *MongoManager to obtain the database connection,
// then calls model.New to initialize a Model[DataSource],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.DataSourceRepository {
	mdl := model.New[entities.DataSource](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

// Create inserts a new DataSource entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a DataSource entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created DataSource (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.DataSource) (*entities.DataSource, error) {
	event, err := r.model.CreateOne(ctx, u)
	return event, err
}

// FindById retrieves a DataSource entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the event ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - dataSourceId: A string representing the unique identifier of the DataSource to be retrieved.
//
// Returns:
//   - A pointer to the DataSource entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, dataSourceId *string) (*entities.DataSource, error) {
	return r.model.FindByID(ctx, *dataSourceId)
}

// FindByIdAndUpdate updates a DataSource entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the event ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - dataSourceId: A pointer to the string representing the unique identifier of the DataSource to be updated.
//   - payload: A map containing the fields and their new values to update in the DataSource entity.
//
// Returns:
//   - A pointer to the updated DataSource entity, populated with the new values from the database.
//   - An error if the update operation fails or if the DataSource is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, dataSourceId *string, payload map[string]any) (*entities.DataSource, error) {
	_id, err := model.ToObjectID(*dataSourceId)
	if err != nil {
		return nil, err
	}

	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	return r.model.FindOneAndUpdate(ctx, &query, &update, &options)
}

// DeleteById removes a DataSource entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the event ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - dataSourceId: A pointer to the string representing the unique identifier of the DataSource to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the DataSource is not found.
//   - nil if the DataSource is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, dataSourceId *string) error {
	_id, err := model.ToObjectID(*dataSourceId)
	if err != nil {
		return err
	}

	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

// FindWithFilters retrieves a paginated list of data sources based on provided filters.
// This method supports complex queries including organization filtering, name search, and pagination.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - filters: MongoDB query filters (e.g., orgId, pathKey, name, enabled, mode, protocol)
//   - pagination: Pagination options (page, perPage)
//   - projection: Fields to include/exclude in the result
//
// Returns:
//   - PaginatedResult: Contains the list of data sources and pagination metadata
//   - error: If the query operation fails
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.DataSource], error) {
	// Build options with projection
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}

	// Execute paginated query
	result, err := r.model.FindByOffset(ctx, filters, pagination, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}
