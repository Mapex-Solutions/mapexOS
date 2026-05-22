package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaMqttConnectivitySensor returns the MQTT password-mode asset
// payload the connectivity-action journey uses. Differs from
// SagaMqttTemperatureSensor in that HealthMonitor wires distinct
// online/offline route groups (the auth/connectivity journey reuses
// the same RG for both, which is fine for the MQTT plumbing test but
// blurs which transition fired the action).
//
// Inputs:
//   - runID         saga run id embedded into name/uuid
//   - templateID    AssetTemplate id created earlier in the saga
//   - onlineRGID    Mongo ObjectID hex of the RG fired on offline→online
//   - offlineRGID   Mongo ObjectID hex of the RG fired on online→offline
//
// Defaults:
//   - Protocol: MQTT password
//   - HealthMonitor: enabled, implicit mode (MQTT presence drives it
//     via broker advisories regardless of the operator setting; we
//     still set implicit so the field validator passes).
func SagaMqttConnectivitySensor(runID, templateID, onlineRGID, offlineRGID string) *AssetCreateBuilder {
	uuid := fmt.Sprintf("saga-mqtt-conn-%s", runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-mqtt-connectivity-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       uuid,
			AssetTemplateID: templateID,
			// Required by the asset contract — the connectivity journey
			// does not exercise the regular RouteGroupIds path, so reuse
			// the offline RG as a placeholder satisfying min=1.
			RouteGroupIds: []string{offlineRGID},
			Protocol: contracts.ProtocolType{
				Type: "mqtt",
				Mqtt: &contracts.MqttConfig{
					ClientId: uuid,
					Username: uuid,
					AuthType: "password",
					Password: SagaMqttDefaultPassword,
				},
			},
			HealthMonitor: &contracts.HealthMonitorConfig{
				Enabled:              zerovalue.Ptr(true),
				ThresholdMinutes:     zerovalue.Ptr(10),
				RequiredMisses:       zerovalue.Ptr(1),
				HeartbeatMode:        zerovalue.Ptr("implicit"),
				OnlineRouteGroupIds:  []string{onlineRGID},
				OfflineRouteGroupIds: []string{offlineRGID},
			},
		},
	}
}

// SagaHttpConnectivitySensor returns the HTTP-protocol asset payload
// the connectivity-action HTTP journey uses. The asset receives
// explicit heartbeats (POST /api/v1/heartbeat) which drives online,
// and the saga calls the assets internal force-offline endpoint to
// drive offline within the run budget.
//
// Inputs:
//   - runID         saga run id embedded into name/uuid
//   - templateID    AssetTemplate id created earlier in the saga
//   - onlineRGID    Mongo ObjectID hex of the RG fired on offline→online
//   - offlineRGID   Mongo ObjectID hex of the RG fired on online→offline
//
// Defaults:
//   - Protocol: HTTP
//   - HealthMonitor: explicit mode, thresholdMinutes=10 (contract min)
func SagaHttpConnectivitySensor(runID, templateID, onlineRGID, offlineRGID string) *AssetCreateBuilder {
	uuid := fmt.Sprintf("saga-http-conn-%s", runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-http-connectivity-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       uuid,
			AssetTemplateID: templateID,
			RouteGroupIds:   []string{offlineRGID},
			Protocol: contracts.ProtocolType{
				Type: "http",
				Http: &contracts.NoneConfig{},
			},
			HealthMonitor: &contracts.HealthMonitorConfig{
				Enabled:              zerovalue.Ptr(true),
				ThresholdMinutes:     zerovalue.Ptr(10),
				RequiredMisses:       zerovalue.Ptr(1),
				HeartbeatMode:        zerovalue.Ptr("explicit"),
				OnlineRouteGroupIds:  []string{onlineRGID},
				OfflineRouteGroupIds: []string{offlineRGID},
			},
		},
	}
}
