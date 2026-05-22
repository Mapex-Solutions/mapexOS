package handlers

import (
	"github.com/gofiber/fiber/v2"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateRouteGroup returns a Fiber handler that creates a new routegroup.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
//
// It expects a validated DTO of type dtos.RouteGroupCreateDTO to be stored
// in the Fiber context under the key "bodyDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The RouteGroupServicePort interface for route group business operations
//
// Returns:
//   - A Fiber handler function that processes the route group creation request
func CreateRouteGroup(service ports.RouteGroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.RouteGroupCreateDTO](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateRouteGroup(ctx, requestContext, bodyData)

		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetRouteGroupById returns a Fiber handler that retrieves a route group by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.RouteGroupIdDTO to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The RouteGroupServicePort interface for route group business operations
//
// Returns:
//   - A Fiber handler function that processes the route group retrieval request
func GetRouteGroupById(service ports.RouteGroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		routegroup, _ := requestValidation.GetDTO[*dtos.RouteGroupIdDTO](c, "paramsDTO")
		retData, err := service.GetRouteGroupById(ctx, &routegroup.RouteGroupId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateRouteGroupById returns a Fiber handler that updates a route group by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects two validated DTOs:
//   - dtos.RouteGroupIdDTO stored in the Fiber context under the key "paramsDTO"
//   - dtos.RouteGroupUpdateDTO stored in the Fiber context under the key "bodyDTO"
//
// (Both are usually populated by requestValidation middleware)
//
// Parameters:
//   - service: The RouteGroupServicePort interface for route group business operations
//
// Returns:
//   - A Fiber handler function that processes the route group update request
func UpdateRouteGroupById(service ports.RouteGroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		routegroup, _ := requestValidation.GetDTO[*dtos.RouteGroupIdDTO](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.RouteGroupUpdateDTO](c, "bodyDTO")

		// Override orgId from requestContext for security (client cannot specify their own orgId)
		if bodyData.OrgId != nil && requestContext.OrgContext != nil {
			bodyData.OrgId = requestContext.OrgContext
		}

		retData, err := service.UpdateRouteGroupById(ctx, &routegroup.RouteGroupId, bodyData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetRouteGroups returns a Fiber handler that retrieves a paginated and filtered list of route groups.
// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
//
// The handler extracts RequestContext from c.Locals("requestContext") which was set by InjectRequestContext middleware.
// This provides access to:
//   - ScopedOrgIds: All accessible organization IDs
//   - OrgContext: Optional org filter from X-Org-Context header
//   - OrgContextData: Detailed org data with PathKey
//   - CoverageOrgs: Full coverage data with hierarchical information
//
// Query supports hierarchical filtering via includeChildren parameter:
//   - OrgContext + includeChildren=true: Returns org and all descendants (PathKey range)
//   - OrgContext + includeChildren=false: Returns specific org only
//   - No OrgContext: Returns all accessible orgs
//
// It expects a validated DTO of type dtos.RouteGroupQueryDTO in the "queryDTO" context key
// (populated by requestValidation middleware) containing optional filters such as:
//   - name, enabled, version (filters)
//   - page, perPage (pagination)
//   - projection (field selection)
//   - includeChildren (hierarchical query flag)
//
// Returns:
//   - 200 OK with paginated route group list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure or requestContext not found
func GetRouteGroups(service ports.RouteGroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.RouteGroupQueryDTO](c, "queryDTO")
		retData, err := service.GetRouteGroups(ctx, requestContext, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteRouteGroupById returns a Fiber handler that deletes a route group by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// Multi-tenant isolation is handled by the coverage middleware before reaching this handler.
//
// It expects a validated DTO of type dtos.RouteGroupIdDTO to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The RouteGroupServicePort interface for route group business operations
//
// Returns:
//   - A Fiber handler function that processes the route group deletion request
func DeleteRouteGroupById(service ports.RouteGroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		routegroup, _ := requestValidation.GetDTO[*dtos.RouteGroupIdDTO](c, "paramsDTO")
		retData, err := service.DeleteRouteGroupById(ctx, &routegroup.RouteGroupId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetRouteGroupCount returns a Fiber handler that returns the total count of route groups.
// Uses cached count with 6h TTL, invalidated on create/delete.
//
// Parameters:
//   - service: The RouteGroupServicePort interface for route group business operations
//
// Returns:
//   - A Fiber handler function that processes the route group count request
func GetRouteGroupCount(service ports.RouteGroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountRouteGroups(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}
