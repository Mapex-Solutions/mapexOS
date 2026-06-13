package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/definitions/interfaces/http/handlers"
)

// RegisterInternalRoutes registers internal HTTP routes with API key authentication.
//
// These endpoints are used by js-workflow-executor as a TieredCache fallback
// when L2 (MinIO) cache miss occurs. They fetch script source from MongoDB
// and repopulate L2 for future requests.
//
// Security: All routes are protected by API Key middleware (X-API-Key header)
//
// Endpoints:
//
//	GET /:definitionId/scripts/:nodeId - Fetch code node script source
func RegisterInternalRoutes(group web.Router, service ports.DefinitionServicePort) {
	group.Get("/:definitionId/scripts/:nodeId", handlers.GetNodeScript(service))
}
