package stringops

import (
	"strings"

	"workflow/src/modules/engine/domain/operators"
)

/**
 * StartsWithOperator checks if a string starts with a prefix.
 * Case-insensitive by default.
 *
 * Examples:
 * - "Hello World" startsWith "Hello" -> true
 * - "Hello World" startsWith "hello" -> true (case-insensitive)
 * - "Hello World" startsWith "World" -> false
 */
type StartsWithOperator struct{}

// Ensure StartsWithOperator implements the interface
var _ operators.ConditionOperator = (*StartsWithOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *StartsWithOperator) Name() string {
	return "startsWith"
}

/**
 * Evaluate checks if fieldValue starts with compareValue.
 *
 * @param timezone - IANA timezone (unused for string ops, kept for interface compliance)
 * @param fieldValue - The string to check
 * @param compareValue - The prefix to look for
 * @returns (true, nil) if fieldValue starts with compareValue
 */
func (o *StartsWithOperator) Evaluate(
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
	return strings.HasPrefix(
		strings.ToLower(fieldStr),
		strings.ToLower(compareStr),
	), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *StartsWithOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "startsWith",
		Category:      "string",
		Description:   "Checks if a string starts with a prefix (case-insensitive)",
		AcceptedTypes: []string{"string"},
		IsBetween:     false,
	}
}
