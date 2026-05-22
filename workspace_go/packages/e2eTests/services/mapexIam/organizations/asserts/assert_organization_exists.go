// Package asserts holds saga oracles for the IAM organizations module.
// Asserts never mutate the system; they read the live stack and fail the
// test when the expected state is not observed.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	orgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/steps"
)

type orgGetResponse struct {
	Data struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Enabled bool   `json:"enabled"`
		PathKey string `json:"pathKey"`
	} `json:"data"`
}

// AssertOrganizationExists fetches the org by id from the bag and verifies
// the API returns a populated entity with enabled=true.
//
// Reads (bag):
//   - orgSteps.BagKeyOrgID  string  set by CreateOrganization or SeedAdminLogin
func AssertOrganizationExists() saga.Assert {
	return saga.Assert{
		Name: "mapexIam/organizations.AssertOrganizationExists",
		Check: func(c *saga.Context) error {
			id := c.MustGetString(orgSteps.BagKeyOrgID)
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodGet, "/api/v1/organizations/"+id, nil)
			if err != nil {
				return fmt.Errorf("get organization: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("get organization: unexpected status %d", resp.StatusCode)
			}
			var out orgGetResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode get-org response: %w", err)
			}
			if out.Data.ID != id {
				return fmt.Errorf("get organization: id mismatch want %q got %q", id, out.Data.ID)
			}
			if !out.Data.Enabled {
				return fmt.Errorf("get organization %s: expected enabled=true", id)
			}
			return nil
		},
	}
}
