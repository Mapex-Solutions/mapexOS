package collection

import (
	"context"

	"mapexIam/src/modules/lists/domain/entities"
	"mapexIam/src/modules/lists/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the XXX entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[XXX],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.ListRepository {
	mdl := model.New[entities.List](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new List entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a List entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created List (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.List) (*entities.List, error) {
	list, err := r.model.CreateOne(ctx, u)
	return list, err
}

// FindById retrieves a List entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the list ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - listId: A string representing the unique identifier of the List to be retrieved.
//
// Returns:
//   - A pointer to the List entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, listId *string) (*entities.List, error) {
	retData, _ := r.model.FindByID(ctx, *listId)
	return retData, nil
}

// FindByIdAndUpdate updates a List entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the list ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - listId: A pointer to the string representing the unique identifier of the List to be updated.
//   - payload: A map containing the fields and their new values to update in the List entity.
//
// Returns:
//   - A pointer to the updated List entity, populated with the new values from the database.
//   - An error if the update operation fails or if the List is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, listId *string, payload map[string]any) (*entities.List, error) {

	_id, _ := model.ToObjectID(*listId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a List entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the list ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - listId: A pointer to the string representing the unique identifier of the List to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the List is not found.
//   - nil if the List is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, listId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*listId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindByEmail retrieves a List entity from the repository by its email address.
// It accepts a context for cancellation and timeouts, and a string representing the list's email.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - email: A string representing the email address of the List to be retrieved.
//
// Returns:
//   - A pointer to the List entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindByEmail(ctx context.Context, email *string) (*entities.List, error) {
	query := model.Map{"email": email}
	retData, _ := r.model.FindOne(ctx, &query)
	return retData, nil
}

// FindWithFilters retrieves a paginated list of List entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"type": "assetType", "isSystem": true}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"name": 1, "type": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching List entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.List], error) {
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

// Compile-time check to ensure repository implements ListRepository
var _ repositories.ListRepository = (*repository)(nil)
