package ports

import (
	"context"

	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
)

// OTAJobsPort fetches a device's pending OTA job from the Asset MS internal
// API. Driven port — implemented by the assets HTTP client adapter — so the
// application layer never depends on the transport.
type OTAJobsPort interface {
	// FetchPendingJob returns the asset's pending OTA command (with a fresh
	// presigned download URL minted by the Asset MS) or nil when none.
	FetchPendingJob(ctx context.Context, assetUUID string) (*downlink.OTAUpdateCommand, error)
}
