package workflow_execution

import (
	runtimeContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/runtime"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream name for workflow execution commands.
// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
var Stream = runtimeContract.StreamWorkflowExecution

// Subject pattern — routing is done by mode in payload, never by subject.
// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
var Subject = runtimeContract.SubjectWorkflowExecution

// Durable name for this consumer.
var Durable = config.Durable("workflow", "execution")

// EventType for DLQ metadata.
const EventType = "workflow-execution"

// DLQServiceType is the service-type tag attached to DLQ messages produced
// by this consumer. Identifies the workflow service in cross-service DLQ
// inspection tools.
const DLQServiceType = "workflow"
