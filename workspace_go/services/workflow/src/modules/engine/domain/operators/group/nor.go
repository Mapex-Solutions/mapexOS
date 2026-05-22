package group

import (
	"workflow/src/modules/engine/domain/operators"
)

/*
 * NorOperator combines multiple condition results with logical NOR (NOT OR).
 * All conditions must be false for the group to match.
 * Supports short-circuit evaluation - returns false on first true.
 *
 * Truth table:
 * - [false, false, false] → true (none are true)
 * - [true, false, false] → false (at least one is true)
 * - [] (empty) → true (no true condition exists)
 *
 * Use case: "None of these error conditions are present"
 */
type NorOperator struct{}

// Ensure NorOperator implements the interface
var _ operators.GroupOperator = (*NorOperator)(nil)

/*
 * Name returns the operator identifier.
 */
func (o *NorOperator) Name() string {
	return "NOR"
}

/*
 * Evaluate combines results with logical NOR.
 *
 * @param results - Array of condition evaluation results
 * @returns true if NO result is true (all are false)
 */
func (o *NorOperator) Evaluate(results []bool) bool {
	// Empty array is true (no true condition exists)
	if len(results) == 0 {
		return true
	}

	for _, result := range results {
		if result {
			return false
		}
	}
	return true
}

/*
 * SupportsShortCircuit indicates that NOR supports short-circuit evaluation.
 */
func (o *NorOperator) SupportsShortCircuit() bool {
	return true
}

/*
 * ShouldShortCircuit returns true if evaluation should stop.
 * For NOR: stop on first true (result is already determined as false).
 */
func (o *NorOperator) ShouldShortCircuit(currentResult bool) bool {
	return currentResult // Short-circuit on true (means NOR is false)
}

/*
 * Metadata returns information about this operator.
 */
func (o *NorOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:        "NOR",
		Category:    "group",
		Description: "Logical NOR - none of the conditions can be true",
		IsBetween:   false,
	}
}
