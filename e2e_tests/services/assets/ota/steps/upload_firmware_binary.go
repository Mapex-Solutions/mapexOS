package steps

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	otaPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/payloads"
)

// UploadFirmwareBinary PUTs the firmware bytes to the presigned URL InitFirmware
// returned — direct to object storage, NOT through the assets service (the bytes
// never traverse the API). Uses a plain HTTP client because the presigned URL is
// self-authenticating and expects raw bytes, not the JSON envelope the service
// clients send.
//
// Reads (bag):
//   - BagKeyFirmwareUploadURL  string  presigned PUT URL (set by InitFirmware)
//
// Compensate: no-op — the object is cleaned up when the firmware/plan lifecycle
// tears down; there is nothing local to undo.
func UploadFirmwareBinary() saga.Step {
	return saga.Step{
		Name: "assets/ota.UploadFirmwareBinary",
		Do: func(c *saga.Context) error {
			// Same deterministic artifact InitFirmware declared (built from RunID).
			fw := otaPayloads.NewFirmwareArtifact(c.RunID)
			url := c.MustGetString(BagKeyFirmwareUploadURL)
			req, err := http.NewRequestWithContext(c.Stdctx, http.MethodPut, url, bytes.NewReader(fw.Bytes()))
			if err != nil {
				return fmt.Errorf("build firmware upload request: %w", err)
			}
			req.ContentLength = fw.Size()
			// The presigned PUT is signed against x-amz-checksum-sha256; MinIO returns
			// 403 SignatureDoesNotMatch if this header is missing or does not match.
			req.Header.Set("x-amz-checksum-sha256", fw.SHA256())
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("upload firmware binary: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("upload firmware binary: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}
