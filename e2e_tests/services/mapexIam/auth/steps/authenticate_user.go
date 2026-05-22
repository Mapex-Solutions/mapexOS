// Package steps holds saga steps that exercise the IAM auth module.
package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/auth"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	onboardingPL "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/onboarding/payloads"
	onboardingSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/onboarding/steps"
	orgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/steps"
)

type loginResponse struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

// AuthenticateUser POSTs to /auth/login with the saga test user
// credentials and reconfigures every HTTP client with the resulting JWT
// plus the saga org id as X-Org-Context. After this step every subsequent
// HTTP call inside the saga is authenticated as the saga test user.
//
// Note: the auth module mounts /auth/login at the service root (no
// /api/v1 prefix). Other IAM endpoints — organizations, roles,
// onboarding — DO live under /api/v1.
//
// Reads (bag):
//   - onboardingSteps.BagKeyUserEmail  string  set by CreateUserWithMemberships
//   - orgSteps.BagKeyOrgID             string  set by CreateOrganization
//
// Writes (bag):
//   - BagKeyUserJWT  string  bearer token attached to subsequent calls
//
// Compensate: server-side session is invalidated when the org cascade
// removes the user during teardown, so the local Compensate is a no-op.
func AuthenticateUser() saga.Step {
	return saga.Step{
		Name: "mapexIam/auth.AuthenticateUser",
		Do: func(c *saga.Context) error {
			email := c.MustGetString(onboardingSteps.BagKeyUserEmail)
			orgID := c.MustGetString(orgSteps.BagKeyOrgID)

			payload := contracts.LoginDTO{
				Email:         email,
				Password:      onboardingPL.SagaIoTAdminUserPassword,
				KeepConnected: false,
			}
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodPost, "/auth/login", payload)
			if err != nil {
				return fmt.Errorf("login: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("login: unexpected status %d", resp.StatusCode)
			}
			var out loginResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode login response: %w", err)
			}
			if out.Data.AccessToken == "" {
				return fmt.Errorf("login: empty access token in response")
			}

			c.Clients.SetBearer(out.Data.AccessToken)
			c.Clients.SetOrgContext(orgID)
			c.Set(BagKeyUserJWT, out.Data.AccessToken)
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
