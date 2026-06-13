package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/interfaces/http/handlers"
)

// RegisterAssetAuthInternalRoutes registers the broker-plugin L3
// fallback endpoint for the slim auth projection. Lives on its own
// group `/internal/asset_auth` (sibling of `/internal/assets`) so
// the broker's auth path is decoupled from the full read-model
// endpoint and the bucket the projection lives in.
//
// Endpoints:
//   - GET /internal/asset_auth/:assetUUID — Fetch AuthProjection by UUID
//
// Parameters:
//   - group: Fiber router group to register routes on (should have
//     apiKey middleware applied)
//   - service: Asset service port interface implementation
func RegisterAssetAuthInternalRoutes(group web.Router, service ports.AssetServicePort) {
	group.Get("/:assetUUID", handlers.GetAssetAuth(service))
}
