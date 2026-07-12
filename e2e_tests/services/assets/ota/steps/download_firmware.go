package steps

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
)

// downloadAndVerifyFirmware fetches the firmware from the command's presigned URL
// and verifies its size + sha256 against the command — a real device refuses to
// apply a mismatched image. Shared by both device sims (MQTT + HTTP): the bytes
// come straight from object storage, never through the platform API.
func downloadAndVerifyFirmware(ctx context.Context, cmd downlink.OTAUpdateCommand) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cmd.DownloadURL, nil)
	if err != nil {
		return fmt.Errorf("build firmware download request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download firmware: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download firmware: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read firmware body: %w", err)
	}
	if int64(len(body)) != cmd.Size {
		return fmt.Errorf("firmware size mismatch: got %d want %d", len(body), cmd.Size)
	}
	if got := sha256Hex(body); got != cmd.Checksum {
		return fmt.Errorf("firmware checksum mismatch: got %s want %s", got, cmd.Checksum)
	}
	return nil
}

// sha256Hex returns the lowercase-hex sha256 of b.
func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
