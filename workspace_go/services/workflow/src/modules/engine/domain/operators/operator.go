package operators

/*
 * ConditionOperator defines the interface for operators that evaluate
 * a single condition by comparing two values.
 *
 * Examples: equals, greaterThan, contains, beforeDate, etc.
 *
 * Each operator receives:
 * - timezone: IANA timezone for date/time operations
 * - fieldValue: The left-hand side value (from event, state, variable, or literal)
 * - compareValue: The right-hand side value to compare against
 *
 * Operators should:
 * - Handle type conversions internally
 * - Return (false, nil) for comparison failures (not errors)
 * - Return (_, error) only for actual errors (invalid config, etc.)
 */
type ConditionOperator interface {
	// Name returns the operator identifier used in rule conditions
	// Example: "equals", "greaterThan", "contains"
	Name() string

	// Evaluate compares fieldValue against compareValue
	// Returns (matched, error)
	Evaluate(timezone string, fieldValue, compareValue interface{}) (bool, error)
}

/**
 * GroupOperator defines the interface for operators that combine
 * multiple condition results into a single boolean.
 *
 * Examples: and, or, nor, nand
 *
 * Group operators implement logical combination:
 * - AND: All conditions must be true
 * - OR: At least one condition must be true
 * - NOR: None of the conditions can be true (NOT OR)
 * - NAND: Not all conditions are true (NOT AND)
 *
 * These operators support short-circuit evaluation:
 * - AND returns false immediately on first false
 * - OR returns true immediately on first true
 */
type GroupOperator interface {
	// Name returns the operator identifier used in rule groups
	// Example: "and", "or", "nor", "nand"
	Name() string

	// Evaluate combines multiple condition results
	// The results slice contains the evaluated result of each condition/group
	Evaluate(results []bool) bool

	// SupportsShortCircuit returns true if this operator can exit early
	// Used by the engine to optimize evaluation
	SupportsShortCircuit() bool

	// ShouldShortCircuit checks if evaluation should stop based on current result
	// For AND: returns true if result is false
	// For OR: returns true if result is true
	ShouldShortCircuit(currentResult bool) bool
}

/**
 * BetweenOperator is a specialized interface for operators that
 * compare a value against a range (min, max).
 *
 * Examples: between, betweenDate, betweenTime
 */
type BetweenOperator interface {
	// Name returns the operator identifier
	Name() string

	// EvaluateRange checks if fieldValue is within the range [minValue, maxValue]
	// The inclusive parameter determines if boundaries are included
	EvaluateRange(
		timezone string,
		fieldValue interface{},
		minValue interface{},
		maxValue interface{},
		inclusive bool,
	) (bool, error)
}

/**
 * OperatorMetadata provides information about an operator for
 * documentation and validation purposes.
 */
type OperatorMetadata struct {
	// Operator name (e.g., "equals", "greaterThan")
	Name string

	// Category for grouping (e.g., "comparison", "string", "datetime", "group")
	Category string

	// Human-readable description
	Description string

	// Expected types for fieldValue and compareValue
	// Empty means any type is accepted
	AcceptedTypes []string

	// Whether this operator supports the "between" pattern
	IsBetween bool
}

/**
 * OperatorWithMetadata extends the basic operator interfaces
 * with metadata information.
 */
type OperatorWithMetadata interface {
	// Metadata returns information about this operator
	Metadata() OperatorMetadata
}
