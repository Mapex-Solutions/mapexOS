package services

import (
	"context"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// NewReconciler returns a reconciler over the given dependencies.
func NewReconciler(deps ReconcilerDeps) *Reconciler {
	return &Reconciler{deps: deps}
}

// PendingJob serves the HTTP device poll: the first actionable execution for
// the asset inside an IN_PROGRESS plan, with a freshly minted download URL.
// QUEUED (or retryable FAILED) executions advance to INITIATED on first serve;
// in-flight ones are re-served without attempt accounting (a device may re-poll
// after a crash).
func (r *Reconciler) PendingJob(c context.Context, assetUUID string) (*ports.OTACommandPayload, error) {
	if assetUUID == "" {
		return nil, nil
	}
	result, err := r.deps.ExecutionRepo.FindWithFilters(c, model.Map{
		"assetUUID": assetUUID,
		"$or": []model.Map{
			{"state": model.Map{"$nin": []string{
				string(entities.ExecUpdated), string(entities.ExecFailed), string(entities.ExecTimedOut),
			}}},
			{"state": string(entities.ExecFailed), "attempts": model.Map{"$lt": r.deps.MaxAttempts}},
		},
	}, &model.PaginationOpts{Page: 1, PerPage: 5}, nil)
	if err != nil {
		return nil, err
	}
	for i := range result.Items {
		if payload, ok := r.serveExecution(c, &result.Items[i]); ok {
			return payload, nil
		}
	}
	return nil, nil
}

// serveExecution builds the poll command for one execution when its plan is
// live; QUEUED/FAILED advance to INITIATED with an attempt bump.
func (r *Reconciler) serveExecution(c context.Context, exec *entities.OTAExecution) (*ports.OTACommandPayload, bool) {
	planID := exec.PlanID.Hex()
	plan, err := r.deps.PlanRepo.FindById(c, &planID)
	if err != nil || plan == nil || plan.ID.IsZero() || plan.Status != entities.PlanInProgress {
		return nil, false
	}
	firmwareID := plan.FirmwareID.Hex()
	fw, err := r.deps.FirmwareRepo.FindById(c, &firmwareID)
	if err != nil || fw == nil || fw.ID.IsZero() {
		return nil, false
	}
	url, err := r.deps.Store.PresignGet(c, fw.ObjectKey, r.deps.PresignTTL)
	if err != nil {
		return nil, false
	}
	if exec.State == entities.ExecQueued || exec.State == entities.ExecFailed {
		r.updateExecution(c, exec.ID.Hex(), map[string]any{
			"state":     string(entities.ExecInitiated),
			"attempts":  exec.Attempts + 1,
			"startedAt": time.Now(),
			"updated":   time.Now(),
		})
	}
	return &ports.OTACommandPayload{
		PlanID:          plan.ID.Hex(),
		ExecutionID:     exec.ID.Hex(),
		TargetVersion:   fw.Version,
		DownloadURL:     url,
		Checksum:        fw.Checksum,
		Size:            fw.Size,
		ReportTransport: protocolHTTP,
		ReportTarget:    httpStatusReportPath,
	}, true
}

// Compile-time check that the reconciler serves the device poll port.
var _ ports.DeviceJobPort = (*Reconciler)(nil)

// Tick runs one reconciliation pass for a plan: abort-check, then dispatch the
// next paced batch of actionable executions to online devices.
func (r *Reconciler) Tick(c context.Context, planID string) error {
	plan, err := r.deps.PlanRepo.FindById(c, &planID)
	if err != nil {
		return err
	}
	if plan == nil || plan.ID.IsZero() || plan.Status != entities.PlanInProgress {
		return nil
	}

	if r.shouldAbort(plan) {
		_, _ = r.deps.PlanRepo.FindByIdAndUpdate(c, &planID,
			map[string]any{"status": string(entities.PlanAborted), "updated": time.Now()})
		return nil
	}

	firmwareID := plan.FirmwareID.Hex()
	fw, err := r.deps.FirmwareRepo.FindById(c, &firmwareID)
	if err != nil {
		return err
	}
	if fw == nil || fw.ID.IsZero() {
		return nil
	}

	rate := plan.RolloutConfig.RatePerMinute
	if rate <= 0 {
		rate = r.deps.DefaultRate
	}

	batch, err := r.actionableBatch(c, planID, rate)
	if err != nil {
		return err
	}
	for i := range batch {
		r.dispatchOne(c, plan, fw, &batch[i])
	}
	return nil
}

// actionableBatch fetches up to `rate` executions that need action: QUEUED, or
// FAILED with retries remaining.
func (r *Reconciler) actionableBatch(c context.Context, planID string, rate int) ([]entities.OTAExecution, error) {
	filters := model.Map{
		"$or": []model.Map{
			{"state": string(entities.ExecQueued)},
			{"state": string(entities.ExecFailed), "attempts": model.Map{"$lt": r.deps.MaxAttempts}},
		},
	}
	result, err := r.deps.ExecutionRepo.FindByPlan(c, &planID, filters,
		&model.PaginationOpts{Page: 1, PerPage: int64(rate)})
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// dispatchOne resolves the asset, presence-gates, mints a fresh download URL,
// dispatches to the edge, and advances the execution to INITIATED. HTTP devices
// are left QUEUED (they poll). Any failure leaves the execution untouched for a
// later tick (best-effort).
func (r *Reconciler) dispatchOne(c context.Context, plan *entities.OTAPlan, fw *entities.Firmware, exec *entities.OTAExecution) {
	assetID := exec.AssetID.Hex()
	info, err := r.deps.AssetReader.GetAssetInfo(c, assetID)
	if err != nil || info == nil {
		return
	}

	// HTTP devices poll for their job; just record the resolved protocol so the
	// poll endpoint can serve them.
	if info.Protocol == protocolHTTP {
		if exec.Protocol != info.Protocol || exec.AssetUUID != info.AssetUUID {
			r.updateExecution(c, exec.ID.Hex(), map[string]any{
				"protocol": info.Protocol, "assetUUID": info.AssetUUID, "updated": time.Now(),
			})
		}
		return
	}

	orgID := plan.OrgID.Hex()
	online, err := r.deps.Presence.IsOnline(c, orgID, info.AssetUUID)
	if err != nil || !online {
		return // offline → stays QUEUED; dispatched on next connect / tick
	}

	url, err := r.deps.Store.PresignGet(c, fw.ObjectKey, r.deps.PresignTTL)
	if err != nil {
		return
	}

	payload := ports.OTACommandPayload{
		PlanID:          plan.ID.Hex(),
		ExecutionID:     exec.ID.Hex(),
		TargetVersion:   fw.Version,
		DownloadURL:     url,
		Checksum:        fw.Checksum,
		Size:            fw.Size,
		ReportTransport: info.Protocol,
		ReportTarget:    mqttStatusTopicFor(info.AssetUUID),
	}
	if err := r.deps.Edge.Dispatch(c, info.Protocol, orgID, info.AssetUUID, payload); err != nil {
		return
	}

	r.updateExecution(c, exec.ID.Hex(), map[string]any{
		"state":     string(entities.ExecInitiated),
		"protocol":  info.Protocol,
		"assetUUID": info.AssetUUID,
		"attempts":  exec.Attempts + 1,
		"startedAt": time.Now(),
		"updated":   time.Now(),
	})
}

func (r *Reconciler) updateExecution(c context.Context, id string, fields map[string]any) {
	_, _ = r.deps.ExecutionRepo.FindByIdAndUpdate(c, &id, fields)
}

// shouldAbort reports whether the plan's failure rate has crossed the abort
// threshold after enough devices have executed.
func (r *Reconciler) shouldAbort(plan *entities.OTAPlan) bool {
	cfg := plan.RolloutConfig
	if cfg.AbortThresholdPct <= 0 || cfg.AbortMinExecuted <= 0 {
		return false
	}
	executed := plan.Counters.Succeeded + plan.Counters.Failed + plan.Counters.TimedOut
	if executed < cfg.AbortMinExecuted {
		return false
	}
	return plan.Counters.Failed*100 > cfg.AbortThresholdPct*executed
}
