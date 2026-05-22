package handlers

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/archiver/application/dtos"
	"workflow/src/modules/archiver/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	requestValidation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetExecutions returns a Fiber handler that retrieves workflow executions with pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
func GetExecutions(service ports.ArchiverServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.ExecutionQueryDTO](c, "queryDTO")
		retData, err := service.GetExecutions(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetExecutionById returns a Fiber handler that retrieves a single execution by ID.
func GetExecutionById(service ports.ArchiverServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		paramsData, _ := requestValidation.GetDTO[*dtos.ExecutionIdDTO](c, "paramsDTO")
		retData, err := service.GetExecutionById(ctx, paramsData.ExecutionId)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
