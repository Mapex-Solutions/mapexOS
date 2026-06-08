package asserts

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// gatewayCertRequest mirrors the POST /api/v1/gateway_certs body.
type gatewayCertRequest struct {
	AssetUUID string `json:"assetUUID"`
	Force     bool   `json:"force"`
}

// AssertGatewayCertRejected attempts to issue a gateway mTLS cert for the bag's
// asset and expects the assets service to refuse it. Run against an eui-mode
// gateway, it proves the eligibility guard (assertGatewayCertEligible): only
// cert-mode gateways may obtain a cert, so an eui-mode gateway is turned away with
// 409 Conflict and no cert is minted. Without the guard, an eui gateway would get
// a usable cert it has no business holding.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID  string  set by CreateAsset
func AssertGatewayCertRejected() saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertGatewayCertRejected",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			body := gatewayCertRequest{AssetUUID: uuid, Force: false}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/gateway_certs", body)
			if err != nil {
				return fmt.Errorf("issue gateway cert (expect reject) for %s: %w", uuid, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				raw, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("gateway cert for eui-mode %s: expected rejection, got %d body=%s", uuid, resp.StatusCode, string(raw))
			}
			if resp.StatusCode != http.StatusConflict {
				raw, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("gateway cert for eui-mode %s: expected 409 Conflict, got %d body=%s", uuid, resp.StatusCode, string(raw))
			}
			return nil
		},
	}
}
