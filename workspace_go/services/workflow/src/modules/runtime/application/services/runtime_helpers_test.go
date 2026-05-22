package services

import (
	"strings"
	"testing"

	appConstants "workflow/src/modules/runtime/application/constants"
)

/**
 * Token Generation Tests
 */

func TestGenerateExecutionToken_Deterministic(t *testing.T) {
	token1 := generateExecutionToken("uuid-123", "node-abc", 0)
	token2 := generateExecutionToken("uuid-123", "node-abc", 0)
	if token1 != token2 {
		t.Fatalf("expected same token for same inputs, got %s and %s", token1, token2)
	}
}

func TestGenerateExecutionToken_DifferentAttempt(t *testing.T) {
	token0 := generateExecutionToken("uuid-123", "node-abc", 0)
	token1 := generateExecutionToken("uuid-123", "node-abc", 1)
	if token0 == token1 {
		t.Fatalf("expected different tokens for different attempts, both got %s", token0)
	}
}

func TestGenerateExecutionToken_DifferentNode(t *testing.T) {
	tokenA := generateExecutionToken("uuid-123", "node-a", 0)
	tokenB := generateExecutionToken("uuid-123", "node-b", 0)
	if tokenA == tokenB {
		t.Fatalf("expected different tokens for different nodes, both got %s", tokenA)
	}
}

func TestGenerateExecutionToken_Length(t *testing.T) {
	token := generateExecutionToken("uuid-123", "node-abc", 0)
	if len(token) != 32 {
		t.Fatalf("expected 32 hex chars, got %d: %s", len(token), token)
	}
}

/**
 * Msg-Id Tests
 */

func TestBuildMsgId_Format(t *testing.T) {
	state := map[string]interface{}{appConstants.NodeStateKeyLoopIndex: 3}
	msgId := buildMsgId("uuid-123", "node-abc", 2, state)
	expected := "uuid-123:node-abc:2:3"
	if msgId != expected {
		t.Fatalf("expected %q, got %q", expected, msgId)
	}
}

func TestBuildMsgId_Deterministic(t *testing.T) {
	state := map[string]interface{}{}
	id1 := buildMsgId("uuid-123", "node-abc", 0, state)
	id2 := buildMsgId("uuid-123", "node-abc", 0, state)
	if id1 != id2 {
		t.Fatalf("expected same Msg-Id for same inputs, got %s and %s", id1, id2)
	}
}

func TestBuildMsgId_DifferentLoop(t *testing.T) {
	state0 := map[string]interface{}{appConstants.NodeStateKeyLoopIndex: 0}
	state3 := map[string]interface{}{appConstants.NodeStateKeyLoopIndex: 3}
	id0 := buildMsgId("uuid-123", "node-abc", 0, state0)
	id3 := buildMsgId("uuid-123", "node-abc", 0, state3)
	if id0 == id3 {
		t.Fatalf("expected different Msg-Id for different loop index, both got %s", id0)
	}
}

func TestBuildMsgId_DifferentAttempt(t *testing.T) {
	state := map[string]interface{}{}
	id0 := buildMsgId("uuid-123", "node-abc", 0, state)
	id1 := buildMsgId("uuid-123", "node-abc", 1, state)
	if id0 == id1 {
		t.Fatalf("expected different Msg-Id for different attempt, both got %s", id0)
	}
}

func TestBuildMsgId_NoLoopIndex(t *testing.T) {
	state := map[string]interface{}{}
	msgId := buildMsgId("uuid-123", "node-abc", 0, state)
	if !strings.HasSuffix(msgId, ":0") {
		t.Fatalf("expected loopIndex=0 when absent, got %s", msgId)
	}
}

func TestBuildMsgId_Float64LoopIndex(t *testing.T) {
	state := map[string]interface{}{appConstants.NodeStateKeyLoopIndex: float64(7)}
	msgId := buildMsgId("uuid-123", "node-abc", 0, state)
	expected := "uuid-123:node-abc:0:7"
	if msgId != expected {
		t.Fatalf("expected %q, got %q", expected, msgId)
	}
}

/**
 * getRetryAttempt Tests
 */

func TestGetRetryAttempt_Int(t *testing.T) {
	nodeStates := map[string]map[string]interface{}{
		"node-1": {appConstants.NodeStateKeyInternalRetry: 3},
	}
	result := getRetryAttempt(nodeStates, "node-1")
	if result != 3 {
		t.Fatalf("expected 3, got %d", result)
	}
}

func TestGetRetryAttempt_Float64(t *testing.T) {
	nodeStates := map[string]map[string]interface{}{
		"node-1": {appConstants.NodeStateKeyInternalRetry: float64(5)},
	}
	result := getRetryAttempt(nodeStates, "node-1")
	if result != 5 {
		t.Fatalf("expected 5, got %d", result)
	}
}

func TestGetRetryAttempt_Missing(t *testing.T) {
	nodeStates := map[string]map[string]interface{}{
		"node-1": {appConstants.NodeStateKeyWaitType: "callback"},
	}
	result := getRetryAttempt(nodeStates, "node-1")
	if result != 0 {
		t.Fatalf("expected 0, got %d", result)
	}
}

func TestGetRetryAttempt_NodeNotFound(t *testing.T) {
	nodeStates := map[string]map[string]interface{}{}
	result := getRetryAttempt(nodeStates, "missing-node")
	if result != 0 {
		t.Fatalf("expected 0, got %d", result)
	}
}

/**
 * Subworkflow UUID Tests
 */

func TestGenerateSubworkflowUUID_Deterministic(t *testing.T) {
	uuid1 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	uuid2 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	if uuid1 != uuid2 {
		t.Fatalf("expected same UUID, got %s and %s", uuid1, uuid2)
	}
}

func TestGenerateSubworkflowUUID_DifferentParent(t *testing.T) {
	uuid1 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	uuid2 := generateSubworkflowUUID("660e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	if uuid1 == uuid2 {
		t.Fatal("expected different UUIDs for different parents")
	}
}

func TestGenerateSubworkflowUUID_DifferentNode(t *testing.T) {
	uuid1 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	uuid2 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_2", 0)
	if uuid1 == uuid2 {
		t.Fatal("expected different UUIDs for different nodes")
	}
}

func TestGenerateSubworkflowUUID_DifferentAttempt(t *testing.T) {
	uuid1 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	uuid2 := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 1)
	if uuid1 == uuid2 {
		t.Fatal("expected different UUIDs for different attempts")
	}
}

func TestGenerateSubworkflowUUID_ValidFormat(t *testing.T) {
	result := generateSubworkflowUUID("550e8400-e29b-41d4-a716-446655440000", "subwf_1", 0)
	if len(result) != 36 {
		t.Fatalf("expected 36-char UUID, got %d: %s", len(result), result)
	}
	count := 0
	for _, c := range result {
		if c == '-' {
			count++
		}
	}
	if count != 4 {
		t.Fatalf("expected 4 hyphens in UUID, got %d: %s", count, result)
	}
}

func TestGenerateSubworkflowUUID_InvalidParentFallback(t *testing.T) {
	result := generateSubworkflowUUID("not-a-uuid", "subwf_1", 0)
	if len(result) != 36 {
		t.Fatalf("expected fallback to random UUID (36 chars), got %d: %s", len(result), result)
	}
}
