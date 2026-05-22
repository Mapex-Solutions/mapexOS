package ports

import (
	"context"
)

// GroupQueryServicePort defines the query-only contract for cross-domain group lookups.
// This port exists to break circular dependencies:
//   - GroupService depends on UserService and MembershipService
//   - UserService and MembershipService need group data
//   - GroupQueryService depends ONLY on GroupRepo + GroupMemberRepo (same domain)
//   - Zero circular dependencies!
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (defines the contract)
//   - Adapter: GroupQueryService (implements the contract)
//
// Used by:
//   - UserService (for group enrichment and coverage filtering)
//   - MembershipService (for resolving group memberships)
type GroupQueryServicePort interface {
	// GetGroupBasicInfo retrieves basic group info (ID, Name, Description).
	// Used by UserService to enrich user detail views with group names.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - groupId: The group's unique identifier
	//
	// Returns:
	//   - *GroupBasicInfo: Basic group data (nil if not found)
	//   - error: Error if query fails
	GetGroupBasicInfo(ctx context.Context, groupId string) (*GroupBasicInfo, error)

	// GetAllUserGroupIds retrieves all group IDs a user belongs to across all organizations.
	// Used by MembershipService.GetAllUserMemberships and UserService.getUserGroups.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: The user's unique identifier
	//
	// Returns:
	//   - []string: List of group IDs
	//   - error: Error if query fails
	GetAllUserGroupIds(ctx context.Context, userId string) ([]string, error)

	// GetUserIdsByGroupIds retrieves all user IDs that are members of any of the given groups.
	// Used by UserService.GetUsers for PATH 2 (users via group memberships).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - groupIds: List of group IDs to query
	//
	// Returns:
	//   - []string: List of unique user IDs
	//   - error: Error if query fails
	GetUserIdsByGroupIds(ctx context.Context, groupIds []string) ([]string, error)

	// CountGroupsByUserIds returns a map of userId -> groupCount for batch enrichment.
	// Used by UserService.GetUsers to display groupsCount in list view.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userIds: List of user IDs to count groups for
	//
	// Returns:
	//   - map[string]int: userId -> groupCount mapping
	//   - error: Error if query fails
	CountGroupsByUserIds(ctx context.Context, userIds []string) (map[string]int, error)
}
