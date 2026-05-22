package ports

import (
	"context"

	"mapexIam/src/modules/roles/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// RoleServicePort defines the output port (Hexagonal Architecture) for role-related operations.
// This port allows other application services (like PermissionPopulationService) to depend on
// role operations without coupling to the concrete RoleService implementation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (defines the contract)
//   - Adapter: RoleService (implements the contract)
//
// Used by:
//   - PermissionPopulationService (for resolving role permissions during auth cache population)
//   - HTTP Handlers (for HTTP interface layer)
type RoleServicePort interface {
	// CreateRole creates a new role with proper multi-tenant hierarchical fields.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Role creation data
	//
	// Returns:
	//   - *dtos.RoleResponse: Created role data
	//   - error: Error if creation fails
	CreateRole(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateRoleDto) (*dtos.RoleResponse, error)

	// GetRoleById retrieves a role entity by its unique identifier.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - roleId: Unique identifier of the role
	//
	// Returns:
	//   - *dtos.RoleResponse: Role data DTO
	//   - error: Error if role not found or database error
	GetRoleById(ctx context.Context, roleId *string) (*dtos.RoleResponse, error)

	// UpdateRoleById updates an existing role entity.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - roleId: Unique identifier of the role
	//   - dto: Fields to update
	//
	// Returns:
	//   - *dtos.RoleResponse: Updated role data
	//   - error: Error if update fails
	UpdateRoleById(ctx context.Context, roleId *string, dto *dtos.UpdateRoleDto) (*dtos.RoleResponse, error)

	// DeleteRoleById removes a role entity by its unique ID.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - roleId: Unique identifier of the role
	//
	// Returns:
	//   - map[string]bool: Success status
	//   - error: Error if deletion fails
	DeleteRoleById(ctx context.Context, roleId *string) (map[string]bool, error)

	// GetRoles retrieves a paginated and filtered list of roles.
	// Uses RequestContext for context-aware organization filtering with hierarchical role inheritance.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Filters, pagination, and projection options
	//
	// Returns:
	//   - *model.PaginatedResult[dtos.RoleResponse]: Paginated roles with DTOs
	//   - error: Error if query fails
	GetRoles(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.RoleQueryDto) (*model.PaginatedResult[dtos.RoleResponse], error)
}
