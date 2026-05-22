package repositories

import (
	"context"
	"mapexIam/src/modules/groups/domain/entities"
)

// GroupMemberRepository defines the interface for managing GroupMember entities.
// This repository handles the junction table between groups and users,
// replacing the embedded Members array for better scalability.
type GroupMemberRepository interface {
	// Create inserts a new GroupMember entity.
	Create(ctx context.Context, member *entities.GroupMember) (*entities.GroupMember, error)

	// DeleteByGroupAndUser removes a specific group membership.
	DeleteByGroupAndUser(ctx context.Context, groupId, userId string) error

	// FindByGroupAndUser retrieves a specific group membership.
	FindByGroupAndUser(ctx context.Context, groupId, userId string) (*entities.GroupMember, error)

	// FindByGroupId retrieves all members of a specific group.
	FindByGroupId(ctx context.Context, groupId string) ([]*entities.GroupMember, error)

	// FindByGroupIds retrieves all members of multiple groups (bulk query).
	// Used for efficient org-based user listing (users via group memberships).
	FindByGroupIds(ctx context.Context, groupIds []string) ([]*entities.GroupMember, error)

	// FindByUserId retrieves all group memberships for a specific user.
	FindByUserId(ctx context.Context, userId string) ([]*entities.GroupMember, error)

	// FindByUserIdAndOrgId retrieves all group memberships for a user in a specific organization.
	FindByUserIdAndOrgId(ctx context.Context, userId, orgId string) ([]*entities.GroupMember, error)

	// CountByGroupId returns the number of members in a specific group.
	CountByGroupId(ctx context.Context, groupId string) (int64, error)

	// CountByGroupIds returns the number of members for multiple groups (batch query).
	// Returns a map of groupId -> count for efficient listing.
	CountByGroupIds(ctx context.Context, groupIds []string) (map[string]int64, error)

	// FindByGroupIdPaginated retrieves members of a group with pagination.
	// Used for efficient listing when groups have many members.
	FindByGroupIdPaginated(ctx context.Context, groupId string, page, perPage int64) ([]*entities.GroupMember, int64, error)

	// DeleteByGroupId removes all memberships for a specific group (cascade delete).
	DeleteByGroupId(ctx context.Context, groupId string) error

	// DeleteByUserId removes all group memberships for a specific user.
	DeleteByUserId(ctx context.Context, userId string) error

	// FindByPathKeyRange retrieves all memberships in a hierarchical range (base36 $gte/$lt).
	FindByPathKeyRange(ctx context.Context, pathKeyStart, pathKeyEnd string) ([]*entities.GroupMember, error)

	// DeleteByPathKeyRange removes all memberships in a hierarchical range (cascade delete subtree).
	DeleteByPathKeyRange(ctx context.Context, pathKeyStart, pathKeyEnd string) error
}
