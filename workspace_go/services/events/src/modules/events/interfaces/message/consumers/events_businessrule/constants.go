package events_businessrule

import (
	jsExecutorScripts "github.com/Mapex-Solutions/MapexOS/contracts/services/js_executor/scripts"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsBusinessRule consumer.
 *
 * The wire-level contract (subject/stream) is owned by the js-executor
 * service (workspace_js) and declared in
 * packages/contracts/services/js_executor/scripts. These locals are thin
 * aliases kept so the consumer file stays stable when subjects evolve.
 * Durable is local to this consumer.
 */

// Stream name for business rule execution history events.
var Stream = jsExecutorScripts.StreamEventsBusinessRule

// Subject for business rule events.
var Subject = jsExecutorScripts.SubjectEventsBusinessRule

// Durable name for the events_businessrule consumer.
var Durable = config.Durable("events", "businessrule")

// EventType for DLQ metadata.
const EventType = jsExecutorScripts.EventTypeEventsBusinessRule
