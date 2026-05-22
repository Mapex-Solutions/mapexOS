package workflow_resume

import (
	runtimeContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/runtime"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for WorkflowResume consumer
 */

// Stream name for workflow resume messages.
// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
var Stream = runtimeContract.StreamWorkflowResume

// Subject for workflow resume messages.
// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
var Subject = runtimeContract.SubjectWorkflowResume

// Durable name for this consumer.
var Durable = config.Durable("workflow", "resume")

// EventType for DLQ metadata.
const EventType = "workflow-resume"

// DLQServiceType is the service-type tag attached to DLQ messages produced
// by this consumer. Identifies the workflow service in cross-service DLQ
// inspection tools.
const DLQServiceType = "workflow"
