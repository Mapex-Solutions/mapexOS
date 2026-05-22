package middlewares

import (
	mqttPorts "assets/src/modules/mqttcerts/application/ports"

	"github.com/gofiber/fiber/v2"
)

// RequireCAReady returns 502 Bad Gateway when the in-RAM CA store is
// empty. Cheap atomic load — safe to use on the hot path.
func RequireCAReady(svc mqttPorts.MqttCertsServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !svc.IsCAReady() {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "ca_not_ready",
				"message": "PKI subsystem not ready; mapexVault unreachable or bootstrap pending. Retry shortly.",
			})
		}
		return c.Next()
	}
}
