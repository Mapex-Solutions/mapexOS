package routes

import (
	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"
	"mapexIam/src/modules/onboarding_orchestrator/application/ports"
	"mapexIam/src/modules/onboarding_orchestrator/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers user-onboarding orchestration routes. Base path:
// /api/v1/onboarding.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.UserOnboardingServicePort) {

	r := swagger.Wrap(group).Tag("Onboarding")

	createUserWithMemberships := validation.NewValidation(&dtos.CreateUserWithMembershipsDto{}, nil, nil)
	r.Post("/users", createUserWithMemberships, swagger.Expose,
		permissionMw.RequirePermission(perms.UserCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateUserWithMemberships(service),
	).
		Summary("Onboard user with memberships").
		Description("Provisions a user together with one or more organization memberships (and roles) in a single orchestrated call.").
		Returns(&dtos.UserOnboardingResponse{})

	updateUserWithAccess := validation.NewValidation(&dtos.UpdateUserWithAccessDto{}, nil, &dtos.UpdateUserWithAccessParamsDto{})
	r.Patch("/users/:userId", updateUserWithAccess, swagger.Expose,
		permissionMw.RequirePermission(perms.UserUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateUserWithAccess(service),
	).
		Summary("Update user with access").
		Description("Updates a user along with their organization access (memberships and roles) in a single orchestrated call.").
		Returns(&dtos.UserOnboardingResponse{})
}
