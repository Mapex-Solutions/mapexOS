package services

import (
	"context"
	"fmt"
	"strings"

	dtos "workflow/src/modules/archiver/application/dtos"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildExecutionListFilters builds the Mongo filter for GetExecutions.
// Org scope comes from the request context; status/instanceId/definitionId
// come from the query DTO. Status accepts a comma-separated list ($in).
func (s *ArchiverService) buildExecutionListFilters(rc *reqCtx.RequestContext, query *dtos.ExecutionQueryDTO) (model.Map, error) {
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	if err != nil {
		return nil, err
	}
	filters := orgFilter
	if query.Status != nil {
		statusStr := *query.Status
		if strings.Contains(statusStr, ",") {
			parts := strings.Split(statusStr, ",")
			filters["status"] = model.Map{"$in": parts}
		} else {
			filters["status"] = statusStr
		}
	}
	if query.InstanceID != nil {
		if instObjId, err := model.ToObjectID(*query.InstanceID); err == nil {
			filters["instanceId"] = instObjId
		}
	}
	if query.DefinitionID != nil {
		if defObjId, err := model.ToObjectID(*query.DefinitionID); err == nil {
			filters["definitionId"] = defObjId
		}
	}
	return filters, nil
}

// buildExecutionListPagination derives pagination opts from the query DTO.
func (s *ArchiverService) buildExecutionListPagination(query *dtos.ExecutionQueryDTO) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}

// mapExecutionListResult converts each entity in a paginated query result
// into its response DTO and re-wraps with the original pagination metadata.
func (s *ArchiverService) mapExecutionListResult(result *model.PaginatedResult[runtimePorts.WorkflowExecution]) *model.PaginatedResult[dtos.ExecutionResponseDTO] {
	dtoItems := make([]dtos.ExecutionResponseDTO, 0, len(result.Items))
	for _, exec := range result.Items {
		dtoItems = append(dtoItems, mapExecutionToDTO(&exec))
	}
	return &model.PaginatedResult[dtos.ExecutionResponseDTO]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}
}

// fetchExecutionOr404 loads an execution by id and surfaces the canonical
// not-found error when the repository returns nil.
func (s *ArchiverService) fetchExecutionOr404(ctx context.Context, executionId string) (*runtimePorts.WorkflowExecution, error) {
	exec, err := s.deps.ArchiveRepo.FindExecutionById(ctx, executionId)
	if err != nil {
		return nil, err
	}
	if exec == nil {
		return nil, fmt.Errorf("execution not found: %s", executionId)
	}
	return exec, nil
}

// enrichWithKVState replaces the lightweight stub with the full state from
// NATS KV when the execution is non-terminal. KV miss falls back to the
// MongoDB stub (logged as a warning).
func (s *ArchiverService) enrichWithKVState(exec *runtimePorts.WorkflowExecution) *runtimePorts.WorkflowExecution {
	if exec.Status.IsTerminal() || exec.WorkflowUUID == "" {
		return exec
	}
	kvExec, _, kvErr := s.fetchFullState(exec.WorkflowUUID)
	if kvErr == nil && kvExec != nil {
		return kvExec
	}
	logger.Warn(fmt.Sprintf("[SERVICE:Archiver] KV lookup failed for %s (using MongoDB stub): %v", exec.WorkflowUUID, kvErr))
	return exec
}
