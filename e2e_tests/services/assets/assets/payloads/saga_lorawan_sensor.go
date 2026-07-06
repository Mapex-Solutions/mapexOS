package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// LoRaWAN sensor key material the saga fixtures use. The values are stable across
// runs (deterministic reruns) and shared with the sim-driving steps so the
// simulated device presents the exact identity the asset was provisioned with.
const (
	// SagaSensorAppKey is the OTAA application root key (32 hex / 16 bytes).
	SagaSensorAppKey = "00112233445566778899AABBCCDDEEFF"
	// SagaSensorABPKey is the ABP session key reused for NwkSKey and AppSKey.
	SagaSensorABPKey = "00112233445566778899AABBCCDDEEFF"
	// SagaSensorDevAddr is the ABP fixed device address (8 hex / 4 bytes).
	SagaSensorDevAddr = "26011BDA"
)

// SensorDevEUI derives a stable uppercase 16-hex DevEUI unique to a (runID, label)
// pair, so a journey can provision more than one sensor without a DevEUI collision
// on the LNS. The JoinLorawanSensor step recomputes the same value to build the
// matching simulated device.
func SensorDevEUI(runID, label string) string {
	return euiFromSeed(runID + "-sensor-" + label)
}

// lorawanHealthMonitor is the health config every saga LoRaWAN asset carries so the
// healthmonitor tracks its presence (it drops advisories for assets whose
// HealthMonitor is not active): implicit mode (each uplink / gateway connect is a
// presence signal → online), a 10-minute anti-flap floor (the contract's min), and
// one missed window to flip offline. Shared by the sensor and gateway builders.
func lorawanHealthMonitor() *contracts.HealthMonitorConfig {
	return &contracts.HealthMonitorConfig{
		Enabled:          zerovalue.Ptr(true),
		ThresholdMinutes: zerovalue.Ptr(10),
		RequiredMisses:   zerovalue.Ptr(1),
		HeartbeatMode:    zerovalue.Ptr("implicit"),
	}
}

// SagaLorawanSensorOTAAFor returns a LoRaWAN end-device (kind=device) OTAA asset
// builder scoped to a label. AssetUUID is the deterministic public identity the
// presence/ingestion pipeline keys on (the raw event threadId); DevEUI is the
// separate radio identity. AppKey is request-only — the platform KEK-encrypts it.
func SagaLorawanSensorOTAAFor(label string) func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	return func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
		return &AssetCreateBuilder{
			spec: contracts.AssetCreate{
				Name:            fmt.Sprintf("saga-lorawan-sensor-%s-%s", label, runID),
				Enabled:         true,
				DebugEnabled:    true,
				AssetUUID:       fmt.Sprintf("saga-lorawan-sensor-%s-%s", label, runID),
				AssetTemplateID: templateID,
				RouteGroupIds:   []string{routeGroupID},
				Protocol: contracts.ProtocolType{
					Type: "lorawan",
					Lorawan: &contracts.LorawanConfig{
						Kind:       contracts.LorawanKindDevice,
						DevEUI:     SensorDevEUI(runID, label),
						JoinEUI:    "0000000000000000",
						Region:     "EU868",
						Class:      "A",
						MacVersion: "1.0.3",
						PhyVersion: "1.0.3",
						Activation: contracts.LorawanActivationOTAA,
						AppKey:     SagaSensorAppKey,
					},
				},
				HealthMonitor: lorawanHealthMonitor(),
			},
		}
	}
}

// SagaLorawanSensorABPFor is the ABP analogue: a fixed session (DevAddr + NwkSKey +
// AppSKey) with no OTAA join. Same label-scoped identity derivation.
func SagaLorawanSensorABPFor(label string) func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	return func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
		return &AssetCreateBuilder{
			spec: contracts.AssetCreate{
				Name:            fmt.Sprintf("saga-lorawan-sensor-abp-%s-%s", label, runID),
				Enabled:         true,
				DebugEnabled:    true,
				AssetUUID:       fmt.Sprintf("saga-lorawan-sensor-abp-%s-%s", label, runID),
				AssetTemplateID: templateID,
				RouteGroupIds:   []string{routeGroupID},
				Protocol: contracts.ProtocolType{
					Type: "lorawan",
					Lorawan: &contracts.LorawanConfig{
						Kind:       contracts.LorawanKindDevice,
						DevEUI:     SensorDevEUI(runID, label),
						Region:     "EU868",
						Class:      "A",
						MacVersion: "1.0.3",
						PhyVersion: "1.0.3",
						Activation: contracts.LorawanActivationABP,
						DevAddr:    SagaSensorDevAddr,
						NwkSKey:    SagaSensorABPKey,
						AppSKey:    SagaSensorABPKey,
					},
				},
				HealthMonitor: lorawanHealthMonitor(),
			},
		}
	}
}
