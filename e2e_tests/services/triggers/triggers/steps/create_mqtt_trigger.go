package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

// CreateMqttTrigger POSTs the SagaMqttTrigger payload to the triggers
// service and publishes the returned id on the bag.
//
// Reads (bag):
//   - BagKeyMqttBrokerHost  string  set by StartMqttBroker.
//   - BagKeyMqttBrokerPort  int     set by StartMqttBroker.
//
// Writes (bag):
//   - BagKeyTriggerID  string  Mongo ObjectID hex.
func CreateMqttTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateMqttTrigger",
		Do: func(c *saga.Context) error {
			host := c.MustGetString(BagKeyMqttBrokerHost)
			portVal, ok := c.Get(BagKeyMqttBrokerPort)
			if !ok {
				return fmt.Errorf("create mqtt trigger: bag key %q missing", BagKeyMqttBrokerPort)
			}
			port, ok := portVal.(int)
			if !ok {
				return fmt.Errorf("create mqtt trigger: bag key %q is not int (%T)", BagKeyMqttBrokerPort, portVal)
			}

			spec := payloads.SagaMqttTrigger(c.RunID, host, port)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create mqtt trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create mqtt trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create mqtt trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create mqtt trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: deleteTriggerOnCompensate("mqtt"),
	}
}
