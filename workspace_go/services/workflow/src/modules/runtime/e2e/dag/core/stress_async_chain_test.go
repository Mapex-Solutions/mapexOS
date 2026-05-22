package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

// 1.1 — code → set_state(fromNodeOutput): nodeOutput propagates after async resume
func TestAsyncChain_CodeToSetStateNodeOutput(t *testing.T) {
	def := NewDefinition("AsyncChain_CodeToState").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { text: 'from_code' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("code1", "text"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"text": "from_code"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "from_code")
}

// 1.2 — code → condition(nodeOutput == X): condition evaluates nodeOutput after resume
func TestAsyncChain_CodeToConditionNodeOutput(t *testing.T) {
	def := NewDefinition("AsyncChain_CodeToCond").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { status: 'ok' }", 5000)).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromNodeOutput("code1", "status"), Operator: "equals", Value: Literal("ok")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("true_path"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("false_path"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "cond1").
		AddEdge("code1", "error", "end").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"status": "ok"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "true_path")
}

// 1.3 — code1 → code2: two async nodes chained (resume → resume)
func TestAsyncChain_CodeToCode(t *testing.T) {
	def := NewDefinition("AsyncChain_CodeCode").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { step: '1' }", 5000)).
		AddNode(CodeNode("code2", "return { step: '2' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("code2", "step"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "code2").
		AddEdge("code1", "error", "end").
		AddEdge("code2", "success", "ss1").
		AddEdge("code2", "error", "end").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"step": "1"}, nil),
			"code2": CodeSuccessCallback(map[string]interface{}{"step": "2"}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "2")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code1")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code2")
}

// 1.4 — delay → code: timer suspend then async callback in sequence
func TestAsyncChain_DelayToCode(t *testing.T) {
	def := NewDefinition("AsyncChain_DelayCode").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(DelayNode("delay1", 5, "seconds")).
		AddNode(CodeNode("code1", "return { done: true }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", Literal("completed"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "delay1").
		AddEdge("delay1", "out", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"delay1": SuccessCallback(nil),
			"code1":  CodeSuccessCallback(map[string]interface{}{"done": true}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "completed")
}

// 1.5 — wait_signal → set_state → condition: signal resume drives state then routing
func TestAsyncChain_SignalToStateToCondition(t *testing.T) {
	def := NewDefinition("AsyncChain_SignalStateCond").
		WithState("approved", "string", "no").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(WaitSignalNode("ws1", "approval")).
		AddNode(SetStateNode("ss1", "set", "approved", Literal("yes"))).
		AddNode(ConditionNode("cond1", "AND", []ConditionItem{
			{Field: FromState("approved"), Operator: "equals", Value: Literal("yes")},
		})).
		AddNode(SetStateNode("ss_t", "set", "path", Literal("approved"))).
		AddNode(SetStateNode("ss_f", "set", "path", Literal("rejected"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ws1").
		AddEdge("ws1", "out", "ss1").
		AddEdge("ss1", "out", "cond1").
		AddEdge("cond1", "true", "ss_t").
		AddEdge("cond1", "false", "ss_f").
		AddEdge("ss_t", "out", "end").
		AddEdge("ss_f", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"ws1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "approved", "yes")
	AssertState(t, exec, "path", "approved")
}

// 1.5 — wait_signal with data → signal data appears in NodeOutputs
func TestAsyncChain_SignalDataAppearsInNodeOutputs(t *testing.T) {
	def := NewDefinition("AsyncChain_SignalData").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(WaitSignalNode("ws1", "sensor_data")).
		AddNode(SetStateNode("ss1", "set", "result", FromNodeOutput("ws1", "temperature"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ws1").
		AddEdge("ws1", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"ws1": SignalCallback(map[string]interface{}{"temperature": 25.3, "humidity": 60}),
		},
	)

	AssertCompleted(t, exec)
	AssertNodeOutput(t, exec, "ws1", map[string]interface{}{"temperature": 25.3, "humidity": float64(60)})
	AssertState(t, exec, "result", "25.3")
}

// 1.6 — wait_signal without data → no empty NodeOutput entry
func TestAsyncChain_SignalWithoutDataNoNodeOutput(t *testing.T) {
	def := NewDefinition("AsyncChain_SignalNoData").
		AddNode(StartNode("__start__")).
		AddNode(WaitSignalNode("ws1", "approval")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ws1").
		AddEdge("ws1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"ws1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
	// No signal data → no NodeOutput entry for ws1
	if _, exists := exec.NodeOutputs["ws1"]; exists {
		t.Fatalf("expected no NodeOutput for ws1 when signal has no data, got %v", exec.NodeOutputs["ws1"])
	}
}
