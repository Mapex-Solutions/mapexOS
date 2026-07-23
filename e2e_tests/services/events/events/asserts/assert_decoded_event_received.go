package asserts

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// AssertDecodedEventReceivedAfter proves two things at once: the events service
// CONSUMED the gateway's event, and the asset template's ScriptConversion RAN to
// produce the platform's StandardizedPayload. It polls the js-executor debug feed
// (GET /api/v1/events/jsexec) — the public surface carrying each script run's decoded
// output — for the asset's threadId, and succeeds once a SUCCESSFUL run after
// startTime carries a decoded `data` map matching want.
//
// Unlike AssertLorawanEventDecodedByLabel, it keys off the plain BagKeyAssetUUID that
// CreateAsset / CreateConnectivityAsset writes (not a label), so any HTTP/datasource
// telemetry flow can assert its template decode. Numbers in want compare with a small
// tolerance (JSON decodes them as float64); other values compare for equality.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID  string     threadId, set by CreateAsset
//   - <startTimeBagKey>           time.Time  set by the action that sent the event
func AssertDecodedEventReceivedAfter(startTimeBagKey string, want map[string]any) saga.Assert {
	return AssertDecodedEventReceivedAfterWithTimeout(startTimeBagKey, want, constants.ScaleTimeout(15*time.Second), 500*time.Millisecond)
}

// AssertDecodedEventReceivedAfterWithTimeout overrides the polling budget.
func AssertDecodedEventReceivedAfterWithTimeout(startTimeBagKey string, want map[string]any, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("events/events.AssertDecodedEventReceivedAfter[%s]", startTimeBagKey),
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

			startWithSlack := start.Add(-2 * time.Second).UTC().Format(time.RFC3339Nano)
			query := url.Values{}
			query.Set("threadId", uuid)
			query.Set("startTime", startWithSlack)
			query.Set("limit", "20")

			deadline := time.Now().Add(timeout)
			var lastMismatch, lastScriptErr string
			for {
				items, err := fetchJsExecEvents(c, query.Encode())
				if err == nil {
					for _, it := range items {
						if !it.Success {
							lastScriptErr = fmt.Sprintf("failedAt=%q error=%q", it.FailedAt, it.Error)
							continue
						}
						data, derr := decodedData(it.Event)
						if derr != nil {
							lastMismatch = derr.Error()
							continue
						}
						if miss := diffDecoded(data, want); miss == "" {
							return nil
						} else {
							lastMismatch = miss
						}
					}
				}
				if time.Now().After(deadline) {
					detail := lastMismatch
					if detail == "" && lastScriptErr != "" {
						detail = "script run(s) failed: " + lastScriptErr
					}
					if detail == "" {
						detail = "no jsexec event for the asset after the send"
					}
					return fmt.Errorf("events did not decode %v for threadId=%s within %v (%s)", want, uuid, timeout, detail)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("decoded event poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}
