// Package asserts holds saga oracles for the IAM auth module.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	authSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/auth/steps"
	orgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/steps"
)

// AssertJwtValid performs a structural check on the bag-stored JWT: it
// must be three base64url segments. We deliberately do not verify the
// signature here because validation against the live stack happens
// implicitly when the saga uses the token on subsequent endpoints; failure
// to authenticate would surface there with a clear 401 status. The shape
// check protects against the login step ever publishing a malformed value.
//
// Reads (bag):
//   - authSteps.BagKeyUserJWT  string  set by SeedAdminLogin or AuthenticateUser
func AssertJwtValid() saga.Assert {
	return saga.Assert{
		Name: "mapexIam/auth.AssertJwtValid",
		Check: func(c *saga.Context) error {
			token := c.MustGetString(authSteps.BagKeyUserJWT)
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				return fmt.Errorf("jwt invalid: expected 3 segments, got %d", len(parts))
			}
			for i, p := range parts {
				if p == "" {
					return fmt.Errorf("jwt invalid: segment %d is empty", i)
				}
			}
			return nil
		},
	}
}

type coverageResponse struct {
	Data struct {
		Organizations []struct {
			// The coverage endpoint returns the org id under "id" (not
			// "orgId"); decoding the field with the correct JSON tag
			// keeps the assertion honest when the response shape stays
			// stable.
			ID string `json:"id"`
		} `json:"organizations"`
	} `json:"data"`
}

// AssertJwtHasOrgContext fetches /auth/users/me/coverage and verifies the
// authenticated user has access to the saga organization. This proves the
// onboarding step's NewGroup membership took effect: coverage cache picks
// up the role-grant inside the org, the API returns it, the saga keeps
// going. A coverage miss here points at membership wiring — not at the
// JWT shape.
//
// Reads (bag):
//   - orgSteps.BagKeyOrgID  string  set by SeedAdminLogin or CreateOrganization
func AssertJwtHasOrgContext() saga.Assert {
	return saga.Assert{
		Name: "mapexIam/auth.AssertJwtHasOrgContext",
		Check: func(c *saga.Context) error {
			orgID := c.MustGetString(orgSteps.BagKeyOrgID)
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodGet, "/auth/users/me/coverage", nil)
			if err != nil {
				return fmt.Errorf("get coverage: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("get coverage: unexpected status %d", resp.StatusCode)
			}
			var out coverageResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode coverage response: %w", err)
			}
			for _, o := range out.Data.Organizations {
				if o.ID == orgID {
					return nil
				}
			}
			return fmt.Errorf("coverage missing org %s; got %d orgs", orgID, len(out.Data.Organizations))
		},
	}
}
