package ports

import (
	"context"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// ExecuteResult contains the outcome of a workflow execution triggered via HTTP.
type ExecuteResult struct {
	WorkflowUUID string      `json:"workflowUUID"`
	Status       string      `json:"status"`
	ErrorInfo    interface{} `json:"errorInfo,omitempty"`
}

// RuntimeServicePort defines the contract for workflow DAG execution operations.
type RuntimeServicePort interface {
	// HandleResume processes a workflow resume message from NATS.
	HandleResume(msg *natsModel.Message)

	// HandleExecution processes a workflow execution command from NATS (WORKFLOW-EXECUTION stream).
	// Dispatches by mode: newInstance, signal, signalOrStart, subworkflow.
	HandleExecution(msg *natsModel.Message)

	// HandleScheduleFire processes a fired NATS Schedule delivery on WORKFLOW-SCHEDULE
	// and re-publishes the body to WORKFLOW-RESUME (workflow.resume.timer.{instanceId}).
	// Handles Ack/Nack/Reject internally based on the outcome.
	HandleScheduleFire(msg *natsModel.Message)

	// ExecuteByInstanceID starts a workflow execution for the given instance.
	// Loads instance, validates, prepares definition, runs DAG walker inline.
	// Returns the execution UUID, final status, and error info if failed.
	// Used by the HTTP execute endpoint — same logic as NATS handleNewInstance but without message dependency.
	ExecuteByInstanceID(ctx context.Context, instanceID string, eventPayload map[string]interface{}, workflowUUID string) (*ExecuteResult, error)
}
