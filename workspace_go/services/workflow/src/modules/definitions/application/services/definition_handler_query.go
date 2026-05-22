package services

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/definitions/application/dtos"
	"workflow/src/modules/definitions/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildDefinitionListFilters merges the org-scope filter with the query
// DTO's optional predicates (name, enabled, status, isTemplate, version).
// A filter-build failure increments the list-error metric so the dashboard
// observes failures even before the repo is touched.
func (s *DefinitionService) buildDefinitionListFilters(requestContext *reqCtx.RequestContext, query *dtos.DefinitionQueryDTO, start time.Time) (model.Map, error) {
	filters, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
		Query:      query,
	})
	if err != nil {
		s.trackDefinitionMetrics("list", "error", start)
		return nil, err
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Status != nil && *query.Status != "" {
		filters["status"] = *query.Status
	}
	if query.IsTemplate != nil {
		filters["isTemplate"] = *query.IsTemplate
	}
	if query.DefinitionVersion != nil {
		filters["definitionVersion"] = *query.DefinitionVersion
	}
	return filters, nil
}

// buildDefinitionPagination converts the int-typed query DTO fields into the
// model's int64 PaginationOpts. Helper exists so the orchestration in
// service.go does not have to know the DTO/model field-type mismatch.
func (s *DefinitionService) buildDefinitionPagination(query *dtos.DefinitionQueryDTO) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}

// mapDefinitionPaginatedResultTracked converts the entity page to the
// response DTO page, observes the list-results-count histogram, and records
// the success metric. Centralising the mapping + metric emission keeps
// service.go's GetDefinitions a recipe rather than a procedure.
func (s *DefinitionService) mapDefinitionPaginatedResultTracked(result *model.PaginatedResult[entities.WorkflowDefinition], start time.Time) (*model.PaginatedResult[dtos.DefinitionResponse], error) {
	dtoItems := make([]dtos.DefinitionResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.WorkflowDefinition, dtos.DefinitionResponse](&entity)
		if err != nil {
			return nil, fmt.Errorf("failed to map entity to response: %w", err)
		}
		dtoItems[i] = *dto
	}
	s.trackDefinitionMetrics("list", "success", start)
	s.deps.Metrics.DefinitionListResultsCount.Observe(float64(len(dtoItems)))
	return &model.PaginatedResult[dtos.DefinitionResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}

// toDefinitionResponseTracked maps the entity to the wire DTO and records
// the success metric for the supplied operation. Used by every CRUD method
// that returns a single definition (Create, Update, Read).
func (s *DefinitionService) toDefinitionResponseTracked(entity *entities.WorkflowDefinition, op string, start time.Time) (*dtos.DefinitionResponse, error) {
	resp, err := mapper.EntityToDto[entities.WorkflowDefinition, dtos.DefinitionResponse](entity)
	if err != nil {
		s.trackDefinitionMetrics(op, "error", start)
		return nil, fmt.Errorf("failed to map entity to response: %w", err)
	}
	s.trackDefinitionMetrics(op, "success", start)
	return resp, nil
}

// repopulateNodeScriptL2Async writes the served script back into MinIO from
// a fresh background context (the request context may cancel before the
// write completes). Failures log at warn — the cache miss already drained
// the slow path so a failed repopulate just means the next call will too.
func (s *DefinitionService) repopulateNodeScriptL2Async(def *entities.WorkflowDefinition, definitionId, nodeId, script string) {
	orgId := extractOrgId(def)
	go func() {
		bgCtx := context.Background()
		if err := s.deps.DefinitionStoragePort.WriteNodeScript(bgCtx, orgId, definitionId, nodeId, []byte(script)); err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Definition] Failed to repopulate L2 for node %s: %v", nodeId, err))
			return
		}
		logger.Info(fmt.Sprintf("[SERVICE:Definition] L2 repopulated for node %s (definition %s, org %s)", nodeId, definitionId, orgId))
	}()
}
