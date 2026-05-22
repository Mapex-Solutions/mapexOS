package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestLog_StaticMessage(t *testing.T) {
	def := NewDefinition("LogStatic").
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

func TestLog_StateToken(t *testing.T) {
	def := NewDefinition("LogStateToken").
		WithState("counter", "number", 42).
		AddNode(StartNode("__start__")).
		AddNode(LogNode("log1", "count: ${state.counter}", "info")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "log1")
}

func TestLog_EventToken(t *testing.T) {
	def := NewDefinition("LogEventToken").
		AddNode(StartNode("__start__")).
		AddNode(LogNode("log1", "type: ${event.eventType}", "warn")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{
		"eventType": "temperature_alert",
	})

	AssertCompleted(t, exec)
	AssertPathContains(t, exec, "log1")
}

func TestLog_MultipleTokens(t *testing.T) {
	def := NewDefinition("LogMultiTokens").
		WithState("a", "string", "hello").
		AddNode(StartNode("__start__")).
		AddNode(LogNode("log1", "${state.a} and ${event.b}", "debug")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{"b": "world"})

	AssertCompleted(t, exec)
}

func TestLog_MissingToken(t *testing.T) {
	def := NewDefinition("LogMissingToken").
		AddNode(StartNode("__start__")).
		AddNode(LogNode("log1", "${state.nonexistent}", "info")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	// Token not substituted — no error, just passthrough
	AssertCompleted(t, exec)
}

func TestLog_EmptyMessage(t *testing.T) {
	def := NewDefinition("LogEmpty").
		AddNode(StartNode("__start__")).
		AddNode(LogNode("log1", "", "info")).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "log1").
		AddEdge("log1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
}

func TestLog_AllLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			def := NewDefinition("LogLevel_" + level).
				AddNode(StartNode("__start__")).
				AddNode(LogNode("log1", "test", level)).
				AddNode(EndNode("end")).
				AddEdge("__start__", "out", "log1").
				AddEdge("log1", "out", "end").
				Build()

			h := NewHarness(t, def)
			exec := h.RunSync(map[string]interface{}{})

			AssertCompleted(t, exec)
		})
	}
}
