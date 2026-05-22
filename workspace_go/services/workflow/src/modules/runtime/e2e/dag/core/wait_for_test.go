package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestWaitFor_ImmediateMatch(t *testing.T) {
	def := NewDefinition("WaitForMatch").
		WithState("counter", "number", 5).
		AddNode(StartNode("__start__")).
		AddNode(WaitForLiteral("wf1", "counter", "equals", "5")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "wf1").
		AddEdge("wf1", "matched", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "wf1")
}

func TestWaitFor_Suspends(t *testing.T) {
	def := NewDefinition("WaitForSuspend").
		WithState("counter", "number", 3).
		AddNode(StartNode("__start__")).
		AddNode(WaitForLiteral("wf1", "counter", "equals", "5")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "wf1").
		AddEdge("wf1", "matched", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"wf1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
}

func TestWaitFor_GreaterThan(t *testing.T) {
	def := NewDefinition("WaitForGT").
		WithState("val", "number", 10).
		AddNode(StartNode("__start__")).
		AddNode(WaitForLiteral("wf1", "val", "greaterThan", "5")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "wf1").
		AddEdge("wf1", "matched", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
}

func TestWaitFor_LessThan_Suspends(t *testing.T) {
	def := NewDefinition("WaitForLT").
		WithState("val", "number", 10).
		AddNode(StartNode("__start__")).
		AddNode(WaitForLiteral("wf1", "val", "lessThan", "5")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "wf1").
		AddEdge("wf1", "matched", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"wf1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
}
