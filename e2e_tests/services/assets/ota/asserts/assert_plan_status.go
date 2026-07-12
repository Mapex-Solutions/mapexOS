// Package asserts holds read-only OTA oracles: they poll the public HTTP API
// (never Mongo/NATS/ClickHouse) until the plan/execution/asset reach the expected
// state, or a timeout. Both OTA journeys import them.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	otaSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/steps"
)

type planStatusResponse struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}

// AssertPlanStatusEventually polls GET /api/v1/ota/plans/{planId} until
// data.status == want or the timeout elapses. Default budget 60s / 1s tick — a
// plan reaches a terminal status only after the reconciler dispatches (one scan)
// and the device reports the full progression.
//
// Reads (bag):
//   - otaSteps.BagKeyPlanID  string  set by CreatePlan
func AssertPlanStatusEventually(want string) saga.Assert {
	return AssertPlanStatusEventuallyWithTimeout(want, 60*time.Second, time.Second)
}

// AssertPlanStatusEventuallyWithTimeout overrides the polling budget.
func AssertPlanStatusEventuallyWithTimeout(want string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("assets/ota.AssertPlanStatus[%s]", want),
		Check: func(c *saga.Context) error {
			id := c.MustGetString(otaSteps.BagKeyPlanID)
			deadline := time.Now().Add(timeout)
			var lastSeen string
			for {
				status, err := fetchPlanStatus(c, id)
				if err == nil {
					lastSeen = status
					if status == want {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("ota plan %s status did not become %q within %v (last seen %q)", id, want, timeout, lastSeen)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("ota plan %s status poll cancelled: %w", id, c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

func fetchPlanStatus(c *saga.Context, id string) (string, error) {
	resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, "/api/v1/ota/plans/"+id, nil)
	if err != nil {
		return "", fmt.Errorf("get ota plan: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get ota plan: unexpected status %d", resp.StatusCode)
	}
	var out planStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode get-ota-plan response: %w", err)
	}
	return out.Data.Status, nil
}
