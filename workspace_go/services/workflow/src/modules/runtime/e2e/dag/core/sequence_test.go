package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestSequence_ThreeSteps(t *testing.T) {
	// Executor increments currentStep before emitting handle: step_{currentStep+1}
	// With steps=3: call1→step_2, call2→step_3, call3→done (2 body executions)
	def := NewDefinition("Seq3").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(SequenceNode("seq1", 3)).
		AddNode(SetStateNode("ss2", "increment", "count", Literal("1"))).
		AddNode(SetStateNode("ss3", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "seq1").
		AddEdge("seq1", "step_2", "ss2").
		AddEdge("ss2", "out", "seq1").
		AddEdge("seq1", "step_3", "ss3").
		AddEdge("ss3", "out", "seq1").
		AddEdge("seq1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 2)
}

func TestSequence_SingleStep(t *testing.T) {
	// With steps=1: call1→done immediately (currentStep 0→1, 1>=1 → done)
	def := NewDefinition("Seq1").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SequenceNode("seq1", 1)).
		AddNode(SetStateNode("ss_done", "set", "result", Literal("done"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "seq1").
		AddEdge("seq1", "done", "ss_done").
		AddEdge("ss_done", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "done")
}

func TestSequence_ZeroSteps(t *testing.T) {
	def := NewDefinition("Seq0").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SequenceNode("seq1", 0)).
		AddNode(SetStateNode("ss_done", "set", "result", Literal("done"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "seq1").
		AddEdge("seq1", "done", "ss_done").
		AddEdge("ss_done", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "done")
}
