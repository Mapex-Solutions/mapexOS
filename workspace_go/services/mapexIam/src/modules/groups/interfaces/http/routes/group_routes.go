package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/application/ports"
	"mapexIam/src/modules/groups/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

// RegisterRoutes registers all HTTP routes for the groups module.
// Accepts GroupServicePort interface (Hexagonal Architecture) for loose coupling.
func RegisterRoutes(group fiber.Router, service ports.GroupServicePort) {

	/**
	* CRUD Routes
	 */

	// Get groups with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Includes system groups which are available to all users
	groupQueryDto := validation.NewValidation(nil, &dtos.GroupQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(groupQueryDto),  // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.GroupList), // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),               // 3. Inject context (cache)
		handlers.GetGroups(service),                     // 4. Handler
	)

	// Create a new group
	groupCreateDto := validation.NewValidation(&dtos.CreateGroupDto{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(groupCreateDto),     // 1. Validate DTO
		permissionMw.RequirePermission(perms.GroupCreate),   // 2. Check permission
		coverageMw.InjectRequestContext(),                   // 3. Inject context
		handlers.CreateGroup(service),                       // 4. Handler
	)

	// Count groups (counter endpoint with cache)
	group.Get("/counter",
		permissionMw.RequirePermission(perms.GroupList),
		coverageMw.InjectRequestContext(),
		handlers.GetGroupCount(service),
	)

	// Get group by ID
	getGroupById := validation.NewValidation(nil, nil, &dtos.GroupIdDto{})
	group.Get("/:groupId",
		validation.ValidationMiddleware(getGroupById),
		permissionMw.RequirePermission(perms.GroupRead),
		handlers.GetGroupById(service),
	)

	// Update group by ID
	updateGroupById := validation.NewValidation(&dtos.UpdateGroupDto{}, nil, &dtos.GroupIdDto{})
	group.Patch("/:groupId",
		validation.ValidationMiddleware(updateGroupById),
		permissionMw.RequirePermission(perms.GroupUpdate),
		handlers.UpdateGroupById(service),
	)

	// Delete group by ID
	deleteGroupById := validation.NewValidation(nil, nil, &dtos.GroupIdDto{})
	group.Delete("/:groupId",
		validation.ValidationMiddleware(deleteGroupById),
		permissionMw.RequirePermission(perms.GroupDelete),
		handlers.DeleteGroupById(service),
	)

	// Get group members (paginated - max 100 per page)
	getGroupMembers := validation.NewValidation(nil, &dtos.GroupMembersQueryDto{}, &dtos.GroupIdDto{})
	group.Get("/:groupId/members",
		validation.ValidationMiddleware(getGroupMembers),
		permissionMw.RequirePermission(perms.GroupRead),
		handlers.GetGroupMembers(service),
	)

	// Add member to group
	addGroupMember := validation.NewValidation(&dtos.GroupMemberAddDto{}, nil, &dtos.GroupIdDto{})
	group.Post("/:groupId/members",
		validation.ValidationMiddleware(addGroupMember),
		permissionMw.RequirePermission(perms.GroupUpdate),
		handlers.AddGroupMember(service),
	)

	// Remove member from group
	removeGroupMember := validation.NewValidation(nil, nil, &dtos.GroupMemberIdDto{})
	group.Delete("/:groupId/members/:userId",
		validation.ValidationMiddleware(removeGroupMember),
		permissionMw.RequirePermission(perms.GroupUpdate),
		handlers.RemoveGroupMember(service),
	)
}
