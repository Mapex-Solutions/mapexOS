package routes

import (
	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/application/ports"
	"mapexIam/src/modules/groups/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers group HTTP routes. Base path: /api/v1/groups.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.GroupServicePort) {

	r := swagger.Wrap(group).Tag("Groups")

	groupQueryDto := validation.NewValidation(nil, &dtos.GroupQueryDto{}, nil)
	r.Get("/", groupQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupList),
		coverageMw.InjectRequestContext(),
		handlers.GetGroups(service),
	).
		Summary("List groups").
		Description("Returns a paginated, filterable list of groups scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.GroupResponse]{})

	groupCreateDto := validation.NewValidation(&dtos.CreateGroupDto{}, nil, nil)
	r.Post("/", groupCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateGroup(service),
	).
		Summary("Create group").
		Description("Creates a new group.").
		Returns(&dtos.GroupResponse{})

	noContract := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", noContract, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupList),
		coverageMw.InjectRequestContext(),
		handlers.GetGroupCount(service),
	).
		Summary("Count groups").
		Description("Returns the total number of groups for the caller's organization.").
		Returns(&contractsCommon.CounterResponse{})

	getGroupById := validation.NewValidation(nil, nil, &dtos.GroupIdDto{})
	r.Get("/:groupId", getGroupById, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupRead),
		handlers.GetGroupById(service),
	).
		Summary("Get group by ID").
		Description("Retrieves a single group by its MongoDB ObjectId.").
		Returns(&dtos.GroupResponse{})

	updateGroupById := validation.NewValidation(&dtos.UpdateGroupDto{}, nil, &dtos.GroupIdDto{})
	r.Patch("/:groupId", updateGroupById, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupUpdate),
		handlers.UpdateGroupById(service),
	).
		Summary("Update group").
		Description("Partially updates an existing group. All body fields are optional.").
		Returns(&dtos.GroupResponse{})

	deleteGroupById := validation.NewValidation(nil, nil, &dtos.GroupIdDto{})
	r.Delete("/:groupId", deleteGroupById, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupDelete),
		handlers.DeleteGroupById(service),
	).
		Summary("Delete group").
		Description("Deletes a group by its MongoDB ObjectId.").
		Returns(map[string]bool{})

	/** Group membership routes */

	getGroupMembers := validation.NewValidation(nil, &dtos.GroupMembersQueryDto{}, &dtos.GroupIdDto{})
	r.Get("/:groupId/members", getGroupMembers, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupRead),
		handlers.GetGroupMembers(service),
	).
		Summary("List group members").
		Description("Returns a paginated list of the users that belong to the group.").
		Returns(&model.PaginatedResult[dtos.GroupMemberResponse]{})

	addGroupMember := validation.NewValidation(&dtos.GroupMemberAddDto{}, nil, &dtos.GroupIdDto{})
	r.Post("/:groupId/members", addGroupMember, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupUpdate),
		handlers.AddGroupMember(service),
	).
		Summary("Add group member").
		Description("Adds a user to the group.").
		Returns(map[string]bool{})

	removeGroupMember := validation.NewValidation(nil, nil, &dtos.GroupMemberIdDto{})
	r.Delete("/:groupId/members/:userId", removeGroupMember, swagger.Expose,
		permissionMw.RequirePermission(perms.GroupUpdate),
		handlers.RemoveGroupMember(service),
	).
		Summary("Remove group member").
		Description("Removes a user from the group.").
		Returns(map[string]bool{})
}
