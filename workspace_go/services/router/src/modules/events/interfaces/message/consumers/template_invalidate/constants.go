package template_invalidate

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for TemplateInvalidate FANOUT consumer.
 *
 * FANOUT Pattern:
 * - Each service instance receives a copy of the message
 * - Used for cache invalidation across all replicas
 * - No queue group (each instance processes independently)
 *
 * Published by: Assets service (on template create/update/delete)
 * Consumed by: Router, Events, JS-Executor
 */

// Stream is the NATS stream for FANOUT broadcast messages.
var Stream = contracts.FanoutStreamName

// Subject is the NATS subject for template cache invalidation events.
var Subject = contracts.FanoutTemplateSubject

// Durable name for this consumer.
var Durable = config.Durable("router", "template-invalidate")

// EventType identifies this consumer in DLQ metadata.
const EventType = "fanout.template.invalidate"
