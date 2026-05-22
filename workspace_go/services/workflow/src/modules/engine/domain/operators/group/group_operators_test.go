package group

import (
	"testing"
)

// --- AndOperator tests ---

func TestAndOperator_AllTrue(t *testing.T) {
	op := &AndOperator{}
	if op.Name() != "AND" {
		t.Fatalf("expected name 'AND', got %q", op.Name())
	}
	if !op.Evaluate([]bool{true, true, true}) {
		t.Fatal("expected true for all true")
	}
}

func TestAndOperator_OneFalse(t *testing.T) {
	op := &AndOperator{}
	if op.Evaluate([]bool{true, false, true}) {
		t.Fatal("expected false when one is false")
	}
}

func TestAndOperator_AllFalse(t *testing.T) {
	op := &AndOperator{}
	if op.Evaluate([]bool{false, false, false}) {
		t.Fatal("expected false for all false")
	}
}

func TestAndOperator_Empty(t *testing.T) {
	op := &AndOperator{}
	if !op.Evaluate([]bool{}) {
		t.Fatal("expected true for empty (vacuous truth)")
	}
}

func TestAndOperator_ShortCircuit(t *testing.T) {
	op := &AndOperator{}
	if !op.SupportsShortCircuit() {
		t.Fatal("expected short-circuit support")
	}
	// Should short-circuit on false
	if !op.ShouldShortCircuit(false) {
		t.Fatal("expected short-circuit on false")
	}
	if op.ShouldShortCircuit(true) {
		t.Fatal("should NOT short-circuit on true")
	}
}

// --- OrOperator tests ---

func TestOrOperator_OneTrue(t *testing.T) {
	op := &OrOperator{}
	if op.Name() != "OR" {
		t.Fatalf("expected name 'OR', got %q", op.Name())
	}
	if !op.Evaluate([]bool{false, true, false}) {
		t.Fatal("expected true when one is true")
	}
}

func TestOrOperator_AllFalse(t *testing.T) {
	op := &OrOperator{}
	if op.Evaluate([]bool{false, false, false}) {
		t.Fatal("expected false for all false")
	}
}

func TestOrOperator_AllTrue(t *testing.T) {
	op := &OrOperator{}
	if !op.Evaluate([]bool{true, true, true}) {
		t.Fatal("expected true for all true")
	}
}

func TestOrOperator_Empty(t *testing.T) {
	op := &OrOperator{}
	if op.Evaluate([]bool{}) {
		t.Fatal("expected false for empty")
	}
}

func TestOrOperator_ShortCircuit(t *testing.T) {
	op := &OrOperator{}
	if !op.SupportsShortCircuit() {
		t.Fatal("expected short-circuit support")
	}
	// Should short-circuit on true
	if !op.ShouldShortCircuit(true) {
		t.Fatal("expected short-circuit on true")
	}
	if op.ShouldShortCircuit(false) {
		t.Fatal("should NOT short-circuit on false")
	}
}

// --- NandOperator tests ---

func TestNandOperator_AllTrue(t *testing.T) {
	op := &NandOperator{}
	if op.Name() != "NAND" {
		t.Fatalf("expected name 'NAND', got %q", op.Name())
	}
	// NAND: NOT(AND) — all true → false
	if op.Evaluate([]bool{true, true, true}) {
		t.Fatal("expected false for NAND(all true)")
	}
}

func TestNandOperator_OneFalse(t *testing.T) {
	op := &NandOperator{}
	// NAND: not all true → true
	if !op.Evaluate([]bool{true, false, true}) {
		t.Fatal("expected true for NAND(one false)")
	}
}

func TestNandOperator_AllFalse(t *testing.T) {
	op := &NandOperator{}
	if !op.Evaluate([]bool{false, false, false}) {
		t.Fatal("expected true for NAND(all false)")
	}
}

func TestNandOperator_Empty(t *testing.T) {
	op := &NandOperator{}
	// NAND of empty = false (opposite of AND's vacuous true)
	if op.Evaluate([]bool{}) {
		t.Fatal("expected false for NAND(empty)")
	}
}

func TestNandOperator_ShortCircuit(t *testing.T) {
	op := &NandOperator{}
	if !op.SupportsShortCircuit() {
		t.Fatal("expected short-circuit support")
	}
	// Should short-circuit on false (NAND is true)
	if !op.ShouldShortCircuit(false) {
		t.Fatal("expected short-circuit on false")
	}
	if op.ShouldShortCircuit(true) {
		t.Fatal("should NOT short-circuit on true")
	}
}

// --- NorOperator tests ---

func TestNorOperator_AllFalse(t *testing.T) {
	op := &NorOperator{}
	if op.Name() != "NOR" {
		t.Fatalf("expected name 'NOR', got %q", op.Name())
	}
	// NOR: NOT(OR) — none true → true
	if !op.Evaluate([]bool{false, false, false}) {
		t.Fatal("expected true for NOR(all false)")
	}
}

func TestNorOperator_OneTrue(t *testing.T) {
	op := &NorOperator{}
	if op.Evaluate([]bool{false, true, false}) {
		t.Fatal("expected false for NOR(one true)")
	}
}

func TestNorOperator_AllTrue(t *testing.T) {
	op := &NorOperator{}
	if op.Evaluate([]bool{true, true, true}) {
		t.Fatal("expected false for NOR(all true)")
	}
}

func TestNorOperator_Empty(t *testing.T) {
	op := &NorOperator{}
	// NOR of empty = true (no true condition exists)
	if !op.Evaluate([]bool{}) {
		t.Fatal("expected true for NOR(empty)")
	}
}

func TestNorOperator_ShortCircuit(t *testing.T) {
	op := &NorOperator{}
	if !op.SupportsShortCircuit() {
		t.Fatal("expected short-circuit support")
	}
	// Should short-circuit on true (NOR is false)
	if !op.ShouldShortCircuit(true) {
		t.Fatal("expected short-circuit on true")
	}
	if op.ShouldShortCircuit(false) {
		t.Fatal("should NOT short-circuit on false")
	}
}

// --- Metadata tests ---

func TestOperatorMetadata(t *testing.T) {
	ops := []struct {
		name     string
		operator interface{ Name() string }
	}{
		{"AND", &AndOperator{}},
		{"OR", &OrOperator{}},
		{"NAND", &NandOperator{}},
		{"NOR", &NorOperator{}},
	}
	for _, tt := range ops {
		if tt.operator.Name() != tt.name {
			t.Errorf("expected name %q, got %q", tt.name, tt.operator.Name())
		}
	}
}
