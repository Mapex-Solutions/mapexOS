package group

import (
	"workflow/src/modules/engine/domain/operators"
)

/*
 * NandOperator combines multiple condition results with logical NAND (NOT AND).
 * Not all conditions can be true for the group to match.
 * Supports short-circuit evaluation - returns true on first false.
 *
 * Truth table:
 * - [true, true, true] → false (all are true, so NAND is false)
 * - [true, false, true] → true (not all are true)
 * - [false, false, false] → true (not all are true)
 * - [] (empty) → false (vacuously, NAND of nothing)
 *
 * Use case: "Not all safety conditions are met (trigger alert)"
 */
type NandOperator struct{}

// Ensure NandOperator implements the interface
var _ operators.GroupOperator = (*NandOperator)(nil)

/*
 * Name returns the operator identifier.
 */
func (o *NandOperator) Name() string {
	return "NAND"
}

/*
 * Evaluate combines results with logical NAND.
 *
 * @param results - Array of condition evaluation results
 * @returns true if NOT ALL results are true
 */
func (o *NandOperator) Evaluate(results []bool) bool {
	// Empty array - NAND of empty is false (opposite of AND's vacuous true)
	if len(results) == 0 {
		return false
	}

	// If any result is false, NAND is true
	for _, result := range results {
		if !result {
			return true
		}
	}
	// All are true, so NAND is false
	return false
}

/*
 * SupportsShortCircuit indicates that NAND supports short-circuit evaluation.
 */
func (o *NandOperator) SupportsShortCircuit() bool {
	return true
}

/*
 * ShouldShortCircuit returns true if evaluation should stop.
 * For NAND: stop on first false (result is already determined as true).
 */
func (o *NandOperator) ShouldShortCircuit(currentResult bool) bool {
	return !currentResult // Short-circuit on false (means NAND is true)
}

/*
 * Metadata returns information about this operator.
 */
func (o *NandOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:        "NAND",
		Category:    "group",
		Description: "Logical NAND - not all conditions can be true",
		IsBetween:   false,
	}
}
