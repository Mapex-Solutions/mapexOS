package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

/**
 * Operation: set
 */

func TestSetState_SetLiteral(t *testing.T) {
	def := NewDefinition("SetLiteral").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "target", Literal("hello"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "target", "hello")
}

func TestSetState_SetFromState(t *testing.T) {
	def := NewDefinition("SetFromState").
		WithState("source", "string", "from_state_value").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "target", FromState("source"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "target", "from_state_value")
}

func TestSetState_SetFromEvent(t *testing.T) {
	def := NewDefinition("SetFromEvent").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "target", FromEvent("data.value"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{
		"data": map[string]interface{}{
			"value": "from_event_value",
		},
	})

	AssertCompleted(t, exec)
	AssertState(t, exec, "target", "from_event_value")
}

func TestSetState_SetFromInput(t *testing.T) {
	def := NewDefinition("SetFromInput").
		WithExternalInput("sensorId", "sensor_123").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "target", FromInput("sensorId"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "target", "sensor_123")
}

func TestSetState_SetFromNodeOutput(t *testing.T) {
	def := NewDefinition("SetFromNodeOutput").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(CodeNode("code1", "return { text: 'from_code' }", 5000)).
		AddNode(SetStateNode("ss1", "set", "target", FromNodeOutput("code1", "text"))).
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
			"code1": CodeSuccessCallback(
				map[string]interface{}{"text": "from_code"},
				nil,
			),
		},
	)

	AssertCompleted(t, exec)
	AssertState(t, exec, "target", "from_code")
}

/**
 * Operation: increment
 */

func TestSetState_IncrementNumeric(t *testing.T) {
	def := NewDefinition("IncrementNumeric").
		WithState("counter", "number", 5).
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "increment", "counter", Literal("3"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "counter", 8)
}

func TestSetState_IncrementNonNumeric(t *testing.T) {
	def := NewDefinition("IncrementNonNumeric").
		WithState("counter", "string", "abc").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "increment", "counter", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// "abc" coerces to 0, so 0 + 1 = 1
	AssertState(t, exec, "counter", 1)
}

/**
 * Operation: decrement
 */

func TestSetState_DecrementNumeric(t *testing.T) {
	def := NewDefinition("DecrementNumeric").
		WithState("counter", "number", 10).
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "decrement", "counter", Literal("3"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "counter", 7)
}

func TestSetState_DecrementNonNumeric(t *testing.T) {
	def := NewDefinition("DecrementNonNumeric").
		WithState("counter", "string", "abc").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "decrement", "counter", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// "abc" coerces to 0, so 0 - 1 = -1
	AssertState(t, exec, "counter", -1)
}

/**
 * Operation: append
 */

func TestSetState_AppendToArray(t *testing.T) {
	def := NewDefinition("AppendToArray").
		WithState("items", "array", []interface{}{1, 2}).
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "append", "items", Literal("3"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "items", []interface{}{float64(1), float64(2), "3"})
}

func TestSetState_AppendToNonExistent(t *testing.T) {
	def := NewDefinition("AppendNonExistent").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "append", "items", Literal("1"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	// nil coerces to empty array, so result is ["1"]
	AssertState(t, exec, "items", []interface{}{"1"})
}

/**
 * Operation: remove
 */

func TestSetState_RemoveExistingField(t *testing.T) {
	def := NewDefinition("RemoveField").
		WithState("counter", "number", 5).
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "remove", "counter", nil)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertStateNotExists(t, exec, "counter")
}

func TestSetState_RemoveNonExistentField(t *testing.T) {
	def := NewDefinition("RemoveNonExistent").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "remove", "doesNotExist", nil)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	// No error, no-op
	AssertCompleted(t, exec)
	AssertStateNotExists(t, exec, "doesNotExist")
}

/**
 * Error cases
 */

func TestSetState_UnknownOperation(t *testing.T) {
	def := NewDefinition("UnknownOp").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "multiply", "target", Literal("2"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "EXECUTION_ERROR")
}

func TestSetState_ResolutionFailure(t *testing.T) {
	def := NewDefinition("ResolveFail").
		WithState("target", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(SetStateNode("ss1", "set", "target", FromNodeOutput("nonexistent_node", "field"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertFailed(t, exec, "EXECUTION_ERROR")
}
