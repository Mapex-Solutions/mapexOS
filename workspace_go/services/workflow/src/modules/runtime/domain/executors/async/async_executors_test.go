package async

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	defEntities "workflow/src/modules/definitions/domain/entities"
	"workflow/src/modules/runtime/domain/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/*
 * MOCKS
 */

type mockValueResolver struct {
	value interface{}
	err   error
}

func (m *mockValueResolver) Resolve(
	_ *defEntities.FieldValue,
	_ map[string]interface{},
	_ map[string]interface{},
	_ map[string]interface{},
	_ map[string]interface{},
) (interface{}, error) {
	return m.value, m.err
}

func (m *mockValueResolver) BuildDescription(_ *defEntities.FieldValue) string {
	return ""
}

func newTestContext(nodeType string, config interface{}) *entities.NodeExecutionContext {
	return &entities.NodeExecutionContext{
		InstanceID:     bson.NewObjectID(),
		State:          make(map[string]interface{}),
		EventPayload:   make(map[string]interface{}),
		NodeOutputs:    make(map[string]interface{}),
		NodeStates:     make(map[string]map[string]interface{}),
		ExternalInputs: make(map[string]interface{}),
		NodeID:         "node-1",
		NodeType:       nodeType,
		ParsedConfig:   config,
		Timezone:       "UTC",
	}
}

/*
 * DELAY EXECUTOR
 */

func TestDelayExecutor(t *testing.T) {
	tests := []struct {
		name         string
		config       *entities.DelayNodeConfig
		wantErr      bool
		wantWaitType string
		minExpiry    time.Duration
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name:         "10 seconds delay",
			config:       &entities.DelayNodeConfig{Duration: 10, Unit: "s"},
			wantWaitType: "timer",
			minExpiry:    9 * time.Second,
		},
		{
			name:         "5 minutes delay",
			config:       &entities.DelayNodeConfig{Duration: 5, Unit: "m"},
			wantWaitType: "timer",
			minExpiry:    4 * time.Minute,
		},
		{
			name:         "2 hours delay",
			config:       &entities.DelayNodeConfig{Duration: 2, Unit: "h"},
			wantWaitType: "timer",
			minExpiry:    1 * time.Hour,
		},
		{
			name:         "1 day delay",
			config:       &entities.DelayNodeConfig{Duration: 1, Unit: "d"},
			wantWaitType: "timer",
			minExpiry:    23 * time.Hour,
		},
		{
			name:    "unknown unit returns error",
			config:  &entities.DelayNodeConfig{Duration: 5, Unit: "x"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewDelayExecutor()
			before := time.Now()

			result, err := executor.Execute(context.Background(), newTestContext("core/delay", tt.config))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
			if result.NodeState["waitType"] != tt.wantWaitType {
				t.Fatalf("expected waitType %s, got %v", tt.wantWaitType, result.NodeState["waitType"])
			}
			expiresAt, ok := result.NodeState["expiresAt"].(time.Time)
			if !ok {
				t.Fatal("expected expiresAt as time.Time")
			}

			elapsed := expiresAt.Sub(before)
			if elapsed < tt.minExpiry {
				t.Fatalf("expected expiry >= %v from now, got %v", tt.minExpiry, elapsed)
			}

			if result.OutputHandles[0] != "out" {
				t.Fatalf("expected handle out, got %v", result.OutputHandles)
			}
		})
	}
}

/*
 * WAIT SIGNAL EXECUTOR
 */

func TestWaitSignalExecutor(t *testing.T) {
	tests := []struct {
		name    string
		config  *entities.WaitSignalNodeConfig
		wantErr bool
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name: "valid config",
			config: &entities.WaitSignalNodeConfig{
				SignalName: "approval",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewWaitSignalExecutor()

			result, err := executor.Execute(context.Background(), newTestContext("core/wait_signal", tt.config))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
			if result.NodeState["waitType"] != "signal" {
				t.Fatalf("expected waitType signal, got %v", result.NodeState["waitType"])
			}
			if result.NodeState["signalName"] != tt.config.SignalName {
				t.Fatalf("expected signal %s, got %v", tt.config.SignalName, result.NodeState["signalName"])
			}
			if result.NodeState["expiresAt"] == nil {
				t.Fatal("expected expiresAt in NodeState, got nil")
			}
		})
	}
}

/*
 * CODE EXECUTOR
 */

func TestCodeExecutor(t *testing.T) {
	tests := []struct {
		name    string
		config  *entities.CodeNodeConfig
		wantErr bool
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name: "valid script",
			config: &entities.CodeNodeConfig{
				Script:  "return 42;",
				Timeout: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewCodeExecutor()

			ctx := newTestContext("core/code", tt.config)
			result, err := executor.Execute(context.Background(), ctx)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
			if result.NodeState["waitType"] != "callback" {
				t.Fatalf("expected waitType callback, got %v", result.NodeState["waitType"])
			}
			if _, exists := result.NodeState["script"]; exists {
				t.Fatal("NodeState should NOT contain 'script' key")
			}
			if result.NodeState["timeout"] != tt.config.Timeout {
				t.Fatalf("expected timeout in NodeState, got %v", result.NodeState["timeout"])
			}
		})
	}
}

/*
 * SUBWORKFLOW EXECUTOR
 */

func TestSubworkflowExecutor(t *testing.T) {
	tests := []struct {
		name          string
		config        *entities.SubworkflowNodeConfig
		depth         int
		resolverValue interface{}
		resolverErr   error
		wantErr       bool
		wantErrType   error
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name: "valid config with input mappings",
			config: &entities.SubworkflowNodeConfig{
				WorkflowID:    "wf-123",
				WorkflowName:  "Child Workflow",
				ExecutionMode: "sync",
				InputMappings: []entities.InputMapping{
					{ChildParamName: "param1", Value: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "hello"}},
				},
			},
			resolverValue: "hello",
		},
		{
			name: "max depth exceeded",
			config: &entities.SubworkflowNodeConfig{
				WorkflowID: "wf-123",
			},
			depth:       10,
			wantErr:     true,
			wantErrType: entities.ErrMaxSubworkflowDepth,
		},
		{
			name: "resolver error propagates",
			config: &entities.SubworkflowNodeConfig{
				WorkflowID: "wf-123",
				InputMappings: []entities.InputMapping{
					{ChildParamName: "param1", Value: defEntities.FieldValue{Type: defEntities.FieldValueState, Value: "missing"}},
				},
			},
			resolverErr: errors.New("field not found"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockValueResolver{value: tt.resolverValue, err: tt.resolverErr}
			executor := NewSubworkflowExecutor(resolver)

			ctx := newTestContext("core/subworkflow", tt.config)
			ctx.Depth = tt.depth

			result, err := executor.Execute(context.Background(), ctx)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Fatalf("expected error %v, got %v", tt.wantErrType, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
			if result.NodeState["waitType"] != "callback" {
				t.Fatalf("expected waitType callback, got %v", result.NodeState["waitType"])
			}
			if result.NodeState["workflowId"] != tt.config.WorkflowID {
				t.Fatalf("expected workflowId in NodeState")
			}
			if result.NodeState["depth"] != 1 {
				t.Fatalf("expected depth=1, got %v", result.NodeState["depth"])
			}
		})
	}
}

/*
 * TRIGGER EVENT EXECUTOR
 */

func TestTriggerEventExecutor(t *testing.T) {
	tests := []struct {
		name          string
		config        *entities.TriggerEventNodeConfig
		resolverValue interface{}
		resolverErr   error
		wantErr       bool
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name: "empty eventType returns error",
			config: &entities.TriggerEventNodeConfig{
				EventType: "",
			},
			wantErr: true,
		},
		{
			name: "valid config with payload",
			config: &entities.TriggerEventNodeConfig{
				EventType: "user.created",
				PayloadMapping: []entities.TriggerPayloadField{
					{Key: "name", Value: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "John"}},
				},
			},
			resolverValue: "John",
		},
		{
			name: "resolver error propagates",
			config: &entities.TriggerEventNodeConfig{
				EventType: "user.created",
				PayloadMapping: []entities.TriggerPayloadField{
					{Key: "name", Value: defEntities.FieldValue{Type: defEntities.FieldValueState, Value: "missing"}},
				},
			},
			resolverErr: errors.New("resolve error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockValueResolver{value: tt.resolverValue, err: tt.resolverErr}
			executor := NewTriggerEventExecutor(resolver)

			ctx := newTestContext("core/trigger_event", tt.config)
			result, err := executor.Execute(context.Background(), ctx)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
			if result.NodeState["waitType"] != "callback" {
				t.Fatalf("expected waitType callback, got %v", result.NodeState["waitType"])
			}
			if result.NodeState["eventType"] != tt.config.EventType {
				t.Fatalf("expected eventType in NodeState")
			}

			payload, ok := result.NodeState["payload"].(map[string]interface{})
			if !ok {
				t.Fatal("expected payload map in NodeState")
			}
			if len(tt.config.PayloadMapping) > 0 {
				if payload["name"] != tt.resolverValue {
					t.Fatalf("expected payload[name] = %v, got %v", tt.resolverValue, payload["name"])
				}
			}
		})
	}
}

/*
 * PLUGIN EXECUTOR — resolveFieldValue / resolveConfigFieldValues
 */

func newPluginExecutorWithResolver(resolver *mockValueResolver) *PluginExecutor {
	return &PluginExecutor{
		resolver: resolver,
	}
}

func TestPluginResolveFieldValue_ErrorPropagation(t *testing.T) {
	resolver := &mockValueResolver{value: nil, err: errors.New("field not found: 'list' in path 'item.list'")}
	pe := newPluginExecutorWithResolver(resolver)
	ctx := newTestContext("telegram/message", nil)

	val, err := pe.resolveFieldValue(ctx, map[string]interface{}{
		"type": "nodeOutput", "value": "item.list", "nodeId": "loop1",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if val != nil {
		t.Fatalf("expected nil value on error, got %v", val)
	}
	if !strings.Contains(err.Error(), "item.list") {
		t.Fatalf("expected error to contain 'item.list', got %q", err.Error())
	}
}

func TestPluginResolveFieldValue_PassthroughPlainValues(t *testing.T) {
	resolver := &mockValueResolver{value: nil, err: errors.New("should not be called")}
	pe := newPluginExecutorWithResolver(resolver)
	ctx := newTestContext("telegram/message", nil)

	// String
	val, err := pe.resolveFieldValue(ctx, "hello")
	if err != nil {
		t.Fatalf("unexpected error for string: %v", err)
	}
	if val != "hello" {
		t.Fatalf("expected 'hello', got %v", val)
	}

	// Number
	val, err = pe.resolveFieldValue(ctx, 42)
	if err != nil {
		t.Fatalf("unexpected error for number: %v", err)
	}
	if val != 42 {
		t.Fatalf("expected 42, got %v", val)
	}

	// Nil
	val, err = pe.resolveFieldValue(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error for nil: %v", err)
	}
	if val != nil {
		t.Fatalf("expected nil, got %v", val)
	}
}

func TestPluginResolveFieldValue_PassthroughNonFieldSource(t *testing.T) {
	resolver := &mockValueResolver{value: nil, err: errors.New("should not be called")}
	pe := newPluginExecutorWithResolver(resolver)
	ctx := newTestContext("telegram/message", nil)

	input := map[string]interface{}{"type": "nodeOutput"}
	val, err := pe.resolveFieldValue(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val == nil {
		t.Fatal("expected passthrough value, got nil")
	}
}

func TestPluginResolveFieldValue_FetchOptionsAsLiteral(t *testing.T) {
	resolver := &mockValueResolver{value: "5847825355", err: nil}
	pe := newPluginExecutorWithResolver(resolver)
	ctx := newTestContext("telegram/message", nil)

	val, err := pe.resolveFieldValue(ctx, map[string]interface{}{
		"type": "fetchOptions", "value": "5847825355",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "5847825355" {
		t.Fatalf("expected '5847825355', got %v", val)
	}
}

func TestPluginResolveConfigFieldValues_ErrorIncludesKey(t *testing.T) {
	resolver := &mockValueResolver{value: nil, err: errors.New("field not found")}
	pe := newPluginExecutorWithResolver(resolver)
	ctx := newTestContext("telegram/message", nil)

	result, err := pe.resolveConfigFieldValues(ctx, map[string]interface{}{
		"text": map[string]interface{}{"type": "nodeOutput", "value": "item.list", "nodeId": "loop1"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil map on error, got %v", result)
	}
	if !strings.Contains(err.Error(), "text") {
		t.Fatalf("expected error to contain field key 'text', got %q", err.Error())
	}
}

func TestPluginResolveConfigFieldValues_SuccessPassthrough(t *testing.T) {
	resolver := &mockValueResolver{value: "resolved_value", err: nil}
	pe := newPluginExecutorWithResolver(resolver)
	ctx := newTestContext("telegram/message", nil)

	result, err := pe.resolveConfigFieldValues(ctx, map[string]interface{}{
		"chatId": map[string]interface{}{"type": "literal", "value": "123"},
		"plain":  "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["chatId"] != "resolved_value" {
		t.Fatalf("expected resolved_value for chatId, got %v", result["chatId"])
	}
	if result["plain"] != "hello" {
		t.Fatalf("expected 'hello' for plain, got %v", result["plain"])
	}
}
