package datetime

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * BetweenDateOperator checks if a date is within a date range.
 * Uses the timezone parameter for parsing.
 * Boundaries are inclusive by default.
 *
 * Examples:
 * - "2025-01-10" betweenDate ["2025-01-01", "2025-01-31"] -> true
 * - "2025-01-01" betweenDate ["2025-01-01", "2025-01-31"] -> true (inclusive)
 * - "2025-02-15" betweenDate ["2025-01-01", "2025-01-31"] -> false
 */
type BetweenDateOperator struct{}

// Ensure BetweenDateOperator implements the interfaces
var _ operators.BetweenOperator = (*BetweenDateOperator)(nil)
var _ operators.ConditionOperator = (*BetweenDateOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *BetweenDateOperator) Name() string {
	return "betweenDate"
}

/**
 * Evaluate checks if fieldValue is between compareValue dates.
 * compareValue should be an array with [startDate, endDate].
 *
 * @param timezone - IANA timezone for date parsing
 * @param fieldValue - The date to check
 * @param compareValue - Array with [startDate, endDate]
 * @returns (true, nil) if fieldValue is within the range
 */
func (o *BetweenDateOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	// Extract start and end dates from compareValue
	startDate, endDate, ok := extractDateRange(compareValue)
	if !ok {
		return false, nil
	}

	return o.EvaluateRange(timezone, fieldValue, startDate, endDate, true)
}

/**
 * EvaluateRange checks if fieldValue is within the date range.
 *
 * @param timezone - IANA timezone for date parsing
 * @param fieldValue - The date to check
 * @param minValue - Start date
 * @param maxValue - End date
 * @param inclusive - If true, boundaries are included
 * @returns (true, nil) if within range
 */
func (o *BetweenDateOperator) EvaluateRange(
	timezone string,
	fieldValue interface{},
	minValue interface{},
	maxValue interface{},
	inclusive bool,
) (bool, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	// Parse all dates
	fieldDate, ok := ParseDate(fieldValue, timezone)
	if !ok {
		return false, nil
	}

	startDate, ok := ParseDate(minValue, timezone)
	if !ok {
		return false, nil
	}

	endDate, ok := ParseDate(maxValue, timezone)
	if !ok {
		return false, nil
	}

	// Compare dates only (ignore time)
	fieldDateOnly := DateOnly(fieldDate)
	startDateOnly := DateOnly(startDate)
	endDateOnly := DateOnly(endDate)

	if inclusive {
		// field >= start && field <= end
		return !fieldDateOnly.Before(startDateOnly) && !fieldDateOnly.After(endDateOnly), nil
	}

	// Exclusive: field > start && field < end
	return fieldDateOnly.After(startDateOnly) && fieldDateOnly.Before(endDateOnly), nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *BetweenDateOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "betweenDate",
		Category:      "datetime",
		Description:   "Checks if a date is within a date range [start, end]",
		AcceptedTypes: []string{"date", "string"},
		IsBetween:     true,
	}
}

/**
 * extractDateRange extracts start and end dates from various formats.
 */
func extractDateRange(value interface{}) (start, end interface{}, ok bool) {
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
