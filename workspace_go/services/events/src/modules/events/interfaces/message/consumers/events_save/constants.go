package events_save

import (
	routerEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/router/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsSave consumer.
 *
 * The wire-level contract (subject/stream) is owned by the router service
 * and declared in packages/contracts/services/router/events. These locals
 * are thin aliases kept so the consumer file stays stable when subjects
 * evolve. Durable is local to this consumer.
 */

// Stream name for processed events from router.
var Stream = routerEvents.StreamEvents

// Subject for save events.
var Subject = routerEvents.SubjectSaveEvent

// Durable name for the events_save consumer.
var Durable = config.Durable("events", "save")

// EventType for DLQ metadata.
const EventType = routerEvents.EventTypeSaveEvent
