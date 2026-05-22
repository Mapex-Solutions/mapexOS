package core

import (
	"testing"

	defEntities "workflow/src/modules/definitions/domain/entities"
	. "workflow/src/modules/runtime/e2e/testutil"
	sharedTypes "workflow/src/shared/types"
)

// TestSchedule_DelayPublishesSchedule verifies that a delay node publishes
// a NATS schedule with waitType="timer".
func TestSchedule_DelayPublishesSchedule(t *testing.T) {
	def := NewDefinition("Sched_Delay").
		AddNode(StartNode("__start__")).
		AddNode(DelayNode("delay1", 5, "seconds")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "delay1").
		AddEdge("delay1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"delay1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
	AssertSchedulePublished(t, h.Publisher, "delay1")

	// Verify waitType is "timer" (delay, not timeout)
	event := h.Publisher.FindDispatch("PublishSchedule", "delay1")
	if event == nil {
		t.Fatal("PublishSchedule event not found for delay1")
	}
	if event.Status != "timer" {
		t.Errorf("expected waitType=timer, got %q", event.Status)
	}
}

// TestSchedule_CodeTimeoutPublishesSchedule verifies that a code node publishes
// a NATS schedule with the correct enableOutput flag.
func TestSchedule_CodeTimeoutPublishesSchedule(t *testing.T) {
	def := NewDefinition("Sched_Code").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { ok: true }", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"ok": true}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertSchedulePublished(t, h.Publisher, "code1")

	// Verify waitType is "callback" (timeout safety net)
	event := h.Publisher.FindDispatch("PublishSchedule", "code1")
	if event == nil {
		t.Fatal("PublishSchedule event not found for code1")
	}
	if event.Status != "callback" {
		t.Errorf("expected waitType=callback, got %q", event.Status)
	}
}

// TestSchedule_CallbackPurgesSchedule verifies that when a callback arrives
// before timeout, the pending schedule is purged.
func TestSchedule_CallbackPurgesSchedule(t *testing.T) {
	def := NewDefinition("Sched_Purge").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { ok: true }", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"ok": true}, nil),
		},
	)

	AssertCompleted(t, exec)
	// Schedule was published on suspend
	AssertSchedulePublished(t, h.Publisher, "code1")
	// Schedule was purged when callback arrived (not a timeout)
	AssertSchedulePurged(t, h.Publisher, "code1")
}

// TestSchedule_SignalPurgesSchedule verifies that when a signal arrives
// before timeout, the pending schedule is purged.
func TestSchedule_SignalPurgesSchedule(t *testing.T) {
	def := NewDefinition("Sched_Signal").
		AddNode(StartNode("__start__")).
		AddNode(WaitSignalNode("signal1", "approval")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "signal1").
		AddEdge("signal1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"signal1": SuccessCallback(nil),
		},
	)

	AssertCompleted(t, exec)
	AssertSchedulePublished(t, h.Publisher, "signal1")
	// Signal callback purges schedule
	AssertSchedulePurged(t, h.Publisher, "signal1")
}

// TestSchedule_FanoutMultipleSchedules verifies that fanout with multiple
// async branches publishes one schedule per waiting node.
func TestSchedule_FanoutMultipleSchedules(t *testing.T) {
	def := NewDefinition("Sched_Fanout").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 2, "")).
		AddNode(CodeNode("code_a", "return {}", 5000)).
		AddNode(CodeNode("code_b", "return {}", 5000)).
		AddNode(MergeNode("merge1", 2)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "code_a").
		AddEdge("fan1", "out_2", "code_b").
		AddEdge("code_a", "success", "merge1").
		AddEdge("code_b", "success", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code_a": CodeSuccessCallback(map[string]interface{}{"a": true}, nil),
			"code_b": CodeSuccessCallback(map[string]interface{}{"b": true}, nil),
		},
	)

	AssertCompleted(t, exec)
	// Both async nodes should have published schedules
	AssertSchedulePublished(t, h.Publisher, "code_a")
	AssertSchedulePublished(t, h.Publisher, "code_b")
	AssertScheduleCount(t, h.Publisher, 2)
}

// TestSchedule_RetryTimerPublishesSchedule verifies that when a node fails
// with retry enabled, a schedule is published for the retry timer.
func TestSchedule_RetryTimerPublishesSchedule(t *testing.T) {
	def := NewDefinition("Sched_Retry").
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
					return CodeErrorCallback("TRANSIENT", "temporary failure")(nodeID, ns)
				}
				return CodeSuccessCallback(map[string]interface{}{"ok": true}, nil)(nodeID, ns)
			},
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "success")

	// First call: code node publishes schedule (callback waitType)
	// Then error → retry → publishes retry timer schedule (retryTimer waitType)
	// Then retry fires → code succeeds → schedule published again (callback waitType)
	// Verify retryTimer was published
	found := false
	for _, e := range h.Publisher.Events {
		if e.Method == "PublishSchedule" && e.Status == "retryTimer" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected PublishSchedule with waitType=retryTimer for retry timer")
	}
}

// TestSchedule_CompletePurgesAll verifies that completing a workflow
// calls PurgeAllSchedules for cleanup.
func TestSchedule_CompletePurgesAll(t *testing.T) {
	def := NewDefinition("Sched_Complete").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return {}", 5000)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "code1").
		AddEdge("code1", "success", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunWithCallbacks(
		map[string]interface{}{},
		map[string]CallbackFunc{
			"code1": CodeSuccessCallback(map[string]interface{}{"ok": true}, nil),
		},
	)

	AssertCompleted(t, exec)
	AssertPurgeAllSchedules(t, h.Publisher)
}
