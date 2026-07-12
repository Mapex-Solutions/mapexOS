package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

type planDownloadResponse struct {
	Data struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
	} `json:"data"`
}

// DownloadPlanFirmware mints a short-TTL presigned GET URL for the plan's
// firmware and asserts it is served (200 + non-empty url) — the operator
// re-download path. It does not fetch the bytes here (the device sims do that).
//
// Reads (bag):
//   - BagKeyPlanID  string  the plan id (set by CreatePlan)
//
// Compensate: no-op — read-only.
func DownloadPlanFirmware() saga.Step {
	return saga.Step{
		Name: "assets/ota.DownloadPlanFirmware",
		Do: func(c *saga.Context) error {
			id := c.MustGetString(BagKeyPlanID)
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodGet, "/api/v1/ota/plans/"+id+"/firmware/download", nil)
			if err != nil {
				return fmt.Errorf("download plan firmware: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("download plan firmware: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out planDownloadResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode download-plan-firmware response: %w", err)
			}
			if out.Data.URL == "" {
				return fmt.Errorf("download plan firmware: empty url in response")
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}
