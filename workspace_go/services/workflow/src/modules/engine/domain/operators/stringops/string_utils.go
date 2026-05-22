package stringops

import (
	"fmt"
	"strings"
)

/**
 * toString converts a value to string for string operations.
 * Numbers and other types are converted to their string representation.
 */
func toString(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}

	switch val := v.(type) {
	case string:
		return val, true
	case fmt.Stringer:
		return val.String(), true
	case []byte:
		return string(val), true
	default:
		// Convert other types to string representation
		return fmt.Sprintf("%v", val), true
	}
}

/**
 * normalizeForComparison converts values for case-insensitive comparison.
 */
func normalizeForComparison(s string, caseInsensitive bool) string {
	if caseInsensitive {
		return strings.ToLower(s)
	}
	return s
}
