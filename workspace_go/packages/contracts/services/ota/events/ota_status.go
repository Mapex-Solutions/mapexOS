package events

import "time"

// OTAStatusAdvisory is the normalized per-device OTA status payload published on
// SubjectOTAStatus each time a device reports progress on an OTA update. The
// Asset MS produces it (from the device's MQTT/HTTP status report); the events
// service consumes it into the events_ota_status ClickHouse history table.
//
// Status is one of the device-reported execution states: "downloading",
// "downloaded", "verified", "updating", "updated", "failed". Progress is the
// download/apply percentage (0–100). Error/Message are populated on failure or
// to carry a human-readable detail.
type OTAStatusAdvisory struct {
	OrgID          string    `json:"orgId"`
	AssetUUID      string    `json:"assetUUID"`
	PlanID         string    `json:"planId"`
	OTAExecutionID string    `json:"otaExecutionId"`
	Status         string    `json:"status"`
	Progress       int32     `json:"progress"`
	Error          string    `json:"error,omitempty"`
	Message        string    `json:"message,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}
