package routes

import (
	"mapexIam/src/modules/organizations/application/dtos"
	"mapexIam/src/modules/organizations/application/ports"
	"mapexIam/src/modules/organizations/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	orgHierarchyMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/orghierarchy"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers organization HTTP routes. Base path: /api/v1/organizations.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.OrganizationServicePort) {

	r := swagger.Wrap(group).Tag("Organizations")

	organizationQueryDto := validation.NewValidation(nil, &dtos.OrganizationQueryDto{}, nil)
	r.Get("/", organizationQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OrganizationList),
		coverageMw.InjectRequestContext(),
		handlers.GetOrganizations(service),
	).
		Summary("List organizations").
		Description("Returns a paginated, filterable list of organizations scoped to the caller's organization context.").
		Returns(&model.PaginatedResult[dtos.OrganizationResponse]{})

	treeQueryDto := validation.NewValidation(nil, &dtos.TreeQueryDto{}, nil)
	r.Get("/tree", treeQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OrganizationList),
		handlers.GetOrganizationsTree(service),
	).
		Summary("Get organization tree").
		Description("Returns the hierarchical organization tree rooted at the caller's organization context.").
		Returns(&dtos.TreeResponseDto{})

	organizationCreateDto := validation.NewValidation(&dtos.CreateOrganizationDto{}, nil, nil)
	r.Post("/", organizationCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OrganizationCreate),
		coverageMw.InjectRequestContext(),
		orgHierarchyMw.ValidateOrgHierarchy(func(c *web.Ctx) (orgHierarchyMw.OrganizationCreateContract, error) {
			return validation.GetDTO[*dtos.CreateOrganizationDto](c, "bodyDTO")
		}),
		handlers.CreateOrganization(service),
	).
		Summary("Create organization").
		Description("Creates a new organization. The org-hierarchy middleware validates the parent placement before creation.").
		Returns(&dtos.OrganizationResponse{})

	getOrganizationById := validation.NewValidation(nil, nil, &dtos.OrganizationIdDto{})
	r.Get("/:organizationId", getOrganizationById, swagger.Expose,
		permissionMw.RequirePermission(perms.OrganizationRead),
		handlers.GetOrganizationById(service),
	).
		Summary("Get organization by ID").
		Description("Retrieves a single organization by its MongoDB ObjectId.").
		Returns(&dtos.OrganizationResponse{})

	updateOrganizationById := validation.NewValidation(&dtos.UpdateOrganizationDto{}, nil, &dtos.OrganizationIdDto{})
	r.Patch("/:organizationId", updateOrganizationById, swagger.Expose,
		permissionMw.RequirePermission(perms.OrganizationUpdate),
		handlers.UpdateOrganizationById(service),
	).
		Summary("Update organization").
		Description("Partially updates an existing organization. All body fields are optional.").
		Returns(&dtos.OrganizationResponse{})

	deleteOrganizationById := validation.NewValidation(nil, nil, &dtos.OrganizationIdDto{})
	r.Delete("/:organizationId", deleteOrganizationById, swagger.Expose,
		permissionMw.RequirePermission(perms.OrganizationDelete),
		handlers.DeleteOrganizationById(service),
	).
		Summary("Delete organization").
		Description("Deletes an organization by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
