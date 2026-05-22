package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"
	"mapexIam/src/modules/onboarding_orchestrator/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// CreateUserWithMemberships returns a Fiber handler that creates a user with multiple memberships atomically.
//
// This is the main onboarding endpoint that orchestrates user creation with organization/role assignments
// in a single atomic operation. It ensures data consistency by rolling back all changes if any step fails.
//
// It expects a validated DTO of type dtos.CreateUserWithMembershipsDto in the "bodyDTO" context key
// (populated by requestValidation middleware) containing:
//   - User data (email, firstName, lastName, password, enabled)
//   - Memberships array (roles for direct assignment) OR Groups array (groupId for group assignment)
//
// Organization Context:
//   - OrgID comes from RequestContext.OrgContext (current selected org)
//   - Scope comes from Organization.AccessPolicy.DefaultScope (centralized config)
//
// Orchestration Flow:
//   1. Gets OrgID from context and Scope from org config
//   2. Creates user in Users service
//   3. Creates membership (direct) or adds to group
//   4. Returns complete user data with created membership IDs
//
// Atomic Transaction:
//   - All operations run in a MongoDB transaction
//   - If any step fails: All changes are rolled back
//
// Returns:
//   - 201 Created with user data and membership IDs
//   - 400 Bad Request if validation fails, user already exists, or no org context
//   - 500 Internal Server Error on service failure
func CreateUserWithMemberships(service ports.UserOnboardingServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Retrieve the timeout-aware Context set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from Fiber locals (set by InjectRequestContext middleware)
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok || requestContext == nil {
			return &customErrors.ServerCustomError{
				Code:   status.BAD_REQUEST,
				Errors: []string{"Request context not found"},
			}
		}

		// Validate org context exists
		if requestContext.OrgContext == nil || *requestContext.OrgContext == "" {
			return &customErrors.ServerCustomError{
				Code:   status.BAD_REQUEST,
				Errors: []string{"Organization context is required. Please select an organization."},
			}
		}

		// Get validated DTO from request body
		bodyData, _ := requestValidation.GetDTO[*dtos.CreateUserWithMembershipsDto](c, "bodyDTO")

		// Call the Application Service to orchestrate user + memberships creation
		retData, err := service.CreateUserWithMemberships(ctx, requestContext, bodyData)

		if err != nil {
			return err
		}

		return response.Created(c, retData)
	}
}

// UpdateUserWithAccess returns a Fiber handler that updates a user and replaces their access configuration atomically.
//
// This endpoint orchestrates user data updates with access configuration changes (memberships/groups)
// in a single atomic operation. It ensures data consistency by rolling back all changes if any step fails.
//
// It expects:
//   - Path parameter "userId" (validated by params middleware)
//   - A validated DTO of type dtos.UpdateUserWithAccessDto in the "bodyDTO" context key
//     (populated by requestValidation middleware) containing:
//   - User data updates (firstName, lastName, phone, jobTitle, enabled, avatar, password, changePasswordNextLogin) - all optional
//   - Access configuration: Memberships array (roles for direct assignment) OR Groups array (groupId for group assignment)
//
// Organization Context:
//   - OrgID comes from RequestContext.OrgContext (current selected org)
//   - Scope comes from Organization.AccessPolicy.DefaultScope (centralized config)
//
// Orchestration Flow:
//  1. Gets OrgID from context and Scope from org config
//  2. Updates user data (only provided fields)
//  3. Removes existing memberships and group memberships for user in current org
//  4. Creates new membership (direct) or adds to new group
//  5. Returns complete user data with new membership IDs
//
// Atomic Transaction:
//   - All operations run in a MongoDB transaction
//   - If any step fails: All changes are rolled back
//
// Returns:
//   - 200 OK with updated user data and new membership IDs
//   - 400 Bad Request if validation fails, user not found, or no org context
//   - 500 Internal Server Error on service failure
func UpdateUserWithAccess(service ports.UserOnboardingServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Retrieve the timeout-aware Context set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from Fiber locals (set by InjectRequestContext middleware)
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok || requestContext == nil {
			return &customErrors.ServerCustomError{
				Code:   status.BAD_REQUEST,
				Errors: []string{"Request context not found"},
			}
		}

		// Validate org context exists
		if requestContext.OrgContext == nil || *requestContext.OrgContext == "" {
			return &customErrors.ServerCustomError{
				Code:   status.BAD_REQUEST,
				Errors: []string{"Organization context is required. Please select an organization."},
			}
		}

		// Get validated path params
		pathParams, _ := requestValidation.GetDTO[*dtos.UpdateUserWithAccessParamsDto](c, "paramsDTO")

		// Get validated DTO from request body
		bodyData, _ := requestValidation.GetDTO[*dtos.UpdateUserWithAccessDto](c, "bodyDTO")

		// Call the Application Service to orchestrate user update + access replacement
		retData, err := service.UpdateUserWithAccess(ctx, requestContext, pathParams.UserID, bodyData)

		if err != nil {
			return err
		}

		return response.Success(c, retData)
	}
}
