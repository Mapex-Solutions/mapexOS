package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/lorawansim"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
)

// JoinLorawanSensor creates a simulated end-device on the labelled gateway using
// the exact identity/keys the sensor asset was provisioned with — the DevEUI is
// derived the same way (payloads.SensorDevEUI) so the LNS matches the device — then
// performs the OTAA join against the LNS.
//
// Reads (bag):
//   - LorawanGatewayKey(gwLabel)  *lorawansim.Gateway  set by ConnectLorawanGateway*
//
// Writes (bag):
//   - LorawanSensorKey(sensorLabel)  *lorawansim.Sensor
//
// Compensate: none — the sensor rides the gateway link, released when the gateway
// closes; there is no separate server-side sensor resource to tear down here.
func JoinLorawanSensor(gwLabel, sensorLabel string) saga.Step {
	return saga.Step{
		Name: "assets/assets.JoinLorawanSensor[" + sensorLabel + "]",
		Do: func(c *saga.Context) error {
			v, ok := c.Get(LorawanGatewayKey(gwLabel))
			if !ok {
				return fmt.Errorf("join lorawan sensor %q: gateway %q not connected", sensorLabel, gwLabel)
			}
			gw := v.(*lorawansim.Gateway)

			s, err := gw.NewSensor(lorawansim.SensorConfig{
				DevEUI:     payloads.SensorDevEUI(c.RunID, sensorLabel),
				JoinEUI:    "0000000000000000",
				AppKey:     payloads.SagaSensorAppKey,
				Region:     "EU868",
				Class:      "A",
				MacVersion: "1.0.3",
				Activation: "otaa",
			})
			if err != nil {
				return fmt.Errorf("build lorawan sensor %q: %w", sensorLabel, err)
			}
			if err := s.Join(c.Stdctx); err != nil {
				return fmt.Errorf("join lorawan sensor %q: %w", sensorLabel, err)
			}
			c.Set(LorawanSensorKey(sensorLabel), s)
			return nil
		},
	}
}
