package routes

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/definitions/application/dtos"
	"workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/definitions/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
)

// RegisterRoutes registers all workflow definition HTTP routes.
//
// Routes:
//
//	GET    /           - List definitions (paginated + filters)
//	POST   /           - Create definition
//	GET    /:workflowId - Get definition by ID
//	PATCH  /:workflowId - Update definition
//	DELETE /:workflowId - Delete definition
func RegisterRoutes(group fiber.Router, service ports.DefinitionServicePort) {

	// Count definitions
	group.Get("/counter",
		permissionMw.RequirePermission(perms.WorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetDefinitionCount(service),
	)

	// List definitions with filters and pagination
	definitionQueryDto := validation.NewValidation(nil, &dtos.DefinitionQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(definitionQueryDto),
		permissionMw.RequirePermission(perms.WorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetDefinitions(service),
	)

	// Create a new definition
	definitionCreateDto := validation.NewValidation(&dtos.DefinitionCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(definitionCreateDto),
		permissionMw.RequirePermission(perms.WorkflowCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateDefinition(service),
	)

	// Get definition by ID
	getDefinitionById := validation.NewValidation(nil, nil, &dtos.DefinitionIdDTO{})
	group.Get("/:workflowId",
		validation.ValidationMiddleware(getDefinitionById),
		permissionMw.RequirePermission(perms.WorkflowRead),
		coverageMw.InjectRequestContext(),
		handlers.GetDefinitionById(service),
	)

	// Update definition by ID
	updateDefinitionById := validation.NewValidation(&dtos.DefinitionUpdateDTO{}, nil, &dtos.DefinitionIdDTO{})
	group.Patch("/:workflowId",
		validation.ValidationMiddleware(updateDefinitionById),
		permissionMw.RequirePermission(perms.WorkflowUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateDefinitionById(service),
	)

	// Delete definition by ID
	deleteDefinitionById := validation.NewValidation(nil, nil, &dtos.DefinitionIdDTO{})
	group.Delete("/:workflowId",
		validation.ValidationMiddleware(deleteDefinitionById),
		permissionMw.RequirePermission(perms.WorkflowDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteDefinitionById(service),
	)
}
