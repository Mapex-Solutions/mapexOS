package healthmonitor

// AdminAssetUUID is the URL params DTO for admin endpoints addressed by
// the asset device id (e.g. POST /internal/health-monitor/:assetUUID/force-offline).
type AdminAssetUUID struct {
	AssetUUID string `params:"assetUUID" validate:"required,min=1"`
}

// AdminForceOfflineRequest is the optional body for the force-offline
// admin endpoint. Reason is propagated into the alert log line so e2e
// runs can be traced back to the originating test or operator action.
type AdminForceOfflineRequest struct {
	Reason string `json:"reason,omitempty" validate:"omitempty,max=200"`
}
