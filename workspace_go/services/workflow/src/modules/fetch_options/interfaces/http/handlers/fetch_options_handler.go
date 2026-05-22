package handlers

import (
	"workflow/src/modules/fetch_options/application/ports"

	"github.com/gofiber/fiber/v2"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// FetchOptions returns a Fiber handler for the fetchOptions proxy endpoint.
func FetchOptions(service ports.FetchOptionsServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		_, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}

		var body FetchOptionsRequest
		if err := c.BodyParser(&body); err != nil {
			return response.BadRequest(c, []string{"Invalid request body"})
		}

		if body.CredentialId == "" || body.ResourceKey == "" {
			return response.BadRequest(c, []string{"credentialId and resourceKey are required"})
		}

		items, err := service.FetchOptions(ctx, body.CredentialId, body.PluginId, body.ResourceKey, body.DependsOn)
		if err != nil {
			return err
		}
		return response.Success(c, items)
	}
}
