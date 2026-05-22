package events_dlq

import (
	dlqContract "github.com/Mapex-Solutions/MapexOS/contracts/common/dlq"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for EventsDLQ consumer.
 *
 * The wire-level contract (subject/stream) is platform-wide and declared
 * in packages/contracts/common/dlq. These locals are thin aliases kept
 * so the consumer file stays stable when the DLQ pipeline evolves.
 * Durable is local to this consumer.
 */

// Stream name for Dead Letter Queue events.
var Stream = dlqContract.Stream

// Subject for DLQ events.
var Subject = dlqContract.Subject

// Durable name for the events_dlq consumer.
var Durable = config.Durable("events", "dlq")
