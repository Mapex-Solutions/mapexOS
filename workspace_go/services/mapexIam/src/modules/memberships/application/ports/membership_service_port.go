package ports

import (
	"context"

	"mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/memberships/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// MembershipServicePort defines the output port (Hexagonal Architecture) for membership-related operations.
// This port allows other application services (like AuthorizationCacheService, CoverageCacheService)
// to depend on membership operations without coupling to the concrete MembershipService implementation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (defines the contract)
//   - Adapter: MembershipService (implements the contract)
//
// Used by:
//   - AuthorizationCacheService (for building authorization cache with recursive inheritance)
//   - CoverageCacheService (for building coverage cache)
//   - HTTP Handlers (for HTTP interface layer)
type MembershipServicePort interface {
	// CreateMembership creates a new membership entity.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - dto: Membership creation data
	//
	// Returns:
	//   - *dtos.MembershipResponse: Created membership data
	//   - error: Error if creation fails
	CreateMembership(ctx context.Context, dto *dtos.CreateMembershipDto) (*dtos.MembershipResponse, error)

	// GetMembershipById retrieves a membership entity by its unique identifier.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - membershipId: Unique identifier of the membership
	//
	// Returns:
	//   - *dtos.MembershipResponse: Membership data
	//   - error: Error if not found or database error
	GetMembershipById(ctx context.Context, membershipId *string) (*dtos.MembershipResponse, error)

	// UpdateMembershipById updates an existing membership entity.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - membershipId: Unique identifier of the membership
	//   - dto: Fields to update
	//
	// Returns:
	//   - *dtos.MembershipResponse: Updated membership data
	//   - error: Error if update fails
	UpdateMembershipById(ctx context.Context, membershipId *string, dto *dtos.UpdateMembershipDto) (*dtos.MembershipResponse, error)

	// DeleteMembershipById removes a membership entity by its unique ID.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - membershipId: Unique identifier of the membership
	//
	// Returns:
	//   - map[string]bool: Success status
	//   - error: Error if deletion fails
	DeleteMembershipById(ctx context.Context, membershipId *string) (map[string]bool, error)

	// GetMemberships retrieves a paginated and filtered list of Membership entities.
	// Uses RequestContext from coverage middleware for context-aware org filtering.
	// FOR EXTERNAL USE (HTTP routes with InjectRequestContext middleware).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: Request context with org filtering data (from coverage middleware)
	//   - query: Query parameters for filtering, pagination, and projection
	//
	// Returns:
	//   - *model.PaginatedResult[dtos.MembershipResponse]: Paginated membership list
	//   - error: Error if query fails
	GetMemberships(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.MembershipQueryDto) (*model.PaginatedResult[dtos.MembershipResponse], error)

	// GetAllMemberships retrieves ALL membership entities matching filters using cursor pagination.
	// FOR INTERNAL USE (cache services, consumers, cross-module calls).
	// Does NOT apply coverage filtering - returns all matching records.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - query: Query parameters for filtering (NO pagination - uses cursor internally)
	//
	// Returns:
	//   - []*entities.Membership: All matching membership entities
	//   - error: Error if operation fails
	GetAllMemberships(ctx context.Context, query *dtos.MembershipQueryDto) ([]*entities.Membership, error)

	// GetUserCoverage retrieves coverage data for a user (organizations, groups, customers).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: The unique identifier of the user
	//
	// Returns:
	//   - *dtos.MeCoverageResponse: User coverage data
	//   - error: Error if operation fails
	GetUserCoverage(ctx context.Context, userId string) (*dtos.MeCoverageResponse, error)

	// GetAllUserMemberships retrieves ALL membership entities for a specific user using cursor pagination.
	// This method fetches all memberships in batches without any arbitrary hard limits.
	// Designed for cache services to build authorization and coverage caches.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: The unique identifier of the user
	//
	// Returns:
	//   - []*entities.Membership: All membership entities for the user
	//   - error: Error if operation fails
	GetAllUserMemberships(ctx context.Context, userId string) ([]*entities.Membership, error)

	// GetDirectUserMemberships retrieves only direct memberships (assigneeType='user') for a user.
	// Does NOT include memberships inherited via groups.
	// Used by UserService.getUserMemberships for the direct membership path.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: The unique identifier of the user
	//
	// Returns:
	//   - []*entities.Membership: Direct membership entities for the user
	//   - error: Error if operation fails
	GetDirectUserMemberships(ctx context.Context, userId string) ([]*entities.Membership, error)

	// GetMembershipsByGroupIds retrieves memberships where assigneeType='group' and assigneeId is in groupIds.
	// Used by UserService.getUserMemberships to resolve inherited org access via groups.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - groupIds: List of group IDs to query memberships for
	//
	// Returns:
	//   - []*entities.Membership: Group membership entities
	//   - error: Error if operation fails
	GetMembershipsByGroupIds(ctx context.Context, groupIds []string) ([]*entities.Membership, error)

	// GetAssigneeIdsByOrgIds retrieves assignee IDs matching org filter and assignee type.
	// Used for coverage filtering (find userIds or groupIds with memberships in given orgs).
	// Internally uses batched pagination to handle large datasets.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - orgIds: List of organization ObjectIDs to filter by
	//   - assigneeType: "user" or "group"
	//
	// Returns:
	//   - []string: List of unique assignee IDs (hex strings)
	//   - error: Error if operation fails
	GetAssigneeIdsByOrgIds(ctx context.Context, orgIds []model.ObjectId, assigneeType string) ([]string, error)
}
