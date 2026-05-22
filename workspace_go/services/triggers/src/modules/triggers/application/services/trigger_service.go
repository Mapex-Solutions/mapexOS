package services

import (
	ctx "context"
	"fmt"

	"triggers/src/modules/triggers/application/constants"
	"triggers/src/modules/triggers/application/di"
	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/application/ports"
	"triggers/src/modules/triggers/domain/entities"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// Compile-time check to ensure TriggerService implements TriggerServicePort.
var _ ports.TriggerServicePort = (*TriggerService)(nil)

// New creates and returns a new instance of TriggerService.
func New(deps di.TriggerServiceDependenciesInjection) ports.TriggerServicePort {
	return &TriggerService{deps: deps}
}

// CreateTrigger orchestrates trigger creation: map DTO to entity -> apply
// multi-tenant scope (system / template / org-local) -> persist -> warm
// cache -> invalidate per-org counter cache -> return response DTO.
func (s *TriggerService) CreateTrigger(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateTriggerDto) (*dtos.TriggerResponse, error) {
	entity, err := s.buildTriggerEntity(dto)
	if err != nil {
		return nil, err
	}
	if err := s.applyTriggerScope(requestContext, dto, entity); err != nil {
		return nil, err
	}
	created, err := s.persistTrigger(c, entity)
	if err != nil {
		return nil, err
	}
	s.warmTriggerCacheAfterCreate(c, requestContext, created)
	return s.mapTriggerResponse(created)
}

// GetTriggerById fetches a trigger by id with cache-aside semantics:
// GetOrSetEx hits Redis first and falls back to the repository on miss.
// CacheMetrics is propagated when the caller passes one (used by event
// processing to track per-call hit/miss).
func (s *TriggerService) GetTriggerById(c ctx.Context, triggerId *string, metrics ...*common.CacheMetrics) (*dtos.TriggerResponse, error) {
	var trigger entities.Trigger
	if err := s.loadTriggerCacheAside(c, triggerId, &trigger, metrics...); err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   status.NOT_FOUND,
			Errors: []string{fmt.Sprintf("Trigger with ID %s not found", *triggerId)},
		}
	}
	return s.mapTriggerResponse(&trigger)
}

// GetTriggers returns the paginated, filtered trigger list scoped to the
// caller's org context (own org + system templates + ancestor templates).
func (s *TriggerService) GetTriggers(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.TriggerQueryDto) (*model.PaginatedResult[dtos.TriggerResponse], error) {
	filters := s.buildTriggerListFilters(requestContext, query)
	pagination := s.buildTriggerListPagination(query)
	result, err := s.deps.TriggerRepository.FindWithFilters(c, filters, pagination, nil)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}
	return s.mapTriggerListResult(result), nil
}

// UpdateTriggerById applies a partial update by id and returns the
// resulting response DTO. 404 when the id is unknown.
func (s *TriggerService) UpdateTriggerById(c ctx.Context, requestContext *reqCtx.RequestContext, triggerId *string, dto *dtos.UpdateTriggerDto) (*dtos.TriggerResponse, error) {
	payload, err := s.buildTriggerUpdatePayload(dto)
	if err != nil {
		return nil, err
	}
	updated, err := s.deps.TriggerRepository.FindByIdAndUpdate(c, triggerId, payload)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{fmt.Sprintf("Trigger with ID %s not found", *triggerId)}}
	}
	s.invalidateTriggerCache(c, *triggerId)
	response, err := mapper.EntityToDtoWithOptions[entities.Trigger, dtos.TriggerResponse](updated, mapperResponseOpts)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}
	return response, nil
}

// DeleteTriggerById removes a trigger and invalidates both per-id and
// per-org counter caches. 404 when the id is unknown.
func (s *TriggerService) DeleteTriggerById(c ctx.Context, triggerId *string) (map[string]bool, error) {
	trigger, _ := s.deps.TriggerRepository.FindById(c, triggerId)
	if err := s.deps.TriggerRepository.DeleteById(c, triggerId); err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{fmt.Sprintf("Trigger with ID %s not found", *triggerId)}}
	}
	s.invalidateTriggerCache(c, *triggerId)
	s.invalidateTriggerCounterCache(c, trigger)
	return map[string]bool{"success": true}, nil
}

// CountTriggers returns the per-org trigger count with cache-aside
// semantics (6h TTL); the counter cache is invalidated on Create/Delete.
func (s *TriggerService) CountTriggers(c ctx.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	orgId := s.triggerOrgIdFromContext(requestContext)
	cacheKey := s.deps.CacheKeyBuilder.CounterKey(orgId)
	if count, ok := s.tryCachedTriggerCount(c, cacheKey); ok {
		return count, nil
	}
	count, err := s.countTriggersFromRepo(c, requestContext)
	if err != nil {
		return 0, err
	}
	s.cacheTriggerCount(c, cacheKey, count)
	return count, nil
}

// silence unused-import warning when generic helpers are pruned in handlers.
var _ = constants.TriggerCacheTTL
