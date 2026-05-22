package ports

import (
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	ctx "context"
	"mapexIam/src/modules/groups/application/dtos"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// GroupServicePort defines the contract for group business logic operations.
// This interface follows Hexagonal Architecture (Ports & Adapters) pattern,
// enabling dependency inversion and decoupling from concrete implementations.
//
// The port allows for:
//   - Easy mocking in tests
//   - Swappable implementations
//   - Clear separation between domain logic and infrastructure concerns
type GroupServicePort interface {
	// CreateGroup creates a new group with proper multi-tenant fields.
	// Handles system-wide groups or organization-scoped groups.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Group creation data
	//
	// Returns:
	//   - *dtos.GroupResponse: Created group data
	//   - error: Error if creation fails
	CreateGroup(context ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateGroupDto) (*dtos.GroupResponse, error)

	// GetGroupById retrieves a group entity by its unique identifier.
	GetGroupById(context ctx.Context, groupId *string) (*dtos.GroupResponse, error)

	// UpdateGroupById updates an existing group entity.
	// Invalidates auth cache for affected users.
	UpdateGroupById(context ctx.Context, groupId *string, dto *dtos.UpdateGroupDto) (*dtos.GroupResponse, error)

	// DeleteGroupById removes a group entity by its unique ID.
	// Invalidates auth cache for all group members.
	DeleteGroupById(context ctx.Context, groupId *string) (map[string]bool, error)

	// GetGroups retrieves a paginated and filtered list of Group DTOs.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query parameters (filters, pagination, projection)
	//
	// Returns:
	//   - *model.PaginatedResult[dtos.GroupResponse]: Paginated groups with DTOs
	//   - error: Error if query fails
	GetGroups(context ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.GroupQueryDto) (*model.PaginatedResult[dtos.GroupResponse], error)

	// AddMemberToGroup adds a user to a group's members list.
	// This method prevents duplicate members and is idempotent.
	// Used by onboarding orchestrator to assign users to groups.
	AddMemberToGroup(context ctx.Context, groupId string, userId string) error

	// RemoveMemberFromGroup removes a user from a group's members list.
	// This method is idempotent (no error if user not in group).
	// Used by onboarding orchestrator to remove users from groups.
	RemoveMemberFromGroup(context ctx.Context, groupId string, userId string) error

	// GetUserGroupsInOrg retrieves all groups a user belongs to in a specific organization.
	// Uses the GroupMember junction table for efficient O(log n) queries.
	// Used by onboarding orchestrator to efficiently remove user from all groups.
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - userId: The user's unique identifier
	//   - orgId: The organization's unique identifier
	//
	// Returns:
	//   - []string: List of group IDs the user belongs to in the org
	//   - error: Error if query fails
	GetUserGroupsInOrg(context ctx.Context, userId, orgId string) ([]string, error)

	// GetGroupMembers retrieves a paginated list of members for a group.
	// Supports pagination with configurable page size (max 100).
	// Used for displaying group members in the UI with efficient loading.
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - groupId: The group's unique identifier
	//   - query: Query parameters (page, perPage - max 100)
	//
	// Returns:
	//   - *model.PaginatedResult[dtos.GroupMemberResponse]: Paginated members
	//   - error: Error if query fails
	GetGroupMembers(context ctx.Context, groupId string, query *dtos.GroupMembersQueryDto) (*model.PaginatedResult[dtos.GroupMemberResponse], error)

	// CountGroups returns the total count of groups for the given org context.
	// Implements cache-aside: check Redis first, fallback to MongoDB CountDocuments.
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - requestContext: Contains org access data from coverage middleware
	//
	// Returns:
	//   - int64: Total count of matching groups
	//   - error: If query fails
	CountGroups(context ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)
}
