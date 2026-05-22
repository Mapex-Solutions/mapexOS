// Package steps holds saga steps that exercise the IAM onboarding
// orchestrator endpoint.
package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/onboarding/payloads"
	roleSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/roles/steps"
)

type onboardingResponse struct {
	Data struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Groups []struct {
			ID string `json:"id"`
		} `json:"groups"`
	} `json:"data"`
}

// CreateUserWithMemberships POSTs the canonical SagaIoTAdminUser payload,
// which provisions user + group + membership atomically, and publishes
// the returned ids on the bag. The payload is built from c.RunID and the
// role id read from the bag so the saga assembly does not thread runtime
// values through the constructor.
//
// Reads (bag):
//   - roleSteps.BagKeyRoleID  string  set by CreateRole
//
// Writes (bag):
//   - BagKeyUserID     string  Mongo ObjectID hex of the created user
//   - BagKeyUserEmail  string  email used at creation time
//   - BagKeyGroupID    string  first group id when the orchestrator created
//                              a new group alongside the user (NewGroup mode)
//
// Compensate: cascade-cleanup happens via the org delete in the
// CreateOrganization step's Compensate. This step is a no-op on rollback
// to avoid 404 noise during teardown.
func CreateUserWithMemberships() saga.Step {
	return saga.Step{
		Name: "mapexIam/onboarding.CreateUserWithMemberships",
		Do: func(c *saga.Context) error {
			roleID := c.MustGetString(roleSteps.BagKeyRoleID)
			spec := payloads.SagaIoTAdminUser(c.RunID, roleID).Build()

			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodPost, "/api/v1/onboarding/users", spec)
			if err != nil {
				return fmt.Errorf("onboarding create user: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("onboarding create user: unexpected status %d", resp.StatusCode)
			}
			var out onboardingResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode onboarding response: %w", err)
			}
			if out.Data.User.ID == "" {
				return fmt.Errorf("onboarding create user: empty user id in response")
			}
			c.Set(BagKeyUserID, out.Data.User.ID)
			c.Set(BagKeyUserEmail, out.Data.User.Email)
			if len(out.Data.Groups) > 0 {
				c.Set(BagKeyGroupID, out.Data.Groups[0].ID)
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			// Org cascade removes user + group + membership; leave this as
			// a no-op so we do not race the cascading delete.
			return nil
		},
	}
}
