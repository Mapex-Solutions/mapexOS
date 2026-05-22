package operators

import (
	"fmt"
)

/**
 * OperatorRegistry manages all available operators for the rule engine.
 * Operators are registered ONCE at startup and never modified at runtime.
 *
 * IMPORTANT: This registry is designed to be READ-ONLY after initialization.
 * All operators must be registered via RegisterAllDefaults() during startup,
 * before any rule evaluation begins. This allows lock-free concurrent reads.
 *
 * Usage:
 *   registry := NewOperatorRegistry()
 *   registry.RegisterAllDefaults()  // Called ONCE at startup
 *
 *   // After initialization, only reads happen (no locks needed)
 *   op, err := registry.GetConditionOperator("equals")
 *   result, err := op.Evaluate(timezone, fieldValue, compareValue)
 */
type OperatorRegistry struct {
	// Condition operators (equals, greaterThan, contains, etc.)
	conditionOperators map[string]ConditionOperator

	// Group operators (and, or, nor, nand)
	groupOperators map[string]GroupOperator

	// Between operators (between, betweenDate, betweenTime)
	betweenOperators map[string]BetweenOperator
}

/**
 * NewOperatorRegistry creates a new empty operator registry.
 * Call RegisterAllDefaults() to populate with built-in operators.
 */
func NewOperatorRegistry() *OperatorRegistry {
	return &OperatorRegistry{
		conditionOperators: make(map[string]ConditionOperator),
		groupOperators:     make(map[string]GroupOperator),
		betweenOperators:   make(map[string]BetweenOperator),
	}
}

/* Condition Operators */

/**
 * RegisterCondition adds a condition operator to the registry.
 * MUST be called during initialization only, before any rule evaluation.
 */
func (r *OperatorRegistry) RegisterCondition(op ConditionOperator) {
	r.conditionOperators[op.Name()] = op
}

/**
 * GetConditionOperator retrieves a condition operator by name.
 * Returns an error if the operator is not found.
 * Safe for concurrent reads (no writes happen after initialization).
 */
func (r *OperatorRegistry) GetConditionOperator(name string) (ConditionOperator, error) {
	op, exists := r.conditionOperators[name]
	if !exists {
		return nil, fmt.Errorf("condition operator '%s' not found", name)
	}
	return op, nil
}

/**
 * HasConditionOperator checks if a condition operator exists.
 */
func (r *OperatorRegistry) HasConditionOperator(name string) bool {
	_, exists := r.conditionOperators[name]
	return exists
}

/**
 * ListConditionOperators returns all registered condition operator names.
 */
func (r *OperatorRegistry) ListConditionOperators() []string {
	names := make([]string, 0, len(r.conditionOperators))
	for name := range r.conditionOperators {
		names = append(names, name)
	}
	return names
}

/* Group Operators */

/**
 * RegisterGroup adds a group operator to the registry.
 * MUST be called during initialization only, before any rule evaluation.
 */
func (r *OperatorRegistry) RegisterGroup(op GroupOperator) {
	r.groupOperators[op.Name()] = op
}

/**
 * GetGroupOperator retrieves a group operator by name.
 * Returns an error if the operator is not found.
 * Safe for concurrent reads (no writes happen after initialization).
 */
func (r *OperatorRegistry) GetGroupOperator(name string) (GroupOperator, error) {
	op, exists := r.groupOperators[name]
	if !exists {
		return nil, fmt.Errorf("group operator '%s' not found", name)
	}
	return op, nil
}

/**
 * HasGroupOperator checks if a group operator exists.
 */
func (r *OperatorRegistry) HasGroupOperator(name string) bool {
	_, exists := r.groupOperators[name]
	return exists
}

/**
 * ListGroupOperators returns all registered group operator names.
 */
func (r *OperatorRegistry) ListGroupOperators() []string {
	names := make([]string, 0, len(r.groupOperators))
	for name := range r.groupOperators {
		names = append(names, name)
	}
	return names
}

/* Between Operators */

/**
 * RegisterBetween adds a between operator to the registry.
 * MUST be called during initialization only, before any rule evaluation.
 */
func (r *OperatorRegistry) RegisterBetween(op BetweenOperator) {
	r.betweenOperators[op.Name()] = op
}

/**
 * GetBetweenOperator retrieves a between operator by name.
 * Returns an error if the operator is not found.
 * Safe for concurrent reads (no writes happen after initialization).
 */
func (r *OperatorRegistry) GetBetweenOperator(name string) (BetweenOperator, error) {
	op, exists := r.betweenOperators[name]
	if !exists {
		return nil, fmt.Errorf("between operator '%s' not found", name)
	}
	return op, nil
}

/**
 * HasBetweenOperator checks if a between operator exists.
 */
func (r *OperatorRegistry) HasBetweenOperator(name string) bool {
	_, exists := r.betweenOperators[name]
	return exists
}

/**
 * ListBetweenOperators returns all registered between operator names.
 */
func (r *OperatorRegistry) ListBetweenOperators() []string {
	names := make([]string, 0, len(r.betweenOperators))
	for name := range r.betweenOperators {
		names = append(names, name)
	}
	return names
}

/**
 * RegisterAllDefaults registers all built-in operators to the registry.
 * This should be called during application initialization.
 *
 * Categories:
 * - Comparison: equals, notEquals, greaterThan, greaterThanEquals, lessThan, lessThanEquals
 * - String: contains, notContains, startsWith, endsWith, regex
 * - DateTime: beforeDate, afterDate, beforeTime, afterTime
 * - Between: between, betweenDate, betweenTime
 * - Group: and, or, nor, nand
 */
func (r *OperatorRegistry) RegisterAllDefaults() {
	// Comparison operators - will be registered when implemented
	// r.RegisterCondition(&comparison.EqualsOperator{})
	// r.RegisterCondition(&comparison.NotEqualsOperator{})
	// r.RegisterCondition(&comparison.GreaterThanOperator{})
	// r.RegisterCondition(&comparison.GreaterThanEqualsOperator{})
	// r.RegisterCondition(&comparison.LessThanOperator{})
	// r.RegisterCondition(&comparison.LessThanEqualsOperator{})

	// String operators - will be registered when implemented
	// r.RegisterCondition(&stringops.ContainsOperator{})
	// r.RegisterCondition(&stringops.NotContainsOperator{})
	// r.RegisterCondition(&stringops.StartsWithOperator{})
	// r.RegisterCondition(&stringops.EndsWithOperator{})
	// r.RegisterCondition(&stringops.RegexOperator{})

	// DateTime operators - will be registered when implemented
	// r.RegisterCondition(&datetime.BeforeDateOperator{})
	// r.RegisterCondition(&datetime.AfterDateOperator{})
	// r.RegisterCondition(&datetime.BeforeTimeOperator{})
	// r.RegisterCondition(&datetime.AfterTimeOperator{})

	// Between operators - will be registered when implemented
	// r.RegisterBetween(&comparison.BetweenOperator{})
	// r.RegisterBetween(&datetime.BetweenDateOperator{})
	// r.RegisterBetween(&datetime.BetweenTimeOperator{})

	// Group operators - will be registered when implemented
	// r.RegisterGroup(&group.AndOperator{})
	// r.RegisterGroup(&group.OrOperator{})
	// r.RegisterGroup(&group.NorOperator{})
	// r.RegisterGroup(&group.NandOperator{})
}
