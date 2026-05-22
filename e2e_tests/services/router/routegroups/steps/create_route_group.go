// Package steps holds saga steps that exercise the router routegroups
// module HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
)

type routeGroupCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// BuilderFn returns the RouteGroupCreateBuilder a CreateRouteGroup
// variant POSTs. Callers pick the variant (save_event vs workflow vs
// trigger) at journey-wire time without touching the step's machinery.
type BuilderFn func(c *saga.Context) *payloads.RouteGroupCreateBuilder

// CreateRouteGroup POSTs the canonical SagaSaveEventRouteGroup payload to
// the router service and publishes the returned id on BagKeyRouteGroupID.
// Kept as the original single-RG variant the IoT MQTT pipeline uses; the
// connectivity-action journeys call CreateRouteGroupAt instead so they
// can stand up two RGs (online + offline) on distinct bag keys.
func CreateRouteGroup() saga.Step {
	return CreateRouteGroupAt(BagKeyRouteGroupID, func(c *saga.Context) *payloads.RouteGroupCreateBuilder {
		return payloads.SagaSaveEventRouteGroup(c.RunID)
	})
}

// CreateRouteGroupAt POSTs the route group built by the supplied
// builder and publishes the returned id at the caller-chosen bag key.
// Lets journeys stand up multiple route groups in the same saga
// (typical: one for online, one for offline) without colliding on
// BagKeyRouteGroupID.
//
// Writes (bag):
//   - <bagKey>  string  Mongo ObjectID hex of the new route group
//
// Compensate: DELETE /api/v1/route_groups/{id}. The id is read back
// from the same bag key so the step value is reusable across runs.
func CreateRouteGroupAt(bagKey string, builder BuilderFn) saga.Step {
	return saga.Step{
		Name: fmt.Sprintf("router/routegroups.CreateRouteGroup[%s]", bagKey),
		Do: func(c *saga.Context) error {
			spec := builder(c).Build()
			resp, err := c.Clients.Router.Raw(c.Stdctx, http.MethodPost, "/api/v1/route_groups", spec)
			if err != nil {
				return fmt.Errorf("create route group: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("create route group: unexpected status %d", resp.StatusCode)
			}
			var out routeGroupCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create route group response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create route group: empty id in response")
			}
			c.Set(bagKey, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(bagKey)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Router.Raw(c.Stdctx, http.MethodDelete, "/api/v1/route_groups/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete route group: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete route group: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
