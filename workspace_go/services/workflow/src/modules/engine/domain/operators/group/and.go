package group

import (
	"workflow/src/modules/engine/domain/operators"
)

/*
 * AndOperator combines multiple condition results with logical AND.
 * All conditions must be true for the group to match.
 * Supports short-circuit evaluation - returns false on first false.
 *
 * Truth table:
 * - [true, true, true] → true
 * - [true, false, true] → false
 * - [] (empty) → true (vacuous truth)
 *
 * Use case: "Temperature is high AND humidity is low AND device is online"
 */
type AndOperator struct{}

// Ensure AndOperator implements the interface
var _ operators.GroupOperator = (*AndOperator)(nil)

/*
 * Name returns the operator identifier.
 */
func (o *AndOperator) Name() string {
	return "AND"
}

/*
 * Evaluate combines results with logical AND.
 *
 * @param results - Array of condition evaluation results
 * @returns true if ALL results are true
 */
func (o *AndOperator) Evaluate(results []bool) bool {
	// Empty array is vacuously true
	if len(results) == 0 {
		return true
	}

	for _, result := range results {
		if !result {
			return false
		}
	}
	return true
}

/*
 * SupportsShortCircuit indicates that AND supports short-circuit evaluation.
 */
func (o *AndOperator) SupportsShortCircuit() bool {
	return true
}

/*
 * ShouldShortCircuit returns true if evaluation should stop.
 * For AND: stop on first false (result is already determined).
 */
func (o *AndOperator) ShouldShortCircuit(currentResult bool) bool {
	return !currentResult // Short-circuit on false
}

/*
 * Metadata returns information about this operator.
 */
func (o *AndOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:        "AND",
		Category:    "group",
		Description: "Logical AND - all conditions must be true",
		IsBetween:   false,
	}
}
