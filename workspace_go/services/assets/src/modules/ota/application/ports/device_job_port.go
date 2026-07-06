package ports

import "context"

// DeviceJobPort serves an HTTP device's pending OTA job (the device poll path,
// exposed via the HTTP gateway → the assets internal route). Implemented by the
// reconciler — it owns the dispatch semantics: presence gating and pacing do
// not apply to polls, but fresh URL minting and attempt accounting do.
type DeviceJobPort interface {
	// PendingJob returns the device's pending OTA command with a freshly minted
	// download URL, or nil when the device has no actionable execution.
	PendingJob(ctx context.Context, assetUUID string) (*OTACommandPayload, error)
}
