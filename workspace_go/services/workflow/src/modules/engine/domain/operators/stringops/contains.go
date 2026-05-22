package stringops

import (
	"strings"

	"workflow/src/modules/engine/domain/operators"
)

/**
 * ContainsOperator checks if a string contains a substring.
 * Case-insensitive by default.
 *
 * Examples:
 * - "Hello World" contains "World" -> true
 * - "Hello World" contains "world" -> true (case-insensitive)
 * - "Hello World" contains "xyz" -> false
 */
type ContainsOperator struct{}

// Ensure ContainsOperator implements the interface
var _ operators.ConditionOperator = (*ContainsOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *ContainsOperator) Name() string {
	return "contains"
}

/**
 * Evaluate checks if fieldValue contains compareValue as substring.
 *
 * @param timezone - IANA timezone (unused for string ops, kept for interface compliance)
 * @param fieldValue - The string to search in
 * @param compareValue - The substring to search for
 * @returns (true, nil) if fieldValue contains compareValue
 */
func (o *ContainsOperator) Evaluate(
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
	return strings.Contains(
		strings.ToLower(fieldStr),
		strings.ToLower(compareStr),
	), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *ContainsOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "contains",
		Category:      "string",
		Description:   "Checks if a string contains a substring (case-insensitive)",
		AcceptedTypes: []string{"string"},
		IsBetween:     false,
	}
}
