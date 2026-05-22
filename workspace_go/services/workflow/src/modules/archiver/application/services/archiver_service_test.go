package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"time"

	"workflow/src/modules/archiver/application/di"
	"workflow/src/modules/archiver/application/ports"
	"workflow/src/modules/archiver/domain/repositories"
	runtimeConstants "workflow/src/modules/runtime/application/constants"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	sharedTypes "workflow/src/shared/types"
)

// stateSubject builds the env-prefixed subject the publisher emits for a
// state event. Mirrors the helper in archiver_helpers.go so test fixtures
// stay aligned with the live classification logic.
func stateSubject(status string) string {
	return fmt.Sprintf(runtimeConstants.StatePatternSubject, status)
}

/*
 * MOCKS
 */

type mockArchiveRepo struct {
	mu                     sync.Mutex
	insertedStubs          []repositories.LightweightExecution
	upsertedFull           []*runtimePorts.WorkflowExecution
	waitingUpdates         []repositories.WaitingUpdate
	resumedIDs             []string
	bulkInsertErr          error
	bulkUpsertErr          error
	bulkUpdateWaitingErr   error
	bulkUpdateResumedErr   error
}

func (m *mockArchiveRepo) BulkInsertLightweight(_ context.Context, stubs []repositories.LightweightExecution) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.insertedStubs = append(m.insertedStubs, stubs...)
	return m.bulkInsertErr
}

func (m *mockArchiveRepo) BulkUpsertFull(_ context.Context, executions []*runtimePorts.WorkflowExecution) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.upsertedFull = append(m.upsertedFull, executions...)
	return m.bulkUpsertErr
}

func (m *mockArchiveRepo) BulkUpdateWaiting(_ context.Context, updates []repositories.WaitingUpdate) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.waitingUpdates = append(m.waitingUpdates, updates...)
	return m.bulkUpdateWaitingErr
}

func (m *mockArchiveRepo) BulkUpdateResumed(_ context.Context, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resumedIDs = append(m.resumedIDs, ids...)
	return m.bulkUpdateResumedErr
}

func (m *mockArchiveRepo) FindExecutions(_ context.Context, _ model.Map, _ *model.PaginationOpts) (*model.PaginatedResult[runtimePorts.WorkflowExecution], error) {
	return nil, nil
}

func (m *mockArchiveRepo) FindExecutionById(_ context.Context, _ string) (*runtimePorts.WorkflowExecution, error) {
	return nil, nil
}

type archiverMockKV struct {
	data    map[string]*natsModel.KVEntry
	deleted []string
}

func (m *archiverMockKV) Get(key string) (*natsModel.KVEntry, error) {
	entry, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return entry, nil
}

func (m *archiverMockKV) Put(_ string, _ []byte) (uint64, error)              { return 0, nil }
func (m *archiverMockKV) Create(_ string, _ []byte) (uint64, error)           { return 0, nil }
func (m *archiverMockKV) Update(_ string, _ []byte, _ uint64) (uint64, error) { return 0, nil }
func (m *archiverMockKV) Delete(key string) error {
	m.deleted = append(m.deleted, key)
	return nil
}
func (m *archiverMockKV) Purge(_ string) error                             { return nil }
func (m *archiverMockKV) Keys() ([]string, error)                          { return nil, nil }
func (m *archiverMockKV) Bucket() string                                   { return "test" }

type archiverMockPublisher struct {
	published []natsModel.PublishConfig
}

func (m *archiverMockPublisher) Publish(config natsModel.PublishConfig) error {
	m.published = append(m.published, config)
	return nil
}

/*
 * HELPERS
 */

type mockMongoManager struct{}

func (m *mockMongoManager) GetBackpressureMode() ports.BackpressureMode { return ports.BackpressureNormal }
func (m *mockMongoManager) WriteP99() int64                             { return 0 }
func (m *mockMongoManager) RecordWriteLatency(_ time.Duration)          {}

var testMongoManager ports.MongoManagerPort = &mockMongoManager{}

func newArchiverTestDeps(repo *mockArchiveRepo, kv natsModel.KeyValueStore, pub natsModel.Publisher) di.ArchiverServiceDependenciesInjection {
	return di.ArchiverServiceDependenciesInjection{
		ArchiveRepo:  repo,
		KVStore:      kv,
		Publisher:    pub,
		MongoManager: testMongoManager,
	}
}

func makeTestMessage(subject string, event sharedTypes.StateEvent) *natsModel.Message {
	data, _ := json.Marshal(event)
	msg := natsModel.NewTestMessage(data, 0, &natsModel.TestMessageCallbacks{
		OnAck:    func() error { return nil },
		OnNack:   func(_ error) error { return nil },
		OnReject: func(_ string) error { return nil },
	})
	msg.Subject = subject
	return msg
}

/*
 * TESTS: event classification
 */

func TestEventClassification(t *testing.T) {
	tests := []struct {
		subject   string
		isCreated bool
		isWaiting bool
		isResumed bool
		isTerminal bool
	}{
		{stateSubject("created"), true, false, false, false},
		{stateSubject("waiting"), false, true, false, false},
		{stateSubject("resumed"), false, false, true, false},
		{stateSubject("completed"), false, false, false, true},
		{stateSubject("failed"), false, false, false, true},
		{stateSubject("cancelled"), false, false, false, true},
		{stateSubject("unknown"), false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			if got := isCreatedEvent(tt.subject); got != tt.isCreated {
				t.Fatalf("isCreatedEvent(%q) = %v, want %v", tt.subject, got, tt.isCreated)
			}
			if got := isWaitingEvent(tt.subject); got != tt.isWaiting {
				t.Fatalf("isWaitingEvent(%q) = %v, want %v", tt.subject, got, tt.isWaiting)
			}
			if got := isResumedEvent(tt.subject); got != tt.isResumed {
				t.Fatalf("isResumedEvent(%q) = %v, want %v", tt.subject, got, tt.isResumed)
			}
			if got := isTerminalEvent(tt.subject); got != tt.isTerminal {
				t.Fatalf("isTerminalEvent(%q) = %v, want %v", tt.subject, got, tt.isTerminal)
			}
		})
	}
}

/*
 * TESTS: buildLightweightStub
 */

func TestBuildLightweightStub_Valid(t *testing.T) {
	event := sharedTypes.StateEvent{
		InstanceID:    "exec-uuid-123",
		WorkflowID:    "683f1a2b3c4d5e6f7a8b9c0e",
		OrgID:         "683f1a2b3c4d5e6f7a8b9c0f",
		WorkflowName:  "Test Workflow",
		Status:        "created",
		ActiveNodeIDs: []string{"start-1"},
		Version:       1,
	}

	stub, err := buildLightweightStub(event)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if stub.WorkflowUUID != "exec-uuid-123" {
		t.Fatalf("expected WorkflowUUID 'exec-uuid-123', got %q", stub.WorkflowUUID)
	}
	if stub.WorkflowName != "Test Workflow" {
		t.Fatalf("expected name 'Test Workflow', got %q", stub.WorkflowName)
	}
	if stub.Status != "created" {
		t.Fatalf("expected status 'created', got %q", stub.Status)
	}
	if len(stub.ActiveNodeIDs) != 1 || stub.ActiveNodeIDs[0] != "start-1" {
		t.Fatalf("expected activeNodeIds=['start-1'], got %v", stub.ActiveNodeIDs)
	}
	if stub.OrgID == nil {
		t.Fatal("expected orgId to be set")
	}
}

func TestBuildLightweightStub_InvalidWorkflowID(t *testing.T) {
	event := sharedTypes.StateEvent{
		InstanceID: "exec-uuid-123",
		WorkflowID: "bad-id",
	}

	_, err := buildLightweightStub(event)
	if err == nil {
		t.Fatal("expected error for invalid workflowId")
	}
}

func TestBuildLightweightStub_EmptyOrgID(t *testing.T) {
	event := sharedTypes.StateEvent{
		InstanceID: "exec-uuid-123",
		WorkflowID: "683f1a2b3c4d5e6f7a8b9c0e",
		OrgID:      "",
	}

	stub, err := buildLightweightStub(event)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if stub.OrgID != nil {
		t.Fatal("expected nil orgId for empty string")
	}
}

/*
 * TESTS: buildWaitingUpdate
 */

func TestBuildWaitingUpdate_Valid(t *testing.T) {
	event := sharedTypes.StateEvent{
		InstanceID:    "exec-123",
		Status:        "waiting",
		ActiveNodeIDs: []string{"delay-1"},
	}

	update, err := buildWaitingUpdate(event)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if update.WorkflowUUID != "exec-123" {
		t.Fatalf("expected instanceId 'exec-123', got %q", update.WorkflowUUID)
	}
	if update.Status != "waiting" {
		t.Fatalf("expected status 'waiting', got %q", update.Status)
	}
}

func TestBuildWaitingUpdate_MissingInstanceID(t *testing.T) {
	event := sharedTypes.StateEvent{
		InstanceID: "",
		Status:     "waiting",
	}

	_, err := buildWaitingUpdate(event)
	if err == nil {
		t.Fatal("expected error for missing instanceId")
	}
}

/*
 * TESTS: ProcessStateBatch
 */

func TestProcessStateBatch_CreatedEvents(t *testing.T) {
	repo := &mockArchiveRepo{}
	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{}}
	pub := &archiverMockPublisher{}

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}

	messages := []*natsModel.Message{
		makeTestMessage(stateSubject("created"), sharedTypes.StateEvent{
			InstanceID:    "exec-uuid-1",
			WorkflowID:    "683f1a2b3c4d5e6f7a8b9c0e",
			WorkflowName:  "WF-1",
			Status:        "created",
			ActiveNodeIDs: []string{"start"},
			Version:       1,
		}),
	}

	svc.ProcessStateBatch(messages)

	if len(repo.insertedStubs) != 1 {
		t.Fatalf("expected 1 inserted stub, got %d", len(repo.insertedStubs))
	}
	if repo.insertedStubs[0].WorkflowName != "WF-1" {
		t.Fatalf("expected workflowName 'WF-1', got %q", repo.insertedStubs[0].WorkflowName)
	}
}

func TestProcessStateBatch_WaitingEvents(t *testing.T) {
	repo := &mockArchiveRepo{}
	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{}}
	pub := &archiverMockPublisher{}

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}

	messages := []*natsModel.Message{
		makeTestMessage(stateSubject("waiting"), sharedTypes.StateEvent{
			InstanceID:    "exec-1",
			Status:        "waiting",
			ActiveNodeIDs: []string{"delay-1"},
		}),
	}

	svc.ProcessStateBatch(messages)

	if len(repo.waitingUpdates) != 1 {
		t.Fatalf("expected 1 waiting update, got %d", len(repo.waitingUpdates))
	}
	if repo.waitingUpdates[0].WorkflowUUID != "exec-1" {
		t.Fatalf("expected instanceId 'exec-1', got %q", repo.waitingUpdates[0].WorkflowUUID)
	}
}

func TestProcessStateBatch_ResumedEvents(t *testing.T) {
	repo := &mockArchiveRepo{}
	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{}}
	pub := &archiverMockPublisher{}

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}

	messages := []*natsModel.Message{
		makeTestMessage(stateSubject("resumed"), sharedTypes.StateEvent{
			InstanceID: "exec-1",
			Status:     "running",
		}),
	}

	svc.ProcessStateBatch(messages)

	if len(repo.resumedIDs) != 1 {
		t.Fatalf("expected 1 resumed ID, got %d", len(repo.resumedIDs))
	}
	if repo.resumedIDs[0] != "exec-1" {
		t.Fatalf("expected 'exec-1', got %q", repo.resumedIDs[0])
	}
}

func TestProcessStateBatch_TerminalEvents(t *testing.T) {
	execution := runtimePorts.WorkflowExecution{
		WorkflowUUID: "exec-done",
		Status:       runtimePorts.ExecStatusCompleted,
		WorkflowName: "WF-Done",
	}
	executionData, _ := json.Marshal(execution)

	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{
		"exec.exec-done": {Value: executionData},
	}}
	repo := &mockArchiveRepo{}
	pub := &archiverMockPublisher{}

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}

	messages := []*natsModel.Message{
		makeTestMessage(stateSubject("completed"), sharedTypes.StateEvent{
			InstanceID: "exec-done",
			Status:     "completed",
		}),
	}

	svc.ProcessStateBatch(messages)

	if len(repo.upsertedFull) != 1 {
		t.Fatalf("expected 1 full upsert, got %d", len(repo.upsertedFull))
	}
	if repo.upsertedFull[0].WorkflowName != "WF-Done" {
		t.Fatalf("expected workflowName 'WF-Done', got %q", repo.upsertedFull[0].WorkflowName)
	}
	if len(kv.deleted) != 1 || kv.deleted[0] != "exec.exec-done" {
		t.Fatalf("expected KV key 'exec.exec-done' deleted, got %v", kv.deleted)
	}
}

func TestProcessStateBatch_InvalidJSON(t *testing.T) {
	repo := &mockArchiveRepo{}
	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{}}
	pub := &archiverMockPublisher{}

	rejected := false
	msg := natsModel.NewTestMessage([]byte("not json"), 0, &natsModel.TestMessageCallbacks{
		OnReject: func(_ string) error {
			rejected = true
			return nil
		},
	})
	msg.Subject = stateSubject("created")

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}
	svc.ProcessStateBatch([]*natsModel.Message{msg})

	if !rejected {
		t.Fatal("expected invalid JSON message to be rejected")
	}
	if len(repo.insertedStubs) != 0 {
		t.Fatalf("expected 0 inserts for invalid JSON, got %d", len(repo.insertedStubs))
	}
}

func TestProcessStateBatch_MixedBatch(t *testing.T) {
	execution := runtimePorts.WorkflowExecution{
		WorkflowUUID: "exec-fail-1",
		Status:       runtimePorts.ExecStatusFailed,
		WorkflowName: "WF-Failed",
	}
	executionData, _ := json.Marshal(execution)

	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{
		"exec.exec-fail-1": {Value: executionData},
	}}
	repo := &mockArchiveRepo{}
	pub := &archiverMockPublisher{}

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}

	messages := []*natsModel.Message{
		makeTestMessage(stateSubject("created"), sharedTypes.StateEvent{
			InstanceID:   "exec-uuid-new",
			WorkflowID:   "683f1a2b3c4d5e6f7a8b9c0e",
			WorkflowName: "WF-New",
			Status:       "created",
			Version:      1,
		}),
		makeTestMessage(stateSubject("waiting"), sharedTypes.StateEvent{
			InstanceID:    "exec-wait",
			Status:        "waiting",
			ActiveNodeIDs: []string{"delay-1"},
		}),
		makeTestMessage(stateSubject("resumed"), sharedTypes.StateEvent{
			InstanceID: "exec-resume",
			Status:     "running",
		}),
		makeTestMessage(stateSubject("failed"), sharedTypes.StateEvent{
			InstanceID: "exec-fail-1",
			Status:     "failed",
		}),
	}

	svc.ProcessStateBatch(messages)

	if len(repo.insertedStubs) != 1 {
		t.Fatalf("expected 1 created stub, got %d", len(repo.insertedStubs))
	}
	if len(repo.waitingUpdates) != 1 {
		t.Fatalf("expected 1 waiting update, got %d", len(repo.waitingUpdates))
	}
	if len(repo.resumedIDs) != 1 {
		t.Fatalf("expected 1 resumed ID, got %d", len(repo.resumedIDs))
	}
	if len(repo.upsertedFull) != 1 {
		t.Fatalf("expected 1 terminal upsert, got %d", len(repo.upsertedFull))
	}
}

func TestProcessStateBatch_BulkInsertError_NACKs(t *testing.T) {
	repo := &mockArchiveRepo{bulkInsertErr: fmt.Errorf("db write failed")}
	kv := &archiverMockKV{data: map[string]*natsModel.KVEntry{}}
	pub := &archiverMockPublisher{}

	nacked := false
	data, _ := json.Marshal(sharedTypes.StateEvent{
		InstanceID:   "exec-uuid-fail",
		WorkflowID:   "683f1a2b3c4d5e6f7a8b9c0e",
		WorkflowName: "WF-Fail",
		Status:       "created",
		Version:      1,
	})
	msg := natsModel.NewTestMessage(data, 0, &natsModel.TestMessageCallbacks{
		OnNack: func(_ error) error {
			nacked = true
			return nil
		},
	})
	msg.Subject = stateSubject("created")

	svc := &ArchiverService{deps: newArchiverTestDeps(repo, kv, pub)}
	svc.ProcessStateBatch([]*natsModel.Message{msg})

	if !nacked {
		t.Fatal("expected NACK when bulk insert fails")
	}
}
