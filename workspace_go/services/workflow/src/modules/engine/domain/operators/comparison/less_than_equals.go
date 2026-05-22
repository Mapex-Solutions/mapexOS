package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * LessThanEqualsOperator checks if fieldValue is less than or equal to compareValue.
 * Supports numeric, string (lexicographic), and time comparisons.
 *
 * Examples:
 * - 30 <= 30 -> true
 * - 30 <= 50 -> true
 * - "a" <= "b" -> true (lexicographic)
 */
type LessThanEqualsOperator struct{}

// Ensure LessThanEqualsOperator implements the interface
var _ operators.ConditionOperator = (*LessThanEqualsOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *LessThanEqualsOperator) Name() string {
	return "lessThanEquals"
}

/**
 * Evaluate checks if fieldValue is less than or equal to compareValue.
 *
 * @param timezone - IANA timezone (unused for this operator)
 * @param fieldValue - The value from the event/state/variable
 * @param compareValue - The value to compare against
 * @returns (true, nil) if fieldValue <= compareValue, (false, nil) otherwise
 */
func (o *LessThanEqualsOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	result := Compare(fieldValue, compareValue)
	return result == CompareLess || result == CompareEqual, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *LessThanEqualsOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "lessThanEquals",
		Category:      "comparison",
		Description:   "Checks if the field value is less than or equal to the compare value",
		AcceptedTypes: []string{"number", "string", "date"},
		IsBetween:     false,
	}
}
