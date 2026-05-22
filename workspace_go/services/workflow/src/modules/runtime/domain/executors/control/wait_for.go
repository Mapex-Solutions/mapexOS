package control

import (
	"context"
	"fmt"
	"time"

	enginePorts "workflow/src/modules/engine/application/ports"
	defPorts "workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/runtime/domain/entities"
	"workflow/src/modules/runtime/domain/executors/async"
)

/*
 * WAIT FOR EXECUTOR
 * Evaluates a state condition. If met: outputs ["matched"] inline.
 * If not met: returns NodeState with waitType "condition" to be polled later.
 */

// WaitForExecutor evaluates a state condition and outputs ["matched"] inline if met,
// or returns a NodeState with waitType "condition" to be polled later if not met.
type WaitForExecutor struct {
	evaluator enginePorts.ConditionEvaluatorPort
}

// NewWaitForExecutor creates a new WaitForExecutor with the given condition evaluator.
func NewWaitForExecutor(evaluator enginePorts.ConditionEvaluatorPort) entities.NodeExecutor {
	return &WaitForExecutor{evaluator: evaluator}
}

// NodeType returns "core/wait_for".
func (e *WaitForExecutor) NodeType() string {
	return "core/wait_for"
}

// Execute builds a condition group from the config, evaluates it, and returns ["matched"]
// immediately if satisfied or a NodeState with waitType "condition" for deferred polling.
func (e *WaitForExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.WaitForNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("wait_for: missing or invalid config")
	}

	conditionGroup := &defPorts.ConditionGroup{
		Logic: defPorts.LogicAND,
		Items: []defPorts.ConditionGroupItem{
			{
				Type: "condition",
				Data: defPorts.ConditionItem{
					Field: defPorts.FieldValue{
						Type:  defPorts.FieldValueState,
						Value: cfg.Field,
					},
					Operator: cfg.Operator,
					Value:    cfg.CompareTo,
				},
			},
		},
	}

	met, err := e.evaluator.EvaluateGroup(
		conditionGroup,
		execCtx.Timezone,
		execCtx.EventPayload,
		execCtx.State,
		execCtx.NodeOutputs,
		execCtx.ExternalInputs,
	)
	if err != nil {
		return nil, fmt.Errorf("wait_for: condition evaluation failed: %w", err)
	}

	if met {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"matched"},
		}, nil
	}

	expiresAt := async.CalculateExpiresAt(execCtx.Timeout, 24*time.Hour)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{},
		NodeState: map[string]interface{}{
			"waitType":     "condition",
			"field":        cfg.Field,
			"operator":     cfg.Operator,
			"compareTo":    fmt.Sprintf("%v", cfg.CompareTo.Value),
			"expiresAt":    expiresAt,
			"enableOutput": async.IsEnableOutput(execCtx.Timeout),
		},
	}, nil
}
