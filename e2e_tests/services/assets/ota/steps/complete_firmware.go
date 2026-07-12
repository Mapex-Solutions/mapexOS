package steps

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// CompleteFirmware confirms the uploaded object (size + sha256) and marks the
// firmware READY. It returns 200 only after the object is confirmed in storage,
// so it must run after UploadFirmwareBinary.
//
// Reads (bag):
//   - BagKeyFirmwareID  string  firmware id (set by InitFirmware)
//
// Compensate: no-op — no un-complete endpoint; the plan Compensate cancels the
// plan that consumes the firmware.
func CompleteFirmware() saga.Step {
	return saga.Step{
		Name: "assets/ota.CompleteFirmware",
		Do: func(c *saga.Context) error {
			id := c.MustGetString(BagKeyFirmwareID)
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/ota/firmware/"+id+"/complete", nil)
			if err != nil {
				return fmt.Errorf("complete firmware: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("complete firmware: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}
