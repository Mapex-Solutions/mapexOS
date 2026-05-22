package routes

import (
	"github.com/gofiber/fiber/v2"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"
	"router/src/modules/routegroups/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/router"
)

// RegisterRoutes registers route group HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/route_groups
//
// HTTP Verbs follow REST conventions:
//   GET    /                    - List route groups (paginated, filtered)
//   POST   /                    - Create route group
//   GET    /:routeGroupId       - Get route group by ID
//   PATCH  /:routeGroupId       - Update route group
//   DELETE /:routeGroupId       - Delete route group
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Route group service port interface implementation
func RegisterRoutes(group fiber.Router, service ports.RouteGroupServicePort) {

	/**
	* List Routes
	 */

	// Get route groups with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Includes hierarchical support via PathKey data from coverage cache
	routeGroupQueryDto := validation.NewValidation(nil, &dtos.RouteGroupQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(routeGroupQueryDto),   // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RouteGroupList),  // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                     // 3. Inject context (cache)
		handlers.GetRouteGroups(service),                      // 4. Handler
	)

	/**
	* Counter Route
	 */

	// Count route groups (counter endpoint with cache)
	group.Get("/counter",
		permissionMw.RequirePermission(perms.RouteGroupList),
		coverageMw.InjectRequestContext(),
		handlers.GetRouteGroupCount(service),
	)

	/**
	* CRUD Routes
	 */

	// Create a new routegroup
	listCreateDto := validation.NewValidation(&dtos.RouteGroupCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(listCreateDto),        // 1. Validate DTO
		permissionMw.RequirePermission(perms.RouteGroupCreate), // 2. Check permission
		coverageMw.InjectRequestContext(),                     // 3. Inject context
		handlers.CreateRouteGroup(service),                    // 4. Handler
	)

	// Get routegroup by ID
	getRouteGroupById := validation.NewValidation(nil, nil, &dtos.RouteGroupIdDTO{})
	group.Get("/:routeGroupId",
		validation.ValidationMiddleware(getRouteGroupById),          // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RouteGroupRead),       // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                          // 3. Inject context (cache)
		handlers.GetRouteGroupById(service),                        // 4. Handler
	)

	// Update routegroup by ID
	updateRouteGroupById := validation.NewValidation(&dtos.RouteGroupUpdateDTO{}, nil, &dtos.RouteGroupIdDTO{})
	group.Patch("/:routeGroupId",
		validation.ValidationMiddleware(updateRouteGroupById),      // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RouteGroupUpdate),     // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                          // 3. Inject context (cache)
		handlers.UpdateRouteGroupById(service),                     // 4. Handler
	)

	// Delete routegroup by ID
	deleteRouteGroupById := validation.NewValidation(nil, nil, &dtos.RouteGroupIdDTO{})
	group.Delete("/:routeGroupId",
		validation.ValidationMiddleware(deleteRouteGroupById),      // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RouteGroupDelete),     // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                          // 3. Inject context (cache)
		handlers.DeleteRouteGroupById(service),                     // 4. Handler
	)
}
