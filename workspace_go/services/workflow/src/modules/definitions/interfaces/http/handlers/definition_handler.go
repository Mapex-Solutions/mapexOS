package handlers

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/definitions/application/dtos"
	"workflow/src/modules/definitions/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateDefinition returns a Fiber handler that creates a new workflow definition.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//   - UserId: The user ID from JWT token
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
func CreateDefinition(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.DefinitionCreateDTO](c, "bodyDTO")
		retData, err := service.CreateDefinition(ctx, requestContext, bodyData)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetDefinitionById returns a Fiber handler that retrieves a workflow definition by its ID.
func GetDefinitionById(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.DefinitionIdDTO](c, "paramsDTO")
		retData, err := service.GetDefinitionById(ctx, &params.WorkflowId)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateDefinitionById returns a Fiber handler that updates a workflow definition.
func UpdateDefinitionById(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.DefinitionIdDTO](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.DefinitionUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateDefinitionById(ctx, &params.WorkflowId, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteDefinitionById returns a Fiber handler that deletes a workflow definition.
func DeleteDefinitionById(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.DefinitionIdDTO](c, "paramsDTO")
		retData, err := service.DeleteDefinitionById(ctx, &params.WorkflowId)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetDefinitionCount returns a Fiber handler that returns the total count of workflow definitions.
func GetDefinitionCount(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountDefinitions(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}

// GetDefinitions returns a Fiber handler that retrieves a paginated list of workflow definitions.
// Uses RequestContext from coverage middleware for org filtering with hierarchical support.
func GetDefinitions(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.DefinitionQueryDTO](c, "queryDTO")
		retData, err := service.GetDefinitions(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
