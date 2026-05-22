package entities

import "time"

// AssetStatusEvent is one row in the asset_status_history ClickHouse table.
// Field order mirrors the table's column order for efficient scanning;
// `ch:` tags match the SQL schema seeded in deployment/.../asset_status_history.sql.
//
// Populated by the asset_status save-handler from the FLAT persistence payload
// produced by the healthmonitor alert publisher (subject
// mapexos.events.asset_status_save).
type AssetStatusEvent struct {
	// Created is the transition timestamp (UTC) produced by the publisher.
	Created time.Time `ch:"created"`

	// OrgId is the organization identifier for multi-tenant filtering.
	OrgId string `ch:"org_id"`

	// PathKey is the asset's hierarchical path for org-based queries.
	PathKey string `ch:"path_key"`

	// AssetUUID is the device identifier.
	AssetUUID string `ch:"asset_uuid"`

	// AssetName is the human-readable asset name at transition time.
	AssetName string `ch:"asset_name"`

	// EventId is a UUID generated once per transition, shared with the
	// corresponding route.execute message (if any).
	EventId string `ch:"event_id"`

	// EventType is "offline" or "online".
	EventType string `ch:"event_type"`

	// LastSeenAt is the last heartbeat timestamp (offline transitions only).
	// Nil on online transitions or when the asset has never connected.
	LastSeenAt *time.Time `ch:"last_seen_at"`

	// ThresholdMinutes is the offline threshold that was in effect at
	// transition time (from HealthMonitorConfig.ThresholdMinutes).
	ThresholdMinutes uint16 `ch:"threshold_minutes"`

	// MissCount is how many consecutive heartbeat scans missed before the
	// offline decision fired (offline transitions only).
	MissCount uint16 `ch:"miss_count"`
}
