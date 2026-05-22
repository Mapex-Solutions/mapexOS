package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// healthStatusResponse decodes the subset of /api/v1/assets/{id} the
// assert needs. The full DTO carries many fields; we ignore them.
type healthStatusResponse struct {
	Data struct {
		ID           string  `json:"id"`
		HealthStatus *string `json:"healthStatus,omitempty"`
	} `json:"data"`
}

// AssertHealthStatusEventually polls GET /api/v1/assets/{id} until the
// response body shows healthStatus == want, or the timeout elapses.
//
// This is the saga's externalised view of the presence pipeline: the
// MQTT client connects/disconnects, the assets service consumes the
// presence event from the broker, the cache flips the status field,
// and the next GET reflects the new state. By polling the public HTTP
// surface we avoid coupling the saga to NATS internals or to Redis
// keys.
//
// Defaults:
//   - timeout: 15s (presence flip is typically <1s; the headroom
//     accommodates a slow stack or first-run CI warm-up)
//   - tick:    500ms
//
// Reads (bag):
//   - assetSteps.BagKeyAssetID  string  set by CreateAsset
func AssertHealthStatusEventually(want string) saga.Assert {
	return AssertHealthStatusEventuallyWithTimeout(want, 15*time.Second, 500*time.Millisecond)
}

// AssertHealthStatusEventuallyWithTimeout overrides the polling
// budget. Use when a test deliberately needs a tighter window or
// debugging a slow stack.
func AssertHealthStatusEventuallyWithTimeout(want string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("assets/assets.AssertHealthStatus[%s]", want),
		Check: func(c *saga.Context) error {
			id := c.MustGetString(assetSteps.BagKeyAssetID)
			deadline := time.Now().Add(timeout)
			var lastSeen string

			for {
				status, err := fetchHealthStatus(c, id)
				if err == nil {
					lastSeen = status
					if status == want {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("asset %s healthStatus did not become %q within %v (last seen %q)", id, want, timeout, lastSeen)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("asset %s healthStatus poll cancelled: %w", id, c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// fetchHealthStatus is a single GET pass; returns the current status
// or an error. Empty string is reported as "" (the field is optional
// in the DTO when the asset has never received a presence signal).
func fetchHealthStatus(c *saga.Context, id string) (string, error) {
	resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, "/api/v1/assets/"+id, nil)
	if err != nil {
		return "", fmt.Errorf("get asset: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get asset: unexpected status %d", resp.StatusCode)
	}
	var out healthStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode get-asset response: %w", err)
	}
	if out.Data.HealthStatus == nil {
		return "", nil
	}
	return *out.Data.HealthStatus, nil
}
