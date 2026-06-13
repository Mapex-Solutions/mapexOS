package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/interfaces/http/handlers"
)

// RegisterInternalRoutes registers internal API routes for TieredCache fallback.
//
// These routes are protected by API Key authentication and are used internally
// by other services (JS-Executor, Events) when L2 cache miss occurs.
//
// Routes:
//   - GET /:templateId - Fetch template and repopulate L2 cache
//
// Parameters:
//   - group: The Fiber router group for internal routes (e.g., /internal/templates)
//   - service: The AssetTemplateServicePort interface for business operations
func RegisterInternalRoutes(group web.Router, service ports.AssetTemplateServicePort) {
	group.Get("/:templateId", handlers.GetTemplateForCacheFallback(service))
}
