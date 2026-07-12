package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	otaPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/payloads"
)

type firmwareInitResponse struct {
	Data struct {
		FirmwareID string `json:"firmwareId"`
		UploadURL  string `json:"uploadUrl"`
	} `json:"data"`
}

// InitFirmware creates the firmware artifact record and captures the short-TTL
// presigned PUT URL for the direct-to-store upload. The declared size+sha256
// come from the deterministic FirmwareArtifact so CompleteFirmware's confirm and
// the device sim's post-download verify all agree. targetTemplateIDKey names the
// bag key holding the target template id (the migration destination).
//
// Reads (bag):
//   - targetTemplateIDKey  string  the target template id (set by CreateTemplateWithLabel)
//
// Writes (bag):
//   - BagKeyFirmwareID / BagKeyFirmwareUploadURL / BagKeyFirmwareSHA256 /
//     BagKeyFirmwareSize / BagKeyFirmwareVersion
//
// Compensate: no-op — there is no firmware-delete endpoint; the artifact is
// retained by design (re-downloadable from the closed plan), and the plan
// Compensate cancels the plan that owns it.
func InitFirmware(fw *otaPayloads.FirmwareArtifact, targetTemplateIDKey string) saga.Step {
	return saga.Step{
		Name: "assets/ota.InitFirmware",
		Do: func(c *saga.Context) error {
			req := otaDtos.FirmwareInitRequest{
				TargetTemplateID: c.MustGetString(targetTemplateIDKey),
				Version:          fw.Version(),
				Filename:         fw.Filename(),
				Size:             fw.Size(),
				SHA256:           fw.SHA256(),
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/ota/firmware/init", req)
			if err != nil {
				return fmt.Errorf("init firmware: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("init firmware: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out firmwareInitResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode init-firmware response: %w", err)
			}
			if out.Data.FirmwareID == "" || out.Data.UploadURL == "" {
				return fmt.Errorf("init firmware: empty firmwareId/uploadUrl in response")
			}
			c.Set(BagKeyFirmwareID, out.Data.FirmwareID)
			c.Set(BagKeyFirmwareUploadURL, out.Data.UploadURL)
			c.Set(BagKeyFirmwareSHA256, fw.SHA256())
			c.Set(BagKeyFirmwareSize, fw.Size())
			c.Set(BagKeyFirmwareVersion, fw.Version())
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}
