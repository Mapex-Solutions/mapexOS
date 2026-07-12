package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

type assetTemplateResponse struct {
	Data struct {
		AssetTemplateID string `json:"assetTemplateId"`
	} `json:"data"`
}

// AssertAssetTemplateSwitched polls GET /api/v1/assets/{assetId} until the
// asset's assetTemplateId equals the target template id (the OTA migration
// destination), or the timeout elapses. The switch happens when the device
// reports "updated" and the status handler moves the asset onto the target
// template. Default budget 30s / 500ms.
//
// Reads (bag):
//   - assetIDKey            string  the asset id (set by CreateAsset*)
//   - targetTemplateIDKey   string  the target template id (set by CreateTemplateWithLabel)
func AssertAssetTemplateSwitched(assetIDKey, targetTemplateIDKey string) saga.Assert {
	return saga.Assert{
		Name: "assets/ota.AssertAssetTemplateSwitched",
		Check: func(c *saga.Context) error {
			assetID := c.MustGetString(assetIDKey)
			want := c.MustGetString(targetTemplateIDKey)
			timeout := 30 * time.Second
			tick := 500 * time.Millisecond
			deadline := time.Now().Add(timeout)
			var lastSeen string
			for {
				tid, err := fetchAssetTemplateID(c, assetID)
				if err == nil {
					lastSeen = tid
					if tid == want {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("asset %s template did not switch to %q within %v (last seen %q)", assetID, want, timeout, lastSeen)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("asset %s template poll cancelled: %w", assetID, c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

func fetchAssetTemplateID(c *saga.Context, id string) (string, error) {
	resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, "/api/v1/assets/"+id, nil)
	if err != nil {
		return "", fmt.Errorf("get asset: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get asset: unexpected status %d", resp.StatusCode)
	}
	var out assetTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode get-asset response: %w", err)
	}
	return out.Data.AssetTemplateID, nil
}
