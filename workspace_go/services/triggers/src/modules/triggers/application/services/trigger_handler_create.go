package services

import (
	ctx "context"
	"time"

	"triggers/src/modules/triggers/application/constants"
	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildTriggerEntity maps the create DTO into a domain entity and stamps the
// timestamp fields. Multi-tenant scope is applied separately in
// applyTriggerScope so the orchestration can show both steps.
func (s *TriggerService) buildTriggerEntity(dto *dtos.CreateTriggerDto) (*entities.Trigger, error) {
	entity, err := mapper.DtoToEntity[dtos.CreateTriggerDto, entities.Trigger](dto)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	now := time.Now()
	entity.Created = now
	entity.Updated = now
	return entity, nil
}

// applyTriggerScope enforces multi-tenant scope rules on a Create request
// and stamps OrgID/PathKey on the entity:
//   - System resources clear org context (forced isTemplate=false)
//   - Templates require Vendor/Customer org context
//   - Local resources require a non-empty org context
func (s *TriggerService) applyTriggerScope(rc *reqCtx.RequestContext, dto *dtos.CreateTriggerDto, entity *entities.Trigger) error {
	if dto.IsSystem {
		entity.OrgID = nil
		entity.PathKey = ""
		entity.IsSystem = true
		entity.IsTemplate = false
		return nil
	}
	if err := orgfilter.ValidateOrgContextForNonSystem(rc); err != nil {
		return &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	if dto.IsTemplate {
		if err := orgfilter.ValidateTemplateCreation(rc.OrgContextData.PathKey); err != nil {
			return &customErrors.ServerCustomError{Code: status.FORBIDDEN, Errors: []string{err.Error()}}
		}
	}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		orgObjectId, _ := model.ToObjectID(*rc.OrgContext)
		entity.OrgID = &orgObjectId
	}
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		entity.PathKey = rc.OrgContextData.PathKey
	}
	entity.IsSystem = false
	entity.IsTemplate = dto.IsTemplate
	return nil
}

// persistTrigger writes the entity to the repository and surfaces a 500
// contract error on failure so the public method stays a thin orchestration.
func (s *TriggerService) persistTrigger(c ctx.Context, entity *entities.Trigger) (*entities.Trigger, error) {
	created, err := s.deps.TriggerRepository.Create(c, entity)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}
	return created, nil
}

// warmTriggerCacheAfterCreate stores the new entity under its per-id cache
// key and drops the per-org counter cache so the next CountTriggers hits
// the repository.
func (s *TriggerService) warmTriggerCacheAfterCreate(c ctx.Context, rc *reqCtx.RequestContext, entity *entities.Trigger) {
	cacheKey := s.deps.CacheKeyBuilder.TriggerKey(entity.ID.Hex())
	s.deps.CacheRepository.SetEx(c, cacheKey, entity, constants.TriggerCacheTTL)
	if rc.OrgContext != nil {
		counterKey := s.deps.CacheKeyBuilder.CounterKey(*rc.OrgContext)
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}

// mapTriggerResponse projects a Trigger entity into its response DTO,
// translating mapper failures into the canonical 500 contract error.
func (s *TriggerService) mapTriggerResponse(entity *entities.Trigger) (*dtos.TriggerResponse, error) {
	response, err := mapper.EntityToDtoWithOptions[entities.Trigger, dtos.TriggerResponse](entity, mapperResponseOpts)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}
	return response, nil
}
