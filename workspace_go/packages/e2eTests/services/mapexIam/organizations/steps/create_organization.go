// Package steps holds saga steps that exercise the IAM organizations
// module HTTP endpoints. Steps perform actions against the live stack and
// publish output keys to the saga bag for downstream consumers; they do
// not verify outcomes — that is the asserts/ package responsibility.
package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/payloads"
)

// orgCreateResponse captures the subset of fields the saga needs from the
// API response. The API wraps the entity inside a StandardResponse envelope
// (status/errors/data). Decoding only the fields used downstream keeps this
// step decoupled from envelope evolution.
type orgCreateResponse struct {
	Data struct {
		ID      string `json:"id"`
		PathKey string `json:"pathKey"`
	} `json:"data"`
}

// CreateOrganization POSTs the canonical SagaTestOrg payload to the
// mapexIam service and publishes the returned id and path key on the
// bag. The payload is built from c.RunID inside Do so the saga assembly
// does not thread runtime values through the constructor.
//
// Reads (bag):
//   - none — auth headers live on the HTTP client; runID lives on Context.
//
// Writes (bag):
//   - BagKeyOrgID       string  Mongo ObjectID hex of the new org
//   - BagKeyOrgPathKey  string  hierarchical path key
//
// Compensate: DELETE /api/v1/organizations/{id}. Cascade-deletes children
// (users, groups, roles, memberships) so a single Compensate cleans up the
// IAM bootstrap performed by the saga. The id is read back from the bag
// rather than captured in a closure so Compensate stays idempotent.
func CreateOrganization() saga.Step {
	return saga.Step{
		Name: "mapexIam/organizations.CreateOrganization",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaTestCustomerOrg(c.RunID).Build()
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodPost, "/api/v1/organizations", spec)
			if err != nil {
				return fmt.Errorf("create organization: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("create organization: unexpected status %d", resp.StatusCode)
			}
			var out orgCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-org response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create organization: empty id in response")
			}
			c.Set(BagKeyOrgID, out.Data.ID)
			c.Set(BagKeyOrgPathKey, out.Data.PathKey)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyOrgID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.HTTP.Raw(c.Stdctx, http.MethodDelete, "/api/v1/organizations/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete organization: %w", err)
			}
			defer resp.Body.Close()
			// Cleanup is best-effort: 404 means another step already removed
			// the org; treat anything outside 2xx + 404 as a real failure so
			// we surface unexpected server errors during teardown.
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete organization: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
