package control

import (
	"context"
	"errors"
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
 * FANOUT EXECUTOR
 */

func TestFanoutExecutor(t *testing.T) {
	tests := []struct {
		name        string
		config      *entities.FanoutNodeConfig
		wantErr     bool
		wantErrType error
		wantCount   int
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "zero branches returns error",
			config:  &entities.FanoutNodeConfig{Branches: 0},
			wantErr: true,
		},
		{
			name:      "3 branches",
			config:    &entities.FanoutNodeConfig{Branches: 3},
			wantCount: 3,
		},
		{
			name:        "exceeds max branches",
			config:      &entities.FanoutNodeConfig{Branches: 50},
			wantErr:     true,
			wantErrType: entities.ErrMaxFanoutBranches,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewFanoutExecutor()

			result, err := executor.Execute(context.Background(), newTestContext("core/fanout", tt.config))

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

			if len(result.OutputHandles) != tt.wantCount {
				t.Fatalf("expected %d handles, got %d", tt.wantCount, len(result.OutputHandles))
			}

			for i, h := range result.OutputHandles {
				expected := "out_" + itoa(i+1)
				if h != expected {
					t.Fatalf("handle[%d]: expected %s, got %s", i, expected, h)
				}
			}
		})
	}
}

func itoa(i int) string {
	return string(rune('0'+i)) // works for single digit
}

/*
 * MERGE EXECUTOR
 */

func TestMergeExecutor(t *testing.T) {
	tests := []struct {
		name        string
		config      *entities.MergeNodeConfig
		nodeStates  map[string]map[string]interface{}
		wantErr     bool
		wantProceed bool
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name:        "all strategy — not enough branches",
			config:      &entities.MergeNodeConfig{Branches: 3, Strategy: "all"},
			nodeStates:  map[string]map[string]interface{}{},
			wantProceed: false,
		},
		{
			name:   "all strategy — enough branches",
			config: &entities.MergeNodeConfig{Branches: 2, Strategy: "all"},
			nodeStates: map[string]map[string]interface{}{
				"node-1": {"branchCount": 1},
			},
			wantProceed: true,
		},
		{
			name:        "any strategy — first branch",
			config:      &entities.MergeNodeConfig{Branches: 3, Strategy: "any"},
			nodeStates:  map[string]map[string]interface{}{},
			wantProceed: true,
		},
		{
			name:        "first strategy — first branch",
			config:      &entities.MergeNodeConfig{Branches: 3, Strategy: "first"},
			nodeStates:  map[string]map[string]interface{}{},
			wantProceed: true,
		},
		{
			name:        "default strategy is all",
			config:      &entities.MergeNodeConfig{Branches: 2, Strategy: ""},
			nodeStates:  map[string]map[string]interface{}{},
			wantProceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewMergeExecutor()

			ctx := newTestContext("core/merge", tt.config)
			if tt.nodeStates != nil {
				ctx.NodeStates = tt.nodeStates
			}

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

			if tt.wantProceed {
				if len(result.OutputHandles) != 1 || result.OutputHandles[0] != "out" {
					t.Fatalf("expected [out], got %v", result.OutputHandles)
				}
			} else {
				if len(result.OutputHandles) != 0 {
					t.Fatalf("expected empty handles, got %v", result.OutputHandles)
				}
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
		})
	}
}

/*
 * SEQUENCE EXECUTOR
 */

func TestSequenceExecutor(t *testing.T) {
	tests := []struct {
		name       string
		config     *entities.SequenceNodeConfig
		nodeStates map[string]map[string]interface{}
		wantErr    bool
		wantHandle string
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name:       "first step",
			config:     &entities.SequenceNodeConfig{Steps: 3},
			nodeStates: map[string]map[string]interface{}{},
			wantHandle: "step_2",
		},
		{
			name:   "second step",
			config: &entities.SequenceNodeConfig{Steps: 3},
			nodeStates: map[string]map[string]interface{}{
				"node-1": {"currentStep": 1},
			},
			wantHandle: "step_3",
		},
		{
			name:   "last step completes",
			config: &entities.SequenceNodeConfig{Steps: 2},
			nodeStates: map[string]map[string]interface{}{
				"node-1": {"currentStep": 1},
			},
			wantHandle: "done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewSequenceExecutor()

			ctx := newTestContext("core/sequence", tt.config)
			if tt.nodeStates != nil {
				ctx.NodeStates = tt.nodeStates
			}

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

			if len(result.OutputHandles) != 1 || result.OutputHandles[0] != tt.wantHandle {
				t.Fatalf("expected [%s], got %v", tt.wantHandle, result.OutputHandles)
			}

			if result.NodeState == nil {
				t.Fatal("expected NodeState, got nil")
			}
		})
	}
}

/*
 * LOOP EXECUTOR
 */

func TestLoopExecutor(t *testing.T) {
	tests := []struct {
		name          string
		config        *entities.LoopNodeConfig
		nodeStates    map[string]map[string]interface{}
		resolverValue interface{}
		resolverErr   error
		wantErr       bool
		wantErrType   error
		wantHandle    string
		wantItem      interface{}
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name: "first iteration",
			config: &entities.LoopNodeConfig{
				Source: defEntities.FieldValue{Type: defEntities.FieldValueState, Value: "items"},
			},
			nodeStates:    map[string]map[string]interface{}{},
			resolverValue: []interface{}{"a", "b", "c"},
			wantHandle:    "body",
			wantItem:      "a",
		},
		{
			name: "middle iteration",
			config: &entities.LoopNodeConfig{
				Source: defEntities.FieldValue{Type: defEntities.FieldValueState, Value: "items"},
			},
			nodeStates: map[string]map[string]interface{}{
				"node-1": {"currentIndex": 1},
			},
			resolverValue: []interface{}{"a", "b", "c"},
			wantHandle:    "body",
			wantItem:      "b",
		},
		{
			name: "done — all items processed",
			config: &entities.LoopNodeConfig{
				Source: defEntities.FieldValue{Type: defEntities.FieldValueState, Value: "items"},
			},
			nodeStates: map[string]map[string]interface{}{
				"node-1": {"currentIndex": 3},
			},
			resolverValue: []interface{}{"a", "b", "c"},
			wantHandle:    "done",
		},
		{
			name: "non-array source returns error",
			config: &entities.LoopNodeConfig{
				Source: defEntities.FieldValue{Type: defEntities.FieldValueLiteral, Value: "not-array"},
			},
			resolverValue: "not-array",
			wantErr:       true,
		},
		{
			name: "resolver error propagates",
			config: &entities.LoopNodeConfig{
				Source: defEntities.FieldValue{Type: defEntities.FieldValueState, Value: "items"},
			},
			resolverErr: errors.New("resolve error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockValueResolver{value: tt.resolverValue, err: tt.resolverErr}
			executor := NewLoopExecutor(resolver)

			ctx := newTestContext("core/loop", tt.config)
			if tt.nodeStates != nil {
				ctx.NodeStates = tt.nodeStates
			}

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

			if len(result.OutputHandles) != 1 || result.OutputHandles[0] != tt.wantHandle {
				t.Fatalf("expected [%s], got %v", tt.wantHandle, result.OutputHandles)
			}

			if tt.wantItem != nil {
				if result.StatePatch["loop_item"] != tt.wantItem {
					t.Fatalf("expected loop_item=%v, got %v", tt.wantItem, result.StatePatch["loop_item"])
				}
			}
		})
	}
}

/*
 * WAIT FOR EXECUTOR
 */

func TestWaitForExecutor(t *testing.T) {
	tests := []struct {
		name       string
		config     *entities.WaitForNodeConfig
		evalResult bool
		evalErr    error
		wantErr    bool
		wantWait   bool
	}{
		{
			name:    "nil config returns error",
			config:  nil,
			wantErr: true,
		},
		{
			name: "condition met — no wait",
			config: &entities.WaitForNodeConfig{
				Field:    "status",
				Operator: "==",
				CompareTo: defEntities.FieldValue{
					Type:  defEntities.FieldValueLiteral,
					Value: "done",
				},
			},
			evalResult: true,
			wantWait:   false,
		},
		{
			name: "condition not met — wait",
			config: &entities.WaitForNodeConfig{
				Field:    "status",
				Operator: "==",
				CompareTo: defEntities.FieldValue{
					Type:  defEntities.FieldValueLiteral,
					Value: "done",
				},
			},
			evalResult: false,
			wantWait:   true,
		},
		{
			name: "evaluator error propagates",
			config: &entities.WaitForNodeConfig{
				Field:    "status",
				Operator: "==",
				CompareTo: defEntities.FieldValue{
					Type:  defEntities.FieldValueLiteral,
					Value: "done",
				},
			},
			evalErr: errors.New("eval error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := &mockConditionEvaluator{result: tt.evalResult, err: tt.evalErr}
			executor := NewWaitForExecutor(evaluator)

			ctx := newTestContext("core/wait_for", tt.config)
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

			if tt.wantWait {
				if len(result.OutputHandles) != 0 {
					t.Fatalf("expected empty OutputHandles for wait, got %v", result.OutputHandles)
				}
			} else {
				if len(result.OutputHandles) != 1 || result.OutputHandles[0] != "matched" {
					t.Fatalf("expected [matched], got %v", result.OutputHandles)
				}
			}

			if tt.wantWait {
				if result.NodeState == nil {
					t.Fatal("expected NodeState, got nil")
				}
				if result.NodeState["waitType"] != "condition" {
					t.Fatalf("expected waitType condition, got %v", result.NodeState["waitType"])
				}
				if result.NodeState["expiresAt"] == nil {
					t.Fatal("expected expiresAt in NodeState, got nil")
				}
				if result.NodeState["field"] != tt.config.Field {
					t.Fatalf("expected field %s, got %v", tt.config.Field, result.NodeState["field"])
				}
			} else {
				if result.NodeState != nil {
					t.Fatalf("expected no NodeState, got %v", result.NodeState)
				}
			}
		})
	}
}
