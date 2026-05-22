package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/memberships/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// CreateMembership returns a Fiber handler that creates a new membership assignment.
//
// Memberships represent the relationship between assignees (users or groups) and organizations,
// defining their roles and access scope within the organizational hierarchy.
//
// It expects a validated DTO of type dtos.CreateMembershipDto in the "bodyDTO" context key
// (populated by requestValidation middleware) containing:
//   - assigneeId: ID of user or group being assigned
//   - assigneeType: "user" or "group"
//   - orgId: Organization ID for the membership
//   - roleIds: Array of role IDs to assign
//   - scope: "local" or "global" (inheritance behavior)
//   - enabled: Boolean to enable/disable the membership
//
// Returns:
//   - 201 Created with membership data
//   - 400 Bad Request if validation fails or assignee/org/roles don't exist
//   - 500 Internal Server Error on service failure
func CreateMembership(service ports.MembershipServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.CreateMembershipDto](c, "bodyDTO")
		retData, err := service.CreateMembership(ctx, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// GetMembershipById returns a Fiber handler that retrieves a membership by its unique identifier.
//
// It expects a validated DTO of type dtos.MembershipIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with membership data if found
//   - 404 Not Found if membership doesn't exist
//   - 500 Internal Server Error on service failure
func GetMembershipById(service ports.MembershipServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MembershipIdDto](c, "paramsDTO")
		retData, err := service.GetMembershipById(ctx, &params.MembershipId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateMembershipById returns a Fiber handler that updates an existing membership's information.
//
// It expects validated DTOs:
//   - dtos.MembershipIdDto in "paramsDTO" for the membership identifier
//   - dtos.UpdateMembershipDto in "bodyDTO" for the fields to update (roleIds, scope, enabled)
//
// Both DTOs are populated by requestValidation middleware.
//
// Returns:
//   - 200 OK with updated membership data
//   - 404 Not Found if membership doesn't exist
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func UpdateMembershipById(service ports.MembershipServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MembershipIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.UpdateMembershipDto](c, "bodyDTO")
		retData, err := service.UpdateMembershipById(ctx, &params.MembershipId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// DeleteMembershipById returns a Fiber handler that permanently deletes a membership.
//
// It expects a validated DTO of type dtos.MembershipIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with success flag if deletion succeeds
//   - 404 Not Found if membership doesn't exist
//   - 500 Internal Server Error on service failure
func DeleteMembershipById(service ports.MembershipServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MembershipIdDto](c, "paramsDTO")
		retData, err := service.DeleteMembershipById(ctx, &params.MembershipId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetMemberships returns a Fiber handler that retrieves a paginated and filtered list of memberships.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query memberships within their accessible organizations with zero extra queries.
//
// It expects a validated DTO of type dtos.MembershipQueryDto in the "queryDTO" context key
// (populated by requestValidation middleware) containing optional filters such as:
//   - assigneeId: Filter by specific user or group ID
//   - assigneeType: "user" or "group" (filter by assignee type)
//   - userId: Convenience filter for assigneeId when assigneeType="user"
//   - roleId: Filter by role ID
//   - scope: "local" or "recursive" (filter by inheritance behavior)
//   - enabled: Boolean filter for active/inactive memberships
//   - page, perPage: Pagination parameters
//   - projection: Field selection (comma-separated string)
//   - includeChildren: Include child orgs hierarchically
//
// Returns:
//   - 200 OK with paginated membership list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure
func GetMemberships(service ports.MembershipServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.MembershipQueryDto](c, "queryDTO")
		retData, err := service.GetMemberships(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetMeCoverage returns a Fiber handler that retrieves all customers/organizations the authenticated user has access to.
// It extracts the user ID from the JWT token using the standard auth middleware helper and returns the cached coverage data.
func GetMeCoverage(service ports.MembershipServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Extract userId from JWT using standard middleware helper
		userId, ok := authmw.GetUserIdFromToken(c)
		if !ok {
			return response.Custom(c, status.UNAUTHORIZED, []string{"Missing or invalid userId in token"})
		}

		// Call service to get user coverage
		retData, err := service.GetUserCoverage(ctx, userId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}
