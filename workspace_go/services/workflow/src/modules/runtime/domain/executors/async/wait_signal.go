package async

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/runtime/domain/entities"
)

/*
 * WAIT SIGNAL EXECUTOR
 * Suspends execution until an external signal is received.
 * Returns NodeState with waitType "signal", signalName, expiresAt, and enableOutput.
 * Timeout is configured at the node level (node.timeout), not inside config.
 */

// WaitSignalExecutor suspends execution until an external signal is received.
// Returns a NodeState with waitType "signal" and the signal name.
type WaitSignalExecutor struct{}

// NewWaitSignalExecutor creates a new WaitSignalExecutor.
func NewWaitSignalExecutor() entities.NodeExecutor {
	return &WaitSignalExecutor{}
}

// NodeType returns "core/wait_signal".
func (e *WaitSignalExecutor) NodeType() string {
	return "core/wait_signal"
}

// Execute builds a NodeState with waitType "signal" that the RuntimeService exposes
// via an HTTP endpoint for external signal delivery.
func (e *WaitSignalExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.WaitSignalNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("wait_signal: missing or invalid config")
	}

	expiresAt := CalculateExpiresAt(execCtx.Timeout, 24*time.Hour)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		NodeState: map[string]interface{}{
			"waitType":     "signal",
			"signalName":   cfg.SignalName,
			"expiresAt":    expiresAt,
			"enableOutput": IsEnableOutput(execCtx.Timeout),
		},
	}, nil
}
