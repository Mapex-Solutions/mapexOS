package routes

import (
	"events/src/modules/asset_status/application/dtos"
	"events/src/modules/asset_status/application/ports"
	"events/src/modules/asset_status/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/events"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes mounts the asset_status HTTP endpoints under the parent group
// (/api/v1/events as set by module.go).
//
//	GET /connectivity_history                    (org-wide list)
//	GET /assets/:assetUUID/connectivity_history  (asset-scoped)
//
// Permission `events.asset_status.list` gates both; org access is enforced by the
// coverage middleware. Routes are registered through the swagger wrapper.
func RegisterRoutes(group web.Router, service ports.AssetStatusServicePort) {

	r := swagger.Wrap(group).Tag("Asset Status")

	queryDto := validation.NewValidation(nil, &dtos.AssetConnectivityHistoryQuery{}, nil)

	// Org-wide connectivity history.
	r.Get("/connectivity_history", queryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsAssetStatusList),
		coverageMw.InjectRequestContext(),
		handlers.ListConnectivityHistory(service),
	).
		Summary("List connectivity history").
		Description("Returns cursor-paginated asset connectivity (online/offline) transitions across the caller's organization.").
		Returns(&dtos.AssetConnectivityCursorResult{})

	// Connectivity history for a single asset.
	r.Get("/assets/:assetUUID/connectivity_history", queryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsAssetStatusList),
		coverageMw.InjectRequestContext(),
		handlers.GetConnectivityHistory(service),
	).
		Summary("Get asset connectivity history").
		Description("Returns cursor-paginated connectivity (online/offline) transitions for a single asset identified by assetUUID.").
		Returns(&dtos.AssetConnectivityCursorResult{})
}
