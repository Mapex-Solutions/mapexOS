package routes

import (
	"github.com/gofiber/fiber/v2"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"
	"router/src/modules/routegroups/interfaces/http/handlers"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
)

// RegisterInternalRoutes registers internal API routes for inter-service communication.
// These routes use API Key authentication (X-API-Key header) and trust the calling service.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/internal/v1/routegroups
//
// Internal Routes:
//   GET / - Get multiple route groups by IDs (comma-separated in query)
//
// Authentication: API Key (X-API-Key header)
//
// Example:
//   GET /api/internal/v1/routegroups?ids=id1,id2,id3&projection=name,enabled
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: RouteGroup service port interface implementation
func RegisterInternalRoutes(group fiber.Router, service ports.RouteGroupServicePort) {

	/**
	 * INTERNAL Routes - MS-to-MS communication
	 */

	// Get multiple route groups by IDs
	// GET /?ids=id1,id2,id3&projection=name,enabled
	getByIdsQuery := validation.NewValidation(nil, &dtos.RouteGroupInternalIdsQuery{}, nil)
	group.Get("/", validation.ValidationMiddleware(getByIdsQuery), handlers.GetRouteGroupsByIds(service))
}
