package routes

import (
	"workflow/src/modules/definitions/application/dtos"
	"workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/definitions/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers all workflow definition HTTP routes.
//
// Routes:
//
//	GET    /            - List definitions (paginated + filters)
//	POST   /            - Create definition
//	GET    /:workflowId - Get definition by ID
//	PATCH  /:workflowId - Update definition
//	DELETE /:workflowId - Delete definition
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
func RegisterRoutes(group web.Router, service ports.DefinitionServicePort) {

	r := swagger.Wrap(group).Tag("Definitions")

	// Count definitions (no request contract to validate).
	counterDto := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", counterDto, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetDefinitionCount(service),
	).
		Summary("Count definitions").
		Description("Returns the total number of workflow definitions for the caller's organization.").
		Returns(&contractsCommon.CounterResponse{})

	// List definitions with filters and pagination.
	definitionQueryDto := validation.NewValidation(nil, &dtos.DefinitionQueryDTO{}, nil)
	r.Get("/", definitionQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetDefinitions(service),
	).
		Summary("List definitions").
		Description("Returns a paginated, filterable list of workflow definitions scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.DefinitionResponse]{})

	// Create a new definition.
	definitionCreateDto := validation.NewValidation(&dtos.DefinitionCreateDTO{}, nil, nil)
	r.Post("/", definitionCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateDefinition(service),
	).
		Summary("Create definition").
		Description("Creates a new workflow definition. Org scoping is applied from the request context.").
		Returns(&dtos.DefinitionResponse{})

	// Get definition by ID.
	getDefinitionById := validation.NewValidation(nil, nil, &dtos.DefinitionIdDTO{})
	r.Get("/:workflowId", getDefinitionById, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowRead),
		coverageMw.InjectRequestContext(),
		handlers.GetDefinitionById(service),
	).
		Summary("Get definition by ID").
		Description("Retrieves a single workflow definition by its MongoDB ObjectId.").
		Returns(&dtos.DefinitionResponse{})

	// Update definition by ID.
	updateDefinitionById := validation.NewValidation(&dtos.DefinitionUpdateDTO{}, nil, &dtos.DefinitionIdDTO{})
	r.Patch("/:workflowId", updateDefinitionById, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateDefinitionById(service),
	).
		Summary("Update definition").
		Description("Partially updates an existing workflow definition. Only provided fields are changed.").
		Returns(&dtos.DefinitionResponse{})

	// Delete definition by ID.
	deleteDefinitionById := validation.NewValidation(nil, nil, &dtos.DefinitionIdDTO{})
	r.Delete("/:workflowId", deleteDefinitionById, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteDefinitionById(service),
	).
		Summary("Delete definition").
		Description("Deletes a workflow definition by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
