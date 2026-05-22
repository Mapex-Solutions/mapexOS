package dtos

import (
	"encoding/json"
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * Events DLQ DTOs
 * Dead Letter Queue events from all services
 */

// EventsDLQQuery represents query parameters for listing DLQ events.
// Uses cursor-based pagination for efficient querying of large datasets.
//
// Standard fields (from CursorQueryDTO):
//   - Cursor: timestamp to start from (RFC3339 format)
//   - Direction: "next" (older items) or "prev" (newer items)
//   - Limit: max items to return (default: 20, max: 100)
//   - SortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - IncludeChildren: include child orgs hierarchically (default: false)
//
// Module-specific filters:
//   - ServiceName: filter by service name (e.g., "events-service", "router-service")
//   - ServiceType: filter by service type (e.g., "processor", "gateway")
//   - EventType: filter by event type (e.g., "raw", "jsexec", "route.execute")
//   - StartTime: filter events after this timestamp
//   - EndTime: filter events before this timestamp
type EventsDLQQuery struct {
	query.CursorQueryDTO

	// Module-specific filters
	EventTrackerId *string    `query:"eventTrackerId" validate:"omitempty"`
	ServiceName    *string    `query:"serviceName" validate:"omitempty"`
	ServiceType    *string    `query:"serviceType" validate:"omitempty"`
	EventType      *string    `query:"eventType" validate:"omitempty"`
	LastError      *string    `query:"lastError" validate:"omitempty"`
	StartTime      *time.Time `query:"startTime" validate:"omitempty"`
	EndTime        *time.Time `query:"endTime" validate:"omitempty"`
}

// EventsDLQResponse represents a DLQ event response.
type EventsDLQResponse struct {
	// Created is when the message was sent to DLQ
	Created time.Time `json:"created"`

	// EventTrackerId is the UUID for tracking event across the pipeline
	EventTrackerId string `json:"eventTrackerId,omitempty"`

	// ID is the unique identifier for this DLQ message
	ID string `json:"id"`

	// OrgId is the organization identifier (MANDATORY for multi-tenant filtering)
	OrgId string `json:"orgId"`

	// PathKey is the hierarchical path (MANDATORY for multi-tenant filtering)
	PathKey string `json:"pathKey,omitempty"`

	// ServiceName identifies the service that sent to DLQ
	ServiceName string `json:"serviceName"`

	// ServiceType categorizes the service
	ServiceType string `json:"serviceType"`

	// EventType describes what kind of events this consumer handles
	EventType string `json:"eventType"`

	// OriginalSubject is the NATS subject the message was originally sent to
	OriginalSubject string `json:"originalSubject"`

	// OriginalStream is the JetStream stream the message came from
	OriginalStream string `json:"originalStream"`

	// OriginalData contains the original message payload (JSON string)
	OriginalData string `json:"originalData"`

	// OriginalHeaders contains the original message headers (JSON string)
	OriginalHeaders string `json:"originalHeaders,omitempty"`

	// LastError contains the last error message that caused the DLQ
	LastError string `json:"lastError"`

	// ErrorCount is the number of errors/retries before sending to DLQ
	ErrorCount uint32 `json:"errorCount"`

	// FirstDelivery is the timestamp of the first delivery attempt
	FirstDelivery time.Time `json:"firstDelivery"`

	// LastDelivery is the timestamp of the last delivery attempt
	LastDelivery time.Time `json:"lastDelivery"`

	// TotalDeliveries is the total number of delivery attempts
	TotalDeliveries uint32 `json:"totalDeliveries"`

	// ConsumerName is the name of the consumer that sent to DLQ
	ConsumerName string `json:"consumerName"`

	// RetentionDays specifies how long this event should be kept
	RetentionDays uint16 `json:"retentionDays"`
}

// EventsDLQCursorResult represents the cursor-paginated response for DLQ events.
type EventsDLQCursorResult struct {
	// Items contains the list of DLQ events
	Items []EventsDLQResponse `json:"items"`

	// NextCursor is the timestamp to use for fetching the next page
	NextCursor *time.Time `json:"nextCursor,omitempty"`

	// PrevCursor is the timestamp to use for fetching the previous page
	PrevCursor *time.Time `json:"prevCursor,omitempty"`

	// HasNext indicates if there are more items after the current page
	HasNext bool `json:"hasNext"`

	// HasPrevious indicates if there are more items before the current page
	HasPrevious bool `json:"hasPrevious"`
}

// EventsDLQCountsQuery represents query parameters for counting DLQ entries by service type.
type EventsDLQCountsQuery struct {
	query.BaseQueryDTO
	StartTime *time.Time `query:"startTime" validate:"omitempty"`
	EndTime   *time.Time `query:"endTime" validate:"omitempty"`
}

// EventsDLQServiceCount represents a single service type count.
type EventsDLQServiceCount struct {
	ServiceType string `json:"serviceType"`
	Count       uint64 `json:"count"`
}

// EventsDLQCountsResult represents the response for DLQ counts by service type.
type EventsDLQCountsResult struct {
	Counts []EventsDLQServiceCount `json:"counts"`
	Total  uint64                  `json:"total"`
}

// DLQEventIncomingDTO represents the payload received from NATS MAPEXOS-DLQ stream.
// This is the incoming message structure with all DLQ metadata.
type DLQEventIncomingDTO struct {
	// ID is the unique identifier
	ID string `json:"id"`

	// EventTrackerId is the UUID for tracking event across the pipeline
	EventTrackerId string `json:"eventTrackerId"`

	// Tenant context (MANDATORY for multi-tenant filtering)
	OrgId   string `json:"orgId"`
	PathKey string `json:"pathKey"`

	// Service context (for filtering)
	ServiceName string `json:"serviceName"`
	ServiceType string `json:"serviceType"`
	EventType   string `json:"eventType"`

	// Original message
	OriginalSubject string            `json:"originalSubject"`
	OriginalStream  string            `json:"originalStream"`
	OriginalData    json.RawMessage   `json:"originalData"`
	OriginalHeaders map[string]string `json:"originalHeaders,omitempty"`

	// Error information
	LastError  string `json:"lastError"`
	ErrorCount int    `json:"errorCount"`

	// Delivery tracking
	FirstDelivery   time.Time `json:"firstDelivery"`
	LastDelivery    time.Time `json:"lastDelivery"`
	TotalDeliveries int       `json:"totalDeliveries"`

	// Consumer context
	ConsumerName string `json:"consumerName"`

	// Timestamps
	SentToDLQAt time.Time `json:"sentToDLQAt"`
}
