package services

import (
	ctx "context"
	"time"

	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/domain/entities"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// Compile-time check
var _ ports.AssetTemplateServicePort = (*AssetTemplateService)(nil)

// New creates and returns a new instance of AssetTemplateService.
func New(deps di.AssetTemplateServiceDependenciesInjection) ports.AssetTemplateServicePort {
	return &AssetTemplateService{
		deps: deps,
	}
}

// ProcessL2WriteRetry is the public entry point for the L2 sync
// fallback consumer (template_l2sync). Re-fetches the current
// template state from Mongo and re-runs syncTemplateL2. On success
// emits the existing FANOUT invalidation so caches downstream
// refresh. Returns an error if Mongo lookup fails — the consumer
// NAKs and NATS retries with backoff.
func (s *AssetTemplateService) ProcessL2WriteRetry(c ctx.Context, templateId string) error {
	t, err := s.fetchTemplateByID(c, templateId)
	if err != nil {
		return err
	}
	if t == nil {
		return nil // template deleted between failure and retry — drop
	}
	s.syncTemplateL2(c, t)
	s.publishTemplateInvalidate(c, t)
	return nil
}

// CreateAssetTemplate orchestrates template creation:
// resolve scope (system / template / org-local) and apply org context ->
// build the entity with classification ObjectIDs and EVA FieldIds ->
// persist -> fan out side effects (MinIO scripts, FANOUT, counter cache)
// -> return the response DTO.
func (s *AssetTemplateService) CreateAssetTemplate(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.AssetTemplateCreateDTO) (*dtos.AssetTemplateResponse, error) {
	start := time.Now()

	if err := s.applyTemplateScope(requestContext, dto); err != nil {
		s.recordTemplateOp("create", "error", start)
		return nil, err
	}
	entity := s.buildTemplateEntity(dto)
	persisted, err := s.deps.AssetTemplateRepo.Create(c, entity)
	if err != nil {
		s.recordTemplateOp("create", "error", start)
		return nil, err
	}
	s.fanoutTemplateCreate(c, requestContext, persisted)

	s.recordTemplateOp("create", "success", start)
	resp, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](persisted)
	return resp, nil
}

// GetAssetTemplateById fetches a template by id and returns its DTO.
// Returns 404 when the id is unknown.
func (s *AssetTemplateService) GetAssetTemplateById(c ctx.Context, templateId *string) (*dtos.AssetTemplateResponse, error) {
	start := time.Now()

	template, err := s.fetchTemplateById(c, templateId)
	if err != nil {
		s.recordTemplateOp("read", "error", start)
		return nil, err
	}
	s.recordTemplateOp("read", "success", start)
	resp, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](template)
	return resp, nil
}

// UpdateAssetTemplateById orchestrates a partial template update:
// load the prior entity (404 on miss) -> assemble the $set map (handles
// EVA dynamic-field FieldId preservation when DynamicFields is in the
// patch) -> apply -> fan out side effects (MinIO scripts, FANOUT) -> DTO.
func (s *AssetTemplateService) UpdateAssetTemplateById(c ctx.Context, templateId *string, dto *dtos.AssetTemplateUpdateDTO) (*dtos.AssetTemplateResponse, error) {
	start := time.Now()

	existing, err := s.fetchTemplateById(c, templateId)
	if err != nil {
		s.recordTemplateOp("update", "error", start)
		return nil, err
	}
	patch := s.buildTemplateUpdate(existing, dto)
	updated, _ := s.deps.AssetTemplateRepo.FindByIdAndUpdate(c, templateId, patch)
	if updated.ID.IsZero() {
		s.recordTemplateOp("update", "error", start)
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Asset Template not found"}}
	}
	s.fanoutTemplateUpdate(c, updated)

	s.recordTemplateOp("update", "success", start)
	resp, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](updated)
	return resp, nil
}

// DeleteAssetTemplateById orchestrates template deletion:
// load (404 on miss) -> delete from Mongo -> tear down MinIO scripts ->
// FANOUT invalidate -> drop the counter cache key.
func (s *AssetTemplateService) DeleteAssetTemplateById(c ctx.Context, templateId *string) (map[string]bool, error) {
	start := time.Now()

	template, err := s.fetchTemplateById(c, templateId)
	if err != nil {
		s.recordTemplateOp("delete", "error", start)
		return nil, err
	}
	if err := s.deps.AssetTemplateRepo.DeleteById(c, templateId); err != nil {
		s.recordTemplateOp("delete", "error", start)
		return nil, err
	}
	s.fanoutTemplateDelete(c, template)

	s.recordTemplateOp("delete", "success", start)
	return map[string]bool{"success": true}, nil
}

// GetAssetTemplates orchestrates the paginated list:
// build filter conditions ($or over org / system / ancestor templates) ->
// run the repository query -> map entities to DTOs.
func (s *AssetTemplateService) GetAssetTemplates(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.AssetTemplateQueryDto) (*model.PaginatedResult[dtos.AssetTemplateResponse], error) {
	start := time.Now()

	filters := s.buildTemplateListFilters(requestContext, query)
	pagination := &model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}
	projection := s.buildTemplateProjection(query)

	result, err := s.deps.AssetTemplateRepo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		s.recordTemplateOp("list", "error", start)
		return nil, err
	}

	dtoItems := s.mapTemplateEntitiesToDtos(result.Items)
	s.recordTemplateOp("list", "success", start)
	s.deps.Metrics.TemplateListResultsCount.Observe(float64(len(dtoItems)))
	return &model.PaginatedResult[dtos.AssetTemplateResponse]{Items: dtoItems, Pagination: result.Pagination}, nil
}

// CountAssetTemplates returns the total count for the caller's org
// context with cache-aside semantics: try Redis -> fallback to Mongo
// CountDocuments using the list-mode filter -> re-cache on miss.
func (s *AssetTemplateService) CountAssetTemplates(c ctx.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	start := time.Now()

	orgId := ""
	if requestContext.OrgContext != nil {
		orgId = *requestContext.OrgContext
	}
	cacheKey := s.deps.CacheKeyBuilder.BuildCounterCacheKey(orgId)

	if count, ok := s.tryCachedTemplateCount(c, cacheKey); ok {
		s.recordTemplateOp("count", "success", start)
		return count, nil
	}

	count, err := s.countTemplatesFromRepo(c, requestContext)
	if err != nil {
		s.recordTemplateOp("count", "error", start)
		return 0, err
	}
	s.cacheTemplateCount(c, cacheKey, count)
	s.recordTemplateOp("count", "success", start)
	return count, nil
}

// UpdateManufacturerName denormalizes the new manufacturer name across
// every template that references the given manufacturer list id.
func (s *AssetTemplateService) UpdateManufacturerName(c ctx.Context, manufacturerId string, newName string) error {
	start := time.Now()
	objectId, err := s.parseClassificationId(manufacturerId, "manufacturer")
	if err != nil {
		s.recordTemplateOp("update_manufacturer", "error", start)
		return err
	}
	matched, err := s.denormalizeTemplateField(c, "manufacturerId", objectId, "manufacturerName", newName)
	if err != nil {
		s.recordTemplateOp("update_manufacturer", "error", start)
		return err
	}
	s.recordTemplateOp("update_manufacturer", "success", start)
	s.logClassificationDenorm("manufacturer", matched)
	return nil
}

// UpdateModelName denormalizes the new model name across every template
// that references the given model list id.
func (s *AssetTemplateService) UpdateModelName(c ctx.Context, modelId string, newName string) error {
	start := time.Now()
	objectId, err := s.parseClassificationId(modelId, "model")
	if err != nil {
		s.recordTemplateOp("update_model", "error", start)
		return err
	}
	matched, err := s.denormalizeTemplateField(c, "modelId", objectId, "modelName", newName)
	if err != nil {
		s.recordTemplateOp("update_model", "error", start)
		return err
	}
	s.recordTemplateOp("update_model", "success", start)
	s.logClassificationDenorm("model", matched)
	return nil
}

// UpdateCategoryName denormalizes the new category name across every
// template that references the given category list id.
func (s *AssetTemplateService) UpdateCategoryName(c ctx.Context, categoryId string, newName string) error {
	start := time.Now()
	objectId, err := s.parseClassificationId(categoryId, "category")
	if err != nil {
		s.recordTemplateOp("update_category", "error", start)
		return err
	}
	matched, err := s.denormalizeTemplateField(c, "categoryId", objectId, "categoryName", newName)
	if err != nil {
		s.recordTemplateOp("update_category", "error", start)
		return err
	}
	s.recordTemplateOp("update_category", "success", start)
	s.logClassificationDenorm("category", matched)
	return nil
}

// HandleListNameUpdated decodes a list-name-updated NATS message and
// dispatches to the matching denormalization handler. Owns Ack/Nack/
// Reject internally so the consumer file stays pure wiring.
func (s *AssetTemplateService) HandleListNameUpdated(msg *natsModel.Message) {
	event, ok := s.parseListNameUpdatedEvent(msg)
	if !ok {
		return
	}
	msg.OrgId = event.OrgId
	s.logListNameUpdatedReceived(event)
	skipped, err := s.dispatchListNameUpdate(ctx.Background(), event)
	s.completeListNameUpdated(msg, event, skipped, err)
}

// GetAvailableFields returns only the AvailableFields slice for one
// template id. Used by the field-source picker.
func (s *AssetTemplateService) GetAvailableFields(c ctx.Context, templateId *string, requestContext *reqCtx.RequestContext) (map[string]interface{}, error) {
	template, err := s.deps.AssetTemplateRepo.FindById(c, templateId)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Asset Template not found"}}
	}
	return map[string]interface{}{"availableFields": template.AvailableFields}, nil
}

// GetTemplateByIdForCacheFallback fetches a template and repopulates the
// MinIO L2 scripts cache before returning the DTO. Used by the internal
// fallback endpoint when consuming services miss in their TieredCache.
func (s *AssetTemplateService) GetTemplateByIdForCacheFallback(c ctx.Context, templateId string) (*dtos.AssetTemplateResponse, error) {
	template, err := s.fetchTemplateById(c, &templateId)
	if err != nil {
		return nil, err
	}
	s.writeScripts(c, template)
	resp, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](template)
	return resp, nil
}

