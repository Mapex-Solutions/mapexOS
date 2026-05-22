package routes

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/assets"
)

// RegisterRoutes registers all asset HTTP routes with the given router group.
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Asset service port interface implementation
func RegisterRoutes(group fiber.Router, service ports.AssetServicePort) {

	/**
	* CRUD Routes
	 */

	// Get assets with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Includes hierarchical support via PathKey data from coverage cache
	assetQueryDto := validation.NewValidation(nil, &dtos.AssetQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(assetQueryDto),
		permissionMw.RequirePermission(perms.AssetList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssets(service),
	)

	// Create a new asset
	// Uses InjectRequestContext middleware to populate PathKey from coverage cache
	assetCreateDto := validation.NewValidation(&dtos.AssetCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(assetCreateDto),
		permissionMw.RequirePermission(perms.AssetCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateAsset(service),
	)

	// Count assets (counter endpoint with cache)
	group.Get("/counter",
		permissionMw.RequirePermission(perms.AssetList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetCount(service),
	)

	// Get asset by ID
	getAssetById := validation.NewValidation(nil, nil, &dtos.AssetIdDto{})
	group.Get("/:assetId",
		validation.ValidationMiddleware(getAssetById),
		permissionMw.RequirePermission(perms.AssetRead),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetById(service),
	)

	// Update asset by ID
	updateAssetById := validation.NewValidation(&dtos.AssetUpdateDTO{}, nil, &dtos.AssetIdDto{})
	group.Patch("/:assetId",
		validation.ValidationMiddleware(updateAssetById),
		permissionMw.RequirePermission(perms.AssetUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateAssetById(service),
	)

	// Delete asset by ID
	deleteAssetById := validation.NewValidation(nil, nil, &dtos.AssetIdDto{})
	group.Delete("/:assetId",
		validation.ValidationMiddleware(deleteAssetById),
		permissionMw.RequirePermission(perms.AssetDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteAssetById(service),
	)

	// Generate a strong random alphanumeric MQTT password for the
	// operator to drop into the asset form. Stateless — it does not
	// touch any asset record. The operator may also type a custom
	// password instead; the platform validates only at create or
	// change-password time. Permission is AssetCreate so operators with
	// the create grant can call it before the asset itself exists.
	group.Get("/_generate_mqtt_password",
		permissionMw.RequirePermission(perms.AssetCreate),
		coverageMw.InjectRequestContext(),
		handlers.GenerateMqttPassword(service),
	)
}
