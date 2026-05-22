// Package asserts holds saga oracles for the assets/assets module.
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

type assetGetResponse struct {
	Data struct {
		ID            string   `json:"id"`
		AssetUUID     string   `json:"assetUUID"`
		Enabled       bool     `json:"enabled"`
		RouteGroupIds []string `json:"routeGroupIds"`
	} `json:"data"`
}

// AssertAssetExists fetches the asset by id from the bag and verifies the
// API returns the entity with enabled=true and at least one route group
// attached. The route-group check guards against a regression where the
// service silently drops the binding during persistence.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetID  string  set by CreateAsset
func AssertAssetExists() saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertAssetExists",
		Check: func(c *saga.Context) error {
			id := c.MustGetString(assetSteps.BagKeyAssetID)
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, "/api/v1/assets/"+id, nil)
			if err != nil {
				return fmt.Errorf("get asset: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("get asset: unexpected status %d", resp.StatusCode)
			}
			var out assetGetResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode get-asset response: %w", err)
			}
			if out.Data.ID != id {
				return fmt.Errorf("get asset: id mismatch want %q got %q", id, out.Data.ID)
			}
			if !out.Data.Enabled {
				return fmt.Errorf("get asset %s: expected enabled=true", id)
			}
			if len(out.Data.RouteGroupIds) == 0 {
				return fmt.Errorf("get asset %s: expected at least one route group", id)
			}
			return nil
		},
	}
}
