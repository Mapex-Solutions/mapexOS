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

	/** Field Vocabulary Route */

	// Get the curated, multi-tenant field vocabulary for the authoring UI.
	// coverage middleware injects org filtering so org-scoped entries union
	// with the system standard. Read-only; lang selects hint/label locale.
	fieldVocabularyDto := validation.NewValidation(nil, &dtos.FieldVocabularyQuery{}, nil)
	r.Get("/field-vocabulary", fieldVocabularyDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateList),
		coverageMw.InjectRequestContext(),
		handlers.GetFieldVocabulary(service),
	).
		Summary("Get field vocabulary").
		Description("Returns the curated, multi-tenant list of canonical dynamic-field names for the asset-template authoring UI, grouped by category in a fixed order. Each field hint and group label are resolved to the requested language (en-US fallback).").
		Returns(&dtos.FieldVocabularyResponse{})

	/** Migration Routes */

	// Registered before the /:assetTemplateId wildcard routes so the static
	// /migrations prefix is not captured as an assetTemplateId. Within the group,
	// the static /migrations paths precede the /migrations/:id parameterized ones.

	// Create a migration plan. coverage middleware scopes the plan to the caller's org.
	createMigrationDto := validation.NewValidation(&dtos.MigrationPlanCreateRequest{}, nil, nil)
	r.Post("/migrations", createMigrationDto, swagger.Expose,
		permissionMw.RequirePermission(perms.TemplateMigrationCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateMigrationPlan(service),
	).
		Summary("Create migration plan").
		Description("Creates a scheduled template migration plan and one pending execution per affected asset. Organization scoping is applied automatically from the request context.").
		Returns(&dtos.MigrationPlanResponse{})

	// List migration plans. coverage middleware injects context-aware org filtering.
	listMigrationsDto := validation.NewValidation(nil, &dtos.MigrationPlanQueryDTO{}, nil)
	r.Get("/migrations", listMigrationsDto, swagger.Expose,
		permissionMw.RequirePermission(perms.TemplateMigrationList),
		coverageMw.InjectRequestContext(),
		handlers.ListMigrationPlans(service),
	).
		Summary("List migration plans").
		Description("Returns a paginated, filterable list of template migration plans scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.MigrationPlanResponse]{})

	// Get a migration plan by ID.
	getMigrationDto := validation.NewValidation(nil, nil, &dtos.MigrationPlanIdDto{})
	r.Get("/migrations/:id", getMigrationDto, swagger.Expose,
		permissionMw.RequirePermission(perms.TemplateMigrationRead),
		handlers.GetMigrationPlan(service),
	).
		Summary("Get migration plan by ID").
		Description("Retrieves a single template migration plan by its MongoDB ObjectId.").
		Returns(&dtos.MigrationPlanResponse{})

	// Update a migration plan by ID.
	updateMigrationDto := validation.NewValidation(&dtos.MigrationPlanUpdateRequest{}, nil, &dtos.MigrationPlanIdDto{})
	r.Patch("/migrations/:id", updateMigrationDto, swagger.Expose,
		permissionMw.RequirePermission(perms.TemplateMigrationUpdate),
		handlers.UpdateMigrationPlan(service),
	).
		Summary("Update migration plan").
		Description("Partially updates an editable migration plan and re-arms the start timer when the schedule changes. All body fields are optional; only provided fields are changed. Returns 409 when the plan is no longer editable.").
		Returns(&dtos.MigrationPlanResponse{})

	// Cancel a migration plan by ID.
	cancelMigrationDto := validation.NewValidation(nil, nil, &dtos.MigrationPlanIdDto{})
	r.Delete("/migrations/:id", cancelMigrationDto, swagger.Expose,
		permissionMw.RequirePermission(perms.TemplateMigrationDelete),
		handlers.CancelMigrationPlan(service),
	).
		Summary("Cancel migration plan").
		Description("Cancels a pending or scheduled migration plan by its MongoDB ObjectId. Returns 409 when the plan is already running or in a terminal state.").
		Returns(map[string]bool{})

	// List the per-asset executions of a migration plan.
	listMigrationExecutionsDto := validation.NewValidation(nil, &dtos.MigrationExecutionQueryDTO{}, &dtos.MigrationPlanIdDto{})
	r.Get("/migrations/:id/executions", listMigrationExecutionsDto, swagger.Expose,
		permissionMw.RequirePermission(perms.TemplateMigrationRead),
		handlers.ListMigrationExecutions(service),
	).
		Summary("List migration executions").
		Description("Returns a paginated, filterable list of the per-asset executions for a migration plan.").
		Returns(&model.PaginatedResult[dtos.MigrationExecutionResponse]{})

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

	// Install/uninstall a marketplace template into the caller's org. Registered
	// before the /:assetTemplateId wildcard so a three-segment install path can
	// never be captured as a template id.
	installDto := validation.NewValidation(&dtos.InstallBody{}, nil, &dtos.InstallParams{})
	r.Post("/:marketplaceVendor/:marketplaceSlug/install", installDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateCreate),
		coverageMw.InjectRequestContext(),
		handlers.InstallFromMarketplace(service),
	).
		Summary("Install a marketplace asset template").
		Description("Fetches a template from the mapexMarketplace catalog, hard-verifies its sha256, resolves org-scoped classification, caches the shared content, and creates the caller organization's per-org link. Rejects with 422 on a checksum mismatch.").
		Returns(&dtos.AssetTemplateResponse{})

	r.Delete("/:marketplaceVendor/:marketplaceSlug/install", installDto, swagger.Expose,
		permissionMw.RequirePermission(perms.AssetTemplateDelete),
		coverageMw.InjectRequestContext(),
		handlers.UninstallFromMarketplace(service),
	).
		Summary("Uninstall a marketplace asset template").
		Description("Removes the caller organization's installation of a marketplace template. The shared content and other organizations' installs are untouched. Returns 404 when the organization has not installed it.").
		Returns(map[string]bool{})

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
