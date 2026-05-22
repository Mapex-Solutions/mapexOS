package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * NotEqualsOperator checks if two values are NOT equal.
 * Supports numeric, string, boolean, and time comparisons with type coercion.
 *
 * Examples:
 * - 42 != 43 -> true
 * - "hello" != "world" -> true
 * - 42 != "42" -> false (numeric coercion, they ARE equal)
 */
type NotEqualsOperator struct{}

// Ensure NotEqualsOperator implements the interface
var _ operators.ConditionOperator = (*NotEqualsOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *NotEqualsOperator) Name() string {
	return "notEquals"
}

/**
 * Evaluate checks if fieldValue does NOT equal compareValue.
 *
 * @param timezone - IANA timezone (unused for this operator)
 * @param fieldValue - The value from the event/state/variable
 * @param compareValue - The value to compare against
 * @returns (true, nil) if values are NOT equal, (false, nil) if they are equal
 */
func (o *NotEqualsOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	return !Equals(fieldValue, compareValue), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *NotEqualsOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "notEquals",
		Category:      "comparison",
		Description:   "Checks if two values are NOT equal (with type coercion)",
		AcceptedTypes: []string{"number", "string", "boolean", "date"},
		IsBetween:     false,
	}
}
