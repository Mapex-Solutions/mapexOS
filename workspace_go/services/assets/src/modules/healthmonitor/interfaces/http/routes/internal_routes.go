package routes

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/healthmonitor/application/ports"
	"assets/src/modules/healthmonitor/interfaces/http/handlers"
)

// RegisterInternalRoutes registers the healthmonitor internal HTTP
// routes. Routes are protected by the standard ApiKeyAuthMiddleware
// (X-API-Key header) applied by the caller on the parent group.
//
// The endpoint drives an asset to offline state immediately, bypassing
// the scheduled scan + threshold cycle. Used by e2e journeys to
// assert offline-action route group firing without waiting the
// configured scan interval.
//
// Endpoints:
//   - POST /:assetUUID/force-offline — Force-transition asset to offline.
func RegisterInternalRoutes(group fiber.Router, service ports.HealthAdminPort) {
	group.Post("/:assetUUID/force-offline", handlers.ForceOffline(service))
}
