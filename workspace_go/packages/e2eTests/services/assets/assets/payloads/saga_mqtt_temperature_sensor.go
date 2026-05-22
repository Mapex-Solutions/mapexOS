// Package payloads holds canonical AssetCreate fixtures for the
// assets/assets module.
//
// Assets bind together a template, route groups, and a transport protocol.
// The canonical SagaMqttTemperatureSensor fixture wires an MQTT-protocol
// asset attached to the saga temperature template and the saga save_event
// route group, providing the simplest end-to-end IoT happy-path the saga
// suite exercises.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaMqttDefaultPassword is the canonical plaintext password the saga
// fixture sends on create. The platform bcrypt-hashes it server-side;
// the connect step retrieves the same plaintext from the bag to present
// to the broker. Long enough to clear the contract's min=8 validator,
// stable across runs so reruns are deterministic.
const SagaMqttDefaultPassword = "saga-mqtt-temp-pwd-32chars-aaaa"

// AssetCreateBuilder wraps the contract DTO with fluent overrides journeys
// use when they need a tailored variant.
type AssetCreateBuilder struct {
	spec contracts.AssetCreate
}

// Build returns the contracts payload ready for POST /api/v1/assets.
func (b *AssetCreateBuilder) Build() contracts.AssetCreate { return b.spec }

// MqttPassword returns the plaintext password the builder embedded in
// the spec, so the create step can publish it on the bag for the
// connect step to consume.
func (b *AssetCreateBuilder) MqttPassword() string {
	if b.spec.Protocol.Mqtt == nil {
		return ""
	}
	return b.spec.Protocol.Mqtt.Password
}

// WithName overrides the asset name. Tests that assert by name override
// the runID-stamped default.
func (b *AssetCreateBuilder) WithName(name string) *AssetCreateBuilder {
	b.spec.Name = name
	return b
}

// WithRouteGroups replaces the route groups attached to the asset. Most
// callers will inject the route group id created earlier in the saga via
// a helper that reads the bag, but this method is here for tests that
// need to validate multi-route-group behavior.
func (b *AssetCreateBuilder) WithRouteGroups(ids ...string) *AssetCreateBuilder {
	b.spec.RouteGroupIds = ids
	return b
}

// WithMqttPassword overrides the plaintext password on the spec. Tests
// that exercise rotation set a different password to check the new
// hash propagates to the broker auth callout cache.
func (b *AssetCreateBuilder) WithMqttPassword(pwd string) *AssetCreateBuilder {
	if b.spec.Protocol.Mqtt != nil {
		b.spec.Protocol.Mqtt.Password = pwd
	}
	return b
}

// SagaMqttTemperatureSensor returns the canonical MQTT temperature sensor
// payload bound to the supplied template id and route group id.
//
// Inputs:
//   - runID         journey identifier embedded into name and assetUUID
//   - templateID    id of an asset template that exists in the saga org
//   - routeGroupID  id of a route group that exists in the saga org
//
// Defaults:
//   - Protocol:      MQTT (clientId/username seeded from runID; password
//                    is the canonical SagaMqttDefaultPassword which the
//                    platform bcrypt-hashes server-side)
//   - Enabled:       true
//   - HealthMonitor: explicit mode, threshold 10 min, 1 missed = offline
//
// Note: assetUUID drives the MQTT identity exposed in the broker plugin
// presence advisories. Keeping it deterministic (runID-based) lets the
// healthmonitor saga match on the assetUUID when verifying cache
// population and offline transitions.
func SagaMqttTemperatureSensor(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	uuid := fmt.Sprintf("saga-mqtt-temp-%s", runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-mqtt-temperature-%s", runID),
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
				Enabled: zerovalue.Ptr(true),
				// ThresholdMinutes carries a min=10 anti-flap floor on the
				// contract; the saga uses the floor so payloads pass DTO
				// validation while keeping the offline window short.
				ThresholdMinutes: zerovalue.Ptr(10),
				RequiredMisses:   zerovalue.Ptr(1),
				HeartbeatMode:    zerovalue.Ptr("implicit"),
			},
		},
	}
}

// SagaMqttCertTemperatureSensor returns the MQTT cert-mode variant of
// the canonical temperature sensor. AuthType=cert + no plaintext
// password — the device cert is issued in a separate step
// (POST /api/v1/mqtt_certs) after asset create. CertTTL is short (1
// day) so the journey doesn't lean on the platform-wide default and
// the cert metadata round-trip is exercised end-to-end.
func SagaMqttCertTemperatureSensor(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	uuid := fmt.Sprintf("saga-mqtt-cert-%s", runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-mqtt-cert-temperature-%s", runID),
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
					AuthType: "cert",
					CertTTL:  &contracts.CertTTLConfig{Value: 1, Unit: "day"},
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
