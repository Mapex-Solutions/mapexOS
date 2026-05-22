package collection

import (
	"context"

	"mapexIam/src/modules/groups/domain/entities"
	"mapexIam/src/modules/groups/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the XXX entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[XXX],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.GroupRepository {
	mdl := model.New[entities.Group](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new Group entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Group entity to be persisted.
// Internally, it calls the underlying model's CreateOne method to store the document in MongoDB.
// It returns the created Group (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.Group) (*entities.Group, error) {
	group, err := r.model.CreateOne(ctx, u)
	return group, err
}

// FindById retrieves a Group entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the group ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - groupId: A string representing the unique identifier of the Group to be retrieved.
//
// Returns:
//   - A pointer to the Group entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, groupId *string) (*entities.Group, error) {
	retData, _ := r.model.FindByID(ctx, *groupId)
	return retData, nil
}

// FindByIdAndUpdate updates a Group entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the group ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - groupId: A pointer to the string representing the unique identifier of the Group to be updated.
//   - payload: A map containing the fields and their new values to update in the Group entity.
//
// Returns:
//   - A pointer to the updated Group entity, populated with the new values from the database.
//   - An error if the update operation fails or if the Group is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, groupId *string, payload map[string]any) (*entities.Group, error) {

	_id, _ := model.ToObjectID(*groupId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a Group entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the group ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - groupId: A pointer to the string representing the unique identifier of the Group to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the Group is not found.
//   - nil if the Group is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, groupId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*groupId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of Group entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"enabled": true}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"name": 1, "enabled": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching Group entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Group], error) {
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

// CountDocuments counts documents matching the provided filters.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - filters: A map of filters to apply to the count query
//
// Returns:
//   - int64: The number of matching documents
//   - error: If the count operation fails
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Compile-time check to ensure repository implements GroupRepository
var _ repositories.GroupRepository = (*repository)(nil)
