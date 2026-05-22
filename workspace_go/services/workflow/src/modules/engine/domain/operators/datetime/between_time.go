package datetime

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * BetweenTimeOperator checks if a time is within a time range (time-of-day only).
 * Compares only the time portion, ignoring the date.
 * Uses the timezone parameter for parsing.
 * Boundaries are inclusive by default.
 *
 * Examples:
 * - "10:30" betweenTime ["08:00", "17:00"] -> true
 * - "08:00" betweenTime ["08:00", "17:00"] -> true (inclusive)
 * - "20:00" betweenTime ["08:00", "17:00"] -> false
 *
 * Note: Does NOT handle overnight ranges (e.g., "22:00" to "06:00").
 * For overnight ranges, use two conditions with OR.
 */
type BetweenTimeOperator struct{}

// Ensure BetweenTimeOperator implements the interfaces
var _ operators.BetweenOperator = (*BetweenTimeOperator)(nil)
var _ operators.ConditionOperator = (*BetweenTimeOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *BetweenTimeOperator) Name() string {
	return "betweenTime"
}

/**
 * Evaluate checks if fieldValue time is between compareValue times.
 * compareValue should be an array with [startTime, endTime].
 *
 * @param timezone - IANA timezone for time parsing
 * @param fieldValue - The time to check
 * @param compareValue - Array with [startTime, endTime]
 * @returns (true, nil) if fieldValue time is within the range
 */
func (o *BetweenTimeOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	// Extract start and end times from compareValue
	startTime, endTime, ok := extractTimeRange(compareValue)
	if !ok {
		return false, nil
	}

	return o.EvaluateRange(timezone, fieldValue, startTime, endTime, true)
}

/**
 * EvaluateRange checks if fieldValue time is within the time range.
 *
 * @param timezone - IANA timezone for time parsing
 * @param fieldValue - The time to check
 * @param minValue - Start time
 * @param maxValue - End time
 * @param inclusive - If true, boundaries are included
 * @returns (true, nil) if within range
 */
func (o *BetweenTimeOperator) EvaluateRange(
	timezone string,
	fieldValue interface{},
	minValue interface{},
	maxValue interface{},
	inclusive bool,
) (bool, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	// Parse all times
	fieldTime, ok := ParseTime(fieldValue, timezone)
	if !ok {
		return false, nil
	}

	startTime, ok := ParseTime(minValue, timezone)
	if !ok {
		return false, nil
	}

	endTime, ok := ParseTime(maxValue, timezone)
	if !ok {
		return false, nil
	}

	// Compare time portions
	fieldVsStart := CompareTimes(fieldTime, startTime)
	fieldVsEnd := CompareTimes(fieldTime, endTime)

	if inclusive {
		// field >= start && field <= end
		return fieldVsStart >= 0 && fieldVsEnd <= 0, nil
	}

	// Exclusive: field > start && field < end
	return fieldVsStart > 0 && fieldVsEnd < 0, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *BetweenTimeOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "betweenTime",
		Category:      "datetime",
		Description:   "Checks if a time-of-day is within a time range [start, end]",
		AcceptedTypes: []string{"time", "string"},
		IsBetween:     true,
	}
}

/**
 * extractTimeRange extracts start and end times from various formats.
 */
func extractTimeRange(value interface{}) (start, end interface{}, ok bool) {
	switch v := value.(type) {
	case []interface{}:
		if len(v) >= 2 {
			return v[0], v[1], true
		}
	case [2]interface{}:
		return v[0], v[1], true
	case []string:
		if len(v) >= 2 {
			return v[0], v[1], true
		}
	case map[string]interface{}:
		// Support {"start": x, "end": y} or {"min": x, "max": y} format
		startVal, hasStart := v["start"]
		endVal, hasEnd := v["end"]
		if hasStart && hasEnd {
			return startVal, endVal, true
		}
		minVal, hasMin := v["min"]
		maxVal, hasMax := v["max"]
		if hasMin && hasMax {
			return minVal, maxVal, true
		}
	}
	return nil, nil, false
}
