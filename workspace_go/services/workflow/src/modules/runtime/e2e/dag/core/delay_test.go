package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestDelay_Suspends(t *testing.T) {
	def := NewDefinition("DelaySuspend").
		AddNode(StartNode("__start__")).
		AddNode(DelayNode("delay1", 5, "seconds")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "delay1").
		AddEdge("delay1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "delay1")
}

func TestDelay_Minutes(t *testing.T) {
	def := NewDefinition("DelayMinutes").
		AddNode(StartNode("__start__")).
		AddNode(DelayNode("delay1", 1, "minutes")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "delay1").
		AddEdge("delay1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
}

func TestDelay_AbbreviatedUnit(t *testing.T) {
	def := NewDefinition("DelayAbbrev").
		AddNode(StartNode("__start__")).
		AddNode(DelayNode("delay1", 10, "s")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "delay1").
		AddEdge("delay1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
}

func TestDelay_UnknownUnit_Error(t *testing.T) {
	def := NewDefinition("DelayBadUnit").
		AddNode(StartNode("__start__")).
		AddNode(DelayNode("delay1", 5, "weeks")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "delay1").
		AddEdge("delay1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "EXECUTION_ERROR")
}
