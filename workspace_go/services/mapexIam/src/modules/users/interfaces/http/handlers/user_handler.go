package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateUser returns a Fiber handler that creates a new user.
//
// It expects a validated DTO of type dtos.UserCreateDTO to be stored
// in the Fiber context under the key "bodyDTO" (usually populated by
// requestValidation middleware).
func CreateUser(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.UserCreateDTO](c, "bodyDTO")
		retData, err := service.CreateUser(ctx, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// GetUserById returns a Fiber handler that retrieves a user by their unique identifier.
//
// It expects a validated DTO of type dtos.UserIdDTO to be stored in the Fiber
// context under the key "paramsDTO" (usually populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with user data if found
//   - 404 Not Found if user doesn't exist
//   - 500 Internal Server Error on service failure
func GetUserById(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		user, _ := requestValidation.GetDTO[*dtos.UserIdDTO](c, "paramsDTO")
		retData, err := service.GetUserById(ctx, &user.UserId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateUserById returns a Fiber handler that updates an existing user's information.
//
// It expects validated DTOs:
//   - dtos.UserIdDTO in "paramsDTO" for the user identifier
//   - dtos.UserUpdateDTO in "bodyDTO" for the fields to update
//
// Both DTOs are populated by requestValidation middleware.
//
// Returns:
//   - 200 OK with updated user data
//   - 404 Not Found if user doesn't exist
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func UpdateUserById(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		user, _ := requestValidation.GetDTO[*dtos.UserIdDTO](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.UserUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateUserById(ctx, &user.UserId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// DeleteUserById returns a Fiber handler that permanently deletes a user.
//
// It expects a validated DTO of type dtos.UserIdDTO to be stored in the Fiber
// context under the key "paramsDTO" (usually populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with success flag if deletion succeeds
//   - 404 Not Found if user doesn't exist
//   - 500 Internal Server Error on service failure
func DeleteUserById(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		user, _ := requestValidation.GetDTO[*dtos.UserIdDTO](c, "paramsDTO")
		retData, err := service.DeleteUserById(ctx, &user.UserId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// Myself returns a Fiber handler that retrieves the current authenticated user's information.
//
// It extracts the user ID from the JWT claims stored in the Fiber context
// (populated by auth middleware) and retrieves the user's data from the service layer.
func Myself(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get the user from JWT
		userId, ok := authmw.GetUserIdFromToken(c)
		if !ok {
			return response.InternalServerError(c, "invalid claims type", nil)
		}

		retData, err := service.GetUserById(ctx, &userId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateMyself returns a Fiber handler that updates the current authenticated user's information.
//
// It extracts the user ID from the JWT claims (populated by auth middleware) and
// expects a validated DTO of type dtos.UserUpdateDTO in the "bodyDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with updated user data
//   - 400 Bad Request if validation fails
//   - 401 Unauthorized if user ID cannot be extracted from token
//   - 404 Not Found if user doesn't exist
//   - 500 Internal Server Error on service failure
func UpdateMyself(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.UserUpdateDTO](c, "bodyDTO")
		userId, _ := authmw.GetUserIdFromToken(c)

		retData, err := service.UpdateUserById(ctx, &userId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// DisableMyTour returns a Fiber handler that disables the onboarding tour for the
// current authenticated user by setting startTour to false.
//
// It extracts the user ID from the JWT claims (populated by auth middleware).
// No request body is needed — the endpoint always sets startTour=false.
//
// Returns:
//   - 200 OK with updated user data
//   - 401 Unauthorized if user ID cannot be extracted from token
//   - 404 Not Found if user doesn't exist
//   - 500 Internal Server Error on service failure
func DisableMyTour(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		userId, ok := authmw.GetUserIdFromToken(c)
		if !ok {
			return response.InternalServerError(c, "invalid claims type", nil)
		}

		startTour := false
		dto := &dtos.UserUpdateDTO{
			StartTour: &startTour,
		}

		retData, err := service.UpdateUserById(ctx, &userId, dto)
		if err != nil {
			return err
		}

		return response.Success(c, retData)
	}
}

// GetUsers returns a Fiber handler that retrieves a paginated and filtered list of users.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query users within their accessible organizations with zero extra queries.
// Since users don't have orgId field, the service queries memberships by orgId to extract userIds.
//
// It expects a validated DTO of type dtos.UserQueryDto in the "queryDTO" context key
// (populated by requestValidation middleware) containing optional filters such as:
//   - email, firstName, lastName (partial match filters)
//   - enabled (boolean filter)
//   - page, perPage (pagination)
//   - sort (sorting option)
//   - projection (field selection)
//
// Returns:
//   - 200 OK with paginated user list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure or requestContext not found
func GetUsers(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.UserQueryDto](c, "queryDTO")
		retData, err := service.GetUsers(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetUserCount returns a Fiber handler that returns the total count of users.
// Uses cached count with 6h TTL, invalidated on create/delete.
//
// Parameters:
//   - service: The UserServicePort interface for user business operations
//
// Returns:
//   - A Fiber handler function that processes the user count request
func GetUserCount(service ports.UserServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountUsers(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}
