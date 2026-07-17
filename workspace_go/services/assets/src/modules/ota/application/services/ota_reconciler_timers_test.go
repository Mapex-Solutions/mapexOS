package services

import (
	"context"
	"testing"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"
)

func TestReconcilerTimers_OnStart_Table(t *testing.T) {
	tests := []struct {
		name          string
		status        entities.PlanStatus
		wantStarted   bool
		wantScheduled bool
	}{
		{name: "scheduled plan starts and schedules the close timer", status: entities.PlanScheduled, wantStarted: true, wantScheduled: true},
		{name: "already in-progress plan is a no-op (idempotent)", status: entities.PlanInProgress},
		{name: "canceled plan is a no-op", status: entities.PlanCanceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &entities.OTAPlan{ID: newTestObjectID(), Status: tt.status, MaxTime: time.Now().Add(time.Hour)}
			planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
			sched := &mockOTAScheduler{}
			timers := NewReconcilerTimers(ReconcilerTimersDeps{
				PlanRepo: planRepo, ExecutionRepo: &mockExecutionRepo{}, FirmwareRepo: &mockFirmwareRepo{},
				Store: &mockStore{}, Scheduler: sched,
			})

			if err := timers.OnStart(context.Background(), plan.ID.Hex()); err != nil {
				t.Fatalf("OnStart error: %v", err)
			}
			if tt.wantStarted {
				if len(planRepo.updates) != 1 || planRepo.updates[0]["status"] != string(entities.PlanInProgress) {
					t.Fatalf("expected IN_PROGRESS update, got %v", planRepo.updates)
				}
			} else if len(planRepo.updates) != 0 {
				t.Fatalf("no-op case must not update the plan")
			}
			if tt.wantScheduled != (len(sched.closes) == 1) {
				t.Fatalf("close-timer scheduling mismatch: %v", sched.closes)
			}
		})
	}
}

func TestReconcilerTimers_ClosePlan_Table(t *testing.T) {
	fw := &entities.Firmware{ID: newTestObjectID(), ObjectKey: "org/tpl/fw.bin"}

	tests := []struct {
		name      string
		status    entities.PlanStatus
		counters  entities.PlanCounters
		reason    string
		wantFinal string
		wantClose bool
		wantNoOp  bool
	}{
		{
			name:   "maxTime close with stragglers -> CLOSED, .bin retained, stragglers timed out",
			status: entities.PlanInProgress, counters: entities.PlanCounters{Total: 5, Succeeded: 3},
			reason: CloseReasonMaxTime, wantFinal: string(entities.PlanClosed), wantClose: true,
		},
		{
			name:   "all devices succeeded -> COMPLETED",
			status: entities.PlanInProgress, counters: entities.PlanCounters{Total: 5, Succeeded: 5},
			reason: CloseReasonAllTerminal, wantFinal: string(entities.PlanCompleted), wantClose: true,
		},
		{
			name:   "operator cancel -> CANCELED",
			status: entities.PlanScheduled, counters: entities.PlanCounters{Total: 5},
			reason: CloseReasonCancel, wantFinal: string(entities.PlanCanceled), wantClose: true,
		},
		{
			name:   "already terminal plan is a no-op (idempotent close)",
			status: entities.PlanClosed, reason: CloseReasonMaxTime, wantNoOp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &entities.OTAPlan{ID: newTestObjectID(), FirmwareID: fw.ID, Status: tt.status, Counters: tt.counters}
			planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
			fwRepo := &mockFirmwareRepo{byID: map[string]*entities.Firmware{fw.ID.Hex(): fw}}
			// Final status derives from the execution counts (source of truth), not
			// the plan.Counters cache — mirror the intended success tally here.
			execRepo := &mockExecutionRepo{counts: map[entities.ExecutionState]int64{entities.ExecUpdated: int64(tt.counters.Succeeded)}}
			store := &mockStore{}
			timers := NewReconcilerTimers(ReconcilerTimersDeps{
				PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: fwRepo,
				Store: store, Scheduler: &mockOTAScheduler{},
			})

			if err := timers.ClosePlan(context.Background(), plan.ID.Hex(), tt.reason); err != nil {
				t.Fatalf("ClosePlan error: %v", err)
			}

			if tt.wantNoOp {
				if len(planRepo.updates) != 0 || len(store.deleted) != 0 || len(execRepo.updatedMany) != 0 {
					t.Fatalf("terminal plan must not be re-closed")
				}
				return
			}

			if len(execRepo.updatedMany) != 1 {
				t.Fatalf("stragglers must be timed out via one UpdateMany")
			}
			if !tt.wantClose {
				return
			}
			// Firmware is retained on close: the .bin is NOT deleted and the
			// artifact status is left untouched (no PURGED marking).
			if len(store.deleted) != 0 {
				t.Fatalf(".bin must be retained, got deletes %v", store.deleted)
			}
			if len(fwRepo.updates) != 0 {
				t.Fatalf("firmware status must be untouched, got %v", fwRepo.updates)
			}
			if len(planRepo.updates) != 1 || planRepo.updates[0]["status"] != tt.wantFinal {
				t.Fatalf("final status = %v, want %s", planRepo.updates, tt.wantFinal)
			}
		})
	}
}

func TestReconcilerTimers_RunScan_EarlyClosesAllTerminalPlan(t *testing.T) {
	// Every device succeeded -> the scan sees all executions terminal (from the
	// counts, the source of truth) and early-closes the plan COMPLETED.
	plan := &entities.OTAPlan{ID: newTestObjectID(), Status: entities.PlanInProgress, Counters: entities.PlanCounters{Total: 2}}
	planRepo := &mockPlanRepo{
		byID:  map[string]*entities.OTAPlan{plan.ID.Hex(): plan},
		found: []entities.OTAPlan{*plan},
	}
	execRepo := &mockExecutionRepo{counts: map[entities.ExecutionState]int64{entities.ExecUpdated: 2}}
	sched := &mockOTAScheduler{}
	timers := NewReconcilerTimers(ReconcilerTimersDeps{
		PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: &mockFirmwareRepo{},
		Store: &mockStore{}, Scheduler: sched,
		Reconciler: NewReconciler(ReconcilerDeps{
			PlanRepo: &mockPlanRepo{}, ExecutionRepo: &mockExecutionRepo{}, FirmwareRepo: &mockFirmwareRepo{},
			Store: &mockStore{}, Presence: &mockPresence{}, Edge: &mockEdge{}, AssetReader: &mockAssetReader{},
		}),
	})

	timers.RunScan(context.Background())

	if len(execRepo.updatedMany) != 1 {
		t.Fatalf("all-terminal plan must be early-closed (one straggler UpdateMany), got %d", len(execRepo.updatedMany))
	}
	if len(planRepo.updates) != 1 || planRepo.updates[0]["status"] != string(entities.PlanCompleted) {
		t.Fatalf("plan must finalize COMPLETED, got %v", planRepo.updates)
	}
	// RunScan owns no cadence: it must never schedule anything (no self-rescheduling loop).
	if len(sched.starts) != 0 || len(sched.closes) != 0 {
		t.Fatalf("RunScan must not schedule timers, got starts=%v closes=%v", sched.starts, sched.closes)
	}
}

func TestReconcilerTimers_ClosePlan_DerivesFromExecutionsNotCounter(t *testing.T) {
	// The counter cache is stale (succeeded=0) but the execution reached UPDATED.
	// The close must still finalize COMPLETED — the decision comes from CountByState,
	// immunizing against a swallowed IncrementCounter.
	plan := &entities.OTAPlan{ID: newTestObjectID(), Status: entities.PlanInProgress, Counters: entities.PlanCounters{Total: 1, Succeeded: 0}}
	planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
	execRepo := &mockExecutionRepo{counts: map[entities.ExecutionState]int64{entities.ExecUpdated: 1}}
	timers := NewReconcilerTimers(ReconcilerTimersDeps{
		PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: &mockFirmwareRepo{},
		Store: &mockStore{}, Scheduler: &mockOTAScheduler{},
	})

	if err := timers.ClosePlan(context.Background(), plan.ID.Hex(), CloseReasonAllTerminal); err != nil {
		t.Fatalf("ClosePlan error: %v", err)
	}
	if len(planRepo.updates) != 1 || planRepo.updates[0]["status"] != string(entities.PlanCompleted) {
		t.Fatalf("close must derive COMPLETED from the executions despite counters.succeeded=0, got %v", planRepo.updates)
	}
}

func TestReconcilerTimers_RunScan_ReconcilesInProgressPlanWithoutClosing(t *testing.T) {
	fw := &entities.Firmware{ID: newTestObjectID(), Version: "1.0.0", ObjectKey: "org/tpl/fw.bin"}
	exec := &entities.OTAExecution{ID: newTestObjectID(), AssetID: newTestObjectID(), State: entities.ExecQueued}
	plan := &entities.OTAPlan{ID: newTestObjectID(), FirmwareID: fw.ID, Status: entities.PlanInProgress, Counters: entities.PlanCounters{Total: 1}}

	planRepo := &mockPlanRepo{
		byID:  map[string]*entities.OTAPlan{plan.ID.Hex(): plan},
		found: []entities.OTAPlan{*plan},
	}
	execRepo := &mockExecutionRepo{}
	edge := &mockEdge{}
	timers := NewReconcilerTimers(ReconcilerTimersDeps{
		PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: &mockFirmwareRepo{},
		Store: &mockStore{}, Scheduler: &mockOTAScheduler{},
		// The reconciler resolves this plan and dispatches its QUEUED execution to an online MQTT device.
		Reconciler: NewReconciler(ReconcilerDeps{
			PlanRepo:      &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}},
			ExecutionRepo: &mockExecutionRepo{byPlan: []entities.OTAExecution{*exec}, byID: map[string]*entities.OTAExecution{exec.ID.Hex(): exec}},
			FirmwareRepo:  &mockFirmwareRepo{byID: map[string]*entities.Firmware{fw.ID.Hex(): fw}},
			Store:         &mockStore{},
			Presence:      &mockPresence{online: map[string]bool{"uuid-1": true}},
			Edge:          edge,
			AssetReader:   &mockAssetReader{infoByID: map[string]*ports.OTAAssetInfo{exec.AssetID.Hex(): {Protocol: "mqtt", AssetUUID: "uuid-1"}}},
		}),
	})

	timers.RunScan(context.Background())

	// Tick ran for the plan and pushed the update to the device.
	if len(edge.dispatched) != 1 {
		t.Fatalf("Tick must dispatch the QUEUED execution to the online device, got %d dispatches", len(edge.dispatched))
	}
	// 0 of 1 terminal, so the plan must NOT be early-closed.
	if len(execRepo.updatedMany) != 0 || len(planRepo.updates) != 0 {
		t.Fatalf("non-terminal plan must not be early-closed: updatedMany=%d updates=%v", len(execRepo.updatedMany), planRepo.updates)
	}
}
