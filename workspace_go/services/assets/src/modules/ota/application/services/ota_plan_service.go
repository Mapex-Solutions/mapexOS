package services

import (
	"context"
	"maps"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"

	localDtos "assets/src/modules/ota/application/dtos"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	mapper "github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	orgfilter "github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// NewPlanService returns a plan service over the given dependencies.
func NewPlanService(deps PlanServiceDeps) ports.OTAPlanServicePort {
	return &PlanService{deps: deps}
}

// CreatePlan validates the firmware is READY, persists the plan (SCHEDULED),
// bulk-inserts one QUEUED execution per asset, and schedules the startAt timer.
func (s *PlanService) CreatePlan(c context.Context, rc *reqCtx.RequestContext, req *otaDtos.OTAPlanCreateRequest) (*otaDtos.OTAPlanResponse, error) {
	orgID, err := resolveOrgID(rc)
	if err != nil {
		return nil, err
	}

	fw, err := s.deps.FirmwareRepo.FindById(c, &req.FirmwareID)
	if err != nil {
		return nil, err
	}
	if fw == nil || fw.ID.IsZero() {
		return nil, notFound("Firmware not found")
	}
	if !fw.CanReference() {
		return nil, badRequest("firmware is not READY")
	}
	// A firmware may back only one active (SCHEDULED/IN_PROGRESS) plan at a time
	// so a rollout's lifecycle and counters map to a single plan. The binary is
	// retained on close, so this is a lifecycle guard, not a purge-safety one.
	if err := s.assertFirmwareUnclaimed(c, fw.ID); err != nil {
		return nil, err
	}

	// Map request → entity via the shared mapper (Name, Description, StartAt,
	// MaxTime, RolloutConfig, and string→ObjectId for FirmwareID/SourceTemplateID).
	plan, err := mapper.DtoToEntityWithOptions[otaDtos.OTAPlanCreateRequest, entities.OTAPlan](
		req, mapper.MapperOptions{StringToObjectId: true},
	)
	if err != nil || plan.SourceTemplateID.IsZero() {
		return nil, badRequest("invalid plan request")
	}

	id := model.NewObjectID()
	now := time.Now()
	plan.ID = id
	plan.OrgID = orgID
	plan.TargetTemplateID = fw.TargetTemplateID
	plan.Status = entities.PlanScheduled
	plan.Counters = entities.PlanCounters{Total: len(req.AssetIDs)}
	plan.Created = now
	plan.Updated = now
	// An absent startAt means "run immediately": stamp now so the start timer
	// fires at once — one scheduling path for both modes.
	if plan.StartAt.IsZero() {
		plan.StartAt = now
	}
	if !plan.MaxTime.After(plan.StartAt) {
		return nil, badRequest("maxTime must be after startAt")
	}
	if _, err := s.deps.PlanRepo.Create(c, plan); err != nil {
		return nil, err
	}

	execs := make([]*entities.OTAExecution, 0, len(req.AssetIDs))
	for _, aid := range req.AssetIDs {
		assetID, err := model.ToObjectID(aid)
		if err != nil {
			return nil, badRequest("invalid assetId: "+aid)
		}
		// Protocol is resolved at dispatch (the reconciler reads the asset for
		// presence anyway), so it is left empty at creation.
		execs = append(execs, &entities.OTAExecution{
			ID:       model.NewObjectID(),
			PlanID:   id,
			AssetID:  assetID,
			OrgID:    orgID,
			State:    entities.ExecQueued,
			QueuedAt: now,
			Updated:  now,
		})
	}
	if _, err := s.deps.ExecutionRepo.BulkInsert(c, execs); err != nil {
		return nil, err
	}

	_ = s.deps.Scheduler.ScheduleStart(id.Hex(), plan.StartAt)

	return toPlanResponse(plan)
}

// ListPlans returns the org-scoped, filtered, paginated list of plans.
func (s *PlanService) ListPlans(c context.Context, rc *reqCtx.RequestContext, query *localDtos.OTAPlanQueryDTO) (*model.PaginatedResult[otaDtos.OTAPlanResponse], error) {
	filters, err := orgScopedFilter(rc)
	if err != nil {
		return nil, err
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.Status != nil && *query.Status != "" {
		filters["status"] = *query.Status
	}

	result, err := s.deps.PlanRepo.FindWithFilters(c, filters,
		&model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}, nil)
	if err != nil {
		return nil, err
	}

	items := make([]otaDtos.OTAPlanResponse, 0, len(result.Items))
	for i := range result.Items {
		dto, err := toPlanResponse(&result.Items[i])
		if err != nil {
			return nil, err
		}
		items = append(items, *dto)
	}
	return &model.PaginatedResult[otaDtos.OTAPlanResponse]{Items: items, Pagination: result.Pagination}, nil
}

// GetPlan returns a single plan by id.
func (s *PlanService) GetPlan(c context.Context, rc *reqCtx.RequestContext, id string) (*otaDtos.OTAPlanResponse, error) {
	plan, err := s.deps.PlanRepo.FindById(c, &id)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.ID.IsZero() {
		return nil, notFound("OTA plan not found")
	}
	resp, err := toPlanResponse(plan)
	if err != nil {
		return nil, err
	}
	// Detail view: denormalize the firmware artifact so the UI shows the file
	// facts (version/name/size/checksum) without a second fetch.
	resp.Firmware = s.firmwareInfo(c, plan.FirmwareID.Hex())
	return resp, nil
}

// firmwareInfo loads the referenced firmware and maps it to the detail block.
// Returns nil when the artifact is missing (best-effort enrichment).
func (s *PlanService) firmwareInfo(c context.Context, firmwareID string) *otaDtos.OTAFirmwareInfo {
	fw, err := s.deps.FirmwareRepo.FindById(c, &firmwareID)
	if err != nil || fw == nil || fw.ID.IsZero() {
		return nil
	}
	return &otaDtos.OTAFirmwareInfo{
		ID:                fw.ID.Hex(),
		Version:           fw.Version,
		Filename:          fw.Filename,
		Size:              fw.Size,
		Checksum:          fw.Checksum,
		ChecksumAlgorithm: fw.ChecksumAlgorithm,
		Status:            string(fw.Status),
		TargetTemplateID:  fw.TargetTemplateID.Hex(),
	}
}

// DownloadFirmware confirms the plan's firmware object still exists and returns
// a short-TTL presigned GET URL for it. The bytes are fetched by the browser
// directly from object storage; the Asset MS only mints the URL.
func (s *PlanService) DownloadFirmware(c context.Context, rc *reqCtx.RequestContext, planID string) (*otaDtos.OTAFirmwareDownloadResponse, error) {
	plan, err := s.deps.PlanRepo.FindById(c, &planID)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.ID.IsZero() {
		return nil, notFound("OTA plan not found")
	}

	firmwareID := plan.FirmwareID.Hex()
	fw, err := s.deps.FirmwareRepo.FindById(c, &firmwareID)
	if err != nil {
		return nil, err
	}
	if fw == nil || fw.ID.IsZero() {
		return nil, notFound("firmware not found")
	}

	// Confirm the object is still in storage before minting a URL, so a gone
	// artifact surfaces as a clean 404 instead of a dead link.
	if _, _, err := s.deps.Store.Stat(c, fw.ObjectKey); err != nil {
		return nil, notFound("firmware not found")
	}

	url, err := s.deps.Store.PresignGet(c, fw.ObjectKey, s.deps.PresignTTL)
	if err != nil {
		return nil, err
	}
	return &otaDtos.OTAFirmwareDownloadResponse{
		URL:       url,
		Filename:  fw.Filename,
		ExpiresAt: time.Now().Add(s.deps.PresignTTL),
	}, nil
}

// UpdatePlan applies a partial update (only provided fields) and returns the
// refreshed plan.
func (s *PlanService) UpdatePlan(c context.Context, rc *reqCtx.RequestContext, id string, req *localDtos.OTAPlanUpdateRequest) (*otaDtos.OTAPlanResponse, error) {
	fields, err := mapper.DtoToMap(req)
	if err != nil {
		return nil, badRequest("invalid update request")
	}
	fields["updated"] = time.Now()

	plan, err := s.deps.PlanRepo.FindByIdAndUpdate(c, &id, fields)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.ID.IsZero() {
		return nil, notFound("OTA plan not found")
	}
	return toPlanResponse(plan)
}

// DeletePlan cancels the plan (marks CANCELED and purges its timers). The full
// close routine (executions TIMED_OUT + .bin deletion) is the reconciler-timers
// ClosePlan, wired in a later task.
func (s *PlanService) DeletePlan(c context.Context, rc *reqCtx.RequestContext, id string) error {
	fields := map[string]any{"status": string(entities.PlanCanceled), "updated": time.Now()}
	plan, err := s.deps.PlanRepo.FindByIdAndUpdate(c, &id, fields)
	if err != nil {
		return err
	}
	if plan == nil || plan.ID.IsZero() {
		return notFound("OTA plan not found")
	}
	// The pending close timer fires later and ClosePlan is idempotent (it sees
	// the CANCELED status and finalizes once) — no per-plan purge needed.
	return nil
}

// ListExecutions returns the paginated per-device executions of a plan. The
// stored execution state is the current state (status updates persist to Mongo),
// so this is the live view.
func (s *PlanService) ListExecutions(c context.Context, rc *reqCtx.RequestContext, planID string, query *localDtos.OTAExecutionQueryDTO) (*model.PaginatedResult[otaDtos.OTAExecutionResponse], error) {
	filters := model.Map{}
	if query.State != nil && *query.State != "" {
		filters["state"] = *query.State
	}
	if query.AssetID != nil && *query.AssetID != "" {
		if assetID, err := model.ToObjectID(*query.AssetID); err == nil {
			filters["assetId"] = assetID
		}
	}

	result, err := s.deps.ExecutionRepo.FindByPlan(c, &planID, filters,
		&model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())})
	if err != nil {
		return nil, err
	}

	items := make([]otaDtos.OTAExecutionResponse, 0, len(result.Items))
	for i := range result.Items {
		dto, err := mapper.EntityToDtoWithOptions[entities.OTAExecution, otaDtos.OTAExecutionResponse](
			&result.Items[i], mapper.MapperOptions{ObjectIdToString: true},
		)
		if err != nil {
			return nil, err
		}
		items = append(items, *dto)
	}
	return &model.PaginatedResult[otaDtos.OTAExecutionResponse]{Items: items, Pagination: result.Pagination}, nil
}

// assertFirmwareUnclaimed rejects the create when another live (non-closed)
// plan already references the firmware.
func (s *PlanService) assertFirmwareUnclaimed(c context.Context, firmwareID model.ObjectId) error {
	existing, err := s.deps.PlanRepo.FindWithFilters(c, model.Map{
		"firmwareId": firmwareID,
		"status": model.Map{"$in": []string{
			string(entities.PlanScheduled), string(entities.PlanInProgress),
		}},
	}, &model.PaginationOpts{Page: 1, PerPage: 1}, nil)
	if err != nil {
		return err
	}
	if len(existing.Items) > 0 {
		return badRequest("firmware is already referenced by an active plan")
	}
	return nil
}

// toPlanResponse maps a plan entity to its response DTO via the shared mapper
// (ids → strings, nested counters copied field-for-field).
func toPlanResponse(plan *entities.OTAPlan) (*otaDtos.OTAPlanResponse, error) {
	return mapper.EntityToDtoWithOptions[entities.OTAPlan, otaDtos.OTAPlanResponse](
		plan, mapper.MapperOptions{ObjectIdToString: true},
	)
}

// orgScopedFilter builds the base org filter from the request context.
func orgScopedFilter(rc *reqCtx.RequestContext) (model.Map, error) {
	of, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return nil, err
	}
	filters := model.Map{}
	maps.Copy(filters, of)
	return filters, nil
}

// resolveOrgID extracts the request's organization id (shared by the OTA
// services in this package).
func resolveOrgID(rc *reqCtx.RequestContext) (model.ObjectId, error) {
	if rc == nil || rc.OrgContext == nil || *rc.OrgContext == "" {
		return model.ObjectId{}, badRequest("organization context is required")
	}
	orgID, err := model.ToObjectID(*rc.OrgContext)
	if err != nil {
		return model.ObjectId{}, badRequest("invalid organization context")
	}
	return orgID, nil
}

// Compile-time check that PlanService implements the port.
var _ ports.OTAPlanServicePort = (*PlanService)(nil)
