package services

import (
	"reflect"
	"testing"

	"router/src/modules/events/application/constants"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
)

func TestResolveRouteGroupIds(t *testing.T) {
	offlineIds := []string{"offline-1", "offline-2"}
	onlineIds := []string{"online-1"}
	defaultIds := []string{"default-1", "default-2"}

	assetWithHM := &assetsContract.AssetReadModel{
		RouteGroupIds: defaultIds,
		HealthMonitor: &assetsContract.HealthMonitorConfig{
			OfflineRouteGroupIds: offlineIds,
			OnlineRouteGroupIds:  onlineIds,
		},
	}
	assetWithoutHM := &assetsContract.AssetReadModel{
		RouteGroupIds: defaultIds,
		HealthMonitor: nil,
	}

	tests := []struct {
		name        string
		asset       *assetsContract.AssetReadModel
		eventSource string
		event       map[string]interface{}
		want        []string
	}{
		{
			name:        "assetEvent default returns RouteGroupIds",
			asset:       assetWithHM,
			eventSource: constants.EventSourceAssetEvent,
			event:       map[string]interface{}{},
			want:        defaultIds,
		},
		{
			name:        "healthStatus + offline returns OfflineRouteGroupIds",
			asset:       assetWithHM,
			eventSource: constants.EventSourceHealthStatus,
			event:       map[string]interface{}{"eventType": constants.HealthStatusOffline},
			want:        offlineIds,
		},
		{
			name:        "healthStatus + online returns OnlineRouteGroupIds",
			asset:       assetWithHM,
			eventSource: constants.EventSourceHealthStatus,
			event:       map[string]interface{}{"eventType": constants.HealthStatusOnline},
			want:        onlineIds,
		},
		{
			name:        "healthStatus + unknown eventType returns nil",
			asset:       assetWithHM,
			eventSource: constants.EventSourceHealthStatus,
			event:       map[string]interface{}{"eventType": "foobar"},
			want:        nil,
		},
		{
			name:        "healthStatus + nil HealthMonitor returns nil",
			asset:       assetWithoutHM,
			eventSource: constants.EventSourceHealthStatus,
			event:       map[string]interface{}{"eventType": constants.HealthStatusOffline},
			want:        nil,
		},
		{
			name:        "empty eventSource falls through to RouteGroupIds",
			asset:       assetWithHM,
			eventSource: "",
			event:       map[string]interface{}{},
			want:        defaultIds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveRouteGroupIds(tt.asset, tt.eventSource, tt.event)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("resolveRouteGroupIds(%s) = %v, want %v", tt.eventSource, got, tt.want)
			}
		})
	}
}

func TestIsAllowedKindForHealthStatus(t *testing.T) {
	tests := []struct {
		kind string
		want bool
	}{
		{kind: "trigger", want: true},
		{kind: "workflow", want: true},
		{kind: "save_event", want: false},
		{kind: "notification", want: false},
		{kind: "lake_house", want: false},
		{kind: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			if got := isAllowedKindForHealthStatus(tt.kind); got != tt.want {
				t.Fatalf("isAllowedKindForHealthStatus(%q) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}
