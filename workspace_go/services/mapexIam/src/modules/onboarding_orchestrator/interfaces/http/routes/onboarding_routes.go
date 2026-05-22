package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"
	"mapexIam/src/modules/onboarding_orchestrator/application/ports"
	"mapexIam/src/modules/onboarding_orchestrator/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

// RegisterRoutes registers all onboarding routes in the provided Fiber router group.
// Includes user creation and update with memberships endpoints.
// Accepts UserOnboardingServicePort interface (Hexagonal Architecture) for loose coupling.
//
// Middleware order (important):
//  1. ValidationMiddleware - Validate DTO first (fail fast)
//  2. RequirePermission - Check user has permission (uses auth cache)
//  3. InjectRequestContext - Inject org context (uses coverage cache)
//  4. Handler - Process request
func RegisterRoutes(group fiber.Router, service ports.UserOnboardingServicePort) {

	/**
	* User Onboarding Routes
	*/

	// Create user with multiple memberships
	// POST /api/v1/onboarding/users
	createUserWithMemberships := validation.NewValidation(&dtos.CreateUserWithMembershipsDto{}, nil, nil)
	group.Post("/users",
		validation.ValidationMiddleware(createUserWithMemberships),
		permissionMw.RequirePermission(perms.UserCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateUserWithMemberships(service),
	)

	// Update user with access configuration
	// PATCH /api/v1/onboarding/users/:userId
	updateUserWithAccess := validation.NewValidation(&dtos.UpdateUserWithAccessDto{}, nil, &dtos.UpdateUserWithAccessParamsDto{})
	group.Patch("/users/:userId",
		validation.ValidationMiddleware(updateUserWithAccess),
		permissionMw.RequirePermission(perms.UserUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateUserWithAccess(service),
	)
}
