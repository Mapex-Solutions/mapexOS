package services

import (
	"context"
	"testing"
	"time"

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
				Store: &mockStore{}, Scheduler: sched, ScanInterval: time.Minute,
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
			execRepo := &mockExecutionRepo{}
			store := &mockStore{}
			timers := NewReconcilerTimers(ReconcilerTimersDeps{
				PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: fwRepo,
				Store: store, Scheduler: &mockOTAScheduler{}, ScanInterval: time.Minute,
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

func TestReconcilerTimers_OnScanTick_Reschedules(t *testing.T) {
	sched := &mockOTAScheduler{}
	timers := NewReconcilerTimers(ReconcilerTimersDeps{
		PlanRepo: &mockPlanRepo{}, ExecutionRepo: &mockExecutionRepo{}, FirmwareRepo: &mockFirmwareRepo{},
		Store: &mockStore{}, Scheduler: sched, ScanInterval: time.Minute,
		Reconciler: NewReconciler(ReconcilerDeps{
			PlanRepo: &mockPlanRepo{}, ExecutionRepo: &mockExecutionRepo{}, FirmwareRepo: &mockFirmwareRepo{},
			Store: &mockStore{}, Presence: &mockPresence{}, Edge: &mockEdge{}, AssetReader: &mockAssetReader{},
		}),
	})

	if err := timers.OnScanTick(context.Background()); err != nil {
		t.Fatalf("OnScanTick error: %v", err)
	}
	if sched.scans != 1 {
		t.Fatalf("the scan must always reschedule itself (self-perpetuating loop)")
	}
}
