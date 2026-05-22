package datetime

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * BeforeTimeOperator checks if a time is before another time (time-of-day only).
 * Compares only the time portion, ignoring the date.
 * Uses the timezone parameter for parsing.
 *
 * Examples:
 * - "08:00" beforeTime "12:00" -> true
 * - "14:30" beforeTime "12:00" -> false
 * - "12:00" beforeTime "12:00" -> false (not before, equal)
 */
type BeforeTimeOperator struct{}

// Ensure BeforeTimeOperator implements the interface
var _ operators.ConditionOperator = (*BeforeTimeOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *BeforeTimeOperator) Name() string {
	return "beforeTime"
}

/**
 * Evaluate checks if fieldValue time is before compareValue time.
 * Only compares the time portion (hours, minutes, seconds).
 *
 * @param timezone - IANA timezone for time parsing
 * @param fieldValue - The time to check
 * @param compareValue - The reference time
 * @returns (true, nil) if fieldValue time is before compareValue time
 */
func (o *BeforeTimeOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	// Parse both times
	fieldTime, ok := ParseTime(fieldValue, timezone)
	if !ok {
		return false, nil
	}

	compareTime, ok := ParseTime(compareValue, timezone)
	if !ok {
		return false, nil
	}

	// Compare only time portions
	return CompareTimes(fieldTime, compareTime) < 0, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *BeforeTimeOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "beforeTime",
		Category:      "datetime",
		Description:   "Checks if a time-of-day is before another time (ignores date)",
		AcceptedTypes: []string{"time", "string"},
		IsBetween:     false,
	}
}
