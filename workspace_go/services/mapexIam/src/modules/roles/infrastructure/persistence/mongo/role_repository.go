package collection

import (
	"context"

	"mapexIam/src/modules/roles/domain/entities"
	"mapexIam/src/modules/roles/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the Role entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[Role],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.RoleRepository {
	mdl := model.New[entities.Role](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new Role entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Role entity to be persisted.
// Internally, it calls the underlying model's CreateOne method to store the document in MongoDB.
// It returns the created Role (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.Role) (*entities.Role, error) {
	role, err := r.model.CreateOne(ctx, u)
	return role, err
}

// FindById retrieves a Role entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the role ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - roleId: A string representing the unique identifier of the Role to be retrieved.
//
// Returns:
//   - A pointer to the Role entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, roleId *string) (*entities.Role, error) {
	retData, _ := r.model.FindByID(ctx, *roleId)
	return retData, nil
}

// FindByIdAndUpdate updates a Role entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the role ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - roleId: A pointer to the string representing the unique identifier of the Role to be updated.
//   - payload: A map containing the fields and their new values to update in the Role entity.
//
// Returns:
//   - A pointer to the updated Role entity, populated with the new values from the database.
//   - An error if the update operation fails or if the Role is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, roleId *string, payload map[string]any) (*entities.Role, error) {

	_id, _ := model.ToObjectID(*roleId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a Role entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the role ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - roleId: A pointer to the string representing the unique identifier of the Role to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the Role is not found.
//   - nil if the Role is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, roleId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*roleId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of Role entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"isSystem": true, "permissions": "read_users"}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"name": 1, "permissions": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching Role entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Role], error) {
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

// Compile-time check to ensure repository implements RoleRepository
var _ repositories.RoleRepository = (*repository)(nil)
