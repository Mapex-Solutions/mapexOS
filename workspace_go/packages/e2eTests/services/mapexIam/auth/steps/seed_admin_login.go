package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	orgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/steps"
)

// SeedAdminLogin authenticates as the canonical seed admin user shipped by
// mongodb-init (constants.RootUserEmail / constants.RootUserPassword),
// propagates the resulting bearer plus the seed root org as
// X-Org-Context across every per-service client in the saga ClientSet,
// and publishes the (jwt, orgID) pair on the bag.
//
// Every journey starts with this step (or an equivalent that creates a
// distinct tenant first). Centralising it here keeps the building-block
// rule the journey README codifies: cross-journey reuse only at
// services/{svc}/{mod}/{steps,asserts,payloads}, never by importing
// another journey's phase.
//
// Reads (bag): none.
// Writes (bag):
//   - BagKeyUserJWT             string  bearer attached to every HTTP client
//   - orgSteps.BagKeyOrgID      string  MapexosOrgID; the seed root used as saga org
//
// Compensate: server-side session is short-lived; nothing local to undo.
func SeedAdminLogin() saga.Step {
	return saga.Step{
		Name: "mapexIam/auth.SeedAdminLogin",
		Do: func(c *saga.Context) error {
			payload := map[string]any{
				"email":         constants.RootUserEmail,
				"password":      constants.RootUserPassword,
				"keepConnected": false,
			}
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodPost, "/auth/login", payload)
			if err != nil {
				return fmt.Errorf("seed admin login: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("seed admin login: unexpected status %d", resp.StatusCode)
			}

			var out struct {
				Data struct {
					AccessToken string `json:"access_token"`
				} `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode seed admin login response: %w", err)
			}
			if out.Data.AccessToken == "" {
				return fmt.Errorf("seed admin login: empty access_token in response")
			}

			c.Clients.SetBearer(out.Data.AccessToken)
			c.Clients.SetOrgContext(constants.MapexosOrgID)
			c.Set(BagKeyUserJWT, out.Data.AccessToken)
			c.Set(orgSteps.BagKeyOrgID, constants.MapexosOrgID)
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
