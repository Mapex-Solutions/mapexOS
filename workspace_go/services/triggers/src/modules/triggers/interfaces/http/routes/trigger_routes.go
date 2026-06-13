package routes

import (
	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/application/ports"
	"triggers/src/modules/triggers/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	permissions "github.com/Mapex-Solutions/MapexOS/permissions/triggers"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers trigger HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/triggers
//
// HTTP Verbs follow REST conventions:
//
//	GET    /         - List triggers (paginated, filtered)
//	GET    /counter  - Count triggers (cached)
//	POST   /         - Create trigger
//	GET    /:id      - Get trigger by ID
//	PATCH  /:id      - Update trigger
//	DELETE /:id      - Delete trigger
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
func RegisterRoutes(group web.Router, service ports.TriggerServicePort) {

	r := swagger.Wrap(group).Tag("Triggers")

	// List triggers with filters, pagination, and projection.
	// coverage middleware injects context-aware org filtering (hierarchical via PathKey).
	triggerQueryDto := validation.NewValidation(nil, &dtos.TriggerQueryDto{}, nil)
	r.Get("/", triggerQueryDto, swagger.Expose,
		permissionMw.RequirePermission(permissions.TriggerList),
		coverageMw.InjectRequestContext(),
		handlers.GetTriggers(service),
	).
		Summary("List triggers").
		Description("Returns a paginated, filterable list of triggers scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.TriggerResponse]{})

	// Count triggers (cached). No request contract to validate.
	counterDto := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", counterDto, swagger.Expose,
		permissionMw.RequirePermission(permissions.TriggerList),
		coverageMw.InjectRequestContext(),
		handlers.GetTriggerCount(service),
	).
		Summary("Count triggers").
		Description("Returns the total number of triggers for the caller's organization.").
		Returns(&contractsCommon.CounterResponse{})

	// Create a new trigger.
	triggerCreateDto := validation.NewValidation(&dtos.CreateTriggerDto{}, nil, nil)
	r.Post("/", triggerCreateDto, swagger.Expose,
		permissionMw.RequirePermission(permissions.TriggerCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateTrigger(service),
	).
		Summary("Create trigger").
		Description("Creates a new trigger. Org scoping is applied from the request context.").
		Returns(&dtos.TriggerResponse{})

	// Get trigger by ID. The :id path param is read directly by the handler and
	// is not validated as a contract, so no params DTO is declared here.
	getTriggerById := validation.NewValidation(nil, nil, nil)
	r.Get("/:id", getTriggerById, swagger.Expose,
		permissionMw.RequirePermission(permissions.TriggerRead),
		handlers.GetTriggerById(service),
	).
		Summary("Get trigger by ID").
		Description("Retrieves a single trigger by its MongoDB ObjectId.").
		Returns(&dtos.TriggerResponse{})

	// Update trigger by ID.
	triggerUpdateDto := validation.NewValidation(&dtos.UpdateTriggerDto{}, nil, nil)
	r.Patch("/:id", triggerUpdateDto, swagger.Expose,
		permissionMw.RequirePermission(permissions.TriggerUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateTriggerById(service),
	).
		Summary("Update trigger").
		Description("Partially updates an existing trigger. Only provided fields are changed.").
		Returns(&dtos.TriggerResponse{})

	// Delete trigger by ID.
	deleteTriggerById := validation.NewValidation(nil, nil, nil)
	r.Delete("/:id", deleteTriggerById, swagger.Expose,
		permissionMw.RequirePermission(permissions.TriggerDelete),
		handlers.DeleteTriggerById(service),
	).
		Summary("Delete trigger").
		Description("Deletes a trigger by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
