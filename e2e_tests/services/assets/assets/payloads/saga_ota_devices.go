package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaOtaHttpDevice returns the minimal HTTP-protocol asset the ota_http journey
// targets: it sits on the OTA source template and is migrated to the firmware's
// target template by the plan. The HTTP OTA path is NOT presence-gated (the device
// polls GET /ota/jobs), so the HealthMonitor here is just the standard explicit
// config that keeps the asset a valid, well-formed connectivity asset — its online
// status is irrelevant to the poll flow.
//
// Inputs:
//   - runID         journey identifier embedded into name/assetUUID
//   - templateID    the OTA source AssetTemplate id (the plan's sourceTemplateId)
//   - routeGroupID  a route group id (the asset contract requires min=1)
func SagaOtaHttpDevice(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	uuid := fmt.Sprintf("saga-ota-http-%s", runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-ota-http-device-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       uuid,
			AssetTemplateID: templateID,
			RouteGroupIds:   []string{routeGroupID},
			Protocol: contracts.ProtocolType{
				Type: "http",
				Http: &contracts.NoneConfig{},
			},
			HealthMonitor: &contracts.HealthMonitorConfig{
				Enabled:          zerovalue.Ptr(true),
				ThresholdMinutes: zerovalue.Ptr(10),
				RequiredMisses:   zerovalue.Ptr(1),
				HeartbeatMode:    zerovalue.Ptr("explicit"),
			},
		},
	}
}

// SagaOtaMqttDevice returns the MQTT password-mode asset the ota_mqtt journey
// targets. The MQTT OTA path IS presence-gated — the reconciler pushes the
// ota_update only to an online device — so HealthMonitor is enabled in implicit
// mode (MQTT broker presence advisories drive online/offline) and a plaintext
// password is set so ConnectMqttPassword can authenticate the sim to the broker.
//
// Inputs:
//   - runID         journey identifier embedded into name/assetUUID
//   - templateID    the OTA source AssetTemplate id (the plan's sourceTemplateId)
//   - routeGroupID  a route group id (the asset contract requires min=1)
func SagaOtaMqttDevice(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	uuid := fmt.Sprintf("saga-ota-mqtt-%s", runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-ota-mqtt-device-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       uuid,
			AssetTemplateID: templateID,
			RouteGroupIds:   []string{routeGroupID},
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
				Enabled:          zerovalue.Ptr(true),
				ThresholdMinutes: zerovalue.Ptr(10),
				RequiredMisses:   zerovalue.Ptr(1),
				HeartbeatMode:    zerovalue.Ptr("implicit"),
			},
		},
	}
}
