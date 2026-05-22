package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/roles/application/dtos"
	"mapexIam/src/modules/roles/application/ports"
	"mapexIam/src/modules/roles/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

func RegisterRoutes(group fiber.Router, service ports.RoleServicePort) {

	/**
	* CRUD Routes
	 */

	// Get roles with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Implements hierarchical role inheritance (system + MAPEX exclusive + local + global from ancestors)
	getRolesQuery := validation.NewValidation(nil, &dtos.RoleQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(getRolesQuery),  // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RoleList),  // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),               // 3. Inject context (cache)
		handlers.GetRoles(service),                      // 4. Handler
	)

	// Create a new role
	roleCreateDto := validation.NewValidation(&dtos.CreateRoleDto{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(roleCreateDto),     // 1. Validate DTO
		permissionMw.RequirePermission(perms.RoleCreate),   // 2. Check permission
		coverageMw.InjectRequestContext(),                  // 3. Inject context
		handlers.CreateRole(service),                       // 4. Handler
	)

	// Get role by ID
	getRoleById := validation.NewValidation(nil, nil, &dtos.RoleIdDto{})
	group.Get("/:roleId",
		validation.ValidationMiddleware(getRoleById),
		permissionMw.RequirePermission(perms.RoleRead),
		handlers.GetRoleById(service),
	)

	// Update role by ID
	updateRoleById := validation.NewValidation(&dtos.UpdateRoleDto{}, nil, &dtos.RoleIdDto{})
	group.Patch("/:roleId",
		validation.ValidationMiddleware(updateRoleById),
		permissionMw.RequirePermission(perms.RoleUpdate),
		handlers.UpdateRoleById(service),
	)

	// Delete role by ID
	deleteRoleById := validation.NewValidation(nil, nil, &dtos.RoleIdDto{})
	group.Delete("/:roleId",
		validation.ValidationMiddleware(deleteRoleById),
		permissionMw.RequirePermission(perms.RoleDelete),
		handlers.DeleteRoleById(service),
	)
}
