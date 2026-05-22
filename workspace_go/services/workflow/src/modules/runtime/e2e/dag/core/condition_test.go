package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestCondition_True(t *testing.T) {
	def := NewDefinition("CondTrue").
		WithState("temp", "number", 30).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("temp"), Operator: "greaterThan", Value: Literal("20")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("true_path"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("false_path"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "true_path")
}

func TestCondition_False(t *testing.T) {
	def := NewDefinition("CondFalse").
		WithState("temp", "number", 10).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("temp"), Operator: "greaterThan", Value: Literal("20")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("true_path"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("false_path"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "false_path")
}

func TestCondition_AND_AllTrue(t *testing.T) {
	def := NewDefinition("CondANDAllTrue").
		WithState("a", "number", 10).
		WithState("b", "number", 20).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("a"), Operator: "greaterThan", Value: Literal("5")},
			{Field: FromState("b"), Operator: "greaterThan", Value: Literal("15")},
		})).
		AddNode(SetStateNode("ss_t", "set", "result", Literal("true"))).
		AddNode(SetStateNode("ss_f", "set", "result", Literal("false"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "true")
}

func TestCondition_AND_OneFalse(t *testing.T) {
	def := NewDefinition("CondANDOneFalse").
		WithState("a", "number", 10).
		WithState("b", "number", 5).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("a"), Operator: "greaterThan", Value: Literal("5")},
			{Field: FromState("b"), Operator: "greaterThan", Value: Literal("15")},
		})).
		AddNode(SetStateNode("ss_t", "set", "result", Literal("true"))).
		AddNode(SetStateNode("ss_f", "set", "result", Literal("false"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "false")
}

func TestCondition_OR_OneTrue(t *testing.T) {
	def := NewDefinition("CondOROneTrue").
		WithState("a", "number", 10).
		WithState("b", "number", 5).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "OR", []ConditionItem{
			{Field: FromState("a"), Operator: "greaterThan", Value: Literal("50")},
			{Field: FromState("b"), Operator: "equals", Value: Literal("5")},
		})).
		AddNode(SetStateNode("ss_t", "set", "result", Literal("true"))).
		AddNode(SetStateNode("ss_f", "set", "result", Literal("false"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "true")
}

func TestCondition_OR_AllFalse(t *testing.T) {
	def := NewDefinition("CondORAllFalse").
		WithState("a", "number", 1).
		WithState("b", "number", 2).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "OR", []ConditionItem{
			{Field: FromState("a"), Operator: "greaterThan", Value: Literal("50")},
			{Field: FromState("b"), Operator: "greaterThan", Value: Literal("50")},
		})).
		AddNode(SetStateNode("ss_t", "set", "result", Literal("true"))).
		AddNode(SetStateNode("ss_f", "set", "result", Literal("false"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "false")
}

func TestCondition_FieldFromEvent(t *testing.T) {
	def := NewDefinition("CondFromEvent").
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromEvent("data.value"), Operator: "equals", Value: Literal("42")},
		})).
		AddNode(SetStateNode("ss_t", "set", "result", Literal("true"))).
		AddNode(SetStateNode("ss_f", "set", "result", Literal("false"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{
		"data": map[string]interface{}{"value": 42},
	})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "true")
}

func TestCondition_FieldFromInput(t *testing.T) {
	def := NewDefinition("CondFromInput").
		WithExternalInput("threshold", "20").
		WithState("temp", "number", 25).
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("temp"), Operator: "greaterThan", Value: FromInput("threshold")},
		})).
		AddNode(SetStateNode("ss_t", "set", "result", Literal("true"))).
		AddNode(SetStateNode("ss_f", "set", "result", Literal("false"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "true")
}
