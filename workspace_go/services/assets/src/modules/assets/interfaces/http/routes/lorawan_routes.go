package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/interfaces/http/handlers"
)

// RegisterLorawanInternalRoutes registers the LNS-facing LoRaWAN resolution
// endpoints. Mounted on its own group `/internal/lorawan` (sibling of
// `/internal/assets` and `/internal/asset_auth`) so the LNS resolves a device by the
// DevAddr an uplink carries.
//
//   - GET /internal/lorawan/by-devaddr/:devAddr — resolve devices by DevAddr
//
// Parameters:
//   - group: Router group for internal routes (API-key middleware applied by caller)
//   - service: The AssetServicePort for asset business operations
func RegisterLorawanInternalRoutes(group web.Router, service ports.AssetServicePort) {
	group.Get("/by-devaddr/:devAddr", handlers.GetLorawanByDevAddr(service))
}
