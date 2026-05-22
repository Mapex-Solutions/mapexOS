package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * GreaterThanOperator checks if fieldValue is greater than compareValue.
 * Supports numeric, string (lexicographic), and time comparisons.
 *
 * Examples:
 * - 50 > 30 -> true
 * - "b" > "a" -> true (lexicographic)
 * - "2025-01-15" > "2025-01-01" -> true (date comparison)
 */
type GreaterThanOperator struct{}

// Ensure GreaterThanOperator implements the interface
var _ operators.ConditionOperator = (*GreaterThanOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *GreaterThanOperator) Name() string {
	return "greaterThan"
}

/**
 * Evaluate checks if fieldValue is greater than compareValue.
 *
 * @param timezone - IANA timezone (unused for this operator)
 * @param fieldValue - The value from the event/state/variable
 * @param compareValue - The value to compare against
 * @returns (true, nil) if fieldValue > compareValue, (false, nil) otherwise
 */
func (o *GreaterThanOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	result := Compare(fieldValue, compareValue)
	return result == CompareGreater, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *GreaterThanOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "greaterThan",
		Category:      "comparison",
		Description:   "Checks if the field value is greater than the compare value",
		AcceptedTypes: []string{"number", "string", "date"},
		IsBetween:     false,
	}
}
