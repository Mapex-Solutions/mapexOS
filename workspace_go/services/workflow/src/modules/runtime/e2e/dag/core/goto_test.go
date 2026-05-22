package core

import (
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
)

func TestGoto_SenderToReceiver(t *testing.T) {
	def := NewDefinition("GotoBasic").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(GotoSender("sender1", "portal_X")).
		AddNode(GotoReceiver("receiver1", "portal_X")).
		AddNode(SetStateNode("ss1", "set", "result", Literal("teleported"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "sender1").
		AddEdge("receiver1", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "teleported")
	AssertPathContains(t, exec, "sender1")
	AssertPathContains(t, exec, "receiver1")
}

func TestGoto_ReceiverPassthrough(t *testing.T) {
	def := NewDefinition("GotoReceiver").
		WithState("result", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(GotoReceiver("receiver1", "portal_X")).
		AddNode(SetStateNode("ss1", "set", "result", Literal("passed"))).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "receiver1").
		AddEdge("receiver1", "out", "ss1").
		AddEdge("ss1", "out", "end").
		Build()

	h := NewHarness(t, def)
	exec := h.RunSync(map[string]interface{}{})

	AssertCompleted(t, exec)
	AssertState(t, exec, "result", "passed")
}
