package services

import (
	"fmt"
	"reflect"
	"strings"
)

/**
 * MATCH EVALUATOR
 * Domain service for evaluating conditional routing rules.
 * Contains pure business logic for match conditions (eq, neq, gt, lt, in, nin)
 * and policy evaluation (all/any). No I/O or infrastructure dependencies.
 */

func NewMatchEvaluator() *MatchEvaluator {
	return &MatchEvaluator{
		operatorTextMap: map[string]string{
			"eq":  "equals",
			"neq": "not equals",
			"gt":  "greater than",
			"gte": "greater than or equals",
			"lt":  "less than",
			"lte": "less than or equals",
			"in":  "in",
			"nin": "not in",
		},
	}
}

// Evaluate evaluates all match rules based on policy (all/any).
// Returns structured results with individual condition evaluations.
func (m *MatchEvaluator) Evaluate(event interface{}, matchConfig *MatchConfig) (*EvaluationResult, error) {
	if matchConfig == nil {
		return &EvaluationResult{
			ShouldProcess: true,
			Conditions:    []ConditionResult{},
			History:       []string{"no match rules defined → allowed"},
		}, nil
	}

	result := &EvaluationResult{
		ShouldProcess: false,
		Conditions:    make([]ConditionResult, 0),
		History:       make([]string, 0),
	}

	passedCount := 0
	totalRules := len(matchConfig.Rules)

	for _, rule := range matchConfig.Rules {
		condResult, historyLine := m.evaluateRule(event, rule)
		result.Conditions = append(result.Conditions, condResult)
		result.History = append(result.History, historyLine)

		if condResult.Passed {
			passedCount++
		}
	}

	switch matchConfig.Policy {
	case "all":
		result.ShouldProcess = passedCount == totalRules
		if result.ShouldProcess {
			result.History = append(result.History, fmt.Sprintf("policy: all → %d/%d rules passed → allowed", passedCount, totalRules))
		} else {
			result.History = append(result.History, fmt.Sprintf("policy: all → %d/%d rules passed → denied", passedCount, totalRules))
		}

	case "any":
		result.ShouldProcess = passedCount > 0
		if result.ShouldProcess {
			result.History = append(result.History, fmt.Sprintf("policy: any → %d/%d rules passed → allowed", passedCount, totalRules))
		} else {
			result.History = append(result.History, fmt.Sprintf("policy: any → %d/%d rules passed → denied", passedCount, totalRules))
		}

	default:
		result.ShouldProcess = passedCount == totalRules
	}

	return result, nil
}

func (m *MatchEvaluator) evaluateRule(event interface{}, rule MatchRule) (ConditionResult, string) {
	condResult := ConditionResult{
		Field:    rule.Field,
		Operator: rule.Operator,
		Expected: rule.Value,
		Actual:   nil,
		Passed:   false,
	}

	fieldValue, err := m.getFieldValue(event, rule.Field)
	if err != nil {
		historyLine := fmt.Sprintf("\"%s\" %s %v → denied (field not found)", rule.Field, m.operatorTextMap[rule.Operator], rule.Value)
		return condResult, historyLine
	}

	condResult.Actual = fieldValue
	condResult.Passed = m.CompareValues(fieldValue, rule.Operator, rule.Value)

	status := "allowed"
	if !condResult.Passed {
		status = "denied"
	}

	historyLine := fmt.Sprintf("\"%s\" %s %v → %s", rule.Field, m.operatorTextMap[rule.Operator], rule.Value, status)
	return condResult, historyLine
}

func (m *MatchEvaluator) getFieldValue(event interface{}, fieldPath string) (interface{}, error) {
	eventMap, ok := event.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("event is not a map")
	}

	parts := strings.Split(fieldPath, ".")

	current := interface{}(eventMap)
	for i, part := range parts {
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field path '%s' not found at segment '%s'", fieldPath, strings.Join(parts[:i+1], "."))
		}

		value, exists := currentMap[part]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found in path '%s'", part, fieldPath)
		}

		current = value
	}

	return current, nil
}

// CompareValues compares two values using the specified operator.
func (m *MatchEvaluator) CompareValues(actualValue interface{}, operator string, expectedValue interface{}) bool {
	switch operator {
	case "eq":
		return m.equals(actualValue, expectedValue)
	case "neq":
		return !m.equals(actualValue, expectedValue)
	case "gt":
		return m.greaterThan(actualValue, expectedValue)
	case "gte":
		return m.greaterThan(actualValue, expectedValue) || m.equals(actualValue, expectedValue)
	case "lt":
		return m.lessThan(actualValue, expectedValue)
	case "lte":
		return m.lessThan(actualValue, expectedValue) || m.equals(actualValue, expectedValue)
	case "in":
		return m.in(actualValue, expectedValue)
	case "nin":
		return !m.in(actualValue, expectedValue)
	default:
		return false
	}
}

func (m *MatchEvaluator) equals(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (m *MatchEvaluator) greaterThan(a, b interface{}) bool {
	aFloat, aOk := m.toFloat64(a)
	bFloat, bOk := m.toFloat64(b)

	if !aOk || !bOk {
		return false
	}

	return aFloat > bFloat
}

func (m *MatchEvaluator) lessThan(a, b interface{}) bool {
	aFloat, aOk := m.toFloat64(a)
	bFloat, bOk := m.toFloat64(b)

	if !aOk || !bOk {
		return false
	}

	return aFloat < bFloat
}

func (m *MatchEvaluator) in(value interface{}, array interface{}) bool {
	arrayValue := reflect.ValueOf(array)
	if arrayValue.Kind() != reflect.Slice && arrayValue.Kind() != reflect.Array {
		return false
	}

	for i := 0; i < arrayValue.Len(); i++ {
		if m.equals(value, arrayValue.Index(i).Interface()) {
			return true
		}
	}

	return false
}

func (m *MatchEvaluator) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}
