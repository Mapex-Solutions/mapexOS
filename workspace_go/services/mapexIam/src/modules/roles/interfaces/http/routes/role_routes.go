package routes

import (
	"mapexIam/src/modules/roles/application/dtos"
	"mapexIam/src/modules/roles/application/ports"
	"mapexIam/src/modules/roles/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers role HTTP routes. Base path: /api/v1/roles.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.RoleServicePort) {

	r := swagger.Wrap(group).Tag("Roles")

	getRolesQuery := validation.NewValidation(nil, &dtos.RoleQueryDto{}, nil)
	r.Get("/", getRolesQuery, swagger.Expose,
		permissionMw.RequirePermission(perms.RoleList),
		coverageMw.InjectRequestContext(),
		handlers.GetRoles(service),
	).
		Summary("List roles").
		Description("Returns a paginated, filterable list of roles scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.RoleResponse]{})

	roleCreateDto := validation.NewValidation(&dtos.CreateRoleDto{}, nil, nil)
	r.Post("/", roleCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RoleCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateRole(service),
	).
		Summary("Create role").
		Description("Creates a new role with its permission set.").
		Returns(&dtos.RoleResponse{})

	getRoleById := validation.NewValidation(nil, nil, &dtos.RoleIdDto{})
	r.Get("/:roleId", getRoleById, swagger.Expose,
		permissionMw.RequirePermission(perms.RoleRead),
		handlers.GetRoleById(service),
	).
		Summary("Get role by ID").
		Description("Retrieves a single role by its MongoDB ObjectId.").
		Returns(&dtos.RoleResponse{})

	updateRoleById := validation.NewValidation(&dtos.UpdateRoleDto{}, nil, &dtos.RoleIdDto{})
	r.Patch("/:roleId", updateRoleById, swagger.Expose,
		permissionMw.RequirePermission(perms.RoleUpdate),
		handlers.UpdateRoleById(service),
	).
		Summary("Update role").
		Description("Partially updates an existing role. All body fields are optional.").
		Returns(&dtos.RoleResponse{})

	deleteRoleById := validation.NewValidation(nil, nil, &dtos.RoleIdDto{})
	r.Delete("/:roleId", deleteRoleById, swagger.Expose,
		permissionMw.RequirePermission(perms.RoleDelete),
		handlers.DeleteRoleById(service),
	).
		Summary("Delete role").
		Description("Deletes a role by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
