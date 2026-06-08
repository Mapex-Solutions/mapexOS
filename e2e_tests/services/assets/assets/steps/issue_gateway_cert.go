package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// IssueGatewayCert calls POST /api/v1/gateway_certs for the saga's gateway asset
// and captures the PEM bundle + serial on the bag. It reuses the same request and
// response shapes as IssueCert (the endpoint shares the signer and CN=assetUUID),
// the only difference being the route and the server-side eligibility guard that
// requires a cert-mode gateway asset.
//
// The cert lifecycle reflection runs server-side synchronously: the assets MS
// persists currentCert and rewrites the mapex-asset-auth projection before the
// request returns, so the next read of /internal/asset-auth/:uuid already carries
// the new serial.
//
// Reads (bag):
//   - BagKeyAssetUUID  string  set by CreateAsset (cert-mode gateway)
//
// Writes (bag):
//   - BagKeyAssetCertPEM      []byte  PEM-encoded gateway cert
//   - BagKeyAssetKeyPEM       []byte  PEM-encoded private key
//   - BagKeyAssetCAChainPEM   []byte  PEM-encoded CA chain
//   - BagKeyAssetCertSerial   string  uppercase hex serial of the issued cert
//
// Compensate: no-op. CreateAsset's compensation deletes the asset, which clears
// currentCert.
func IssueGatewayCert() saga.Step {
	return saga.Step{
		Name: "assets/assets.IssueGatewayCert",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			body := issueCertRequest{AssetUUID: uuid, Force: false}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/gateway_certs", body)
			if err != nil {
				return fmt.Errorf("issue gateway cert for %s: %w", uuid, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				raw, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("issue gateway cert for %s: unexpected status %d body=%s", uuid, resp.StatusCode, string(raw))
			}
			var out issueCertResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode issue-gateway-cert response: %w", err)
			}
			if len(out.CertPEM) == 0 || len(out.KeyPEM) == 0 {
				return fmt.Errorf("issue gateway cert for %s: empty cert or key in response", uuid)
			}
			if out.Serial == "" {
				return fmt.Errorf("issue gateway cert for %s: empty serial in response", uuid)
			}
			c.Set(BagKeyAssetCertPEM, out.CertPEM)
			c.Set(BagKeyAssetKeyPEM, out.KeyPEM)
			c.Set(BagKeyAssetCAChainPEM, out.CAChainPEM)
			c.Set(BagKeyAssetCertSerial, out.Serial)
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
