package evaluators

import (
	"fmt"

	defPorts "workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/engine/domain/operators"
)

/*
 * CONDITION EVALUATOR
 * Evaluates ConditionGroup trees recursively with short-circuit logic.
 * Pure evaluation — no logging, no actions, no side effects.
 */

type ConditionEvaluator struct {
	registry      *operators.OperatorRegistry
	valueResolver *ValueResolver
}

/*
 * NewConditionEvaluator creates a new ConditionEvaluator with dependencies.
 */
func NewConditionEvaluator(registry *operators.OperatorRegistry) *ConditionEvaluator {
	return &ConditionEvaluator{
		registry:      registry,
		valueResolver: NewValueResolver(),
	}
}

/*
 * EvaluateGroup evaluates a ConditionGroup recursively with short-circuit logic.
 *
 * Group operators:
 * - AND: all items must match (short-circuits on first false)
 * - OR: at least one item must match (short-circuits on first true)
 * - NAND: NOT(AND) — at least one item must NOT match
 * - NOR: NOT(OR) — all items must NOT match
 */
func (e *ConditionEvaluator) EvaluateGroup(
	group *defPorts.ConditionGroup,
	timezone string,
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (bool, error) {
	if group == nil || len(group.Items) == 0 {
		return false, ErrEmptyGroup
	}

	switch group.Logic {
	case defPorts.LogicAND:
		return e.evaluateAND(group.Items, timezone, eventPayload, state, nodeOutputs, externalInputs)

	case defPorts.LogicOR:
		return e.evaluateOR(group.Items, timezone, eventPayload, state, nodeOutputs, externalInputs)

	case defPorts.LogicNAND:
		result, err := e.evaluateAND(group.Items, timezone, eventPayload, state, nodeOutputs, externalInputs)
		if err != nil {
			return false, err
		}
		return !result, nil

	case defPorts.LogicNOR:
		result, err := e.evaluateOR(group.Items, timezone, eventPayload, state, nodeOutputs, externalInputs)
		if err != nil {
			return false, err
		}
		return !result, nil

	default:
		return false, fmt.Errorf("unknown group logic operator: %s", group.Logic)
	}
}

/*
 * evaluateAND evaluates items with AND logic (short-circuits on first false).
 */
func (e *ConditionEvaluator) evaluateAND(
	items []defPorts.ConditionGroupItem,
	timezone string,
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (bool, error) {
	for i := range items {
		result, err := e.evaluateItem(&items[i], timezone, eventPayload, state, nodeOutputs, externalInputs)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

/*
 * evaluateOR evaluates items with OR logic (short-circuits on first true).
 */
func (e *ConditionEvaluator) evaluateOR(
	items []defPorts.ConditionGroupItem,
	timezone string,
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (bool, error) {
	for i := range items {
		result, err := e.evaluateItem(&items[i], timezone, eventPayload, state, nodeOutputs, externalInputs)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

/*
 * evaluateItem dispatches evaluation based on the item type.
 * Supports both typed structs and pointer variants.
 */
func (e *ConditionEvaluator) evaluateItem(
	item *defPorts.ConditionGroupItem,
	timezone string,
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (bool, error) {
	switch item.Type {
	case "condition":
		ci, err := asConditionItem(item.Data)
		if err != nil {
			return false, err
		}
		return e.evaluateCondition(ci, timezone, eventPayload, state, nodeOutputs, externalInputs)

	case "group":
		cg, err := asConditionGroup(item.Data)
		if err != nil {
			return false, err
		}
		return e.EvaluateGroup(cg, timezone, eventPayload, state, nodeOutputs, externalInputs)

	default:
		return false, fmt.Errorf("%w: unknown item type '%s'", ErrInvalidGroupItem, item.Type)
	}
}

/*
 * evaluateCondition evaluates a single condition by resolving values and applying an operator.
 * Field resolution failures are treated as non-match (not hard errors).
 */
func (e *ConditionEvaluator) evaluateCondition(
	condition *defPorts.ConditionItem,
	timezone string,
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (bool, error) {
	// Resolve field value
	fieldValue, err := e.valueResolver.Resolve(&condition.Field, eventPayload, state, nodeOutputs, externalInputs)
	if err != nil {
		return false, nil
	}

	// Check for between operators (dual interface with min/max)
	if isBetweenOperator(condition.Operator) {
		return e.evaluateBetween(condition, timezone, fieldValue, eventPayload, state, nodeOutputs, externalInputs)
	}

	// Resolve compare value
	compareValue, err := e.valueResolver.Resolve(&condition.Value, eventPayload, state, nodeOutputs, externalInputs)
	if err != nil {
		return false, nil
	}

	// Get standard operator from registry
	op, err := e.registry.GetConditionOperator(condition.Operator)
	if err != nil {
		return false, fmt.Errorf("%w: '%s'", ErrOperatorNotFound, condition.Operator)
	}

	return op.Evaluate(timezone, fieldValue, compareValue)
}

/*
 * evaluateBetween handles between operators that require min/max values.
 * The compare value resolves to a map with "min", "max", and optional "inclusive" keys.
 */
func (e *ConditionEvaluator) evaluateBetween(
	condition *defPorts.ConditionItem,
	timezone string,
	fieldValue interface{},
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (bool, error) {
	op, err := e.registry.GetBetweenOperator(condition.Operator)
	if err != nil {
		return false, fmt.Errorf("%w: '%s'", ErrOperatorNotFound, condition.Operator)
	}

	compareValue, err := e.valueResolver.Resolve(&condition.Value, eventPayload, state, nodeOutputs, externalInputs)
	if err != nil {
		return false, nil
	}

	minVal, maxVal, inclusive, err := extractRange(compareValue)
	if err != nil {
		return false, nil
	}

	return op.EvaluateRange(timezone, fieldValue, minVal, maxVal, inclusive)
}

/*
 * DATA TYPE ASSERTIONS
 * Handle typed structs and pointer variants from ConditionGroupItem.Data.
 */

func asConditionItem(data interface{}) (*defPorts.ConditionItem, error) {
	switch v := data.(type) {
	case *defPorts.ConditionItem:
		return v, nil
	case defPorts.ConditionItem:
		return &v, nil
	default:
		return nil, fmt.Errorf("%w: expected ConditionItem, got %T", ErrInvalidGroupItem, data)
	}
}

func asConditionGroup(data interface{}) (*defPorts.ConditionGroup, error) {
	switch v := data.(type) {
	case *defPorts.ConditionGroup:
		return v, nil
	case defPorts.ConditionGroup:
		return &v, nil
	default:
		return nil, fmt.Errorf("%w: expected ConditionGroup, got %T", ErrInvalidGroupItem, data)
	}
}

/*
 * HELPER FUNCTIONS
 */

/*
 * isBetweenOperator checks if the operator name is a between-type operator.
 */
func isBetweenOperator(name string) bool {
	switch name {
	case "between", "betweenDate", "betweenTime":
		return true
	default:
		return false
	}
}

/*
 * extractRange extracts min, max and inclusive from a range value.
 * Accepts:
 *   - map[string]interface{} with "min", "max", optional "inclusive" keys
 *   - []interface{} with [min, max] (assumes inclusive=true)
 */
func extractRange(value interface{}) (min, max interface{}, inclusive bool, err error) {
	switch v := value.(type) {
	case map[string]interface{}:
		min, minOk := v["min"]
		max, maxOk := v["max"]
		if !minOk || !maxOk {
			return nil, nil, false, fmt.Errorf("between operator requires both min and max")
		}
		// Default inclusive to true when the key is absent so callers can
		// omit the field for the common closed-interval case.
		inclusive = true
		if raw, exists := v["inclusive"]; exists {
			if b, ok := raw.(bool); ok {
				inclusive = b
			}
		}
		return min, max, inclusive, nil

	case []interface{}:
		if len(v) < 2 {
			return nil, nil, false, fmt.Errorf("between operator array requires at least 2 elements")
		}
		return v[0], v[1], true, nil

	default:
		return nil, nil, false, fmt.Errorf("between operator requires map with min/max keys or [min, max] array, got %T", value)
	}
}
