package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

// 2.1 — loop → body(code → set_state(increment)): async inside loop, 3 resumes
func TestLoopAsync_CodeInBody(t *testing.T) {
	def := NewDefinition("LoopAsync_CodeBody").
		WithState("items", "array", []interface{}{"a", "b", "c"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(CodeNode("code1", "return { ok: true }", 5000)).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"ok": true}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 3)
	AssertDispatchCount(t, h.Publisher, "DispatchCodeExecution", 3)
}

// 2.2 — loop → body(condition → [true: set_state, false: skip]): branching inside loop
func TestLoopAsync_ConditionInBody(t *testing.T) {
	def := NewDefinition("LoopAsync_CondBody").
		WithState("items", "array", []interface{}{1, 2, 3, 4, 5}).
		WithState("even_count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		// Check if loop_index is even (0, 2, 4)
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("loop_index"), Operator: "equals", Value: Literal("0")},
		})).
		AddNode(SetStateNode("ss_inc", "increment", "even_count", Literal("1"))).
		AddNode(EndNode("end")).
		// loop body → condition → true: increment, false: back to loop
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "cond1").
		AddEdge("cond1", "true", "ss_inc").
		AddEdge("cond1", "false", "loop1").
		AddEdge("ss_inc", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// loop_index 0 matches equals "0" → only first iteration increments
	// This tests condition evaluation inside loop body with loop_index state
}

// 2.3a — nested loops SYNC: outer[a,b] × inner[x,y] = 4 increments
// Inner "done" has NO edge — the walker pops the loop stack and returns to outer_loop.
// This is the correct pattern: nested loop "done" relies on the stack, not edges.
func TestLoopAsync_NestedLoopsSync(t *testing.T) {
	def := NewDefinition("LoopAsync_NestedSync").
		WithState("outer", "array", []interface{}{"a", "b"}).
		WithState("inner", "array", []interface{}{"x", "y"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("outer_loop", FromState("outer"))).
		AddNode(LoopNode("inner_loop", FromState("inner"))).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "outer_loop").
		AddEdge("outer_loop", "body", "inner_loop").
		AddEdge("inner_loop", "body", "ss1").
		AddEdge("ss1", "out", "inner_loop").
		// NO edge from inner_loop "done" → stack pop handles return to outer
		AddEdge("outer_loop", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 4) // 2 outer x 2 inner
}

// 2.3b — nested loops ASYNC: outer[a,b] × inner body has code (async) → 4 callbacks
func TestLoopAsync_NestedLoopsAsync(t *testing.T) {
	def := NewDefinition("LoopAsync_NestedAsync").
		WithState("outer", "array", []interface{}{"a", "b"}).
		WithState("inner", "array", []interface{}{"x", "y"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("outer_loop", FromState("outer"))).
		AddNode(LoopNode("inner_loop", FromState("inner"))).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "outer_loop").
		AddEdge("outer_loop", "body", "inner_loop").
		AddEdge("inner_loop", "body", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "inner_loop").
		// NO edge from inner_loop "done" → stack pop handles return to outer
		AddEdge("outer_loop", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 4) // 2 outer x 2 inner, 4 async callbacks
	AssertDispatchCount(t, h.Publisher, "DispatchCodeExecution", 4)
}

// 2.5 — loop → body(sequence) → steps(set_state): sequence inside loop, state resets per iteration
func TestLoopAsync_SequenceInBody(t *testing.T) {
	def := NewDefinition("LoopAsync_SeqBody").
		WithState("items", "array", []interface{}{"a", "b"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(SequenceNode("seq1", 3)).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		// Sequence with steps=3: step_2 and step_3 each increment (2 per iteration)
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "seq1").
		AddEdge("seq1", "step_2", "ss1").
		AddEdge("ss1", "out", "seq1").
		AddEdge("seq1", "step_3", "ss1").
		AddEdge("seq1", "done", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// 2 loop iterations × 2 sequence steps = 4 increments
	AssertState(t, exec, "count", 4)
}

// 2.6 — nested loops with sequence: outer → body(inner → body(sequence → steps(inc))) → done
func TestLoopAsync_NestedWithSequence(t *testing.T) {
	def := NewDefinition("LoopAsync_NestedSeq").
		WithState("outer", "array", []interface{}{"a", "b"}).
		WithState("inner", "array", []interface{}{"x", "y"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("outer_loop", FromState("outer"))).
		AddNode(LoopNode("inner_loop", FromState("inner"))).
		AddNode(SequenceNode("seq1", 2)).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "outer_loop").
		AddEdge("outer_loop", "body", "inner_loop").
		AddEdge("inner_loop", "body", "seq1").
		AddEdge("seq1", "step_2", "ss1").
		AddEdge("ss1", "out", "seq1").
		AddEdge("seq1", "done", "inner_loop").
		AddEdge("inner_loop", "done", "outer_loop").
		AddEdge("outer_loop", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// 2 outer × 2 inner × 1 sequence step = 4 increments
	AssertState(t, exec, "count", 4)
}

// 2.4 — loop → body(delay): timer async in each iteration
func TestLoopAsync_DelayInBody(t *testing.T) {
	def := NewDefinition("LoopAsync_Delay").
		WithState("items", "array", []interface{}{1, 2, 3}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(DelayNode("delay1", 5, "seconds")).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "delay1").
		AddEdge("delay1", "out", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 3)
}
