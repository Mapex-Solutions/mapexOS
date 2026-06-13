package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"
	"router/src/modules/routegroups/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/router"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
)

// RegisterRoutes registers route group HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/route_groups
//
// HTTP Verbs follow REST conventions:
//
//	GET    /                    - List route groups (paginated, filtered)
//	POST   /                    - Create route group
//	GET    /:routeGroupId       - Get route group by ID
//	PATCH  /:routeGroupId       - Update route group
//	DELETE /:routeGroupId       - Delete route group
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Route group service port interface implementation
func RegisterRoutes(group web.Router, service ports.RouteGroupServicePort) {

	r := swagger.Wrap(group).Tag("Route Groups")

	/** List Routes */

	// Get route groups with filters, pagination, and projection.
	// coverage middleware injects context-aware org filtering (hierarchical via PathKey).
	routeGroupQueryDto := validation.NewValidation(nil, &dtos.RouteGroupQueryDTO{}, nil)
	r.Get("/", routeGroupQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RouteGroupList),
		coverageMw.InjectRequestContext(),
		handlers.GetRouteGroups(service),
	).
		Summary("List route groups").
		Description("Returns a paginated, filterable list of route groups scoped to the caller's organization. Supports hierarchical queries via includeChildren.").
		Returns(&model.PaginatedResult[dtos.RouteGroupResponse]{})

	/** Counter Route */

	// Count route groups (cached, 6h TTL). No request contract to validate.
	counterDto := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", counterDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RouteGroupList),
		coverageMw.InjectRequestContext(),
		handlers.GetRouteGroupCount(service),
	).
		Summary("Count route groups").
		Description("Returns the total number of route groups for the caller's organization. Uses a cached count with a 6h TTL, invalidated on create and delete.").
		Returns(&contractsCommon.CounterResponse{})

	/** CRUD Routes */

	// Create a new route group.
	listCreateDto := validation.NewValidation(&dtos.RouteGroupCreateDTO{}, nil, nil)
	r.Post("/", listCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RouteGroupCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateRouteGroup(service),
	).
		Summary("Create route group").
		Description("Creates a new route group. Organization scoping (orgId, pathKey) is applied automatically from the request context.").
		Returns(&dtos.RouteGroupResponse{})

	// Get route group by ID.
	getRouteGroupById := validation.NewValidation(nil, nil, &dtos.RouteGroupIdDTO{})
	r.Get("/:routeGroupId", getRouteGroupById, swagger.Expose,
		permissionMw.RequirePermission(perms.RouteGroupRead),
		coverageMw.InjectRequestContext(),
		handlers.GetRouteGroupById(service),
	).
		Summary("Get route group by ID").
		Description("Retrieves a single route group by its MongoDB ObjectId.").
		Returns(&dtos.RouteGroupResponse{})

	// Update route group by ID.
	updateRouteGroupById := validation.NewValidation(&dtos.RouteGroupUpdateDTO{}, nil, &dtos.RouteGroupIdDTO{})
	r.Patch("/:routeGroupId", updateRouteGroupById, swagger.Expose,
		permissionMw.RequirePermission(perms.RouteGroupUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateRouteGroupById(service),
	).
		Summary("Update route group").
		Description("Partially updates an existing route group. All body fields are optional; only provided fields are changed.").
		Returns(&dtos.RouteGroupResponse{})

	// Delete route group by ID.
	deleteRouteGroupById := validation.NewValidation(nil, nil, &dtos.RouteGroupIdDTO{})
	r.Delete("/:routeGroupId", deleteRouteGroupById, swagger.Expose,
		permissionMw.RequirePermission(perms.RouteGroupDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteRouteGroupById(service),
	).
		Summary("Delete route group").
		Description("Deletes a route group by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
