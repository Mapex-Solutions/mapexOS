package collection

import (
	"context"

	"mapexIam/src/modules/users/domain/entities"
	"mapexIam/src/modules/users/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a repository for the User entity.
// It initializes a Model[User] with the configured collection name and indexes.
//
// Indexes created:
//   - idx_email_unique: Unique index on email field
func New(m *manager.MongoManager) repositories.UserRepository {
	mdl := model.New[entities.User](m.GetDatabase(), CollectionName, model.Config{
		Indexes: Indexes,
	})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new User entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a User entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created User (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.User) (*entities.User, error) {
	user, err := r.model.CreateOne(ctx, u)
	return user, err
}

// FindById retrieves a User entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the user ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - userId: A string representing the unique identifier of the User to be retrieved.
//
// Returns:
//   - A pointer to the User entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, userId *string) (*entities.User, error) {
	retData, _ := r.model.FindByID(ctx, *userId)
	return retData, nil
}

// FindByIdAndUpdate updates a User entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the user ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - userId: A pointer to the string representing the unique identifier of the User to be updated.
//   - payload: A map containing the fields and their new values to update in the User entity.
//
// Returns:
//   - A pointer to the updated User entity, populated with the new values from the database.
//   - An error if the update operation fails or if the User is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, userId *string, payload map[string]any) (*entities.User, error) {

	_id, _ := model.ToObjectID(*userId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a User entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the user ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - userId: A pointer to the string representing the unique identifier of the User to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the User is not found.
//   - nil if the User is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, userId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*userId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindByEmail retrieves a User entity from the repository by its email address.
// It accepts a context for cancellation and timeouts, and a string representing the user's email.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - email: A string representing the email address of the User to be retrieved.
//
// Returns:
//   - A pointer to the User entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindByEmail(ctx context.Context, email *string) (*entities.User, error) {
	query := model.Map{"email": email}
	retData, _ := r.model.FindOne(ctx, &query)
	return retData, nil
}

// FindWithFilters retrieves a paginated list of User entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"email": "user@example.com", "enabled": true}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"email": 1, "firstName": 1, "lastName": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching User entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.User], error) {
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

// Compile-time check to ensure repository implements UserRepository
var _ repositories.UserRepository = (*repository)(nil)
