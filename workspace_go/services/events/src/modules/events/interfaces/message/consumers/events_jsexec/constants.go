package events_jsexec

import (
	jsExecutorScripts "github.com/Mapex-Solutions/MapexOS/contracts/services/js_executor/scripts"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsJsExec consumer.
 *
 * The wire-level contract (subject/stream) is owned by the js-executor
 * service (workspace_js) and declared in
 * packages/contracts/services/js_executor/scripts. These locals are thin
 * aliases kept so the consumer file stays stable when subjects evolve.
 * Durable is local to this consumer.
 */

// Stream name for JS Executor debug events.
var Stream = jsExecutorScripts.StreamEventsJsExecutor

// Subject for JS Executor events.
var Subject = jsExecutorScripts.SubjectEventsJsExecutor

// Durable name for the events_jsexec consumer.
var Durable = config.Durable("events", "jsexec")

// EventType for DLQ metadata.
const EventType = jsExecutorScripts.EventTypeEventsJsExecutor
