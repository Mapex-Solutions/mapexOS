package services

import (
	ctx "context"

	"http_gateway/src/modules/datasources/application/constants"
	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildDataSourceFilters assembles the MongoDB filter map for the list query.
// It combines the org filter derived from the RequestContext with module-specific
// filters (name, enabled, mode, protocol) from the query DTO.
func (s *DataSourceService) buildDataSourceFilters(
	requestContext *reqCtx.RequestContext,
	query *dtos.DataSourceQueryDTO,
) (model.Map, error) {
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
		Query:      query,
	})
	if err != nil {
		return nil, err
	}

	filters := model.Map{}
	for k, v := range orgFilter {
		filters[k] = v
	}

	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Mode != nil && *query.Mode != "" {
		filters["mode"] = *query.Mode
	}
	if query.Protocol != nil && *query.Protocol != "" {
		filters["protocol"] = *query.Protocol
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}

	return filters, nil
}

// buildDataSourceListOptions builds pagination options and the projection map from the query DTO.
func (s *DataSourceService) buildDataSourceListOptions(query *dtos.DataSourceQueryDTO) (*model.PaginationOpts, model.Map) {
	pagination := &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
	projection := orgfilter.BuildProjection(query.Projection)
	return pagination, projection
}

// mapDataSourceListResult converts a paginated entity result into a paginated DTO result.
func (s *DataSourceService) mapDataSourceListResult(
	result *model.PaginatedResult[entities.DataSource],
) *model.PaginatedResult[dtos.DataSourceResponse] {
	dtoItems := make([]dtos.DataSourceResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.DataSource, dtos.DataSourceResponse](&entity)
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.DataSourceResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}
}

// loadDataSourceCache reads a DataSource via cache-aside: returns cached value when present,
// otherwise invokes the callback to hydrate from the repository and populate the cache.
// The returned boolean is true when the cache was missed (callback executed).
func (s *DataSourceService) loadDataSourceCache(c ctx.Context, dataSourceId *string) (*entities.DataSource, bool, error) {
	cacheKey := s.deps.CacheKeyBuilder.BuildKey(*dataSourceId)

	var entity entities.DataSource
	cacheMiss := false

	_, err := s.deps.CacheRepo.GetOrSetEx(common.GetOrSetParams{
		Ctx:      c,
		CacheKey: cacheKey,
		CacheTTL: int(constants.DataSourceCacheTTL.Seconds()),
		Dest:     &entity,
		Callback: func() (interface{}, error) {
			cacheMiss = true
			return s.deps.DataSourceRepo.FindById(c, dataSourceId)
		},
	})

	return &entity, cacheMiss, err
}

// recordDataSourceCacheOutcome records the cache hit/miss metric for a read.
func (s *DataSourceService) recordDataSourceCacheOutcome(cacheMiss bool) {
	if cacheMiss {
		s.deps.Metrics.DsCacheTotal.WithLabelValues("miss").Inc()
		return
	}
	s.deps.Metrics.DsCacheTotal.WithLabelValues("hit").Inc()
}

// dataSourceNotFoundError returns the standardized NOT_FOUND error used across query flows.
func (s *DataSourceService) dataSourceNotFoundError() error {
	return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Data source not found"}}
}

// dataSourceMappingError returns the standardized mapping failure error used on DTO conversion.
func (s *DataSourceService) dataSourceMappingError() error {
	return &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"failed to map data source response"}}
}
