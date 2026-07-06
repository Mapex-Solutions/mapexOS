package entities

import (
	"time"
)

// OTAStatusEvent is one per-device OTA status transition persisted to the
// events_ota_status ClickHouse table. It is the history projection of the
// OTAStatusAdvisory the Asset MS publishes on each device progress report
// (downloading -> ... -> updated/failed), kept for audit and rollout analytics.
//
// ClickHouse table: events_ota_status.
//
// Note: field order matches ClickHouse column order for efficient scanning.
type OTAStatusEvent struct {
	// Created is the device report time (advisory Timestamp; now() if absent).
	Created time.Time `ch:"created"`

	// OrgId is the organization identifier for multi-tenancy.
	OrgId string `ch:"org_id"`

	// AssetUuid is the device-facing asset identity carried by the advisory.
	AssetUuid string `ch:"asset_uuid"`

	// PlanId is the OTA plan (rollout) this transition belongs to.
	PlanId string `ch:"plan_id"`

	// OTAExecutionId is the per-device execution record id (stable per device
	// for the whole update).
	OTAExecutionId string `ch:"ota_execution_id"`

	// Status is the device-reported execution state (downloading, downloaded,
	// verified, updating, updated, failed).
	Status string `ch:"status"`

	// Progress is the download/apply percentage (0-100).
	Progress int32 `ch:"progress"`

	// Error carries the failure detail when Status is failed (optional).
	Error string `ch:"error"`

	// Message carries an optional human-readable detail.
	Message string `ch:"message"`

	// RetentionDays specifies how long this row should be kept.
	RetentionDays uint16 `ch:"retention_days"`
}
