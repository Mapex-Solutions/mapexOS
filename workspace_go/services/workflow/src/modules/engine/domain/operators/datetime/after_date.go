package datetime

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * AfterDateOperator checks if a date is after another date.
 * Uses the timezone parameter for parsing.
 *
 * Examples:
 * - "2025-01-15" afterDate "2025-01-01" -> true
 * - "2025-01-01" afterDate "2025-01-15" -> false
 * - "2025-01-15" afterDate "2025-01-15" -> false (not after, equal)
 */
type AfterDateOperator struct{}

// Ensure AfterDateOperator implements the interface
var _ operators.ConditionOperator = (*AfterDateOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *AfterDateOperator) Name() string {
	return "afterDate"
}

/**
 * Evaluate checks if fieldValue date is after compareValue date.
 *
 * @param timezone - IANA timezone for date parsing
 * @param fieldValue - The date to check
 * @param compareValue - The reference date
 * @returns (true, nil) if fieldValue is after compareValue
 */
func (o *AfterDateOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	fieldDate, ok := ParseDate(fieldValue, timezone)
	if !ok {
		return false, nil
	}

	compareDate, ok := ParseDate(compareValue, timezone)
	if !ok {
		return false, nil
	}

	// Compare dates only (ignore time)
	fieldDateOnly := DateOnly(fieldDate)
	compareDateOnly := DateOnly(compareDate)

	return fieldDateOnly.After(compareDateOnly), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *AfterDateOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "afterDate",
		Category:      "datetime",
		Description:   "Checks if a date is after another date (date only, ignores time)",
		AcceptedTypes: []string{"date", "string"},
		IsBetween:     false,
	}
}
