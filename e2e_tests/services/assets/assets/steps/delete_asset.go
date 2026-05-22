package steps

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// DeleteAsset issues DELETE /api/v1/assets/{id} as an explicit journey
// step (not as Compensate). The MQTT broker journey uses it to prove
// the FANOUT invalidation path: after delete the broker's L1 + L2
// entries are wiped, and a subsequent CONNECT for the same identity
// must fail with auth deny.
//
// Reads (bag):
//   - BagKeyAssetID  string  set by CreateAsset
//
// Writes (bag):
//   - BagKeyAssetDeleted  bool  true so CreateAsset.Compensate no-ops
//
// Compensate: no-op. The delete is the test action; reinstating the
// asset would invalidate the assertion the next step makes.
func DeleteAsset() saga.Step {
	return saga.Step{
		Name: "assets/assets.DeleteAsset",
		Do: func(c *saga.Context) error {
			id := c.MustGetString(BagKeyAssetID)
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, "/api/v1/assets/"+id, nil)
			if err != nil {
				return fmt.Errorf("delete asset %s: %w", id, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				c.Set(BagKeyAssetDeleted, true)
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("delete asset %s: unexpected status %d body=%s", id, resp.StatusCode, string(body))
			}
			c.Set(BagKeyAssetDeleted, true)
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
