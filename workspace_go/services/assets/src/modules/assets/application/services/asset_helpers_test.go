package services

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"assets/src/modules/assets/domain/entities"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// fakeRouteGroupPort is an inline mock of ports.RouteGroupPort shared between
// asset_service_test.go and asset_helpers_test.go.
type fakeRouteGroupPort struct {
	kindsByGroup map[string][]string
	kindsErr     error
	namesByIdsFn func(ctx context.Context, ids []string) ([]string, error)
}

func (f *fakeRouteGroupPort) GetNamesByIds(ctx context.Context, ids []string) ([]string, error) {
	if f.namesByIdsFn != nil {
		return f.namesByIdsFn(ctx, ids)
	}
	return nil, nil
}

func (f *fakeRouteGroupPort) GetRouterKindsByIds(ctx context.Context, ids []string) (map[string][]string, error) {
	if f.kindsErr != nil {
		return nil, f.kindsErr
	}
	return f.kindsByGroup, nil
}

func boolPtr(b bool) *bool { return &b }

func strPtr(s string) *string { return &s }

// lorawanProtocol builds a LoRaWAN ProtocolType for the given device/gateway
// kind, used to exercise the heartbeat-mode guard.
func lorawanProtocol(kind string) *contracts.ProtocolType {
	return &contracts.ProtocolType{
		Type:    "lorawan",
		Lorawan: &contracts.LorawanConfig{Kind: kind},
	}
}

func TestValidateHealthMonitorConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		protocol     *contracts.ProtocolType
		hm           *contracts.HealthMonitorConfig
		kindsByGroup map[string][]string
		kindsErr     error
		wantErr      bool
		wantCode     int
		wantMsgSub   string // substring expected in the error message (first entry)
	}{
		{
			name:    "nil HealthMonitor passes",
			hm:      nil,
			wantErr: false,
		},
		{
			name:       "LoRaWAN device with explicit heartbeat fails 422",
			protocol:   lorawanProtocol(contracts.LorawanKindDevice),
			hm:         &contracts.HealthMonitorConfig{Enabled: boolPtr(true), HeartbeatMode: strPtr("explicit")},
			wantErr:    true,
			wantCode:   status.UNPROCESSABLE_ENTITY,
			wantMsgSub: "implicit heartbeat only",
		},
		{
			name:     "LoRaWAN device with implicit heartbeat passes",
			protocol: lorawanProtocol(contracts.LorawanKindDevice),
			hm:       &contracts.HealthMonitorConfig{Enabled: boolPtr(true), HeartbeatMode: strPtr("implicit")},
			wantErr:  false,
		},
		{
			name:     "LoRaWAN gateway with explicit heartbeat passes (LNS drives presence)",
			protocol: lorawanProtocol(contracts.LorawanKindGateway),
			hm:       &contracts.HealthMonitorConfig{Enabled: boolPtr(true), HeartbeatMode: strPtr("explicit")},
			wantErr:  false,
		},
		{
			name:    "Enabled is nil passes",
			hm:      &contracts.HealthMonitorConfig{},
			wantErr: false,
		},
		{
			name:    "Enabled=false passes regardless of arrays",
			hm:      &contracts.HealthMonitorConfig{Enabled: boolPtr(false)},
			wantErr: false,
		},
		{
			name: "Enabled=true with empty arrays is valid (monitor-only)",
			hm: &contracts.HealthMonitorConfig{
				Enabled: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "Enabled=true with offline trigger passes",
			hm: &contracts.HealthMonitorConfig{
				Enabled:              boolPtr(true),
				OfflineRouteGroupIds: []string{"rg-off"},
			},
			kindsByGroup: map[string][]string{
				"rg-off": {"trigger"},
			},
			wantErr: false,
		},
		{
			name: "Enabled=true with online workflow passes",
			hm: &contracts.HealthMonitorConfig{
				Enabled:             boolPtr(true),
				OnlineRouteGroupIds: []string{"rg-on"},
			},
			kindsByGroup: map[string][]string{
				"rg-on": {"workflow"},
			},
			wantErr: false,
		},
		{
			name: "Enabled=true with save_event router fails 422",
			hm: &contracts.HealthMonitorConfig{
				Enabled:              boolPtr(true),
				OfflineRouteGroupIds: []string{"rg-off"},
			},
			kindsByGroup: map[string][]string{
				"rg-off": {"trigger", "save_event"},
			},
			wantErr:    true,
			wantCode:   status.UNPROCESSABLE_ENTITY,
			wantMsgSub: "save_event",
		},
		{
			name: "Enabled=true with notification router fails 422",
			hm: &contracts.HealthMonitorConfig{
				Enabled:             boolPtr(true),
				OnlineRouteGroupIds: []string{"rg-on"},
			},
			kindsByGroup: map[string][]string{
				"rg-on": {"notification"},
			},
			wantErr:    true,
			wantCode:   status.UNPROCESSABLE_ENTITY,
			wantMsgSub: "notification",
		},
		{
			name: "RouteGroupPort critical failure returns wrapped error",
			hm: &contracts.HealthMonitorConfig{
				Enabled:              boolPtr(true),
				OfflineRouteGroupIds: []string{"rg-off"},
			},
			kindsErr: errors.New("router HTTP 500"),
			wantErr:  true,
			// wantCode left 0 — wrapped error is not a ServerCustomError.
			wantMsgSub: "router HTTP 500",
		},
		{
			name: "Missing group in port response skips validation (no error)",
			hm: &contracts.HealthMonitorConfig{
				Enabled:              boolPtr(true),
				OfflineRouteGroupIds: []string{"rg-unknown"},
			},
			kindsByGroup: map[string][]string{}, // port returned empty map
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := &fakeRouteGroupPort{kindsByGroup: tt.kindsByGroup, kindsErr: tt.kindsErr}
			err := validateHealthMonitorConfig(ctx, port, tt.protocol, tt.hm)

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.wantCode != 0 {
				var sce *customErrors.ServerCustomError
				if !errors.As(err, &sce) {
					t.Fatalf("expected *ServerCustomError, got %T: %v", err, err)
				}
				if sce.Code != tt.wantCode {
					t.Fatalf("expected code %d, got %d", tt.wantCode, sce.Code)
				}
				if len(sce.Errors) == 0 || !strings.Contains(sce.Errors[0], tt.wantMsgSub) {
					t.Fatalf("expected error containing %q, got %v", tt.wantMsgSub, sce.Errors)
				}
			} else if !strings.Contains(err.Error(), tt.wantMsgSub) {
				t.Fatalf("expected error containing %q, got %v", tt.wantMsgSub, err)
			}
		})
	}
}

func TestConvertHealthMonitor(t *testing.T) {
	t.Run("nil entity returns nil contract", func(t *testing.T) {
		if got := convertHealthMonitor(nil); got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("value-only entity wraps fields as pointers", func(t *testing.T) {
		got := convertHealthMonitor(&entities.HealthMonitorConfig{
			Enabled:          true,
			ThresholdMinutes: 10,
			RequiredMisses:   3,
		})
		if got == nil {
			t.Fatal("expected non-nil result")
		}
		if got.Enabled == nil || *got.Enabled != true {
			t.Errorf("Enabled = %v, want *true", got.Enabled)
		}
		if got.ThresholdMinutes == nil || *got.ThresholdMinutes != 10 {
			t.Errorf("ThresholdMinutes = %v, want *10", got.ThresholdMinutes)
		}
		if got.RequiredMisses == nil || *got.RequiredMisses != 3 {
			t.Errorf("RequiredMisses = %v, want *3", got.RequiredMisses)
		}
	})

	t.Run("empty route group slices are left nil (omitempty)", func(t *testing.T) {
		got := convertHealthMonitor(&entities.HealthMonitorConfig{
			Enabled: true,
		})
		if got.OfflineRouteGroupIds != nil {
			t.Errorf("OfflineRouteGroupIds = %v, want nil", got.OfflineRouteGroupIds)
		}
		if got.OnlineRouteGroupIds != nil {
			t.Errorf("OnlineRouteGroupIds = %v, want nil", got.OnlineRouteGroupIds)
		}
	})

	t.Run("populated route group slices are deep-copied", func(t *testing.T) {
		offline := []string{"rg-offline-1", "rg-offline-2"}
		online := []string{"rg-online-1"}
		entity := &entities.HealthMonitorConfig{
			Enabled:              true,
			OfflineRouteGroupIds: offline,
			OnlineRouteGroupIds:  online,
		}

		got := convertHealthMonitor(entity)

		if !reflect.DeepEqual(got.OfflineRouteGroupIds, offline) {
			t.Errorf("OfflineRouteGroupIds = %v, want %v", got.OfflineRouteGroupIds, offline)
		}
		if !reflect.DeepEqual(got.OnlineRouteGroupIds, online) {
			t.Errorf("OnlineRouteGroupIds = %v, want %v", got.OnlineRouteGroupIds, online)
		}

		// Mutating the source must NOT affect the contract copy.
		entity.OfflineRouteGroupIds[0] = "mutated"
		if got.OfflineRouteGroupIds[0] != "rg-offline-1" {
			t.Fatal("converted slice shares backing array with source — append([]string(nil), ...) expected")
		}
	})

	t.Run("Enabled=false still produces *bool(false), not nil", func(t *testing.T) {
		got := convertHealthMonitor(&entities.HealthMonitorConfig{
			Enabled: false,
		})
		if got.Enabled == nil || *got.Enabled != false {
			t.Errorf("Enabled = %v, want *false", got.Enabled)
		}
	})

	_ = &contracts.HealthMonitorConfig{} // ensure contracts import is used
}

func TestValidateAssetTopologyRequirements(t *testing.T) {
	tests := []struct {
		name       string
		protocol   *contracts.ProtocolType
		templateID string
		routeGroup []string
		wantErr    bool
	}{
		{
			name:       "lorawan gateway is exempt with no template and no route group",
			protocol:   lorawanProtocol(contracts.LorawanKindGateway),
			templateID: "",
			routeGroup: nil,
		},
		{
			name:       "lorawan device with template and a route group passes",
			protocol:   lorawanProtocol(contracts.LorawanKindDevice),
			templateID: "691bb4071e717d77a2430b46",
			routeGroup: []string{"rg-1"},
		},
		{
			name:       "lorawan device without template denies",
			protocol:   lorawanProtocol(contracts.LorawanKindDevice),
			templateID: "",
			routeGroup: []string{"rg-1"},
			wantErr:    true,
		},
		{
			name:       "mqtt without a route group denies",
			protocol:   &contracts.ProtocolType{Type: "mqtt"},
			templateID: "691bb4071e717d77a2430b46",
			routeGroup: nil,
			wantErr:    true,
		},
		{
			name:       "http without template denies",
			protocol:   &contracts.ProtocolType{Type: "http"},
			templateID: "",
			routeGroup: []string{"rg-1"},
			wantErr:    true,
		},
		{
			name:       "mqtt with template and a route group passes",
			protocol:   &contracts.ProtocolType{Type: "mqtt"},
			templateID: "691bb4071e717d77a2430b46",
			routeGroup: []string{"rg-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAssetTopologyRequirements(tt.protocol, tt.templateID, tt.routeGroup)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}
}

func TestValidateAssetAttributes(t *testing.T) {
	geo := func(lat, lon any) map[string]any { return map[string]any{"lat": lat, "lon": lon} }

	tests := []struct {
		name    string
		attrs   []contracts.AssetAttribute
		wantErr bool
	}{
		{
			name:    "nil list is valid",
			attrs:   nil,
			wantErr: false,
		},
		{
			name: "one of each kind is valid",
			attrs: []contracts.AssetAttribute{
				{Label: "site", Kind: contracts.AssetAttributeKindString, Value: "Lisbon-DC1"},
				{Label: "floor", Kind: contracts.AssetAttributeKindInteger, Value: float64(42)},
				{Label: "critical", Kind: contracts.AssetAttributeKindBoolean, Value: true},
				{Label: "installed_at", Kind: contracts.AssetAttributeKindDate, Value: "2027-01-01T00:00:00Z"},
				{Label: "location", Kind: contracts.AssetAttributeKindGeo, Value: geo(38.7, -9.1)},
			},
			wantErr: false,
		},
		{
			name:    "integer with fractional value is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "n", Kind: contracts.AssetAttributeKindInteger, Value: float64(1.5)}},
			wantErr: true,
		},
		{
			name:    "integer with string value is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "n", Kind: contracts.AssetAttributeKindInteger, Value: "abc"}},
			wantErr: true,
		},
		{
			name:    "date not RFC3339 is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "d", Kind: contracts.AssetAttributeKindDate, Value: "2027-13-40"}},
			wantErr: true,
		},
		{
			name:    "geo missing lon is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "g", Kind: contracts.AssetAttributeKindGeo, Value: map[string]any{"lat": float64(10)}}},
			wantErr: true,
		},
		{
			name:    "geo lat out of range is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "g", Kind: contracts.AssetAttributeKindGeo, Value: geo(float64(100), float64(0))}},
			wantErr: true,
		},
		{
			name:    "boolean with non-bool value is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "b", Kind: contracts.AssetAttributeKindBoolean, Value: "true"}},
			wantErr: true,
		},
		{
			name: "duplicate labels are rejected",
			attrs: []contracts.AssetAttribute{
				{Label: "dup", Kind: contracts.AssetAttributeKindString, Value: "a"},
				{Label: "dup", Kind: contracts.AssetAttributeKindString, Value: "b"},
			},
			wantErr: true,
		},
		{
			name:    "reserved mapex prefix is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "mapex.owner", Kind: contracts.AssetAttributeKindString, Value: "x"}},
			wantErr: true,
		},
		{
			name:    "label with space fails the pattern",
			attrs:   []contracts.AssetAttribute{{Label: "has space", Kind: contracts.AssetAttributeKindString, Value: "x"}},
			wantErr: true,
		},
		{
			name:    "unknown kind is rejected",
			attrs:   []contracts.AssetAttribute{{Label: "x", Kind: "json", Value: "x"}},
			wantErr: true,
		},
		{
			name:    "more than max attributes is rejected",
			attrs:   makeAttributes(contracts.AssetAttributeMaxCount + 1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAssetAttributes(tt.attrs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateAssetAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// makeAttributes builds n valid string attributes with unique labels.
func makeAttributes(n int) []contracts.AssetAttribute {
	attrs := make([]contracts.AssetAttribute, n)
	for i := range attrs {
		attrs[i] = contracts.AssetAttribute{
			Label: "attr_" + strconv.Itoa(i),
			Kind:  contracts.AssetAttributeKindString,
			Value: "v",
		}
	}
	return attrs
}
