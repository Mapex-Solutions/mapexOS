package routes

import (
	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	"assets/src/modules/mqttcerts/interfaces/http/handlers"
	localMw "assets/src/modules/mqttcerts/interfaces/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes mounts /api/v1/mqtt_certs/*. The router passed in
// MUST already be JWT-gated + coverage-injected by the caller
// (module.go composes that group). Permissions wrapping (MqttCertCreate
// etc.) is added per-route below — the caller need not pre-apply.
func RegisterRoutes(router fiber.Router, h *handlers.MqttCertsHandler, svc mqttPorts.MqttCertsServicePort) {
	g := router.Group("/", localMw.RequireCAReady(svc))
	g.Post("/", h.IssueCert)
	g.Delete("/:serial", h.RevokeCert)
	g.Get("/", h.ListByAsset)
}
