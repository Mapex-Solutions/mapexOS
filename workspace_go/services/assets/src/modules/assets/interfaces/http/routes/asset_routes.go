package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/assets"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
)

// RegisterRoutes registers all asset HTTP routes with the given router group.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/assets
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Asset service port interface implementation
func RegisterRoutes(group web.Router, service ports.AssetServicePort) {

	r := swagger.Wrap(group).Tag("Assets")

	/** List Routes */

	// Get assets with filters, pagination, and projection. coverage middleware
	// injects context-aware org filtering (hierarchical via PathKey).
	assetQueryDto := validation.NewValidation(nil, &dtos.AssetQueryDTO{}, nil)
	r.Get("/", assetQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssets(service),
	).
		Summary("List assets").
		Description("Returns a paginated, filterable list of assets scoped to the caller's organization. Supports hierarchical queries via includeChildren.").
		Returns(&model.PaginatedResult[dtos.AssetResponse]{})

	/** Counter Route */

	// Count assets (cached). No request contract to validate.
	counterDto := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", counterDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetCount(service),
	).
		Summary("Count assets").
		Description("Returns the total number of assets for the caller's organization, using a cached count.").
		Returns(&contractsCommon.CounterResponse{})

	/** CRUD Routes */

	// Create a new asset. coverage middleware populates PathKey from the cache.
	assetCreateDto := validation.NewValidation(&dtos.AssetCreateDTO{}, nil, nil)
	r.Post("/", assetCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateAsset(service),
	).
		Summary("Create asset").
		Description("Creates a new asset. Organization scoping (orgId, pathKey) is applied automatically from the request context.").
		Returns(&dtos.AssetResponse{})

	// Get asset by ID.
	getAssetById := validation.NewValidation(nil, nil, &dtos.AssetIdDto{})
	r.Get("/:assetId", getAssetById, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetRead),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetById(service),
	).
		Summary("Get asset by ID").
		Description("Retrieves a single asset by its MongoDB ObjectId.").
		Returns(&dtos.AssetResponse{})

	// Update asset by ID.
	updateAssetById := validation.NewValidation(&dtos.AssetUpdateDTO{}, nil, &dtos.AssetIdDto{})
	r.Patch("/:assetId", updateAssetById, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateAssetById(service),
	).
		Summary("Update asset").
		Description("Partially updates an existing asset. All body fields are optional; only provided fields are changed.").
		Returns(&dtos.AssetResponse{})

	// Delete asset by ID.
	deleteAssetById := validation.NewValidation(nil, nil, &dtos.AssetIdDto{})
	r.Delete("/:assetId", deleteAssetById, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteAssetById(service),
	).
		Summary("Delete asset").
		Description("Deletes an asset by its MongoDB ObjectId.").
		Returns(map[string]bool{})

	/** Utility Routes */

	// Generate a strong random alphanumeric MQTT password for the operator to
	// drop into the asset form. Stateless — it does not touch any asset record.
	// The operator may also type a custom password instead; the platform
	// validates only at create or change-password time. Permission is AssetCreate
	// so operators with the create grant can call it before the asset exists.
	generateMqttPasswordDto := validation.NewValidation(nil, nil, nil)
	r.Get("/_generate_mqtt_password", generateMqttPasswordDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetCreate),
		coverageMw.InjectRequestContext(),
		handlers.GenerateMqttPassword(service),
	).
		Summary("Generate MQTT password").
		Description("Returns a strong random alphanumeric MQTT password suggestion for the asset form. Stateless — no asset record is created or modified.").
		Returns(&dtos.GenerateMqttPasswordResponseDTO{})
}
