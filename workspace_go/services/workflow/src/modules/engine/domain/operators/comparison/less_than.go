package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * LessThanOperator checks if fieldValue is less than compareValue.
 * Supports numeric, string (lexicographic), and time comparisons.
 *
 * Examples:
 * - 30 < 50 -> true
 * - "a" < "b" -> true (lexicographic)
 * - "2025-01-01" < "2025-01-15" -> true (date comparison)
 */
type LessThanOperator struct{}

// Ensure LessThanOperator implements the interface
var _ operators.ConditionOperator = (*LessThanOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *LessThanOperator) Name() string {
	return "lessThan"
}

/**
 * Evaluate checks if fieldValue is less than compareValue.
 *
 * @param timezone - IANA timezone (unused for this operator)
 * @param fieldValue - The value from the event/state/variable
 * @param compareValue - The value to compare against
 * @returns (true, nil) if fieldValue < compareValue, (false, nil) otherwise
 */
func (o *LessThanOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	result := Compare(fieldValue, compareValue)
	return result == CompareLess, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *LessThanOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "lessThan",
		Category:      "comparison",
		Description:   "Checks if the field value is less than the compare value",
		AcceptedTypes: []string{"number", "string", "date"},
		IsBetween:     false,
	}
}
