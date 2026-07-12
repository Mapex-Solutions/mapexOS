package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	otaSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/steps"
)

type executionsResponse struct {
	Data struct {
		Items []struct {
			AssetID    string `json:"assetId"`
			State      string `json:"state"`
			Percentage int32  `json:"percentage"`
		} `json:"items"`
	} `json:"data"`
}

// AssertExecutionStateEventually polls GET /api/v1/ota/plans/{planId}/executions
// until the execution for the asset (assetIDKey) reaches wantState with
// percentage >= wantPct, or the timeout elapses. Default budget 60s / 1s. The
// per-device execution advances as the device reports its progression.
//
// Reads (bag):
//   - otaSteps.BagKeyPlanID  string  set by CreatePlan
//   - assetIDKey             string  the asset id (set by CreateAsset*)
func AssertExecutionStateEventually(assetIDKey, wantState string, wantPct int) saga.Assert {
	return AssertExecutionStateEventuallyWithTimeout(assetIDKey, wantState, wantPct, 60*time.Second, time.Second)
}

// AssertExecutionStateEventuallyWithTimeout overrides the polling budget.
func AssertExecutionStateEventuallyWithTimeout(assetIDKey, wantState string, wantPct int, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("assets/ota.AssertExecutionState[%s]", wantState),
		Check: func(c *saga.Context) error {
			planID := c.MustGetString(otaSteps.BagKeyPlanID)
			assetID := c.MustGetString(assetIDKey)
			deadline := time.Now().Add(timeout)
			var lastState string
			var lastPct int
			for {
				st, pct, found, err := fetchExecutionState(c, planID, assetID)
				if err == nil && found {
					lastState, lastPct = st, pct
					if st == wantState && pct >= wantPct {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("ota execution for asset %s in plan %s did not reach %q@%d within %v (last seen %q@%d)",
						assetID, planID, wantState, wantPct, timeout, lastState, lastPct)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("ota execution poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

func fetchExecutionState(c *saga.Context, planID, assetID string) (state string, pct int, found bool, err error) {
	resp, e := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, "/api/v1/ota/plans/"+planID+"/executions", nil)
	if e != nil {
		return "", 0, false, fmt.Errorf("get ota executions: %w", e)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", 0, false, fmt.Errorf("get ota executions: unexpected status %d", resp.StatusCode)
	}
	var out executionsResponse
	if e := json.NewDecoder(resp.Body).Decode(&out); e != nil {
		return "", 0, false, fmt.Errorf("decode get-ota-executions response: %w", e)
	}
	for _, it := range out.Data.Items {
		if it.AssetID == assetID {
			return it.State, int(it.Percentage), true, nil
		}
	}
	return "", 0, false, nil
}
