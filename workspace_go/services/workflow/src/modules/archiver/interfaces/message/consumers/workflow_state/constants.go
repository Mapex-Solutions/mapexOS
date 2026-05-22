package workflow_state

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the JetStream stream for workflow state lifecycle events.
// Resolved at package init — e.g. "DEV-MAPEXOS-WORKFLOW-STATE".
var Stream = config.StreamName("WORKFLOW", "STATE")

// Subject is the subject pattern for all state lifecycle events. Resolved
// at package init — e.g. "dev.mapexos.workflow.state.>".
var Subject = config.Subject("workflow", "state") + ".>"

// Durable name for this archiver consumer.
var Durable = config.Durable("workflow", "state-archiver")

// EventType identifies this consumer type for DLQ metadata.
const EventType = "state.archiver"

// DLQServiceType is the service-type tag attached to DLQ messages produced
// by this consumer. Identifies the workflow service in cross-service DLQ
// inspection tools.
const DLQServiceType = "workflow"

// BatchSize is the number of state messages fetched per batch. Sized for
// per-message ACK/NACK semantics via BatchMessageHandlerV2.
const BatchSize = 500
