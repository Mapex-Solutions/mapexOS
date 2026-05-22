package routes

import (
	"github.com/gofiber/fiber/v2"

	"events/src/modules/asset_status/application/dtos"
	"events/src/modules/asset_status/application/ports"
	"events/src/modules/asset_status/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/events"
)

// RegisterRoutes mounts the asset_status HTTP endpoints under the parent
// group (/api/v1/events as set by module.go).
//
// Two endpoints:
//   - GET /connectivity_history                         (org-wide list)
//   - GET /assets/:assetUUID/connectivity_history       (asset-scoped)
//
// Permission `events.asset_status.list` gates both. Org access is enforced
// by the coverage middleware (InjectRequestContext).
func RegisterRoutes(group fiber.Router, service ports.AssetStatusServicePort) {
	queryDto := validation.NewValidation(nil, &dtos.AssetConnectivityHistoryQuery{}, nil)

	group.Get("/connectivity_history",
		validation.ValidationMiddleware(queryDto),
		permissionMw.RequirePermission(perms.EventsAssetStatusList),
		coverageMw.InjectRequestContext(),
		handlers.ListConnectivityHistory(service),
	)

	group.Get("/assets/:assetUUID/connectivity_history",
		validation.ValidationMiddleware(queryDto),
		permissionMw.RequirePermission(perms.EventsAssetStatusList),
		coverageMw.InjectRequestContext(),
		handlers.GetConnectivityHistory(service),
	)
}
