package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"assets/src/modules/ota/domain/entities"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

func TestPlanService_CreatePlan_Table(t *testing.T) {
	fwID := model.NewObjectID()
	sourceID := model.NewObjectID()
	targetID := model.NewObjectID()

	validRequest := func(assetCount int) *otaDtos.OTAPlanCreateRequest {
		assetIDs := make([]string, 0, assetCount)
		for range assetCount {
			assetIDs = append(assetIDs, model.NewObjectID().Hex())
		}
		return &otaDtos.OTAPlanCreateRequest{
			FirmwareID:       fwID.Hex(),
			SourceTemplateID: sourceID.Hex(),
			AssetIDs:         assetIDs,
			StartAt:          time.Now().Add(time.Hour),
			MaxTime:          time.Now().Add(2 * time.Hour),
			Name:             "rollout",
			RolloutConfig:    otaDtos.OTARolloutConfigDTO{RatePerMinute: 60},
		}
	}

	tests := []struct {
		name           string
		firmwareStatus entities.FirmwareStatus
		activePlans    []entities.OTAPlan
		assetCount     int
		wantErr        bool
	}{
		{
			name:           "READY firmware creates the plan, queues N executions, schedules startAt",
			firmwareStatus: entities.FirmwareReady,
			assetCount:     3,
		},
		{
			name:           "non-READY firmware is rejected",
			firmwareStatus: entities.FirmwarePendingUpload,
			assetCount:     1,
			wantErr:        true,
		},
		{
			name:           "firmware already claimed by an active plan is rejected",
			firmwareStatus: entities.FirmwareReady,
			activePlans:    []entities.OTAPlan{{ID: model.NewObjectID(), Status: entities.PlanInProgress}},
			assetCount:     1,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fwRepo := &mockFirmwareRepo{byID: map[string]*entities.Firmware{
				fwID.Hex(): {ID: fwID, Status: tt.firmwareStatus, TargetTemplateID: targetID},
			}}
			planRepo := &mockPlanRepo{byID: map[string]*entities.OTAPlan{}, found: tt.activePlans}
			execRepo := &mockExecutionRepo{}
			sched := &mockOTAScheduler{}
			svc := NewPlanService(PlanServiceDeps{
				PlanRepo: planRepo, ExecutionRepo: execRepo, FirmwareRepo: fwRepo, Scheduler: sched,
			})

			resp, err := svc.CreatePlan(context.Background(), testRequestContext(t), validRequest(tt.assetCount))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got plan %+v", resp)
				}
				if len(planRepo.created) != 0 || len(execRepo.bulk) != 0 || len(sched.starts) != 0 {
					t.Fatalf("nothing may be persisted or scheduled on rejection")
				}
				return
			}
			if err != nil {
				t.Fatalf("CreatePlan error: %v", err)
			}

			if len(planRepo.created) != 1 {
				t.Fatalf("expected 1 plan created")
			}
			plan := planRepo.created[0]
			if plan.Status != entities.PlanScheduled {
				t.Fatalf("plan status = %s, want SCHEDULED", plan.Status)
			}
			if plan.TargetTemplateID != targetID {
				t.Fatalf("target template must come from the firmware")
			}
			if plan.Counters.Total != tt.assetCount {
				t.Fatalf("counters.total = %d, want %d", plan.Counters.Total, tt.assetCount)
			}

			if len(execRepo.bulk) != 1 || len(execRepo.bulk[0]) != tt.assetCount {
				t.Fatalf("expected %d queued executions in one bulk insert", tt.assetCount)
			}
			for _, exec := range execRepo.bulk[0] {
				if exec.State != entities.ExecQueued {
					t.Fatalf("execution state = %s, want QUEUED", exec.State)
				}
			}

			if len(sched.starts) != 1 {
				t.Fatalf("startAt timer not scheduled")
			}
		})
	}
}

func TestPlanService_DownloadFirmware_Table(t *testing.T) {
	fwID := model.NewObjectID()
	planID := model.NewObjectID()
	plan := &entities.OTAPlan{ID: planID, FirmwareID: fwID}
	fw := &entities.Firmware{ID: fwID, Filename: "fw.bin", ObjectKey: "org/tpl/fw.bin", Status: entities.FirmwareReady}

	tests := []struct {
		name     string
		planByID map[string]*entities.OTAPlan
		fwByID   map[string]*entities.Firmware
		statErr  error
		wantErr  bool
	}{
		{
			name:     "existing object -> presigned GET url + filename",
			planByID: map[string]*entities.OTAPlan{planID.Hex(): plan},
			fwByID:   map[string]*entities.Firmware{fwID.Hex(): fw},
		},
		{
			name:     "object gone from storage -> not found",
			planByID: map[string]*entities.OTAPlan{planID.Hex(): plan},
			fwByID:   map[string]*entities.Firmware{fwID.Hex(): fw},
			statErr:  errors.New("no such key"),
			wantErr:  true,
		},
		{
			name:     "missing firmware record -> not found",
			planByID: map[string]*entities.OTAPlan{planID.Hex(): plan},
			fwByID:   map[string]*entities.Firmware{},
			wantErr:  true,
		},
		{
			name:     "missing plan -> not found",
			planByID: map[string]*entities.OTAPlan{},
			fwByID:   map[string]*entities.Firmware{fwID.Hex(): fw},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{putURL: "https://store/fw.bin", statErr: tt.statErr}
			svc := NewPlanService(PlanServiceDeps{
				PlanRepo:     &mockPlanRepo{byID: tt.planByID},
				FirmwareRepo: &mockFirmwareRepo{byID: tt.fwByID},
				Store:        store,
				PresignTTL:   time.Minute,
			})

			resp, err := svc.DownloadFirmware(context.Background(), testRequestContext(t), planID.Hex())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", resp)
				}
				return
			}
			if err != nil {
				t.Fatalf("DownloadFirmware error: %v", err)
			}
			if resp.URL != "https://store/fw.bin" {
				t.Fatalf("url = %q, want the presigned GET url", resp.URL)
			}
			if resp.Filename != fw.Filename {
				t.Fatalf("filename = %q, want %q", resp.Filename, fw.Filename)
			}
			if resp.ExpiresAt.IsZero() {
				t.Fatalf("expiresAt must be set")
			}
		})
	}
}
