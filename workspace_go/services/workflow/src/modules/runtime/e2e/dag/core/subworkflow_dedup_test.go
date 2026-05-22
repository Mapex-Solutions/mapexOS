package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
	"workflow/src/modules/runtime/domain/entities"
)

func TestSubworkflow_CompleteAndParentResumes(t *testing.T) {
	def := NewDefinition("SubwfParentResume").
		AddNode(StartNode("__start__")).
		AddNode(SubworkflowNode("subwf_1", "child-def-id")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "subwf_1").
		AddEdge("subwf_1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"subwf_1": SubworkflowSuccessCallback(map[string]interface{}{"result": "ok"}),
		},
	)

	AssertCompleted(t, exec)
	AssertDispatched(t, h.Publisher, "DispatchSubworkflowExecution", "subwf_1")
}

func TestSubworkflow_TokenStoredInNodeState(t *testing.T) {
	def := NewDefinition("SubwfToken").
		AddNode(StartNode("__start__")).
		AddNode(SubworkflowNode("subwf_1", "child-def-id")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "subwf_1").
		AddEdge("subwf_1", "out", "end").
		Build()

	h := NewHarness(t, def)
	h.RunSync(map[string]interface{}{})
	exec := h.StateRepo.GetLatest()

	if exec.Status != entities.ExecStatusWaiting {
		t.Fatalf("expected waiting, got %s", exec.Status)
	}

	ns := exec.NodeStates["subwf_1"]
	if ns == nil {
		t.Fatal("expected NodeState for subwf_1")
	}
	token, ok := ns["executionToken"].(string)
	if !ok || token == "" {
		t.Fatal("expected executionToken in NodeState")
	}
	if len(token) != 32 {
		t.Fatalf("expected 32-char token, got %d: %s", len(token), token)
	}

	t.Logf("Subworkflow execution token: %s", token)
}

func TestSubworkflow_DispatchCaptured(t *testing.T) {
	def := NewDefinition("SubwfDispatch").
		AddNode(StartNode("__start__")).
		AddNode(SubworkflowNode("subwf_1", "child-def-id")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "subwf_1").
		AddEdge("subwf_1", "out", "end").
		Build()

	h := NewHarness(t, def)
	h.RunSync(map[string]interface{}{})

	AssertDispatched(t, h.Publisher, "DispatchSubworkflowExecution", "subwf_1")
	AssertDispatchCount(t, h.Publisher, "DispatchSubworkflowExecution", 1)
}

func TestSubworkflow_ErrorAndParentResumes(t *testing.T) {
	def := NewDefinition("SubwfError").
		AddNode(StartNode("__start__")).
		AddNode(SubworkflowNode("subwf_1", "child-def-id")).
		AddNode(EndNode("end_ok")).
		AddNode(EndNode("end_err")).
		AddEdge("__start__", "out", "subwf_1").
		AddEdge("subwf_1", "out", "end_ok").
		AddEdge("subwf_1", "error", "end_err").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"subwf_1": ErrorCallback("CHILD_FAILED", "child execution failed"),
		},
	)

	// Error follows the "error" handle → end_err → completes
	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "end_err")
}
