package services

import (
	"context"
	"errors"
	"testing"

	assetPorts "assets/src/modules/assets/application/ports"
)

// newAdminService composes a HealthMonitorService with the three mocked
// ports and the throwaway metrics fixture shared with the heartbeat
// suite. Kept separate so admin tests do not borrow heartbeat-specific
// service plumbing if it ever diverges.
func newAdminService(
	healthRepo *mockHealthRepo,
	assetRepo *mockAssetRepo,
	publisher *mockAlertPublisher,
) *HealthMonitorService {
	return newHeartbeatService(healthRepo, assetRepo, publisher)
}

// TestForceOfflineByAssetUUID drives the admin offline-transition
// endpoint through every short-circuit the service enforces, plus the
// happy path that fires the offline alert publish. Each subtest seeds
// the three port mocks and asserts the exact set of downstream calls.
func TestForceOfflineByAssetUUID(t *testing.T) {
	const (
		orgHex    = "68f5bbce1aef22967c3ebb30"
		assetUUID = "saga-uuid-admin"
	)

	enabledAsset := func() *assetPorts.Asset {
		return monitoredAsset(orgHex, assetUUID, &assetPorts.HealthMonitorConfig{
			Enabled:              true,
			ThresholdMinutes:     10,
			RequiredMisses:       1,
			OfflineRouteGroupIds: []string{"rg-offline-1"},
			OnlineRouteGroupIds:  []string{"rg-online-1"},
		})
	}

	tests := []struct {
		name                    string
		reason                  string
		repoErr                 error
		repoAsset               *assetPorts.Asset
		isAlerted               bool
		expectFindCalls         int
		expectMarkAlerted       int
		expectUpdateStatusCalls int
		expectUpdateStatusValue string
		expectPublishOffline    int
		expectRouteGroupIds     []string
	}{
		{
			name:            "asset_not_found_no_side_effects",
			reason:          "ci",
			repoAsset:       nil,
			expectFindCalls: 1,
		},
		{
			name:            "asset_lookup_error_no_side_effects",
			reason:          "ci",
			repoErr:         errors.New("mongo down"),
			expectFindCalls: 1,
		},
		{
			name:            "disabled_asset_skipped",
			reason:          "ci",
			repoAsset:       monitoredAsset(orgHex, assetUUID, &assetPorts.HealthMonitorConfig{Enabled: false}),
			expectFindCalls: 1,
		},
		{
			name:            "already_alerted_is_idempotent",
			reason:          "ci",
			repoAsset:       enabledAsset(),
			isAlerted:       true,
			expectFindCalls: 1,
		},
		{
			name:                    "happy_path_fires_offline_publish_with_route_groups",
			reason:                  "ci",
			repoAsset:               enabledAsset(),
			expectFindCalls:         1,
			expectMarkAlerted:       1,
			expectUpdateStatusCalls: 1,
			expectUpdateStatusValue: "offline",
			expectPublishOffline:    1,
			expectRouteGroupIds:     []string{"rg-offline-1"},
		},
		{
			name:                    "empty_reason_defaults_and_still_fires",
			reason:                  "",
			repoAsset:               enabledAsset(),
			expectFindCalls:         1,
			expectMarkAlerted:       1,
			expectUpdateStatusCalls: 1,
			expectUpdateStatusValue: "offline",
			expectPublishOffline:    1,
			expectRouteGroupIds:     []string{"rg-offline-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			healthRepo := &mockHealthRepo{isAlertedReturn: tt.isAlerted}
			publisher := &mockAlertPublisher{}
			assetRepo := &mockAssetRepo{asset: tt.repoAsset, findErr: tt.repoErr}

			svc := newAdminService(healthRepo, assetRepo, publisher)

			if err := svc.ForceOfflineByAssetUUID(context.Background(), assetUUID, tt.reason); err != nil {
				t.Fatalf("ForceOfflineByAssetUUID: unexpected error: %v", err)
			}

			if assetRepo.findByAssetUUIDCalls != tt.expectFindCalls {
				t.Errorf("FindByAssetUUID calls: want %d, got %d", tt.expectFindCalls, assetRepo.findByAssetUUIDCalls)
			}
			if healthRepo.markAlertedCalls != tt.expectMarkAlerted {
				t.Errorf("MarkAlerted calls: want %d, got %d", tt.expectMarkAlerted, healthRepo.markAlertedCalls)
			}
			if assetRepo.updateHealthStatusCalls != tt.expectUpdateStatusCalls {
				t.Errorf("UpdateHealthStatusWithChangedAt calls: want %d, got %d",
					tt.expectUpdateStatusCalls, assetRepo.updateHealthStatusCalls)
			}
			if tt.expectUpdateStatusValue != "" && assetRepo.updateHealthStatusLastStatus != tt.expectUpdateStatusValue {
				t.Errorf("UpdateHealthStatusWithChangedAt status arg: want %q, got %q",
					tt.expectUpdateStatusValue, assetRepo.updateHealthStatusLastStatus)
			}
			if publisher.publishOfflineCalls != tt.expectPublishOffline {
				t.Errorf("PublishOffline calls: want %d, got %d", tt.expectPublishOffline, publisher.publishOfflineCalls)
			}
			if tt.expectPublishOffline > 0 {
				if publisher.lastOfflineEvent == nil {
					t.Fatal("lastOfflineEvent: expected captured event, got nil")
				}
				if !equalStrings(publisher.lastOfflineEvent.RouteGroupIds, tt.expectRouteGroupIds) {
					t.Errorf("RouteGroupIds on offline event: want %v, got %v",
						tt.expectRouteGroupIds, publisher.lastOfflineEvent.RouteGroupIds)
				}
				if publisher.lastOfflineEvent.AssetUUID != assetUUID {
					t.Errorf("AssetUUID on offline event: want %q, got %q",
						assetUUID, publisher.lastOfflineEvent.AssetUUID)
				}
				if publisher.lastOfflineEvent.OrgId != orgHex {
					t.Errorf("OrgId on offline event: want %q, got %q",
						orgHex, publisher.lastOfflineEvent.OrgId)
				}
			}
			if publisher.publishOnlineCalls != 0 {
				t.Errorf("PublishOnline MUST NOT be called from force-offline path, got %d", publisher.publishOnlineCalls)
			}
		})
	}
}

// equalStrings is a small slice comparator that treats nil and empty
// slices as equal — the publisher does not normalize between the two
// and either is acceptable on the wire.
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
