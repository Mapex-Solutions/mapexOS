package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/application/ports"
	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateGroup returns a Fiber handler that creates a new group.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
//
// It expects a validated DTO of type dtos.CreateGroupDto in the "bodyDTO" context key
// (populated by requestValidation middleware) containing:
//   - name: Group name
//   - description: Optional description
//   - enabled: Boolean to enable/disable the group
//   - roleIds: Array of role IDs for the group membership
//
// Multi-tenant fields (orgId, pathKey, scope) are automatically populated by the service
// based on RequestContext.
//
// Returns:
//   - 201 Created with group data
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func CreateGroup(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.CreateGroupDto](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateGroup(ctx, requestContext, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// GetGroupById returns a Fiber handler that retrieves a group by its unique identifier.
//
// It expects a validated DTO of type dtos.GroupIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with group data if found
//   - 404 Not Found if group doesn't exist
//   - 500 Internal Server Error on service failure
func GetGroupById(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.GroupIdDto](c, "paramsDTO")
		retData, err := service.GetGroupById(ctx, &params.GroupId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateGroupById returns a Fiber handler that updates an existing group's information.
//
// It expects validated DTOs:
//   - dtos.GroupIdDto in "paramsDTO" for the group identifier
//   - dtos.UpdateGroupDto in "bodyDTO" for the fields to update
//
// Both DTOs are populated by requestValidation middleware.
//
// Returns:
//   - 200 OK with updated group data
//   - 404 Not Found if group doesn't exist
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func UpdateGroupById(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.GroupIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.UpdateGroupDto](c, "bodyDTO")
		retData, err := service.UpdateGroupById(ctx, &params.GroupId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// DeleteGroupById returns a Fiber handler that permanently deletes a group.
//
// It expects a validated DTO of type dtos.GroupIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with success flag if deletion succeeds
//   - 404 Not Found if group doesn't exist
//   - 500 Internal Server Error on service failure
func DeleteGroupById(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.GroupIdDto](c, "paramsDTO")
		retData, err := service.DeleteGroupById(ctx, &params.GroupId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetGroups returns a Fiber handler that retrieves a paginated and filtered list of groups.
// Uses RequestContext from coverage middleware for context-aware org filtering.
func GetGroups(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.GroupQueryDto](c, "queryDTO")
		retData, err := service.GetGroups(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetGroupMembers returns a Fiber handler that retrieves paginated members of a group.
//
// It expects validated DTOs:
//   - dtos.GroupIdDto in "paramsDTO" for the group identifier
//   - dtos.GroupMembersQueryDto in "queryDTO" for pagination (page, perPage - max 100)
//
// Returns:
//   - 200 OK with paginated member list
//   - 404 Not Found if group doesn't exist
//   - 500 Internal Server Error on service failure
func GetGroupMembers(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.GroupIdDto](c, "paramsDTO")
		queryData, _ := requestValidation.GetDTO[*dtos.GroupMembersQueryDto](c, "queryDTO")

		retData, err := service.GetGroupMembers(ctx, params.GroupId, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// AddGroupMember returns a Fiber handler that adds a user to a group.
//
// It expects validated DTOs:
//   - dtos.GroupIdDto in "paramsDTO" for the group identifier
//   - dtos.GroupMemberAddDto in "bodyDTO" for the user to add
//
// Returns:
//   - 201 Created on success
//   - 404 Not Found if group doesn't exist
//   - 500 Internal Server Error on service failure
func AddGroupMember(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.GroupIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.GroupMemberAddDto](c, "bodyDTO")

		err := service.AddMemberToGroup(ctx, params.GroupId, bodyData.UserID)
		if err != nil {
			return err
		}
		return response.Created(c, map[string]bool{"success": true})
	}
}

// RemoveGroupMember returns a Fiber handler that removes a user from a group.
//
// It expects validated DTO:
//   - dtos.GroupMemberIdDto in "paramsDTO" for groupId and userId
//
// Returns:
//   - 200 OK on success
//   - 404 Not Found if group doesn't exist
//   - 500 Internal Server Error on service failure
func RemoveGroupMember(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.GroupMemberIdDto](c, "paramsDTO")

		err := service.RemoveMemberFromGroup(ctx, params.GroupId, params.UserId)
		if err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// GetGroupCount returns a Fiber handler that returns the total count of groups.
// Uses cached count with 6h TTL, invalidated on create/delete.
//
// Parameters:
//   - service: The GroupServicePort interface for group business operations
//
// Returns:
//   - A Fiber handler function that processes the group count request
func GetGroupCount(service ports.GroupServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountGroups(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}
