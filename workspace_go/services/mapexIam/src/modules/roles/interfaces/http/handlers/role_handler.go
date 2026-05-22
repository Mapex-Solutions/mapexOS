package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/roles/application/dtos"
	"mapexIam/src/modules/roles/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateRole returns a Fiber handler that creates a new role.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
//
// It expects a validated DTO of type dtos.CreateRoleDto to be stored
// in the Fiber context under the key "bodyDTO" (usually populated by
// requestValidation middleware).
func CreateRole(service ports.RoleServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.CreateRoleDto](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateRole(ctx, requestContext, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// GetRoleById returns a Fiber handler that retrieves a role by its unique identifier.
//
// It expects a validated DTO of type dtos.RoleIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with role data if found
//   - 404 Not Found if role doesn't exist
//   - 500 Internal Server Error on service failure
func GetRoleById(service ports.RoleServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.RoleIdDto](c, "paramsDTO")
		retData, err := service.GetRoleById(ctx, &params.RoleId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateRoleById returns a Fiber handler that updates an existing role's information.
//
// It expects validated DTOs:
//   - dtos.RoleIdDto in "paramsDTO" for the role identifier
//   - dtos.UpdateRoleDto in "bodyDTO" for the fields to update
//
// Both DTOs are populated by requestValidation middleware.
//
// Returns:
//   - 200 OK with updated role data
//   - 404 Not Found if role doesn't exist
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func UpdateRoleById(service ports.RoleServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.RoleIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.UpdateRoleDto](c, "bodyDTO")
		retData, err := service.UpdateRoleById(ctx, &params.RoleId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// DeleteRoleById returns a Fiber handler that permanently deletes a role.
//
// It expects a validated DTO of type dtos.RoleIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with success flag if deletion succeeds
//   - 404 Not Found if role doesn't exist
//   - 500 Internal Server Error on service failure
func DeleteRoleById(service ports.RoleServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.RoleIdDto](c, "paramsDTO")
		retData, err := service.DeleteRoleById(ctx, &params.RoleId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetRoles returns a Fiber handler that retrieves a paginated and filtered list of roles.
// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical role inheritance.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query roles within their accessible organizations with zero extra queries.
// Implements hierarchical role query (system + MAPEX exclusive + local + global from ancestors).
//
// Returns:
//   - 200 OK with paginated role list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure or requestContext not found
func GetRoles(service ports.RoleServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.RoleQueryDto](c, "queryDTO")
		retData, err := service.GetRoles(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}
