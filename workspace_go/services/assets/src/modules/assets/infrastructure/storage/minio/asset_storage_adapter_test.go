package minio

import (
	"testing"

	"assets/src/modules/assets/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// TestBuildReadModel_HeartbeatMode_Explicit asserts that when the persisted
// entity has HealthMonitor.HeartbeatMode='explicit', the read model written
// to MinIO (L2 cache) carries the same value — so js-executor reads the
// correct mode and the gate (publish vs skip) behaves as configured.
//
// Guards against regressions where a future copier change or a missed
// converter field would silently drop heartbeatMode from the L2 cache and
// force every device into implicit mode.
func TestBuildReadModel_HeartbeatMode_Explicit(t *testing.T) {
	adapter := &AssetStorageAdapter{}

	orgId, _ := model.ToObjectID("0000000000000000000aa001")
	asset := &entities.Asset{
		ID:        model.ObjectId(orgId),
		AssetUUID: "uuid-1",
		OrgID:     orgId,
		HealthMonitor: &entities.HealthMonitorConfig{
			Enabled:          true,
			ThresholdMinutes: 10,
			RequiredMisses:   3,
			HeartbeatMode:    "explicit",
		},
	}

	rm := adapter.buildReadModel(asset, "")
	if rm == nil || rm.HealthMonitor == nil {
		t.Fatal("expected non-nil read model + healthMonitor")
	}
	if rm.HealthMonitor.HeartbeatMode == nil {
		t.Fatal("HealthMonitor.HeartbeatMode is nil — should be set")
	}
	if got := *rm.HealthMonitor.HeartbeatMode; got != "explicit" {
		t.Fatalf("HealthMonitor.HeartbeatMode = %q, want %q", got, "explicit")
	}
}

// TestBuildReadModel_HeartbeatMode_Implicit asserts that an entity with the
// default empty-string HeartbeatMode is normalized to 'implicit' in the
// read model (via ResolvedMode), preserving the back-compat default for
// assets created before the field existed.
func TestBuildReadModel_HeartbeatMode_Implicit(t *testing.T) {
	adapter := &AssetStorageAdapter{}

	orgId, _ := model.ToObjectID("0000000000000000000aa001")
	asset := &entities.Asset{
		ID:        model.ObjectId(orgId),
		AssetUUID: "uuid-2",
		OrgID:     orgId,
		HealthMonitor: &entities.HealthMonitorConfig{
			Enabled:          true,
			ThresholdMinutes: 10,
			RequiredMisses:   3,
			// HeartbeatMode left empty — represents legacy assets stored in
			// Mongo before the field was introduced.
		},
	}

	rm := adapter.buildReadModel(asset, "")
	if rm == nil || rm.HealthMonitor == nil {
		t.Fatal("expected non-nil read model + healthMonitor")
	}
	if rm.HealthMonitor.HeartbeatMode == nil {
		t.Fatal("HealthMonitor.HeartbeatMode is nil — should default to 'implicit'")
	}
	if got := *rm.HealthMonitor.HeartbeatMode; got != "implicit" {
		t.Fatalf("HealthMonitor.HeartbeatMode = %q, want %q (back-compat default)", got, "implicit")
	}
}

// TestBuildReadModel_HeartbeatMode_NoHealthMonitor asserts that an entity
// without a HealthMonitor at all produces a read model with HealthMonitor=nil
// (no implicit promotion of the field for non-monitored assets).
func TestBuildReadModel_HeartbeatMode_NoHealthMonitor(t *testing.T) {
	adapter := &AssetStorageAdapter{}

	orgId, _ := model.ToObjectID("0000000000000000000aa001")
	asset := &entities.Asset{
		ID:        model.ObjectId(orgId),
		AssetUUID: "uuid-3",
		OrgID:     orgId,
		// HealthMonitor: nil
	}

	rm := adapter.buildReadModel(asset, "")
	if rm == nil {
		t.Fatal("expected non-nil read model")
	}
	if rm.HealthMonitor != nil {
		t.Fatalf("HealthMonitor should be nil for non-monitored asset; got %+v", rm.HealthMonitor)
	}
}
