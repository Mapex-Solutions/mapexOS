package events

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// NATS subjects for http_gateway events outbound cross-service messages.
// Stream and subject names resolve at package init from GO_ENV so the same
// binary serves multiple environments on a shared NATS cluster.

// SubjectProcessorJSExecute is the NATS subject used by http_gateway to
// dispatch HTTP-ingested events to the js-executor pipeline (workspace_js).
// Payload: ProcessorExecutePayload. Resolved at package init —
// e.g. "dev.mapexos.processor.js.execute".
var SubjectProcessorJSExecute = config.Subject("processor", "js.execute")

// SubjectEventsRaw is the NATS subject used to publish raw events from
// HTTP/MQTT gateways and the js-executor pipeline. Consumed by the
// events service for ClickHouse raw debug storage. Resolved at package
// init — e.g. "dev.mapexos.events.raw".
var SubjectEventsRaw = config.Subject("events", "raw")

// StreamEventsRaw is the NATS JetStream stream that carries raw gateway
// events for the events service consumer (events_raw). Resolved at
// package init — e.g. "DEV-MAPEXOS-EVENTS-RAW".
var StreamEventsRaw = config.StreamName("EVENTS", "RAW")

// EventTypeEventsRaw tags DLQ messages produced by the events_raw
// consumer in the events service.
const EventTypeEventsRaw = "raw"
