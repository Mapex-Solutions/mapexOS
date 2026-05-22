package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
	"workflow/src/modules/runtime/domain/entities"
	sharedTypes "workflow/src/shared/types"
)

func TestResilience_AsyncCallback_CorrectToken(t *testing.T) {
	def := NewDefinition("TokenCorrect").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"ok": true}, nil),
		},
	)

	// Token is echoed automatically by harness — should resume and complete
	AssertCompleted(t, exec)
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code1")
}

func TestResilience_AsyncCallback_WrongToken(t *testing.T) {
	def := NewDefinition("TokenWrong").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)

	// Trigger execution — it will suspend at code1
	h.RunSync(map[string]interface{}{})
	exec := h.StateRepo.GetLatest()
	if exec.Status != entities.ExecStatusWaiting {
		t.Fatalf("expected waiting, got %s", exec.Status)
	}

	// Send callback with WRONG token — should be rejected
	h.SendResumeWithToken(exec, "code1", "wrong-token-12345", sharedTypes.ResumeMessage{
		NodeID:       "code1",
		Status: "success",
		Output:       map[string]interface{}{"ok": true},
	})

	// Execution should still be waiting (wrong token rejected)
	exec = h.StateRepo.GetLatest()
	if exec.Status != entities.ExecStatusWaiting {
		t.Fatalf("expected execution to still be waiting after wrong token, got %s", exec.Status)
	}

	// Now send with correct token — read from NodeState
	token := ""
	if ns := exec.NodeStates["code1"]; ns != nil {
		if t, ok := ns["executionToken"].(string); ok {
			token = t
		}
	}
	if token == "" {
		t.Fatal("expected executionToken in NodeState")
	}

	h.SendResumeWithToken(exec, "code1", token, sharedTypes.ResumeMessage{
		NodeID:       "code1",
		Status: "success",
		Output:       map[string]interface{}{"ok": true},
	})

	exec = h.StateRepo.GetLatest()
	AssertCompleted(t, exec)
}

func TestResilience_AsyncCallback_EmptyToken(t *testing.T) {
	def := NewDefinition("TokenEmpty").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)

	// Trigger — suspends at code1
	h.RunSync(map[string]interface{}{})
	exec := h.StateRepo.GetLatest()
	if exec.Status != entities.ExecStatusWaiting {
		t.Fatalf("expected waiting, got %s", exec.Status)
	}

	// Send callback with EMPTY token — should be accepted (backward compat)
	h.SendResumeWithToken(exec, "code1", "", sharedTypes.ResumeMessage{
		NodeID:       "code1",
		Status: "success",
		Output:       map[string]interface{}{"ok": true},
	})

	exec = h.StateRepo.GetLatest()
	AssertCompleted(t, exec)
}

func TestResilience_TokenStoredInNodeState(t *testing.T) {
	def := NewDefinition("TokenInNodeState").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)

	// Trigger — suspends at code1
	h.RunSync(map[string]interface{}{})
	exec := h.StateRepo.GetLatest()
	if exec.Status != entities.ExecStatusWaiting {
		t.Fatalf("expected waiting, got %s", exec.Status)
	}

	// Verify token is stored in NodeState
	ns := exec.NodeStates["code1"]
	if ns == nil {
		t.Fatal("expected NodeState for code1")
	}
	token, ok := ns["executionToken"].(string)
	if !ok || token == "" {
		t.Fatal("expected executionToken in NodeState, got empty or missing")
	}
	if len(token) != 32 {
		t.Fatalf("expected 32-char hex token, got %d chars: %s", len(token), token)
	}

	t.Logf("Token stored in NodeState: %s", token)
}

func TestResilience_DispatchFailure_NoCheckpoint(t *testing.T) {
	def := NewDefinition("DispatchFail").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "end").
		Build()

	h := NewHarness(t, def)

	// Make dispatch fail
	h.Publisher.FailOnMethod("DispatchCodeExecution")

	// Trigger with tolerance for Nack (dispatch failure causes Nack, not Fatal)
	h.RunSyncAllowNack(map[string]interface{}{})
	exec := h.StateRepo.GetLatest()

	// Execution should NOT be in waiting state (dispatch failed → no checkpoint for waiting)
	if exec.Status == entities.ExecStatusWaiting {
		t.Fatal("expected execution to NOT be waiting after dispatch failure — checkpoint should not have been written")
	}

	t.Logf("Execution status after dispatch failure: %s", exec.Status)
}

func TestResilience_Fanout4Async_ConcurrentCallbacks(t *testing.T) {
	def := NewDefinition("Fanout4Async").
		WithState("results", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 4, "")).
		AddNode(CodeNode("code_a", "return { v: 'a' }", 5000)).
		AddNode(CodeNode("code_b", "return { v: 'b' }", 5000)).
		AddNode(CodeNode("code_c", "return { v: 'c' }", 5000)).
		AddNode(CodeNode("code_d", "return { v: 'd' }", 5000)).
		AddNode(MergeNode("merge1", 4)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code_a").
		AddEdge("fan1", "out_2", "code_b").
		AddEdge("fan1", "out_3", "code_c").
		AddEdge("fan1", "out_4", "code_d").
		AddEdge("code_a", "success", "merge1").
		AddEdge("code_b", "success", "merge1").
		AddEdge("code_c", "success", "merge1").
		AddEdge("code_d", "success", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code_a": CodeSuccessCallback(map[string]interface{}{"v": "a"}, nil),
			"code_b": CodeSuccessCallback(map[string]interface{}{"v": "b"}, nil),
			"code_c": CodeSuccessCallback(map[string]interface{}{"v": "c"}, nil),
			"code_d": CodeSuccessCallback(map[string]interface{}{"v": "d"}, nil),
		},
	)

	AssertCompleted(t, exec)

	// Verify all 4 outputs merged
	for _, nodeID := range []string{"code_a", "code_b", "code_c", "code_d"} {
		if exec.NodeOutputs[nodeID] == nil {
			t.Fatalf("expected output for %s, got nil", nodeID)
		}
	}

	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code_a")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code_b")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code_c")
	AssertDispatched(t, h.Publisher, "DispatchCodeExecution", "code_d")
}
