package comparison

import (
	"workflow/src/modules/engine/domain/operators"
)

/**
 * BetweenOperator checks if a value is within a range.
 * Supports numeric, string (lexicographic), and time comparisons.
 *
 * Examples:
 * - 45 between [30, 50] (inclusive) -> true
 * - 30 between [30, 50] (exclusive) -> false
 * - "b" between ["a", "c"] -> true
 */
type BetweenOperator struct{}

// Ensure BetweenOperator implements the interfaces
var _ operators.BetweenOperator = (*BetweenOperator)(nil)
var _ operators.ConditionOperator = (*BetweenOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *BetweenOperator) Name() string {
	return "between"
}

/**
 * Evaluate checks if fieldValue is between compareValue (used as a two-element array).
 * This implements the ConditionOperator interface for simple usage.
 *
 * compareValue should be a slice/array with [min, max] values.
 * Assumes inclusive boundaries.
 *
 * @param timezone - IANA timezone for date/time operations
 * @param fieldValue - The value to check
 * @param compareValue - Array/slice with [min, max]
 * @returns (true, nil) if value is within range
 */
func (o *BetweenOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	// Extract min and max from compareValue
	minValue, maxValue, ok := extractMinMax(compareValue)
	if !ok {
		return false, nil
	}

	return o.EvaluateRange(timezone, fieldValue, minValue, maxValue, true)
}

/**
 * EvaluateRange checks if fieldValue is within the range [minValue, maxValue].
 * This implements the BetweenOperator interface for explicit range specification.
 *
 * @param timezone - IANA timezone for date/time operations
 * @param fieldValue - The value to check
 * @param minValue - Minimum boundary
 * @param maxValue - Maximum boundary
 * @param inclusive - If true, boundaries are included in the range
 * @returns (true, nil) if value is within range
 */
func (o *BetweenOperator) EvaluateRange(
	timezone string,
	fieldValue interface{},
	minValue interface{},
	maxValue interface{},
	inclusive bool,
) (bool, error) {
	// Compare with minimum
	minResult := Compare(fieldValue, minValue)
	if minResult == CompareError {
		return false, nil
	}

	// Compare with maximum
	maxResult := Compare(fieldValue, maxValue)
	if maxResult == CompareError {
		return false, nil
	}

	if inclusive {
		// value >= min && value <= max
		return (minResult == CompareGreater || minResult == CompareEqual) &&
			(maxResult == CompareLess || maxResult == CompareEqual), nil
	}

	// Exclusive: value > min && value < max
	return minResult == CompareGreater && maxResult == CompareLess, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *BetweenOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "between",
		Category:      "comparison",
		Description:   "Checks if the field value is within a range [min, max]",
		AcceptedTypes: []string{"number", "string", "date"},
		IsBetween:     true,
	}
}

/**
 * extractMinMax extracts min and max values from a slice or array.
 */
func extractMinMax(value interface{}) (min, max interface{}, ok bool) {
	switch v := value.(type) {
	case []interface{}:
		if len(v) >= 2 {
			return v[0], v[1], true
		}
	case [2]interface{}:
		return v[0], v[1], true
	case []float64:
		if len(v) >= 2 {
			return v[0], v[1], true
		}
	case []int:
		if len(v) >= 2 {
			return v[0], v[1], true
		}
	case []string:
		if len(v) >= 2 {
			return v[0], v[1], true
		}
	case map[string]interface{}:
		// Support {"min": x, "max": y} format
		minVal, hasMin := v["min"]
		maxVal, hasMax := v["max"]
		if hasMin && hasMax {
			return minVal, maxVal, true
		}
	}
	return nil, nil, false
}
