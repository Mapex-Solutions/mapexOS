package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestLoop_ThreeItems(t *testing.T) {
	def := NewDefinition("Loop3").
		WithState("items", "array", []interface{}{"a", "b", "c"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 3)
}

func TestLoop_EmptyArray(t *testing.T) {
	def := NewDefinition("LoopEmpty").
		WithState("items", "array", []interface{}{}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 0)
}

func TestLoop_SingleItem(t *testing.T) {
	def := NewDefinition("LoopSingle").
		WithState("items", "array", []interface{}{"only"}).
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 1)
}

func TestLoop_SetsLoopItem(t *testing.T) {
	def := NewDefinition("LoopItem").
		WithState("items", "array", []interface{}{"hello", "world"}).
		WithState("last_item", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(SetStateNode("ss1", "set", "last_item", FromState("loop_item"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// Last iteration sets loop_item to "world"
	AssertState(t, exec, "last_item", "world")
}

func TestLoop_SourceFromEvent(t *testing.T) {
	def := NewDefinition("LoopEvent").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromEvent("data.list"))).
		AddNode(SetStateNode("ss1", "increment", "count", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "ss1").
		AddEdge("ss1", "out", "loop1").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{
		"data": map[string]interface{}{
			"list": []interface{}{1, 2, 3, 4},
		},
	})

	AssertCompleted(t, exec)
	AssertState(t, exec, "count", 4)
}

func TestLoop_NonArraySource_Error(t *testing.T) {
	def := NewDefinition("LoopNonArray").
		WithState("items", "string", "not_an_array").
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "end").
		AddEdge("loop1", "done", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "EXECUTION_ERROR")
}

// Error handle is a recovery mechanism — when connected, the workflow completes via the error path.
func TestLoop_ErrorHandle_ToEndNode_Completes(t *testing.T) {
	def := NewDefinition("LoopErrorHandle").
		WithState("items", "string", "not_an_array").
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(EndNode("end_ok")).
		AddNode(EndNode("end_err")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "end_ok").
		AddEdge("loop1", "done", "end_ok").
		AddEdge("loop1", "error", "end_err").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	// Error handle catches the error → workflow completes (recovery pattern)
	AssertCompleted(t, exec)

	// Path entry for the loop node should show error status
	found := false
	for _, pe := range exec.ExecutionPath {
		if pe.NodeID == "loop1" && pe.Status == "error" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected loop1 path entry with status=error")
	}
}

// Error handle → End with terminateWithError=true → workflow fails with the End node's error.
func TestLoop_ErrorHandle_WithErrorEnd_FailsExecution(t *testing.T) {
	def := NewDefinition("LoopErrorHandleWithErrorEnd").
		WithState("items", "string", "not_an_array").
		AddNode(StartNode("__start__")).
		AddNode(LoopNode("loop1", FromState("items"))).
		AddNode(EndNode("end_ok")).
		AddNode(EndNodeWithError("end_err", "CUSTOM_ERROR", "custom message")).
		AddEdge("__start__", "out", "loop1").
		AddEdge("loop1", "body", "end_ok").
		AddEdge("loop1", "done", "end_ok").
		AddEdge("loop1", "error", "end_err").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	// Error handle catches → routes to End with terminateWithError → fails with the End node's error
	AssertFailed(t, exec, "CUSTOM_ERROR")
}
