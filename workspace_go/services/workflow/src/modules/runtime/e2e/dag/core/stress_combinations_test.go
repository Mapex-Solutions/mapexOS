package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

/**
 * 4. Goto combinations
 */

// 4.1 — goto → code → set_state: teleport to async node
func TestGotoStress_TeleportToAsync(t *testing.T) {
	def := NewDefinition("Goto_Async").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(GotoSender("sender1", "PORTAL")).
		AddNode(GotoReceiver("receiver1", "PORTAL")).
		AddNode(CodeNode("code1", "return { val: 'teleported' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("code1", "val"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sender1").
		AddEdge("receiver1", "out", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"val": "teleported"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "teleported")
}

// 4.2 — goto → set_state → condition: teleport to routing logic
func TestGotoStress_TeleportToRouting(t *testing.T) {
	def := NewDefinition("Goto_Routing").
		WithState("val", "number", 0).
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(GotoSender("sender1", "ROUTE")).
		AddNode(GotoReceiver("receiver1", "ROUTE")).
		AddNode(SetStateNode("ss1", "set", "val", Literal("10"))).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("val"), Operator: "greaterThan", Value: Literal("5")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("high"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("low"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sender1").
		AddEdge("receiver1", "out", "ss1").
		AddEdge("ss1", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "high")
}

// 4.4 — set_state → goto → receiver → condition(state): state persists across teleport
func TestGotoStress_StatePersistsAcrossTeleport(t *testing.T) {
	def := NewDefinition("Goto_StatePersist").
		WithState("marker", "string", "").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "marker", Literal("set_before_teleport"))).
		AddNode(GotoSender("sender1", "CHECK")).
		AddNode(GotoReceiver("receiver1", "CHECK")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("marker"), Operator: "equals", Value: Literal("set_before_teleport")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("persisted"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("lost"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "sender1").
		AddEdge("receiver1", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "persisted")
}

/**
 * 5. Sequence with async
 */

// 5.1 — sequence(2) → [step: code, step: set_state] → done: mixed async+sync steps
func TestSequenceStress_AsyncStep(t *testing.T) {
	def := NewDefinition("Seq_Async").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SequenceNode("seq1", 3)).
		AddNode(CodeNode("code1", "return { val: 'step1' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("code1", "val"))).
		AddNode(EndNode("end")).
		// step_2 → code (async), step_3 → set_state (sync)
		AddEdge("__start__", "out", "seq1").
		AddEdge("seq1", "step_2", "code1").
		AddEdge("code1", "success", "seq1").
		AddEdge("code1", "error", "end").
		AddEdge("seq1", "step_3", "ss1").
		AddEdge("ss1", "out", "seq1").
		AddEdge("seq1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"val": "step1"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "step1")
}

/**
 * 6. Switch with async
 */

// 6.1 — set_state → switch(first) → [case: code → end, default: end]:
// switch routes to async branch
func TestSwitchStress_AsyncCase(t *testing.T) {
	def := NewDefinition("Switch_Async").
		WithState("val", "number", 10).
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SwitchNode("sw1", "first", []SwitchCaseItem{
			{ID: "high", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromState("val"), Operator: "greaterThan", Value: Literal("5")},
			}},
		})).
		AddNode(CodeNode("code1", "return { msg: 'high_branch' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("code1", "msg"))).
		AddNode(SetStateNode("ss_d", "set", "result", Literal("default"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sw1").
		AddEdge("sw1", "case_high", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "end").
		AddEdge("sw1", "default", "ss_d").
		AddEdge("ss_d", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"msg": "high_branch"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "high_branch")
}

// 6.2 — code → switch based on nodeOutput
func TestSwitchStress_NodeOutputRouting(t *testing.T) {
	def := NewDefinition("Switch_NodeOutput").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { status: 'warning' }", 5000)).
		AddNode(SwitchNode("sw1", "first", []SwitchCaseItem{
			{ID: "ok", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromNodeOutput("code1", "status"), Operator: "equals", Value: Literal("ok")},
			}},
			{ID: "warn", Logic: "AND", Conditions: []ConditionItem{
				{Field: FromNodeOutput("code1", "status"), Operator: "equals", Value: Literal("warning")},
			}},
		})).
		AddNode(SetStateNode("ss_ok", "set", "result", Literal("was_ok"))).
		AddNode(SetStateNode("ss_warn", "set", "result", Literal("was_warning"))).
		AddNode(SetStateNode("ss_d", "set", "result", Literal("was_default"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "sw1").
		AddEdge("code1", "error", "end").
		AddEdge("sw1", "case_ok", "ss_ok").
		AddEdge("sw1", "case_warn", "ss_warn").
		AddEdge("sw1", "default", "ss_d").
		AddEdge("ss_ok", "out", "end").
		AddEdge("ss_warn", "out", "end").
		AddEdge("ss_d", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"status": "warning"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "was_warning")
}

/**
 * 7. Wait_for + state mutations
 */

// 7.1 — set_state → wait_for: immediate match after state change
func TestWaitForStress_ImmediateAfterSetState(t *testing.T) {
	def := NewDefinition("WaitFor_Immediate").
		WithState("counter", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "counter", Literal("5"))).
		AddNode(WaitForLiteral("wf1", "counter", "equals", "5")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "wf1").
		AddEdge("wf1", "matched", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
}

/**
 * 8. End-to-end data propagation
 */

// 8.1 — event → condition → set_state → log: full data flow chain
func TestDataFlow_EventToConditionToStateToLog(t *testing.T) {
	def := NewDefinition("DataFlow_Full").
		WithState("alert", "string", "").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromEvent("temp"), Operator: "greaterThan", Value: Literal("25")},
		})).
		AddNode(SetStateNode("ss_t", "set", "alert", Literal("hot"))).
		AddNode(SetStateNode("ss_f", "set", "alert", Literal("normal"))).
		AddNode(LogNode("log1", "Alert: ${state.alert}", "info")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "log1").
		AddEdge("ss_f", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{"temp": 30})

	AssertCompleted(t, exec)
	AssertState(t, exec, "alert", "hot")
}

// 8.2 — code(return list) → loop(fromNodeOutput) → set_state(inc): code output feeds loop
func TestDataFlow_CodeOutputFeedsLoop(t *testing.T) {
	def := NewDefinition("DataFlow_CodeLoop").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { list: [1,2,3] }", 5000)).
		AddNode(LoopNode("loop1", FromNodeOutput("code1", "list"))).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "loop1").
		AddEdge("code1", "error", "end").
		AddEdge("loop1", "body", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{
				"list": []interface{}{1, 2, 3},
			}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 3)
}

// 8.3 — externalInput → set_state → condition: input drives state drives routing
func TestDataFlow_InputToStateToCondition(t *testing.T) {
	def := NewDefinition("DataFlow_Input").
		WithExternalInput("threshold", "10").
		WithState("val", "number", 0).
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "val", FromInput("threshold"))).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("val"), Operator: "greaterThan", Value: Literal("5")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("above"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("below"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "above")
}

/**
 * 9. Errors inside structures
 */

// 9.1 — code error → end(terminateWithError): error follows "error" handle → end with error
func TestErrorStress_CodeErrorToEndError(t *testing.T) {
	def := NewDefinition("Error_CodeToEnd").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(EndNodeWithError("end_err", "WORKFLOW_FAIL", "code failed")).
		AddNode(EndNode("end_ok")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end_ok").
		AddEdge("code1", "error", "end_err").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("SCRIPT_ERROR", "fail"),
		},
	)

	// Error follows "error" handle → end_err (terminateWithError) → fails with WORKFLOW_FAIL
	AssertFailed(t, exec, "WORKFLOW_FAIL")
}

/**
 * 10. Real-world complex scenarios
 */

// 10.1 — The user's real workflow pattern:
// start → fanout(3, firstCompleted) → [goto→code→end, goto→wait_signal→end, goto→end]
func TestComplex_FanoutGotoFirstCompleted(t *testing.T) {
	def := NewDefinition("Complex_FanoutGoto").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 3, "firstCompleted")).
		// Branch 1: goto → code → set_state → end
		AddNode(GotoSender("sender_notif", "NOTIFICATION")).
		AddNode(GotoReceiver("receiver_notif", "NOTIFICATION")).
		AddNode(CodeNode("code1", "return { text: 'notified' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("code1", "text"))).
		AddNode(EndNode("end1")).
		// Branch 2: goto → wait_signal → end
		AddNode(GotoSender("sender_signal", "SIGNAL")).
		AddNode(GotoReceiver("receiver_signal", "SIGNAL")).
		AddNode(WaitSignalNode("ws1", "continue")).
		AddNode(EndNode("end2")).
		// Branch 3: goto → end (immediate)
		AddNode(GotoSender("sender_loop", "LOOP")).
		AddNode(GotoReceiver("receiver_loop", "LOOP")).
		AddNode(EndNode("end3")).
		// Fanout edges
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "sender_notif").
		AddEdge("fan1", "out_2", "sender_signal").
		AddEdge("fan1", "out_3", "sender_loop").
		// Branch 1 edges
		AddEdge("receiver_notif", "out", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end1").
		AddEdge("ss1", "out", "end1").
		// Branch 2 edges
		AddEdge("receiver_signal", "out", "ws1").
		AddEdge("ws1", "out", "end2").
		// Branch 3 edges
		AddEdge("receiver_loop", "out", "end3").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"text": "notified"}, nil),
			// ws1 never gets callback — branch 3 completes first (immediate goto→end)
		},
	)

	// firstCompleted: branch 3 (goto→end) finishes immediately,
	// cancels the other waiting branches
	AssertCompleted(t, exec)
}

// 10.2 — Full stress: fanout(2, waitAll) → [loop→code→set_state, code→condition→set_state] → merge → end
func TestComplex_FullStress(t *testing.T) {
	def := NewDefinition("Complex_Full").
		WithState("loop_count", "number", 0).
		WithState("items", "array", []interface{}{"x", "y"}).
		WithState("branch2_path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		// Branch 1: loop → code → set_state(increment)
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(CodeNode("code_loop", "return { ok: true }", 5000)).
		AddNode(SetStateNode("ss_inc", "increment", "loop_count", Literal("1"))).
		// Branch 2: code → condition → set_state
		AddNode(CodeNode("code_branch", "return { status: 'ok' }", 5000)).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromNodeOutput("code_branch", "status"), Operator: "equals", Value: Literal("ok")},
		})).
		AddNode(SetStateNode("ss_t", "set", "branch2_path", Literal("ok_path"))).
		AddNode(SetStateNode("ss_f", "set", "branch2_path", Literal("fail_path"))).
		// Merge + end
		AddNode(MergeNode("merge1", 2)).
		AddNode(EndNode("end")).
		// Fanout edges
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "loop1").
		AddEdge("fan1", "out_2", "code_branch").
		// Branch 1: loop body
		AddEdge("loop1", "body", "code_loop").
		AddEdge("code_loop", "success", "ss_inc").
		AddEdge("code_loop", "error", "merge1").
		AddEdge("ss_inc", "out", "loop1").
		AddEdge("loop1", "done", "merge1").
		// Branch 2: code → condition
		AddEdge("code_branch", "success", "cond1").
		AddEdge("code_branch", "error", "merge1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "merge1").
		AddEdge("ss_f", "out", "merge1").
		// Merge → end
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code_loop":   CodeSuccessCallback(map[string]interface{}{"ok": true}, nil),
			"code_branch": CodeSuccessCallback(map[string]interface{}{"status": "ok"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "branch2_path", "ok_path")
}
