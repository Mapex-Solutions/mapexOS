package stringops

import (
	"strings"

	"workflow/src/modules/engine/domain/operators"
)

/**
 * EndsWithOperator checks if a string ends with a suffix.
 * Case-insensitive by default.
 *
 * Examples:
 * - "Hello World" endsWith "World" -> true
 * - "Hello World" endsWith "world" -> true (case-insensitive)
 * - "Hello World" endsWith "Hello" -> false
 */
type EndsWithOperator struct{}

// Ensure EndsWithOperator implements the interface
var _ operators.ConditionOperator = (*EndsWithOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *EndsWithOperator) Name() string {
	return "endsWith"
}

/**
 * Evaluate checks if fieldValue ends with compareValue.
 *
 * @param timezone - IANA timezone (unused for string ops, kept for interface compliance)
 * @param fieldValue - The string to check
 * @param compareValue - The suffix to look for
 * @returns (true, nil) if fieldValue ends with compareValue
 */
func (o *EndsWithOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	// Convert both values to strings
	fieldStr, ok := toString(fieldValue)
	if !ok {
		return false, nil
	}

	compareStr, ok := toString(compareValue)
	if !ok {
		return false, nil
	}

	// Case-insensitive comparison
	return strings.HasSuffix(
		strings.ToLower(fieldStr),
		strings.ToLower(compareStr),
	), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *EndsWithOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "endsWith",
		Category:      "string",
		Description:   "Checks if a string ends with a suffix (case-insensitive)",
		AcceptedTypes: []string{"string"},
		IsBetween:     false,
	}
}
