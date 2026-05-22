package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// issueCertRequest mirrors POST /api/v1/mqtt_certs body shape.
type issueCertRequest struct {
	AssetUUID string `json:"assetUUID"`
	Force     bool   `json:"force"`
}

// issueCertResponse mirrors the platform's HTTP envelope wrapping
// the cert payload. CertPEM/KeyPEM/CAChainPEM are `[]byte` on the Go
// contract; the wire JSON carries them as base64 strings, and Go's
// stdlib json decoder un-base64s automatically into the byte slices.
type issueCertResponse struct {
	Serial      string `json:"serial"`
	Fingerprint string `json:"fingerprint"`
	SubjectCN   string `json:"subjectCN"`
	CertPEM     []byte `json:"certPEM"`
	KeyPEM      []byte `json:"keyPEM"`
	CAChainPEM  []byte `json:"caChainPEM"`
}

// IssueCert calls POST /api/v1/mqtt_certs for the saga's asset and
// captures the PEM bundle on the bag. The cert lifecycle reflection
// runs server-side (assets MS persists currentCert + fans out the
// L2 invalidation), so the next read of the asset will surface the
// new cert metadata.
//
// Reads (bag):
//   - BagKeyAssetUUID  string  set by CreateAsset (cert mode)
//
// Writes (bag):
//   - BagKeyAssetCertPEM      []byte  PEM-encoded device cert
//   - BagKeyAssetKeyPEM       []byte  PEM-encoded private key
//   - BagKeyAssetCAChainPEM   []byte  PEM-encoded CA chain (intermediate)
//   - BagKeyAssetCertSerial   string  uppercase hex serial of the issued cert
//
// Compensate: no-op. CreateAsset's compensation deletes the asset
// which cascades to currentCert + revoked rows.
func IssueCert() saga.Step {
	return saga.Step{
		Name: "assets/assets.IssueCert",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			body := issueCertRequest{AssetUUID: uuid, Force: false}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/mqtt_certs", body)
			if err != nil {
				return fmt.Errorf("issue cert for %s: %w", uuid, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				raw, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("issue cert for %s: unexpected status %d body=%s", uuid, resp.StatusCode, string(raw))
			}
			var out issueCertResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode issue-cert response: %w", err)
			}
			if len(out.CertPEM) == 0 || len(out.KeyPEM) == 0 {
				return fmt.Errorf("issue cert for %s: empty cert or key in response", uuid)
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
