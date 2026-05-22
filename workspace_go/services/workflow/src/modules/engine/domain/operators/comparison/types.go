package comparison

/*
 * COMPARISON TYPES
 * Type aliases and constants for comparison operations.
 */

// CompareResult represents the result of a comparison operation.
type CompareResult int

const (
	CompareLess    CompareResult = -1
	CompareEqual   CompareResult = 0
	CompareGreater CompareResult = 1
	CompareError   CompareResult = -999 // Indicates comparison is not possible
)
