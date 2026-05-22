package services

import (
	ctx "context"
	"time"

	"http_gateway/src/modules/datasources/application/constants"
	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// recordDataSourceOperation records Prometheus counters and latency histogram for a DataSource operation.
// op is the logical operation (list/create/read/update/delete); outcome is success/error; start is the timer origin.
func (s *DataSourceService) recordDataSourceOperation(op, outcome string, start time.Time) {
	s.deps.Metrics.DsOperations.WithLabelValues(op, outcome).Inc()
	s.deps.Metrics.DsOperationDuration.WithLabelValues(op).Observe(time.Since(start).Seconds())
}

// recordDataSourceListSuccess emits the per-page result-count histogram on top
// of the standard list-success metrics so the orchestration in service.go
// stays a single named-step call.
func (s *DataSourceService) recordDataSourceListSuccess(start time.Time, dtoResult *model.PaginatedResult[dtos.DataSourceResponse]) {
	s.recordDataSourceOperation("list", "success", start)
	s.deps.Metrics.DsListResultsCount.Observe(float64(len(dtoResult.Items)))
}

// cacheDataSourceEntity stores the given DataSource entity in Redis using the module's cache key builder.
func (s *DataSourceService) cacheDataSourceEntity(c ctx.Context, id string, entity *entities.DataSource) {
	cacheKey := s.deps.CacheKeyBuilder.BuildKey(id)
	s.deps.CacheRepo.SetEx(c, cacheKey, entity, constants.DataSourceCacheTTL)
}

// mapDataSourceResponse converts a DataSource entity to its response DTO.
func (s *DataSourceService) mapDataSourceResponse(entity *entities.DataSource) (*dtos.DataSourceResponse, error) {
	return mapper.EntityToDto[entities.DataSource, dtos.DataSourceResponse](entity)
}
