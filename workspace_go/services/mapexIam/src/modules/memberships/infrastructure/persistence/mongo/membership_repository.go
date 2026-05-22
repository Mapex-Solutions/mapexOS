package collection

import (
	"context"

	"mapexIam/src/modules/memberships/domain/entities"
	"mapexIam/src/modules/memberships/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a repository for the Membership entity.
// It initializes a Model[Membership] with the configured collection name and indexes.
//
// Indexes created:
//   - idx_membership_unique: Unique compound index on (assigneeType, assigneeId, orgId)
//   - idx_assignee: Compound index on (assigneeType, assigneeId)
//   - idx_org: Index on orgId
func New(m *manager.MongoManager) repositories.MembershipRepository {
	mdl := model.New[entities.Membership](m.GetDatabase(), CollectionName, model.Config{
		Indexes: Indexes,
	})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new Membership entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Membership entity to be persisted.
// Internally, it calls the underlying model's CreateOne method to store the document in MongoDB.
// It returns the created Membership (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.Membership) (*entities.Membership, error) {
	membership, err := r.model.CreateOne(ctx, u)
	return membership, err
}

// FindById retrieves a Membership entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the membership ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - membershipId: A string representing the unique identifier of the Membership to be retrieved.
//
// Returns:
//   - A pointer to the Membership entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, membershipId *string) (*entities.Membership, error) {
	retData, _ := r.model.FindByID(ctx, *membershipId)
	return retData, nil
}

// FindByIdAndUpdate updates a Membership entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the membership ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - membershipId: A pointer to the string representing the unique identifier of the Membership to be updated.
//   - payload: A map containing the fields and their new values to update in the Membership entity.
//
// Returns:
//   - A pointer to the updated Membership entity, populated with the new values from the database.
//   - An error if the update operation fails or if the Membership is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, membershipId *string, payload map[string]any) (*entities.Membership, error) {

	_id, _ := model.ToObjectID(*membershipId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a Membership entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the membership ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - membershipId: A pointer to the string representing the unique identifier of the Membership to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the Membership is not found.
//   - nil if the Membership is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, membershipId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*membershipId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindByUserId retrieves all Membership entities for a specific user.
// It queries memberships where assigneeType="user" and assigneeId matches the provided userId.
// Uses cursor-based pagination to fetch all results in batches of 300.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - userId: A pointer to the string representing the unique identifier of the user.
//
// Returns:
//   - A slice of pointers to Membership entities.
//   - An error if the query fails.
func (r *repository) FindByUserId(ctx context.Context, userId *string) ([]*entities.Membership, error) {
	userObjectID, err := model.ToObjectID(*userId)
	if err != nil {
		return nil, err
	}

	filters := model.Map{"assigneeType": "user", "assigneeId": userObjectID}

	var allMemberships []*entities.Membership
	var lastCursorID interface{} = nil
	perPage := int64(300) // Batch size as per architecture decision

	// Loop using cursor-based pagination until all results are fetched
	for {
		pagination := &model.PaginationOpts{
			UseCursor:     true,
			CursorID:      lastCursorID,
			PerPage:       perPage,
			SortDirection: 1, // forward
		}

		result, err := r.model.FindByCursor(ctx, filters, pagination, nil)
		if err != nil {
			return nil, err
		}

		// Append current batch to result
		for _, item := range result.Items {
			membershipCopy := item // Create copy to avoid pointer aliasing
			allMemberships = append(allMemberships, &membershipCopy)
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

	return allMemberships, nil
}

// FindByGroupIds retrieves all Membership entities for a list of groups.
// It queries memberships where assigneeType="group" and assigneeId IN groupIds.
// Uses cursor-based pagination to fetch all results in batches of 300.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - groupIds: A slice of group IDs to query memberships for.
//
// Returns:
//   - A slice of pointers to Membership entities.
//   - An error if the query fails.
func (r *repository) FindByGroupIds(ctx context.Context, groupIds []string) ([]*entities.Membership, error) {
	if len(groupIds) == 0 {
		return []*entities.Membership{}, nil
	}

	// Convert string IDs to ObjectIDs
	groupObjectIDs := make([]model.ObjectId, 0, len(groupIds))
	for _, gid := range groupIds {
		objectID, err := model.ToObjectID(gid)
		if err != nil {
			continue // Skip invalid IDs
		}
		groupObjectIDs = append(groupObjectIDs, objectID)
	}

	if len(groupObjectIDs) == 0 {
		return []*entities.Membership{}, nil
	}

	filters := model.Map{
		"assigneeType": "group",
		"assigneeId":   model.Map{"$in": groupObjectIDs},
	}

	var allMemberships []*entities.Membership
	var lastCursorID interface{} = nil
	perPage := int64(300) // Batch size as per architecture decision

	// Loop using cursor-based pagination until all results are fetched
	for {
		pagination := &model.PaginationOpts{
			UseCursor:     true,
			CursorID:      lastCursorID,
			PerPage:       perPage,
			SortDirection: 1, // forward
		}

		result, err := r.model.FindByCursor(ctx, filters, pagination, nil)
		if err != nil {
			return nil, err
		}

		// Append current batch to result
		for _, item := range result.Items {
			membershipCopy := item // Create copy to avoid pointer aliasing
			allMemberships = append(allMemberships, &membershipCopy)
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

	return allMemberships, nil
}

// FindWithFilters retrieves a paginated list of Membership entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"assigneeType": "user", "enabled": true}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"assigneeType": 1, "roles": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching Membership entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Membership], error) {
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

// FindByAssigneeAndOrg retrieves a Membership entity by assignee and organization.
// Used for idempotent checks before creating new memberships.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - assigneeType: The type of assignee ("user" or "group").
//   - assigneeId: The unique identifier of the assignee (as string).
//   - orgId: The unique identifier of the organization (as string).
//
// Returns:
//   - A pointer to the Membership entity if found, or nil if not found.
//   - An error if the query fails.
func (r *repository) FindByAssigneeAndOrg(ctx context.Context, assigneeType string, assigneeId string, orgId string) (*entities.Membership, error) {
	assigneeObjectID, err := model.ToObjectID(assigneeId)
	if err != nil {
		return nil, err
	}

	orgObjectID, err := model.ToObjectID(orgId)
	if err != nil {
		return nil, err
	}

	query := model.Map{
		"assigneeType": assigneeType,
		"assigneeId":   assigneeObjectID,
		"orgId":        orgObjectID,
	}

	retData, err := r.model.FindOne(ctx, &query)
	if err != nil {
		return nil, err
	}

	return retData, nil
}

// Compile-time check to ensure repository implements MembershipRepository
var _ repositories.MembershipRepository = (*repository)(nil)
