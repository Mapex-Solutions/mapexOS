package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/application/ports"
	"mapexIam/src/modules/users/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

func RegisterRoutes(group fiber.Router, service ports.UserServicePort) {

	/**
	* Myself Routes
	 */

	// Get my informations
	group.Get("/me", handlers.Myself(service))

	// Update my informations
	userUpdateDto := validation.NewValidation(&dtos.UserUpdateDTO{}, nil, nil)
	group.Patch("/me",
		validation.ValidationMiddleware(userUpdateDto),
		handlers.UpdateMyself(service),
	)

	// Disable my onboarding tour (no permission required — auth-only)
	group.Patch("/me/tour", handlers.DisableMyTour(service))

	/**
	* CRUD Routes
	 */

	// Get users with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Since users don't have orgId, service queries memberships by orgId to extract userIds
	getUsersQuery := validation.NewValidation(nil, &dtos.UserQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(getUsersQuery),    // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.UserList),    // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                 // 3. Inject context (cache)
		handlers.GetUsers(service),                        // 4. Handler
	)

	// Create a new user
	userCreateDto := validation.NewValidation(&dtos.UserCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(userCreateDto),
		permissionMw.RequirePermission(perms.UserCreate),
		handlers.CreateUser(service),
	)

	// Count users (counter endpoint with cache)
	group.Get("/counter",
		permissionMw.RequirePermission(perms.UserList),
		coverageMw.InjectRequestContext(),
		handlers.GetUserCount(service),
	)

	// Get user by ID
	getUserById := validation.NewValidation(nil, nil, &dtos.UserIdDTO{})
	group.Get("/:userId",
		validation.ValidationMiddleware(getUserById),
		permissionMw.RequirePermission(perms.UserRead),
		handlers.GetUserById(service),
	)

	// Update user by ID
	updateUserById := validation.NewValidation(&dtos.UserUpdateDTO{}, nil, &dtos.UserIdDTO{})
	group.Patch("/:userId",
		validation.ValidationMiddleware(updateUserById),
		permissionMw.RequirePermission(perms.UserUpdate),
		handlers.UpdateUserById(service),
	)

	// Delete user by ID
	deleteUserById := validation.NewValidation(nil, nil, &dtos.UserIdDTO{})
	group.Delete("/:userId",
		validation.ValidationMiddleware(deleteUserById),
		permissionMw.RequirePermission(perms.UserDelete),
		handlers.DeleteUserById(service),
	)
}
