package handlers

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/instances/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetInstanceCount returns a Fiber handler that returns the total count of workflow instances.
func GetInstanceCount(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountInstances(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}
