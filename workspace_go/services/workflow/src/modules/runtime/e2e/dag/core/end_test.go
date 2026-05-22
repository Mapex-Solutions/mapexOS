package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestEnd_NormalCompletion(t *testing.T) {
	def := NewDefinition("EndNormal").
		AddNode(StartNode("__start__")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "end")
	AssertPathNodeStatus(t, exec, "end", "completed")
}

func TestEnd_ErrorLiteral(t *testing.T) {
	def := NewDefinition("EndErrorLiteral").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithError("end", "E001", "something broke")).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailedWithMessage(t, exec, "E001", "something broke")
}

func TestEnd_ErrorMessageFromState(t *testing.T) {
	def := NewDefinition("EndErrorState").
		WithState("errorMsg", "string", "state error happened").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithErrorSource("end", "E002", FromState("errorMsg"))).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailedWithMessage(t, exec, "E002", "state error happened")
}

func TestEnd_ErrorMessageFromEvent(t *testing.T) {
	def := NewDefinition("EndErrorEvent").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithErrorSource("end", "E003", FromEvent("data.error"))).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{
		"data": map[string]interface{}{
			"error": "event error happened",
		},
	})

	AssertFailedWithMessage(t, exec, "E003", "event error happened")
}

func TestEnd_ErrorMessageFromInput(t *testing.T) {
	def := NewDefinition("EndErrorInput").
		WithExternalInput("reason", "input error happened").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithErrorSource("end", "E004", FromInput("reason"))).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailedWithMessage(t, exec, "E004", "input error happened")
}

func TestEnd_ErrorMessageFromNodeOutput(t *testing.T) {
	def := NewDefinition("EndErrorNodeOutput").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { result: 'node output error' }", 5000)).
		AddNode(EndNodeWithErrorSource("end", "E005", FromNodeOutput("code1", "result"))).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(
				map[string]interface{}{"result": "node output error"},
				nil,
			),
		},
	)

	AssertFailedWithMessage(t, exec, "E005", "node output error")
}

func TestEnd_ErrorResolutionFailure_FallbackToCode(t *testing.T) {
	def := NewDefinition("EndErrorFallback").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithErrorSource("end", "FALLBACK_CODE", FromState("nonexistent"))).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	// When resolution fails, errorMessage falls back to errorCode
	AssertFailedWithMessage(t, exec, "FALLBACK_CODE", "FALLBACK_CODE")
}

func TestEnd_ErrorHasMetadata(t *testing.T) {
	def := NewDefinition("EndErrorMeta").
		AddNode(StartNode("__start__")).
		AddNode(EndNodeWithError("end", "META_ERR", "test metadata")).
		AddEdge("__start__", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "META_ERR")
	AssertErrorHasMetadata(t, exec)
	if exec.ErrorInfo.NodeID != "end" {
		t.Fatalf("expected errorInfo.nodeId=%q, got %q", "end", exec.ErrorInfo.NodeID)
	}
	if exec.ErrorInfo.NodeType != "core/end" {
		t.Fatalf("expected errorInfo.nodeType=%q, got %q", "core/end", exec.ErrorInfo.NodeType)
	}
}
