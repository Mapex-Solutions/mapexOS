// Package steps holds saga steps that exercise the triggers module
// HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

type triggerCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CreateTrigger POSTs the canonical SagaSimpleTrigger payload to the
// triggers service and publishes the returned id on the bag.
//
// Writes (bag):
//   - BagKeyTriggerID  string  Mongo ObjectID hex
//
// Compensate: DELETE /api/v1/triggers/{id}. Idempotent.
func CreateTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateTrigger",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaSimpleTrigger(c.RunID)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyTriggerID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodDelete, "/api/v1/triggers/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete trigger: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
