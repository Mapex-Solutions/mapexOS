package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestFanout_TwoBranches_Sync(t *testing.T) {
	def := NewDefinition("Fanout2Sync").
		WithState("a", "string", "").
		WithState("b", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(SetStateNode("ss_a", "set", "a", Literal("branch_a"))).
		AddNode(SetStateNode("ss_b", "set", "b", Literal("branch_b"))).
		AddNode(MergeNode("merge1", 2)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "ss_a").
		AddEdge("fan1", "out_2", "ss_b").
		AddEdge("ss_a", "out", "merge1").
		AddEdge("ss_b", "out", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "a", "branch_a")
	AssertState(t, exec, "b", "branch_b")
}

func TestFanout_ThreeBranches_Sync(t *testing.T) {
	def := NewDefinition("Fanout3Sync").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 3, "")).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(SetStateNode("ss2", "increment", "count", Literal("1"))).
		AddNode(SetStateNode("ss3", "increment", "count", Literal("1"))).
		AddNode(MergeNode("merge1", 3)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "ss1").
		AddEdge("fan1", "out_2", "ss2").
		AddEdge("fan1", "out_3", "ss3").
		AddEdge("ss1", "out", "merge1").
		AddEdge("ss2", "out", "merge1").
		AddEdge("ss3", "out", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
}

func TestFanout_SingleBranch(t *testing.T) {
	def := NewDefinition("Fanout1").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 1, "")).
		AddNode(SetStateNode("ss1", "set", "result", Literal("done"))).
		AddNode(MergeNode("merge1", 1)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "ss1").
		AddEdge("ss1", "out", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "done")
}

func TestFanout_AsyncBranches(t *testing.T) {
	def := NewDefinition("FanoutAsync").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(CodeNode("code1", "return {r:1}", 5000)).
		AddNode(CodeNode("code2", "return {r:2}", 5000)).
		AddNode(MergeNode("merge1", 2)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code1").
		AddEdge("fan1", "out_2", "code2").
		AddEdge("code1", "success", "merge1").
		AddEdge("code1", "error", "merge1").
		AddEdge("code2", "success", "merge1").
		AddEdge("code2", "error", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"r": 1}, nil),
			"code2": CodeSuccessCallback(map[string]interface{}{"r": 2}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code1")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code2")
}

func TestFanout_ZeroBranches_Error(t *testing.T) {
	def := NewDefinition("Fanout0").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 0, "")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "EXECUTION_ERROR")
}

func TestFanout_NegativeBranches_Error(t *testing.T) {
	def := NewDefinition("FanoutNeg").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", -1, "")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "EXECUTION_ERROR")
}
