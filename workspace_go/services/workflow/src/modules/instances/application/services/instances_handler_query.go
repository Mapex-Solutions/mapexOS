package services

import (
	"fmt"

	"workflow/src/modules/instances/application/dtos"
	"workflow/src/modules/instances/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildInstanceListFilters merges the org-scope filter (from request context)
// with module-specific predicates from the query DTO. Only fields the caller
// supplied are added so an empty query returns the full org scope.
func (s *InstancesService) buildInstanceListFilters(requestContext *reqCtx.RequestContext, query *dtos.InstanceQueryDTO) (model.Map, error) {
	filters, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
		Query:      query,
	})
	if err != nil {
		return nil, err
	}

	if query.DefinitionID != nil && *query.DefinitionID != "" {
		filters["definitionId"] = *query.DefinitionID
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.UniqueExecution != nil {
		filters["uniqueExecution"] = *query.UniqueExecution
	}
	return filters, nil
}

// buildInstancePagination converts the int-typed query DTO fields into the
// model's int64 PaginationOpts. Helper exists so the orchestration in
// service.go does not have to know the DTO/model field-type mismatch.
func (s *InstancesService) buildInstancePagination(query *dtos.InstanceQueryDTO) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}

// buildInstanceProjection translates the comma-separated projection field
// from the query string into a Mongo projection map. A nil result means
// "return all fields" — preserved as the no-op default.
func (s *InstancesService) buildInstanceProjection(query *dtos.InstanceQueryDTO) model.Map {
	if query.Projection == nil || *query.Projection == "" {
		return nil
	}
	return model.StringToProjection(*query.Projection)
}

// mapInstancePaginatedResult walks the entity list returned by the repo,
// mapping each through EntityToDto and rewrapping in the response page so
// the wire format matches the rest of the API surface.
func (s *InstancesService) mapInstancePaginatedResult(result *model.PaginatedResult[entities.WorkflowInstance]) (*model.PaginatedResult[dtos.InstanceResponse], error) {
	dtoItems := make([]dtos.InstanceResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.WorkflowInstance, dtos.InstanceResponse](&entity)
		if err != nil {
			return nil, fmt.Errorf("failed to map entity to response: %w", err)
		}
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.InstanceResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}
