package services

import (
	"context"
	"fmt"

	"workflow/src/modules/instances/application/di"
	"workflow/src/modules/instances/application/dtos"
	"workflow/src/modules/instances/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// Compile-time check
var _ ports.InstancesServicePort = (*InstancesService)(nil)

// New creates and returns a new InstancesService.
func New(deps di.InstancesServiceDependenciesInjection) ports.InstancesServicePort {
	return &InstancesService{deps: deps}
}

// CreateInstance materialises a new workflow instance config from the create
// DTO, stamps it with the org context (orgId + pathKey) for multi-tenant
// scoping, persists it, and returns the response DTO.
func (s *InstancesService) CreateInstance(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.InstanceCreateDTO) (*dtos.InstanceResponse, error) {
	entity := s.buildInstanceEntity(dto)
	s.applyOrgContextToInstance(entity, requestContext)
	created, err := s.deps.InstanceRepo.Create(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}
	return s.toInstanceResponse(created)
}

// UpdateInstanceById applies a partial update (only fields set on the DTO are
// changed) and invalidates the loader cache so the next read sees the
// updated entity instead of the prior cached snapshot.
func (s *InstancesService) UpdateInstanceById(ctx context.Context, instanceId *string, dto *dtos.InstanceUpdateDTO) (*dtos.InstanceResponse, error) {
	payload := s.buildInstanceUpdatePayload(dto)
	updated, err := s.deps.InstanceRepo.FindByIdAndUpdate(ctx, instanceId, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update instance: %w", err)
	}
	_ = s.deps.InstanceLoader.Invalidate(ctx, *instanceId)
	return s.toInstanceResponse(updated)
}

// DeleteInstanceById removes the instance from the repository and clears its
// loader-cache entry so subsequent reads do not resurrect a stale copy.
func (s *InstancesService) DeleteInstanceById(ctx context.Context, instanceId *string) error {
	if err := s.deps.InstanceRepo.DeleteById(ctx, instanceId); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	_ = s.deps.InstanceLoader.Invalidate(ctx, *instanceId)
	return nil
}

// GetInstanceById reads through the TieredCache loader for the given id.
// Returns (nil, nil) when not found — the caller distinguishes "missing" from
// errors via the second return value.
func (s *InstancesService) GetInstanceById(ctx context.Context, instanceId *string) (*dtos.InstanceResponse, error) {
	entity, err := s.deps.InstanceLoader.GetInstance(ctx, *instanceId)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	return s.toInstanceResponse(entity)
}

// GetInstances returns a paginated, filtered, projected list scoped by org
// context. The orchestration delegates filter assembly, pagination, and
// projection to dedicated helpers so the public method stays as a recipe.
func (s *InstancesService) GetInstances(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.InstanceQueryDTO) (*model.PaginatedResult[dtos.InstanceResponse], error) {
	filters, err := s.buildInstanceListFilters(requestContext, query)
	if err != nil {
		return nil, err
	}
	pagination := s.buildInstancePagination(query)
	projection := s.buildInstanceProjection(query)

	result, err := s.deps.InstanceRepo.FindWithFilters(ctx, filters, pagination, projection)
	if err != nil {
		return nil, err
	}
	return s.mapInstancePaginatedResult(result)
}

// CountInstances returns the count of instances visible to the request's org
// scope. Errors are logged at the service layer (not just bubbled) because
// the dashboard endpoint expects a 0 fallback for a usable degraded UX.
func (s *InstancesService) CountInstances(ctx context.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	filters, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
	})
	if err != nil {
		logger.Error(err, "[SERVICE:Instances] Failed to build org filter for counter")
		return 0, err
	}
	count, err := s.deps.InstanceRepo.CountDocuments(ctx, filters)
	if err != nil {
		logger.Error(err, "[SERVICE:Instances] Failed to count instances")
		return 0, err
	}
	return count, nil
}

// ExecuteInstance is the synchronous HTTP-driven path: forward the call to
// the RuntimeService port and reshape the runtime result into the instances
// response DTO. The body shows the cross-module composition explicitly.
func (s *InstancesService) ExecuteInstance(ctx context.Context, instanceId string, eventPayload map[string]interface{}, workflowUUID string) (*dtos.ExecuteResponseDTO, error) {
	result, err := s.deps.RuntimeService.ExecuteByInstanceID(ctx, instanceId, eventPayload, workflowUUID)
	if err != nil {
		return nil, err
	}
	return &dtos.ExecuteResponseDTO{
		WorkflowUUID: result.WorkflowUUID,
		Status:       result.Status,
		ErrorInfo:    result.ErrorInfo,
	}, nil
}
