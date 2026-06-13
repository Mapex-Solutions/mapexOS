package routes

import (
	"workflow/src/modules/archiver/application/dtos"
	"workflow/src/modules/archiver/application/ports"
	"workflow/src/modules/archiver/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers the workflow execution HTTP routes.
//
// Routes:
//
//	GET /              - List executions (paginated + filters)
//	GET /:executionId  - Get execution by ID
func RegisterRoutes(group web.Router, service ports.ArchiverServicePort) {

	r := swagger.Wrap(group).Tag("Executions")

	// List executions with filters and pagination.
	executionQueryDto := validation.NewValidation(nil, &dtos.ExecutionQueryDTO{}, nil)
	r.Get("/", executionQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowExecutionList),
		coverageMw.InjectRequestContext(),
		handlers.GetExecutions(service),
	).
		Summary("List executions").
		Description("Returns a paginated, filterable list of archived workflow executions scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.ExecutionResponseDTO]{})

	// Get execution by ID.
	getExecutionById := validation.NewValidation(nil, nil, &dtos.ExecutionIdDTO{})
	r.Get("/:executionId", getExecutionById, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowExecutionRead),
		handlers.GetExecutionById(service),
	).
		Summary("Get execution by ID").
		Description("Retrieves a single archived workflow execution by its ID.").
		Returns(&dtos.ExecutionResponseDTO{})
}
