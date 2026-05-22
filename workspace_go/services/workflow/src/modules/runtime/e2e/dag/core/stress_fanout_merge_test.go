package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

// 3.1 — fanout(2) → [code1, code2] → merge: both async, waitAll
func TestFanoutStress_TwoAsyncBranches(t *testing.T) {
	def := NewDefinition("Fanout_2Async").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(CodeNode("code1", "return {r:1}", 5000)).
		AddNode(CodeNode("code2", "return {r:2}", 5000)).
		AddNode(MergeNode("merge1", 2)).
		AddNode(SetStateNode("ss1", "set", "result", Literal("merged"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code1").
		AddEdge("fan1", "out_2", "code2").
		AddEdge("code1", "success", "merge1").
		AddEdge("code1", "error", "merge1").
		AddEdge("code2", "success", "merge1").
		AddEdge("code2", "error", "merge1").
		AddEdge("merge1", "out", "ss1").
		AddEdge("ss1", "out", "end").
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
	AssertState(t, exec, "result", "merged")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code1")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code2")
}

// 3.2 — fanout(2) → [code, delay] → merge: mixed async types
func TestFanoutStress_MixedAsyncTypes(t *testing.T) {
	def := NewDefinition("Fanout_Mixed").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(DelayNode("delay1", 5, "seconds")).
		AddNode(MergeNode("merge1", 2)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code1").
		AddEdge("fan1", "out_2", "delay1").
		AddEdge("code1", "success", "merge1").
		AddEdge("code1", "error", "merge1").
		AddEdge("delay1", "out", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1":  CodeSuccessCallback(map[string]interface{}{}, nil),
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
}

// 3.3 — fanout(2) → [set_state("a"), set_state("b")] → merge → condition:
// sync branches, state from both visible after merge
func TestFanoutStress_SyncBranchesStateMerge(t *testing.T) {
	def := NewDefinition("Fanout_SyncState").
		WithState("a", "string", "").
		WithState("b", "string", "").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(SetStateNode("ss_a", "set", "a", Literal("branch_a"))).
		AddNode(SetStateNode("ss_b", "set", "b", Literal("branch_b"))).
		AddNode(MergeNode("merge1", 2)).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("a"), Operator: "equals", Value: Literal("branch_a")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("both_set"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("missing"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "ss_a").
		AddEdge("fan1", "out_2", "ss_b").
		AddEdge("ss_a", "out", "merge1").
		AddEdge("ss_b", "out", "merge1").
		AddEdge("merge1", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "a", "branch_a")
	AssertState(t, exec, "b", "branch_b")
	AssertState(t, exec, "path", "both_set")
}

// 3.4 — fanout(2, firstCompleted) → [code_fast, code_slow]:
// first branch completes, second is cancelled
func TestFanoutStress_FirstCompleted(t *testing.T) {
	def := NewDefinition("Fanout_FirstCompleted").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "firstCompleted")).
		AddNode(CodeNode("fast", "return {}", 5000)).
		AddNode(CodeNode("slow", "return {}", 30000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "fast").
		AddEdge("fan1", "out_2", "slow").
		AddEdge("fast", "success", "end").
		AddEdge("fast", "error", "end").
		AddEdge("slow", "success", "end").
		AddEdge("slow", "error", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"fast": CodeSuccessCallback(map[string]interface{}{}, nil),
			// slow never gets a callback — should be cancelled by firstCompleted
		},
	)

	AssertCompleted(t, exec)
}

// 3.5 — fanout(3) → [code, wait_signal, delay] → merge: 3 different async types
func TestFanoutStress_ThreeAsyncTypes(t *testing.T) {
	def := NewDefinition("Fanout_3Async").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 3, "")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(WaitSignalNode("ws1", "test_signal")).
		AddNode(DelayNode("delay1", 10, "seconds")).
		AddNode(MergeNode("merge1", 3)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code1").
		AddEdge("fan1", "out_2", "ws1").
		AddEdge("fan1", "out_3", "delay1").
		AddEdge("code1", "success", "merge1").
		AddEdge("code1", "error", "merge1").
		AddEdge("ws1", "out", "merge1").
		AddEdge("delay1", "out", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1":  CodeSuccessCallback(map[string]interface{}{}, nil),
			"ws1":    SuccessCallback(nil),
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
}

// 3.6 — fanout(2) → [code_success, code_error] → merge: error branch follows "error" handle
func TestFanoutStress_OneBranchError(t *testing.T) {
	def := NewDefinition("Fanout_OneError").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(CodeNode("code_ok", "return {}", 5000)).
		AddNode(CodeNode("code_err", "throw new Error('boom')", 5000)).
		AddNode(MergeNode("merge1", 2)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code_ok").
		AddEdge("fan1", "out_2", "code_err").
		AddEdge("code_ok", "success", "merge1").
		AddEdge("code_ok", "error", "merge1").
		AddEdge("code_err", "success", "merge1").
		AddEdge("code_err", "error", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code_ok":  CodeSuccessCallback(map[string]interface{}{}, nil),
			"code_err": CodeErrorCallback("SCRIPT_ERROR", "boom"),
		},
	)

	// Error follows "error" handle → merge → end → completes (user handled the error)
	AssertCompleted(t, exec)
}
