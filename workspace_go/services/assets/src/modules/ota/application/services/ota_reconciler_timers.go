package services

import (
	"context"
	"time"

	"assets/src/modules/ota/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// NewReconcilerTimers returns a timer service over the given dependencies.
func NewReconcilerTimers(deps ReconcilerTimersDeps) *ReconcilerTimers {
	return &ReconcilerTimers{deps: deps}
}

// OnStart handles a plan's start timer: SCHEDULED → IN_PROGRESS and schedules the
// close (maxTime) timer. Idempotent (no-op if the plan is not SCHEDULED).
func (t *ReconcilerTimers) OnStart(c context.Context, planID string) error {
	plan, err := t.deps.PlanRepo.FindById(c, &planID)
	if err != nil {
		return err
	}
	if plan == nil || plan.ID.IsZero() || plan.Status != entities.PlanScheduled {
		return nil
	}
	if _, err := t.deps.PlanRepo.FindByIdAndUpdate(c, &planID, map[string]any{
		"status": string(entities.PlanInProgress), "updated": time.Now(),
	}); err != nil {
		return err
	}
	return t.deps.Scheduler.ScheduleClose(planID, plan.MaxTime)
}

// OnClose handles a plan's close (maxTime) timer.
func (t *ReconcilerTimers) OnClose(c context.Context, planID string) error {
	return t.ClosePlan(c, planID, CloseReasonMaxTime)
}

// RunScan performs one pacing sweep, driven by the elected leader's ticker:
// reconcile every IN_PROGRESS plan, then early-close the ones whose executions
// are all terminal. It schedules nothing — the leader owns the cadence.
func (t *ReconcilerTimers) RunScan(c context.Context) {
	result, err := t.deps.PlanRepo.FindWithFilters(c,
		model.Map{"status": string(entities.PlanInProgress)},
		&model.PaginationOpts{Page: 1, PerPage: 1000}, nil)
	if err != nil {
		return
	}
	for i := range result.Items {
		plan := &result.Items[i]
		planID := plan.ID.Hex()
		_ = t.deps.Reconciler.Tick(c, planID)
		if plan.Counters.Total > 0 && t.terminalCount(c, planID) >= int64(plan.Counters.Total) {
			_ = t.ClosePlan(c, planID, CloseReasonAllTerminal)
		}
	}
}

// terminalCount counts the plan's executions in a terminal state, read from the
// executions (the source of truth) rather than the denormalized plan counters —
// so a missed counter increment cannot skew the close decision.
func (t *ReconcilerTimers) terminalCount(c context.Context, planID string) int64 {
	var total int64
	for _, s := range entities.TerminalStates() {
		n, _ := t.deps.ExecutionRepo.CountByState(c, &planID, s)
		total += n
	}
	return total
}

// ClosePlan is the single idempotent close routine: mark stragglers TIMED_OUT
// and set the terminal status. The firmware .bin is RETAINED so operators can
// re-download it from the closed plan's detail page (the artifacts are tiny).
func (t *ReconcilerTimers) ClosePlan(c context.Context, planID, reason string) error {
	plan, err := t.deps.PlanRepo.FindById(c, &planID)
	if err != nil {
		return err
	}
	if plan == nil || plan.ID.IsZero() || isPlanTerminal(plan.Status) {
		return nil
	}

	// Mark non-terminal executions as TIMED_OUT.
	_, _ = t.deps.ExecutionRepo.UpdateMany(c,
		model.Map{"planId": plan.ID, "state": model.Map{"$nin": entities.TerminalStateStrings()}},
		model.Map{"$set": model.Map{"state": string(entities.ExecTimedOut), "updated": time.Now()}},
	)

	// Final status is derived from the executions (count of UPDATED), not the
	// denormalized counter cache — a missed increment cannot force a wrong COMPLETED.
	succeeded, _ := t.deps.ExecutionRepo.CountByState(c, &planID, entities.ExecUpdated)
	_, err = t.deps.PlanRepo.FindByIdAndUpdate(c, &planID, map[string]any{
		"status": string(deriveFinalStatus(plan, reason, succeeded)), "updated": time.Now(),
	})
	return err
}

// deriveFinalStatus picks the plan's terminal status from the close reason and the
// success tally (succeeded = count of UPDATED executions, from the source of truth).
func deriveFinalStatus(plan *entities.OTAPlan, reason string, succeeded int64) entities.PlanStatus {
	if reason == CloseReasonCancel {
		return entities.PlanCanceled
	}
	if plan.Counters.Total > 0 && succeeded >= int64(plan.Counters.Total) {
		return entities.PlanCompleted
	}
	return entities.PlanClosed
}

func isPlanTerminal(s entities.PlanStatus) bool {
	switch s {
	case entities.PlanCompleted, entities.PlanClosed, entities.PlanCanceled, entities.PlanAborted:
		return true
	}
	return false
}
