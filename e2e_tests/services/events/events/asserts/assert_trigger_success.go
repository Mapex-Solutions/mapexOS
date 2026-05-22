package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// triggerExecRecord captures the fields the success / content-key
// asserts read out of /api/v1/events/trigger. Defined separately from
// triggerCursorEnvelope so it can carry the requestData and error
// fields the content-key oracle inspects.
type triggerExecRecord struct {
	TriggerId   string    `json:"triggerId"`
	TriggerName string    `json:"triggerName"`
	TriggerType string    `json:"triggerType"`
	Success     bool      `json:"success"`
	Error       string    `json:"error"`
	RequestData string    `json:"requestData"`
	Created     time.Time `json:"created"`
}

type triggerExecEnvelope struct {
	Data struct {
		Items []triggerExecRecord `json:"items"`
	} `json:"data"`
}

// AssertTriggerExecutedSuccessfullyEventually polls
// /api/v1/events/trigger filtered by the saga's triggerId and waits
// until at least `expected` records carry success=true. Used by the
// MQTT / NATS / RabbitMQ / WebSocket smokes — these protocols would
// otherwise require a per-protocol subscriber to verify delivery; the
// events service publishes the success flag only after the trigger's
// publish call returns nil, so it is a valid downstream oracle.
//
// Reads (bag):
//   - triggerSteps.BagKeyTriggerID  string  set by CreateXTrigger
//
// Default polling budget mirrors the existing
// AssertTriggerEventReceivedAfter (60 s / 1 s tick). The ClickHouse
// insert pipeline carries a batch buffer so the first observation
// can take 10–30 s in DEV.
func AssertTriggerExecutedSuccessfullyEventually(expected int) saga.Assert {
	return AssertTriggerExecutedSuccessfullyEventuallyWithTimeout(expected, 60*time.Second, 1*time.Second)
}

// AssertTriggerExecutedSuccessfullyEventuallyWithTimeout overrides
// the polling budget.
func AssertTriggerExecutedSuccessfullyEventuallyWithTimeout(expected int, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("events/events.AssertTriggerExecutedSuccessfullyEventually[%d]", expected),
		Check: func(c *saga.Context) error {
			triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
			query := url.Values{}
			query.Set("triggerId", triggerID)
			query.Set("limit", "50")

			deadline := time.Now().Add(timeout)
			var lastErr string
			for {
				items, err := fetchTriggerExecRecords(c, query.Encode())
				if err == nil {
					ok := 0
					lastErr = ""
					for _, it := range items {
						if it.Success {
							ok++
						} else if it.Error != "" {
							lastErr = it.Error
						}
					}
					if ok >= expected {
						return nil
					}
				}
				if time.Now().After(deadline) {
					hint := ""
					if lastErr != "" {
						hint = fmt.Sprintf(" — last failure: %s", lastErr)
					}
					return fmt.Errorf("trigger %s: want >=%d successful executions in %v%s",
						triggerID, expected, timeout, hint)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("trigger success poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// AssertLastTriggerRequestDataContains validates that the most recent
// successful execution record's requestData (serialised resolved
// trigger config) contains every expected substring. Content-key
// oracle for the same protocols that use
// AssertTriggerExecutedSuccessfullyEventually: it lets the smoke
// assert message contents (topic, subject, body, etc.) without
// having to subscribe to the broker directly.
//
// Substring match is intentional — the resolved config is a JSON
// blob, an exact compare would be fragile across template tweaks.
//
// Reads (bag):
//   - triggerSteps.BagKeyTriggerID  string  set by CreateXTrigger
func AssertLastTriggerRequestDataContains(expectedSubstrings ...string) saga.Assert {
	return saga.Assert{
		Name: "events/events.AssertLastTriggerRequestDataContains",
		Check: func(c *saga.Context) error {
			triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
			query := url.Values{}
			query.Set("triggerId", triggerID)
			query.Set("limit", "20")

			items, err := fetchTriggerExecRecords(c, query.Encode())
			if err != nil {
				return fmt.Errorf("fetch trigger records: %w", err)
			}
			var latest *triggerExecRecord
			for i := range items {
				if items[i].Success {
					if latest == nil || items[i].Created.After(latest.Created) {
						latest = &items[i]
					}
				}
			}
			if latest == nil {
				return fmt.Errorf("trigger %s: no successful executions to inspect requestData", triggerID)
			}
			for _, want := range expectedSubstrings {
				if !strings.Contains(latest.RequestData, want) {
					return fmt.Errorf("trigger %s requestData: want substring %q in %q",
						triggerID, want, latest.RequestData)
				}
			}
			return nil
		},
	}
}

// fetchTriggerExecRecords performs one GET pass and returns the items
// slice. Defined here (not reused from assert_trigger_event_received.go)
// because the existing helper's envelope is narrower and would lose
// the requestData / error fields these asserts need.
func fetchTriggerExecRecords(c *saga.Context, query string) ([]triggerExecRecord, error) {
	resp, err := c.Clients.Events.Raw(c.Stdctx, http.MethodGet, "/api/v1/events/trigger?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("get events/trigger: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get events/trigger: unexpected status %d", resp.StatusCode)
	}
	var out triggerExecEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode events/trigger response: %w", err)
	}
	return out.Data.Items, nil
}
