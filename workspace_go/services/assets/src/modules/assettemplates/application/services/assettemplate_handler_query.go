package services

import (
	ctx "context"
	"fmt"

	"assets/src/modules/assettemplates/application/constants"
	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildTemplateListFilters assembles the Mongo filter map for the
// template list. Templates are visible when ANY of (org-scope, system,
// ancestor) match — encoded via $or — minus the explicit isSystem=false
// or isTemplate=false exclusions, plus the optional name/manufacturer/
// model filters.
func (s *AssetTemplateService) buildTemplateListFilters(rc *reqCtx.RequestContext, query *dtos.AssetTemplateQueryDto) model.Map {
	orConditions := []model.Map{}
	if orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query}); len(orgFilter) > 0 {
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
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.ManufacturerId != nil && *query.ManufacturerId != "" {
		if id, err := model.ToObjectID(*query.ManufacturerId); err == nil {
			filters["manufacturerId"] = id
		}
	}
	if query.ModelId != nil && *query.ModelId != "" {
		if id, err := model.ToObjectID(*query.ModelId); err == nil {
			filters["modelId"] = id
		}
	}
	return filters
}

// buildTemplateProjection wraps orgfilter.BuildProjection so the service
// signature stays focused on its inputs.
func (s *AssetTemplateService) buildTemplateProjection(query *dtos.AssetTemplateQueryDto) model.Map {
	return orgfilter.BuildProjection(query.Projection)
}

// mapTemplateEntitiesToDtos converts each entity to a response DTO.
func (s *AssetTemplateService) mapTemplateEntitiesToDtos(items []entities.Assettemplate) []dtos.AssetTemplateResponse {
	out := make([]dtos.AssetTemplateResponse, len(items))
	for i := range items {
		entity := items[i]
		dto, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](&entity)
		out[i] = *dto
	}
	return out
}

// tryCachedTemplateCount looks up the cached counter; returns
// (count, true) on hit and tracks the cache_total metric.
func (s *AssetTemplateService) tryCachedTemplateCount(c ctx.Context, cacheKey string) (int64, bool) {
	var count int64
	if err := s.deps.AppCache.Get(c, cacheKey, &count); err == nil {
		s.deps.Metrics.TemplateCacheTotal.WithLabelValues("hit").Inc()
		return count, true
	}
	s.deps.Metrics.TemplateCacheTotal.WithLabelValues("miss").Inc()
	logger.Debug(fmt.Sprintf("[SERVICE:AssetTemplate] Counter cache miss for key=%s", cacheKey))
	return 0, false
}

// countTemplatesFromRepo counts via Mongo using the same filter shape as
// the list path so the user-visible counter and the result list stay
// consistent (org-scope / system / ancestor union).
func (s *AssetTemplateService) countTemplatesFromRepo(c ctx.Context, rc *reqCtx.RequestContext) (int64, error) {
	orConditions := []model.Map{}
	if orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc}); len(orgFilter) > 0 {
		orConditions = append(orConditions, orgFilter)
	}
	orConditions = append(orConditions, model.Map{"isSystem": true})

	filters := model.Map{}
	if len(orConditions) > 0 {
		filters["$or"] = orConditions
	}
	return s.deps.AssetTemplateRepo.CountDocuments(c, filters)
}

// cacheTemplateCount writes the freshly computed count back to Redis
// with the standard TTL. Best-effort.
func (s *AssetTemplateService) cacheTemplateCount(c ctx.Context, cacheKey string, count int64) {
	_ = s.deps.AppCache.SetEx(c, cacheKey, count, constants.CounterCacheTTL)
}
