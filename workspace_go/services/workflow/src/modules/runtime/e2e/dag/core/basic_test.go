package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestMinimal_Start_End(t *testing.T) {
	def := NewDefinition("Minimal").
		AddNode(StartNode("__start__")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "__start__")
	AssertPathContains(t, exec, "end")
	AssertPathLength(t, exec, 2)
}

func TestStart_SetState_End(t *testing.T) {
	def := NewDefinition("SetState").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "count", Literal("42"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", "42")
}

func TestStart_Log_End(t *testing.T) {
	def := NewDefinition("Log").
		AddNode(StartNode("__start__")).
		AddNode(LogNode("log1", "hello world", "info")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "log1")
}

func TestEnd_TerminateWithError(t *testing.T) {
	def := NewDefinition("EndError").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithError("end", "CUSTOM_ERROR", "something broke")).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "CUSTOM_ERROR")
}

func TestCode_Success(t *testing.T) {
	def := NewDefinition("CodeSuccess").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { text: 'ok' }", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(
				map[string]interface{}{"text": "ok"},
				nil,
			),
		},
	)

	AssertCompleted(t, exec)
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code1")
}

func TestCode_Error(t *testing.T) {
	def := NewDefinition("CodeError").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(EndNode("end_ok")).
		AddNode(EndNode("end_err")).
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

	// Error follows the "error" handle → end_err → completes (user handled the error)
	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "end_err")
}
