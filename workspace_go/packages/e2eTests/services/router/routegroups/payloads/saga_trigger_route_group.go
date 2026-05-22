package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/router/routegroups"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaTriggerRouteGroup returns a single-router route group whose
// router kind=trigger targets the supplied trigger. Used by the
// connectivity-action journeys to wire online/offline transitions to
// a trigger execution the saga can then observe via
// /api/v1/events/trigger.
//
// Inputs:
//   - runID      saga run identifier embedded in the name/description
//   - role       human label (e.g. "online" / "offline") for self-
//                describing names in dev listings
//   - triggerID  Mongo ObjectID hex of an existing trigger created by
//                the saga in a previous step
//
// Defaults:
//   - Version: 1.0.0
//   - Enabled: true
//   - Routers: one trigger router with empty metadata (the trigger
//             carries its own resolved config server-side).
func SagaTriggerRouteGroup(runID, role, triggerID string) *RouteGroupCreateBuilder {
	return &RouteGroupCreateBuilder{
		spec: contracts.RouteGroupCreate{
			Version:     "1.0.0",
			Name:        fmt.Sprintf("saga-trigger-%s-%s", role, runID),
			Description: zerovalue.Ptr(fmt.Sprintf("Saga route group; fires trigger on %s", role)),
			Enabled:     zerovalue.Ptr(true),
			Routers: &[]contracts.Router{
				{
					Kind: "trigger",
					Trigger: &contracts.TriggerData{
						TriggerId: triggerID,
					},
				},
			},
		},
	}
}
