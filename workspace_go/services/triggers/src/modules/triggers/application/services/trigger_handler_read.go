package services

import (
	ctx "context"
	"fmt"

	"triggers/src/modules/triggers/application/constants"
	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/domain/entities"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// loadTriggerCacheAside runs the GetOrSetEx cache-aside pattern for the
// per-id Redis cache. Optional CacheMetrics is propagated when provided
// so the caller can record per-call hit/miss observability.
func (s *TriggerService) loadTriggerCacheAside(c ctx.Context, triggerId *string, dest *entities.Trigger, metrics ...*common.CacheMetrics) error {
	cacheKey := s.deps.CacheKeyBuilder.TriggerKey(*triggerId)
	params := common.GetOrSetParams{
		Ctx:      c,
		CacheKey: cacheKey,
		CacheTTL: int(constants.TriggerCacheTTL.Seconds()),
		Dest:     dest,
		Callback: func() (interface{}, error) {
			return s.deps.TriggerRepository.FindById(c, triggerId)
		},
	}
	if len(metrics) > 0 && metrics[0] != nil {
		params.Metrics = metrics[0]
	}
	_, err := s.deps.CacheRepository.GetOrSetEx(params)
	return err
}

// buildTriggerListFilters assembles the Mongo filter for GetTriggers:
// org-scope $or (own org + system templates + ancestor templates) plus
// per-field filters (id, name, type, category, enabled).
func (s *TriggerService) buildTriggerListFilters(rc *reqCtx.RequestContext, query *dtos.TriggerQueryDto) model.Map {
	orConditions := []model.Map{}
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	if err == nil && len(orgFilter) > 0 {
		orConditions = append(orConditions, orgFilter)
	}
	if query.IsSystem == nil || *query.IsSystem {
		orConditions = append(orConditions, model.Map{"isSystem": true})
	}
	if query.IsTemplate != nil && *query.IsTemplate {
		if templateFilter, err := orgfilter.BuildTemplateAncestorFilter(rc); err == nil && len(templateFilter) > 0 {
			orConditions = append(orConditions, templateFilter)
		}
	}
	filters := model.Map{}
	if len(orConditions) > 0 {
		filters["$or"] = orConditions
	}
	if query.IsSystem != nil && !*query.IsSystem {
		filters["isSystem"] = false
	}
	if query.IsTemplate != nil && !*query.IsTemplate {
		filters["isTemplate"] = false
	}
	if query.ID != nil {
		objId, _ := model.ToObjectID(*query.ID)
		filters["_id"] = objId
	}
	if query.Name != nil {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.TriggerType != nil {
		filters["triggerType"] = *query.TriggerType
	}
	if query.Category != nil {
		filters["category"] = *query.Category
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	return filters
}

// buildTriggerListPagination derives the Mongo pagination opts from the
// query DTO using the shared per-page defaults.
func (s *TriggerService) buildTriggerListPagination(query *dtos.TriggerQueryDto) *model.PaginationOpts {
	pagination := &model.PaginationOpts{Page: 1, PerPage: 10}
	if query.Page != nil && *query.Page > 0 {
		pagination.Page = int64(*query.Page)
	}
	if query.PerPage != nil && *query.PerPage > 0 {
		pagination.PerPage = int64(*query.PerPage)
	}
	return pagination
}

// mapTriggerListResult converts each entity in a paginated query result to
// its response DTO and re-wraps with the original pagination metadata.
func (s *TriggerService) mapTriggerListResult(result *model.PaginatedResult[entities.Trigger]) *model.PaginatedResult[dtos.TriggerResponse] {
	responses := make([]dtos.TriggerResponse, 0, len(result.Items))
	for _, trigger := range result.Items {
		response, err := mapper.EntityToDtoWithOptions[entities.Trigger, dtos.TriggerResponse](&trigger, mapperResponseOpts)
		if err != nil {
			continue
		}
		responses = append(responses, *response)
	}
	return &model.PaginatedResult[dtos.TriggerResponse]{
		Items: responses,
		Pagination: model.Pagination{
			Page:       result.Pagination.Page,
			PerPage:    result.Pagination.PerPage,
			TotalItems: result.Pagination.TotalItems,
			TotalPages: result.Pagination.TotalPages,
		},
	}
}

// triggerOrgIdFromContext extracts the orgId string from the request
// context (or "" when absent) for the per-org counter cache key.
func (s *TriggerService) triggerOrgIdFromContext(rc *reqCtx.RequestContext) string {
	if rc.OrgContext != nil {
		return *rc.OrgContext
	}
	return ""
}

// tryCachedTriggerCount returns (count, true) on a Redis hit; (0, false)
// when the caller must fall back to the repository.
func (s *TriggerService) tryCachedTriggerCount(c ctx.Context, cacheKey string) (int64, bool) {
	var count int64
	if err := s.deps.AppCache.Get(c, cacheKey, &count); err == nil {
		return count, true
	}
	return 0, false
}

// countTriggersFromRepo runs the org-scoped CountDocuments query.
func (s *TriggerService) countTriggersFromRepo(c ctx.Context, rc *reqCtx.RequestContext) (int64, error) {
	logger.Debug(fmt.Sprintf("[SERVICE:Trigger] Counter cache miss for orgId=%s", s.triggerOrgIdFromContext(rc)))
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return 0, err
	}
	return s.deps.TriggerRepository.CountDocuments(c, orgFilter)
}

// cacheTriggerCount stores a fresh count under the per-org cache key with
// the standard counter TTL. Best-effort; failures are intentionally ignored.
func (s *TriggerService) cacheTriggerCount(c ctx.Context, cacheKey string, count int64) {
	_ = s.deps.AppCache.SetEx(c, cacheKey, count, constants.CounterCacheTTL)
}
