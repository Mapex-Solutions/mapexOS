package services

import (
	"sort"
	"testing"

	"workflow/src/modules/definitions/domain/entities"
)

func node(id, nodeType string) entities.WorkflowNode {
	return entities.WorkflowNode{ID: id, Type: nodeType}
}

func edge(source, target string) entities.WorkflowEdge {
	return entities.WorkflowEdge{Source: source, Target: target}
}

func edgeWithHandles(source, sourceHandle, target, targetHandle string) entities.WorkflowEdge {
	return entities.WorkflowEdge{
		Source:       source,
		SourceHandle: sourceHandle,
		Target:       target,
		TargetHandle: targetHandle,
	}
}

func sorted(ids []string) []string {
	if ids == nil {
		return nil
	}
	s := make([]string, len(ids))
	copy(s, ids)
	sort.Strings(s)
	return s
}

func TestNoCycles(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/start"),
		node("B", "core/code"),
		node("C", "core/end"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edge("B", "C"),
	}

	result := DetectTightCycles(nodes, edges)
	if len(result) != 0 {
		t.Fatalf("expected no tight cycles, got %v", result)
	}
}

func TestTightCycle(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
		node("B", "core/set_state"),
		node("C", "core/condition"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edge("B", "C"),
		edge("C", "A"),
	}

	result := sorted(DetectTightCycles(nodes, edges))
	expected := []string{"A", "B", "C"}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	}
}

func TestCycleWithAsync(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
		node("B", "core/delay"),
		node("C", "core/set_state"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edge("B", "C"),
		edge("C", "A"),
	}

	result := DetectTightCycles(nodes, edges)
	if len(result) != 0 {
		t.Fatalf("expected no tight cycles (cycle has async node), got %v", result)
	}
}

func TestLoopBodyExcluded(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("loop1", "core/loop"),
		node("body1", "core/code"),
	}
	edges := []entities.WorkflowEdge{
		edgeWithHandles("loop1", "body", "body1", ""),
		edge("body1", "loop1"),
	}

	// The "body" handle edge is excluded → body1→loop1 is the only edge.
	// That alone doesn't form a cycle since loop1→body1 was excluded.
	result := DetectTightCycles(nodes, edges)
	if len(result) != 0 {
		t.Fatalf("expected no tight cycles (loop body back-edge excluded), got %v", result)
	}
}

func TestAnnotationsExcluded(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
		node("B", "core/code"),
		node("note1", "core/text_note"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edgeWithHandles("A", "__note_out", "note1", "__note"),
	}

	result := DetectTightCycles(nodes, edges)
	if len(result) != 0 {
		t.Fatalf("expected no tight cycles, got %v", result)
	}
}

func TestMultipleTightCycles(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
		node("B", "core/code"),
		node("C", "core/code"),
		node("D", "core/code"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edge("B", "A"), // cycle 1: A↔B
		edge("C", "D"),
		edge("D", "C"), // cycle 2: C↔D
	}

	result := sorted(DetectTightCycles(nodes, edges))
	expected := []string{"A", "B", "C", "D"}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	}
}

func TestMixedCycles(t *testing.T) {
	// Cycle 1: A→B→A (tight, no async)
	// Cycle 2: C→D→E→C (has async wait_signal, allowed)
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
		node("B", "core/set_state"),
		node("C", "core/code"),
		node("D", "core/wait_signal"),
		node("E", "core/code"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edge("B", "A"),
		edge("C", "D"),
		edge("D", "E"),
		edge("E", "C"),
	}

	result := sorted(DetectTightCycles(nodes, edges))
	expected := []string{"A", "B"}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	}
}

func TestSelfLoop(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "A"),
	}

	result := DetectTightCycles(nodes, edges)
	if len(result) != 1 || result[0] != "A" {
		t.Fatalf("expected [A], got %v", result)
	}
}

func TestEmptyGraph(t *testing.T) {
	result := DetectTightCycles(nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty, got %v", result)
	}
}

func TestGroupFrameExcluded(t *testing.T) {
	nodes := []entities.WorkflowNode{
		node("A", "core/code"),
		node("B", "core/code"),
		node("frame1", "core/group_frame"),
	}
	edges := []entities.WorkflowEdge{
		edge("A", "B"),
		edge("frame1", "A"), // edge from annotation → ignored
	}

	result := DetectTightCycles(nodes, edges)
	if len(result) != 0 {
		t.Fatalf("expected no tight cycles, got %v", result)
	}
}
