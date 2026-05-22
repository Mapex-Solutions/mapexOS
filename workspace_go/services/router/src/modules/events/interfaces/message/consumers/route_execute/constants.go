package route_execute

import (
	hmContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for RouteExecute consumer.
 *
 * Stream and subject are owned by the cross-service contract in
 * packages/contracts/services/assets/healthmonitor (RouterStream and
 * RouterSubject) — already env-prefixed via the gokit config helpers.
 * Durable is local to this consumer.
 */

// Stream name for route execution events.
var Stream = hmContracts.RouterStream

// Subject for route execution events.
var Subject = hmContracts.RouterSubject

// Durable name for this consumer.
var Durable = config.Durable("router", "execute")

// EventType for DLQ metadata.
const EventType = "route.execute"
