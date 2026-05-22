// Package steps holds saga steps that exercise the assets/assettemplates
// module HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
)

type templateCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CreateTemplate POSTs the canonical SagaTemperatureTemplate payload to
// the assets service and publishes the returned id on the bag. The
// payload is built from c.RunID inside Do so the saga assembly does not
// thread runtime values through the constructor.
//
// Reads (bag):
//   - none — auth headers live on the HTTP client; runID lives on Context.
//
// Writes (bag):
//   - BagKeyTemplateID  string  Mongo ObjectID hex of the new template
//
// Compensate: DELETE /api/v1/asset_templates/{id}. The id is read back
// from the bag rather than captured in a closure so Compensate stays
// idempotent.
func CreateTemplate() saga.Step {
	return saga.Step{
		Name: "assets/assettemplates.CreateTemplate",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaTemperatureTemplate(c.RunID).Build()
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/asset_templates", spec)
			if err != nil {
				return fmt.Errorf("create template: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("create template: unexpected status %d", resp.StatusCode)
			}
			var out templateCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-template response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create template: empty id in response")
			}
			c.Set(BagKeyTemplateID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyTemplateID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, "/api/v1/asset_templates/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete template: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete template: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
