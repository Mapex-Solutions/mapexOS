package handlers

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assettemplates/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetTemplateForCacheFallback returns a Fiber handler that retrieves a template by ID
// and repopulates L2 (MinIO) cache.
//
// This endpoint is used internally by other services (JS-Executor, Events) when
// TieredCache L2 miss occurs. It fetches the template from MongoDB (source of truth)
// and repopulates the L2 cache for future requests.
//
// Security: Protected by API Key authentication (internal service-to-service only).
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the cache fallback request
func GetTemplateForCacheFallback(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		templateId := c.Params("templateId")
		if templateId == "" {
			return response.BadRequest(c, []string{"templateId is required"})
		}

		retData, err := service.GetTemplateByIdForCacheFallback(ctx, templateId)
		if err != nil {
			return err
		}

		return response.Success(c, retData)
	}
}
