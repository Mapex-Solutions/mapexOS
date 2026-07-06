package events

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// NATS subjects and streams for OTA status-history cross-service messages.
// Stream and subject names resolve at package init from GO_ENV so the same
// binary serves multiple environments on a shared NATS cluster.

// SubjectOTAStatus is the NATS subject the Asset MS publishes on each per-device
// OTA status transition (downloading → … → updated/failed). Consumed by the
// events service for ClickHouse history (events_ota_status). Resolved at
// package init — e.g. "dev.mapexos.events.ota.status".
var SubjectOTAStatus = config.Subject("events", "ota.status")

// StreamOTAStatus is the NATS JetStream stream that carries OTA status events
// for the events service consumer (ota_status). Resolved at package init —
// e.g. "DEV-MAPEXOS-EVENTS-OTA-STATUS".
var StreamOTAStatus = config.StreamName("EVENTS", "OTA-STATUS")

// SubjectOTAStatusAdvisory is the STATIC inbound subject where the edges
// publish a device's normalized OTA status advisory (payload OTAStatusAdvisory)
// for the Asset MS to consume: the HTTP gateway from POST /api/v1/ota/status,
// and mapexMQTTBroker forwarding the device-facing `ota/status` topic. Device
// identity + executionId travel in the PAYLOAD, never in the subject.
var SubjectOTAStatusAdvisory = config.Subject("ota", "status.advisory")

// EventTypeOTAStatus tags DLQ messages produced by the ota_status consumer in
// the events service.
const EventTypeOTAStatus = "ota_status"
