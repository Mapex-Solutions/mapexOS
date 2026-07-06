package services

import (
	"context"
	"strings"
	"testing"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"
)

func inProgressPlan(fwID *entities.Firmware) *entities.OTAPlan {
	planID := newTestObjectID()
	return &entities.OTAPlan{
		ID:            planID,
		OrgID:         newTestObjectID(),
		FirmwareID:    fwID.ID,
		Status:        entities.PlanInProgress,
		RolloutConfig: entities.OTARolloutConfig{RatePerMinute: 10},
	}
}

func TestReconciler_Tick_DispatchesOnlyOnlineWithFreshURLs(t *testing.T) {
	fw := &entities.Firmware{ID: newTestObjectID(), Version: "2.0.0", Checksum: "sum", Size: 10, ObjectKey: "k"}
	plan := inProgressPlan(fw)

	onlineExec := entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecQueued}
	offlineExec := entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecQueued}
	secondOnline := entities.OTAExecution{ID: newTestObjectID(), PlanID: plan.ID, AssetID: newTestObjectID(), State: entities.ExecQueued}

	assetReader := &mockAssetReader{infoByID: map[string]*ports.OTAAssetInfo{
		onlineExec.AssetID.Hex():   {AssetUUID: "dev-online-1", Protocol: "mqtt"},
		offlineExec.AssetID.Hex():  {AssetUUID: "dev-offline", Protocol: "mqtt"},
		secondOnline.AssetID.Hex(): {AssetUUID: "dev-online-2", Protocol: "mqtt"},
	}}
	presence := &mockPresence{online: map[string]bool{"dev-online-1": true, "dev-online-2": true}}
	store := &mockStore{getURLs: []string{"https://get/one", "https://get/two"}}
	edge := &mockEdge{}
	execRepo := &mockExecutionRepo{
		byPlan: []entities.OTAExecution{onlineExec, offlineExec, secondOnline},
		byID:   map[string]*entities.OTAExecution{},
	}
	planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
	fwRepo := &mockFirmwareRepo{byID: map[string]*entities.Firmware{fw.ID.Hex(): fw}}

	rec := NewReconciler(ReconcilerDeps{
		PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: fwRepo,
		Store: store, Presence: presence, Edge: edge, AssetReader: assetReader,
		MaxAttempts: 2, DefaultRate: 10,
	})

	if err := rec.Tick(context.Background(), plan.ID.Hex()); err != nil {
		t.Fatalf("Tick error: %v", err)
	}

	if len(edge.dispatched) != 2 {
		t.Fatalf("dispatched = %d, want 2 (offline device must be skipped)", len(edge.dispatched))
	}
	if edge.dispatched[0].DownloadURL == edge.dispatched[1].DownloadURL {
		t.Fatalf("each dispatch must mint a FRESH url, both got %q", edge.dispatched[0].DownloadURL)
	}
	if edge.dispatched[0].TargetVersion != "2.0.0" {
		t.Fatalf("command must carry the firmware version")
	}
	if !strings.HasPrefix(edge.dispatched[0].ReportTarget, "events/") ||
		!strings.HasSuffix(edge.dispatched[0].ReportTarget, "/ota_status") {
		t.Fatalf("mqtt reportTarget must follow events/{assetUUID}/ota_status, got %q", edge.dispatched[0].ReportTarget)
	}

	// The execution update advances state + attempts and NEVER stores the URL.
	if len(execRepo.updates) != 2 {
		t.Fatalf("expected 2 execution updates, got %d", len(execRepo.updates))
	}
	for _, update := range execRepo.updates {
		if update["state"] != string(entities.ExecInitiated) {
			t.Fatalf("dispatch must advance to INITIATED, got %v", update["state"])
		}
		for key, value := range update {
			if s, ok := value.(string); ok && strings.Contains(s, "https://get/") {
				t.Fatalf("download URL persisted under %q — the URL must never be stored", key)
			}
		}
	}
}

func TestReconciler_Tick_AutoAbortsPastThreshold(t *testing.T) {
	fw := &entities.Firmware{ID: newTestObjectID(), ObjectKey: "k"}
	plan := inProgressPlan(fw)
	plan.RolloutConfig = entities.OTARolloutConfig{RatePerMinute: 10, AbortThresholdPct: 20, AbortMinExecuted: 5}
	plan.Counters = entities.PlanCounters{Total: 10, Succeeded: 4, Failed: 3}

	planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{plan.ID.Hex(): plan}}
	edge := &mockEdge{}
	rec := NewReconciler(ReconcilerDeps{
		PlanRepo: planRepo, ExecutionRepo: &mockExecutionRepo{}, FirmwareRepo: &mockFirmwareRepo{},
		Store: &mockStore{}, Presence: &mockPresence{}, Edge: edge, AssetReader: &mockAssetReader{},
		MaxAttempts: 2, DefaultRate: 10,
	})

	if err := rec.Tick(context.Background(), plan.ID.Hex()); err != nil {
		t.Fatalf("Tick error: %v", err)
	}
	if len(edge.dispatched) != 0 {
		t.Fatalf("aborted plan must not dispatch")
	}
	if len(planRepo.updates) != 1 || planRepo.updates[0]["status"] != string(entities.PlanAborted) {
		t.Fatalf("expected ABORTED update, got %v", planRepo.updates)
	}
}

func TestReconciler_PendingJob_Table(t *testing.T) {
	fw := &entities.Firmware{ID: newTestObjectID(), Version: "2.0.0", Checksum: "sum", Size: 10, ObjectKey: "k"}
	plan := inProgressPlan(fw)
	closedPlan := inProgressPlan(fw)
	closedPlan.Status = entities.PlanClosed

	tests := []struct {
		name        string
		assetUUID   string
		executions  []entities.OTAExecution
		planForExec *entities.OTAPlan
		wantJob     bool
		wantState   string
	}{
		{
			name:      "queued execution in a live plan is served and advances to INITIATED",
			assetUUID: "dev-1",
			executions: []entities.OTAExecution{
				{ID: newTestObjectID(), PlanID: plan.ID, AssetUUID: "dev-1", State: entities.ExecQueued},
			},
			planForExec: plan,
			wantJob:     true,
			wantState:   string(entities.ExecInitiated),
		},
		{
			name:      "execution whose plan is closed is not served",
			assetUUID: "dev-2",
			executions: []entities.OTAExecution{
				{ID: newTestObjectID(), PlanID: closedPlan.ID, AssetUUID: "dev-2", State: entities.ExecQueued},
			},
			planForExec: closedPlan,
		},
		{
			name:      "no actionable execution returns nil",
			assetUUID: "dev-3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{}}
			if tt.planForExec != nil {
				planRepo.byID[tt.planForExec.ID.Hex()] = tt.planForExec
			}
			execRepo := &mockExecutionRepo{byFilters: tt.executions, byID: map[string]*entities.OTAExecution{}}
			store := &mockStore{getURLs: []string{"https://get/poll"}}
			rec := NewReconciler(ReconcilerDeps{
				PlanRepo: planRepo, ExecutionRepo: execRepo,
				FirmwareRepo: &mockFirmwareRepo{byID: map[string]*entities.Firmware{fw.ID.Hex(): fw}},
				Store:        store, Presence: &mockPresence{}, Edge: &mockEdge{}, AssetReader: &mockAssetReader{},
				MaxAttempts: 2, DefaultRate: 10,
			})

			job, err := rec.PendingJob(context.Background(), tt.assetUUID)
			if err != nil {
				t.Fatalf("PendingJob error: %v", err)
			}
			if !tt.wantJob {
				if job != nil {
					t.Fatalf("expected no job, got %+v", job)
				}
				return
			}
			if job == nil {
				t.Fatalf("expected a job")
			}
			if job.DownloadURL != "https://get/poll" {
				t.Fatalf("job must carry a freshly minted URL")
			}
			if len(execRepo.updates) != 1 || execRepo.updates[0]["state"] != tt.wantState {
				t.Fatalf("first poll must advance QUEUED to INITIATED, got %v", execRepo.updates)
			}
		})
	}
}
