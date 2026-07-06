package services

import (
	ctx "context"
	"maps"
	"math"
	"time"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	orgfilter "github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// defaultMigrationBatchSize is the fallback batch size when neither the request
// nor config supplies one.
const defaultMigrationBatchSize = 100

// resolveMigrationOrgID extracts the request's organization id, required to
// scope the plan.
func (s *AssetTemplateService) resolveMigrationOrgID(rc *reqCtx.RequestContext) (model.ObjectId, error) {
	if rc == nil || rc.OrgContext == nil || *rc.OrgContext == "" {
		return model.ObjectId{}, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"organization context is required"}}
	}
	orgID, err := model.ToObjectID(*rc.OrgContext)
	if err != nil {
		return model.ObjectId{}, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"invalid organization context"}}
	}
	return orgID, nil
}

// resolveMigrationTemplateNames asserts both templates exist (fetchTemplateById
// reports 404) and returns their names, denormalized onto the plan so the list
// shows names without a per-row lookup.
func (s *AssetTemplateService) resolveMigrationTemplateNames(c ctx.Context, dto *dtos.MigrationPlanCreateRequest) (string, string, error) {
	from, err := s.fetchTemplateById(c, &dto.FromTemplateID)
	if err != nil {
		return "", "", err
	}
	to, err := s.fetchTemplateById(c, &dto.ToTemplateID)
	if err != nil {
		return "", "", err
	}
	return from.Name, to.Name, nil
}

// buildMigrationPlan assembles the scheduled plan and one pending execution per
// asset, converting every string id to an ObjectId.
func (s *AssetTemplateService) buildMigrationPlan(orgID model.ObjectId, dto *dtos.MigrationPlanCreateRequest, fromName, toName string) (*entities.MigrationPlan, []*entities.MigrationExecution, error) {
	fromID, err := model.ToObjectID(dto.FromTemplateID)
	if err != nil {
		return nil, nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"invalid fromTemplateId"}}
	}
	toID, err := model.ToObjectID(dto.ToTemplateID)
	if err != nil {
		return nil, nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"invalid toTemplateId"}}
	}
	assetIDs := make([]model.ObjectId, 0, len(dto.AssetIDs))
	for _, aid := range dto.AssetIDs {
		oid, err := model.ToObjectID(aid)
		if err != nil {
			return nil, nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"invalid assetId: " + aid}}
		}
		assetIDs = append(assetIDs, oid)
	}

	id := model.NewObjectID()
	now := time.Now()
	plan := &entities.MigrationPlan{
		ID:             id,
		OrgID:          orgID,
		Name:             dto.Name,
		Description:      dto.Description,
		FromTemplateID:   fromID,
		FromTemplateName: fromName,
		ToTemplateID:     toID,
		ToTemplateName:   toName,
		AssetIDs:         assetIDs,
		ScheduleAt:     s.resolveMigrationScheduleAt(dto.ScheduleAt),
		BatchSize:      s.resolveMigrationBatchSize(dto.BatchSize),
		Status:         entities.PlanScheduled,
		Counters:       entities.PlanCounters{Total: len(assetIDs)},
		Created:        now,
		Updated:        now,
	}

	execs := make([]*entities.MigrationExecution, 0, len(assetIDs))
	for _, assetID := range assetIDs {
		execs = append(execs, &entities.MigrationExecution{
			ID:      model.NewObjectID(),
			PlanID:  id,
			AssetID: assetID,
			OrgID:   orgID,
			Status:  entities.ExecPending,
			Created: now,
			Updated: now,
		})
	}
	return plan, execs, nil
}

// resolveMigrationBatchSize prefers the request value, then the configured
// default, then a hard fallback.
func (s *AssetTemplateService) resolveMigrationBatchSize(requested *int) int {
	if requested != nil && *requested > 0 {
		return *requested
	}
	if def, err := config.GetIntValue("template_migration_batch_size"); err == nil && def > 0 {
		return def
	}
	return defaultMigrationBatchSize
}

// resolveMigrationScheduleAt runs immediately when the requested time is nil or
// already in the past.
func (s *AssetTemplateService) resolveMigrationScheduleAt(requested *time.Time) time.Time {
	now := time.Now()
	if requested != nil && requested.After(now) {
		return *requested
	}
	return now
}

// persistMigrationPlan writes the plan aggregate and bulk-inserts its
// executions.
func (s *AssetTemplateService) persistMigrationPlan(c ctx.Context, plan *entities.MigrationPlan, execs []*entities.MigrationExecution) error {
	if _, err := s.deps.MigrationPlanRepo.Create(c, plan); err != nil {
		return err
	}
	if _, err := s.deps.MigrationExecutionRepo.BulkInsert(c, execs); err != nil {
		return err
	}
	return nil
}

// buildMigrationPlanFilters builds the org-scoped filter plus optional
// case-insensitive name and exact status matches.
func (s *AssetTemplateService) buildMigrationPlanFilters(rc *reqCtx.RequestContext, query *dtos.MigrationPlanQueryDTO) (model.Map, error) {
	of, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return nil, err
	}
	filters := model.Map{}
	maps.Copy(filters, of)
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.Status != nil && *query.Status != "" {
		filters["status"] = *query.Status
	}
	return filters, nil
}

// mapMigrationPlansToDtos maps a page of plan entities to response DTOs.
func (s *AssetTemplateService) mapMigrationPlansToDtos(plans []entities.MigrationPlan) []dtos.MigrationPlanResponse {
	items := make([]dtos.MigrationPlanResponse, 0, len(plans))
	for i := range plans {
		items = append(items, *s.toMigrationPlanResponse(&plans[i]))
	}
	return items
}

// fetchMigrationPlan loads a plan by id, reporting 404 when unknown.
func (s *AssetTemplateService) fetchMigrationPlan(c ctx.Context, planId *string) (*entities.MigrationPlan, error) {
	plan, err := s.deps.MigrationPlanRepo.FindById(c, planId)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Migration plan not found"}}
	}
	return plan, nil
}

// buildMigrationPlanPatch converts the partial update to a $set map with a fresh
// updated timestamp.
func (s *AssetTemplateService) buildMigrationPlanPatch(dto *dtos.MigrationPlanUpdateRequest) (map[string]any, error) {
	patch, err := mapper.DtoToMap(dto)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"invalid update request"}}
	}
	patch["updated"] = time.Now()
	return patch, nil
}

// rescheduleMigrationIfNeeded re-arms the start timer when ScheduleAt moves. The
// prior timer message may still fire, but the runner's status guards make a
// duplicate start a no-op.
func (s *AssetTemplateService) rescheduleMigrationIfNeeded(planId *string, dto *dtos.MigrationPlanUpdateRequest) {
	if dto.ScheduleAt != nil {
		_ = s.deps.MigrationScheduler.ScheduleStart(*planId, *dto.ScheduleAt)
	}
}

// markMigrationPlanCancelled transitions the plan to cancelled.
func (s *AssetTemplateService) markMigrationPlanCancelled(c ctx.Context, planId *string) error {
	if _, err := s.deps.MigrationPlanRepo.FindByIdAndUpdate(c, planId, map[string]any{
		"status":  string(entities.PlanCancelled),
		"updated": time.Now(),
	}); err != nil {
		return err
	}
	// Cancelling a plan also halts its still-pending per-asset executions, so the
	// assets in the plan stop reading as queued once the plan is cancelled.
	oid, err := model.ToObjectID(*planId)
	if err != nil {
		return nil
	}
	_, err = s.deps.MigrationExecutionRepo.CancelPending(c, oid)
	return err
}

// buildMigrationExecutionFilters builds the status/asset filter for a plan's
// executions.
func (s *AssetTemplateService) buildMigrationExecutionFilters(query *dtos.MigrationExecutionQueryDTO) model.Map {
	filters := model.Map{}
	if query.Status != nil && *query.Status != "" {
		filters["status"] = *query.Status
	}
	if query.AssetID != nil && *query.AssetID != "" {
		if assetID, err := model.ToObjectID(*query.AssetID); err == nil {
			filters["assetId"] = assetID
		}
	}
	return filters
}

// mapMigrationExecutionsToDtos maps a page of execution entities to response
// DTOs (ids as strings).
func (s *AssetTemplateService) mapMigrationExecutionsToDtos(execs []entities.MigrationExecution) ([]dtos.MigrationExecutionResponse, error) {
	items := make([]dtos.MigrationExecutionResponse, 0, len(execs))
	for i := range execs {
		dto, err := mapper.EntityToDtoWithOptions[entities.MigrationExecution, dtos.MigrationExecutionResponse](
			&execs[i], mapper.MapperOptions{ObjectIdToString: true},
		)
		if err != nil {
			return nil, err
		}
		items = append(items, *dto)
	}
	return items, nil
}

// toMigrationPlanResponse maps a plan entity to its response DTO, flattening the
// counters and deriving progress.
func (s *AssetTemplateService) toMigrationPlanResponse(plan *entities.MigrationPlan) *dtos.MigrationPlanResponse {
	return &dtos.MigrationPlanResponse{
		ID:             plan.ID.Hex(),
		OrgID:          plan.OrgID.Hex(),
		Name:           plan.Name,
		Description:    plan.Description,
		FromTemplateID:   plan.FromTemplateID.Hex(),
		FromTemplateName: plan.FromTemplateName,
		ToTemplateID:     plan.ToTemplateID.Hex(),
		ToTemplateName:   plan.ToTemplateName,
		ScheduleAt:       plan.ScheduleAt,
		BatchSize:      plan.BatchSize,
		Status:         string(plan.Status),
		Total:          plan.Counters.Total,
		Migrated:       plan.Counters.Migrated,
		Failed:         plan.Counters.Failed,
		ProgressPct:    s.migrationProgressPct(plan.Counters),
		Created:        plan.Created,
		Updated:        plan.Updated,
	}
}

// migrationProgressPct returns the rounded completion percentage, guarding the
// zero-asset plan against division by zero.
func (s *AssetTemplateService) migrationProgressPct(c entities.PlanCounters) int {
	if c.Total == 0 {
		return 0
	}
	return int(math.Round(float64((c.Migrated+c.Failed)*100) / float64(c.Total)))
}
