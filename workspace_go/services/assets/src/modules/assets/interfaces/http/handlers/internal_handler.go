package handlers

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assets/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetAssetReadModel returns a Fiber handler that retrieves an asset read model by UUID.
//
// This is the internal L3 fallback endpoint used by every consumer that
// caches AssetReadModel (Router, JS-Executor, Events, mapex-mqtt-broker
// plugin). The endpoint fetches the asset from MongoDB, writes it to
// MinIO (L2) for future requests, and returns the read model.
//
// The broker plugin reads `PasswordHash` + `CurrentCert` off the
// returned AssetReadModel to decide MQTT CONNECTs locally (bcrypt
// compare or cert serial match) — there is NO separate auth callout.
//
// Security: Protected by API Key authentication (X-API-Key header)
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler function that processes the asset read model request
func GetAssetReadModel(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get assetUUID from URL path parameter
		assetUUID := c.Params("assetUUID")
		if assetUUID == "" {
			return response.BadRequest(c, []string{"assetUUID is required"})
		}

		// Fetch read model from service (which handles MongoDB lookup and MinIO write)
		readModel, err := service.GetAssetReadModelByUUID(ctx, assetUUID)
		if err != nil {
			return err
		}

		return response.Success(c, readModel)
	}
}
