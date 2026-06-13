package routes

import (
	"mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/memberships/application/ports"
	"mapexIam/src/modules/memberships/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers membership HTTP routes. Base path: /api/v1/memberships.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.MembershipServicePort) {

	r := swagger.Wrap(group).Tag("Memberships")

	membershipQueryDto := validation.NewValidation(nil, &dtos.MembershipQueryDto{}, nil)
	r.Get("/", membershipQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.MembershipList),
		coverageMw.InjectRequestContext(),
		handlers.GetMemberships(service),
	).
		Summary("List memberships").
		Description("Returns a paginated, filterable list of memberships (the user↔organization relation) scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.MembershipResponse]{})

	membershipCreateDto := validation.NewValidation(&dtos.CreateMembershipDto{}, nil, nil)
	r.Post("/", membershipCreateDto, swagger.Expose, handlers.CreateMembership(service)).
		Summary("Create membership").
		Description("Creates a membership binding a user to an organization with one or more roles.").
		Returns(&dtos.MembershipResponse{})

	getMembershipById := validation.NewValidation(nil, nil, &dtos.MembershipIdDto{})
	r.Get("/:membershipId", getMembershipById, swagger.Expose, handlers.GetMembershipById(service)).
		Summary("Get membership by ID").
		Description("Retrieves a single membership by its MongoDB ObjectId.").
		Returns(&dtos.MembershipResponse{})

	updateMembershipById := validation.NewValidation(&dtos.UpdateMembershipDto{}, nil, &dtos.MembershipIdDto{})
	r.Patch("/:membershipId", updateMembershipById, swagger.Expose, handlers.UpdateMembershipById(service)).
		Summary("Update membership").
		Description("Partially updates an existing membership (e.g. its role set). All body fields are optional.").
		Returns(&dtos.MembershipResponse{})

	deleteMembershipById := validation.NewValidation(nil, nil, &dtos.MembershipIdDto{})
	r.Delete("/:membershipId", deleteMembershipById, swagger.Expose, handlers.DeleteMembershipById(service)).
		Summary("Delete membership").
		Description("Deletes a membership by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}

// RegisterMeRoutes registers the authenticated caller's own context routes.
// Base path: /api/v1/me.
func RegisterMeRoutes(group web.Router, service ports.MembershipServicePort) {

	r := swagger.Wrap(group).Tag("Me")

	noContract := validation.NewValidation(nil, nil, nil)
	r.Get("/coverage", noContract, swagger.Expose, handlers.GetMeCoverage(service)).
		Summary("Get my coverage").
		Description("Returns the authenticated caller's authorization coverage — the organizations and scopes the caller can act on.").
		Returns(&dtos.MeCoverageResponse{})
}
