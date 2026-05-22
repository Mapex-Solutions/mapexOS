package control

import (
	"context"
	"fmt"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * LOOP EXECUTOR
 * Iterates over a source array, executing the body for each item.
 * Tracks current index in NodeStates[nodeId]["currentIndex"].
 * Injects loop_item and loop_index into state for body nodes.
 * If currentIndex < totalItems: output ["body"]
 * If currentIndex == totalItems: output ["done"]
 * Validates totalItems <= MaxLoopIterations.
 */

// LoopExecutor iterates over a source array, executing the body for each item.
// Tracks current index in NodeStates and injects loop_item/loop_index into state.
// Validates that total items does not exceed MaxLoopIterations.
type LoopExecutor struct {
	resolver enginePorts.ValueResolverPort
}

// NewLoopExecutor creates a new LoopExecutor with the given value resolver.
func NewLoopExecutor(resolver enginePorts.ValueResolverPort) entities.NodeExecutor {
	return &LoopExecutor{resolver: resolver}
}

// NodeType returns "core/loop".
func (e *LoopExecutor) NodeType() string {
	return "core/loop"
}

// Execute resolves the source array, advances the loop index, and outputs ["body"]
// for the next iteration or ["done"] when all items have been processed.
func (e *LoopExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.LoopNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("loop: missing or invalid config")
	}

	resolved, err := e.resolver.Resolve(
		&cfg.Source,
		execCtx.EventPayload,
		execCtx.State,
		execCtx.NodeOutputs,
		execCtx.ExternalInputs,
	)
	if err != nil {
		return nil, fmt.Errorf("loop: failed to resolve source: %w", err)
	}

	items, ok := resolved.([]interface{})
	if !ok {
		return nil, fmt.Errorf("loop: source must be an array, got %T", resolved)
	}

	totalItems := len(items)
	if totalItems > constants.MaxLoopIterations {
		return nil, fmt.Errorf("%w: %d items", entities.ErrMaxLoopIterations, totalItems)
	}

	currentIndex := 0
	if myState := execCtx.NodeStates[execCtx.NodeID]; myState != nil {
		if v, ok := myState["currentIndex"]; ok {
			currentIndex = toInt(v)
		}
	}

	if currentIndex < 0 || currentIndex >= totalItems {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"done"},
			NodeState: map[string]interface{}{
				"currentIndex": currentIndex,
			},
		}, nil
	}

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"body"},
		NodeState: map[string]interface{}{
			"currentIndex": currentIndex + 1,
		},
		StatePatch: map[string]interface{}{
			"loop_item":  items[currentIndex],
			"loop_index": currentIndex,
		},
		NodeOutput: map[string]interface{}{
			"item":  items[currentIndex],
			"index": currentIndex,
		},
	}, nil
}
