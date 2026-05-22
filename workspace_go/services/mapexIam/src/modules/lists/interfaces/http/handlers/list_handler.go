package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateList returns a Fiber handler that creates a new list.
//
// Lists support hierarchical multi-tenant architecture with isSystem flag for system/org-scoped resources.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
//
// It expects a validated DTO of type dtos.ListCreateDTO in the "bodyDTO" context key
// (populated by requestValidation middleware) containing:
//   - type: List type (assetGroup, assetType)
//   - name: List name
//   - value: List value
//   - isSystem: If true, creates system-level list (no org context)
//
// Multi-tenant fields (orgId, pathKey) are automatically populated by the service
// based on RequestContext.
//
// Returns:
//   - 201 Created with list data
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func CreateList(service ports.ListServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.ListCreateDTO](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateList(ctx, requestContext, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// GetListById returns a Fiber handler that retrieves a list by its unique identifier.
//
// It expects a validated DTO of type dtos.ListIdDTO in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with list data if found
//   - 404 Not Found if list doesn't exist
//   - 500 Internal Server Error on service failure
func GetListById(service ports.ListServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		list, _ := requestValidation.GetDTO[*dtos.ListIdDTO](c, "paramsDTO")
		retData, err := service.GetListById(ctx, &list.ListId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateListById returns a Fiber handler that updates an existing list's information.
//
// It expects validated DTOs:
//   - dtos.ListIdDTO in "paramsDTO" for the list identifier
//   - dtos.ListUpdateDTO in "bodyDTO" for the fields to update
//
// Both DTOs are populated by requestValidation middleware.
//
// Returns:
//   - 200 OK with updated list data
//   - 404 Not Found if list doesn't exist
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func UpdateListById(service ports.ListServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		list, _ := requestValidation.GetDTO[*dtos.ListIdDTO](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.ListUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateListById(ctx, &list.ListId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// DeleteListById returns a Fiber handler that permanently deletes a list.
//
// It expects a validated DTO of type dtos.ListIdDTO in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with success flag if deletion succeeds
//   - 404 Not Found if list doesn't exist
//   - 500 Internal Server Error on service failure
func DeleteListById(service ports.ListServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		list, _ := requestValidation.GetDTO[*dtos.ListIdDTO](c, "paramsDTO")
		retData, err := service.DeleteListById(ctx, &list.ListId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetLists returns a Fiber handler that retrieves a paginated and filtered list of lists.
// Uses scopedOrgIds from coverage middleware for multi-tenant access control.
//
// The handler extracts scopedOrgIds from c.Locals() which was set by the coverage middleware.
// This ensures users can only query lists within their accessible organizations.
// GetLists returns a Fiber handler that retrieves a paginated and filtered list of lists.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query lists within their accessible organizations with zero extra queries.
func GetLists(service ports.ListServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.ListQueryDTO](c, "queryDTO")
		retData, err := service.GetLists(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}
