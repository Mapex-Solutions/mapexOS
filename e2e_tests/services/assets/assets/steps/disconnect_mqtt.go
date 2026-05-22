package steps

import (
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// DisconnectMqtt cleanly closes the MQTT handle stored on the bag.
// A graceful close (250ms quiesce) flushes any inflight publishes and
// triggers the broker's clean-disconnect path so the asset's
// healthStatus flips to "offline" through the normal presence flow.
//
// Reads (bag):
//   - BagKeyMqttClient  *mqttclient.Client  set by ConnectMqtt
//
// Writes (bag):
//   - BagKeyMqttDisconnectedAt  time.Time  captured before close so
//     downstream asserts can scope event searches by this value.
//
// Compensate: no-op. The connection has already been closed; the saga
// runner does not call Disconnect twice.
func DisconnectMqtt() saga.Step {
	return saga.Step{
		Name: "assets/assets.DisconnectMqtt",
		Do: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyMqttClient)
			if !ok {
				return nil
			}
			cli, ok := v.(*mqttclient.Client)
			if !ok {
				return nil
			}
			c.Set(BagKeyMqttDisconnectedAt, time.Now().UTC())
			cli.Disconnect(250)
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
