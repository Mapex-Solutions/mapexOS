package services

import (
	"context"
	"testing"

	"assets/src/modules/ota/domain/entities"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
)

func TestStatusHandler_HandleStatus_Table(t *testing.T) {
	targetTemplate := newTestObjectID()

	tests := []struct {
		name          string
		execState     entities.ExecutionState
		execAssetUUID string
		advisory      otaEvents.OTAStatusAdvisory
		wantUpdated   bool
		wantCounter   string
		wantSwitch    bool
	}{
		{
			name:          "downloading updates the execution, live state, and history",
			execState:     entities.ExecInitiated,
			execAssetUUID: "dev-1",
			advisory:      otaEvents.OTAStatusAdvisory{AssetUUID: "dev-1", Status: "downloading", Progress: 40},
			wantUpdated:   true,
		},
		{
			name:          "updated bumps the succeeded counter and switches the template",
			execState:     entities.ExecUpdating,
			execAssetUUID: "dev-1",
			advisory:      otaEvents.OTAStatusAdvisory{AssetUUID: "dev-1", Status: "updated", Progress: 100},
			wantUpdated:   true,
			wantCounter:   "counters.succeeded",
			wantSwitch:    true,
		},
		{
			name:          "failed bumps the failed counter without switching",
			execState:     entities.ExecDownloading,
			execAssetUUID: "dev-1",
			advisory:      otaEvents.OTAStatusAdvisory{AssetUUID: "dev-1", Status: "failed", Error: "flash error"},
			wantUpdated:   true,
			wantCounter:   "counters.failed",
		},
		{
			name:          "terminal execution ignores late reports (idempotent)",
			execState:     entities.ExecUpdated,
			execAssetUUID: "dev-1",
			advisory:      otaEvents.OTAStatusAdvisory{AssetUUID: "dev-1", Status: "downloading"},
		},
		{
			name:          "foreign device cannot touch another device's execution",
			execState:     entities.ExecDownloading,
			execAssetUUID: "dev-1",
			advisory:      otaEvents.OTAStatusAdvisory{AssetUUID: "dev-EVIL", Status: "failed"},
		},
		{
			name:          "unknown status string is dropped",
			execState:     entities.ExecDownloading,
			execAssetUUID: "dev-1",
			advisory:      otaEvents.OTAStatusAdvisory{AssetUUID: "dev-1", Status: "exploded"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &entities.OTAPlan{ID: newTestObjectID(), TargetTemplateID: targetTemplate, Status: entities.PlanInProgress}
			exec := &entities.OTAExecution{
				ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(),
				State: tt.execState, AssetUUID: tt.execAssetUUID,
			}
			tt.advisory.OTAExecutionID = exec.ID.Hex()

			execRepo := &mockExecutionRepo{byID: map[string]*entities.OTAExecution{exec.ID.Hex(): exec}}
			planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
			live := &mockLiveState{}
			history := &mockHistory{}
			switcher := &mockTemplateSwitcher{}
			handler := NewStatusHandler(StatusHandlerDeps{
				ExecutionRepo: execRepo, PlanRepo: planRepo,
				LiveState: live, History: history, TemplateSwitcher: switcher,
			})

			if err := handler.HandleStatus(context.Background(), tt.advisory); err != nil {
				t.Fatalf("HandleStatus error: %v", err)
			}

			if !tt.wantUpdated {
				if len(execRepo.updatedWhere) != 0 || len(history.published) != 0 || len(switcher.switched) != 0 {
					t.Fatalf("dropped report must have no side effects")
				}
				return
			}

			if len(execRepo.updatedWhere) != 1 {
				t.Fatalf("expected 1 execution update, got %d", len(execRepo.updatedWhere))
			}
			if len(live.set) != 1 {
				t.Fatalf("live state must be refreshed")
			}
			if len(history.published) != 1 {
				t.Fatalf("advisory must be forwarded to the history stream")
			}

			if tt.wantCounter != "" {
				if len(planRepo.increments) != 1 || planRepo.increments[0] != tt.wantCounter {
					t.Fatalf("counter increments = %v, want [%s]", planRepo.increments, tt.wantCounter)
				}
			} else if len(planRepo.increments) != 0 {
				t.Fatalf("non-terminal transitions must not bump counters, got %v", planRepo.increments)
			}

			if tt.wantSwitch {
				if len(switcher.switched) != 1 || switcher.switched[0][1] != targetTemplate.Hex() {
					t.Fatalf("UPDATED must switch the asset to the plan's target template, got %v", switcher.switched)
				}
			} else if len(switcher.switched) != 0 {
				t.Fatalf("template must switch ONLY on updated")
			}
		})
	}
}

// newStatusHandlerFor wires a handler over one execution, returning the shared
// mocks so a test can drive a sequence of advisories against it.
func newStatusHandlerFor(exec *entities.OTAExecution, plan *entities.OTAPlan) (*StatusHandler, *mockExecutionRepo, *mockPlanRepo) {
	execRepo := &mockExecutionRepo{byID: map[string]*entities.OTAExecution{exec.ID.Hex(): exec}}
	planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
	handler := NewStatusHandler(StatusHandlerDeps{
		ExecutionRepo: execRepo, PlanRepo: planRepo,
		LiveState: &mockLiveState{}, History: &mockHistory{}, TemplateSwitcher: &mockTemplateSwitcher{},
	})
	return handler, execRepo, planRepo
}

func advisoryFor(exec *entities.OTAExecution, status string, progress int32) otaEvents.OTAStatusAdvisory {
	return otaEvents.OTAStatusAdvisory{OTAExecutionID: exec.ID.Hex(), AssetUUID: exec.AssetUUID, Status: status, Progress: progress}
}

func TestStatusHandler_ForwardOnly_NoRegression(t *testing.T) {
	plan := &entities.OTAPlan{ID: newTestObjectID(), Status: entities.PlanInProgress}

	t.Run("terminal UPDATED is never regressed by a late lower-status advisory", func(t *testing.T) {
		exec := &entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecUpdating, Percentage: 80, AssetUUID: "dev-1"}
		h, _, planRepo := newStatusHandlerFor(exec, plan)

		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "updated", 100)) // applies
		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "updating", 80)) // reordered/late → must no-op

		if exec.State != entities.ExecUpdated || exec.Percentage != 100 {
			t.Fatalf("state regressed to %s@%d, want UPDATED@100", exec.State, exec.Percentage)
		}
		if len(planRepo.increments) != 1 || planRepo.increments[0] != "counters.succeeded" {
			t.Fatalf("succeeded must bump exactly once, got %v", planRepo.increments)
		}
	})

	t.Run("non-terminal regression is blocked by the percentage precondition", func(t *testing.T) {
		exec := &entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecUpdating, Percentage: 80, AssetUUID: "dev-1"}
		h, _, _ := newStatusHandlerFor(exec, plan)

		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "downloaded", 40)) // earlier status, lower % → no-op

		if exec.State != entities.ExecUpdating || exec.Percentage != 80 {
			t.Fatalf("earlier-status advisory regressed the state to %s@%d", exec.State, exec.Percentage)
		}
	})

	t.Run("duplicate UPDATED (at-least-once) does not double-count", func(t *testing.T) {
		exec := &entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecUpdating, Percentage: 80, AssetUUID: "dev-1"}
		h, _, planRepo := newStatusHandlerFor(exec, plan)

		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "updated", 100))
		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "updated", 100)) // redelivery

		if len(planRepo.increments) != 1 {
			t.Fatalf("redelivered UPDATED must bump once, got %v", planRepo.increments)
		}
	})
}

func TestStatusHandler_FailedOrthogonality(t *testing.T) {
	plan := &entities.OTAPlan{ID: newTestObjectID(), Status: entities.PlanInProgress}

	t.Run("failed applies over higher non-terminal progress", func(t *testing.T) {
		exec := &entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecUpdating, Percentage: 80, AssetUUID: "dev-1"}
		h, _, planRepo := newStatusHandlerFor(exec, plan)

		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "failed", 40))

		if exec.State != entities.ExecFailed {
			t.Fatalf("failed at a lower %% must still apply over updating@80, got %s", exec.State)
		}
		if len(planRepo.increments) != 1 || planRepo.increments[0] != "counters.failed" {
			t.Fatalf("failed must bump the failed counter, got %v", planRepo.increments)
		}
	})

	t.Run("failed never overwrites a terminal UPDATED (first terminal wins)", func(t *testing.T) {
		exec := &entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecUpdating, Percentage: 80, AssetUUID: "dev-1"}
		h, _, planRepo := newStatusHandlerFor(exec, plan)

		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "updated", 100))
		_ = h.HandleStatus(context.Background(), advisoryFor(exec, "failed", 40)) // must no-op

		if exec.State != entities.ExecUpdated {
			t.Fatalf("terminal UPDATED must not be overwritten by a late failed, got %s", exec.State)
		}
		if len(planRepo.increments) != 1 || planRepo.increments[0] != "counters.succeeded" {
			t.Fatalf("only the succeeded bump should stand, got %v", planRepo.increments)
		}
	})
}
