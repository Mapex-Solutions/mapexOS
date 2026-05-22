package services

import (
	"testing"

	"workflow/src/modules/definitions/domain/entities"
)

func configNode(id, nodeType string, config map[string]interface{}) entities.WorkflowNode {
	return entities.WorkflowNode{ID: id, Type: nodeType, Config: config}
}

// --- Valid nodes (all types) ---

func TestValidNode_Start(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/start", nil)}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_End(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/end", nil)}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Log(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/log", map[string]interface{}{"message": "hello"})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Condition(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/condition", map[string]interface{}{
		"logic": "AND",
		"items": []interface{}{map[string]interface{}{"type": "condition"}},
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Code(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/code", map[string]interface{}{
		"script": "console.log('hello')",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_SetState(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/set_state", map[string]interface{}{
		"targetField": "count",
		"operation":   "increment",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Switch(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/switch", map[string]interface{}{
		"cases": []interface{}{map[string]interface{}{"id": "c1"}},
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Delay(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/delay", map[string]interface{}{
		"duration": 5,
		"unit":     "minutes",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_WaitSignal(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/wait_signal", map[string]interface{}{
		"signalName": "approval",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Loop(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/loop", map[string]interface{}{
		"source": map[string]interface{}{"type": "state", "value": "items"},
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Fanout(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/fanout", map[string]interface{}{
		"branches": 3,
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Merge(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/merge", map[string]interface{}{
		"branches": 3,
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Sequence(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/sequence", map[string]interface{}{
		"steps": 5,
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_TriggerEvent(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/trigger_event", map[string]interface{}{
		"eventType": "order.created",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_WaitFor(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/wait_for", map[string]interface{}{
		"field":    "status",
		"operator": "equals",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Goto(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/goto", map[string]interface{}{
		"role":      "sender",
		"pairLabel": "error-handler",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidNode_Subworkflow(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/subworkflow", map[string]interface{}{
		"workflowId": "abc123",
	})}
	if errs := ValidateNodes(nodes); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

// --- Invalid node tests ---

func TestCondition_MissingItems(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/condition", map[string]interface{}{
		"logic": "AND",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "items is required" {
		t.Fatalf("expected 'items is required', got %v", errs)
	}
}

func TestCondition_EmptyConfig(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/condition", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "items is required" {
		t.Fatalf("expected 'items is required', got %v", errs)
	}
}

func TestCode_EmptyScript(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/code", map[string]interface{}{
		"script": "",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "script is required" {
		t.Fatalf("expected 'script is required', got %v", errs)
	}
}

func TestSetState_MissingFields(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/set_state", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 {
		t.Fatalf("expected 1 node with errors, got %d", len(errs))
	}
	if len(errs[0].Errors) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs[0].Errors)
	}
}

func TestSetState_InvalidOperation(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/set_state", map[string]interface{}{
		"targetField": "count",
		"operation":   "multiply",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || len(errs[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", errs)
	}
}

func TestSwitch_EmptyCases(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/switch", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "cases is required and must not be empty" {
		t.Fatalf("expected 'cases is required and must not be empty', got %v", errs)
	}
}

func TestDelay_InvalidUnit(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/delay", map[string]interface{}{
		"duration": 5,
		"unit":     "weeks",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || len(errs[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", errs)
	}
}

func TestDelay_ZeroDuration(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/delay", map[string]interface{}{
		"duration": 0,
		"unit":     "seconds",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "duration must be greater than 0" {
		t.Fatalf("expected 'duration must be greater than 0', got %v", errs)
	}
}

func TestFanout_OutOfRange(t *testing.T) {
	// branches = 0
	nodes := []entities.WorkflowNode{configNode("1", "core/fanout", map[string]interface{}{
		"branches": 0,
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "branches must be greater than 0" {
		t.Fatalf("expected 'branches must be greater than 0', got %v", errs)
	}

	// branches = 25
	nodes = []entities.WorkflowNode{configNode("1", "core/fanout", map[string]interface{}{
		"branches": 25,
	})}
	errs = ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "branches must not exceed 20" {
		t.Fatalf("expected 'branches must not exceed 20', got %v", errs)
	}
}

func TestGoto_InvalidRole(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/goto", map[string]interface{}{
		"role":      "observer",
		"pairLabel": "label1",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "role must be one of: sender, receiver" {
		t.Fatalf("expected role validation error, got %v", errs)
	}
}

func TestGoto_MissingPairLabel(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/goto", map[string]interface{}{
		"role":      "sender",
		"pairLabel": "",
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "pairLabel is required" {
		t.Fatalf("expected 'pairLabel is required', got %v", errs)
	}
}

func TestMultipleNodesMultipleErrors(t *testing.T) {
	nodes := []entities.WorkflowNode{
		configNode("1", "core/code", map[string]interface{}{}),
		configNode("2", "core/set_state", map[string]interface{}{}),
		configNode("3", "core/start", nil), // valid
	}
	errs := ValidateNodes(nodes)
	if len(errs) != 2 {
		t.Fatalf("expected 2 nodes with errors, got %d", len(errs))
	}
}

func TestVisualNodesSkipped(t *testing.T) {
	nodes := []entities.WorkflowNode{
		configNode("1", "core/text_note", nil),
		configNode("2", "core/group_frame", nil),
	}
	errs := ValidateNodes(nodes)
	if len(errs) != 0 {
		t.Fatalf("expected no errors for visual nodes, got %v", errs)
	}
}

func TestWaitFor_MissingFields(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/wait_for", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || len(errs[0].Errors) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs)
	}
}

func TestSubworkflow_MissingWorkflowId(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/subworkflow", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "workflowId is required" {
		t.Fatalf("expected 'workflowId is required', got %v", errs)
	}
}

func TestWaitSignal_MissingSignalName(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/wait_signal", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "signalName is required" {
		t.Fatalf("expected 'signalName is required', got %v", errs)
	}
}

func TestLoop_MissingSource(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/loop", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "source is required" {
		t.Fatalf("expected 'source is required', got %v", errs)
	}
}

func TestTriggerEvent_MissingEventType(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/trigger_event", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "eventType is required" {
		t.Fatalf("expected 'eventType is required', got %v", errs)
	}
}

func TestMerge_ZeroBranches(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/merge", map[string]interface{}{
		"branches": 0,
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "branches must be greater than 0" {
		t.Fatalf("expected 'branches must be greater than 0', got %v", errs)
	}
}

func TestSequence_ZeroSteps(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "core/sequence", map[string]interface{}{
		"steps": 0,
	})}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "steps must be greater than 0" {
		t.Fatalf("expected 'steps must be greater than 0', got %v", errs)
	}
}

func TestUnknownNodeType(t *testing.T) {
	nodes := []entities.WorkflowNode{configNode("1", "custom/unknown", map[string]interface{}{})}
	errs := ValidateNodes(nodes)
	if len(errs) != 0 {
		t.Fatalf("expected no errors for unknown node type, got %v", errs)
	}
}

func TestNilConfig(t *testing.T) {
	// Nodes that require config but have nil — should still produce errors
	nodes := []entities.WorkflowNode{configNode("1", "core/code", nil)}
	errs := ValidateNodes(nodes)
	if len(errs) != 1 || errs[0].Errors[0] != "script is required" {
		t.Fatalf("expected 'script is required', got %v", errs)
	}
}
