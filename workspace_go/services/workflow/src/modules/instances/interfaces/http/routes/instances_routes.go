package routes

import (
	"workflow/src/modules/instances/application/dtos"
	"workflow/src/modules/instances/application/ports"
	"workflow/src/modules/instances/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers all workflow instance config HTTP routes.
//
// Routes:
//
//	GET    /                      - List instance configs (paginated + filters)
//	GET    /:instanceId           - Get instance config by ID
//	POST   /                      - Create instance config
//	PUT    /:instanceId           - Update instance config
//	POST   /:instanceId/execute   - Execute a workflow instance
//	DELETE /:instanceId           - Delete instance config
func RegisterRoutes(group web.Router, service ports.InstancesServicePort) {

	r := swagger.Wrap(group).Tag("Instances")

	// Count instances (no request contract to validate).
	counterDto := validation.NewValidation(nil, nil, nil)
	r.Get("/counter", counterDto, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceList),
		coverageMw.InjectRequestContext(),
		handlers.GetInstanceCount(service),
	).
		Summary("Count instances").
		Description("Returns the total number of workflow instances for the caller's organization.").
		Returns(&contractsCommon.CounterResponse{})

	// List instance configs with filters and pagination.
	instanceQueryDto := validation.NewValidation(nil, &dtos.InstanceQueryDTO{}, nil)
	r.Get("/", instanceQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceList),
		coverageMw.InjectRequestContext(),
		handlers.GetInstances(service),
	).
		Summary("List instances").
		Description("Returns a paginated, filterable list of workflow instances scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.InstanceResponse]{})

	// Get instance config by ID.
	getInstanceById := validation.NewValidation(nil, nil, &dtos.InstanceIdDTO{})
	r.Get("/:instanceId", getInstanceById, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceRead),
		coverageMw.InjectRequestContext(),
		handlers.GetInstanceById(service),
	).
		Summary("Get instance by ID").
		Description("Retrieves a single workflow instance by its MongoDB ObjectId.").
		Returns(&dtos.InstanceResponse{})

	// Create instance config.
	createInstance := validation.NewValidation(&dtos.InstanceCreateDTO{}, nil, nil)
	r.Post("/", createInstance, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateInstance(service),
	).
		Summary("Create instance").
		Description("Creates a new workflow instance config. Org scoping is applied from the request context.").
		Returns(&dtos.InstanceResponse{})

	// Update instance config.
	updateInstance := validation.NewValidation(&dtos.InstanceUpdateDTO{}, nil, &dtos.InstanceIdDTO{})
	r.Put("/:instanceId", updateInstance, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateInstanceById(service),
	).
		Summary("Update instance").
		Description("Replaces an existing workflow instance config by its MongoDB ObjectId.").
		Returns(&dtos.InstanceResponse{})

	// Execute a workflow instance on demand.
	executeInstance := validation.NewValidation(nil, nil, &dtos.InstanceIdDTO{})
	r.Post("/:instanceId/execute", executeInstance, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceExecute),
		coverageMw.InjectRequestContext(),
		handlers.ExecuteInstance(service),
	).
		Summary("Execute instance").
		Description("Starts an execution of the given workflow instance and returns the execution identifiers.").
		Returns(&dtos.ExecuteResponseDTO{})

	// Delete instance config.
	deleteInstance := validation.NewValidation(nil, nil, &dtos.InstanceIdDTO{})
	r.Delete("/:instanceId", deleteInstance, swagger.Expose,
		permissionMw.RequirePermission(perms.WorkflowInstanceDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteInstanceById(service),
	).
		Summary("Delete instance").
		Description("Deletes a workflow instance config by its MongoDB ObjectId.")
}
