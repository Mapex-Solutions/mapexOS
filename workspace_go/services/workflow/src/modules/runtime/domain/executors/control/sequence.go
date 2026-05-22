package control

import (
	"context"
	"fmt"

	"workflow/src/modules/runtime/domain/entities"
)

/*
 * SEQUENCE EXECUTOR
 * Executes steps in order: step_1 → step_2 → ... → done.
 * Tracks current step in NodeStates[nodeId]["currentStep"].
 * If currentStep < totalSteps: output ["step_{currentStep+1}"]
 * If currentStep == totalSteps: output ["done"]
 */

// SequenceExecutor executes steps in order: step_1, step_2, ..., done.
// Tracks the current step in NodeStates.
type SequenceExecutor struct{}

// NewSequenceExecutor creates a new SequenceExecutor.
func NewSequenceExecutor() entities.NodeExecutor {
	return &SequenceExecutor{}
}

// NodeType returns "core/sequence".
func (e *SequenceExecutor) NodeType() string {
	return "core/sequence"
}

// Execute advances the step counter and outputs ["step_{n}"] for the next step
// or ["done"] when all steps have been completed.
func (e *SequenceExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.SequenceNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("sequence: missing or invalid config")
	}

	currentStep := 0
	if myState := execCtx.NodeStates[execCtx.NodeID]; myState != nil {
		if v, ok := myState["currentStep"]; ok {
			currentStep = toInt(v)
		}
	}

	currentStep++

	nodeState := map[string]interface{}{
		"currentStep": currentStep,
	}

	if currentStep < 0 || currentStep >= cfg.Steps {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"done"},
			NodeState:     nodeState,
		}, nil
	}

	return &entities.NodeExecutionResult{
		OutputHandles: []string{fmt.Sprintf("step_%d", currentStep+1)},
		NodeState:     nodeState,
	}, nil
}
