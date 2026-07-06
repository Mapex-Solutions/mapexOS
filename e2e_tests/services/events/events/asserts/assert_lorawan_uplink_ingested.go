package asserts

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// AssertLorawanUplinkIngestedByLabel proves a LoRaWAN uplink from the labelled
// sensor reached MapexOS intact. It polls the events service's public raw feed
// (GET /api/v1/events/raw?threadId={assetUUID}) for an event that arrived after the
// uplink and whose payload equals the fired hex.
//
// What arrives is not just the hex: mapexLNS strips only the LoRaWAN/AES layer, so
// the raw event carries the undecoded application bytes plus the frame metadata
// (fPort/fCnt/rxInfo). The application-codec decode — the template's first script,
// run by js-executor and surfaced on route.execute — is a downstream concern and is
// NOT asserted here. The assert therefore checks: a raw event for the sensor's
// threadId after the uplink, whose event.bytes == bytes(payloadHex) and which
// carries fPort/fCnt/rxInfo.
//
// Reads (bag):
//   - assetSteps.AssetUUIDKey(label)           string     threadId, set by CreateAssetWithLabel
//   - assetSteps.LorawanUplinkSentAtKey(label)  time.Time  set by FireLorawanUplink
func AssertLorawanUplinkIngestedByLabel(label, payloadHex string) saga.Assert {
	return AssertLorawanUplinkIngestedByLabelWithTimeout(label, payloadHex, 20*time.Second, 500*time.Millisecond)
}

// AssertLorawanUplinkIngestedByLabelWithTimeout overrides the polling budget.
func AssertLorawanUplinkIngestedByLabelWithTimeout(label, payloadHex string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: "events/events.AssertLorawanUplinkIngested[" + label + "]",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.AssetUUIDKey(label))
			want, err := hex.DecodeString(payloadHex)
			if err != nil {
				return fmt.Errorf("bad payloadHex %q: %w", payloadHex, err)
			}

			startVal, ok := c.Get(assetSteps.LorawanUplinkSentAtKey(label))
			if !ok {
				return fmt.Errorf("bag key %q missing; fire the uplink first", assetSteps.LorawanUplinkSentAtKey(label))
			}
			start, ok := startVal.(time.Time)
			if !ok {
				return fmt.Errorf("uplink sent-at is not time.Time (%T)", startVal)
			}

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
					for _, it := range items {
						if lorawanBytesMatch(it.Event["bytes"], want) &&
							fieldPresent(it.Event, "fPort") &&
							fieldPresent(it.Event, "fCnt") &&
							fieldPresent(it.Event, "rxInfo") {
							return nil
						}
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("no ingested lorawan uplink matching %s for threadId=%s after %s within %v (last poll returned %d events)",
						payloadHex, uuid, startWithSlack, timeout, lastCount)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("lorawan ingest poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// lorawanBytesMatch reports whether a JSON number array (the raw event's
// event.bytes) equals the wanted byte slice.
func lorawanBytesMatch(v any, want []byte) bool {
	arr, ok := v.([]any)
	if !ok || len(arr) != len(want) {
		return false
	}
	for i, e := range arr {
		f, ok := e.(float64)
		if !ok || byte(int(f)) != want[i] {
			return false
		}
	}
	return true
}

// fieldPresent reports whether the map carries a non-nil value at key k.
func fieldPresent(m map[string]any, k string) bool {
	v, ok := m[k]
	return ok && v != nil
}
