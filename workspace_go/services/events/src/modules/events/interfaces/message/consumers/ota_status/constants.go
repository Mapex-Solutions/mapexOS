package ota_status

import (
	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Constants for the OTA status-history consumer.
//
// The wire-level contract (subject/stream) is owned by the OTA bounded context
// and declared in packages/contracts/services/ota/events. These locals are thin
// aliases so the consumer stays stable when subjects evolve. Durable is local.

// Stream name for OTA status-history events.
var Stream = otaEvents.StreamOTAStatus

// Subject for OTA status-history events.
var Subject = otaEvents.SubjectOTAStatus

// Durable name for the ota_status consumer.
var Durable = config.Durable("events", "ota_status")

// EventType for DLQ metadata.
const EventType = otaEvents.EventTypeOTAStatus
