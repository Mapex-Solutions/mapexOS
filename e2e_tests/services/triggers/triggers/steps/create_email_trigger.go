package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

// CreateEmailTrigger POSTs the SagaEmailTrigger payload to the
// triggers service and publishes the returned id on the bag.
//
// Reads (bag):
//   - BagKeySmtpHost / Port  the SMTP sink address (StartSmtpSink).
//
// Writes (bag):
//   - BagKeyTriggerID  string  Mongo ObjectID hex.
//
// Compensate: DELETE /api/v1/triggers/{id}. Idempotent (tolerates 404).
func CreateEmailTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateEmailTrigger",
		Do: func(c *saga.Context) error {
			host, port, err := smtpSinkAddr(c)
			if err != nil {
				return fmt.Errorf("create email trigger: %w", err)
			}
			spec := payloads.SagaEmailTrigger(c.RunID, host, port)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create email trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create email trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create email trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create email trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: deleteTriggerOnCompensate("email"),
	}
}
