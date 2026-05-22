package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/memberships/application/ports"
	"mapexIam/src/modules/memberships/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

func RegisterRoutes(group fiber.Router, service ports.MembershipServicePort) {

	/**
	* CRUD Routes
	 */

	// Get memberships with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Supports hierarchical filtering via includeChildren parameter and X-Org-Context header
	membershipQueryDto := validation.NewValidation(nil, &dtos.MembershipQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(membershipQueryDto),    // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.MembershipList),   // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                      // 3. Inject context (cache)
		handlers.GetMemberships(service),                       // 4. Handler
	)

	// Create a new membership
	membershipCreateDto := validation.NewValidation(&dtos.CreateMembershipDto{}, nil, nil)
	group.Post("/", validation.ValidationMiddleware(membershipCreateDto), handlers.CreateMembership(service))

	// Get membership by ID
	getMembershipById := validation.NewValidation(nil, nil, &dtos.MembershipIdDto{})
	group.Get("/:membershipId", validation.ValidationMiddleware(getMembershipById), handlers.GetMembershipById(service))

	// Update membership by ID
	updateMembershipById := validation.NewValidation(&dtos.UpdateMembershipDto{}, nil, &dtos.MembershipIdDto{})
	group.Patch("/:membershipId", validation.ValidationMiddleware(updateMembershipById), handlers.UpdateMembershipById(service))

	// Delete membership by ID
	deleteMembershipById := validation.NewValidation(nil, nil, &dtos.MembershipIdDto{})
	group.Delete("/:membershipId", validation.ValidationMiddleware(deleteMembershipById), handlers.DeleteMembershipById(service))
}

func RegisterMeRoutes(group fiber.Router, service ports.MembershipServicePort) {

	/**
	* /me Routes - User coverage and context
	 */

	// Get user coverage (customers/organizations accessible by the authenticated user)
	group.Get("/coverage", handlers.GetMeCoverage(service))
}
