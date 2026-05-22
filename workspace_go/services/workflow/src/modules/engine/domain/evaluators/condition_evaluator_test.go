package evaluators

import (
	"errors"
	"testing"

	defEntities "workflow/src/modules/definitions/domain/entities"
	"workflow/src/modules/engine/domain/operators"
	"workflow/src/modules/engine/domain/operators/comparison"
	"workflow/src/modules/engine/domain/operators/stringops"
)

// newTestRegistry creates a registry with basic operators for testing.
func newTestRegistry() *operators.OperatorRegistry {
	r := operators.NewOperatorRegistry()
	r.RegisterCondition(&comparison.EqualsOperator{})
	r.RegisterCondition(&comparison.NotEqualsOperator{})
	r.RegisterCondition(&comparison.GreaterThanOperator{})
	r.RegisterCondition(&comparison.LessThanOperator{})
	r.RegisterCondition(&stringops.ContainsOperator{})
	r.RegisterBetween(&comparison.BetweenOperator{})
	return r
}

// condItem creates a ConditionGroupItem wrapping a ConditionItem.
func condItem(field defEntities.FieldValue, op string, value defEntities.FieldValue) defEntities.ConditionGroupItem {
	return defEntities.ConditionGroupItem{
		Type: "condition",
		Data: defEntities.ConditionItem{
			ID:       "c1",
			Field:    field,
			Operator: op,
			Value:    value,
		},
	}
}

// groupItem creates a ConditionGroupItem wrapping a nested ConditionGroup.
func groupItem(logic defEntities.GroupLogicOperator, items ...defEntities.ConditionGroupItem) defEntities.ConditionGroupItem {
	return defEntities.ConditionGroupItem{
		Type: "group",
		Data: defEntities.ConditionGroup{
			ID:    "g1",
			Logic: logic,
			Items: items,
		},
	}
}

func litFV(val string) defEntities.FieldValue {
	return defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: val}
}

func stateFV(path string) defEntities.FieldValue {
	return defEntities.FieldValue{Type: defEntities.FieldValueState, Value: path}
}

func TestEvaluateGroup_EmptyGroup(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	_, err := eval.EvaluateGroup(nil, "", nil, nil, nil, nil)
	if !errors.Is(err, ErrEmptyGroup) {
		t.Fatalf("expected ErrEmptyGroup, got %v", err)
	}

	emptyGroup := &defEntities.ConditionGroup{Logic: defEntities.LogicAND, Items: nil}
	_, err = eval.EvaluateGroup(emptyGroup, "", nil, nil, nil, nil)
	if !errors.Is(err, ErrEmptyGroup) {
		t.Fatalf("expected ErrEmptyGroup, got %v", err)
	}
}

func TestEvaluateGroup_AND_AllTrue(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello", "y": "hello"}
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicAND,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "equals", litFV("hello")),
			condItem(stateFV("y"), "equals", litFV("hello")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true (AND: all match)")
	}
}

func TestEvaluateGroup_AND_OneFalse(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello", "y": "world"}
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicAND,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "equals", litFV("hello")),
			condItem(stateFV("y"), "equals", litFV("hello")), // false
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false (AND: one doesn't match)")
	}
}

func TestEvaluateGroup_OR_OneTrue(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello", "y": "world"}
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicOR,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "equals", litFV("nope")),  // false
			condItem(stateFV("y"), "equals", litFV("world")), // true
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true (OR: one matches)")
	}
}

func TestEvaluateGroup_OR_NoneTrue(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello"}
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicOR,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "equals", litFV("a")),
			condItem(stateFV("x"), "equals", litFV("b")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false (OR: none matches)")
	}
}

func TestEvaluateGroup_NAND(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello", "y": "hello"}
	// NAND: NOT(AND) → all true → NAND returns false
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicNAND,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "equals", litFV("hello")),
			condItem(stateFV("y"), "equals", litFV("hello")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false (NAND: all true → false)")
	}
}

func TestEvaluateGroup_NOR(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello"}
	// NOR: NOT(OR) → none true → NOR returns true
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicNOR,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "equals", litFV("a")),
			condItem(stateFV("x"), "equals", litFV("b")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true (NOR: none match → true)")
	}
}

func TestEvaluateGroup_NestedGroup(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello", "y": "world"}
	// OR(AND(x=hello, y=world), x=nope) → true
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicOR,
		Items: []defEntities.ConditionGroupItem{
			groupItem(defEntities.LogicAND,
				condItem(stateFV("x"), "equals", litFV("hello")),
				condItem(stateFV("y"), "equals", litFV("world")),
			),
			condItem(stateFV("x"), "equals", litFV("nope")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true (nested AND inside OR)")
	}
}

func TestEvaluateGroup_UnknownOperator(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello"}
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicAND,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("x"), "unknownOp", litFV("hello")),
		},
	}
	_, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if !errors.Is(err, ErrOperatorNotFound) {
		t.Fatalf("expected ErrOperatorNotFound, got %v", err)
	}
}

func TestEvaluateGroup_FieldNotFound_NonMatch(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"x": "hello"}
	// Field "missing" not in state → resolve fails → treated as non-match (false, nil)
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicAND,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("missing"), "equals", litFV("hello")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false (field not found → non-match)")
	}
}

func TestEvaluateGroup_StringContains(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"msg": "Hello World"}
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicAND,
		Items: []defEntities.ConditionGroupItem{
			condItem(stateFV("msg"), "contains", litFV("world")),
		},
	}
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected true (case-insensitive contains)")
	}
}

func TestEvaluateGroup_Between(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	state := map[string]interface{}{"temp": 25}
	// Between: value must be a map with min/max
	group := &defEntities.ConditionGroup{
		Logic: defEntities.LogicAND,
		Items: []defEntities.ConditionGroupItem{
			{
				Type: "condition",
				Data: defEntities.ConditionItem{
					ID:       "c1",
					Field:    stateFV("temp"),
					Operator: "between",
					Value: defEntities.FieldValue{
						Type:  defEntities.FieldValueLiteral,
						Value: "",
					},
				},
			},
		},
	}
	// Between requires a map value, literal string won't work → resolve gives string → non-match
	result, err := eval.EvaluateGroup(group, "", nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Between expects map with min/max, literal gives string → extractRange fails → non-match
	if result {
		t.Fatal("expected false (between with invalid range)")
	}
}

func TestEvaluateGroup_UnknownLogic(t *testing.T) {
	eval := NewConditionEvaluator(newTestRegistry())
	group := &defEntities.ConditionGroup{
		Logic: "XNOR",
		Items: []defEntities.ConditionGroupItem{
			condItem(litFV("a"), "equals", litFV("a")),
		},
	}
	_, err := eval.EvaluateGroup(group, "", nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for unknown logic operator")
	}
}
