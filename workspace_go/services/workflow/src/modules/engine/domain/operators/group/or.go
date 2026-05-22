package group

import (
	"workflow/src/modules/engine/domain/operators"
)

/*
 * OrOperator combines multiple condition results with logical OR.
 * At least one condition must be true for the group to match.
 * Supports short-circuit evaluation - returns true on first true.
 *
 * Truth table:
 * - [true, false, false] → true
 * - [false, false, false] → false
 * - [] (empty) → false
 *
 * Use case: "Temperature is high OR humidity is critical OR pressure is dangerous"
 */
type OrOperator struct{}

// Ensure OrOperator implements the interface
var _ operators.GroupOperator = (*OrOperator)(nil)

/*
 * Name returns the operator identifier.
 */
func (o *OrOperator) Name() string {
	return "OR"
}

/*
 * Evaluate combines results with logical OR.
 *
 * @param results - Array of condition evaluation results
 * @returns true if ANY result is true
 */
func (o *OrOperator) Evaluate(results []bool) bool {
	// Empty array is false (no true condition found)
	if len(results) == 0 {
		return false
	}

	for _, result := range results {
		if result {
			return true
		}
	}
	return false
}

/*
 * SupportsShortCircuit indicates that OR supports short-circuit evaluation.
 */
func (o *OrOperator) SupportsShortCircuit() bool {
	return true
}

/*
 * ShouldShortCircuit returns true if evaluation should stop.
 * For OR: stop on first true (result is already determined).
 */
func (o *OrOperator) ShouldShortCircuit(currentResult bool) bool {
	return currentResult // Short-circuit on true
}

/*
 * Metadata returns information about this operator.
 */
func (o *OrOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:        "OR",
		Category:    "group",
		Description: "Logical OR - at least one condition must be true",
		IsBetween:   false,
	}
}
