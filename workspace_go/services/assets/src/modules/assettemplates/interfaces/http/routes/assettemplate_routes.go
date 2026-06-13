package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/assets"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
)

// RegisterRoutes registers asset template HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/asset_templates
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Asset template service port interface implementation
func RegisterRoutes(group web.Router, service ports.AssetTemplateServicePort) {

	r := swagger.Wrap(group).Tag("Asset Templates")

	/** List Routes */

	// List asset templates with pagination and filters. coverage middleware
	// injects context-aware org filtering (hierarchical via PathKey).
	listDto := validation.NewValidation(nil, &dtos.AssetTemplateQueryDto{}, nil)
	r.Get("/", listDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetTemplates(service),
	).
		Summary("List asset templates").
		Description("Returns a paginated, filterable list of asset templates scoped to the caller's organization. Supports hierarchical queries via includeChildren.").
		Returns(&model.PaginatedResult[dtos.AssetTemplateResponse]{})

	/** Counter Route */

	// Count asset templates (cached). No request contract to validate.
	counterDto := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", counterDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetTemplateCount(service),
	).
		Summary("Count asset templates").
		Description("Returns the total number of asset templates for the caller's organization, using a cached count.").
		Returns(&contractsCommon.CounterResponse{})

	/** CRUD Routes */

	// Create a new asset template. coverage middleware populates orgId and pathKey.
	createDto := validation.NewValidation(&dtos.AssetTemplateCreateDTO{}, nil, nil)
	r.Post("/", createDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateAssetTemplate(service),
	).
		Summary("Create asset template").
		Description("Creates a new asset template, including its decode/validate/transform scripts. Organization scoping is applied automatically from the request context.").
		Returns(&dtos.AssetTemplateResponse{})

	// Get asset template by ID.
	getAssetTemplateById := validation.NewValidation(nil, nil, &dtos.AssetTemplateIdDto{})
	r.Get("/:assetTemplateId", getAssetTemplateById, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateRead),
		handlers.GetAssetTemplateById(service),
	).
		Summary("Get asset template by ID").
		Description("Retrieves a single asset template by its MongoDB ObjectId.").
		Returns(&dtos.AssetTemplateResponse{})

	// Update asset template by ID.
	updateAssetTemplateById := validation.NewValidation(&dtos.AssetTemplateUpdateDTO{}, nil, &dtos.AssetTemplateIdDto{})
	r.Patch("/:assetTemplateId", updateAssetTemplateById, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateUpdate),
		handlers.UpdateAssetTemplateById(service),
	).
		Summary("Update asset template").
		Description("Partially updates an existing asset template. All body fields are optional; only provided fields are changed.").
		Returns(&dtos.AssetTemplateResponse{})

	// Delete asset template by ID.
	deleteAssetTemplateById := validation.NewValidation(nil, nil, &dtos.AssetTemplateIdDto{})
	r.Delete("/:assetTemplateId", deleteAssetTemplateById, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateDelete),
		handlers.DeleteAssetTemplateById(service),
	).
		Summary("Delete asset template").
		Description("Deletes an asset template by its MongoDB ObjectId.").
		Returns(map[string]bool{})

	/** Special Routes */

	// Get available fields for an asset template (for Rule autocomplete). Backed
	// by a Redis cache (24h TTL), invalidated when the template is created or updated.
	getAvailableFields := validation.NewValidation(nil, nil, &dtos.AssetTemplateIdDto{})
	r.Get("/:assetTemplateId/available_fields", getAvailableFields, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateRead),
		coverageMw.InjectRequestContext(),
		handlers.GetAvailableFields(service),
	).
		Summary("Get available fields").
		Description("Returns the available output fields produced by the template's scripts, used for rule-builder autocomplete. Cached with a 24h TTL.").
		Returns(map[string]interface{}{})
}
