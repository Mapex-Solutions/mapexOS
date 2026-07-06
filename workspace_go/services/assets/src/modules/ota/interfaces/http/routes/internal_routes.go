package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/interfaces/http/handlers"
)

// RegisterInternalRoutes registers the OTA internal HTTP routes. Routes are
// protected by the standard ApiKeyAuthMiddleware (X-API-Key header) applied by
// the caller on the parent group.
//
// The HTTP gateway (the owner of device data-source auth) relays the device
// poll here — the Asset MS never hosts device-facing auth itself.
//
// Endpoints:
//   - GET /jobs?assetUUID={assetUUID} — the asset's pending OTA command with a
//     fresh presigned download URL, or data=null when none.
func RegisterInternalRoutes(group web.Router, jobs ports.DeviceJobPort) {
	group.Get("/jobs", handlers.GetPendingJob(jobs))
}
