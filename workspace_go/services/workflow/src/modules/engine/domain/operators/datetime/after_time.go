package datetime

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * AfterTimeOperator checks if a time is after another time (time-of-day only).
 * Compares only the time portion, ignoring the date.
 * Uses the timezone parameter for parsing.
 *
 * Examples:
 * - "14:30" afterTime "12:00" -> true
 * - "08:00" afterTime "12:00" -> false
 * - "12:00" afterTime "12:00" -> false (not after, equal)
 */
type AfterTimeOperator struct{}

// Ensure AfterTimeOperator implements the interface
var _ operators.ConditionOperator = (*AfterTimeOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *AfterTimeOperator) Name() string {
	return "afterTime"
}

/**
 * Evaluate checks if fieldValue time is after compareValue time.
 * Only compares the time portion (hours, minutes, seconds).
 *
 * @param timezone - IANA timezone for time parsing
 * @param fieldValue - The time to check
 * @param compareValue - The reference time
 * @returns (true, nil) if fieldValue time is after compareValue time
 */
func (o *AfterTimeOperator) Evaluate(
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
	return CompareTimes(fieldTime, compareTime) > 0, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *AfterTimeOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "afterTime",
		Category:      "datetime",
		Description:   "Checks if a time-of-day is after another time (ignores date)",
		AcceptedTypes: []string{"time", "string"},
		IsBetween:     false,
	}
}
