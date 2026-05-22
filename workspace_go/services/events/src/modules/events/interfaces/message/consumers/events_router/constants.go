package events_router

import (
	routerEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/router/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsRouter consumer.
 *
 * The wire-level contract (subject/stream) is owned by the router service
 * and declared in packages/contracts/services/router/events. These locals
 * are thin aliases kept so the consumer file stays stable when subjects
 * evolve. Durable is local to this consumer.
 */

// Stream name for router execution history events.
var Stream = routerEvents.StreamRouterHistory

// Subject for router events.
var Subject = routerEvents.SubjectRouterHistory

// Durable name for the events_router consumer.
var Durable = config.Durable("events", "router")

// EventType for DLQ metadata.
const EventType = routerEvents.EventTypeRouterHistory
