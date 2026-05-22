package inline

import (
	"context"
	"fmt"
	"time"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * END EXECUTOR
 * Terminal node of a workflow execution.
 * If terminateWithError is true, resolves errorMessage and returns an ExecutionError.
 * Otherwise, returns empty OutputHandles (signals workflow completion).
 */

// EndExecutor is the terminal node of a workflow execution.
// If terminateWithError is true, it resolves errorMessage and returns an ExecutionError.
// Otherwise, it returns empty OutputHandles signaling workflow completion.
type EndExecutor struct {
	resolver enginePorts.ValueResolverPort
}

// NewEndExecutor creates a new EndExecutor with the given value resolver.
func NewEndExecutor(resolver enginePorts.ValueResolverPort) entities.NodeExecutor {
	return &EndExecutor{resolver: resolver}
}

// NodeType returns "core/end".
func (e *EndExecutor) NodeType() string {
	return "core/end"
}

// Execute terminates the workflow, optionally producing an ExecutionError
// when terminateWithError is configured.
func (e *EndExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.EndNodeConfig)
	if !ok || cfg == nil {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{},
		}, nil
	}

	if !cfg.TerminateWithError {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{},
		}, nil
	}

	errorMsg := cfg.ErrorCode
	if cfg.ErrorMessage.Value != "" {
		resolved, err := e.resolver.Resolve(
			&cfg.ErrorMessage,
			execCtx.EventPayload,
			execCtx.State,
			execCtx.NodeOutputs,
			execCtx.ExternalInputs,
		)
		if err == nil {
			errorMsg = fmt.Sprintf("%v", resolved)
		}
	}

	return &entities.NodeExecutionResult{
		OutputHandles: []string{},
		Error: &entities.ExecutionError{
			Code:      cfg.ErrorCode,
			Message:   errorMsg,
			NodeID:    execCtx.NodeID,
			NodeType:  execCtx.NodeType,
			Timestamp: time.Now(),
		},
	}, nil
}
