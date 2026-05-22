package ports

import (
	"workflow/src/modules/runtime/domain/entities"
	sharedTypes "workflow/src/shared/types"
)

// Port-level type aliases — expose domain entities through the port boundary.
// Other modules import these types from ports, NEVER from domain/entities directly.

// ResumeMessage is the NATS WORKFLOW-RESUME payload. Exposed here so the
// application layer consumes it via ports instead of importing interfaces/message
// (Hexagonal layering). The intra-service alias in interfaces/message still points to the
// same shared type.
type ResumeMessage = sharedTypes.ResumeMessage

// WorkflowExecution is the aggregate root for a single execution run.
type WorkflowExecution = entities.WorkflowExecution

// ExecutionStatus represents the terminal state of an execution.
type ExecutionStatus = entities.ExecutionStatus

// PathEntry records a single step in the execution path.
type PathEntry = entities.PathEntry

// ExecutionError captures structured error info for failed executions.
type ExecutionError = entities.ExecutionError

// Re-export status constants used by cross-module consumers.
const (
	ExecStatusCreated   = entities.ExecStatusCreated
	ExecStatusRunning   = entities.ExecStatusRunning
	ExecStatusWaiting   = entities.ExecStatusWaiting
	ExecStatusCompleted = entities.ExecStatusCompleted
	ExecStatusFailed    = entities.ExecStatusFailed
	ExecStatusCancelled = entities.ExecStatusCancelled
)
