package control

import (
	"context"
	"fmt"

	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * FANOUT EXECUTOR
 * Splits execution into N parallel branches.
 * Outputs ["out_1", "out_2", ..., "out_N"] (dynamic handles).
 * Validates branches <= MaxFanoutBranches.
 * The RuntimeService spawns the branches — this executor only declares intent.
 */

// FanoutExecutor splits execution into N parallel branches.
// Outputs dynamic handles ["out_1", "out_2", ..., "out_N"] and validates
// that branches does not exceed MaxFanoutBranches.
type FanoutExecutor struct{}

// NewFanoutExecutor creates a new FanoutExecutor.
func NewFanoutExecutor() entities.NodeExecutor {
	return &FanoutExecutor{}
}

// NodeType returns "core/fanout".
func (e *FanoutExecutor) NodeType() string {
	return "core/fanout"
}

// Execute declares the parallel branch intent by returning N output handles.
// The RuntimeService is responsible for spawning the actual branches.
func (e *FanoutExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.FanoutNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("fanout: missing or invalid config")
	}

	if cfg.Branches <= 0 {
		return nil, fmt.Errorf("fanout: branches must be > 0")
	}
	if cfg.Branches > constants.MaxFanoutBranches {
		return nil, fmt.Errorf("%w: %d branches", entities.ErrMaxFanoutBranches, cfg.Branches)
	}

	handles := make([]string, cfg.Branches)
	for i := range handles {
		handles[i] = fmt.Sprintf("out_%d", i+1)
	}

	return &entities.NodeExecutionResult{
		OutputHandles: handles,
	}, nil
}
