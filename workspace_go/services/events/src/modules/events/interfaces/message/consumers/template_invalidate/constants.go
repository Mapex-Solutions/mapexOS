package template_invalidate

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for TemplateInvalidate FANOUT consumer.
 *
 * FANOUT pattern:
 *  - Each service instance receives a copy of the message
 *  - Used for cache invalidation across all replicas
 *  - No queue group (each instance processes independently)
 */

// Stream is the JetStream stream name (re-exported from contracts).
var Stream = contracts.FanoutStreamName

// Subject is the FANOUT NATS subject (re-exported from contracts).
var Subject = contracts.FanoutTemplateSubject

// Durable name for the template_invalidate consumer.
var Durable = config.Durable("events", "template-invalidate")

// EventType is the local DLQ tag for this consumer.
const EventType = "fanout.template.invalidate"
