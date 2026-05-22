package ports

import (
	ctx "context"

	"workflow/src/modules/instances/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// InstancesServicePort defines the contract for workflow instance config CRUD operations.
type InstancesServicePort interface {
	// CreateInstance creates a new workflow instance config.
	// requestContext provides orgId and pathKey from the coverage middleware.
	CreateInstance(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.InstanceCreateDTO) (*dtos.InstanceResponse, error)

	// GetInstanceById retrieves a workflow instance config by its unique identifier.
	GetInstanceById(ctx ctx.Context, instanceId *string) (*dtos.InstanceResponse, error)

	// UpdateInstanceById updates a workflow instance config.
	UpdateInstanceById(ctx ctx.Context, instanceId *string, dto *dtos.InstanceUpdateDTO) (*dtos.InstanceResponse, error)

	// DeleteInstanceById deletes a workflow instance config.
	DeleteInstanceById(ctx ctx.Context, instanceId *string) error

	// CountInstances returns the total count of workflow instances for the given org context.
	CountInstances(ctx ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)

	// GetInstances retrieves a paginated and filtered list of workflow instance configs.
	GetInstances(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.InstanceQueryDTO) (*model.PaginatedResult[dtos.InstanceResponse], error)

	// ExecuteInstance executes a workflow instance by ID. Returns execution UUID, status, and error info.
	ExecuteInstance(ctx ctx.Context, instanceId string, eventPayload map[string]interface{}, workflowUUID string) (*dtos.ExecuteResponseDTO, error)
}
