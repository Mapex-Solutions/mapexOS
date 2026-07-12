package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

// CreateTeamsTrigger POSTs the SagaTeamsTrigger payload to the
// triggers service and publishes the returned id on the bag.
//
// Reads (bag):
//   - BagKeyTriggerSinkHost / Port  the HTTP sink address (StartTestSink).
//
// Writes (bag):
//   - BagKeyTriggerID  string  Mongo ObjectID hex.
//
// Compensate: DELETE /api/v1/triggers/{id}.
func CreateTeamsTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateTeamsTrigger",
		Do: func(c *saga.Context) error {
			url, err := httpSinkURL(c)
			if err != nil {
				return fmt.Errorf("create teams trigger: %w", err)
			}
			spec := payloads.SagaTeamsTrigger(c.RunID, url)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create teams trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create teams trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create teams trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create teams trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: deleteTriggerOnCompensate("teams"),
	}
}
