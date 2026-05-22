package datetime

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * BeforeDateOperator checks if a date is before another date.
 * Uses the timezone parameter for parsing.
 *
 * Examples:
 * - "2025-01-01" beforeDate "2025-01-15" -> true
 * - "2025-01-15" beforeDate "2025-01-01" -> false
 * - "2025-01-15" beforeDate "2025-01-15" -> false (not before, equal)
 */
type BeforeDateOperator struct{}

// Ensure BeforeDateOperator implements the interface
var _ operators.ConditionOperator = (*BeforeDateOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *BeforeDateOperator) Name() string {
	return "beforeDate"
}

/**
 * Evaluate checks if fieldValue date is before compareValue date.
 *
 * @param timezone - IANA timezone for date parsing
 * @param fieldValue - The date to check
 * @param compareValue - The reference date
 * @returns (true, nil) if fieldValue is before compareValue
 */
func (o *BeforeDateOperator) Evaluate(
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

	return fieldDateOnly.Before(compareDateOnly), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *BeforeDateOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "beforeDate",
		Category:      "datetime",
		Description:   "Checks if a date is before another date (date only, ignores time)",
		AcceptedTypes: []string{"date", "string"},
		IsBetween:     false,
	}
}
