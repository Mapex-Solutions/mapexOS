// Package triggers holds the cross-service contract types and constants
// emitted by the triggers service triggers module.
//
// These constants are the wire-level contract for messages published by
// the triggers service and consumed by other services:
//   - mapexos.events.trigger → events MS (EVENTS-TRIGGER stream) — trigger
//     execution history sink.
//
// Ownership: triggers service (publisher).
// Consumers (Go): events.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/triggers/triggers.
//
// Contracts stay leaf-level — no imports from services/.
package triggers

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// SubjectTriggerEvents is the NATS subject for publishing trigger execution
// events to the events service for ClickHouse persistence. Resolved at
// package init — e.g. "dev.mapexos.events.trigger".
var SubjectTriggerEvents = config.Subject("events", "trigger")

// StreamTriggerEvents is the NATS JetStream stream that carries trigger
// execution history events for the events service consumer (events_trigger).
// Resolved at package init — e.g. "DEV-MAPEXOS-EVENTS-TRIGGERS-LOGS".
var StreamTriggerEvents = config.StreamName("EVENTS", "TRIGGERS-LOGS")

// EventTypeTriggerEvents tags DLQ messages produced by the events_trigger
// consumer.
const EventTypeTriggerEvents = "trigger"
