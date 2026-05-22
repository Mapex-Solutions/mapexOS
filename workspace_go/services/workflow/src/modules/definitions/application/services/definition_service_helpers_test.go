package services

import (
	"testing"

	"workflow/src/modules/definitions/domain/entities"

	contractDefs "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/definitions"
	mongoModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// --- findNode tests ---

func TestFindNode_Found(t *testing.T) {
	nodes := []entities.WorkflowNode{
		{ID: "n1", Type: "core/start"},
		{ID: "n2", Type: "core/code", Config: map[string]interface{}{"script": "console.log('hi')"}},
		{ID: "n3", Type: "core/end"},
	}
	node, ok := findNode(nodes, "n2")
	if !ok {
		t.Fatal("expected to find node n2")
	}
	if node.Type != "core/code" {
		t.Fatalf("expected core/code, got %s", node.Type)
	}
}

func TestFindNode_NotFound(t *testing.T) {
	nodes := []entities.WorkflowNode{
		{ID: "n1", Type: "core/start"},
	}
	_, ok := findNode(nodes, "missing")
	if ok {
		t.Fatal("expected ok=false for missing node")
	}
}

func TestFindNode_EmptySlice(t *testing.T) {
	_, ok := findNode(nil, "n1")
	if ok {
		t.Fatal("expected ok=false for empty slice")
	}
}

// --- getNodeScript tests ---

func TestGetNodeScript_HasScript(t *testing.T) {
	node := entities.WorkflowNode{
		ID:     "n1",
		Type:   "core/code",
		Config: map[string]interface{}{"script": "return 42;"},
	}
	script := getNodeScript(node)
	if script != "return 42;" {
		t.Fatalf("expected 'return 42;', got %q", script)
	}
}

func TestGetNodeScript_NoScript(t *testing.T) {
	node := entities.WorkflowNode{
		ID:     "n1",
		Type:   "core/start",
		Config: nil,
	}
	script := getNodeScript(node)
	if script != "" {
		t.Fatalf("expected empty string, got %q", script)
	}
}

func TestGetNodeScript_EmptyConfig(t *testing.T) {
	node := entities.WorkflowNode{
		ID:     "n1",
		Type:   "core/code",
		Config: map[string]interface{}{},
	}
	script := getNodeScript(node)
	if script != "" {
		t.Fatalf("expected empty string, got %q", script)
	}
}

// --- getCodeNodeIds tests ---

func TestGetCodeNodeIds_HasCodeNodes(t *testing.T) {
	def := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/start"},
			{ID: "n2", Type: "core/code"},
			{ID: "n3", Type: "core/end"},
			{ID: "n4", Type: "core/code"},
		},
	}
	ids := getCodeNodeIds(def)
	if len(ids) != 2 {
		t.Fatalf("expected 2 code nodes, got %d", len(ids))
	}
	if ids[0] != "n2" || ids[1] != "n4" {
		t.Fatalf("expected [n2, n4], got %v", ids)
	}
}

func TestGetCodeNodeIds_NoCodeNodes(t *testing.T) {
	def := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/start"},
			{ID: "n2", Type: "core/end"},
		},
	}
	ids := getCodeNodeIds(def)
	if len(ids) != 0 {
		t.Fatalf("expected 0 code nodes, got %d", len(ids))
	}
}

// --- extractOrgId tests ---

func TestExtractOrgId_HasOrg(t *testing.T) {
	orgId := mongoModel.NewObjectID()
	def := &entities.WorkflowDefinition{OrgID: &orgId}
	result := extractOrgId(def)
	if result != orgId.Hex() {
		t.Fatalf("expected %s, got %s", orgId.Hex(), result)
	}
}

func TestExtractOrgId_NilOrg(t *testing.T) {
	def := &entities.WorkflowDefinition{}
	result := extractOrgId(def)
	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

// --- diffCodeNodes tests ---

func TestDiffCodeNodes_Added(t *testing.T) {
	before := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
		},
	}
	after := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
			{ID: "n2", Type: "core/code", Config: map[string]interface{}{"script": "b"}},
		},
	}
	added, removed, modified := diffCodeNodes(before, after)
	if len(added) != 1 || added[0] != "n2" {
		t.Fatalf("expected added=[n2], got %v", added)
	}
	if len(removed) != 0 {
		t.Fatalf("expected no removed, got %v", removed)
	}
	if len(modified) != 0 {
		t.Fatalf("expected no modified, got %v", modified)
	}
}

func TestDiffCodeNodes_Removed(t *testing.T) {
	before := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
			{ID: "n2", Type: "core/code", Config: map[string]interface{}{"script": "b"}},
		},
	}
	after := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
		},
	}
	added, removed, modified := diffCodeNodes(before, after)
	if len(added) != 0 {
		t.Fatalf("expected no added, got %v", added)
	}
	if len(removed) != 1 || removed[0] != "n2" {
		t.Fatalf("expected removed=[n2], got %v", removed)
	}
	if len(modified) != 0 {
		t.Fatalf("expected no modified, got %v", modified)
	}
}

func TestDiffCodeNodes_Modified(t *testing.T) {
	before := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "v1"}},
		},
	}
	after := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "v2"}},
		},
	}
	added, removed, modified := diffCodeNodes(before, after)
	if len(added) != 0 {
		t.Fatalf("expected no added, got %v", added)
	}
	if len(removed) != 0 {
		t.Fatalf("expected no removed, got %v", removed)
	}
	if len(modified) != 1 || modified[0] != "n1" {
		t.Fatalf("expected modified=[n1], got %v", modified)
	}
}

func TestDiffCodeNodes_NoChanges(t *testing.T) {
	nodes := []entities.WorkflowNode{
		{ID: "n1", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
	}
	before := &entities.WorkflowDefinition{Nodes: nodes}
	after := &entities.WorkflowDefinition{Nodes: nodes}
	added, removed, modified := diffCodeNodes(before, after)
	if len(added)+len(removed)+len(modified) != 0 {
		t.Fatalf("expected no changes, got added=%v removed=%v modified=%v", added, removed, modified)
	}
}

func TestDiffCodeNodes_IgnoresNonCodeNodes(t *testing.T) {
	before := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/start"},
			{ID: "n2", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
		},
	}
	after := &entities.WorkflowDefinition{
		Nodes: []entities.WorkflowNode{
			{ID: "n1", Type: "core/start"},
			{ID: "n3", Type: "core/condition"}, // non-code, should be ignored
			{ID: "n2", Type: "core/code", Config: map[string]interface{}{"script": "a"}},
		},
	}
	added, removed, modified := diffCodeNodes(before, after)
	if len(added)+len(removed)+len(modified) != 0 {
		t.Fatalf("expected no code node changes, got added=%v removed=%v modified=%v", added, removed, modified)
	}
}

// --- contractNodesToEntity tests ---

func TestContractNodesToEntity(t *testing.T) {
	contractNodes := []contractDefs.WorkflowNode{
		{
			ID:    "n1",
			Type:  "core/start",
			Label: "Start",
			Position: contractDefs.Position{
				X: 100,
				Y: 200,
			},
			Config:       map[string]interface{}{"key": "val"},
			ParentNodeID: "parent1",
		},
	}
	result := contractNodesToEntity(contractNodes)
	if len(result) != 1 {
		t.Fatalf("expected 1 node, got %d", len(result))
	}
	if result[0].ID != "n1" {
		t.Fatalf("expected ID 'n1', got %q", result[0].ID)
	}
	if result[0].Position.X != 100 || result[0].Position.Y != 200 {
		t.Fatalf("expected position (100,200), got (%v,%v)", result[0].Position.X, result[0].Position.Y)
	}
	if result[0].ParentNodeID != "parent1" {
		t.Fatalf("expected parentNodeID 'parent1', got %q", result[0].ParentNodeID)
	}
}

func TestContractNodesToEntity_Empty(t *testing.T) {
	result := contractNodesToEntity(nil)
	if len(result) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(result))
	}
}

// --- contractEdgesToEntity tests ---

func TestContractEdgesToEntity(t *testing.T) {
	contractEdges := []contractDefs.WorkflowEdge{
		{
			ID:           "e1",
			Source:       "n1",
			SourceHandle: "out",
			Target:       "n2",
			TargetHandle: "in",
			Label:        "edge1",
			PathOffsetX:  10,
			PathOffsetY:  20,
		},
	}
	result := contractEdgesToEntity(contractEdges)
	if len(result) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(result))
	}
	if result[0].Source != "n1" || result[0].Target != "n2" {
		t.Fatalf("expected source=n1 target=n2, got source=%q target=%q", result[0].Source, result[0].Target)
	}
	if result[0].PathOffsetX != 10 || result[0].PathOffsetY != 20 {
		t.Fatalf("expected offsets (10,20), got (%v,%v)", result[0].PathOffsetX, result[0].PathOffsetY)
	}
}

func TestContractEdgesToEntity_Empty(t *testing.T) {
	result := contractEdgesToEntity(nil)
	if len(result) != 0 {
		t.Fatalf("expected 0 edges, got %d", len(result))
	}
}
