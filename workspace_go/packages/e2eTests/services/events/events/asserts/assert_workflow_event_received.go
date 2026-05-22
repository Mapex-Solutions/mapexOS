package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	workflowInstanceSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/workflow/instances/steps"
)

// workflowCursorEnvelope mirrors the events service workflow listing
// response. Only the field the assert reads is decoded — everything
// else passes through verbatim.
type workflowCursorEnvelope struct {
	Data struct {
		Items []struct {
			ExecutionId string    `json:"executionId"`
			InstanceId  string    `json:"instanceId"`
			OrgId       string    `json:"orgId"`
			Status      string    `json:"status"`
			Success     bool      `json:"success"`
			Created     time.Time `json:"created"`
		} `json:"items"`
	} `json:"data"`
}

// AssertWorkflowEventReceivedAfter polls GET /api/v1/events/workflow
// filtered by the workflow instance the saga provisioned and a
// startTime read from the bag (the timestamp the action that should
// have fired the workflow ran at). Returns success on the first
// non-empty response or fails when the polling budget expires.
//
// startTimeBagKey lets callers point this assert at the right
// timestamp on the bag — assetSteps.BagKeyHeartbeatSentAt for online
// transitions and assetSteps.BagKeyForceOfflineSentAt for offline.
//
// Reads (bag):
//   - workflowInstanceSteps.BagKeyInstanceID  string     set by CreateInstance
//   - <startTimeBagKey>                       time.Time  set by the action under test
func AssertWorkflowEventReceivedAfter(startTimeBagKey string) saga.Assert {
	return AssertWorkflowEventReceivedAfterWithTimeout(startTimeBagKey, 60*time.Second, 1*time.Second)
}

// AssertWorkflowEventReceivedAfterWithTimeout overrides the polling
// budget. Workflow runtime cold starts can take a few seconds on the
// dev compose, so the default is intentionally higher than the raw
// events assert.
func AssertWorkflowEventReceivedAfterWithTimeout(startTimeBagKey string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("events/events.AssertWorkflowEventReceivedAfter[%s]", startTimeBagKey),
		Check: func(c *saga.Context) error {
			instanceID := c.MustGetString(workflowInstanceSteps.BagKeyInstanceID)

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
			query.Set("instanceId", instanceID)
			query.Set("startTime", startWithSlack)
			query.Set("limit", "20")

			deadline := time.Now().Add(timeout)
			lastCount := 0
			for {
				items, err := fetchWorkflowEvents(c, query.Encode())
				if err == nil {
					lastCount = len(items)
					if lastCount > 0 {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no workflow events for instanceId=%s after %s within %v (last poll returned %d)",
						instanceID, startWithSlack, timeout, lastCount)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("workflow event poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// fetchWorkflowEvents performs one GET pass and returns the items slice.
func fetchWorkflowEvents(c *saga.Context, query string) ([]struct {
	ExecutionId string    `json:"executionId"`
	InstanceId  string    `json:"instanceId"`
	OrgId       string    `json:"orgId"`
	Status      string    `json:"status"`
	Success     bool      `json:"success"`
	Created     time.Time `json:"created"`
}, error,
) {
	resp, err := c.Clients.Events.Raw(c.Stdctx, http.MethodGet, "/api/v1/events/workflow?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("get events/workflow: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get events/workflow: unexpected status %d", resp.StatusCode)
	}
	var out workflowCursorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode events/workflow response: %w", err)
	}
	return out.Data.Items, nil
}
