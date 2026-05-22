package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

// CreateRabbitmqTrigger POSTs the SagaRabbitmqTrigger payload to the
// triggers service and publishes the returned id on the bag.
//
// Reads (bag):
//   - BagKeyRabbitmqHost / Port / User / Pass  set by StartRabbitmqContainer.
//
// Writes (bag):
//   - BagKeyTriggerID  string  Mongo ObjectID hex.
func CreateRabbitmqTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateRabbitmqTrigger",
		Do: func(c *saga.Context) error {
			host := c.MustGetString(BagKeyRabbitmqHost)
			portVal, ok := c.Get(BagKeyRabbitmqPort)
			if !ok {
				return fmt.Errorf("create rabbitmq trigger: bag key %q missing", BagKeyRabbitmqPort)
			}
			port, ok := portVal.(int)
			if !ok {
				return fmt.Errorf("create rabbitmq trigger: bag key %q is not int (%T)", BagKeyRabbitmqPort, portVal)
			}
			user := c.MustGetString(BagKeyRabbitmqUser)
			pass := c.MustGetString(BagKeyRabbitmqPass)

			spec := payloads.SagaRabbitmqTrigger(c.RunID, host, port, user, pass)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create rabbitmq trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create rabbitmq trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create rabbitmq trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create rabbitmq trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: deleteTriggerOnCompensate("rabbitmq"),
	}
}
