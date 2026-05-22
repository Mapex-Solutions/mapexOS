// Package steps holds saga steps that exercise the workflow definitions
// module HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/workflow/definitions/payloads"
)

type definitionCreateResponse struct {
	Data struct {
		ID      string `json:"_id"`
		Version int    `json:"definitionVersion"`
	} `json:"data"`
}

// CreateDefinition POSTs the canonical SagaSimpleDefinition payload to
// the workflow service and publishes id + version on the bag. The
// CreateInstance step reads both to populate the InstanceCreate body
// (definitionId + definitionVersion).
//
// Writes (bag):
//   - BagKeyDefinitionID       string  Mongo ObjectID hex
//   - BagKeyDefinitionVersion  int     version returned by the API
//
// Compensate: DELETE /api/v1/workflow_definitions/{id}. Idempotent; reads id
// back from the bag instead of capturing it in a closure so the Step
// value is safe to reuse across runs.
func CreateDefinition() saga.Step {
	return saga.Step{
		Name: "workflow/definitions.CreateDefinition",
		Do: func(c *saga.Context) error {
			spec := payloads.SagaSimpleDefinition(c.RunID)
			resp, err := c.Clients.Workflow.Raw(c.Stdctx, http.MethodPost, "/api/v1/workflow_definitions", spec)
			if err != nil {
				return fmt.Errorf("create definition: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create definition: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out definitionCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create definition response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create definition: empty id in response")
			}
			c.Set(BagKeyDefinitionID, out.Data.ID)
			c.Set(BagKeyDefinitionVersion, out.Data.Version)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyDefinitionID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Workflow.Raw(c.Stdctx, http.MethodDelete, "/api/v1/workflow_definitions/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete definition: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete definition: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
