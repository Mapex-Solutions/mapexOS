package handlers

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/healthmonitor/application/dtos"
	"assets/src/modules/healthmonitor/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// ForceOffline returns a Fiber handler that transitions an asset to
// offline immediately, bypassing the scan + threshold cycle the
// scheduler enforces by default.
//
// The endpoint is addressed by device UUID; orgId is resolved
// server-side. The handler is idempotent — calling it for an asset
// already in the alerted (offline) state is a no-op.
//
// Used by e2e journeys to assert offline-action route group wiring
// without waiting the configured scan interval (default 600s on prod
// runtimes). Protected by the internal API key middleware applied on
// the parent route group.
func ForceOffline(service ports.HealthAdminPort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		var params dtos.AdminAssetUUIDDto
		if err := c.ParamsParser(&params); err != nil {
			return response.BadRequest(c, []string{err.Error()})
		}
		if params.AssetUUID == "" {
			return response.BadRequest(c, []string{"assetUUID is required"})
		}

		var body dtos.AdminForceOfflineRequestDto
		if len(c.Body()) > 0 {
			if err := c.BodyParser(&body); err != nil {
				return response.BadRequest(c, []string{err.Error()})
			}
		}

		if err := service.ForceOfflineByAssetUUID(ctx, params.AssetUUID, body.Reason); err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
