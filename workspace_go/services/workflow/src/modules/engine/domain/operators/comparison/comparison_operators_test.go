package comparison

import (
	"testing"
	"time"
)

// --- Compare utility tests ---

func TestCompare_Numeric(t *testing.T) {
	if r := Compare(10, 20); r != CompareLess {
		t.Fatalf("expected CompareLess, got %d", r)
	}
	if r := Compare(20, 10); r != CompareGreater {
		t.Fatalf("expected CompareGreater, got %d", r)
	}
	if r := Compare(10, 10); r != CompareEqual {
		t.Fatalf("expected CompareEqual, got %d", r)
	}
}

func TestCompare_NumericCoercion(t *testing.T) {
	// int vs float64
	if r := Compare(10, 10.0); r != CompareEqual {
		t.Fatalf("expected CompareEqual for int vs float64, got %d", r)
	}
	// string number vs int
	if r := Compare("5", 10); r != CompareLess {
		t.Fatalf("expected CompareLess for string '5' vs int 10, got %d", r)
	}
}

func TestCompare_Strings(t *testing.T) {
	if r := Compare("apple", "banana"); r != CompareLess {
		t.Fatalf("expected CompareLess, got %d", r)
	}
	if r := Compare("banana", "apple"); r != CompareGreater {
		t.Fatalf("expected CompareGreater, got %d", r)
	}
	if r := Compare("same", "same"); r != CompareEqual {
		t.Fatalf("expected CompareEqual, got %d", r)
	}
}

func TestCompare_NilCases(t *testing.T) {
	if r := Compare(nil, nil); r != CompareEqual {
		t.Fatalf("expected CompareEqual for nil vs nil, got %d", r)
	}
	if r := Compare(nil, 1); r != CompareError {
		t.Fatalf("expected CompareError for nil vs 1, got %d", r)
	}
	if r := Compare(1, nil); r != CompareError {
		t.Fatalf("expected CompareError for 1 vs nil, got %d", r)
	}
}

func TestCompare_Booleans(t *testing.T) {
	if r := Compare(true, true); r != CompareEqual {
		t.Fatalf("expected CompareEqual for true==true, got %d", r)
	}
	// TryFloat64 converts booleans: true→1.0, false→0.0
	// So Compare(true, false) = Compare(1.0, 0.0) = CompareGreater
	if r := Compare(true, false); r != CompareGreater {
		t.Fatalf("expected CompareGreater for true(1) vs false(0), got %d", r)
	}
	if r := Compare(false, true); r != CompareLess {
		t.Fatalf("expected CompareLess for false(0) vs true(1), got %d", r)
	}
}

func TestCompare_Times(t *testing.T) {
	t1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	if r := Compare(t1, t2); r != CompareLess {
		t.Fatalf("expected CompareLess, got %d", r)
	}
	if r := Compare(t2, t1); r != CompareGreater {
		t.Fatalf("expected CompareGreater, got %d", r)
	}
	if r := Compare(t1, t1); r != CompareEqual {
		t.Fatalf("expected CompareEqual, got %d", r)
	}
}

func TestCompare_IncompatibleTypes(t *testing.T) {
	// Structs get string-compared via TryString/Sprintf fallback
	// so they become lexicographic comparison of "{1}" vs "{2}"
	type custom struct{ x int }
	r := Compare(custom{1}, custom{2})
	if r == CompareError {
		t.Fatalf("expected string fallback comparison, got CompareError")
	}
}

// --- Equals utility tests ---

func TestEquals_SameType(t *testing.T) {
	if !Equals("hello", "hello") {
		t.Fatal("expected equal strings")
	}
	if !Equals(42, 42) {
		t.Fatal("expected equal ints")
	}
	if Equals("a", "b") {
		t.Fatal("expected not equal strings")
	}
}

func TestEquals_NumericCoercion(t *testing.T) {
	if !Equals(42, 42.0) {
		t.Fatal("expected int == float64")
	}
	if !Equals(float32(42), 42) {
		t.Fatal("expected float32 == int")
	}
}

func TestEquals_BoolCoercion(t *testing.T) {
	if !Equals(true, true) {
		t.Fatal("expected true == true")
	}
	if Equals(true, false) {
		t.Fatal("expected true != false")
	}
}

func TestEquals_NilCases(t *testing.T) {
	if !Equals(nil, nil) {
		t.Fatal("expected nil == nil")
	}
	if Equals(nil, "something") {
		t.Fatal("expected nil != string")
	}
	if Equals("something", nil) {
		t.Fatal("expected string != nil")
	}
}

// --- EqualsOperator tests ---

func TestEqualsOperator_Match(t *testing.T) {
	op := &EqualsOperator{}
	if op.Name() != "equals" {
		t.Fatalf("expected name 'equals', got %q", op.Name())
	}
	result, err := op.Evaluate("", "hello", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true for equal values")
	}
}

func TestEqualsOperator_NoMatch(t *testing.T) {
	op := &EqualsOperator{}
	result, err := op.Evaluate("", "hello", "world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for different values")
	}
}

func TestEqualsOperator_NumericCoercion(t *testing.T) {
	op := &EqualsOperator{}
	result, err := op.Evaluate("", 42, 42.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true for int vs float64")
	}
}

// --- NotEqualsOperator tests ---

func TestNotEqualsOperator_Match(t *testing.T) {
	op := &NotEqualsOperator{}
	if op.Name() != "notEquals" {
		t.Fatalf("expected name 'notEquals', got %q", op.Name())
	}
	result, err := op.Evaluate("", "hello", "world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true for different values")
	}
}

func TestNotEqualsOperator_SameValues(t *testing.T) {
	op := &NotEqualsOperator{}
	result, err := op.Evaluate("", "hello", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for equal values")
	}
}

// --- GreaterThanOperator tests ---

func TestGreaterThanOperator_Numbers(t *testing.T) {
	op := &GreaterThanOperator{}
	if op.Name() != "greaterThan" {
		t.Fatalf("expected name 'greaterThan', got %q", op.Name())
	}

	result, err := op.Evaluate("", 20, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 20 > 10")
	}

	result, err = op.Evaluate("", 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 10 NOT > 20")
	}

	result, err = op.Evaluate("", 10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 10 NOT > 10 (equal)")
	}
}

func TestGreaterThanOperator_Strings(t *testing.T) {
	op := &GreaterThanOperator{}
	result, err := op.Evaluate("", "banana", "apple")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected banana > apple lexicographically")
	}
}

// --- GreaterThanEqualsOperator tests ---

func TestGreaterThanEqualsOperator(t *testing.T) {
	op := &GreaterThanEqualsOperator{}
	if op.Name() != "greaterThanEquals" {
		t.Fatalf("expected name 'greaterThanEquals', got %q", op.Name())
	}

	// Greater
	result, err := op.Evaluate("", 20, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 20 >= 10")
	}

	// Equal
	result, err = op.Evaluate("", 10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 10 >= 10")
	}

	// Less
	result, err = op.Evaluate("", 5, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 5 NOT >= 10")
	}
}

// --- LessThanOperator tests ---

func TestLessThanOperator_Numbers(t *testing.T) {
	op := &LessThanOperator{}
	if op.Name() != "lessThan" {
		t.Fatalf("expected name 'lessThan', got %q", op.Name())
	}

	result, err := op.Evaluate("", 5, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 5 < 10")
	}

	result, err = op.Evaluate("", 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 10 NOT < 5")
	}

	result, err = op.Evaluate("", 10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 10 NOT < 10 (equal)")
	}
}

// --- LessThanEqualsOperator tests ---

func TestLessThanEqualsOperator(t *testing.T) {
	op := &LessThanEqualsOperator{}
	if op.Name() != "lessThanEquals" {
		t.Fatalf("expected name 'lessThanEquals', got %q", op.Name())
	}

	// Less
	result, err := op.Evaluate("", 5, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 5 <= 10")
	}

	// Equal
	result, err = op.Evaluate("", 10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 10 <= 10")
	}

	// Greater
	result, err = op.Evaluate("", 20, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 20 NOT <= 10")
	}
}

// --- BetweenOperator tests ---

func TestBetweenOperator_InRange_Slice(t *testing.T) {
	op := &BetweenOperator{}
	if op.Name() != "between" {
		t.Fatalf("expected name 'between', got %q", op.Name())
	}

	// value=15 in [10, 20] → true
	result, err := op.Evaluate("", 15, []interface{}{10, 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 15 between [10, 20]")
	}
}

func TestBetweenOperator_AtBoundary_Inclusive(t *testing.T) {
	op := &BetweenOperator{}
	// Evaluate uses inclusive by default
	result, err := op.Evaluate("", 10, []interface{}{10, 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 10 between [10, 20] inclusive")
	}

	result, err = op.Evaluate("", 20, []interface{}{10, 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 20 between [10, 20] inclusive")
	}
}

func TestBetweenOperator_OutOfRange(t *testing.T) {
	op := &BetweenOperator{}
	result, err := op.Evaluate("", 25, []interface{}{10, 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 25 NOT between [10, 20]")
	}
}

func TestBetweenOperator_MapFormat(t *testing.T) {
	op := &BetweenOperator{}
	result, err := op.Evaluate("", 15, map[string]interface{}{"min": 10, "max": 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 15 between {min:10, max:20}")
	}
}

func TestBetweenOperator_InvalidRange(t *testing.T) {
	op := &BetweenOperator{}
	// Single element slice — extractMinMax fails
	result, err := op.Evaluate("", 15, []interface{}{10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for invalid range (single element)")
	}
}

func TestBetweenOperator_EvaluateRange_Exclusive(t *testing.T) {
	op := &BetweenOperator{}
	// Exclusive: boundary NOT included
	result, err := op.EvaluateRange("", 10, 10, 20, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 10 NOT in (10, 20) exclusive")
	}

	result, err = op.EvaluateRange("", 15, 10, 20, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 15 in (10, 20) exclusive")
	}
}

func TestBetweenOperator_Strings(t *testing.T) {
	op := &BetweenOperator{}
	result, err := op.Evaluate("", "b", []interface{}{"a", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 'b' between ['a', 'c']")
	}
}

func TestBetweenOperator_Metadata(t *testing.T) {
	op := &BetweenOperator{}
	m := op.Metadata()
	if m.Name != "between" {
		t.Fatalf("expected name 'between', got %q", m.Name)
	}
	if !m.IsBetween {
		t.Fatal("expected IsBetween=true")
	}
}

// --- extractMinMax tests ---

func TestExtractMinMax_FloatSlice(t *testing.T) {
	min, max, ok := extractMinMax([]float64{1.5, 3.5})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if min != 1.5 || max != 3.5 {
		t.Fatalf("expected [1.5, 3.5], got [%v, %v]", min, max)
	}
}

func TestExtractMinMax_IntSlice(t *testing.T) {
	min, max, ok := extractMinMax([]int{1, 10})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if min != 1 || max != 10 {
		t.Fatalf("expected [1, 10], got [%v, %v]", min, max)
	}
}

func TestExtractMinMax_StringSlice(t *testing.T) {
	min, max, ok := extractMinMax([]string{"a", "z"})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if min != "a" || max != "z" {
		t.Fatalf("expected [a, z], got [%v, %v]", min, max)
	}
}

func TestExtractMinMax_InvalidType(t *testing.T) {
	_, _, ok := extractMinMax("not a slice")
	if ok {
		t.Fatal("expected ok=false for string")
	}
}
