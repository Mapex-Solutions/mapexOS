// Package steps holds saga steps that exercise the IAM roles module.
package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/roles/payloads"
)

type roleCreateResponse struct {
	Data struct {
		ID    string `json:"id"`
		OrgID string `json:"orgId"`
	} `json:"data"`
}

// CreateRole POSTs the canonical SagaIoTAdminRole payload to the
// mapexIam service and publishes the returned id on the bag. The
// payload is built from c.RunID inside Do so the saga assembly does not
// thread runtime values through the constructor.
//
// Reads (bag):
//   - none — auth headers live on the HTTP client; runID lives on Context.
//
// Writes (bag):
//   - BagKeyRoleID  string  Mongo ObjectID hex of the new role
//
// Compensate: DELETE /api/v1/roles/{id} explicitly. The role is created
// inside the seed parent org (where the bootstrap actor has coverage),
// not the saga child org, so the child-org cascade does not reach it.
// The id is read back from the bag rather than captured in a closure so
// Compensate stays idempotent.
func CreateRole() saga.Step {
	return saga.Step{
		Name: "mapexIam/roles.CreateRole",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaIoTAdminRole(c.RunID).Build()
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodPost, "/api/v1/roles", spec)
			if err != nil {
				return fmt.Errorf("create role: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create role: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out roleCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-role response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create role: empty id in response")
			}
			c.Set(BagKeyRoleID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyRoleID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodDelete, "/api/v1/roles/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete role: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete role: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
