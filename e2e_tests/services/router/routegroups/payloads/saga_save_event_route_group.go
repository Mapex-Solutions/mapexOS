// Package payloads holds canonical RouteGroupCreate fixtures for the
// router routegroups module.
//
// Every asset MUST belong to at least one route group (RouteGroupIds is
// validated min=1,max=3). Saga journeys that create assets compose
// SagaSaveEventRouteGroup before the asset step so the asset has a valid
// reference. The canonical fixture wires a single save_event router so
// telemetry produced by the asset is persisted by the events service —
// the simplest end-to-end pipeline that still exercises router routing.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/router/routegroups"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// RouteGroupCreateBuilder wraps contracts.RouteGroupCreate so journeys can
// override individual fields without redeclaring the canonical baseline.
type RouteGroupCreateBuilder struct {
	spec contracts.RouteGroupCreate
}

// Build returns the contracts payload ready for POST /api/v1/routegroups.
func (b *RouteGroupCreateBuilder) Build() contracts.RouteGroupCreate {
	return b.spec
}

// WithName overrides the route group name. Useful when the journey wants
// a deterministic identifier the test can assert on.
func (b *RouteGroupCreateBuilder) WithName(name string) *RouteGroupCreateBuilder {
	b.spec.Name = name
	return b
}

// SagaSaveEventRouteGroup returns the canonical save_event-only route group
// for IoT pipeline tests. The single router persists incoming events to
// ClickHouse via the events service, giving the saga a verifiable side
// effect downstream of asset telemetry without dragging triggers or
// workflows into Phase 1.
//
// Defaults:
//   - Version: 1.0.0
//   - Enabled: true
//   - Routers: one save_event router with no Match conditions (matches
//     every event the asset emits)
func SagaSaveEventRouteGroup(runID string) *RouteGroupCreateBuilder {
	return &RouteGroupCreateBuilder{
		spec: contracts.RouteGroupCreate{
			Version:     "1.0.0",
			Name:        fmt.Sprintf("saga-save-event-%s", runID),
			Description: zerovalue.Ptr("Saga route group; persists events through save_event router"),
			Enabled:     zerovalue.Ptr(true),
			Routers: &[]contracts.Router{
				{
					Kind:      "save_event",
					SaveEvent: &contracts.SaveEventData{},
				},
			},
		},
	}
}
