package services

import (
	"testing"

	
)

// createMatchConfig creates a MatchConfig for testing.
func createMatchConfig(policy string, rules []MatchRule) *MatchConfig {
	return &MatchConfig{
		Policy: policy,
		Rules:  rules,
	}
}

/**
 * TEST: Evaluate with nil config
 */

func TestEvaluate_NilConfig(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"temperature": 30.0,
	}

	result, err := evaluator.Evaluate(event, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !result.ShouldProcess {
		t.Error("Expected ShouldProcess to be true for nil config")
	}

	if len(result.Conditions) != 0 {
		t.Errorf("Expected 0 conditions, got: %d", len(result.Conditions))
	}
}

/**
 * TEST: Policy "all" (AND logic)
 */

func TestEvaluate_PolicyAll_AllPass(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"temperature": 30.0,
		"humidity":    50.0,
	}

	matchConfig := createMatchConfig("all", []MatchRule{
		{Field: "temperature", Operator: "gt", Value: float64(25)},
		{Field: "humidity", Operator: "lt", Value: float64(60)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !result.ShouldProcess {
		t.Error("Expected ShouldProcess to be true when all conditions pass")
	}

	if len(result.Conditions) != 2 {
		t.Errorf("Expected 2 conditions, got: %d", len(result.Conditions))
	}

	for _, cond := range result.Conditions {
		if !cond.Passed {
			t.Errorf("Expected condition %s to pass", cond.Field)
		}
	}
}

func TestEvaluate_PolicyAll_SomeFail(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"temperature": 20.0, // < 25, should fail
		"humidity":    50.0,
	}

	matchConfig := createMatchConfig("all", []MatchRule{
		{Field: "temperature", Operator: "gt", Value: float64(25)},
		{Field: "humidity", Operator: "lt", Value: float64(60)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.ShouldProcess {
		t.Error("Expected ShouldProcess to be false when not all conditions pass")
	}
}

/**
 * TEST: Policy "any" (OR logic)
 */

func TestEvaluate_PolicyAny_SomePass(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"temperature": 20.0, // < 25, should fail
		"humidity":    50.0, // < 60, should pass
	}

	matchConfig := createMatchConfig("any", []MatchRule{
		{Field: "temperature", Operator: "gt", Value: float64(25)},
		{Field: "humidity", Operator: "lt", Value: float64(60)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !result.ShouldProcess {
		t.Error("Expected ShouldProcess to be true when at least one condition passes")
	}
}

func TestEvaluate_PolicyAny_NonePasses(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"temperature": 20.0, // < 25, should fail
		"humidity":    70.0, // > 60, should fail
	}

	matchConfig := createMatchConfig("any", []MatchRule{
		{Field: "temperature", Operator: "gt", Value: float64(25)},
		{Field: "humidity", Operator: "lt", Value: float64(60)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.ShouldProcess {
		t.Error("Expected ShouldProcess to be false when no conditions pass")
	}
}

/**
 * TEST: Operator "eq"
 */

func TestCompareValues_Eq_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if !evaluator.CompareValues("active", "eq", "active") {
		t.Error("Expected 'active' eq 'active' to be true")
	}

	if !evaluator.CompareValues(float64(100), "eq", float64(100)) {
		t.Error("Expected 100 eq 100 to be true")
	}
}

func TestCompareValues_Eq_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues("active", "eq", "inactive") {
		t.Error("Expected 'active' eq 'inactive' to be false")
	}
}

/**
 * TEST: Operator "neq"
 */

func TestCompareValues_Neq_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if !evaluator.CompareValues("active", "neq", "inactive") {
		t.Error("Expected 'active' neq 'inactive' to be true")
	}
}

func TestCompareValues_Neq_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues("active", "neq", "active") {
		t.Error("Expected 'active' neq 'active' to be false")
	}
}

/**
 * TEST: Operator "gt"
 */

func TestCompareValues_Gt_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if !evaluator.CompareValues(float64(30), "gt", float64(25)) {
		t.Error("Expected 30 > 25 to be true")
	}

	if !evaluator.CompareValues(int(30), "gt", int(25)) {
		t.Error("Expected int 30 > 25 to be true")
	}
}

func TestCompareValues_Gt_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues(float64(20), "gt", float64(25)) {
		t.Error("Expected 20 > 25 to be false")
	}

	if evaluator.CompareValues(float64(25), "gt", float64(25)) {
		t.Error("Expected 25 > 25 to be false")
	}
}

/**
 * TEST: Operator "gte"
 */

func TestCompareValues_Gte_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if !evaluator.CompareValues(float64(30), "gte", float64(25)) {
		t.Error("Expected 30 >= 25 to be true")
	}

	if !evaluator.CompareValues(float64(25), "gte", float64(25)) {
		t.Error("Expected 25 >= 25 to be true")
	}
}

func TestCompareValues_Gte_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues(float64(20), "gte", float64(25)) {
		t.Error("Expected 20 >= 25 to be false")
	}
}

/**
 * TEST: Operator "lt"
 */

func TestCompareValues_Lt_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if !evaluator.CompareValues(float64(20), "lt", float64(25)) {
		t.Error("Expected 20 < 25 to be true")
	}
}

func TestCompareValues_Lt_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues(float64(30), "lt", float64(25)) {
		t.Error("Expected 30 < 25 to be false")
	}

	if evaluator.CompareValues(float64(25), "lt", float64(25)) {
		t.Error("Expected 25 < 25 to be false")
	}
}

/**
 * TEST: Operator "lte"
 */

func TestCompareValues_Lte_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if !evaluator.CompareValues(float64(20), "lte", float64(25)) {
		t.Error("Expected 20 <= 25 to be true")
	}

	if !evaluator.CompareValues(float64(25), "lte", float64(25)) {
		t.Error("Expected 25 <= 25 to be true")
	}
}

func TestCompareValues_Lte_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues(float64(30), "lte", float64(25)) {
		t.Error("Expected 30 <= 25 to be false")
	}
}

/**
 * TEST: Operator "in"
 */

func TestCompareValues_In_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	list := []interface{}{"active", "pending", "review"}
	if !evaluator.CompareValues("active", "in", list) {
		t.Error("Expected 'active' in list to be true")
	}
}

func TestCompareValues_In_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	list := []interface{}{"active", "pending", "review"}
	if evaluator.CompareValues("inactive", "in", list) {
		t.Error("Expected 'inactive' in list to be false")
	}
}

/**
 * TEST: Operator "nin"
 */

func TestCompareValues_Nin_Match(t *testing.T) {
	evaluator := NewMatchEvaluator()

	list := []interface{}{"error", "failed"}
	if !evaluator.CompareValues("success", "nin", list) {
		t.Error("Expected 'success' not in list to be true")
	}
}

func TestCompareValues_Nin_NoMatch(t *testing.T) {
	evaluator := NewMatchEvaluator()

	list := []interface{}{"error", "failed"}
	if evaluator.CompareValues("error", "nin", list) {
		t.Error("Expected 'error' not in list to be false")
	}
}

/**
 * TEST: Unknown operator
 */

func TestCompareValues_UnknownOperator(t *testing.T) {
	evaluator := NewMatchEvaluator()

	if evaluator.CompareValues("value", "unknown", "value") {
		t.Error("Expected unknown operator to return false")
	}
}

/**
 * TEST: Nested field path
 */

func TestEvaluate_NestedFieldPath(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"payload": map[string]interface{}{
			"sensor": map[string]interface{}{
				"temperature": 30.0,
			},
		},
	}

	matchConfig := createMatchConfig("all", []MatchRule{
		{Field: "payload.sensor.temperature", Operator: "gt", Value: float64(25)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !result.ShouldProcess {
		t.Error("Expected ShouldProcess to be true for nested field match")
	}

	if len(result.Conditions) != 1 {
		t.Errorf("Expected 1 condition, got: %d", len(result.Conditions))
	}

	if result.Conditions[0].Actual != float64(30) {
		t.Errorf("Expected actual value 30, got: %v", result.Conditions[0].Actual)
	}
}

/**
 * TEST: Field not found
 */

func TestEvaluate_FieldNotFound(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"humidity": 50.0,
	}

	matchConfig := createMatchConfig("all", []MatchRule{
		{Field: "temperature", Operator: "gt", Value: float64(25)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.ShouldProcess {
		t.Error("Expected ShouldProcess to be false when field not found")
	}

	if len(result.Conditions) != 1 {
		t.Errorf("Expected 1 condition, got: %d", len(result.Conditions))
	}

	if result.Conditions[0].Passed {
		t.Error("Expected condition to fail when field not found")
	}
}

/**
 * TEST: History entries
 */

func TestEvaluate_HistoryEntries(t *testing.T) {
	evaluator := NewMatchEvaluator()

	event := map[string]interface{}{
		"temperature": 30.0,
	}

	matchConfig := createMatchConfig("all", []MatchRule{
		{Field: "temperature", Operator: "gt", Value: float64(25)},
	})

	result, err := evaluator.Evaluate(event, matchConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Should have history entries for each rule + final policy decision
	if len(result.History) < 2 {
		t.Errorf("Expected at least 2 history entries, got: %d", len(result.History))
	}
}

/**
 * TEST: Non-numeric comparison
 */

func TestCompareValues_NonNumericGt(t *testing.T) {
	evaluator := NewMatchEvaluator()

	// String comparison with gt should fail
	if evaluator.CompareValues("abc", "gt", "def") {
		t.Error("Expected string gt comparison to return false")
	}
}

/**
 * TEST: Type conversions
 */

func TestCompareValues_TypeConversions(t *testing.T) {
	evaluator := NewMatchEvaluator()

	// int vs float64
	if !evaluator.CompareValues(int(30), "gt", float64(25)) {
		t.Error("Expected int 30 > float64 25 to be true")
	}

	// int64
	if !evaluator.CompareValues(int64(30), "gt", float64(25)) {
		t.Error("Expected int64 30 > float64 25 to be true")
	}

	// uint
	if !evaluator.CompareValues(uint(30), "gt", float64(25)) {
		t.Error("Expected uint 30 > float64 25 to be true")
	}
}
