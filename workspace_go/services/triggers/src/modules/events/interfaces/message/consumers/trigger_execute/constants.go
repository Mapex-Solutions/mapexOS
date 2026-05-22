package trigger_execute

/**
 * Constants for TriggerExecute consumer.
 *
 * Subject is re-exported from the router service cross-service contract —
 * never redefined locally. Stream is intra-service (consumed only by
 * triggers) and lives in interfaces/message/constants.go. EventType is
 * consumer-local DLQ metadata.
 */

import (
	"triggers/src/modules/events/interfaces/message"

	routerEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/router/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the JetStream stream consumed by this consumer (intra-service).
var Stream = message.StreamTriggers

// Subject is the NATS subject for trigger execution events from the Router
// Service (cross-service contract).
var Subject = routerEvents.SubjectTriggerRouterExecute

// Durable name for this consumer.
var Durable = config.Durable("triggers", "router-execute")

// EventType is the DLQ metadata label for this consumer (consumer-local).
const EventType = "trigger.execute"
