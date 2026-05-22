package services

import (
	ctx "context"
	"time"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// applyTemplateScope decides between system / template-ancestor / org-local
// scope and applies the corresponding multi-tenant fields to the create DTO.
// Returns a 4xx custom error when the caller cannot create the requested
// scope (e.g. non-vendor trying to publish a template).
func (s *AssetTemplateService) applyTemplateScope(rc *reqCtx.RequestContext, dto *dtos.AssetTemplateCreateDTO) error {
	if dto.IsSystem {
		dto.OrgID = nil
		dto.PathKey = nil
		return nil
	}
	if dto.IsTemplate {
		if err := orgfilter.ValidateTemplateCreation(rc.OrgContextData.PathKey); err != nil {
			return &customErrors.ServerCustomError{Code: httpStatus.FORBIDDEN, Errors: []string{err.Error()}}
		}
		s.populateOrgContext(rc, dto)
		return nil
	}
	if err := orgfilter.ValidateOrgContextForNonSystem(rc); err != nil {
		return &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	s.populateOrgContext(rc, dto)
	return nil
}

// buildTemplateEntity maps the create DTO to a domain entity, converts
// classification ids to ObjectIDs, and seeds FieldId/Status for each
// dynamic field (EVA pattern: FieldId=1..N, Status=1, NextFieldId=N+1).
func (s *AssetTemplateService) buildTemplateEntity(dto *dtos.AssetTemplateCreateDTO) *entities.Assettemplate {
	entity, _ := mapper.DtoToEntity[dtos.AssetTemplateCreateDTO, entities.Assettemplate](dto)

	if dto.CategoryId != nil && *dto.CategoryId != "" {
		if id, err := model.ToObjectID(*dto.CategoryId); err == nil {
			entity.CategoryId = &id
		}
	}
	if dto.ManufacturerId != nil && *dto.ManufacturerId != "" {
		if id, err := model.ToObjectID(*dto.ManufacturerId); err == nil {
			entity.ManufacturerId = &id
		}
	}
	if dto.ModelId != nil && *dto.ModelId != "" {
		if id, err := model.ToObjectID(*dto.ModelId); err == nil {
			entity.ModelId = &id
		}
	}

	if len(entity.DynamicFields) > 0 {
		var nextFieldId uint16 = 1
		for i := range entity.DynamicFields {
			entity.DynamicFields[i].FieldId = nextFieldId
			entity.DynamicFields[i].Status = 1
			nextFieldId++
		}
		entity.NextFieldId = nextFieldId
	}
	return entity
}

// fanoutTemplateCreate publishes the post-create side effects: write
// scripts to MinIO, FANOUT cache invalidation, drop the org counter key.
func (s *AssetTemplateService) fanoutTemplateCreate(c ctx.Context, rc *reqCtx.RequestContext, template *entities.Assettemplate) {
	s.writeScripts(c, template)
	s.publishTemplateInvalidate(c, template)
	if rc.OrgContext != nil {
		counterKey := s.deps.CacheKeyBuilder.BuildCounterCacheKey(*rc.OrgContext)
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}

// buildTemplateUpdate translates the update DTO into a Mongo $set map.
// When DynamicFields is part of the patch, runs the EVA preservation
// helper so existing fields keep their FieldId and removed ones get
// Status=0. Always updates the timestamp and converts string ids.
func (s *AssetTemplateService) buildTemplateUpdate(existing *entities.Assettemplate, dto *dtos.AssetTemplateUpdateDTO) map[string]interface{} {
	if dto.DynamicFields != nil {
		processed := s.processDynamicFieldsUpdate(existing, dto.DynamicFields)
		dto.DynamicFields = nil
		fields, _ := mapper.DtoToMap(dto)
		fields["dynamicFields"] = processed.Fields
		fields["nextFieldId"] = processed.NextFieldId
		s.convertIdFieldsInMap(fields)
		fields["updated"] = time.Now()
		return fields
	}
	fields, _ := mapper.DtoToMap(dto)
	s.convertIdFieldsInMap(fields)
	fields["updated"] = time.Now()
	return fields
}

// fanoutTemplateUpdate publishes the post-update side effects: rewrite
// scripts to MinIO and broadcast FANOUT cache invalidation.
func (s *AssetTemplateService) fanoutTemplateUpdate(c ctx.Context, template *entities.Assettemplate) {
	s.writeScripts(c, template)
	s.publishTemplateInvalidate(c, template)
}

// fanoutTemplateDelete tears down all template-scoped caches: delete the
// MinIO scripts blob, broadcast FANOUT invalidation, and drop the org
// counter cache key.
func (s *AssetTemplateService) fanoutTemplateDelete(c ctx.Context, template *entities.Assettemplate) {
	s.deleteScripts(c, template)
	s.publishTemplateInvalidate(c, template)
	if template.OrgID != nil && !template.OrgID.IsZero() {
		counterKey := s.deps.CacheKeyBuilder.BuildCounterCacheKey(template.OrgID.Hex())
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}
