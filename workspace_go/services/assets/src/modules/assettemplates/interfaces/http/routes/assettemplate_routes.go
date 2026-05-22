package routes

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/assets"
)

// RegisterRoutes registers asset template HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/asset_templates
//
// HTTP Verbs follow REST conventions:
//
//	GET    /                      - List asset templates (paginated, filtered)
//	POST   /                      - Create asset template
//	GET    /:assetTemplateId      - Get asset template by ID
//	PATCH  /:assetTemplateId      - Update asset template
//	DELETE /:assetTemplateId      - Delete asset template
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Asset template service port interface implementation
func RegisterRoutes(group fiber.Router, service ports.AssetTemplateServicePort) {

	/**
	* List Routes
	 */

	// List asset templates with pagination and filters
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Includes hierarchical support via PathKey data from coverage cache
	listDto := validation.NewValidation(nil, &dtos.AssetTemplateQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(listDto),                // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.AssetTemplateList), // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                       // 3. Inject context (cache)
		handlers.GetAssetTemplates(service),                     // 4. Handler
	)

	/**
	* CRUD Routes
	 */

	// Create a new assettemplate
	// Uses InjectRequestContext middleware to populate orgId and pathKey
	dataSourceCreateDto := validation.NewValidation(&dtos.AssetTemplateCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(dataSourceCreateDto),      // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.AssetTemplateCreate), // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                         // 3. Inject context (cache)
		handlers.CreateAssetTemplate(service),                     // 4. Handler
	)

	// Count asset templates (counter endpoint with cache)
	group.Get("/counter",
		permissionMw.RequirePermission(perms.AssetTemplateList),
		coverageMw.InjectRequestContext(),
		handlers.GetAssetTemplateCount(service),
	)

	// Get assettemplate by ID
	getAssetTemplateById := validation.NewValidation(nil, nil, &dtos.AssetTemplateIdDto{})
	group.Get("/:assetTemplateId",
		validation.ValidationMiddleware(getAssetTemplateById),
		permissionMw.RequirePermission(perms.AssetTemplateRead),
		handlers.GetAssetTemplateById(service),
	)

	// Update assettemplate by ID
	updateAssetTemplateById := validation.NewValidation(&dtos.AssetTemplateUpdateDTO{}, nil, &dtos.AssetTemplateIdDto{})
	group.Patch("/:assetTemplateId",
		validation.ValidationMiddleware(updateAssetTemplateById),
		permissionMw.RequirePermission(perms.AssetTemplateUpdate),
		handlers.UpdateAssetTemplateById(service),
	)

	// Delete assettemplate by ID
	deleteAssetTemplateById := validation.NewValidation(nil, nil, &dtos.AssetTemplateIdDto{})
	group.Delete("/:assetTemplateId",
		validation.ValidationMiddleware(deleteAssetTemplateById),
		permissionMw.RequirePermission(perms.AssetTemplateDelete),
		handlers.DeleteAssetTemplateById(service),
	)

	/**
	* Special Routes
	 */

	// Get available fields for an asset template (for Rule autocomplete)
	// Uses Redis cache with 24-hour TTL for performance
	// Cache is invalidated when template is created or updated
	getAvailableFields := validation.NewValidation(nil, nil, &dtos.AssetTemplateIdDto{})
	group.Get("/:assetTemplateId/available_fields",
		validation.ValidationMiddleware(getAvailableFields),
		permissionMw.RequirePermission(perms.AssetTemplateRead),
		coverageMw.InjectRequestContext(),
		handlers.GetAvailableFields(service),
	)
}
