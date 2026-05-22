package stringops

import (
	"strings"

	"workflow/src/modules/engine/domain/operators"
)

/**
 * NotContainsOperator checks if a string does NOT contain a substring.
 * Case-insensitive by default.
 *
 * Examples:
 * - "Hello World" notContains "xyz" -> true
 * - "Hello World" notContains "World" -> false
 * - "Hello World" notContains "world" -> false (case-insensitive)
 */
type NotContainsOperator struct{}

// Ensure NotContainsOperator implements the interface
var _ operators.ConditionOperator = (*NotContainsOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *NotContainsOperator) Name() string {
	return "notContains"
}

/**
 * Evaluate checks if fieldValue does NOT contain compareValue as substring.
 *
 * @param timezone - IANA timezone (unused for string ops, kept for interface compliance)
 * @param fieldValue - The string to search in
 * @param compareValue - The substring to search for
 * @returns (true, nil) if fieldValue does NOT contain compareValue
 */
func (o *NotContainsOperator) Evaluate(
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

	// Case-insensitive comparison - returns true if NOT found
	return !strings.Contains(
		strings.ToLower(fieldStr),
		strings.ToLower(compareStr),
	), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *NotContainsOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "notContains",
		Category:      "string",
		Description:   "Checks if a string does NOT contain a substring (case-insensitive)",
		AcceptedTypes: []string{"string"},
		IsBetween:     false,
	}
}
