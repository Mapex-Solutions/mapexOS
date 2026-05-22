package inline

import (
	"context"
	"fmt"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/entities"

	typeconv "github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
)

/*
 * SET STATE EXECUTOR
 * Modifies workflow instance state using resolved values.
 * Operations: set, increment, decrement, append, remove
 * - remove: deletes the key from state (StatePatch value = nil → runtime deletes key)
 * - append: adds resolved value to existing array
 * Outputs ["out"] with StatePatch containing the state delta.
 */

// SetStateExecutor modifies workflow instance state using resolved values.
// Supports operations: set, increment, decrement, append, remove.
type SetStateExecutor struct {
	resolver enginePorts.ValueResolverPort
}

// NewSetStateExecutor creates a new SetStateExecutor with the given value resolver.
func NewSetStateExecutor(resolver enginePorts.ValueResolverPort) entities.NodeExecutor {
	return &SetStateExecutor{resolver: resolver}
}

// NodeType returns "core/set_state".
func (e *SetStateExecutor) NodeType() string {
	return "core/set_state"
}

// Execute resolves the value source and applies the configured operation to the target field,
// returning a StatePatch with the state delta.
func (e *SetStateExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.SetStateNodeConfig)
	if !ok || cfg == nil {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"out"},
		}, nil
	}

	// remove operation: delete the key from state (no value needed)
	if cfg.Operation == "remove" {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"out"},
			StatePatch: map[string]interface{}{
				cfg.TargetField: nil, // nil sentinel → runtime deletes this key
			},
		}, nil
	}

	resolved, err := e.resolver.Resolve(
		&cfg.ValueSource,
		execCtx.EventPayload,
		execCtx.State,
		execCtx.NodeOutputs,
		execCtx.ExternalInputs,
	)
	if err != nil {
		return nil, fmt.Errorf("set_state: failed to resolve value: %w", err)
	}

	newValue, err := applyOperation(cfg.Operation, cfg.TargetField, resolved, execCtx.State)
	if err != nil {
		return nil, fmt.Errorf("set_state: operation %q failed: %w", cfg.Operation, err)
	}

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		StatePatch: map[string]interface{}{
			cfg.TargetField: newValue,
		},
	}, nil
}

func applyOperation(operation, targetField string, resolved interface{}, state map[string]interface{}) (interface{}, error) {
	switch operation {
	case "set", "":
		return resolved, nil

	case "increment":
		current := toFloat64(state[targetField])
		delta := toFloat64(resolved)
		return current + delta, nil

	case "decrement":
		current := toFloat64(state[targetField])
		delta := toFloat64(resolved)
		return current - delta, nil

	case "append":
		current := toSlice(state[targetField])
		return append(current, resolved), nil

	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

// toFloat64 converts a value to float64, returning 0 on failure.
func toFloat64(v interface{}) float64 {
	f, _ := typeconv.ToFloat64(v)
	return f
}

func toSlice(v interface{}) []interface{} {
	if v == nil {
		return []interface{}{}
	}
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return []interface{}{}
}

