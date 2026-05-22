package steps

import (
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// deleteTriggerOnCompensate returns the Compensate body every
// CreateXTrigger step shares: read BagKeyTriggerID, DELETE the
// trigger, tolerate 404. `label` is the trigger type used only in
// error messages so a failing rollback identifies the type that broke.
func deleteTriggerOnCompensate(label string) func(c *saga.Context) error {
	return func(c *saga.Context) error {
		id, ok := c.Get(BagKeyTriggerID)
		if !ok {
			return nil
		}
		resp, err := c.Clients.Triggers.Raw(c.Stdctx, http.MethodDelete, "/api/v1/triggers/"+id.(string), nil)
		if err != nil {
			return fmt.Errorf("delete %s trigger: %w", label, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("delete %s trigger: unexpected status %d", label, resp.StatusCode)
		}
		return nil
	}
}
