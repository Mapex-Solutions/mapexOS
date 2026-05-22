package core

import (
	"fmt"
	"testing"

	. "workflow/src/modules/runtime/e2e/testutil"
	"workflow/src/modules/runtime/domain/entities"
)

func TestShutdown_MidChain_20Nodes(t *testing.T) {
	b := NewDefinition("ShutdownChain20").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__"))

	for i := 1; i <= 20; i++ {
		b.AddNode(SetStateNode(fmt.Sprintf("ss%d", i), "increment", "count", Literal("1")))
	}
	b.AddNode(EndNode("end"))

	// Wire: start → ss1 → ss2 → ... → ss20 → end
	b.AddEdge("__start__", "out", "ss1")
	for i := 1; i < 20; i++ {
		b.AddEdge(fmt.Sprintf("ss%d", i), "out", fmt.Sprintf("ss%d", i+1))
	}
	b.AddEdge("ss20", "out", "end")

	def := b.Build()
	h := NewHarness(t, def)

	// Set shutdown flag after 10 checkpoints (checkpoint happens after each node advance)
	h.StateRepo.OnCheckpoint = func(count int) {
		if count == 10 {
			h.ShutdownMgr.SetShuttingDown(true)
		}
	}

	exec := h.RunSync(map[string]interface{}{})

	// Should NOT be completed — walker stopped mid-chain
	if exec.Status == entities.ExecStatusCompleted {
		t.Fatal("expected execution to NOT complete — shutdown should have stopped it mid-chain")
	}

	// Path should have more than 1 node (start ran) but less than 22 (didn't finish)
	pathLen := len(exec.ExecutionPath)
	if pathLen < 2 {
		t.Fatalf("expected at least 2 nodes in path, got %d", pathLen)
	}
	if pathLen >= 22 {
		t.Fatalf("expected less than 22 nodes in path (shutdown should stop), got %d", pathLen)
	}

	t.Logf("Shutdown stopped chain at %d path entries (expected ~10-11)", pathLen)
}

func TestShutdown_MidLoop_10Iterations(t *testing.T) {
	items := make([]interface{}, 10)
	for i := range items {
		items[i] = fmt.Sprintf("item%d", i)
	}

	def := NewDefinition("ShutdownLoop10").
		WithState("items", "array", items).
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

	// Set shutdown after 8 checkpoints (loop generates checkpoints each iteration)
	h.StateRepo.OnCheckpoint = func(count int) {
		if count == 8 {
			h.ShutdownMgr.SetShuttingDown(true)
		}
	}

	exec := h.RunSync(map[string]interface{}{})

	if exec.Status == entities.ExecStatusCompleted {
		t.Fatal("expected execution to NOT complete — shutdown should have stopped mid-loop")
	}

	// Count should be > 0 (at least one iteration ran) and < 10 (didn't finish all)
	countVal, ok := exec.State["count"]
	if !ok {
		t.Fatal("expected state 'count' to exist")
	}

	count := toInt(countVal)
	if count <= 0 {
		t.Fatalf("expected count > 0, got %d", count)
	}
	if count >= 10 {
		t.Fatalf("expected count < 10 (shutdown should stop loop), got %d", count)
	}

	t.Logf("Shutdown stopped loop at iteration %d of 10", count)
}

func TestShutdown_MidFanout_3Branches(t *testing.T) {
	def := NewDefinition("ShutdownFanout3").
		WithState("a", "string", "").
		WithState("b", "string", "").
		WithState("c", "string", "").
		AddNode(StartNode("__start__")).
		AddNode(FanoutNode("fan1", 3, "")).
		AddNode(SetStateNode("ss_a", "set", "a", Literal("done_a"))).
		AddNode(SetStateNode("ss_b", "set", "b", Literal("done_b"))).
		AddNode(SetStateNode("ss_c", "set", "c", Literal("done_c"))).
		AddNode(MergeNode("merge1", 3)).
		AddNode(EndNode("end")).
		AddEdge("__start__", "out", "fan1").
		AddEdge("fan1", "out_1", "ss_a").
		AddEdge("fan1", "out_2", "ss_b").
		AddEdge("fan1", "out_3", "ss_c").
		AddEdge("ss_a", "out", "merge1").
		AddEdge("ss_b", "out", "merge1").
		AddEdge("ss_c", "out", "merge1").
		AddEdge("merge1", "out", "end").
		Build()

	h := NewHarness(t, def)

	// Set shutdown after fanout checkpoint (fanout internally checkpoints)
	h.StateRepo.OnCheckpoint = func(count int) {
		if count == 3 {
			h.ShutdownMgr.SetShuttingDown(true)
		}
	}

	exec := h.RunSync(map[string]interface{}{})

	// Execution should exist in KV with state saved
	if exec == nil {
		t.Fatal("expected execution to exist in KV")
	}

	t.Logf("Shutdown during fanout: status=%s, path=%d entries", exec.Status, len(exec.ExecutionPath))
}

func TestShutdown_ResumeAfterRestart(t *testing.T) {
	b := NewDefinition("ShutdownResume").
		WithState("count", "number", 0).
		AddNode(StartNode("__start__"))

	for i := 1; i <= 20; i++ {
		b.AddNode(SetStateNode(fmt.Sprintf("ss%d", i), "increment", "count", Literal("1")))
	}
	b.AddNode(EndNode("end"))

	b.AddEdge("__start__", "out", "ss1")
	for i := 1; i < 20; i++ {
		b.AddEdge(fmt.Sprintf("ss%d", i), "out", fmt.Sprintf("ss%d", i+1))
	}
	b.AddEdge("ss20", "out", "end")

	def := b.Build()

	// Phase 1: Run with shutdown at checkpoint 10
	h1 := NewHarness(t, def)
	h1.StateRepo.OnCheckpoint = func(count int) {
		if count == 10 {
			h1.ShutdownMgr.SetShuttingDown(true)
		}
	}

	exec1 := h1.RunSync(map[string]interface{}{})

	if exec1.Status == entities.ExecStatusCompleted {
		t.Fatal("phase 1: expected execution to NOT complete")
	}

	stoppedAt := len(exec1.ExecutionPath)
	stoppedNode := ""
	if len(exec1.ActiveNodeIDs) > 0 {
		stoppedNode = exec1.ActiveNodeIDs[0]
	}

	t.Logf("Phase 1: stopped at %d path entries, active node: %s", stoppedAt, stoppedNode)

	// Phase 2: Create new harness (simulate pod restart) — copy KV state
	h2 := NewHarness(t, def)

	// Seed the new harness KV with the execution from phase 1
	if err := h2.StateRepo.Create(exec1); err != nil {
		// Already exists from the initial trigger — overwrite via Save
		if saveErr := h2.StateRepo.Save(exec1); saveErr != nil {
			t.Fatalf("failed to seed phase 2 KV: %v", saveErr)
		}
	}

	// Send resume to continue from where we stopped
	h2.ResumeExecution(exec1)

	exec2 := h2.StateRepo.GetLatest()
	if exec2 == nil {
		t.Fatal("phase 2: no execution found after resume")
	}

	AssertCompleted(t, exec2)
	AssertState(t, exec2, "count", 20)

	t.Logf("Phase 2: completed with count=%v, path=%d entries", exec2.State["count"], len(exec2.ExecutionPath))
}

// toInt converts interface{} to int (handles float64 from JSON and int).
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case int64:
		return int(val)
	default:
		return 0
	}
}
