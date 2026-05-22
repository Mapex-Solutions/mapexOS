package entities

import (
	"time"
)

/**
 * DLQEvent represents a Dead Letter Queue event stored in ClickHouse.
 * DLQ events are stored for debugging and reprocessing failed messages.
 *
 * This entity mirrors the DLQMessage structure from the NATS package
 * with additional fields required for persistence.
 *
 * ClickHouse table: events_dlq
 *
 * Note: Field order matches ClickHouse column order for efficient scanning.
 */
type DLQEvent struct {
	// Created when the message was sent to DLQ
	Created time.Time `ch:"created"`

	// EventTrackerId is the UUID for end-to-end event tracking across services
	EventTrackerId string `ch:"event_tracker_id"`

	// ID is the unique identifier for this DLQ message
	ID string `ch:"id"`

	// OrgId is the organization identifier for multi-tenant filtering (MANDATORY)
	OrgId string `ch:"org_id"`

	// PathKey is the hierarchical path for multi-tenant range queries (MANDATORY)
	PathKey string `ch:"path_key"`

	// ServiceName identifies the service that sent to DLQ (e.g., "events-service", "router-service")
	ServiceName string `ch:"service_name"`

	// ServiceType categorizes the service (e.g., "processor", "gateway", "worker")
	ServiceType string `ch:"service_type"`

	// EventType describes what kind of events this consumer handles (e.g., "raw", "jsexec", "route.execute")
	EventType string `ch:"event_type"`

	// OriginalSubject is the NATS subject the message was originally sent to
	OriginalSubject string `ch:"original_subject"`

	// OriginalStream is the JetStream stream the message came from
	OriginalStream string `ch:"original_stream"`

	// OriginalData contains the original message payload (JSON string)
	OriginalData string `ch:"original_data"`

	// OriginalHeaders contains the original message headers (JSON string)
	OriginalHeaders string `ch:"original_headers"`

	// LastError contains the last error message that caused the DLQ
	LastError string `ch:"last_error"`

	// ErrorCount is the number of errors/retries before sending to DLQ
	ErrorCount uint32 `ch:"error_count"`

	// FirstDelivery is the timestamp of the first delivery attempt
	FirstDelivery time.Time `ch:"first_delivery"`

	// LastDelivery is the timestamp of the last delivery attempt
	LastDelivery time.Time `ch:"last_delivery"`

	// TotalDeliveries is the total number of delivery attempts
	TotalDeliveries uint32 `ch:"total_deliveries"`

	// ConsumerName is the name of the consumer that sent to DLQ
	ConsumerName string `ch:"consumer_name"`

	// RetentionDays specifies how long this event should be kept
	RetentionDays uint16 `ch:"retention_days"`
}

// DLQServiceCount represents a service type count from a GROUP BY query.
type DLQServiceCount struct {
	ServiceType string `ch:"service_type"`
	Count       uint64 `ch:"count"`
}
