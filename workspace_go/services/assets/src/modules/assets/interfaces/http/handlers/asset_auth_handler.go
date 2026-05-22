package handlers

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assets/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetAssetAuth returns a Fiber handler that retrieves the slim auth
// projection for an asset by UUID — the L3 fallback the broker
// plugin hits when L1 (Pebble) and L2 (MinIO mapex-asset-auth) both
// miss. The handler also spawns an async warm-up writing the
// projection back to MinIO so the next broker lookup hits L2.
//
// Security: Protected by API Key authentication (X-API-Key header)
// applied at the parent group.
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler that returns AuthProjection JSON or 404.
func GetAssetAuth(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		assetUUID := c.Params("assetUUID")
		if assetUUID == "" {
			return response.BadRequest(c, []string{"assetUUID is required"})
		}

		projection, err := service.GetAuthProjectionByUUID(ctx, assetUUID)
		if err != nil {
			return err
		}

		return response.Success(c, projection)
	}
}
