package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

// CreateNatsTrigger POSTs the SagaNatsTrigger payload to the triggers
// service and publishes the returned id on the bag.
//
// Reads (bag):
//   - BagKeyNatsURL  string  set by StartNatsServer.
//
// Writes (bag):
//   - BagKeyTriggerID  string  Mongo ObjectID hex.
func CreateNatsTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateNatsTrigger",
		Do: func(c *saga.Context) error {
			natsURL := c.MustGetString(BagKeyNatsURL)

			spec := payloads.SagaNatsTrigger(c.RunID, natsURL)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create nats trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create nats trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create nats trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create nats trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: deleteTriggerOnCompensate("nats"),
	}
}
