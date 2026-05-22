package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/router/routegroups"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaWorkflowRouteGroup returns a single-router route group whose
// router kind=workflow targets the supplied workflow instance in
// newInstance mode. Used by the connectivity-action journeys to wire
// online/offline transitions to a workflow execution that the saga
// can then observe via /api/v1/events/workflow.
//
// Inputs:
//   - runID       saga run identifier embedded in the name/description
//   - role        human label (e.g. "online" / "offline") so listing
//                 the RGs in dev produces self-describing names
//   - instanceID  Mongo ObjectID hex of an existing workflow instance
//                 created by the saga in a previous step
//
// Defaults:
//   - Version:    1.0.0
//   - Enabled:    true
//   - Routers:    one workflow router, mode=newInstance, data.instanceId=instanceID
//
// The workflow data is the minimum the router contract accepts for
// mode=newInstance — instanceId only; the workflow runtime fills the
// rest from the instance config on dispatch.
func SagaWorkflowRouteGroup(runID, role, instanceID string) *RouteGroupCreateBuilder {
	return &RouteGroupCreateBuilder{
		spec: contracts.RouteGroupCreate{
			Version:     "1.0.0",
			Name:        fmt.Sprintf("saga-workflow-%s-%s", role, runID),
			Description: zerovalue.Ptr(fmt.Sprintf("Saga route group; fires workflow instance on %s", role)),
			Enabled:     zerovalue.Ptr(true),
			Routers: &[]contracts.Router{
				{
					Kind: "workflow",
					Workflow: &contracts.WorkflowData{
						Mode: "newInstance",
						Data: map[string]any{
							"instanceId": instanceID,
						},
					},
				},
			},
		},
	}
}
