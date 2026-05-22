package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * EqualsOperator checks if two values are equal.
 * Supports numeric, string, boolean, and time comparisons with type coercion.
 *
 * Examples:
 * - 42 == 42 -> true
 * - "hello" == "hello" -> true
 * - 42 == "42" -> true (numeric coercion)
 * - true == "true" -> true (boolean coercion)
 */
type EqualsOperator struct{}

// Ensure EqualsOperator implements the interface
var _ operators.ConditionOperator = (*EqualsOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *EqualsOperator) Name() string {
	return "equals"
}

/**
 * Evaluate checks if fieldValue equals compareValue.
 *
 * @param timezone - IANA timezone (unused for this operator)
 * @param fieldValue - The value from the event/state/variable
 * @param compareValue - The value to compare against
 * @returns (true, nil) if values are equal, (false, nil) otherwise
 */
func (o *EqualsOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	return Equals(fieldValue, compareValue), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *EqualsOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "equals",
		Category:      "comparison",
		Description:   "Checks if two values are equal (with type coercion)",
		AcceptedTypes: []string{"number", "string", "boolean", "date"},
		IsBetween:     false,
	}
}
