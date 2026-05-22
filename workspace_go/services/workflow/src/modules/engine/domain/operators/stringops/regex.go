package stringops

import (
	"regexp"
	"sync"
	"sync/atomic"

	"workflow/src/modules/engine/domain/operators"
)

// maxRegexCacheSize is the maximum number of compiled regex patterns to cache.
// When exceeded, the cache is cleared to prevent unbounded memory growth.
const maxRegexCacheSize = 1000

/**
 * RegexOperator checks if a string matches a regular expression pattern.
 * Uses Go's regexp package (RE2 syntax).
 *
 * Examples:
 * - "hello123" matches "^[a-z]+[0-9]+$" -> true
 * - "user@example.com" matches "^[\w.-]+@[\w.-]+$" -> true
 * - "invalid" matches "^[0-9]+$" -> false
 *
 * Note: Compiled patterns are cached for performance (max 1000 entries).
 */
type RegexOperator struct {
	cache sync.Map
	size  atomic.Int64
}

// Ensure RegexOperator implements the interface
var _ operators.ConditionOperator = (*RegexOperator)(nil)

/**
 * Name returns the operator identifier.
 */
func (o *RegexOperator) Name() string {
	return "regex"
}

/**
 * Evaluate checks if fieldValue matches the regex pattern in compareValue.
 *
 * @param timezone - IANA timezone (unused for string ops, kept for interface compliance)
 * @param fieldValue - The string to test
 * @param compareValue - The regex pattern
 * @returns (true, nil) if fieldValue matches the pattern, (false, nil) if pattern is invalid
 */
func (o *RegexOperator) Evaluate(
	timezone string,
	fieldValue, compareValue interface{},
) (bool, error) {
	// Convert field value to string
	fieldStr, ok := toString(fieldValue)
	if !ok {
		return false, nil
	}

	// Get the pattern
	patternStr, ok := toString(compareValue)
	if !ok {
		return false, nil
	}

	// Get or compile the regex
	re, err := o.getOrCompile(patternStr)
	if err != nil {
		// Invalid regex pattern - return false, no error
		// This allows the rule to continue evaluation
		return false, nil
	}

	return re.MatchString(fieldStr), nil
}

/**
 * getOrCompile retrieves a cached compiled regex or compiles a new one.
 * Evicts all entries when cache exceeds maxRegexCacheSize to prevent unbounded growth.
 */
func (o *RegexOperator) getOrCompile(pattern string) (*regexp.Regexp, error) {
	// Check cache first
	if cached, ok := o.cache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}

	// Compile the pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Evict all if cache is full (simple reset strategy)
	if o.size.Load() >= maxRegexCacheSize {
		o.cache.Range(func(key, _ any) bool {
			o.cache.Delete(key)
			return true
		})
		o.size.Store(0)
	}

	// Store in cache
	o.cache.Store(pattern, re)
	o.size.Add(1)
	return re, nil
}

/**
 * Metadata returns information about this operator.
 */
func (o *RegexOperator) Metadata() operators.OperatorMetadata {
	return operators.OperatorMetadata{
		Name:          "regex",
		Category:      "string",
		Description:   "Checks if a string matches a regular expression pattern (RE2 syntax)",
		AcceptedTypes: []string{"string"},
		IsBetween:     false,
	}
}

/**
 * NewRegexOperator creates a new RegexOperator with an empty cache.
 */
func NewRegexOperator() *RegexOperator {
	return &RegexOperator{}
}
