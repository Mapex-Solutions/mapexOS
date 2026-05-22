package services

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"triggers/src/bootstrap"
	"triggers/src/modules/events/application/di"
	"triggers/src/modules/events/application/ports"
	triggerDtos "triggers/src/modules/triggers/application/dtos"

	"github.com/prometheus/client_golang/prometheus"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	triggers "github.com/Mapex-Solutions/MapexOS/contracts/services/triggers/triggers"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
)

/**
 * Mock Implementations
 */

// MockTriggerServicePort is a mock implementation of TriggerServicePort
type MockTriggerServicePort struct {
	GetTriggerByIdFunc func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error)
}

func (m *MockTriggerServicePort) CreateTrigger(ctx context.Context, requestContext *reqCtx.RequestContext, dto *triggerDtos.CreateTriggerDto) (*triggerDtos.TriggerResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *MockTriggerServicePort) GetTriggerById(ctx context.Context, triggerId *string, metrics ...*common.CacheMetrics) (*triggerDtos.TriggerResponse, error) {
	if m.GetTriggerByIdFunc != nil {
		return m.GetTriggerByIdFunc(ctx, triggerId)
	}
	return nil, errors.New("not found")
}

func (m *MockTriggerServicePort) UpdateTriggerById(ctx context.Context, requestContext *reqCtx.RequestContext, triggerId *string, dto *triggerDtos.UpdateTriggerDto) (*triggerDtos.TriggerResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *MockTriggerServicePort) GetTriggers(ctx context.Context, requestContext *reqCtx.RequestContext, query *triggerDtos.TriggerQueryDto) (*model.PaginatedResult[triggerDtos.TriggerResponse], error) {
	return nil, errors.New("not implemented")
}

func (m *MockTriggerServicePort) DeleteTriggerById(ctx context.Context, triggerId *string) (map[string]bool, error) {
	return nil, errors.New("not implemented")
}

func (m *MockTriggerServicePort) CountTriggers(ctx context.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	return 0, nil
}

// MockExecutorRegistry is a mock implementation of ExecutorRegistry
type MockExecutorRegistry struct {
	GetExecutorFunc func(triggerType string) (ports.TriggerExecutor, bool)
}

func (m *MockExecutorRegistry) GetExecutor(triggerType string) (ports.TriggerExecutor, bool) {
	if m.GetExecutorFunc != nil {
		return m.GetExecutorFunc(triggerType)
	}
	return nil, false
}

// MockTriggerExecutor is a mock implementation of TriggerExecutor
type MockTriggerExecutor struct {
	ExecuteFunc func(ctx context.Context, config map[string]interface{}) error
	TypeName    string
}

func (m *MockTriggerExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, config)
	}
	return nil
}

func (m *MockTriggerExecutor) GetType() string {
	return m.TypeName
}


// MockCorePublisher is a mock implementation of natsModel.CorePublisher
type MockCorePublisher struct {
	PublishCoreFunc   func(config natsModel.PublishCoreConfig) error
	FlushConnectionFunc func() error
}

func (m *MockCorePublisher) PublishCore(config natsModel.PublishCoreConfig) error {
	if m.PublishCoreFunc != nil {
		return m.PublishCoreFunc(config)
	}
	return nil
}

func (m *MockCorePublisher) FlushConnection() error {
	if m.FlushConnectionFunc != nil {
		return m.FlushConnectionFunc()
	}
	return nil
}

/**
 * Helper Functions
 */

func createTestMetrics() *bootstrap.TriggerMetrics {
	reg := metrics.NewRegistry("triggers_test")

	return &bootstrap.TriggerMetrics{
		Registry: reg,
		TriggersProcessed: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "trigger", Name: "processed_total", Help: "test",
		}, []string{"status"}),
		TriggerProcessingDuration: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "trigger", Name: "processing_duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}),
		TriggersBatchSize: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "trigger", Name: "batch_size", Help: "test",
			Buckets: prometheus.DefBuckets,
		}),
		MessagesTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "message", Name: "total", Help: "test",
		}, []string{"result"}),
		ExecutorDuration: reg.NewHistogramVec(metrics.HistogramOpts{
			Subsystem: "executor", Name: "duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}, []string{"type"}),
		ExecutorTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "executor", Name: "total", Help: "test",
		}, []string{"type", "status"}),
		EventsPublished: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "event", Name: "published_total", Help: "test",
		}, []string{"status"}),
		PublishDuration: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "publish", Name: "duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}),
		TriggerCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "trigger", Name: "cache_total", Help: "test",
		}, []string{"result"}),
		PlaceholderResolutions: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "placeholder", Name: "resolutions_total", Help: "test",
		}, []string{"status"}),
	}
}

func createTestEventService(triggerSvc *MockTriggerServicePort, registry *MockExecutorRegistry) *EventService {
	return &EventService{
		deps: di.EventServiceDependenciesInjection{
			TriggerService:   triggerSvc,
			ExecutorRegistry: registry,
			NatsBus:          &MockCorePublisher{},
			Metrics:          createTestMetrics(),
		},
		workers: 4,
	}
}

func createTestTriggerResponse() *triggerDtos.TriggerResponse {
	name := "Test HTTP Trigger"
	triggerType := "http"
	category := "technical"
	enabled := true

	return &triggerDtos.TriggerResponse{
		Name:        &name,
		TriggerType: &triggerType,
		Category:    &category,
		Enabled:     &enabled,
		Config: &triggers.TriggerConfig{
			Http: &triggers.HttpConfig{
				Endpoint: "https://api.example.com/webhook",
				Method:   "POST",
				Body: map[string]interface{}{
					"message": "Alert: {{message}}",
				},
			},
		},
	}
}

func createTriggerExecuteEventJSON(triggerId string, payload map[string]interface{}) []byte {
	event := map[string]interface{}{
		"triggerId":   triggerId,
		"executionId": "exec-123",
		"payload":     payload,
		"orgId":       "507f1f77bcf86cd799439011",
		"pathKey":     "mapex.vendor.customer",
		"timestamp":   "2024-01-01T00:00:00Z",
	}
	data, _ := json.Marshal(event)
	return data
}

/**
 * ProcessTriggerExecution Tests (V1 Legacy)
 */

func TestEventService_ProcessTriggerExecution_Success(t *testing.T) {
	executorCalled := false

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			executorCalled = true
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			if triggerType == "http" {
				return mockExecutor, true
			}
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{
		"message": "Test alert",
	})

	err := service.ProcessTriggerExecution(eventData)

	if err != nil {
		t.Fatalf("ProcessTriggerExecution() unexpected error: %v", err)
	}

	if !executorCalled {
		t.Error("ProcessTriggerExecution() should call the executor")
	}
}

func TestEventService_ProcessTriggerExecution_InvalidJSON(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	invalidJSON := []byte("not valid json")

	err := service.ProcessTriggerExecution(invalidJSON)

	if err == nil {
		t.Fatal("ProcessTriggerExecution() expected error for invalid JSON, got nil")
	}
}

func TestEventService_ProcessTriggerExecution_TriggerNotFound(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return nil, errors.New("trigger not found")
		},
	}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("nonexistent", map[string]interface{}{})

	err := service.ProcessTriggerExecution(eventData)

	if err == nil {
		t.Fatal("ProcessTriggerExecution() expected error for non-existent trigger, got nil")
	}
}

func TestEventService_ProcessTriggerExecution_TriggerDisabled(t *testing.T) {
	disabledTrigger := createTestTriggerResponse()
	enabled := false
	disabledTrigger.Enabled = &enabled

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return disabledTrigger, nil
		},
	}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})

	// Disabled triggers should NOT return an error - they should be silently skipped
	err := service.ProcessTriggerExecution(eventData)

	if err != nil {
		t.Errorf("ProcessTriggerExecution() should not error for disabled trigger, got: %v", err)
	}
}

func TestEventService_ProcessTriggerExecution_ExecutorNotFound(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			trigger := createTestTriggerResponse()
			triggerType := "unknown_type"
			trigger.TriggerType = &triggerType
			return trigger, nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return nil, false // No executor found
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})

	err := service.ProcessTriggerExecution(eventData)

	if err == nil {
		t.Fatal("ProcessTriggerExecution() expected error when executor not found, got nil")
	}
}

func TestEventService_ProcessTriggerExecution_ExecutorFails(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			return errors.New("HTTP request failed")
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return mockExecutor, true
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})

	err := service.ProcessTriggerExecution(eventData)

	if err == nil {
		t.Fatal("ProcessTriggerExecution() expected error when executor fails, got nil")
	}
}

/**
 * ProcessTriggerExecutionBatch Tests (V2 Recommended)
 */

// createTestMessage creates a natsModel.Message via NewTestMessage with tracking callbacks.
// Returns the message and a tracker struct to verify Ack/Nack/Reject calls.
type messageTracker struct {
	AckCalled    bool
	NackCalled   bool
	NackErr      error
	RejectCalled bool
	RejectReason string
}

func createTestMessage(data []byte, index int) (*natsModel.Message, *messageTracker) {
	tracker := &messageTracker{}
	msg := natsModel.NewTestMessage(data, index, &natsModel.TestMessageCallbacks{
		OnAck: func() error {
			tracker.AckCalled = true
			return nil
		},
		OnNack: func(err error) error {
			tracker.NackCalled = true
			tracker.NackErr = err
			return nil
		},
		OnReject: func(reason string) error {
			tracker.RejectCalled = true
			tracker.RejectReason = reason
			return nil
		},
	})
	return msg, tracker
}

func TestEventService_ProcessTriggerExecutionBatch_EmptyBatch(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	// Empty batch should return nil without processing
	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should not error on empty batch, got: %v", err)
	}
}

func TestEventService_ProcessTriggerExecutionBatch_InvalidJSON(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	msg, tracker := createTestMessage([]byte("not valid json"), 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil (handles errors internally), got: %v", err)
	}

	if !tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should Reject message with invalid JSON")
	}

	if tracker.AckCalled || tracker.NackCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Reject (not Ack or Nack) for invalid JSON")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_TriggerNotFound(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return nil, errors.New("trigger not found")
		},
	}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("nonexistent", map[string]interface{}{})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	if !tracker.NackCalled {
		t.Error("ProcessTriggerExecutionBatch() should Nack when trigger not found (retry with backoff)")
	}

	if tracker.AckCalled || tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Nack (not Ack or Reject) for trigger not found")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_TriggerDisabled(t *testing.T) {
	disabledTrigger := createTestTriggerResponse()
	enabled := false
	disabledTrigger.Enabled = &enabled

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return disabledTrigger, nil
		},
	}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	if !tracker.AckCalled {
		t.Error("ProcessTriggerExecutionBatch() should Ack disabled triggers (skip silently, don't retry)")
	}

	if tracker.NackCalled || tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Ack (not Nack or Reject) for disabled triggers")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_ExecutorNotFound(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			trigger := createTestTriggerResponse()
			triggerType := "unknown_type"
			trigger.TriggerType = &triggerType
			return trigger, nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{
		"message": "Test alert",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	if !tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should Reject when executor not found (fatal, no retry)")
	}

	if tracker.AckCalled || tracker.NackCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Reject (not Ack or Nack) for missing executor")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_ExecutionFailure(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			return errors.New("HTTP request failed")
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return mockExecutor, true
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	// Provide payload that resolves the {{message}} placeholder in the trigger config
	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{
		"message": "Test alert",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	if !tracker.NackCalled {
		t.Error("ProcessTriggerExecutionBatch() should Nack when execution fails (retry with backoff)")
	}

	if tracker.AckCalled || tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Nack (not Ack or Reject) for execution failure")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_Success(t *testing.T) {
	executorCalled := false

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			executorCalled = true
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			if triggerType == "http" {
				return mockExecutor, true
			}
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{
		"message": "Test alert",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	if !executorCalled {
		t.Error("ProcessTriggerExecutionBatch() should call the executor")
	}

	if !tracker.AckCalled {
		t.Error("ProcessTriggerExecutionBatch() should Ack on success")
	}

	if tracker.NackCalled || tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Ack (not Nack or Reject) on success")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_MultipleMixedMessages(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			if *triggerId == "not-found" {
				return nil, errors.New("trigger not found")
			}
			if *triggerId == "disabled" {
				trigger := createTestTriggerResponse()
				enabled := false
				trigger.Enabled = &enabled
				return trigger, nil
			}
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			if triggerType == "http" {
				return mockExecutor, true
			}
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	// Message 0: Invalid JSON -> Reject
	msg0, tracker0 := createTestMessage([]byte("invalid"), 0)

	// Message 1: Trigger not found -> Nack
	eventData1 := createTriggerExecuteEventJSON("not-found", map[string]interface{}{})
	msg1, tracker1 := createTestMessage(eventData1, 1)

	// Message 2: Trigger disabled -> Ack (skip)
	eventData2 := createTriggerExecuteEventJSON("disabled", map[string]interface{}{})
	msg2, tracker2 := createTestMessage(eventData2, 2)

	// Message 3: Success -> Ack
	eventData3 := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{"message": "test"})
	msg3, tracker3 := createTestMessage(eventData3, 3)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg0, msg1, msg2, msg3})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	// Message 0: Invalid JSON -> Reject
	if !tracker0.RejectCalled {
		t.Error("Message 0 (invalid JSON) should be Rejected")
	}

	// Message 1: Trigger not found -> Nack
	if !tracker1.NackCalled {
		t.Error("Message 1 (trigger not found) should be Nacked")
	}

	// Message 2: Trigger disabled -> Ack
	if !tracker2.AckCalled {
		t.Error("Message 2 (trigger disabled) should be Acked")
	}

	// Message 3: Success -> Ack
	if !tracker3.AckCalled {
		t.Error("Message 3 (success) should be Acked")
	}
}

/**
 * Placeholder Resolution Integration Tests
 */

func TestEventService_ProcessTriggerExecution_PlaceholderResolution(t *testing.T) {
	var capturedConfig map[string]interface{}

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			trigger := createTestTriggerResponse()
			trigger.Config = &triggers.TriggerConfig{
				Http: &triggers.HttpConfig{
					Endpoint: "https://api.example.com/webhook",
					Method:   "POST",
					Body: map[string]interface{}{
						"alertMessage": "Sensor {{sensorId}} reported {{temperature}}°C",
						"severity":     "{{severity}}",
					},
				},
			}
			return trigger, nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			capturedConfig = config
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return mockExecutor, true
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	payload := map[string]interface{}{
		"sensorId":    "SENSOR-001",
		"temperature": 42.5,
		"severity":    "high",
	}
	eventData := createTriggerExecuteEventJSON("trigger123", payload)

	err := service.ProcessTriggerExecution(eventData)

	if err != nil {
		t.Fatalf("ProcessTriggerExecution() unexpected error: %v", err)
	}

	// Verify that placeholders were resolved
	if capturedConfig == nil {
		t.Fatal("Executor should have received config")
	}

	httpConfig, ok := capturedConfig["http"].(map[string]interface{})
	if !ok {
		t.Fatal("Config should contain 'http' field")
	}

	body, ok := httpConfig["body"].(map[string]interface{})
	if !ok {
		t.Fatal("HTTP config should contain 'body' field")
	}

	alertMessage, ok := body["alertMessage"].(string)
	if !ok {
		t.Fatal("Body should contain 'alertMessage' field")
	}

	expectedMessage := "Sensor SENSOR-001 reported 42.5°C"
	if alertMessage != expectedMessage {
		t.Errorf("Placeholder resolution failed: got %q, want %q", alertMessage, expectedMessage)
	}
}

/**
 * Edge Cases
 */

func TestEventService_ProcessTriggerExecution_NilEnabledField(t *testing.T) {
	trigger := createTestTriggerResponse()
	trigger.Enabled = nil // nil instead of false

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return trigger, nil
		},
	}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})

	// Nil Enabled should be treated as disabled - no error, just skip
	err := service.ProcessTriggerExecution(eventData)

	if err != nil {
		t.Errorf("ProcessTriggerExecution() should not error for nil Enabled field, got: %v", err)
	}
}

func TestEventService_ProcessTriggerExecution_NilConfig(t *testing.T) {
	trigger := createTestTriggerResponse()
	trigger.Config = nil // nil Config

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return trigger, nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			if triggerType == "http" {
				return mockExecutor, true
			}
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})

	// Nil Config should be handled gracefully (TriggerConfigToMap returns empty map for nil)
	err := service.ProcessTriggerExecution(eventData)

	if err != nil {
		t.Errorf("ProcessTriggerExecution() should not error for nil Config, got: %v", err)
	}
}

func TestEventService_ProcessTriggerExecution_NilTriggerType(t *testing.T) {
	trigger := createTestTriggerResponse()
	trigger.TriggerType = nil // nil TriggerType

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return trigger, nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})

	// Nil TriggerType should cause executor lookup to fail
	err := service.ProcessTriggerExecution(eventData)

	if err == nil {
		t.Fatal("ProcessTriggerExecution() expected error for nil TriggerType, got nil")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_NilConfig_Reject(t *testing.T) {
	trigger := createTestTriggerResponse()
	trigger.Config = nil

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return trigger, nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			if triggerType == "http" {
				return mockExecutor, true
			}
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	// Nil Config produces an empty map via TriggerConfigToMap — executor receives empty config.
	// The executor mock accepts any config, so this should Ack.
	if !tracker.AckCalled {
		t.Error("ProcessTriggerExecutionBatch() should Ack when nil Config produces empty map (executor succeeds)")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_ConcurrentSuccess(t *testing.T) {
	var mu sync.Mutex
	executionCount := 0

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			mu.Lock()
			executionCount++
			mu.Unlock()
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			if triggerType == "http" {
				return mockExecutor, true
			}
			return nil, false
		},
	}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	// Create 10 messages to exercise concurrency
	batchSize := 10
	messages := make([]*natsModel.Message, batchSize)
	trackers := make([]*messageTracker, batchSize)

	for i := 0; i < batchSize; i++ {
		eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{
			"message": "Test alert",
		})
		messages[i], trackers[i] = createTestMessage(eventData, i)
	}

	err := service.ProcessTriggerExecutionBatch(messages)

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	// All messages should be Acked
	for i, tracker := range trackers {
		if !tracker.AckCalled {
			t.Errorf("Message %d should be Acked", i)
		}
		if tracker.NackCalled || tracker.RejectCalled {
			t.Errorf("Message %d should only be Acked (not Nack or Reject)", i)
		}
	}

	// All executors should have been called
	mu.Lock()
	defer mu.Unlock()
	if executionCount != batchSize {
		t.Errorf("Executor should be called %d times, got %d", batchSize, executionCount)
	}
}

func TestEventService_ProcessTriggerExecutionBatch_NilEnabledField_Ack(t *testing.T) {
	trigger := createTestTriggerResponse()
	trigger.Enabled = nil // nil instead of false

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return trigger, nil
		},
	}
	mockRegistry := &MockExecutorRegistry{}

	service := createTestEventService(mockTriggerSvc, mockRegistry)

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil, got: %v", err)
	}

	// Nil Enabled should be treated as disabled -> Ack (skip, don't retry)
	if !tracker.AckCalled {
		t.Error("ProcessTriggerExecutionBatch() should Ack when Enabled is nil (treat as disabled)")
	}

	if tracker.NackCalled || tracker.RejectCalled {
		t.Error("ProcessTriggerExecutionBatch() should only Ack (not Nack or Reject) for nil Enabled")
	}
}

func TestEventService_ProcessTriggerExecutionBatch_FlushConnectionError(t *testing.T) {
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTestTriggerResponse(), nil
		},
	}

	mockExecutor := &MockTriggerExecutor{
		TypeName: "http",
		ExecuteFunc: func(ctx context.Context, config map[string]interface{}) error {
			return nil
		},
	}

	mockRegistry := &MockExecutorRegistry{
		GetExecutorFunc: func(triggerType string) (ports.TriggerExecutor, bool) {
			return mockExecutor, true
		},
	}

	service := &EventService{
		deps: di.EventServiceDependenciesInjection{
			TriggerService:   mockTriggerSvc,
			ExecutorRegistry: mockRegistry,
			NatsBus: &MockCorePublisher{
				FlushConnectionFunc: func() error {
					return errors.New("NATS flush connection failed")
				},
			},
			Metrics: createTestMetrics(),
		},
		workers: 4,
	}

	eventData := createTriggerExecuteEventJSON("trigger123", map[string]interface{}{
		"message": "Test alert",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})

	// Batch should still return nil — flush error is logged but not propagated
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() should return nil even with flush error, got: %v", err)
	}

	// Message should still be Acked — execution succeeded
	if !tracker.AckCalled {
		t.Error("ProcessTriggerExecutionBatch() should Ack even when flush fails (execution succeeded)")
	}
}
