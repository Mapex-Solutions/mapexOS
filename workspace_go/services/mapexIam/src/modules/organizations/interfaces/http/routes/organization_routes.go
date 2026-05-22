package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/organizations/application/dtos"
	"mapexIam/src/modules/organizations/application/ports"
	"mapexIam/src/modules/organizations/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	orgHierarchyMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/orghierarchy"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

// RegisterRoutes registers all HTTP routes for the organizations module.
// Accepts OrganizationServicePort interface (Hexagonal Architecture) for loose coupling.
func RegisterRoutes(group fiber.Router, service ports.OrganizationServicePort) {

	/**
	* CRUD Routes
	 */

	// Get organizations with filters, pagination, and projection
	// Uses coverage middleware to inject RequestContext for context-aware org filtering
	organizationQueryDto := validation.NewValidation(nil, &dtos.OrganizationQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(organizationQueryDto),
		permissionMw.RequirePermission(perms.OrganizationList),
		coverageMw.InjectRequestContext(),
		handlers.GetOrganizations(service),
	)

	// Get organizations tree with cursor pagination (for UI navigation)
	treeQueryDto := validation.NewValidation(nil, &dtos.TreeQueryDto{}, nil)
	group.Get("/tree",
		validation.ValidationMiddleware(treeQueryDto),
		permissionMw.RequirePermission(perms.OrganizationList),
		handlers.GetOrganizationsTree(service),
	)

	// Create a new organization
	// Chain: Validation → Permission → Coverage → Hierarchy Validation → Handler
	organizationCreateDto := validation.NewValidation(&dtos.CreateOrganizationDto{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(organizationCreateDto),
		permissionMw.RequirePermission(perms.OrganizationCreate),
		coverageMw.InjectRequestContext(),
		orgHierarchyMw.ValidateOrgHierarchy(func(c *fiber.Ctx) (orgHierarchyMw.OrganizationCreateContract, error) {
			return validation.GetDTO[*dtos.CreateOrganizationDto](c, "bodyDTO")
		}),
		handlers.CreateOrganization(service),
	)

	// Get organization by ID
	getOrganizationById := validation.NewValidation(nil, nil, &dtos.OrganizationIdDto{})
	group.Get("/:organizationId",
		validation.ValidationMiddleware(getOrganizationById),
		permissionMw.RequirePermission(perms.OrganizationRead),
		handlers.GetOrganizationById(service),
	)

	// Update organization by ID
	updateOrganizationById := validation.NewValidation(&dtos.UpdateOrganizationDto{}, nil, &dtos.OrganizationIdDto{})
	group.Patch("/:organizationId",
		validation.ValidationMiddleware(updateOrganizationById),
		permissionMw.RequirePermission(perms.OrganizationUpdate),
		handlers.UpdateOrganizationById(service),
	)

	// Delete organization by ID
	deleteOrganizationById := validation.NewValidation(nil, nil, &dtos.OrganizationIdDto{})
	group.Delete("/:organizationId",
		validation.ValidationMiddleware(deleteOrganizationById),
		permissionMw.RequirePermission(perms.OrganizationDelete),
		handlers.DeleteOrganizationById(service),
	)
}
