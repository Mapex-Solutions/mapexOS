package steps

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// PublishTelemetry sends a deterministic JSON payload over the active
// MQTT connection on topic "events/<assetUUID>/temperature". The broker
// plugin's ACL accepts only the bare-assetUUID topic shape; orgId no
// longer travels on the wire (the auth projection's orgId is the
// trust anchor server-side). The events service's events_raw consumer
// lifts the message into the pipeline.
//
// The payload carries the saga's runID so an assert that polls the
// events HTTP API can disambiguate this telemetry from any other event
// flowing through the stack at the same time.
//
// Reads (bag):
//   - BagKeyMqttClient   *mqttclient.Client  set by ConnectMqttPassword / ConnectMqttCert
//   - BagKeyAssetUUID    string              set by CreateAsset
//
// Writes (bag):
//   - BagKeyTelemetrySentAt  time.Time       publish timestamp
//
// Compensate: no-op. The message is in-flight; there is nothing local
// to undo.
func PublishTelemetry() saga.Step {
	return saga.Step{
		Name: "assets/assets.PublishTelemetry",
		Do: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyMqttClient)
			if !ok {
				return fmt.Errorf("publish telemetry: mqtt client missing on bag")
			}
			cli, ok := v.(*mqttclient.Client)
			if !ok {
				return fmt.Errorf("publish telemetry: bag[%s] is not *mqttclient.Client (%T)", BagKeyMqttClient, v)
			}

			uuid := c.MustGetString(BagKeyAssetUUID)
			topic := fmt.Sprintf("events/%s/temperature", uuid)

			payload := map[string]any{
				"runId":     c.RunID,
				"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
				"value":     23.5,
				"unit":      "celsius",
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("publish telemetry: marshal: %w", err)
			}

			if err := cli.Publish(c.Stdctx, topic, 1, false, body); err != nil {
				return fmt.Errorf("publish telemetry on %q: %w", topic, err)
			}
			c.Set(BagKeyTelemetrySentAt, time.Now().UTC())
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
