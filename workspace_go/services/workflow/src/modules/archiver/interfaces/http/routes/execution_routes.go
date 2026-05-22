package routes

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/archiver/application/dtos"
	"workflow/src/modules/archiver/application/ports"
	"workflow/src/modules/archiver/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
)

// RegisterRoutes registers the workflow execution HTTP routes.
func RegisterRoutes(group fiber.Router, service ports.ArchiverServicePort) {

	// List executions with filters and pagination
	executionQueryDto := validation.NewValidation(nil, &dtos.ExecutionQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(executionQueryDto),
		permissionMw.RequirePermission(perms.WorkflowExecutionList),
		coverageMw.InjectRequestContext(),
		handlers.GetExecutions(service),
	)

	// Get execution by ID
	getExecutionById := validation.NewValidation(nil, nil, &dtos.ExecutionIdDTO{})
	group.Get("/:executionId",
		validation.ValidationMiddleware(getExecutionById),
		permissionMw.RequirePermission(perms.WorkflowExecutionRead),
		handlers.GetExecutionById(service),
	)
}
