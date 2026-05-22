package routes

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/interfaces/http/handlers"
)

// RegisterInternalRoutes registers internal HTTP routes with API key authentication.
//
// The single endpoint here is the L3 read-model fallback for every
// consumer that caches the AssetReadModel (Router, JS-Executor, Events,
// mapex-mqtt-broker plugin). When the consumer's local cache (L1) and
// the shared L2 (MinIO) both miss, it falls back to this endpoint —
// the handler also repopulates L2 on the way out so the next read hits
// the cache.
//
// The broker plugin uses the AssetReadModel's `Protocol.Mqtt.PasswordHash`
// and `CurrentCert.Serial` to decide MQTT CONNECTs locally (bcrypt
// compare for password mode, serial-equality for cert mode). There is
// no separate MQTT auth callout — auth is always plugin-local.
//
// Security: Routes are protected by the standard ApiKeyAuthMiddleware
// (X-API-Key header), applied by the caller on the parent group.
//
// Endpoints:
//   - GET /internal/assets/:assetUUID - Fetch asset read model by UUID
//
// Parameters:
//   - group: Fiber router group to register routes on (should have apiKey middleware applied)
//   - service: Asset service port interface implementation
func RegisterInternalRoutes(group fiber.Router, service ports.AssetServicePort) {
	group.Get("/:assetUUID", handlers.GetAssetReadModel(service))
}
