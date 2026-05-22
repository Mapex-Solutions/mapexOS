package events_raw

import (
	httpGatewayEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsRaw consumer.
 *
 * The wire-level contract (subject/stream) is owned by the http_gateway
 * service (and also published by the js-executor pipeline) and declared
 * in packages/contracts/services/http_gateway/events. These locals are
 * thin aliases kept so the consumer file stays stable when subjects evolve.
 * Durable is local to this consumer.
 */

// Stream name for raw events from HTTP/MQTT gateways.
var Stream = httpGatewayEvents.StreamEventsRaw

// Subject for raw events.
var Subject = httpGatewayEvents.SubjectEventsRaw

// Durable name for the events_raw consumer.
var Durable = config.Durable("events", "raw")

// EventType for DLQ metadata.
const EventType = httpGatewayEvents.EventTypeEventsRaw
