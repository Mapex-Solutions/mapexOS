package services

import (
	ctx "context"
	"fmt"

	"assets/src/modules/assets/application/constants"
	"assets/src/modules/assets/application/dtos"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildListFilters assembles the asset-collection and template-collection
// filter maps from the request context + query DTO. Returns a non-nil
// error only when the org filter cannot be derived (e.g. missing context).
func (s *AssetService) buildListFilters(rc *reqCtx.RequestContext, query *dtos.AssetQueryDTO) (model.Map, model.Map, error) {
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: rc,
		Query:      query,
	})
	if err != nil {
		return nil, nil, err
	}

	assetFilters := orgFilter
	if query.Enabled != nil {
		assetFilters["enabled"] = *query.Enabled
	}
	if query.AssetTemplateID != nil && *query.AssetTemplateID != "" {
		if id, err := model.ToObjectID(*query.AssetTemplateID); err == nil {
			assetFilters["assetTemplateId"] = id
		}
	}
	if query.AssetUUID != nil && *query.AssetUUID != "" {
		assetFilters["assetUUID"] = *query.AssetUUID
	}
	if query.Name != nil && *query.Name != "" {
		assetFilters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.HealthStatus != nil && *query.HealthStatus != "" {
		assetFilters["healthStatus"] = *query.HealthStatus
	}

	templateFilters := model.Map{}
	if query.CategoryId != nil && *query.CategoryId != "" {
		if id, err := model.ToObjectID(*query.CategoryId); err == nil {
			templateFilters["categoryId"] = id
		}
	}
	if query.ManufacturerId != nil && *query.ManufacturerId != "" {
		if id, err := model.ToObjectID(*query.ManufacturerId); err == nil {
			templateFilters["manufacturerId"] = id
		}
	}
	if query.ModelId != nil && *query.ModelId != "" {
		if id, err := model.ToObjectID(*query.ModelId); err == nil {
			templateFilters["modelId"] = id
		}
	}
	return assetFilters, templateFilters, nil
}

// tryCachedCount reads the cached counter for the active org. Returns
// (count, true) only on a fresh hit; cache misses fall through silently.
func (s *AssetService) tryCachedCount(c ctx.Context, cacheKey string) (int64, bool) {
	var count int64
	if err := s.deps.AppCache.Get(c, cacheKey, &count); err == nil {
		s.deps.Metrics.AssetCacheTotal.WithLabelValues("hit").Inc()
		return count, true
	}
	s.deps.Metrics.AssetCacheTotal.WithLabelValues("miss").Inc()
	logger.Debug(fmt.Sprintf("[SERVICE:Asset] Counter cache miss for key=%s", cacheKey))
	return 0, false
}

// countFromRepo computes the count via Mongo CountDocuments using the
// org filter derived from the request context.
func (s *AssetService) countFromRepo(c ctx.Context, rc *reqCtx.RequestContext) (int64, error) {
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return 0, err
	}
	return s.deps.AssetRepo.CountDocuments(c, orgFilter)
}

// cacheCount writes the freshly computed count back to Redis with the
// standard TTL. Best-effort — failures do not bubble up.
func (s *AssetService) cacheCount(c ctx.Context, cacheKey string, count int64) {
	_ = s.deps.AppCache.SetEx(c, cacheKey, count, constants.CounterCacheTTL)
}
