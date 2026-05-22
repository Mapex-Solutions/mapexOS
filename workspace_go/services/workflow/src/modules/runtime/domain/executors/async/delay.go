package async

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/runtime/domain/entities"
)

/*
 * DELAY EXECUTOR
 * Suspends execution for a specified duration.
 * Returns NodeState with waitType "timer" and computed expiresAt.
 * The RuntimeService checkpoints to KV and publishes a NATS Schedule for timer delivery.
 */

// DelayExecutor suspends execution for a specified duration.
// Returns a NodeState with waitType "timer" and the computed expiration time.
type DelayExecutor struct{}

// NewDelayExecutor creates a new DelayExecutor.
func NewDelayExecutor() entities.NodeExecutor {
	return &DelayExecutor{}
}

// NodeType returns "core/delay".
func (e *DelayExecutor) NodeType() string {
	return "core/delay"
}

// Execute parses the configured duration and returns a NodeState with waitType "timer"
// that the RuntimeService checkpoints to KV. A NATS Schedule delivers the resume on expiry.
func (e *DelayExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.DelayNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("delay: missing or invalid config")
	}

	duration, err := ParseDuration(cfg.Duration, cfg.Unit)
	if err != nil {
		return nil, fmt.Errorf("delay: %w", err)
	}

	expiresAt := time.Now().Add(duration)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		NodeState: map[string]interface{}{
			"waitType":  "timer",
			"expiresAt": expiresAt,
		},
	}, nil
}

// parseDuration is now consolidated in timeout_helper.go as ParseDuration
