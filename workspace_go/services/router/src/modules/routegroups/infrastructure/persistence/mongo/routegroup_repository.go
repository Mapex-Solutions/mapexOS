package collection

import (
	"context"

	"router/src/modules/routegroups/domain/entities"
	"router/src/modules/routegroups/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the RouteGroup entity.
// It accepts a *MongoManager to obtain the database connection,
// then calls model.New to initialize a Model[RouteGroup],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.RouteGroupRepository {
	mdl := model.New[entities.RouteGroup](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

// Create inserts a new RouteGroup entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a RouteGroup entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created RouteGroup (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.RouteGroup) (*entities.RouteGroup, error) {
	routegroup, err := r.model.CreateOne(ctx, u)
	return routegroup, err
}

// FindById retrieves a RouteGroup entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the routegroup ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - listId: A string representing the unique identifier of the RouteGroup to be retrieved.
//
// Returns:
//   - A pointer to the RouteGroup entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, listId *string) (*entities.RouteGroup, error) {
	retData, err := r.model.FindByID(ctx, *listId)
	if err != nil {
		return nil, err
	}
	return retData, nil
}

// FindByIdAndUpdate updates a RouteGroup entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the routegroup ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - listId: A pointer to the string representing the unique identifier of the RouteGroup to be updated.
//   - payload: A map containing the fields and their new values to update in the RouteGroup entity.
//
// Returns:
//   - A pointer to the updated RouteGroup entity, populated with the new values from the database.
//   - An error if the update operation fails or if the RouteGroup is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, listId *string, payload map[string]any) (*entities.RouteGroup, error) {
	_id, err := model.ToObjectID(*listId)
	if err != nil {
		return nil, err
	}

	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, err := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	if err != nil {
		return nil, err
	}
	return retData, nil
}

// DeleteById removes a RouteGroup entity from the repository by its ID.
// It accepts a context for cancellation and timeouts and a pointer to the routegroup ID.
// Multi-tenant isolation is handled by the coverage middleware before reaching this method.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - listId: A pointer to the string representing the unique identifier of the RouteGroup to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the RouteGroup is not found.
//   - nil if the RouteGroup is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, listId *string) error {
	// Convert the string ID to an ObjectID for MongoDB
	_id, err := model.ToObjectID(*listId)
	if err != nil {
		return err
	}

	// Query using only _id (coverage middleware already validated access)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	return r.model.DeleteOne(ctx, &query)
}

// FindWithFilters retrieves a paginated list of route groups based on provided filters.
// This method supports complex queries including organization filtering, name search, and pagination.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - filters: MongoDB query filters (e.g., orgId, pathKey, name, enabled)
//   - pagination: Pagination options (page, perPage)
//   - projection: Fields to include/exclude in the result
//
// Returns:
//   - PaginatedResult: Contains the list of route groups and pagination metadata
//   - error: If the query operation fails
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.RouteGroup], error) {
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

// CountDocuments returns the total number of route group documents matching the given filters.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - filters: MongoDB query filters to narrow the count
//
// Returns:
//   - int64: The count of matching documents
//   - error: If the count operation fails
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}
