package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/payloads"
)

// CreateWebsocketTrigger POSTs the SagaWebsocketTrigger payload to
// the triggers service and publishes the returned id on the bag.
//
// Reads (bag):
//   - BagKeyWsHost / Port  the WS sink address (StartWebsocketSink).
func CreateWebsocketTrigger() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.CreateWebsocketTrigger",
		Do: func(c *saga.Context) error {
			url, err := wsSinkURL(c)
			if err != nil {
				return fmt.Errorf("create websocket trigger: %w", err)
			}
			spec := payloads.SagaWebsocketTrigger(c.RunID, url)
			resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodPost, "/api/v1/triggers", spec)
			if err != nil {
				return fmt.Errorf("create websocket trigger: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create websocket trigger: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out triggerCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create websocket trigger response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create websocket trigger: empty id in response")
			}
			c.Set(BagKeyTriggerID, out.Data.ID)
			return nil
		},
		Compensate: deleteTriggerOnCompensate("websocket"),
	}
}
