package handlers

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/plugins/application/dtos"
	"workflow/src/modules/plugins/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreatePlugin returns a Fiber handler that creates a new plugin manifest.
// Expects a validated bodyDTO of type *dtos.PluginManifestResponse from ValidationMiddleware.
// Uses RequestContext from coverage middleware for multi-tenant org population.
func CreatePlugin(service ports.PluginServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.PluginManifestResponse](c, "bodyDTO")

		retData, err := service.CreatePlugin(ctx, requestContext, bodyData)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetPluginById returns a Fiber handler that retrieves a plugin manifest by its MongoDB ObjectId.
// Expects a validated paramsDTO of type *dtos.PluginIdDTO from ValidationMiddleware.
func GetPluginById(service ports.PluginServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.PluginIdDTO](c, "paramsDTO")

		retData, err := service.GetPluginById(ctx, &params.PluginId)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdatePlugin returns a Fiber handler that updates a plugin manifest.
// Expects a validated bodyDTO of type *dtos.PluginManifestUpdate and
// paramsDTO of type *dtos.PluginIdDTO from ValidationMiddleware.
func UpdatePlugin(service ports.PluginServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.PluginIdDTO](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.PluginManifestUpdate](c, "bodyDTO")

		retData, err := service.UpdatePluginById(ctx, &params.PluginId, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeletePlugin returns a Fiber handler that deletes a plugin manifest.
// Expects a validated paramsDTO of type *dtos.PluginIdDTO from ValidationMiddleware.
func DeletePlugin(service ports.PluginServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.PluginIdDTO](c, "paramsDTO")

		retData, err := service.DeletePluginById(ctx, &params.PluginId)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetPlugins returns a Fiber handler that retrieves a paginated list of plugin manifests.
// Respects multi-tenant visibility: system + template (ancestor) + local (org).
// Expects a validated queryDTO of type *dtos.PluginQueryDTO from ValidationMiddleware.
// Uses RequestContext from coverage middleware for org filtering.
func GetPlugins(service ports.PluginServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.PluginQueryDTO](c, "queryDTO")

		retData, err := service.GetPlugins(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetEnabledPlugins returns a Fiber handler that retrieves all enabled plugins.
// Used by the frontend editor boot sequence. No validation needed.
func GetEnabledPlugins(service ports.PluginServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		retData, err := service.GetEnabledPlugins(ctx)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
