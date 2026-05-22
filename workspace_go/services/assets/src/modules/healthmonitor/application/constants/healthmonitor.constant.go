package constants

// Health status values surfaced on AssetResponse.healthStatus and persisted
// in Mongo. Values are part of the cross-module surface — the assets module
// uses the same constants when enriching responses from Redis.
const (
	StatusOnline  = "online"
	StatusOffline = "offline"
	StatusUnknown = "unknown"
)
