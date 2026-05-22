package services

import (
	ctx "context"
	"time"

	"router/src/modules/routegroups/application/di"
	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"
	"router/src/modules/routegroups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// Compile-time check to ensure RouteGroupService implements RouteGroupServicePort.
var _ ports.RouteGroupServicePort = (*RouteGroupService)(nil)

// New creates and returns a new instance of RouteGroupService.
func New(deps di.RouteGroupServiceDependenciesInjection) ports.RouteGroupServicePort {
	return &RouteGroupService{deps: deps}
}

// CreateRouteGroup orchestrates RouteGroup creation: enforce multi-tenant
// scope -> map DTO to entity -> persist -> warm caches -> return response DTO.
func (s *RouteGroupService) CreateRouteGroup(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.RouteGroupCreateDTO) (*dtos.RouteGroupResponse, error) {
	start := time.Now()
	if err := s.applyRouteGroupScope(requestContext, dto); err != nil {
		return nil, s.recordRouteGroupError("create", start, err)
	}
	entity, err := s.persistNewRouteGroup(c, dto)
	if err != nil {
		return nil, s.recordRouteGroupError("create", start, err)
	}
	s.warmCachesAfterCreate(c, requestContext, entity)
	response, err := s.mapRouteGroupResponse(entity)
	if err != nil {
		return nil, s.recordRouteGroupError("create", start, err)
	}
	s.recordRouteGroupOp("create", "success", start)
	return response, nil
}

// GetRouteGroupById fetches a RouteGroup by id (cache-first) and returns its
// response DTO.
func (s *RouteGroupService) GetRouteGroupById(c ctx.Context, routeGroupId *string) (*dtos.RouteGroupResponse, error) {
	entity, err := s.GetRouteGroupEntityById(c, routeGroupId)
	if err != nil {
		return nil, err
	}
	return s.mapRouteGroupResponse(entity)
}

// GetRouteGroupEntityById is the cross-module entry point that returns the
// raw RouteGroup entity; cache-aside with metric instrumentation per call.
func (s *RouteGroupService) GetRouteGroupEntityById(c ctx.Context, routeGroupId *string) (*entities.RouteGroup, error) {
	start := time.Now()
	entity, cacheMiss, err := s.loadCachedRouteGroup(c, routeGroupId)
	s.recordRouteGroupCacheOutcome(cacheMiss)
	if err != nil || entity.ID.IsZero() {
		s.recordRouteGroupOp("read", "error", start)
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"route group not found"}}
	}
	s.recordRouteGroupOp("read", "success", start)
	return entity, nil
}

// UpdateRouteGroupById applies a partial update by id and returns the
// resulting response DTO. 404 when the id is unknown.
func (s *RouteGroupService) UpdateRouteGroupById(c ctx.Context, listId *string, dto *dtos.RouteGroupUpdateDTO) (*dtos.RouteGroupResponse, error) {
	start := time.Now()
	entity, err := s.applyRouteGroupUpdateFromDto(c, listId, dto)
	if err != nil {
		return nil, s.recordRouteGroupError("update", start, err)
	}
	s.cacheRouteGroupEntity(c, *listId, entity)
	response, err := s.mapRouteGroupResponse(entity)
	if err != nil {
		return nil, s.recordRouteGroupError("update", start, err)
	}
	s.recordRouteGroupOp("update", "success", start)
	return response, nil
}

// DeleteRouteGroupById removes a RouteGroup, invalidates its caches, and
// drops the per-org counter cache. 404 when the id is unknown.
func (s *RouteGroupService) DeleteRouteGroupById(c ctx.Context, listId *string) (map[string]bool, error) {
	start := time.Now()
	routeGroup, err := s.fetchRouteGroupForDelete(c, listId)
	if err != nil {
		s.recordRouteGroupOp("delete", "error", start)
		return nil, err
	}
	if err := s.deleteRouteGroupFromRepo(c, listId); err != nil {
		s.recordRouteGroupOp("delete", "error", start)
		return nil, err
	}
	s.deleteRouteGroupCache(c, listId)
	s.invalidateCounterCacheForDelete(c, routeGroup)
	s.recordRouteGroupOp("delete", "success", start)
	return map[string]bool{"success": true}, nil
}

// GetRouteGroupsByIds resolves a list of RouteGroups by id; missing entries
// are silently skipped so the caller still receives the partial set.
func (s *RouteGroupService) GetRouteGroupsByIds(c ctx.Context, ids []string) ([]dtos.RouteGroupResponse, error) {
	results := make([]dtos.RouteGroupResponse, 0, len(ids))
	for _, id := range ids {
		routeGroup, err := s.GetRouteGroupById(c, &id)
		if err != nil {
			continue
		}
		results = append(results, *routeGroup)
	}
	return results, nil
}

// GetRouteGroups returns a paginated, filtered RouteGroup list scoped to the
// caller's org context (own org + system templates + ancestor templates).
func (s *RouteGroupService) GetRouteGroups(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.RouteGroupQueryDTO) (*model.PaginatedResult[dtos.RouteGroupResponse], error) {
	start := time.Now()
	filters := s.buildRouteGroupListFilters(requestContext, query)
	pagination := s.buildRouteGroupListPagination(query)
	projection := s.buildRouteGroupListProjection(query)
	result, err := s.deps.RouteGroupRepo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, s.recordRouteGroupError("list", start, err)
	}
	dtoResult, err := s.mapRouteGroupListResult(result)
	if err != nil {
		return nil, s.recordRouteGroupError("list", start, err)
	}
	s.recordRouteGroupListSuccess(start, dtoResult)
	return dtoResult, nil
}

// CountRouteGroups returns the per-org RouteGroup count with a cache-aside
// 6h TTL — counter cache is invalidated on every Create/Delete.
func (s *RouteGroupService) CountRouteGroups(c ctx.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	orgId := ""
	if requestContext.OrgContext != nil {
		orgId = *requestContext.OrgContext
	}
	cacheKey := s.buildCounterCacheKey(orgId)
	if count, ok := s.tryCachedRouteGroupCount(c, cacheKey); ok {
		return count, nil
	}
	count, err := s.countRouteGroupsFromRepo(c, requestContext)
	if err != nil {
		return 0, err
	}
	s.cacheRouteGroupCount(c, cacheKey, count)
	return count, nil
}
