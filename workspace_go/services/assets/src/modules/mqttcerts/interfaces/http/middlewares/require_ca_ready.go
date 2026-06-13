package middlewares

import (
	mqttPorts "assets/src/modules/mqttcerts/application/ports"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RequireCAReady returns 502 Bad Gateway when the in-RAM CA store is
// empty. Cheap atomic load — safe to use on the hot path.
func RequireCAReady(svc mqttPorts.MqttCertsServicePort) web.Handler {
	return func(c *web.Ctx) error {
		if !svc.IsCAReady() {
			return c.Status(web.StatusBadGateway).JSON(web.Map{
				"error":   "ca_not_ready",
				"message": "PKI subsystem not ready; mapexVault unreachable or bootstrap pending. Retry shortly.",
			})
		}
		return c.Next()
	}
}
