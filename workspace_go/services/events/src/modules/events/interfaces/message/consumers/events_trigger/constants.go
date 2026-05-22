package events_trigger

import (
	triggersContract "github.com/Mapex-Solutions/MapexOS/contracts/services/triggers/triggers"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsTrigger consumer.
 *
 * The wire-level contract (subject/stream) is owned by the triggers service
 * and declared in packages/contracts/services/triggers/triggers. These
 * locals are thin aliases kept so the consumer file stays stable when
 * subjects evolve. Durable is local to this consumer.
 */

// Stream name for trigger execution history events.
var Stream = triggersContract.StreamTriggerEvents

// Subject for trigger events.
var Subject = triggersContract.SubjectTriggerEvents

// Durable name for the events_trigger consumer.
var Durable = config.Durable("events", "trigger")

// EventType for DLQ metadata.
const EventType = triggersContract.EventTypeTriggerEvents
