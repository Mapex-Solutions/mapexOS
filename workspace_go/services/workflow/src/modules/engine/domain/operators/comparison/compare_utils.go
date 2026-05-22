package comparison

import (
	"reflect"
	"time"

	typeconv "github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
)

/**
 * Compare compares two values and returns their relative order.
 * Handles numeric, string, and time comparisons.
 *
 * Returns:
 * - CompareLess (-1) if a < b
 * - CompareEqual (0) if a == b
 * - CompareGreater (1) if a > b
 * - CompareError if comparison is not possible
 */
func Compare(a, b interface{}) CompareResult {
	// Handle nil cases
	if a == nil && b == nil {
		return CompareEqual
	}
	if a == nil || b == nil {
		return CompareError
	}

	// Try numeric comparison first
	aNum, aOk := typeconv.TryFloat64(a)
	bNum, bOk := typeconv.TryFloat64(b)
	if aOk && bOk {
		return compareNumbers(aNum, bNum)
	}

	// Try string comparison
	aStr, aOk := typeconv.TryString(a)
	bStr, bOk := typeconv.TryString(b)
	if aOk && bOk {
		return compareStrings(aStr, bStr)
	}

	// Try time comparison
	aTime, aOk := tryTimeUTC(a)
	bTime, bOk := tryTimeUTC(b)
	if aOk && bOk {
		return compareTimes(aTime, bTime)
	}

	// Try boolean comparison (equal only)
	aBool, aOk := typeconv.TryBool(a)
	bBool, bOk := typeconv.TryBool(b)
	if aOk && bOk {
		if aBool == bBool {
			return CompareEqual
		}
		return CompareError // Booleans can only be equal/not equal
	}

	return CompareError
}

/**
 * Equals checks if two values are equal.
 * Handles type coercion for common cases.
 */
func Equals(a, b interface{}) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Try direct comparison with reflect
	if reflect.DeepEqual(a, b) {
		return true
	}

	// Try numeric comparison
	aNum, aOk := typeconv.TryFloat64(a)
	bNum, bOk := typeconv.TryFloat64(b)
	if aOk && bOk {
		return aNum == bNum
	}

	// Try string comparison
	aStr, aOk := typeconv.TryString(a)
	bStr, bOk := typeconv.TryString(b)
	if aOk && bOk {
		return aStr == bStr
	}

	// Try boolean comparison
	aBool, aOk := typeconv.TryBool(a)
	bBool, bOk := typeconv.TryBool(b)
	if aOk && bOk {
		return aBool == bBool
	}

	return false
}

/* Type Conversion Helpers */

// tryTimeUTC is a thin wrapper around typeconv.TryTime using UTC as default timezone.
// Comparison operations don't have a user-specified timezone, so UTC is used.
func tryTimeUTC(v interface{}) (time.Time, bool) {
	return typeconv.TryTime(v, "")
}

/* Comparison Functions */

func compareNumbers(a, b float64) CompareResult {
	if a < b {
		return CompareLess
	}
	if a > b {
		return CompareGreater
	}
	return CompareEqual
}

func compareStrings(a, b string) CompareResult {
	if a < b {
		return CompareLess
	}
	if a > b {
		return CompareGreater
	}
	return CompareEqual
}

func compareTimes(a, b time.Time) CompareResult {
	if a.Before(b) {
		return CompareLess
	}
	if a.After(b) {
		return CompareGreater
	}
	return CompareEqual
}
