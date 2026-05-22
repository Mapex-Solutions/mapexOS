package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// triggerCursorEnvelope mirrors the events service trigger listing
// response. Only the fields the assert reads are decoded.
type triggerCursorEnvelope struct {
	Data struct {
		Items []struct {
			TriggerId   string    `json:"triggerId"`
			TriggerName string    `json:"triggerName"`
			TriggerType string    `json:"triggerType"`
			Success     bool      `json:"success"`
			Created     time.Time `json:"created"`
		} `json:"items"`
	} `json:"data"`
}

// AssertTriggerEventReceivedAfter polls GET /api/v1/events/trigger
// filtered by the trigger the saga created and a startTime read from
// the bag (the timestamp the action that should have fired the
// trigger ran at). Returns success on the first non-empty response or
// fails when the polling budget expires.
//
// Reads (bag):
//   - triggerSteps.BagKeyTriggerID  string     set by CreateTrigger
//   - <startTimeBagKey>             time.Time  set by the action under test
func AssertTriggerEventReceivedAfter(startTimeBagKey string) saga.Assert {
	return AssertTriggerEventReceivedAfterWithTimeout(startTimeBagKey, 60*time.Second, 1*time.Second)
}

// AssertTriggerEventReceivedAfterWithTimeout overrides the polling
// budget.
func AssertTriggerEventReceivedAfterWithTimeout(startTimeBagKey string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("events/events.AssertTriggerEventReceivedAfter[%s]", startTimeBagKey),
		Check: func(c *saga.Context) error {
			triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)

			startVal, ok := c.Get(startTimeBagKey)
			if !ok {
				return fmt.Errorf("bag key %q missing; cannot bound the search window", startTimeBagKey)
			}
			start, ok := startVal.(time.Time)
			if !ok {
				return fmt.Errorf("bag key %q is not time.Time (%T)", startTimeBagKey, startVal)
			}

			startWithSlack := start.Add(-2 * time.Second).UTC().Format(time.RFC3339Nano)
			query := url.Values{}
			query.Set("triggerId", triggerID)
			query.Set("startTime", startWithSlack)
			query.Set("limit", "20")

			deadline := time.Now().Add(timeout)
			lastCount := 0
			for {
				items, err := fetchTriggerEvents(c, query.Encode())
				if err == nil {
					lastCount = len(items)
					if lastCount > 0 {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no trigger events for triggerId=%s after %s within %v (last poll returned %d)",
						triggerID, startWithSlack, timeout, lastCount)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("trigger event poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// fetchTriggerEvents performs one GET pass and returns the items slice.
func fetchTriggerEvents(c *saga.Context, query string) ([]struct {
	TriggerId   string    `json:"triggerId"`
	TriggerName string    `json:"triggerName"`
	TriggerType string    `json:"triggerType"`
	Success     bool      `json:"success"`
	Created     time.Time `json:"created"`
}, error,
) {
	resp, err := c.Clients.Events.Raw(c.Stdctx, http.MethodGet, "/api/v1/events/trigger?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("get events/trigger: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get events/trigger: unexpected status %d", resp.StatusCode)
	}
	var out triggerCursorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode events/trigger response: %w", err)
	}
	return out.Data.Items, nil
}
