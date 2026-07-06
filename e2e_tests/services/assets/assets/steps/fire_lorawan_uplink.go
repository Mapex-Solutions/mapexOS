package steps

import (
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/lorawansim"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// FireLorawanUplink sends one uplink from the labelled sensor over its gateway
// link. It records the send time (LorawanUplinkSentAtKey) immediately before, so
// the ingestion assert can scope its /events/raw search to "after this uplink".
//
// Reads (bag):
//   - LorawanSensorKey(sensorLabel)  *lorawansim.Sensor  set by JoinLorawanSensor
//
// Writes (bag):
//   - LorawanUplinkSentAtKey(sensorLabel)  time.Time
//
// Compensate: none — an uplink is fire-and-forget; nothing to undo.
func FireLorawanUplink(sensorLabel, payloadHex string) saga.Step {
	return saga.Step{
		Name: "assets/assets.FireLorawanUplink[" + sensorLabel + "]",
		Do: func(c *saga.Context) error {
			v, ok := c.Get(LorawanSensorKey(sensorLabel))
			if !ok {
				return fmt.Errorf("fire lorawan uplink %q: sensor not joined", sensorLabel)
			}
			s := v.(*lorawansim.Sensor)

			c.Set(LorawanUplinkSentAtKey(sensorLabel), time.Now().UTC())
			if err := s.SendUplink(2, payloadHex, false); err != nil {
				return fmt.Errorf("fire lorawan uplink %q: %w", sensorLabel, err)
			}
			return nil
		},
	}
}
