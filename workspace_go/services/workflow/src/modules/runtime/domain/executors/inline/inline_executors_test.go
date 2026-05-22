package inline

import (
	"context"
	"testing"

	defEntities "workflow/src/modules/definitions/domain/entities"
	"workflow/src/modules/runtime/domain/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/*
 * MOCKS
 */

type mockConditionEvaluator struct {
	result bool
	err    error
}

func (m *mockConditionEvaluator) EvaluateGroup(
	_ *defEntities.ConditionGroup,
	_ string,
	_ map[string]interface{},
	_ map[string]interface{},
	_ map[string]interface{},
	_ map[string]interface{},
) (bool, error) {
	return m.result, m.err
}

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
 * START EXECUTOR
 */

func TestStartExecutor(t *testing.T) {
	executor := NewStartExecutor()

	if executor.NodeType() != "core/start" {
		t.Fatalf("expected NodeType core/start, got %s", executor.NodeType())
	}

	result, err := executor.Execute(context.Background(), newTestContext("core/start", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.OutputHandles) != 1 || result.OutputHandles[0] != "out" {
		t.Fatalf("expected [out], got %v", result.OutputHandles)
	}
}

/*
 * END EXECUTOR
 */

func TestEndExecutor(t *testing.T) {
	tests := []struct {
		name           string
		config         *entities.EndNodeConfig
		resolverValue  interface{}
		wantHandles    []string
		wantError      bool
		wantErrorCode  string
	}{
		{
			name:        "nil config returns empty handles",
			config:      nil,
			wantHandles: []string{},
		},
		{
			name:        "no error terminates normally",
			config:      &entities.EndNodeConfig{TerminateWithError: false},
			wantHandles: []string{},
		},
		{
			name: "terminate with error",
			config: &entities.EndNodeConfig{
				TerminateWithError: true,
				ErrorCode:          "ERR_001",
				ErrorMessage: defEntities.FieldValue{
					Type:  defEntities.FieldValueLiteral,
					Value: "something failed",
				},
			},
			resolverValue: "something failed",
			wantHandles:   []string{},
			wantError:     true,
			wantErrorCode: "ERR_001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockValueResolver{value: tt.resolverValue}
			executor := NewEndExecutor(resolver)

			ctx := newTestContext("core/end", tt.config)
			result, err := executor.Execute(context.Background(), ctx)
			if err != nil {
				t.Fatalf("unexpected execution error: %v", err)
			}

			if len(result.OutputHandles) != len(tt.wantHandles) {
				t.Fatalf("expected handles %v, got %v", tt.wantHandles, result.OutputHandles)
			}

			if tt.wantError {
				if result.Error == nil {
					t.Fatal("expected execution error, got nil")
				}
				if result.Error.Code != tt.wantErrorCode {
					t.Fatalf("expected error code %s, got %s", tt.wantErrorCode, result.Error.Code)
				}
			} else {
				if result.Error != nil {
					t.Fatalf("unexpected execution error: %v", result.Error)
				}
			}
		})
	}
}

// TestGotoExecutor tests GoTo executor for receiver passthrough, sender with edge, and sender without edge.
func TestGotoExecutor(t *testing.T) {
	executor := NewGotoExecutor()

	if executor.NodeType() != "core/goto" {
		t.Fatalf("expected NodeType core/goto, got %s", executor.NodeType())
	}

	// Receiver: passthrough — returns ["out"]
	receiverCtx := newTestContext("core/goto", &entities.GotoNodeConfig{
		Role:      "receiver",
		PairLabel: "error-handler",
	})
	receiverCtx.Graph = &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{},
	}
	result, err := executor.Execute(context.Background(), receiverCtx)
	if err != nil {
		t.Fatalf("receiver: unexpected error: %v", err)
	}
	if len(result.OutputHandles) != 1 || result.OutputHandles[0] != "out" {
		t.Fatalf("receiver: expected [out], got %v", result.OutputHandles)
	}

	// Sender with matching edge: passthrough — returns ["out"]
	senderCtx := newTestContext("core/goto", &entities.GotoNodeConfig{
		Role:      "sender",
		PairLabel: "error-handler",
	})
	senderCtx.Graph = &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{
			"node-1": {"out": "receiver-1"},
		},
	}
	result, err = executor.Execute(context.Background(), senderCtx)
	if err != nil {
		t.Fatalf("sender with edge: unexpected error: %v", err)
	}
	if len(result.OutputHandles) != 1 || result.OutputHandles[0] != "out" {
		t.Fatalf("sender with edge: expected [out], got %v", result.OutputHandles)
	}

	// Sender without matching edge: returns error result
	orphanCtx := newTestContext("core/goto", &entities.GotoNodeConfig{
		Role:      "sender",
		PairLabel: "missing-receiver",
	})
	orphanCtx.Graph = &entities.ExecutionGraph{
		Adjacency: map[string]map[string]string{},
	}
	result, err = executor.Execute(context.Background(), orphanCtx)
	if err != nil {
		t.Fatalf("orphan sender: unexpected error: %v", err)
	}
	if result.Error == nil {
		t.Fatal("orphan sender: expected error result, got nil")
	}
	if result.Error.Code != "GOTO_NO_RECEIVER" {
		t.Fatalf("orphan sender: expected GOTO_NO_RECEIVER, got %s", result.Error.Code)
	}
}

/*
 * CONDITION EXECUTOR
 */

func TestConditionExecutor(t *testing.T) {
	tests := []struct {
		name        string
		config      *entities.ConditionNodeConfig
		evalResult  bool
		wantHandle  string
	}{
		{
			name:       "nil config returns false",
			config:     nil,
			wantHandle: "false",
		},
		{
			name: "condition true",
			config: &entities.ConditionNodeConfig{
				Condition: defEntities.ConditionGroup{Logic: defEntities.LogicAND},
			},
			evalResult: true,
			wantHandle: "true",
		},
		{
			name: "condition false",
			config: &entities.ConditionNodeConfig{
				Condition: defEntities.ConditionGroup{Logic: defEntities.LogicAND},
			},
			evalResult: false,
			wantHandle: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := &mockConditionEvaluator{result: tt.evalResult}
			executor := NewConditionExecutor(evaluator)

			ctx := newTestContext("core/condition", tt.config)
			result, err := executor.Execute(context.Background(), ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.OutputHandles) != 1 || result.OutputHandles[0] != tt.wantHandle {
				t.Fatalf("expected [%s], got %v", tt.wantHandle, result.OutputHandles)
			}
		})
	}
}

/*
 * SWITCH EXECUTOR
 */

func TestSwitchExecutor(t *testing.T) {
	tests := []struct {
		name        string
		config      *entities.SwitchNodeConfig
		evalResults []bool
		wantHandles []string
	}{
		{
			name:        "nil config returns default",
			config:      nil,
			wantHandles: []string{"default"},
		},
		{
			name: "first match mode — returns first",
			config: &entities.SwitchNodeConfig{
				MatchMode: "first",
				Cases: []defEntities.SwitchCase{
					{ID: "a", Condition: defEntities.ConditionGroup{}},
					{ID: "b", Condition: defEntities.ConditionGroup{}},
				},
			},
			evalResults: []bool{true, true},
			wantHandles: []string{"case_a"},
		},
		{
			name: "all match mode — returns all",
			config: &entities.SwitchNodeConfig{
				MatchMode: "all",
				Cases: []defEntities.SwitchCase{
					{ID: "a", Condition: defEntities.ConditionGroup{}},
					{ID: "b", Condition: defEntities.ConditionGroup{}},
					{ID: "c", Condition: defEntities.ConditionGroup{}},
				},
			},
			evalResults: []bool{true, false, true},
			wantHandles: []string{"case_a", "case_c"},
		},
		{
			name: "no match returns default",
			config: &entities.SwitchNodeConfig{
				MatchMode: "first",
				Cases: []defEntities.SwitchCase{
					{ID: "a", Condition: defEntities.ConditionGroup{}},
				},
			},
			evalResults: []bool{false},
			wantHandles: []string{"default"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callIndex := 0
			evaluator := &sequentialMockEvaluator{
				results: tt.evalResults,
				index:   &callIndex,
			}
			executor := NewSwitchExecutor(evaluator)

			ctx := newTestContext("core/switch", tt.config)
			result, err := executor.Execute(context.Background(), ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.OutputHandles) != len(tt.wantHandles) {
				t.Fatalf("expected handles %v, got %v", tt.wantHandles, result.OutputHandles)
			}
			for i, h := range tt.wantHandles {
				if result.OutputHandles[i] != h {
					t.Fatalf("handle[%d]: expected %s, got %s", i, h, result.OutputHandles[i])
				}
			}
		})
	}
}

// sequentialMockEvaluator returns different results for each call
type sequentialMockEvaluator struct {
	results []bool
	index   *int
}

func (m *sequentialMockEvaluator) EvaluateGroup(
	_ *defEntities.ConditionGroup,
	_ string,
	_ map[string]interface{},
	_ map[string]interface{},
	_ map[string]interface{},
	_ map[string]interface{},
) (bool, error) {
	idx := *m.index
	*m.index++
	if idx < len(m.results) {
		return m.results[idx], nil
	}
	return false, nil
}

/*
 * SET STATE EXECUTOR
 */

func TestSetStateExecutor(t *testing.T) {
	tests := []struct {
		name          string
		config        *entities.SetStateNodeConfig
		state         map[string]interface{}
		resolverValue interface{}
		wantPatchKey  string
		wantPatchVal  interface{}
	}{
		{
			name:   "nil config returns out with no patch",
			config: nil,
		},
		{
			name: "set operation",
			config: &entities.SetStateNodeConfig{
				Operation:   "set",
				TargetField: "count",
				ValueSource: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "42"},
			},
			resolverValue: "42",
			wantPatchKey:  "count",
			wantPatchVal:  "42",
		},
		{
			name: "increment operation",
			config: &entities.SetStateNodeConfig{
				Operation:   "increment",
				TargetField: "counter",
				ValueSource: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "5"},
			},
			state:         map[string]interface{}{"counter": float64(10)},
			resolverValue: float64(5),
			wantPatchKey:  "counter",
			wantPatchVal:  float64(15),
		},
		{
			name: "decrement operation",
			config: &entities.SetStateNodeConfig{
				Operation:   "decrement",
				TargetField: "counter",
				ValueSource: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "3"},
			},
			state:         map[string]interface{}{"counter": float64(10)},
			resolverValue: float64(3),
			wantPatchKey:  "counter",
			wantPatchVal:  float64(7),
		},
		{
			name: "append operation",
			config: &entities.SetStateNodeConfig{
				Operation:   "append",
				TargetField: "tags",
				ValueSource: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "new"},
			},
			state:         map[string]interface{}{"tags": []interface{}{"old"}},
			resolverValue: "new",
			wantPatchKey:  "tags",
			wantPatchVal:  nil, // checked separately
		},
		{
			name: "remove operation deletes key",
			config: &entities.SetStateNodeConfig{
				Operation:   "remove",
				TargetField: "tags",
			},
			state:        map[string]interface{}{"tags": []interface{}{"old", "keep"}},
			wantPatchKey: "tags",
			wantPatchVal: nil, // nil sentinel — runtime will delete this key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockValueResolver{value: tt.resolverValue}
			executor := NewSetStateExecutor(resolver)

			ctx := newTestContext("core/set_state", tt.config)
			if tt.state != nil {
				ctx.State = tt.state
			}

			result, err := executor.Execute(context.Background(), ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.OutputHandles[0] != "out" {
				t.Fatalf("expected handle out, got %v", result.OutputHandles)
			}

			if tt.config == nil {
				return
			}

			val, ok := result.StatePatch[tt.wantPatchKey]
			if !ok {
				t.Fatalf("expected patch key %s", tt.wantPatchKey)
			}

			switch tt.config.Operation {
			case "append":
				arr, ok := val.([]interface{})
				if !ok {
					t.Fatalf("expected slice, got %T", val)
				}
				if len(arr) != 2 || arr[1] != "new" {
					t.Fatalf("expected [old new], got %v", arr)
				}
			case "remove":
				if val != nil {
					t.Fatalf("expected nil (delete sentinel), got %v", val)
				}
			default:
				if val != tt.wantPatchVal {
					t.Fatalf("expected patch value %v, got %v", tt.wantPatchVal, val)
				}
			}
		})
	}
}

/*
 * LOG EXECUTOR
 */

func TestLogExecutor(t *testing.T) {
	tests := []struct {
		name         string
		config       *entities.LogNodeConfig
		state        map[string]interface{}
		event        map[string]interface{}
		wantMessage  string
		wantLevel    entities.LogLevel
	}{
		{
			name:   "nil config returns out with no logs",
			config: nil,
		},
		{
			name:        "simple message",
			config:      &entities.LogNodeConfig{Message: "hello world", Level: "info"},
			wantMessage: "hello world",
			wantLevel:   entities.LogInfo,
		},
		{
			name:        "interpolates state",
			config:      &entities.LogNodeConfig{Message: "count is ${state.count}", Level: "debug"},
			state:       map[string]interface{}{"count": 42},
			wantMessage: "count is 42",
			wantLevel:   entities.LogDebug,
		},
		{
			name:        "interpolates event",
			config:      &entities.LogNodeConfig{Message: "user: ${event.name}", Level: "warn"},
			event:       map[string]interface{}{"name": "John"},
			wantMessage: "user: John",
			wantLevel:   entities.LogWarn,
		},
		{
			name:        "default level is info",
			config:      &entities.LogNodeConfig{Message: "test", Level: ""},
			wantMessage: "test",
			wantLevel:   entities.LogInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewLogExecutor()

			ctx := newTestContext("core/log", tt.config)
			if tt.state != nil {
				ctx.State = tt.state
			}
			if tt.event != nil {
				ctx.EventPayload = tt.event
			}

			result, err := executor.Execute(context.Background(), ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.OutputHandles[0] != "out" {
				t.Fatalf("expected handle out, got %v", result.OutputHandles)
			}

			if tt.config == nil {
				return
			}

			if len(result.LogEntries) != 1 {
				t.Fatalf("expected 1 log entry, got %d", len(result.LogEntries))
			}

			entry := result.LogEntries[0]
			if entry.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, entry.Message)
			}
			if entry.Level != tt.wantLevel {
				t.Fatalf("expected level %s, got %s", tt.wantLevel, entry.Level)
			}
		})
	}
}
