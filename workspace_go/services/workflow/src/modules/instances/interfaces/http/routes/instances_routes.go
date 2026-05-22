package routes

import (
	"github.com/gofiber/fiber/v2"

	"workflow/src/modules/instances/application/dtos"
	"workflow/src/modules/instances/application/ports"
	"workflow/src/modules/instances/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
)

// RegisterRoutes registers all workflow instance config HTTP routes.
//
// Routes:
//
//	GET    /              - List instance configs (paginated + filters)
//	GET    /:instanceId   - Get instance config by ID
//	POST   /              - Create instance config
//	PUT    /:instanceId   - Update instance config
//	DELETE /:instanceId   - Delete instance config
func RegisterRoutes(group fiber.Router, service ports.InstancesServicePort) {

	// Count instances
	group.Get("/counter",
		permissionMw.RequirePermission(perms.WorkflowInstanceList),
		coverageMw.InjectRequestContext(),
		handlers.GetInstanceCount(service),
	)

	// List instance configs with filters and pagination
	instanceQueryDto := validation.NewValidation(nil, &dtos.InstanceQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(instanceQueryDto),
		permissionMw.RequirePermission(perms.WorkflowInstanceList),
		coverageMw.InjectRequestContext(),
		handlers.GetInstances(service),
	)

	// Get instance config by ID
	getInstanceById := validation.NewValidation(nil, nil, &dtos.InstanceIdDTO{})
	group.Get("/:instanceId",
		validation.ValidationMiddleware(getInstanceById),
		permissionMw.RequirePermission(perms.WorkflowInstanceRead),
		coverageMw.InjectRequestContext(),
		handlers.GetInstanceById(service),
	)

	// Create instance config
	createInstance := validation.NewValidation(&dtos.InstanceCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(createInstance),
		permissionMw.RequirePermission(perms.WorkflowInstanceCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateInstance(service),
	)

	// Update instance config
	updateInstance := validation.NewValidation(&dtos.InstanceUpdateDTO{}, nil, &dtos.InstanceIdDTO{})
	group.Put("/:instanceId",
		validation.ValidationMiddleware(updateInstance),
		permissionMw.RequirePermission(perms.WorkflowInstanceUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateInstanceById(service),
	)

	// Execute workflow instance
	executeInstance := validation.NewValidation(nil, nil, &dtos.InstanceIdDTO{})
	group.Post("/:instanceId/execute",
		validation.ValidationMiddleware(executeInstance),
		permissionMw.RequirePermission(perms.WorkflowInstanceExecute),
		coverageMw.InjectRequestContext(),
		handlers.ExecuteInstance(service),
	)

	// Delete instance config
	deleteInstance := validation.NewValidation(nil, nil, &dtos.InstanceIdDTO{})
	group.Delete("/:instanceId",
		validation.ValidationMiddleware(deleteInstance),
		permissionMw.RequirePermission(perms.WorkflowInstanceDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteInstanceById(service),
	)
}
