package asserts

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// AssertLorawanEventDecodedByLabel proves the device CODEC ran end to end: a
// LoRaWAN uplink from the labelled sensor was decoded by the template's first
// script (ScriptConversion) into the device's semantic fields — not left as raw
// bytes. It polls the js-executor debug feed (GET /api/v1/events/jsexec), the
// public surface that carries each script run's StandardizedPayload output, for a
// SUCCESSFUL run after the uplink whose decoded data matches want.
//
// This is the counterpart of AssertLorawanUplinkIngested (which checks the RAW
// bytes): here we assert the DECODED output, so a regression in the codec pipeline
// (js-executor applying the template script) is caught.
//
// want holds the expected decoded fields under the StandardizedPayload's data,
// e.g. {"TempC_SHT": 21.0, "Hum_SHT": 60.0, "Node_type": "LHT65N"}. Numbers are
// compared with a small tolerance (JSON decodes them as float64); other values are
// compared for equality.
//
// Reads (bag):
//   - assetSteps.AssetUUIDKey(label)            string     threadId, set by CreateAssetWithLabel
//   - assetSteps.LorawanUplinkSentAtKey(label)  time.Time  set by FireLorawanUplink
func AssertLorawanEventDecodedByLabel(label string, want map[string]any) saga.Assert {
	return AssertLorawanEventDecodedByLabelWithTimeout(label, want, 30*time.Second, 500*time.Millisecond)
}

// AssertLorawanEventDecodedByLabelWithTimeout overrides the polling budget.
func AssertLorawanEventDecodedByLabelWithTimeout(label string, want map[string]any, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: "events/events.AssertLorawanEventDecoded[" + label + "]",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.AssetUUIDKey(label))

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
			var lastMismatch string
			var lastScriptErr string
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
						detail = "no jsexec event for the sensor after the uplink"
					}
					return fmt.Errorf("lorawan codec did not decode %v for threadId=%s within %v (%s)", want, uuid, timeout, detail)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("lorawan decode poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// jsExecItem is the subset of a /events/jsexec item this assert reads. Event is
// the StandardizedPayload serialized as a JSON string (the codec's output).
type jsExecItem struct {
	Created  time.Time `json:"created"`
	ThreadID string    `json:"threadId"`
	Event    string    `json:"event"`
	Success  bool      `json:"success"`
	FailedAt string    `json:"failedAt"`
	Error    string    `json:"error"`
}

type jsExecCursorEnvelope struct {
	Data struct {
		Items []jsExecItem `json:"items"`
	} `json:"data"`
}

// fetchJsExecEvents pulls one page of js-executor events for the given query.
func fetchJsExecEvents(c *saga.Context, query string) ([]jsExecItem, error) {
	resp, err := c.Clients.Events.Raw(c.Stdctx, http.MethodGet, "/api/v1/events/jsexec?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("get events/jsexec: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get events/jsexec: unexpected status %d", resp.StatusCode)
	}
	var out jsExecCursorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode events/jsexec response: %w", err)
	}
	return out.Data.Items, nil
}

// decodedData parses the StandardizedPayload JSON string and returns its data map.
func decodedData(eventJSON string) (map[string]any, error) {
	if eventJSON == "" {
		return nil, fmt.Errorf("empty jsexec event payload")
	}
	var std struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal([]byte(eventJSON), &std); err != nil {
		return nil, fmt.Errorf("parse standardized payload: %w", err)
	}
	if std.Data == nil {
		return nil, fmt.Errorf("standardized payload has no data object")
	}
	return std.Data, nil
}

// diffDecoded returns "" when every want field is present in data with a matching
// value, or a human-readable description of the first mismatch. Numbers compare
// with a tolerance (JSON numbers decode as float64); other values by equality.
func diffDecoded(data, want map[string]any) string {
	for k, w := range want {
		got, ok := data[k]
		if !ok {
			return fmt.Sprintf("field %q missing (got data=%v)", k, data)
		}
		if wf, wok := toFloat(w); wok {
			gf, gok := toFloat(got)
			if !gok || math.Abs(gf-wf) > 1e-6 {
				return fmt.Sprintf("field %q = %v, want %v", k, got, w)
			}
			continue
		}
		if fmt.Sprint(got) != fmt.Sprint(w) {
			return fmt.Sprintf("field %q = %v, want %v", k, got, w)
		}
	}
	return ""
}

// toFloat coerces JSON-decoded numeric values to float64.
func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}
