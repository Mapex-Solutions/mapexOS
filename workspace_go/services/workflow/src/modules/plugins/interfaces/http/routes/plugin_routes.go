package routes

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/plugins/application/dtos"
	"workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/plugins/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
)

/*
 * RegisterRoutes registers plugin HTTP routes.
 *
 * Base path: /api/v1/plugins
 *
 * HTTP Verbs follow REST conventions:
 *   GET    /                    - List plugins (paginated, filtered, multi-tenant)
 *   POST   /                    - Create plugin manifest
 *   GET    /enabled             - Get all enabled plugins (for editor boot)
 *   GET    /:id                 - Get plugin by MongoDB ObjectId
 *   PATCH  /:id                 - Update plugin manifest
 *   DELETE /:id                 - Delete plugin manifest
 */
func RegisterRoutes(group fiber.Router, service ports.PluginServicePort) {

	// List plugins with filters, pagination, and multi-tenant visibility
	pluginQueryDto := validation.NewValidation(nil, &dtos.PluginQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(pluginQueryDto),
		permissionMw.RequirePermission(perms.PluginList),
		coverageMw.InjectRequestContext(),
		handlers.GetPlugins(service),
	)

	// Get all enabled plugins (for editor boot — no validation needed)
	group.Get("/enabled",
		permissionMw.RequirePermission(perms.PluginRead),
		handlers.GetEnabledPlugins(service),
	)

	// Create a new plugin manifest
	pluginCreateDto := validation.NewValidation(&dtos.PluginManifestResponse{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(pluginCreateDto),
		permissionMw.RequirePermission(perms.PluginCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreatePlugin(service),
	)

	// Get plugin by MongoDB ObjectId
	pluginIdDto := validation.NewValidation(nil, nil, &dtos.PluginIdDTO{})
	group.Get("/:id",
		validation.ValidationMiddleware(pluginIdDto),
		permissionMw.RequirePermission(perms.PluginRead),
		handlers.GetPluginById(service),
	)

	// Update plugin manifest
	pluginUpdateDto := validation.NewValidation(&dtos.PluginManifestUpdate{}, nil, &dtos.PluginIdDTO{})
	group.Patch("/:id",
		validation.ValidationMiddleware(pluginUpdateDto),
		permissionMw.RequirePermission(perms.PluginUpdate),
		handlers.UpdatePlugin(service),
	)

	// Delete plugin manifest
	pluginDeleteDto := validation.NewValidation(nil, nil, &dtos.PluginIdDTO{})
	group.Delete("/:id",
		validation.ValidationMiddleware(pluginDeleteDto),
		permissionMw.RequirePermission(perms.PluginDelete),
		handlers.DeletePlugin(service),
	)
}
