package ports

import (
	ctx "context"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// RouteGroupServicePort defines the contract for route group business operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. Handlers and routes depend on this interface
// rather than the concrete service implementation.
//
// Methods:
//   - CreateRouteGroup: Creates a new route group
//   - GetRouteGroupById: Retrieves a route group by its ID
//   - UpdateRouteGroupById: Updates an existing route group
//   - DeleteRouteGroupById: Deletes a route group
type RouteGroupServicePort interface {
	// CreateRouteGroup creates a new route group from the provided DTO.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Route group creation data
	//
	// Returns:
	//   - RouteGroupResponse: The created route group
	//   - error: If creation fails
	CreateRouteGroup(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.RouteGroupCreateDTO) (*dtos.RouteGroupResponse, error)

	// GetRouteGroupById retrieves a route group by its unique identifier.
	//
	// This method uses a cache-aside pattern: checks cache first, then DB if cache miss.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - routeGroupId: The unique identifier of the route group
	//
	// Returns:
	//   - RouteGroupResponse: The route group if found
	//   - error: If not found or retrieval fails
	GetRouteGroupById(ctx ctx.Context, routeGroupId *string) (*dtos.RouteGroupResponse, error)

	// UpdateRouteGroupById updates an existing route group.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - routeGroupId: The unique identifier of the route group
	//   - dto: Update data (only non-nil fields will be updated)
	//
	// Returns:
	//   - RouteGroupResponse: The updated route group
	//   - error: If update fails or route group not found
	UpdateRouteGroupById(ctx ctx.Context, routeGroupId *string, dto *dtos.RouteGroupUpdateDTO) (*dtos.RouteGroupResponse, error)

	// GetRouteGroups retrieves a paginated and filtered list of route groups.
	// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
	//
	// This method supports multi-tenant isolation by filtering results based on the
	// RequestContext parameter provided by the InjectRequestContext middleware.
	// Only route groups belonging to organizations accessible by the user will be returned.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Request context with org access data (from coverage middleware)
	//   - query: Filters, pagination, and projection options
	//
	// Returns:
	//   - PaginatedResult: Matching route groups and pagination metadata
	//   - error: If query fails
	GetRouteGroups(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.RouteGroupQueryDTO) (*model.PaginatedResult[dtos.RouteGroupResponse], error)

	// DeleteRouteGroupById removes a route group by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - routeGroupId: The unique identifier of the route group to delete
	//
	// Returns:
	//   - map[string]bool: Success indicator ({"success": true})
	//   - error: If deletion fails or route group not found
	DeleteRouteGroupById(ctx ctx.Context, routeGroupId *string) (map[string]bool, error)

	// GetRouteGroupsByIds retrieves multiple route groups by their IDs.
	// This method is designed for internal API calls (MS-to-MS communication).
	//
	// It iterates over the provided IDs and calls GetRouteGroupById for each,
	// leveraging the existing cache-aside pattern for optimal performance.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - ids: Array of route group IDs to retrieve
	//
	// Returns:
	//   - []RouteGroupResponse: Array of found route groups (skips not found)
	//   - error: If any critical error occurs
	GetRouteGroupsByIds(ctx ctx.Context, ids []string) ([]dtos.RouteGroupResponse, error)

	// CountRouteGroups returns the total count of route groups for the given org context.
	// Implements cache-aside: check Redis first, fallback to MongoDB CountDocuments.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains org access data from coverage middleware
	//
	// Returns:
	//   - int64: Total count of matching route groups
	//   - error: If query fails
	CountRouteGroups(ctx ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)

	// GetRouteGroupEntityById retrieves a route group as a domain entity (bson).
	// Used internally by services that process data within the MS — no DTO conversion overhead.
	GetRouteGroupEntityById(ctx ctx.Context, routeGroupId *string) (*entities.RouteGroup, error)
}
