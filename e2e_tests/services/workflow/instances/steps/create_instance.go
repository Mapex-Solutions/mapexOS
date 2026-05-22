// Package steps holds saga steps that exercise the workflow instances
// module HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	definitionSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/workflow/definitions/steps"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/workflow/instances/payloads"
)

type instanceCreateResponse struct {
	Data struct {
		ID string `json:"_id"`
	} `json:"data"`
}

// CreateInstance POSTs the canonical SagaSimpleInstance payload to the
// workflow service, injecting definitionId + definitionVersion from
// the bag so the instance binds to the definition just created.
//
// Reads (bag):
//   - definitionSteps.BagKeyDefinitionID       string  set by CreateDefinition
//   - definitionSteps.BagKeyDefinitionVersion  int     set by CreateDefinition
//
// Writes (bag):
//   - BagKeyInstanceID  string  Mongo ObjectID hex
//
// Compensate: DELETE /api/v1/workflow_instances/{id}. Idempotent; reads id back
// from the bag.
func CreateInstance() saga.Step {
	return saga.Step{
		Name: "workflow/instances.CreateInstance",
		Do: func(c *saga.Context) error {
			defID := c.MustGetString(definitionSteps.BagKeyDefinitionID)
			defVerRaw, ok := c.Get(definitionSteps.BagKeyDefinitionVersion)
			if !ok {
				return fmt.Errorf("missing bag key %q", definitionSteps.BagKeyDefinitionVersion)
			}
			defVer, ok := defVerRaw.(int)
			if !ok {
				return fmt.Errorf("bag key %q is not int (%T)", definitionSteps.BagKeyDefinitionVersion, defVerRaw)
			}

			spec := payloads.SagaSimpleInstance(c.RunID)
			spec["definitionId"] = defID
			spec["definitionVersion"] = defVer
			// Keep definitionName aligned with the definition name we
			// parameterized on runID so the operator listing instances
			// for this saga run sees a consistent label across the
			// definition + instance pair.
			spec["definitionName"] = fmt.Sprintf("saga-workflow-def-%s", c.RunID)

			resp, err := c.Clients.Workflow.Raw(c.Stdctx, http.MethodPost, "/api/v1/workflow_instances", spec)
			if err != nil {
				return fmt.Errorf("create instance: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create instance: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out instanceCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create instance response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create instance: empty id in response")
			}
			c.Set(BagKeyInstanceID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyInstanceID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Workflow.Raw(c.Stdctx, http.MethodDelete, "/api/v1/workflow_instances/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete instance: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete instance: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
