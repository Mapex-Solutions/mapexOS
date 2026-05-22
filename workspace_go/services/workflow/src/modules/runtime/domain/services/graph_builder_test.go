package services

import (
	"testing"

	defEntities "workflow/src/modules/definitions/domain/entities"
	"workflow/src/modules/runtime/domain/entities"
)

func TestBuildGraph_FiltersVisualOnlyNodes(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "start-1", Type: "core/start"},
			{ID: "note-1", Type: "core/text_note"},
			{ID: "frame-1", Type: "core/group_frame"},
			{ID: "end-1", Type: "core/end"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "start-1", Target: "end-1"},
		},
	}

	graph := BuildGraph(def)

	if len(graph.Nodes) != 2 {
		t.Fatalf("expected 2 nodes (start+end), got %d", len(graph.Nodes))
	}
	if _, ok := graph.Nodes["note-1"]; ok {
		t.Fatal("text_note should be filtered out")
	}
	if _, ok := graph.Nodes["frame-1"]; ok {
		t.Fatal("group_frame should be filtered out")
	}
	if _, ok := graph.Nodes["start-1"]; !ok {
		t.Fatal("start node should be in graph")
	}
	if _, ok := graph.Nodes["end-1"]; !ok {
		t.Fatal("end node should be in graph")
	}
}

func TestBuildGraph_AdjacencyFromEdges(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "n1", Type: "core/start"},
			{ID: "n2", Type: "core/delay"},
			{ID: "n3", Type: "core/end"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "n1", SourceHandle: "out", Target: "n2"},
			{ID: "e2", Source: "n2", SourceHandle: "out", Target: "n3"},
		},
	}

	graph := BuildGraph(def)

	// n1 → n2
	if target, ok := graph.Adjacency["n1"]["out"]; !ok || target != "n2" {
		t.Fatalf("expected n1→out→n2, got %v", graph.Adjacency["n1"])
	}
	// n2 → n3
	if target, ok := graph.Adjacency["n2"]["out"]; !ok || target != "n3" {
		t.Fatalf("expected n2→out→n3, got %v", graph.Adjacency["n2"])
	}
}

func TestBuildGraph_DefaultHandleIsOut(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "a", Type: "core/start"},
			{ID: "b", Type: "core/end"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "a", SourceHandle: "", Target: "b"},
		},
	}

	graph := BuildGraph(def)

	if target, ok := graph.Adjacency["a"]["out"]; !ok || target != "b" {
		t.Fatalf("expected empty handle to default to 'out', got %v", graph.Adjacency["a"])
	}
}

func TestBuildGraph_SkipsVisualHandles(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "a", Type: "core/start"},
			{ID: "b", Type: "core/end"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "a", SourceHandle: "__visual", Target: "b"},
		},
	}

	graph := BuildGraph(def)

	if _, ok := graph.Adjacency["a"]; ok {
		t.Fatalf("visual handles (__ prefix) should be skipped, got %v", graph.Adjacency["a"])
	}
}

func TestBuildGraph_SkipsEdgesToVisualNodes(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "a", Type: "core/start"},
			{ID: "note", Type: "core/text_note"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "a", SourceHandle: "out", Target: "note"},
		},
	}

	graph := BuildGraph(def)

	if _, ok := graph.Adjacency["a"]; ok {
		t.Fatalf("edges targeting visual-only nodes should be skipped, got %v", graph.Adjacency["a"])
	}
}

func TestBuildGraph_GotoPairs(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "start", Type: "core/start"},
			{ID: "goto-sender", Type: "core/goto", Config: map[string]interface{}{
				"role":      "sender",
				"pairLabel": "error-handler",
			}},
			{ID: "goto-receiver", Type: "core/goto", Config: map[string]interface{}{
				"role":      "receiver",
				"pairLabel": "error-handler",
			}},
			{ID: "end", Type: "core/end"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "start", SourceHandle: "out", Target: "goto-sender"},
			{ID: "e2", Source: "goto-receiver", SourceHandle: "out", Target: "end"},
		},
	}

	graph := BuildGraph(def)

	// GotoPairs should map pairLabel → receiverNodeId
	if receiverID, ok := graph.GotoPairs["error-handler"]; !ok || receiverID != "goto-receiver" {
		t.Fatalf("expected GotoPairs['error-handler']='goto-receiver', got %v", graph.GotoPairs)
	}

	// Sender should have a logical edge to receiver (injected by step 4)
	if target, ok := graph.Adjacency["goto-sender"]["out"]; !ok || target != "goto-receiver" {
		t.Fatalf("expected goto-sender→out→goto-receiver, got %v", graph.Adjacency["goto-sender"])
	}
}

func TestBuildGraph_ParsedConfigs(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "start", Type: "core/start"},
			{ID: "delay-1", Type: "core/delay", Config: map[string]interface{}{
				"duration": float64(5),
				"unit":     "s",
			}},
		},
	}

	graph := BuildGraph(def)

	// core/start has no config parsing, should be nil
	if _, ok := graph.ParsedConfigs["start"]; ok {
		t.Fatal("core/start should not have a parsed config")
	}

	// core/delay should have a parsed config
	if _, ok := graph.ParsedConfigs["delay-1"]; !ok {
		t.Fatal("core/delay should have a parsed config")
	}
}

func TestBuildGraph_Timezone(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Timezone: defEntities.FieldValue{
			Type:  defEntities.FieldValueLiteral,
			Value: "America/Sao_Paulo",
		},
		Nodes: []defEntities.WorkflowNode{
			{ID: "start", Type: "core/start"},
		},
	}

	graph := BuildGraph(def)

	if graph.Timezone != "America/Sao_Paulo" {
		t.Fatalf("expected timezone 'America/Sao_Paulo', got '%s'", graph.Timezone)
	}
}

func TestBuildGraph_TimezoneNonLiteral(t *testing.T) {
	def := &defEntities.WorkflowDefinition{
		Timezone: defEntities.FieldValue{
			Type:  defEntities.FieldValueState,
			Value: "user_tz",
		},
		Nodes: []defEntities.WorkflowNode{
			{ID: "start", Type: "core/start"},
		},
	}

	graph := BuildGraph(def)

	if graph.Timezone != "" {
		t.Fatalf("expected empty timezone for non-literal, got '%s'", graph.Timezone)
	}
}

func TestBuildGraph_ComplexDAG(t *testing.T) {
	// start → condition → (true_handle→delay, false_handle→end)
	//                       delay → end
	def := &defEntities.WorkflowDefinition{
		Nodes: []defEntities.WorkflowNode{
			{ID: "start", Type: "core/start"},
			{ID: "cond", Type: "core/condition", Config: map[string]interface{}{
				"condition": map[string]interface{}{
					"id":    "g1",
					"logic": "AND",
					"items": []interface{}{},
				},
			}},
			{ID: "delay", Type: "core/delay", Config: map[string]interface{}{
				"duration": float64(10),
				"unit":     "s",
			}},
			{ID: "end", Type: "core/end"},
		},
		Edges: []defEntities.WorkflowEdge{
			{ID: "e1", Source: "start", SourceHandle: "out", Target: "cond"},
			{ID: "e2", Source: "cond", SourceHandle: "true", Target: "delay"},
			{ID: "e3", Source: "cond", SourceHandle: "false", Target: "end"},
			{ID: "e4", Source: "delay", SourceHandle: "out", Target: "end"},
		},
	}

	graph := BuildGraph(def)

	if len(graph.Nodes) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(graph.Nodes))
	}

	// Condition has two handles: true→delay, false→end
	condAdj := graph.Adjacency["cond"]
	if condAdj == nil || condAdj["true"] != "delay" {
		t.Fatalf("expected cond→true→delay, got %v", condAdj)
	}
	if condAdj["false"] != "end" {
		t.Fatalf("expected cond→false→end, got %v", condAdj)
	}

	// Delay → end
	if graph.Adjacency["delay"]["out"] != "end" {
		t.Fatalf("expected delay→out→end, got %v", graph.Adjacency["delay"])
	}

	// ParsedConfigs should have entries for condition and delay
	if graph.ParsedConfigs["cond"] == nil {
		t.Fatal("expected parsed config for condition node")
	}
	if graph.ParsedConfigs["delay"] == nil {
		t.Fatal("expected parsed config for delay node")
	}
}

func TestBuildGraph_EmptyDefinition(t *testing.T) {
	def := &defEntities.WorkflowDefinition{}

	graph := BuildGraph(def)

	if len(graph.Nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(graph.Nodes))
	}
	if len(graph.Adjacency) != 0 {
		t.Fatalf("expected 0 adjacency entries, got %d", len(graph.Adjacency))
	}
}

func TestResolveDefaultHandle_ReturnsOutEdge(t *testing.T) {
	graph := &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{
			"n1": {"out": "n2"},
		},
	}
	if handle := graph.ResolveDefaultHandle("n1"); handle != "out" {
		t.Fatalf("expected 'out', got %q", handle)
	}
}

func TestResolveDefaultHandle_ReturnsSuccessEdge(t *testing.T) {
	graph := &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{
			"n1": {"success": "n2", "error": "n3"},
		},
	}
	handle := graph.ResolveDefaultHandle("n1")
	if handle != "success" {
		t.Fatalf("expected 'success', got %q", handle)
	}
}

func TestResolveDefaultHandle_FallbackWhenOnlyErrorAndTimeout(t *testing.T) {
	graph := &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{
			"n1": {"error": "n3", "timeout": "n4"},
		},
	}
	if handle := graph.ResolveDefaultHandle("n1"); handle != "out" {
		t.Fatalf("expected fallback 'out', got %q", handle)
	}
}

func TestResolveDefaultHandle_FallbackWhenNoEdges(t *testing.T) {
	graph := &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{
			"n1": {},
		},
	}
	if handle := graph.ResolveDefaultHandle("n1"); handle != "out" {
		t.Fatalf("expected fallback 'out', got %q", handle)
	}
}

func TestResolveDefaultHandle_FallbackWhenNodeNotInAdjacency(t *testing.T) {
	graph := &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{},
	}
	if handle := graph.ResolveDefaultHandle("unknown"); handle != "out" {
		t.Fatalf("expected fallback 'out', got %q", handle)
	}
}

func TestIsVisualOnly(t *testing.T) {
	tests := []struct {
		nodeType string
		expected bool
	}{
		{"core/text_note", true},
		{"core/group_frame", true},
		{"core/start", false},
		{"core/end", false},
		{"core/delay", false},
		{"core/condition", false},
	}

	for _, tt := range tests {
		t.Run(tt.nodeType, func(t *testing.T) {
			if got := isVisualOnly(tt.nodeType); got != tt.expected {
				t.Fatalf("isVisualOnly(%q) = %v, want %v", tt.nodeType, got, tt.expected)
			}
		})
	}
}
