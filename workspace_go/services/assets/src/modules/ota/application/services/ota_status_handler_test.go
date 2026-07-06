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
				if len(execRepo.updates) != 0 || len(history.published) != 0 || len(switcher.switched) != 0 {
					t.Fatalf("dropped report must have no side effects")
				}
				return
			}

			if len(execRepo.updates) != 1 {
				t.Fatalf("expected 1 execution update, got %d", len(execRepo.updates))
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
