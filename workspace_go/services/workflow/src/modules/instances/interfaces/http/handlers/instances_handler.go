package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/instances/application/dtos"
	"workflow/src/modules/instances/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// GetInstances returns a Fiber handler that retrieves a paginated list of workflow instance configs.
func GetInstances(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.InstanceQueryDTO](c, "queryDTO")
		retData, err := service.GetInstances(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetInstanceById returns a Fiber handler that retrieves a workflow instance config by its ID.
func GetInstanceById(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.InstanceIdDTO](c, "paramsDTO")
		retData, err := service.GetInstanceById(ctx, &params.InstanceId)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// CreateInstance returns a Fiber handler that creates a new workflow instance config.
func CreateInstance(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.InstanceCreateDTO](c, "bodyDTO")
		retData, err := service.CreateInstance(ctx, requestContext, bodyData)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// UpdateInstanceById returns a Fiber handler that updates a workflow instance config.
func UpdateInstanceById(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.InstanceIdDTO](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.InstanceUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateInstanceById(ctx, &params.InstanceId, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteInstanceById returns a Fiber handler that deletes a workflow instance config.
func DeleteInstanceById(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.InstanceIdDTO](c, "paramsDTO")
		err := service.DeleteInstanceById(ctx, &params.InstanceId)
		if err != nil {
			return err
		}
		return response.Success(c, nil)
	}
}

// ExecuteInstance returns a Fiber handler that executes a workflow instance.
func ExecuteInstance(service ports.InstancesServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.InstanceIdDTO](c, "paramsDTO")

		// Parse optional body
		var body dtos.ExecuteRequestDTO
		_ = c.BodyParser(&body)

		result, err := service.ExecuteInstance(ctx, params.InstanceId, body.EventPayload, body.WorkflowUUID)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not found") {
				return response.NotFound(c, fmt.Errorf("%s", errMsg))
			}
			if strings.Contains(errMsg, "disabled") {
				return response.BadRequest(c, []string{errMsg})
			}
			if strings.Contains(errMsg, "already running") {
				return response.Conflict(c, []string{errMsg})
			}
			return response.InternalServerError(c, errMsg, err)
		}
		return response.Success(c, result)
	}
}
