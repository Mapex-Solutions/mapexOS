package collection

import (
	"context"

	"triggers/src/modules/triggers/domain/entities"
	"triggers/src/modules/triggers/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the Trigger entity.
// It accepts a *MongoManager to obtain the database connection,
// then calls model.New to initialize a Model[Trigger],
// targeting the configured collection name with default settings.
//
// Parameters:
//   - m: The MongoManager instance providing database access
//
// Returns:
//   - repositories.TriggerRepository: The repository interface implementation
func New(m *manager.MongoManager) repositories.TriggerRepository {
	mdl := model.New[entities.Trigger](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

/* REPOSITORY METHODS */

// Create inserts a new Trigger entity into the repository.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - trigger: Pointer to a Trigger entity to be persisted
//
// Returns:
//   - *entities.Trigger: The created Trigger (with database-assigned fields)
//   - error: If something goes wrong during creation
func (r *repository) Create(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error) {
	createdTrigger, err := r.model.CreateOne(ctx, trigger)
	return createdTrigger, err
}

// FindById retrieves a Trigger entity from the repository by its ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - triggerId: String representing the unique identifier of the Trigger
//
// Returns:
//   - *entities.Trigger: The Trigger entity if found, nil if not found
//   - error: If the retrieval operation fails
func (r *repository) FindById(ctx context.Context, triggerId *string) (*entities.Trigger, error) {
	trigger, err := r.model.FindByID(ctx, *triggerId)
	return trigger, err
}

// FindByIdAndUpdate updates a Trigger entity in the repository by its ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - triggerId: String representing the unique identifier of the Trigger
//   - payload: Map containing the fields and their new values to update
//
// Returns:
//   - *entities.Trigger: The updated Trigger entity
//   - error: If the update operation fails or Trigger not found
func (r *repository) FindByIdAndUpdate(
	ctx context.Context,
	triggerId *string,
	payload map[string]any,
) (*entities.Trigger, error) {

	_id, err := model.ToObjectID(*triggerId)
	if err != nil {
		return nil, err
	}

	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)                           // Return the updated document
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	updatedTrigger, err := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return updatedTrigger, err
}

// DeleteById removes a Trigger entity from the repository by its ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - triggerId: String representing the unique identifier of the Trigger
//
// Returns:
//   - error: If the deletion operation fails or Trigger not found
func (r *repository) DeleteById(ctx context.Context, triggerId *string) error {
	_id, err := model.ToObjectID(*triggerId)
	if err != nil {
		return err
	}

	query := model.Map{"_id": _id}
	err = r.model.DeleteOne(ctx, &query)
	return err
}

// CountDocuments returns the total count of Trigger entities matching the provided filters.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - filters: MongoDB query filters (e.g., {"orgId": objectId})
//
// Returns:
//   - int64: The total count of matching documents
//   - error: If the count operation fails
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindWithFilters retrieves a paginated list of Trigger entities matching the provided filters.
//
// This method supports:
//   - Complex filtering (e.g., field equality, regex, range queries)
//   - Pagination (page, perPage)
//   - Sorting (sort field and direction)
//   - Field projection (select specific fields)
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - filters: MongoDB query filters (e.g., {"status": "Active", "category": "technical"})
//   - pagination: Pagination options (page, perPage, sort)
//   - projection: Field projection (nil for all fields)
//
// Returns:
//   - *model.PaginatedResult[entities.Trigger]: Paginated results with metadata
//   - error: If the query operation fails
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Trigger], error) {

	// Execute paginated find using FindByOffset method
	opts := &model.CommonOpts{Projection: &projection}
	result, err := r.model.FindByOffset(ctx, filters, pagination, opts)
	return result, err
}
