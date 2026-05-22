package async

import (
	"context"
	"fmt"
	"time"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * SUBWORKFLOW EXECUTOR
 * Triggers a child workflow and suspends until it completes.
 * Validates depth < MaxSubworkflowDepth to prevent infinite recursion.
 * Resolves inputMappings via ValueResolverPort.
 * Returns NodeState with waitType "callback" — the RuntimeService publishes
 * to WORKFLOW-EXECUTION stream with mode=subworkflow.
 */

// SubworkflowExecutor triggers a child workflow and suspends until it completes.
// Validates depth against MaxSubworkflowDepth to prevent infinite recursion.
type SubworkflowExecutor struct {
	resolver enginePorts.ValueResolverPort
}

// NewSubworkflowExecutor creates a new SubworkflowExecutor with the given value resolver.
func NewSubworkflowExecutor(resolver enginePorts.ValueResolverPort) entities.NodeExecutor {
	return &SubworkflowExecutor{resolver: resolver}
}

// NodeType returns "core/subworkflow".
func (e *SubworkflowExecutor) NodeType() string {
	return "core/subworkflow"
}

// Execute resolves input mappings, validates recursion depth, and returns a NodeState
// with waitType "callback" for the RuntimeService to dispatch via WORKFLOW-EXECUTION stream.
func (e *SubworkflowExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.SubworkflowNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("subworkflow: missing or invalid config")
	}

	currentDepth := execCtx.Depth
	if currentDepth >= constants.MaxSubworkflowDepth {
		return nil, fmt.Errorf("%w: depth=%d", entities.ErrMaxSubworkflowDepth, currentDepth)
	}

	inputData := make(map[string]interface{}, len(cfg.InputMappings))
	for _, mapping := range cfg.InputMappings {
		resolved, err := e.resolver.Resolve(
			&mapping.Value,
			execCtx.EventPayload,
			execCtx.State,
			execCtx.NodeOutputs,
			execCtx.ExternalInputs,
		)
		if err != nil {
			return nil, fmt.Errorf("subworkflow: failed to resolve input %q: %w", mapping.ChildParamName, err)
		}
		inputData[mapping.ChildParamName] = resolved
	}

	expiresAt := CalculateExpiresAt(execCtx.Timeout, 1*time.Hour)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		NodeState: map[string]interface{}{
			"waitType":       "callback",
			"workflowId":     cfg.WorkflowID,
			"workflowName":   cfg.WorkflowName,
			"executionMode":  cfg.ExecutionMode,
			"inputData":      inputData,
			"outputMappings": cfg.OutputMappings,
			"timeout":        cfg.Timeout,
			"depth":          currentDepth + 1,
			"expiresAt":      expiresAt,
			"enableOutput":   IsEnableOutput(execCtx.Timeout),
		},
	}, nil
}
