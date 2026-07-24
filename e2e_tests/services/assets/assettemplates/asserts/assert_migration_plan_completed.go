// Package asserts holds saga oracles for the assettemplates module. They validate
// the platform through its public HTTP API only — never NATS/Mongo/ClickHouse.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	steps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
)

// migrationPlanEnvelope decodes the plan read model out of the {status,errors,data}
// envelope.
type migrationPlanEnvelope struct {
	Data struct {
		Status   string `json:"status"`
		Total    int    `json:"total"`
		Migrated int    `json:"migrated"`
		Failed   int    `json:"failed"`
	} `json:"data"`
}

// AssertMigrationPlanStatusEventually polls GET /migrations/:id until the plan reaches
// the wanted status; for a terminal status it also enforces the counter invariant
// ("complete" ⇒ migrated==total && failed==0; "completed_with_errors" ⇒ failed>=1 &&
// migrated==total-failed). It asserts only terminal states, never the intermediate
// "running", so the immediate-but-async run is not raced.
//
// Reads (bag):
//   - steps.BagKeyMigrationPlanID  string  set by CreateMigrationPlan
func AssertMigrationPlanStatusEventually(want string) saga.Assert {
	return AssertMigrationPlanStatusEventuallyWithTimeout(want, constants.ScaleTimeout(60*time.Second), time.Second)
}

// AssertMigrationPlanStatusEventuallyWithTimeout overrides the polling budget.
func AssertMigrationPlanStatusEventuallyWithTimeout(want string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("assets/assettemplates.AssertMigrationPlanStatus[%s]", want),
		Check: func(c *saga.Context) error {
			planID := c.MustGetString(steps.BagKeyMigrationPlanID)
			path := "/api/v1/asset_templates/migrations/" + planID

			deadline := time.Now().Add(timeout)
			lastStatus := ""
			lastTotal, lastMigrated, lastFailed := 0, 0, 0
			var lastErr error
			for {
				data, err := fetchMigrationPlan(c, path)
				if err != nil {
					lastErr = err
				} else {
					lastStatus = data.Status
					lastTotal, lastMigrated, lastFailed = data.Total, data.Migrated, data.Failed
					if data.Status == want && counterInvariantHolds(want, data.Total, data.Migrated, data.Failed) {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("migration plan %s did not reach status=%q within %v (last status=%q total=%d migrated=%d failed=%d, last error: %v)",
						planID, want, timeout, lastStatus, lastTotal, lastMigrated, lastFailed, lastErr)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("migration plan poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// counterInvariantHolds enforces the migrated/failed math for the terminal statuses;
// non-terminal wanted statuses (should not be used) impose no counter constraint.
func counterInvariantHolds(want string, total, migrated, failed int) bool {
	switch want {
	case "complete":
		return migrated == total && failed == 0
	case "completed_with_errors":
		return failed >= 1 && migrated == total-failed
	default:
		return true
	}
}

// fetchMigrationPlan performs one GET pass and returns the decoded plan data.
func fetchMigrationPlan(c *saga.Context, path string) (struct {
	Status   string `json:"status"`
	Total    int    `json:"total"`
	Migrated int    `json:"migrated"`
	Failed   int    `json:"failed"`
}, error,
) {
	var out migrationPlanEnvelope
	resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, path, nil)
	if err != nil {
		return out.Data, fmt.Errorf("get migration plan: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out.Data, fmt.Errorf("get migration plan: unexpected status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return out.Data, fmt.Errorf("decode migration plan response: %w", err)
	}
	return out.Data, nil
}
