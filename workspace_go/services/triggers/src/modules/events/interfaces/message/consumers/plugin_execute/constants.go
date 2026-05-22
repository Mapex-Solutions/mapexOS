package plugin_execute

/**
 * Constants for WorkflowExecute consumer.
 *
 * Receives all execution requests from the Workflow Service (both trigger
 * entities and plugin actions).
 *
 * Subject is re-exported from the workflow service cross-service contract —
 * never redefined locally. Stream is intra-service (consumed only by
 * triggers) and lives in interfaces/message/constants.go. EventType is
 * consumer-local DLQ metadata.
 */

import (
	"triggers/src/modules/events/interfaces/message"

	workflowRuntime "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/runtime"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the JetStream stream consumed by this consumer (intra-service).
var Stream = message.StreamTriggers

// Subject is the NATS subject for workflow-originated execution requests
// (cross-service contract).
var Subject = workflowRuntime.SubjectTriggerWorkflowExecute

// Durable name for this consumer.
var Durable = config.Durable("triggers", "workflow-execute")

// EventType is the DLQ metadata label for this consumer (consumer-local).
const EventType = "workflow.execute"
