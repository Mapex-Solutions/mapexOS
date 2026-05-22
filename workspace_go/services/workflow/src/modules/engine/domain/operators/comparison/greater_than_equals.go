package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * GreaterThanEqualsOperator checks if fieldValue is greater than or equal to compareValue.
 * Supports numeric, string (lexicographic), and time comparisons.
 *
 * Examples:
 * - 50 >= 50 -> true
 * - 50 >= 30 -> true
 * - "b" >= "a" -> true (lexicographic)
 */
type GreaterThanEqualsOperator struct{}

// Ensure GreaterThanEqualsOperator implements the interface
var _ operators.ConditionOperator = (*GreaterThanEqualsOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *GreaterThanEqualsOperator) Name() string {
	return "greaterThanEquals"
}

/**
 * Evaluate checks if fieldValue is greater than or equal to compareValue.
 *
 * @param timezone - IANA timezone (unused for this operator)
 * @param fieldValue - The value from the event/state/variable
 * @param compareValue - The value to compare against
 * @returns (true, nil) if fieldValue >= compareValue, (false, nil) otherwise
 */
func (o *GreaterThanEqualsOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	result := Compare(fieldValue, compareValue)
	return result == CompareGreater || result == CompareEqual, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *GreaterThanEqualsOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "greaterThanEquals",
		Category:      "comparison",
		Description:   "Checks if the field value is greater than or equal to the compare value",
		AcceptedTypes: []string{"number", "string", "date"},
		IsBetween:     false,
	}
}
