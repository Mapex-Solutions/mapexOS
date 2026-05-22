// Package asserts holds saga oracles for the events/events module.
//
// These asserts query the events service through its public HTTP API
// (GET /api/v1/events/raw and friends) — they NEVER subscribe to NATS
// subjects, peek at Mongo collections, or query ClickHouse directly.
// The point of a saga is to validate the platform from the
// outside-the-stack perspective; observing infra would couple the test
// to internals that should be free to change.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// rawCursorEnvelope mirrors the events service response shape:
// StandardResponse wrapping a CursorResult.
type rawCursorEnvelope struct {
	Data struct {
		Items []struct {
			Created  time.Time      `json:"created"`
			ThreadID string         `json:"threadId"`
			OrgID    string         `json:"orgId"`
			Source   string         `json:"source"`
			Event    map[string]any `json:"event"`
			Success  bool           `json:"success"`
		} `json:"items"`
	} `json:"data"`
}

// AssertRawEventReceivedAfter polls GET /api/v1/events/raw filtered by
// the saga asset's UUID (mapped to threadId at gateway ingest time)
// and a startTime read from the bag, until the response carries at
// least one item or the timeout elapses.
//
// startTimeBagKey lets callers point this assert at the right
// timestamp on the bag — Phase 2 uses BagKeyMqttConnectedAt; Phase 3
// uses BagKeyTelemetrySentAt — so the search window is "events that
// happened after the action under test", never the entire history.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID  string     set by CreateAsset
//   - <startTimeBagKey>           time.Time  set by the action that triggered the event
func AssertRawEventReceivedAfter(startTimeBagKey string) saga.Assert {
	return AssertRawEventReceivedAfterWithTimeout(startTimeBagKey, 15*time.Second, 500*time.Millisecond)
}

// AssertRawEventReceivedAfterWithTimeout overrides the polling budget.
func AssertRawEventReceivedAfterWithTimeout(startTimeBagKey string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("events/events.AssertRawEventReceivedAfter[%s]", startTimeBagKey),
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)

			startVal, ok := c.Get(startTimeBagKey)
			if !ok {
				return fmt.Errorf("bag key %q missing; cannot bound the search window", startTimeBagKey)
			}
			start, ok := startVal.(time.Time)
			if !ok {
				return fmt.Errorf("bag key %q is not time.Time (%T)", startTimeBagKey, startVal)
			}

			// Subtract a small slack so the boundary is inclusive across
			// clock skew between the test runner and the events service.
			startWithSlack := start.Add(-2 * time.Second).UTC().Format(time.RFC3339Nano)
			query := url.Values{}
			query.Set("threadId", uuid)
			query.Set("startTime", startWithSlack)
			query.Set("limit", "20")

			deadline := time.Now().Add(timeout)
			lastCount := 0
			for {
				items, err := fetchRawEvents(c, query.Encode())
				if err == nil {
					lastCount = len(items)
					if lastCount > 0 {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no raw events for threadId=%s after %s within %v (last poll returned %d)",
						uuid, startWithSlack, timeout, lastCount)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("raw event poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// fetchRawEvents performs one GET pass and returns the items slice.
func fetchRawEvents(c *saga.Context, query string) ([]struct {
	Created  time.Time      `json:"created"`
	ThreadID string         `json:"threadId"`
	OrgID    string         `json:"orgId"`
	Source   string         `json:"source"`
	Event    map[string]any `json:"event"`
	Success  bool           `json:"success"`
}, error,
) {
	resp, err := c.Clients.Events.Raw(c.Stdctx, http.MethodGet, "/api/v1/events/raw?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("get events/raw: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get events/raw: unexpected status %d", resp.StatusCode)
	}
	var out rawCursorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode events/raw response: %w", err)
	}
	return out.Data.Items, nil
}
