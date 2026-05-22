package async

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/runtime/domain/entities"
)

/*
 * CODE EXECUTOR
 * Suspends execution to run user-defined JavaScript code externally.
 * Returns NodeState with waitType "callback" containing script and timeout.
 * The RuntimeService publishes the script to the WORKFLOW-JS-CODE stream
 * and the code runner sends the result back via WORKFLOW-RESUME.
 */

// CodeExecutor suspends execution to run user-defined JavaScript code externally.
// Returns a NodeState with waitType "callback" for the code runner response.
type CodeExecutor struct{}

// NewCodeExecutor creates a new CodeExecutor.
func NewCodeExecutor() entities.NodeExecutor {
	return &CodeExecutor{}
}

// NodeType returns "core/code".
func (e *CodeExecutor) NodeType() string {
	return "core/code"
}

// Execute builds a NodeState with waitType "callback" containing the script and timeout,
// delegating actual code execution to the external WORKFLOW-JS-CODE stream.
func (e *CodeExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.CodeNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("code: missing or invalid config")
	}

	expiresAt := CalculateExpiresAt(execCtx.Timeout, 30*time.Second)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		NodeState: map[string]interface{}{
			"waitType":     "callback",
			"timeout":      cfg.Timeout,
			"expiresAt":    expiresAt,
			"enableOutput": IsEnableOutput(execCtx.Timeout),
		},
	}, nil
}
