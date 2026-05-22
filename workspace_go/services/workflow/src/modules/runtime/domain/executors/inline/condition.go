package inline

import (
	"context"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * CONDITION EXECUTOR
 * Evaluates a ConditionGroup and outputs ["true"] or ["false"].
 * Delegates to ConditionEvaluatorPort (engine module) for the actual evaluation.
 */

// ConditionExecutor evaluates a ConditionGroup and outputs ["true"] or ["false"].
// Delegates to ConditionEvaluatorPort (engine module) for the actual evaluation.
type ConditionExecutor struct {
	evaluator enginePorts.ConditionEvaluatorPort
}

// NewConditionExecutor creates a new ConditionExecutor with the given condition evaluator.
func NewConditionExecutor(evaluator enginePorts.ConditionEvaluatorPort) entities.NodeExecutor {
	return &ConditionExecutor{evaluator: evaluator}
}

// NodeType returns "core/condition".
func (e *ConditionExecutor) NodeType() string {
	return "core/condition"
}

// Execute evaluates the condition groups and returns ["true"] or ["false"] output handles.
func (e *ConditionExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.ConditionNodeConfig)
	if !ok || cfg == nil {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"false"},
		}, nil
	}

	result, err := e.evaluator.EvaluateGroup(
		&cfg.Condition,
		execCtx.Timezone,
		execCtx.EventPayload,
		execCtx.State,
		execCtx.NodeOutputs,
		execCtx.ExternalInputs,
	)
	if err != nil {
		return nil, err
	}

	if result {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"true"},
		}, nil
	}
	return &entities.NodeExecutionResult{
		OutputHandles: []string{"false"},
	}, nil
}
