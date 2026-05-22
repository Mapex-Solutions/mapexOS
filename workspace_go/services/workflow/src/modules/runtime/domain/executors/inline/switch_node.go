package inline

import (
	"context"
	"fmt"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * SWITCH EXECUTOR
 * Evaluates multiple cases and outputs matching case handles.
 * matchMode "first": first match → ["case_{id}"]
 * matchMode "all":   all matches → ["case_1", "case_3", ...]
 * No match → ["default"]
 */

// SwitchExecutor evaluates multiple cases and outputs matching case handles.
// In "first" match mode it returns the first match; in "all" mode it returns all matches.
// Falls back to ["default"] when no case matches.
type SwitchExecutor struct {
	evaluator enginePorts.ConditionEvaluatorPort
}

// NewSwitchExecutor creates a new SwitchExecutor with the given condition evaluator.
func NewSwitchExecutor(evaluator enginePorts.ConditionEvaluatorPort) entities.NodeExecutor {
	return &SwitchExecutor{evaluator: evaluator}
}

// NodeType returns "core/switch".
func (e *SwitchExecutor) NodeType() string {
	return "core/switch"
}

// Execute evaluates each case condition and returns the matching output handles
// based on the configured match mode.
func (e *SwitchExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.SwitchNodeConfig)
	if !ok || cfg == nil {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"default"},
		}, nil
	}

	matchMode := cfg.MatchMode
	if matchMode == "" {
		matchMode = "first"
	}

	var matched []string

	for _, c := range cfg.Cases {
		result, err := e.evaluator.EvaluateGroup(
			&c.Condition,
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
			matched = append(matched, fmt.Sprintf("case_%s", c.ID))
			if matchMode == "first" {
				break
			}
		}
	}

	if len(matched) == 0 {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"default"},
		}, nil
	}

	return &entities.NodeExecutionResult{
		OutputHandles: matched,
	}, nil
}
