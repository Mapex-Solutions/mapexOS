package dagwalker

import (
	"fmt"
	"strings"
	"testing"

	"workflow/src/modules/runtime/domain/entities"
)

// AssertCompleted asserts the execution reached "completed" status.
func AssertCompleted(t *testing.T, exec *entities.WorkflowExecution) {
	t.Helper()
	if exec == nil {
		t.Fatal("execution is nil")
	}
	if exec.Status != entities.ExecStatusCompleted {
		t.Fatalf("expected status=completed, got %q (errorInfo=%+v)", exec.Status, exec.ErrorInfo)
	}
}

// AssertFailed asserts the execution reached "failed" status with the given error code.
func AssertFailed(t *testing.T, exec *entities.WorkflowExecution, expectedCode string) {
	t.Helper()
	if exec == nil {
		t.Fatal("execution is nil")
	}
	if exec.Status != entities.ExecStatusFailed {
		t.Fatalf("expected status=failed, got %q", exec.Status)
	}
	if exec.ErrorInfo == nil {
		t.Fatal("expected errorInfo to be set, got nil")
	}
	if exec.ErrorInfo.Code != expectedCode {
		t.Fatalf("expected errorInfo.code=%q, got %q (message=%q)", expectedCode, exec.ErrorInfo.Code, exec.ErrorInfo.Message)
	}
}

// AssertFailedWithMessage asserts the execution failed with the given error code and message.
func AssertFailedWithMessage(t *testing.T, exec *entities.WorkflowExecution, expectedCode, expectedMessage string) {
	t.Helper()
	AssertFailed(t, exec, expectedCode)
	if exec.ErrorInfo.Message != expectedMessage {
		t.Fatalf("expected errorInfo.message=%q, got %q", expectedMessage, exec.ErrorInfo.Message)
	}
}

// AssertErrorHasMetadata asserts that errorInfo has nodeId, nodeType, and a non-zero timestamp.
func AssertErrorHasMetadata(t *testing.T, exec *entities.WorkflowExecution) {
	t.Helper()
	if exec.ErrorInfo == nil {
		t.Fatal("expected errorInfo to be set, got nil")
	}
	if exec.ErrorInfo.NodeID == "" {
		t.Fatal("expected errorInfo.nodeId to be set")
	}
	if exec.ErrorInfo.NodeType == "" {
		t.Fatal("expected errorInfo.nodeType to be set")
	}
	if exec.ErrorInfo.Timestamp.IsZero() {
		t.Fatal("expected errorInfo.timestamp to be non-zero")
	}
}

// AssertWaiting asserts the execution is in "waiting" status.
func AssertWaiting(t *testing.T, exec *entities.WorkflowExecution) {
	t.Helper()
	if exec == nil {
		t.Fatal("execution is nil")
	}
	if exec.Status != entities.ExecStatusWaiting {
		t.Fatalf("expected status=waiting, got %q", exec.Status)
	}
}

// AssertState asserts a specific key in the execution state equals the expected value.
func AssertState(t *testing.T, exec *entities.WorkflowExecution, key string, expected interface{}) {
	t.Helper()
	actual, exists := exec.State[key]
	if !exists {
		t.Fatalf("state key %q does not exist (state=%+v)", key, exec.State)
	}
	// JSON serialization normalizes types — compare as strings for simplicity
	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)
	if actualStr != expectedStr {
		t.Fatalf("state[%q]: expected %v (%T), got %v (%T)", key, expected, expected, actual, actual)
	}
}

// AssertStateNotExists asserts a key does NOT exist in the execution state.
func AssertStateNotExists(t *testing.T, exec *entities.WorkflowExecution, key string) {
	t.Helper()
	if _, exists := exec.State[key]; exists {
		t.Fatalf("state key %q should not exist, but has value %v", key, exec.State[key])
	}
}

// AssertNodeOutput asserts a specific node's output equals the expected value.
func AssertNodeOutput(t *testing.T, exec *entities.WorkflowExecution, nodeID string, expected interface{}) {
	t.Helper()
	actual, exists := exec.NodeOutputs[nodeID]
	if !exists {
		t.Fatalf("nodeOutputs[%q] does not exist", nodeID)
	}
	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)
	if actualStr != expectedStr {
		t.Fatalf("nodeOutputs[%q]: expected %v, got %v", nodeID, expected, actual)
	}
}

// AssertPathContains asserts the execution path contains the given nodeID.
func AssertPathContains(t *testing.T, exec *entities.WorkflowExecution, nodeID string) {
	t.Helper()
	for _, entry := range exec.ExecutionPath {
		if entry.NodeID == nodeID {
			return
		}
	}
	ids := make([]string, len(exec.ExecutionPath))
	for i, e := range exec.ExecutionPath {
		ids[i] = e.NodeID
	}
	t.Fatalf("execution path does not contain %q (path=%v)", nodeID, ids)
}

// AssertPathNodeStatus asserts a specific node in the execution path has the expected status.
func AssertPathNodeStatus(t *testing.T, exec *entities.WorkflowExecution, nodeID, expectedStatus string) {
	t.Helper()
	for _, entry := range exec.ExecutionPath {
		if entry.NodeID == nodeID {
			if entry.Status != expectedStatus {
				t.Fatalf("path node %q: expected status=%q, got %q", nodeID, expectedStatus, entry.Status)
			}
			return
		}
	}
	t.Fatalf("node %q not found in execution path", nodeID)
}

// AssertPathLength asserts the execution path has exactly N entries.
func AssertPathLength(t *testing.T, exec *entities.WorkflowExecution, expected int) {
	t.Helper()
	if len(exec.ExecutionPath) != expected {
		ids := make([]string, len(exec.ExecutionPath))
		for i, e := range exec.ExecutionPath {
			ids[i] = fmt.Sprintf("%s(%s)", e.NodeID, e.Status)
		}
		t.Fatalf("expected path length=%d, got %d (path=%v)", expected, len(exec.ExecutionPath), ids)
	}
}

// AssertDispatched asserts the publisher received a dispatch call for the given method and nodeID.
func AssertDispatched(t *testing.T, pub *CapturingPublisher, method, nodeID string) {
	t.Helper()
	if pub.FindDispatch(method, nodeID) == nil {
		t.Fatalf("expected dispatch %s for node %q, but not found (events=%+v)", method, nodeID, pub.Events)
	}
}

// AssertDispatchCount asserts the publisher received exactly N calls of the given method.
func AssertDispatchCount(t *testing.T, pub *CapturingPublisher, method string, expected int) {
	t.Helper()
	actual := pub.CountMethod(method)
	if actual != expected {
		t.Fatalf("expected %d calls to %s, got %d", expected, method, actual)
	}
}

// AssertSchedulePublished asserts that PublishSchedule was called for the given nodeID.
func AssertSchedulePublished(t *testing.T, publisher *CapturingPublisher, nodeID string) {
	t.Helper()
	if publisher.FindDispatch("PublishSchedule", nodeID) == nil {
		t.Errorf("expected PublishSchedule for node %s, not found", nodeID)
	}
}

// AssertScheduleNotPublished asserts that PublishSchedule was NOT called for the given nodeID.
func AssertScheduleNotPublished(t *testing.T, publisher *CapturingPublisher, nodeID string) {
	t.Helper()
	if publisher.FindDispatch("PublishSchedule", nodeID) != nil {
		t.Errorf("expected NO PublishSchedule for node %s, but found one", nodeID)
	}
}

// AssertSchedulePurged asserts that PurgeSchedule was called for the given nodeID.
func AssertSchedulePurged(t *testing.T, publisher *CapturingPublisher, nodeID string) {
	t.Helper()
	if publisher.FindDispatch("PurgeSchedule", nodeID) == nil {
		t.Errorf("expected PurgeSchedule for node %s, not found", nodeID)
	}
}

// AssertScheduleCount asserts the publisher received exactly N PublishSchedule calls.
func AssertScheduleCount(t *testing.T, publisher *CapturingPublisher, expectedCount int) {
	t.Helper()
	count := publisher.CountMethod("PublishSchedule")
	if count != expectedCount {
		t.Errorf("expected %d PublishSchedule calls, got %d", expectedCount, count)
	}
}

// AssertPurgeAllSchedules asserts that PurgeAllSchedules was called at least once.
func AssertPurgeAllSchedules(t *testing.T, publisher *CapturingPublisher) {
	t.Helper()
	if publisher.CountMethod("PurgeAllSchedules") == 0 {
		t.Error("expected PurgeAllSchedules to be called")
	}
}

// AssertPathNodeError asserts a node's PathEntry has status "error" and Error contains the substring.
func AssertPathNodeError(t *testing.T, exec *entities.WorkflowExecution, nodeID, errorSubstring string) {
	t.Helper()
	for _, entry := range exec.ExecutionPath {
		if entry.NodeID == nodeID {
			if entry.Status != "error" {
				t.Fatalf("path node %q: expected status='error', got %q", nodeID, entry.Status)
			}
			if entry.Error == nil {
				t.Fatalf("path node %q: expected Error to be set, got nil", nodeID)
			}
			if !strings.Contains(*entry.Error, errorSubstring) {
				t.Fatalf("path node %q: expected Error to contain %q, got %q", nodeID, errorSubstring, *entry.Error)
			}
			return
		}
	}
	t.Fatalf("node %q not found in execution path", nodeID)
}
