package entities

import (
	"time"
)

/**
 * RawEvent represents an unprocessed event from gateways stored in ClickHouse.
 * Raw events are stored for debugging purposes with short retention (1-3 days).
 *
 * This entity mirrors the RawEventDTO structure but includes additional
 * fields required for persistence (RetentionDays, Source, Metadata).
 *
 * ClickHouse table: events_raw
 *
 * Note: Field order matches ClickHouse column order for efficient scanning.
 */
type RawEvent struct {
	// Created is the event occurrence time
	Created time.Time `ch:"created"`

	// EventTrackerId is the UUID for tracking event across the pipeline
	EventTrackerId string `ch:"event_tracker_id"`

	// ThreadId is the data source identifier (DataSource ID)
	ThreadId string `ch:"thread_id"`

	// OrgId is the organization identifier for multi-tenancy
	OrgId string `ch:"org_id"`

	// PathKey is the asset's hierarchical path for organizational queries
	PathKey string `ch:"path_key"`

	// Source identifies the originating gateway (http_gateway, mqtt_gateway, etc.)
	Source string `ch:"source"`

	// Name is a friendly name for the data source (optional)
	Name string `ch:"name"`

	// Description is a description of the data source (optional)
	Description string `ch:"description"`

	// Event contains the raw payload from the gateway (stored as JSON in ClickHouse)
	Event map[string]interface{} `ch:"event"`

	// Metadata contains additional context (headers, etc.) (stored as JSON in ClickHouse)
	Metadata map[string]interface{} `ch:"metadata"`

	// Success indicates if gateway auth/validation succeeded
	Success bool `ch:"success"`

	// Error contains the error message if Success is false
	Error string `ch:"error"`

	// RetentionDays specifies how long this event should be kept
	RetentionDays uint16 `ch:"retention_days"`
}
