package collection

import (
	"context"
	"time"

	"events/src/modules/retention/domain/entities"
	"events/src/modules/retention/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the RetentionPolicy entity.
// It accepts a *MongoManager to obtain the database connection,
// then calls model.New to initialize a Model[RetentionPolicy],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.RetentionRepository {
	mdl := model.New[entities.RetentionPolicy](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

// Create inserts a new RetentionPolicy entity into the repository.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - policy: A pointer to a RetentionPolicy entity to be persisted.
//
// Returns:
//   - A pointer to the created RetentionPolicy with database-assigned fields.
//   - An error if something goes wrong.
func (r *repository) Create(ctx context.Context, policy *entities.RetentionPolicy) (*entities.RetentionPolicy, error) {
	created, err := r.model.CreateOne(ctx, policy)
	return created, err
}

// FindById retrieves a RetentionPolicy entity from the repository by its ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - policyId: A string representing the unique identifier of the RetentionPolicy.
//
// Returns:
//   - A pointer to the RetentionPolicy entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, policyId *string) (*entities.RetentionPolicy, error) {
	return r.model.FindByID(ctx, *policyId)
}

// FindByOrgIdAndType retrieves a RetentionPolicy entity by organization ID and type.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - orgId: A string representing the organization ID.
//   - retentionType: The retention type (e.g., "events", "eventsRaw").
//
// Returns:
//   - A pointer to the RetentionPolicy entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindByOrgIdAndType(ctx context.Context, orgId *string, retentionType string) (*entities.RetentionPolicy, error) {
	objId, err := model.ToObjectID(*orgId)
	if err != nil {
		return nil, err
	}
	query := model.Map{"orgId": objId, "type": retentionType}
	return r.model.FindOne(ctx, &query)
}

// Upsert creates or updates a retention policy by organization ID and type.
// Uses FindOneAndUpdate with upsert=true on filter {orgId, type}.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - orgId: The organization ObjectId.
//   - retentionType: The retention type (e.g., "events", "eventsRaw").
//   - policy: The retention policy data to upsert.
//
// Returns:
//   - A pointer to the upserted RetentionPolicy entity.
//   - An error if the upsert operation fails.
func (r *repository) Upsert(ctx context.Context, orgId *model.ObjectId, retentionType string, policy *entities.RetentionPolicy) (*entities.RetentionPolicy, error) {
	query := model.Map{"orgId": orgId, "type": retentionType}

	now := time.Now()
	update := model.Map{
		"$set": model.Map{
			"name":          policy.Name,
			"retentionDays": policy.RetentionDays,
			"pathKey":       policy.PathKey,
			"enabled":       policy.Enabled,
			"updated":       now,
		},
		"$setOnInsert": model.Map{
			"orgId":   orgId,
			"type":    retentionType,
			"created": now,
		},
	}

	upsertTrue := true
	returnDoc := model.ReturnDoc(1) // Return document after update
	options := model.CommonOpts{
		Upsert:         &upsertTrue,
		ReturnDocument: &returnDoc,
	}

	retData, err := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, err
}

// DeleteById removes a RetentionPolicy entity from the repository by its ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - policyId: A pointer to the string representing the unique identifier.
//
// Returns:
//   - An error if the deletion operation fails.
func (r *repository) DeleteById(ctx context.Context, policyId *string) error {
	_id, err := model.ToObjectID(*policyId)
	if err != nil {
		return err
	}
	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

// FindWithFilters retrieves a paginated list of retention policies based on filters.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: MongoDB query filters.
//   - pagination: Pagination options (page, perPage).
//   - projection: Fields to include/exclude in the result.
//
// Returns:
//   - PaginatedResult: Contains the list of retention policies and pagination metadata.
//   - error: If the query operation fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.RetentionPolicy], error) {
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
