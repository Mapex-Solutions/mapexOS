package control

import (
	"context"
	"fmt"

	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"

	typeconv "github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
)

/*
 * MERGE EXECUTOR
 * Joins parallel branches back into a single execution path.
 * Strategy:
 *   "all"   — waits for ALL branches to complete
 *   "any"   — proceeds when ANY branch completes
 *   "first" — proceeds when the FIRST branch completes
 *
 * Branch count is read from NodeStates[nodeId]["branchCount"].
 * If condition is met: outputs ["out"]
 * If not met: returns empty handles (wait for more branches).
 */

// MergeExecutor joins parallel branches back into a single execution path.
// Supports strategies "all" (wait for all), "any" (proceed on any), and "first" (proceed on first).
type MergeExecutor struct{}

// NewMergeExecutor creates a new MergeExecutor.
func NewMergeExecutor() entities.NodeExecutor {
	return &MergeExecutor{}
}

// NodeType returns "core/merge".
func (e *MergeExecutor) NodeType() string {
	return "core/merge"
}

// Execute tracks completed branch count in NodeStates and returns ["out"]
// when the merge strategy condition is satisfied, or empty handles to wait for more branches.
func (e *MergeExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.MergeNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("merge: missing or invalid config")
	}

	if cfg.Branches <= 0 {
		return nil, fmt.Errorf("merge: branches must be > 0, got %d", cfg.Branches)
	}

	strategy := cfg.Strategy
	if strategy == "" {
		strategy = constants.MergeStrategyAll
	}

	completedCount := 0
	if myState := execCtx.NodeStates[execCtx.NodeID]; myState != nil {
		if v, ok := myState[constants.BranchCountKey]; ok {
			completedCount = toInt(v)
		}
	}

	completedCount++

	shouldProceed := false
	switch strategy {
	case constants.MergeStrategyAll:
		shouldProceed = completedCount >= cfg.Branches
	case constants.MergeStrategyAny, constants.MergeStrategyFirst:
		shouldProceed = completedCount >= 1
	}

	nodeState := map[string]interface{}{
		constants.BranchCountKey: completedCount,
	}

	if !shouldProceed {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{},
			NodeState:     nodeState,
		}, nil
	}

	return &entities.NodeExecutionResult{
		OutputHandles: []string{constants.OutputHandleOut},
		NodeState:     nodeState,
	}, nil
}

// toInt converts a value to int, returning 0 on failure.
func toInt(v interface{}) int {
	i, _ := typeconv.ToInt64(v)
	return int(i)
}
