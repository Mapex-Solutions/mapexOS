package core

import (
	"testing"

	defEntities "workflow/src/modules/definitions/domain/entities"
	. "workflow/src/modules/runtime/e2e/testutil"
	sharedTypes "workflow/src/shared/types"
)

/**
 * Error Handler: error output handle (no retry)
 */

// Code error with "error" handle connected → follows error handle, completes
func TestErrorHandler_CodeErrorFollowsErrorHandle(t *testing.T) {
	def := NewDefinition("EH_ErrorHandle").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(SetStateNode("ss_ok", "set", "path", Literal("success"))).
		AddNode(SetStateNode("ss_err", "set", "path", Literal("error_handled"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "ss_ok").
		AddEdge("code1", "error", "ss_err").
		AddEdge("ss_ok", "out", "end").
		AddEdge("ss_err", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("SCRIPT_ERROR", "fail"),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "path", "error_handled")
}

// Code error WITHOUT "error" handle → failExecution (no handle to follow)
func TestErrorHandler_NoErrorHandleFailsExecution(t *testing.T) {
	def := NewDefinition("EH_NoHandle").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		// NO error edge — error has nowhere to go
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("SCRIPT_ERROR", "fail"),
		},
	)

	AssertFailed(t, exec, "SCRIPT_ERROR")
}

/**
 * Error Handler: retry with backoff
 */

// Code fails once, retry succeeds on second attempt → completes via success handle
func TestErrorHandler_RetrySucceedsOnSecondAttempt(t *testing.T) {
	def := NewDefinition("EH_RetrySuccess").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { ok: true }", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", Literal("success"))).
		AddNode(SetStateNode("ss_err", "set", "result", Literal("error"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "ss_err").
		AddEdge("ss1", "out", "end").
		AddEdge("ss_err", "out", "end").
		Build()

	// Enable retry on code1
	for i := range def.Nodes {
		if def.Nodes[i].ID == "code1" {
			def.Nodes[i].ErrorHandler = &defEntities.ErrorHandlerConfig{
				Enabled:           true,
				MaxAttempts:       3,
				InitialInterval:   1,
				IntervalUnit:      "seconds",
				BackoffMultiplier: 1.0,
			}
		}
	}

	callCount := 0
	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": func(nodeID string, ns map[string]interface{}) sharedTypes.ResumeMessage {
				callCount++
				if callCount == 1 {
					// First attempt fails
					return CodeErrorCallback("TRANSIENT", "temporary failure")(nodeID, ns)
				}
				// Second attempt succeeds
				return CodeSuccessCallback(map[string]interface{}{"ok": true}, nil)(nodeID, ns)
			},
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "success")
}

// Code fails all retry attempts → follows error handle
func TestErrorHandler_RetryExhaustedFollowsErrorHandle(t *testing.T) {
	def := NewDefinition("EH_RetryExhausted").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(SetStateNode("ss_ok", "set", "result", Literal("success"))).
		AddNode(SetStateNode("ss_err", "set", "result", Literal("all_retries_failed"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "ss_ok").
		AddEdge("code1", "error", "ss_err").
		AddEdge("ss_ok", "out", "end").
		AddEdge("ss_err", "out", "end").
		Build()

	// Enable retry with max 2 attempts
	for i := range def.Nodes {
		if def.Nodes[i].ID == "code1" {
			def.Nodes[i].ErrorHandler = &defEntities.ErrorHandlerConfig{
				Enabled:           true,
				MaxAttempts:       2,
				InitialInterval:   1,
				IntervalUnit:      "seconds",
				BackoffMultiplier: 1.0,
			}
		}
	}

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("PERSISTENT_ERROR", "always fails"),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "all_retries_failed")
}

// Retry disabled (enabled=false) → error goes directly to error handle
func TestErrorHandler_RetryDisabledGoesDirectToErrorHandle(t *testing.T) {
	def := NewDefinition("EH_RetryDisabled").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(SetStateNode("ss_ok", "set", "result", Literal("success"))).
		AddNode(SetStateNode("ss_err", "set", "result", Literal("direct_error"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "ss_ok").
		AddEdge("code1", "error", "ss_err").
		AddEdge("ss_ok", "out", "end").
		AddEdge("ss_err", "out", "end").
		Build()

	// ErrorHandler present but disabled
	for i := range def.Nodes {
		if def.Nodes[i].ID == "code1" {
			def.Nodes[i].ErrorHandler = &defEntities.ErrorHandlerConfig{
				Enabled:           false,
				MaxAttempts:       3,
				InitialInterval:   1,
				IntervalUnit:      "seconds",
				BackoffMultiplier: 2.0,
			}
		}
	}

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("SCRIPT_ERROR", "fail"),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "direct_error")
}

/**
 * PathEntry verification on error/timeout/retry
 */

// PathEntry shows status="error" when code node fails with error handle
func TestErrorHandler_PathEntryShowsErrorStatus(t *testing.T) {
	def := NewDefinition("EH_PathError").
		WithState("path", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(SetStateNode("ss_err", "set", "path", Literal("handled"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		AddEdge("code1", "error", "ss_err").
		AddEdge("ss_err", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("SCRIPT_ERROR", "something broke"),
		},
	)

	AssertCompleted(t, exec)
	AssertPathNodeStatus(t, exec, "code1", "error")
	AssertPathNodeError(t, exec, "code1", "SCRIPT_ERROR")
}

// PathEntry shows status="error" when code node fails without error handle (failExecution)
func TestErrorHandler_PathEntryShowsErrorOnFail(t *testing.T) {
	def := NewDefinition("EH_PathFail").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "throw new Error('fail')", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeErrorCallback("FATAL", "unrecoverable"),
		},
	)

	AssertFailed(t, exec, "FATAL")
	AssertPathNodeStatus(t, exec, "code1", "error")
	AssertPathNodeError(t, exec, "code1", "FATAL")
}

// PathEntry shows error before retry, then retrying on retry timer
func TestErrorHandler_PathEntryShowsErrorBeforeRetry(t *testing.T) {
	def := NewDefinition("EH_PathRetry").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(SetStateNode("ss1", "set", "result", Literal("ok"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "ss1").
		AddEdge("code1", "error", "end").
		AddEdge("ss1", "out", "end").
		Build()

	for i := range def.Nodes {
		if def.Nodes[i].ID == "code1" {
			def.Nodes[i].ErrorHandler = &defEntities.ErrorHandlerConfig{
				Enabled:           true,
				MaxAttempts:       2,
				InitialInterval:   1,
				IntervalUnit:      "seconds",
				BackoffMultiplier: 1.0,
			}
		}
	}

	callCount := 0
	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": func(nodeID string, ns map[string]interface{}) sharedTypes.ResumeMessage {
				callCount++
				if callCount == 1 {
					return CodeErrorCallback("TRANSIENT", "temp")(nodeID, ns)
				}
				return CodeSuccessCallback(map[string]interface{}{"ok": true}, nil)(nodeID, ns)
			},
		},
	)

	AssertCompleted(t, exec)
	// The first code1 entry should have status "error" (before retry)
	found := false
	for _, entry := range exec.ExecutionPath {
		if entry.NodeID == "code1" && entry.Status == "error" {
			found = true
			if entry.Error == nil {
				t.Fatal("expected error field to be populated on error PathEntry")
			}
			break
		}
	}
	if !found {
		t.Fatal("expected at least one code1 PathEntry with status 'error'")
	}
}
