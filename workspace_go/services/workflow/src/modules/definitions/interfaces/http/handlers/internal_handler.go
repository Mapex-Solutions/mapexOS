package handlers

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/definitions/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetNodeScript returns a Fiber handler that retrieves a code node's script source.
//
// This is an internal endpoint used by TieredCache consumers (js-workflow-executor)
// as a fallback when L2 (MinIO) cache miss occurs. The endpoint fetches the definition
// from MongoDB, extracts the script, repopulates L2, and returns the script source.
//
// Security: Protected by API Key authentication (X-API-Key header)
func GetNodeScript(service ports.DefinitionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		definitionId := c.Params("definitionId")
		if definitionId == "" {
			return response.BadRequest(c, []string{"definitionId is required"})
		}

		nodeId := c.Params("nodeId")
		if nodeId == "" {
			return response.BadRequest(c, []string{"nodeId is required"})
		}

		script, err := service.GetNodeScript(ctx, definitionId, nodeId)
		if err != nil {
			return err
		}

		return response.Success(c, script)
	}
}
