package services

import (
	"context"
	"errors"
	"testing"

	"assets/src/modules/ota/domain/entities"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

func testRequestContext(t *testing.T) *reqCtx.RequestContext {
	t.Helper()
	org := model.NewObjectID().Hex()
	return &reqCtx.RequestContext{OrgContext: &org}
}

func TestFirmwareService_InitUpload_CreatesPendingArtifact(t *testing.T) {
	repo := &mockFirmwareRepo{byID: map[string]*entities.Firmware{}}
	store := &mockStore{putURL: "https://store/put"}
	sched := &mockFirmwareScheduler{}
	svc := NewFirmwareService(FirmwareServiceDeps{FirmwareRepo: repo, Store: store, Scheduler: sched})

	resp, err := svc.InitUpload(context.Background(), testRequestContext(t), &otaDtos.FirmwareInitRequest{
		TargetTemplateID: model.NewObjectID().Hex(),
		Version:          "2.0.0",
		Filename:         "fw.bin",
		Size:             1024,
		SHA256:           "abc123",
	})
	if err != nil {
		t.Fatalf("InitUpload error: %v", err)
	}
	if resp.UploadURL != "https://store/put" || resp.FirmwareID == "" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected 1 created artifact, got %d", len(repo.created))
	}
	fw := repo.created[0]
	if fw.Status != entities.FirmwarePendingUpload {
		t.Fatalf("status = %s, want PENDING_UPLOAD", fw.Status)
	}
	if fw.Version != "2.0.0" || fw.Checksum != "abc123" {
		t.Fatalf("declared facts not persisted: %+v", fw)
	}
	if len(sched.scheduled) != 1 {
		t.Fatalf("abandon-check not scheduled")
	}
}

func TestFirmwareService_CompleteUpload_Table(t *testing.T) {
	fwID := model.NewObjectID()

	tests := []struct {
		name       string
		status     entities.FirmwareStatus
		statSize   int64
		statSum    string
		statErr    error
		wantErr    bool
		wantReady  bool
		wantPurged bool
	}{
		{
			name:   "confirmed object flips to READY and purges the abandon check",
			status: entities.FirmwarePendingUpload, statSize: 1024, statSum: "sum==",
			wantReady: true, wantPurged: true,
		},
		{
			name:   "missing object keeps PENDING_UPLOAD",
			status: entities.FirmwarePendingUpload, statErr: errors.New("no such key"),
			wantErr: true,
		},
		{
			name:   "size mismatch keeps PENDING_UPLOAD",
			status: entities.FirmwarePendingUpload, statSize: 999, statSum: "sum==",
			wantErr: true,
		},
		{
			name:   "checksum mismatch keeps PENDING_UPLOAD",
			status: entities.FirmwarePendingUpload, statSize: 1024, statSum: "other==",
			wantErr: true,
		},
		{
			name:   "non-pending artifact rejected",
			status: entities.FirmwareReady, statSize: 1024, statSum: "sum==",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFirmwareRepo{byID: map[string]*entities.Firmware{
				fwID.Hex(): {ID: fwID, Status: tt.status, Size: 1024, Checksum: "sum==", ObjectKey: "org/tpl/fw.bin"},
			}}
			store := &mockStore{statSize: tt.statSize, statSum: tt.statSum, statErr: tt.statErr}
			sched := &mockFirmwareScheduler{}
			svc := NewFirmwareService(FirmwareServiceDeps{FirmwareRepo: repo, Store: store, Scheduler: sched})

			err := svc.CompleteUpload(context.Background(), testRequestContext(t), fwID.Hex())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if len(repo.updates) != 0 {
					t.Fatalf("artifact must stay untouched on failure, got updates: %v", repo.updates)
				}
				return
			}
			if err != nil {
				t.Fatalf("CompleteUpload error: %v", err)
			}
			if !tt.wantReady {
				return
			}
			if len(repo.updates) != 1 || repo.updates[0]["status"] != string(entities.FirmwareReady) {
				t.Fatalf("expected READY update, got %v", repo.updates)
			}
			if tt.wantPurged && len(sched.purged) != 1 {
				t.Fatalf("abandon-check not purged")
			}
		})
	}
}

func TestFirmwareService_HandleFirmwareAbandon_Table(t *testing.T) {
	fwID := model.NewObjectID()

	tests := []struct {
		name        string
		status      entities.FirmwareStatus
		wantDeleted bool
	}{
		{name: "never-finalized artifact is deleted and marked ABANDONED", status: entities.FirmwarePendingUpload, wantDeleted: true},
		{name: "finalized artifact is a no-op (schedule superseded)", status: entities.FirmwareReady},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFirmwareRepo{byID: map[string]*entities.Firmware{
				fwID.Hex(): {ID: fwID, Status: tt.status, ObjectKey: "org/tpl/fw.bin"},
			}}
			store := &mockStore{}
			svc := NewFirmwareService(FirmwareServiceDeps{FirmwareRepo: repo, Store: store, Scheduler: &mockFirmwareScheduler{}})

			if err := svc.HandleFirmwareAbandon(context.Background(), fwID.Hex()); err != nil {
				t.Fatalf("HandleFirmwareAbandon error: %v", err)
			}
			if tt.wantDeleted {
				if len(store.deleted) != 1 {
					t.Fatalf("orphan object not deleted")
				}
				if len(repo.updates) != 1 || repo.updates[0]["status"] != string(entities.FirmwareAbandoned) {
					t.Fatalf("expected ABANDONED update, got %v", repo.updates)
				}
				return
			}
			if len(store.deleted) != 0 || len(repo.updates) != 0 {
				t.Fatalf("finalized artifact must not be touched")
			}
		})
	}
}
