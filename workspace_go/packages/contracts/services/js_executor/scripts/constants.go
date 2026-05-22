// Package scripts holds the cross-service contract constants emitted by
// the js-executor service scripts module (workspace_js).
//
// These subject constants are the wire-level contract for messages
// published by the js-executor batch pipeline and consumed by other
// services:
//   - ${env}.mapexos.events.logs.jsexecutor → events MS — JS execution debug logs.
//   - ${env}.mapexos.events.businessrule    → events MS — business rule execution history.
//
// Ownership: js-executor service (publisher, workspace_js).
// Consumers (Go): events.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/jsexecutor/scripts.
//
// Contracts stay leaf-level — no imports from services/.
package scripts

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// SubjectEventsJsExecutor is the NATS subject for publishing js-executor
// debug execution logs to the events service for ClickHouse persistence.
// Resolved at package init — e.g. "dev.mapexos.events.logs.jsexecutor".
var SubjectEventsJsExecutor = config.Subject("events", "logs.jsexecutor")

// StreamEventsJsExecutor is the NATS JetStream stream that carries
// js-executor debug logs for the events service consumer (events_jsexec).
// Resolved at package init — e.g. "DEV-MAPEXOS-EVENTS-JSEXECUTOR-LOGS".
var StreamEventsJsExecutor = config.StreamName("EVENTS", "JSEXECUTOR-LOGS")

// EventTypeEventsJsExecutor tags DLQ messages produced by the events_jsexec
// consumer in the events service.
const EventTypeEventsJsExecutor = "jsexec"

// SubjectEventsBusinessRule is the NATS subject for publishing business
// rule execution history events to the events service for ClickHouse
// persistence. Resolved at package init — e.g. "dev.mapexos.events.businessrule".
var SubjectEventsBusinessRule = config.Subject("events", "businessrule")

// StreamEventsBusinessRule is the NATS JetStream stream that carries
// business rule execution history events for the events service consumer
// (events_businessrule). Resolved at package init —
// e.g. "DEV-MAPEXOS-EVENTS-BUSINESSRULE".
var StreamEventsBusinessRule = config.StreamName("EVENTS", "BUSINESSRULE")

// EventTypeEventsBusinessRule tags DLQ messages produced by the
// events_businessrule consumer in the events service.
const EventTypeEventsBusinessRule = "businessrule"
