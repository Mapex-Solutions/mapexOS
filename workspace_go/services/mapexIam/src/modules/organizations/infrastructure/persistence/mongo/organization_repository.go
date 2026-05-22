package collection

import (
	"context"

	"mapexIam/src/modules/organizations/domain/entities"
	"mapexIam/src/modules/organizations/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the XXX entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[XXX],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.OrganizationRepository {
	mdl := model.New[entities.Organization](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new Organization entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Organization entity to be persisted.
// Internally, it calls the underlying model's CreateOne method to store the document in MongoDB.
// It returns the created Organization (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.Organization) (*entities.Organization, error) {
	organization, err := r.model.CreateOne(ctx, u)
	return organization, err
}

// FindById retrieves an Organization entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the organization ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - organizationId: A string representing the unique identifier of the Organization to be retrieved.
//
// Returns:
//   - A pointer to the Organization entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, organizationId *string) (*entities.Organization, error) {
	retData, _ := r.model.FindByID(ctx, *organizationId)
	return retData, nil
}

// FindByIds retrieves multiple Organization entities by their IDs.
// Uses cursor-based pagination to fetch results in batches to avoid DB overload.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - organizationIds: A slice of organization ID strings to retrieve.
//
// Returns:
//   - A slice of pointers to Organization entities.
//   - An error if the query fails.
func (r *repository) FindByIds(ctx context.Context, organizationIds []string) ([]*entities.Organization, error) {
	if len(organizationIds) == 0 {
		return []*entities.Organization{}, nil
	}

	// Convert string IDs to ObjectIDs
	objectIDs := make([]any, 0, len(organizationIds))
	for _, id := range organizationIds {
		objID, err := model.ToObjectID(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIDs = append(objectIDs, objID)
	}

	if len(objectIDs) == 0 {
		return []*entities.Organization{}, nil
	}

	// Query using $in operator
	filters := model.Map{"_id": model.Map{"$in": objectIDs}}

	var allOrganizations []*entities.Organization
	var lastCursorID any = nil
	perPage := int64(1000)

	// Loop using cursor-based pagination
	for {
		pagination := &model.PaginationOpts{
			UseCursor:     true,
			CursorID:      lastCursorID,
			PerPage:       perPage,
			SortDirection: 1,
		}

		result, err := r.model.FindByCursor(ctx, filters, pagination, nil)
		if err != nil {
			return nil, err
		}

		// Append current batch
		for i := range result.Items {
			allOrganizations = append(allOrganizations, &result.Items[i])
		}

		// Check if there are more results
		if result.Pagination.HasNext == nil || !*result.Pagination.HasNext {
			break
		}

		// Get last ID for next iteration
		if len(result.Items) > 0 {
			lastCursorID = result.Items[len(result.Items)-1].ID
		} else {
			break
		}
	}

	return allOrganizations, nil
}

// FindByIdAndUpdate updates an Organization entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the organization ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - organizationId: A pointer to the string representing the unique identifier of the Organization to be updated.
//   - payload: A map containing the fields and their new values to update in the Organization entity.
//
// Returns:
//   - A pointer to the updated Organization entity, populated with the new values from the database.
//   - An error if the update operation fails or if the Organization is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, organizationId *string, payload map[string]any) (*entities.Organization, error) {

	_id, _ := model.ToObjectID(*organizationId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes an Organization entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the organization ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - organizationId: A pointer to the string representing the unique identifier of the Organization to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the Organization is not found.
//   - nil if the Organization is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, organizationId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*organizationId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of Organization entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"type": "customer", "enabled": true}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"name": 1, "type": 1, "pathKey": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching Organization entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Organization], error) {
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

// FindWithCursor retrieves Organization entities using cursor-based pagination.
// This method is optimized for infinite scroll and tree navigation scenarios.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"pathKey": {"$gte": "000001", "$lt": "000002"}}).
//   - cursorOpts: Cursor pagination options (cursor ID, direction, limit, sort).
//   - projection: A map specifying which fields to include in the results.
//
// Returns:
//   - A pointer to CursorResult containing the matching Organization entities and cursor metadata.
//   - An error if the query fails.
func (r *repository) FindWithCursor(
	ctx context.Context,
	filters model.Map,
	cursorOpts *model.CursorOpts,
	projection model.Map,
) (*model.CursorResult[entities.Organization], error) {
	result, err := r.model.FindWithCursor(ctx, filters, cursorOpts, projection)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Compile-time check to ensure repository implements OrganizationRepository
var _ repositories.OrganizationRepository = (*repository)(nil)
