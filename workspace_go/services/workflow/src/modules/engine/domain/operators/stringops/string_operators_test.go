package stringops

import (
	"testing"
)

// --- ContainsOperator tests ---

func TestContainsOperator_Match(t *testing.T) {
	op := &ContainsOperator{}
	if op.Name() != "contains" {
		t.Fatalf("expected name 'contains', got %q", op.Name())
	}
	result, err := op.Evaluate("", "Hello World", "world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected case-insensitive match")
	}
}

func TestContainsOperator_NoMatch(t *testing.T) {
	op := &ContainsOperator{}
	result, err := op.Evaluate("", "Hello World", "xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected no match")
	}
}

func TestContainsOperator_NilField(t *testing.T) {
	op := &ContainsOperator{}
	result, err := op.Evaluate("", nil, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for nil field")
	}
}

func TestContainsOperator_NumericField(t *testing.T) {
	op := &ContainsOperator{}
	// toString converts 12345 → "12345", then check contains "234"
	result, err := op.Evaluate("", 12345, "234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected numeric field converted to string to contain '234'")
	}
}

// --- NotContainsOperator tests ---

func TestNotContainsOperator_NoMatch(t *testing.T) {
	op := &NotContainsOperator{}
	if op.Name() != "notContains" {
		t.Fatalf("expected name 'notContains', got %q", op.Name())
	}
	result, err := op.Evaluate("", "Hello World", "xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true when substring not found")
	}
}

func TestNotContainsOperator_Match(t *testing.T) {
	op := &NotContainsOperator{}
	result, err := op.Evaluate("", "Hello World", "world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false when substring found (case-insensitive)")
	}
}

func TestNotContainsOperator_NilField(t *testing.T) {
	op := &NotContainsOperator{}
	result, err := op.Evaluate("", nil, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for nil field (toString fails)")
	}
}

// --- StartsWithOperator tests ---

func TestStartsWithOperator_Match(t *testing.T) {
	op := &StartsWithOperator{}
	if op.Name() != "startsWith" {
		t.Fatalf("expected name 'startsWith', got %q", op.Name())
	}
	result, err := op.Evaluate("", "Hello World", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected case-insensitive prefix match")
	}
}

func TestStartsWithOperator_NoMatch(t *testing.T) {
	op := &StartsWithOperator{}
	result, err := op.Evaluate("", "Hello World", "World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected no prefix match")
	}
}

func TestStartsWithOperator_ExactMatch(t *testing.T) {
	op := &StartsWithOperator{}
	result, err := op.Evaluate("", "Hello", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected exact string to start with itself")
	}
}

// --- EndsWithOperator tests ---

func TestEndsWithOperator_Match(t *testing.T) {
	op := &EndsWithOperator{}
	if op.Name() != "endsWith" {
		t.Fatalf("expected name 'endsWith', got %q", op.Name())
	}
	result, err := op.Evaluate("", "Hello World", "WORLD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected case-insensitive suffix match")
	}
}

func TestEndsWithOperator_NoMatch(t *testing.T) {
	op := &EndsWithOperator{}
	result, err := op.Evaluate("", "Hello World", "Hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected no suffix match")
	}
}

func TestEndsWithOperator_EmptyString(t *testing.T) {
	op := &EndsWithOperator{}
	result, err := op.Evaluate("", "Hello", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected any string to end with empty suffix")
	}
}

// --- RegexOperator tests ---

func TestRegexOperator_Match(t *testing.T) {
	op := NewRegexOperator()
	if op.Name() != "regex" {
		t.Fatalf("expected name 'regex', got %q", op.Name())
	}
	result, err := op.Evaluate("", "hello123", "^[a-z]+[0-9]+$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected regex match")
	}
}

func TestRegexOperator_NoMatch(t *testing.T) {
	op := NewRegexOperator()
	result, err := op.Evaluate("", "hello", "^[0-9]+$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected no regex match")
	}
}

func TestRegexOperator_InvalidPattern(t *testing.T) {
	op := NewRegexOperator()
	// Invalid regex → returns (false, nil), not error
	result, err := op.Evaluate("", "hello", "[invalid")
	if err != nil {
		t.Fatalf("expected no error for invalid pattern, got %v", err)
	}
	if result {
		t.Fatal("expected false for invalid regex pattern")
	}
}

func TestRegexOperator_CacheHit(t *testing.T) {
	op := NewRegexOperator()
	pattern := "^test[0-9]+$"

	// First call compiles and caches
	result1, _ := op.Evaluate("", "test123", pattern)
	// Second call uses cache
	result2, _ := op.Evaluate("", "test456", pattern)

	if !result1 || !result2 {
		t.Fatal("expected both to match")
	}
}

func TestRegexOperator_NilField(t *testing.T) {
	op := NewRegexOperator()
	result, err := op.Evaluate("", nil, ".*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for nil field")
	}
}

func TestRegexOperator_EmailPattern(t *testing.T) {
	op := NewRegexOperator()
	result, err := op.Evaluate("", "user@example.com", `^[\w.-]+@[\w.-]+\.\w+$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected email regex match")
	}
}

// --- toString utility tests ---

func TestToString_String(t *testing.T) {
	s, ok := toString("hello")
	if !ok || s != "hello" {
		t.Fatalf("expected ('hello', true), got (%q, %v)", s, ok)
	}
}

func TestToString_Nil(t *testing.T) {
	_, ok := toString(nil)
	if ok {
		t.Fatal("expected ok=false for nil")
	}
}

func TestToString_Number(t *testing.T) {
	s, ok := toString(42)
	if !ok || s != "42" {
		t.Fatalf("expected ('42', true), got (%q, %v)", s, ok)
	}
}

func TestToString_ByteSlice(t *testing.T) {
	s, ok := toString([]byte("bytes"))
	if !ok || s != "bytes" {
		t.Fatalf("expected ('bytes', true), got (%q, %v)", s, ok)
	}
}
