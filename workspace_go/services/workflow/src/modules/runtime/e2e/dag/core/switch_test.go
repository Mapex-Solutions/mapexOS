package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestSwitch_FirstMatch(t *testing.T) {
	def := NewDefinition("SwitchFirst").
		WithState("val", "number", 10).
		AddNode(StartNode("__start__")).
		AddNode(SwitchNode("sw1", "first", []SwitchCaseItem{
			{ID: "a", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromState("val"), Operator: "greaterThan", Value: Literal("5")},
			}},
			{ID: "b", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromState("val"), Operator: "greaterThan", Value: Literal("8")},
			}},
		})).
		AddNode(SetStateNode("ss_a", "set", "result", Literal("case_a"))).
		AddNode(SetStateNode("ss_b", "set", "result", Literal("case_b"))).
		AddNode(SetStateNode("ss_d", "set", "result", Literal("default"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sw1").
		AddEdge("sw1", "case_a", "ss_a").
		AddEdge("sw1", "case_b", "ss_b").
		AddEdge("sw1", "default", "ss_d").
		AddEdge("ss_a", "out", "end").
		AddEdge("ss_b", "out", "end").
		AddEdge("ss_d", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// First match mode → stops at case_a even though case_b also matches
	AssertState(t, exec, "result", "case_a")
}

func TestSwitch_NoMatch_Default(t *testing.T) {
	def := NewDefinition("SwitchDefault").
		WithState("val", "number", 1).
		AddNode(StartNode("__start__")).
		AddNode(SwitchNode("sw1", "first", []SwitchCaseItem{
			{ID: "a", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromState("val"), Operator: "greaterThan", Value: Literal("100")},
			}},
		})).
		AddNode(SetStateNode("ss_a", "set", "result", Literal("case_a"))).
		AddNode(SetStateNode("ss_d", "set", "result", Literal("default"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sw1").
		AddEdge("sw1", "case_a", "ss_a").
		AddEdge("sw1", "default", "ss_d").
		AddEdge("ss_a", "out", "end").
		AddEdge("ss_d", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "default")
}

func TestSwitch_SecondMatch_First(t *testing.T) {
	def := NewDefinition("SwitchSecond").
		WithState("val", "number", 10).
		AddNode(StartNode("__start__")).
		AddNode(SwitchNode("sw1", "first", []SwitchCaseItem{
			{ID: "a", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromState("val"), Operator: "greaterThan", Value: Literal("100")},
			}},
			{ID: "b", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromState("val"), Operator: "greaterThan", Value: Literal("5")},
			}},
		})).
		AddNode(SetStateNode("ss_a", "set", "result", Literal("case_a"))).
		AddNode(SetStateNode("ss_b", "set", "result", Literal("case_b"))).
		AddNode(SetStateNode("ss_d", "set", "result", Literal("default"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sw1").
		AddEdge("sw1", "case_a", "ss_a").
		AddEdge("sw1", "case_b", "ss_b").
		AddEdge("sw1", "default", "ss_d").
		AddEdge("ss_a", "out", "end").
		AddEdge("ss_b", "out", "end").
		AddEdge("ss_d", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "case_b")
}

func TestSwitch_EmptyCases_Default(t *testing.T) {
	def := NewDefinition("SwitchEmpty").
		AddNode(StartNode("__start__")).
		AddNode(SwitchNode("sw1", "first", nil)).
		AddNode(SetStateNode("ss_d", "set", "result", Literal("default"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sw1").
		AddEdge("sw1", "default", "ss_d").
		AddEdge("ss_d", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "default")
}
