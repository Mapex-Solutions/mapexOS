package collection

import (
	"context"

	"mapexIam/src/modules/groups/domain/entities"
	"mapexIam/src/modules/groups/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// NewGroupMemberRepository creates and returns a repository for the GroupMember entity.
// It configures the necessary indexes for efficient queries:
//   - Unique compound index on (groupId, userId) to prevent duplicate memberships
//   - Compound index on (userId, orgId) for "user's groups in org" queries
//   - Index on groupId for "group members" queries
func NewGroupMemberRepository(m *manager.MongoManager) repositories.GroupMemberRepository {
	indexes := []model.IndexDefinition{
		{
			Name:   "idx_group_user_unique",
			Keys:   map[string]int{"groupId": 1, "userId": 1},
			Unique: true,
		},
		{
			Name: "idx_user_org",
			Keys: map[string]int{"userId": 1, "orgId": 1},
		},
		{
			Name: "idx_group",
			Keys: map[string]int{"groupId": 1},
		},
		{
			Name: "idx_pathkey",
			Keys: map[string]int{"pathKey": 1},
		},
	}

	mdl := model.New[entities.GroupMember](m.GetDatabase(), groupMemberCollectionName, model.Config{
		Indexes: indexes,
	})
	return &groupMemberRepository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new GroupMember entity into the repository.
// Returns an error if the membership already exists (duplicate groupId+userId).
func (r *groupMemberRepository) Create(ctx context.Context, member *entities.GroupMember) (*entities.GroupMember, error) {
	created, err := r.model.CreateOne(ctx, member)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// DeleteByGroupAndUser removes a specific group membership.
func (r *groupMemberRepository) DeleteByGroupAndUser(ctx context.Context, groupId, userId string) error {
	groupObjId, err := model.ToObjectID(groupId)
	if err != nil {
		return err
	}
	userObjId, err := model.ToObjectID(userId)
	if err != nil {
		return err
	}

	query := &model.Map{
		"groupId": groupObjId,
		"userId":  userObjId,
	}

	return r.model.DeleteOne(ctx, query)
}

// FindByGroupAndUser retrieves a specific group membership.
func (r *groupMemberRepository) FindByGroupAndUser(ctx context.Context, groupId, userId string) (*entities.GroupMember, error) {
	groupObjId, err := model.ToObjectID(groupId)
	if err != nil {
		return nil, err
	}
	userObjId, err := model.ToObjectID(userId)
	if err != nil {
		return nil, err
	}

	query := &model.Map{
		"groupId": groupObjId,
		"userId":  userObjId,
	}

	return r.model.FindOne(ctx, query, nil)
}

// FindByGroupId retrieves all members of a specific group.
func (r *groupMemberRepository) FindByGroupId(ctx context.Context, groupId string) ([]*entities.GroupMember, error) {
	groupObjId, err := model.ToObjectID(groupId)
	if err != nil {
		return nil, err
	}

	query := model.Map{"groupId": groupObjId}

	// Use FindByOffset with standard limit
	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    1,
		PerPage: 300,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Convert []entities.GroupMember to []*entities.GroupMember
	members := make([]*entities.GroupMember, len(result.Items))
	for i := range result.Items {
		members[i] = &result.Items[i]
	}

	return members, nil
}

// FindByGroupIds retrieves all members of multiple groups (bulk query).
// Used for efficient org-based user listing (users via group memberships).
func (r *groupMemberRepository) FindByGroupIds(ctx context.Context, groupIds []string) ([]*entities.GroupMember, error) {
	if len(groupIds) == 0 {
		return []*entities.GroupMember{}, nil
	}

	// Convert string IDs to ObjectIDs
	groupObjIds := make([]model.ObjectId, 0, len(groupIds))
	for _, gid := range groupIds {
		objId, err := model.ToObjectID(gid)
		if err != nil {
			continue // Skip invalid IDs
		}
		groupObjIds = append(groupObjIds, objId)
	}

	if len(groupObjIds) == 0 {
		return []*entities.GroupMember{}, nil
	}

	query := model.Map{
		"groupId": model.Map{"$in": groupObjIds},
	}

	// Use cursor pagination to fetch all results
	var allMembers []*entities.GroupMember
	page := int64(1)
	batchSize := int64(300)

	for {
		result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
			Page:    page,
			PerPage: batchSize,
		}, nil)
		if err != nil {
			return nil, err
		}

		for i := range result.Items {
			allMembers = append(allMembers, &result.Items[i])
		}

		// Check if we've fetched all pages
		if int64(len(result.Items)) < batchSize {
			break
		}
		page++
	}

	return allMembers, nil
}

// FindByUserId retrieves all group memberships for a specific user.
func (r *groupMemberRepository) FindByUserId(ctx context.Context, userId string) ([]*entities.GroupMember, error) {
	userObjId, err := model.ToObjectID(userId)
	if err != nil {
		return nil, err
	}

	query := model.Map{"userId": userObjId}

	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    1,
		PerPage: 300,
	}, nil)
	if err != nil {
		return nil, err
	}

	members := make([]*entities.GroupMember, len(result.Items))
	for i := range result.Items {
		members[i] = &result.Items[i]
	}

	return members, nil
}

// FindByUserIdAndOrgId retrieves all group memberships for a user in a specific organization.
func (r *groupMemberRepository) FindByUserIdAndOrgId(ctx context.Context, userId, orgId string) ([]*entities.GroupMember, error) {
	userObjId, err := model.ToObjectID(userId)
	if err != nil {
		return nil, err
	}
	orgObjId, err := model.ToObjectID(orgId)
	if err != nil {
		return nil, err
	}

	query := model.Map{
		"userId": userObjId,
		"orgId":  orgObjId,
	}

	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    1,
		PerPage: 300,
	}, nil)
	if err != nil {
		return nil, err
	}

	members := make([]*entities.GroupMember, len(result.Items))
	for i := range result.Items {
		members[i] = &result.Items[i]
	}

	return members, nil
}

// CountByGroupId returns the number of members in a specific group.
func (r *groupMemberRepository) CountByGroupId(ctx context.Context, groupId string) (int64, error) {
	groupObjId, err := model.ToObjectID(groupId)
	if err != nil {
		return 0, err
	}

	query := model.Map{"groupId": groupObjId}

	// Use FindByOffset and get the total count from pagination
	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    1,
		PerPage: 1, // Only need count, not data
	}, nil)
	if err != nil {
		return 0, err
	}

	return result.Pagination.TotalItems, nil
}

// DeleteByGroupId removes all memberships for a specific group (cascade delete).
func (r *groupMemberRepository) DeleteByGroupId(ctx context.Context, groupId string) error {
	groupObjId, err := model.ToObjectID(groupId)
	if err != nil {
		return err
	}

	query := model.Map{"groupId": groupObjId}

	_, err = r.model.DeleteMany(ctx, query)
	// Ignore ErrNotFound - it's ok if there are no members to delete
	if err != nil && err.Error() != "document not found" {
		return err
	}
	return nil
}

// DeleteByUserId removes all group memberships for a specific user.
func (r *groupMemberRepository) DeleteByUserId(ctx context.Context, userId string) error {
	userObjId, err := model.ToObjectID(userId)
	if err != nil {
		return err
	}

	query := model.Map{"userId": userObjId}

	_, err = r.model.DeleteMany(ctx, query)
	// Ignore ErrNotFound - it's ok if there are no memberships to delete
	if err != nil && err.Error() != "document not found" {
		return err
	}
	return nil
}

// FindByPathKeyRange retrieves all memberships in a hierarchical range using base36 $gte/$lt.
// This is O(log n) with B-tree index vs O(n) regex full scan.
func (r *groupMemberRepository) FindByPathKeyRange(ctx context.Context, pathKeyStart, pathKeyEnd string) ([]*entities.GroupMember, error) {
	query := model.Map{
		"pathKey": model.Map{
			"$gte": pathKeyStart,
			"$lt":  pathKeyEnd,
		},
	}

	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    1,
		PerPage: 300,
	}, nil)
	if err != nil {
		return nil, err
	}

	members := make([]*entities.GroupMember, len(result.Items))
	for i := range result.Items {
		members[i] = &result.Items[i]
	}

	return members, nil
}

// DeleteByPathKeyRange removes all memberships in a hierarchical range (cascade delete subtree).
func (r *groupMemberRepository) DeleteByPathKeyRange(ctx context.Context, pathKeyStart, pathKeyEnd string) error {
	query := model.Map{
		"pathKey": model.Map{
			"$gte": pathKeyStart,
			"$lt":  pathKeyEnd,
		},
	}

	_, err := r.model.DeleteMany(ctx, query)
	if err != nil && err.Error() != "document not found" {
		return err
	}
	return nil
}

// CountByGroupIds returns the number of members for multiple groups.
// Returns a map of groupId -> count for efficient batch querying.
// Uses individual count queries per group (model doesn't support aggregation).
func (r *groupMemberRepository) CountByGroupIds(ctx context.Context, groupIds []string) (map[string]int64, error) {
	if len(groupIds) == 0 {
		return make(map[string]int64), nil
	}

	countMap := make(map[string]int64)

	// Query count for each group
	// Note: This could be optimized with raw MongoDB aggregation if needed
	for _, groupId := range groupIds {
		count, err := r.CountByGroupId(ctx, groupId)
		if err != nil {
			// Log warning but continue with other groups
			continue
		}
		countMap[groupId] = count
	}

	return countMap, nil
}

// FindByGroupIdPaginated retrieves members of a group with pagination.
// Returns members slice, total count, and error.
func (r *groupMemberRepository) FindByGroupIdPaginated(ctx context.Context, groupId string, page, perPage int64) ([]*entities.GroupMember, int64, error) {
	groupObjId, err := model.ToObjectID(groupId)
	if err != nil {
		return nil, 0, err
	}

	query := model.Map{"groupId": groupObjId}

	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    page,
		PerPage: perPage,
	}, nil)
	if err != nil {
		return nil, 0, err
	}

	members := make([]*entities.GroupMember, len(result.Items))
	for i := range result.Items {
		members[i] = &result.Items[i]
	}

	return members, result.Pagination.TotalItems, nil
}

// Compile-time check to ensure groupMemberRepository implements GroupMemberRepository
var _ repositories.GroupMemberRepository = (*groupMemberRepository)(nil)
