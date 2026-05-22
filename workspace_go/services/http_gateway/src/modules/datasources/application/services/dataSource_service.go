package services

import (
	ctx "context"
	"time"

	"http_gateway/src/modules/datasources/application/di"
	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// Compile-time check to ensure DataSourceService implements DataSourceServicePort interface.
// This will cause a compilation error if the interface is not fully implemented.
var _ ports.DataSourceServicePort = (*DataSourceService)(nil)

// New creates and returns a new instance of DataSourceService.
//
// This constructor follows Hexagonal Architecture by:
//   - Accepting dependencies through a DI struct (single parameter pattern)
//   - Returning the service port interface (not concrete type)
//   - Enabling loose coupling and testability
//
// Parameters:
//   - deps: Aggregated dependencies (repositories, NATS bus) injected by dig
//
// Returns:
//   - DataSourceServicePort: The service port interface implementation
func New(deps di.DataSourceServiceDependenciesInjection) ports.DataSourceServicePort {
	return &DataSourceService{
		deps: deps,
	}
}

// GetDataSources retrieves a paginated and filtered list of data sources.
// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
func (s *DataSourceService) GetDataSources(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.DataSourceQueryDTO) (*model.PaginatedResult[dtos.DataSourceResponse], error) {
	start := time.Now()
	filters, err := s.buildDataSourceFilters(requestContext, query)
	if err != nil {
		s.recordDataSourceOperation("list", "error", start)
		return nil, err
	}
	pagination, projection := s.buildDataSourceListOptions(query)
	result, err := s.deps.DataSourceRepo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		s.recordDataSourceOperation("list", "error", start)
		return nil, err
	}
	dtoResult := s.mapDataSourceListResult(result)
	s.recordDataSourceListSuccess(start, dtoResult)
	return dtoResult, nil
}

// CreateDataSource creates a new DataSource entity from the provided CreateDataSourceDTO
// and persists it using the DataSourceRepository.
func (s *DataSourceService) CreateDataSource(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.DataSourceCreateDTO) (*dtos.DataSourceResponse, error) {
	start := time.Now()
	s.applyCreateRequestContext(requestContext, dto)
	entity, err := s.persistNewDataSource(c, dto)
	if err != nil {
		s.recordDataSourceOperation("create", "error", start)
		return nil, err
	}
	s.cacheDataSourceEntity(c, entity.ID.Hex(), entity)
	response, _ := s.mapDataSourceResponse(entity)
	s.recordDataSourceOperation("create", "success", start)
	return response, nil
}

// GetDataSourceById retrieves a DataSource entity by its unique identifier.
// It uses cache-aside pattern: checks Redis cache first, then DB if cache miss.
func (s *DataSourceService) GetDataSourceById(c ctx.Context, dataSourceId *string) (*dtos.DataSourceResponse, error) {
	start := time.Now()

	entity, cacheMiss, err := s.loadDataSourceCache(c, dataSourceId)
	s.recordDataSourceCacheOutcome(cacheMiss)

	if err != nil || entity.ID.IsZero() {
		s.recordDataSourceOperation("read", "error", start)
		return nil, s.dataSourceNotFoundError()
	}

	response, mapErr := s.mapDataSourceResponse(entity)
	if mapErr != nil {
		s.recordDataSourceOperation("read", "error", start)
		return nil, s.dataSourceMappingError()
	}

	s.recordDataSourceOperation("read", "success", start)
	return response, nil
}

// UpdateDataSourceById updates an existing DataSource entity identified by its unique ID
// with the provided data from the DataSourceUpdateDTO.
func (s *DataSourceService) UpdateDataSourceById(c ctx.Context, dataSourceId *string, dto *dtos.DataSourceUpdateDTO) (*dtos.DataSourceResponse, error) {
	start := time.Now()
	fieldsToUpdate := s.buildUpdateFields(dto)
	entity := s.applyDataSourceUpdate(c, dataSourceId, fieldsToUpdate)
	if entity.ID.IsZero() {
		s.recordDataSourceOperation("update", "error", start)
		return nil, s.dataSourceNotFoundError()
	}
	s.cacheDataSourceEntity(c, *dataSourceId, entity)
	response, _ := s.mapDataSourceResponse(entity)
	s.recordDataSourceOperation("update", "success", start)
	return response, nil
}

// DeleteDataSourceById removes a DataSource entity identified by its unique ID from the repository.
func (s *DataSourceService) DeleteDataSourceById(c ctx.Context, dataSourceId *string) (map[string]bool, error) {
	start := time.Now()

	if err := s.deps.DataSourceRepo.DeleteById(c, dataSourceId); err != nil {
		s.recordDataSourceOperation("delete", "error", start)
		return nil, err
	}

	s.deleteDataSourceCache(c, dataSourceId)
	s.recordDataSourceOperation("delete", "success", start)

	return map[string]bool{"success": true}, nil
}
