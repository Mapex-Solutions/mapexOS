package inline

import (
	"context"

	"workflow/src/modules/runtime/domain/entities"
)

/*
 * START EXECUTOR
 * Entry point of every workflow execution. Passthrough that outputs ["out"].
 * State initialization with variable defaults is done by the RuntimeService
 * when creating the WorkflowInstance, not here.
 */

// StartExecutor is the entry point of every workflow execution.
// It is a passthrough that outputs ["out"].
type StartExecutor struct{}

// NewStartExecutor creates a new StartExecutor.
func NewStartExecutor() entities.NodeExecutor {
	return &StartExecutor{}
}

// NodeType returns "core/start".
func (e *StartExecutor) NodeType() string {
	return "core/start"
}

// Execute returns a passthrough result with output handle ["out"].
func (e *StartExecutor) Execute(_ context.Context, _ *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
	}, nil
}
