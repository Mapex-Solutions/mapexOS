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

// migrationExecutionsEnvelope decodes the per-asset executions out of the paginated
// {status,errors,data:{items}} envelope.
type migrationExecutionsEnvelope struct {
	Data struct {
		Items []struct {
			AssetID  string `json:"assetId"`
			Status   string `json:"status"`
			Error    string `json:"error"`
			Attempts int    `json:"attempts"`
		} `json:"items"`
	} `json:"data"`
}

// AssertMigrationExecutionStatus polls GET /migrations/:id/executions until the row
// for the asset (read from assetIDKey) reaches the wanted status. When want=="failed"
// it also requires a non-empty error and attempts==1, so the failure path is really
// exercised, not just an unattempted row.
//
// Reads (bag):
//   - steps.BagKeyMigrationPlanID  string  set by CreateMigrationPlan
//   - <assetIDKey>                 string  the asset id to match
func AssertMigrationExecutionStatus(assetIDKey, want string) saga.Assert {
	return assertMigrationExecution(func(c *saga.Context) string { return c.MustGetString(assetIDKey) }, want, constants.ScaleTimeout(60*time.Second), time.Second)
}

// AssertMigrationExecutionStatusByID is AssertMigrationExecutionStatus keyed by a
// literal asset id (compile-time wiring) rather than a bag key — used for a
// fabricated, non-existent asset id in the partial-failure phase.
func AssertMigrationExecutionStatusByID(assetID, want string) saga.Assert {
	return assertMigrationExecution(func(*saga.Context) string { return assetID }, want, constants.ScaleTimeout(60*time.Second), time.Second)
}

// assertMigrationExecution is the shared implementation; resolveAssetID defers reading
// the asset id (bag or literal) to Check time.
func assertMigrationExecution(resolveAssetID func(*saga.Context) string, want string, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("assets/assettemplates.AssertMigrationExecutionStatus[%s]", want),
		Check: func(c *saga.Context) error {
			planID := c.MustGetString(steps.BagKeyMigrationPlanID)
			assetID := resolveAssetID(c)
			path := "/api/v1/asset_templates/migrations/" + planID + "/executions"

			deadline := time.Now().Add(timeout)
			lastSeen := "<no row>"
			var lastErr error
			for {
				items, err := fetchMigrationExecutions(c, path)
				if err != nil {
					lastErr = err
				} else {
					for _, it := range items {
						if it.AssetID != assetID {
							continue
						}
						lastSeen = fmt.Sprintf("status=%q error=%q attempts=%d", it.Status, it.Error, it.Attempts)
						if it.Status == want && failureDetailOK(want, it.Error, it.Attempts) {
							return nil
						}
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("migration execution for asset %s did not reach status=%q within %v (last seen %s, last error: %v)",
						assetID, want, timeout, lastSeen, lastErr)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("migration executions poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// failureDetailOK requires a real failure signature (non-empty error, one attempt)
// only when the wanted status is "failed"; other statuses impose no such constraint.
func failureDetailOK(want, execErr string, attempts int) bool {
	if want != "failed" {
		return true
	}
	return execErr != "" && attempts == 1
}

// fetchMigrationExecutions performs one GET pass and returns the execution rows.
func fetchMigrationExecutions(c *saga.Context, path string) ([]struct {
	AssetID  string `json:"assetId"`
	Status   string `json:"status"`
	Error    string `json:"error"`
	Attempts int    `json:"attempts"`
}, error,
) {
	var out migrationExecutionsEnvelope
	resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get migration executions: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get migration executions: unexpected status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode migration executions response: %w", err)
	}
	return out.Data.Items, nil
}
