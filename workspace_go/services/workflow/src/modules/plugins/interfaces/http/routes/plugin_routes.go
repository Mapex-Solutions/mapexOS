package routes

import (
	"workflow/src/modules/plugins/application/dtos"
	"workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/plugins/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers plugin HTTP routes.
//
// Base path: /api/v1/plugins
//
// HTTP Verbs follow REST conventions:
//
//	GET    /         - List plugins (paginated, filtered, multi-tenant)
//	GET    /enabled  - Get all enabled plugins (for editor boot)
//	POST   /         - Create plugin manifest
//	GET    /:id      - Get plugin by MongoDB ObjectId
//	PATCH  /:id      - Update plugin manifest
//	DELETE /:id      - Delete plugin manifest
func RegisterRoutes(group web.Router, service ports.PluginServicePort) {

	r := swagger.Wrap(group).Tag("Plugins")

	// List plugins with filters, pagination, and multi-tenant visibility.
	pluginQueryDto := validation.NewValidation(nil, &dtos.PluginQueryDTO{}, nil)
	r.Get("/", pluginQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.PluginList),
		coverageMw.InjectRequestContext(),
		handlers.GetPlugins(service),
	).
		Summary("List plugins").
		Description("Returns a paginated, filterable list of plugins visible to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.PluginManifestResponse]{})

	// Get all enabled plugins (for editor boot — no validation needed).
	enabledDto := validation.NewValidation(nil, nil, nil)
	r.Get("/enabled", enabledDto, swagger.Expose,
		permissionMw.RequirePermission(perms.PluginRead),
		handlers.GetEnabledPlugins(service),
	).
		Summary("List enabled plugins").
		Description("Returns every enabled plugin manifest, used to boot the workflow editor.").
		Returns(&[]dtos.PluginManifestResponse{})

	// Create a new plugin manifest.
	pluginCreateDto := validation.NewValidation(&dtos.PluginManifestResponse{}, nil, nil)
	r.Post("/", pluginCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.PluginCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreatePlugin(service),
	).
		Summary("Create plugin").
		Description("Registers a new plugin manifest.").
		Returns(&dtos.PluginManifestResponse{})

	// Get plugin by MongoDB ObjectId.
	pluginIdDto := validation.NewValidation(nil, nil, &dtos.PluginIdDTO{})
	r.Get("/:id", pluginIdDto, swagger.Expose,
		permissionMw.RequirePermission(perms.PluginRead),
		handlers.GetPluginById(service),
	).
		Summary("Get plugin by ID").
		Description("Retrieves a single plugin manifest by its MongoDB ObjectId.").
		Returns(&dtos.PluginManifestResponse{})

	// Update plugin manifest.
	pluginUpdateDto := validation.NewValidation(&dtos.PluginManifestUpdate{}, nil, &dtos.PluginIdDTO{})
	r.Patch("/:id", pluginUpdateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.PluginUpdate),
		handlers.UpdatePlugin(service),
	).
		Summary("Update plugin").
		Description("Partially updates an existing plugin manifest. Only provided fields are changed.").
		Returns(&dtos.PluginManifestResponse{})

	// Delete plugin manifest.
	pluginDeleteDto := validation.NewValidation(nil, nil, &dtos.PluginIdDTO{})
	r.Delete("/:id", pluginDeleteDto, swagger.Expose,
		permissionMw.RequirePermission(perms.PluginDelete),
		handlers.DeletePlugin(service),
	).
		Summary("Delete plugin").
		Description("Deletes a plugin manifest by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
