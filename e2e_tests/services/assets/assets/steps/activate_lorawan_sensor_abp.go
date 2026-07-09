package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/lorawansim"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
)

// ActivateLorawanSensorABP creates a simulated ABP end-device on the labelled
// gateway using the fixed session (DevAddr + NwkSKey + AppSKey) the sensor asset
// was provisioned with. DevEUI and DevAddr are derived the same way the payload
// derives them (payloads.SensorDevEUI / payloads.SensorDevAddr) so the LNS matches
// the device. ABP needs no OTAA join: the device is pre-activated, so NewSensor
// binds its downlink handler at once and it can send uplinks immediately.
//
// Reads (bag):
//   - LorawanGatewayKey(gwLabel)  *lorawansim.Gateway  set by ConnectLorawanGateway*
//
// Writes (bag):
//   - LorawanSensorKey(sensorLabel)  *lorawansim.Sensor
//
// Compensate: none — the sensor rides the gateway link, released when the gateway
// closes; there is no separate server-side sensor resource to tear down here.
func ActivateLorawanSensorABP(gwLabel, sensorLabel string) saga.Step {
	return saga.Step{
		Name: "assets/assets.ActivateLorawanSensorABP[" + sensorLabel + "]",
		Do: func(c *saga.Context) error {
			v, ok := c.Get(LorawanGatewayKey(gwLabel))
			if !ok {
				return fmt.Errorf("activate lorawan abp sensor %q: gateway %q not connected", sensorLabel, gwLabel)
			}
			gw := v.(*lorawansim.Gateway)

			s, err := gw.NewSensor(lorawansim.SensorConfig{
				DevEUI:     payloads.SensorDevEUI(c.RunID, sensorLabel),
				JoinEUI:    "0000000000000000",
				Region:     "EU868",
				Class:      "A",
				MacVersion: "1.0.3",
				Activation: "abp",
				DevAddr:    payloads.SensorDevAddr(c.RunID, sensorLabel),
				NwkSKey:    payloads.SagaSensorABPKey,
				AppSKey:    payloads.SagaSensorABPKey,
			})
			if err != nil {
				return fmt.Errorf("build lorawan abp sensor %q: %w", sensorLabel, err)
			}
			c.Set(LorawanSensorKey(sensorLabel), s)
			return nil
		},
	}
}
