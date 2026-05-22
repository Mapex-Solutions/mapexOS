package services

import (
	"testing"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

	defEntities "workflow/src/modules/definitions/domain/entities"
	"workflow/src/modules/runtime/domain/entities"
)

func TestParseNodeConfig(t *testing.T) {
	tests := []struct {
		name     string
		nodeType string
		raw      map[string]interface{}
		verify   func(t *testing.T, result interface{})
	}{
		{
			name:     "nil raw returns nil",
			nodeType: "core/condition",
			raw:      nil,
			verify: func(t *testing.T, result interface{}) {
				if result != nil {
					t.Fatalf("expected nil, got %v", result)
				}
			},
		},
		{
			name:     "unknown type returns PluginNodeConfig",
			nodeType: "telegram/message",
			raw:      map[string]interface{}{"operation": "sendText", "credentialId": "cred-1"},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.PluginNodeConfig)
				if !ok {
					t.Fatalf("expected *PluginNodeConfig, got %T", result)
				}
				if cfg.Operation != "sendText" {
					t.Fatalf("expected operation=sendText, got %s", cfg.Operation)
				}
				if cfg.CredentialID != "cred-1" {
					t.Fatalf("expected credentialId=cred-1, got %s", cfg.CredentialID)
				}
			},
		},
		{
			name:     "core/condition",
			nodeType: "core/condition",
			// Real UI format: logic and items are top-level keys (not wrapped in "condition")
			raw: map[string]interface{}{
				"logic": "AND",
				"items": []interface{}{
					map[string]interface{}{
						"type": "condition",
						"data": map[string]interface{}{
							"id":       "c1",
							"name":     "condition",
							"field":    map[string]interface{}{"type": "literal", "value": "123"},
							"operator": "equals",
							"value":    map[string]interface{}{"type": "state", "value": "counter"},
						},
					},
				},
				"selectedTemplateIds": []interface{}{"tpl-1"},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.ConditionNodeConfig)
				if !ok {
					t.Fatalf("expected *ConditionNodeConfig, got %T", result)
				}
				if cfg.Condition.Logic != defEntities.LogicAND {
					t.Fatalf("expected logic AND, got %s", cfg.Condition.Logic)
				}
				if len(cfg.Condition.Items) != 1 {
					t.Fatalf("expected 1 item, got %d", len(cfg.Condition.Items))
				}
				if cfg.Condition.Items[0].Type != "condition" {
					t.Fatalf("expected item type=condition, got %s", cfg.Condition.Items[0].Type)
				}
				if len(cfg.SelectedTemplateIds) != 1 || cfg.SelectedTemplateIds[0] != "tpl-1" {
					t.Fatalf("expected [tpl-1], got %v", cfg.SelectedTemplateIds)
				}
			},
		},
		{
			name:     "core/switch",
			nodeType: "core/switch",
			raw: map[string]interface{}{
				"matchMode": "all",
				"cases": []interface{}{
					map[string]interface{}{
						"id":   "c1",
						"name": "Case 1",
						"condition": map[string]interface{}{
							"logic": "OR",
							"items": []interface{}{},
						},
					},
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.SwitchNodeConfig)
				if !ok {
					t.Fatalf("expected *SwitchNodeConfig, got %T", result)
				}
				if cfg.MatchMode != "all" {
					t.Fatalf("expected matchMode=all, got %s", cfg.MatchMode)
				}
				if len(cfg.Cases) != 1 {
					t.Fatalf("expected 1 case, got %d", len(cfg.Cases))
				}
				if cfg.Cases[0].ID != "c1" {
					t.Fatalf("expected case ID=c1, got %s", cfg.Cases[0].ID)
				}
			},
		},
		{
			name:     "core/fanout",
			nodeType: "core/fanout",
			raw:      map[string]interface{}{"branches": 3},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.FanoutNodeConfig)
				if !ok {
					t.Fatalf("expected *FanoutNodeConfig, got %T", result)
				}
				if cfg.Branches != 3 {
					t.Fatalf("expected branches=3, got %d", cfg.Branches)
				}
			},
		},
		{
			name:     "core/merge",
			nodeType: "core/merge",
			raw:      map[string]interface{}{"branches": 2, "strategy": "any"},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.MergeNodeConfig)
				if !ok {
					t.Fatalf("expected *MergeNodeConfig, got %T", result)
				}
				if cfg.Branches != 2 {
					t.Fatalf("expected branches=2, got %d", cfg.Branches)
				}
				if cfg.Strategy != "any" {
					t.Fatalf("expected strategy=any, got %s", cfg.Strategy)
				}
			},
		},
		{
			name:     "core/sequence",
			nodeType: "core/sequence",
			raw:      map[string]interface{}{"steps": 5},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.SequenceNodeConfig)
				if !ok {
					t.Fatalf("expected *SequenceNodeConfig, got %T", result)
				}
				if cfg.Steps != 5 {
					t.Fatalf("expected steps=5, got %d", cfg.Steps)
				}
			},
		},
		{
			name:     "core/loop",
			nodeType: "core/loop",
			raw: map[string]interface{}{
				"source": map[string]interface{}{
					"type":  "state",
					"value": "items",
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.LoopNodeConfig)
				if !ok {
					t.Fatalf("expected *LoopNodeConfig, got %T", result)
				}
				if cfg.Source.Type != defEntities.FieldValueState {
					t.Fatalf("expected source type=state, got %s", cfg.Source.Type)
				}
				if cfg.Source.Value != "items" {
					t.Fatalf("expected source value=items, got %s", cfg.Source.Value)
				}
			},
		},
		{
			name:     "core/goto",
			nodeType: "core/goto",
			raw: map[string]interface{}{
				"role":      "sender",
				"pairLabel": "error-handler",
				"pairColor": "#ff0000",
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.GotoNodeConfig)
				if !ok {
					t.Fatalf("expected *GotoNodeConfig, got %T", result)
				}
				if cfg.Role != "sender" {
					t.Fatalf("expected role=sender, got %s", cfg.Role)
				}
				if cfg.PairLabel != "error-handler" {
					t.Fatalf("expected pairLabel=error-handler, got %s", cfg.PairLabel)
				}
			},
		},
		{
			name:     "core/set_state",
			nodeType: "core/set_state",
			raw: map[string]interface{}{
				"operation":   "increment",
				"targetField": "counter",
				"valueSource": map[string]interface{}{
					"type":  "literal",
					"value": "1",
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.SetStateNodeConfig)
				if !ok {
					t.Fatalf("expected *SetStateNodeConfig, got %T", result)
				}
				if cfg.Operation != "increment" {
					t.Fatalf("expected operation=increment, got %s", cfg.Operation)
				}
				if cfg.TargetField != "counter" {
					t.Fatalf("expected targetField=counter, got %s", cfg.TargetField)
				}
				if cfg.ValueSource.Type != defEntities.FieldValueLiteral {
					t.Fatalf("expected valueSource type=literal, got %s", cfg.ValueSource.Type)
				}
			},
		},
		{
			name:     "core/log",
			nodeType: "core/log",
			raw:      map[string]interface{}{"message": "hello ${state.x}", "level": "warn"},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.LogNodeConfig)
				if !ok {
					t.Fatalf("expected *LogNodeConfig, got %T", result)
				}
				if cfg.Message != "hello ${state.x}" {
					t.Fatalf("expected message hello ${state.x}, got %s", cfg.Message)
				}
				if cfg.Level != "warn" {
					t.Fatalf("expected level=warn, got %s", cfg.Level)
				}
			},
		},
		{
			name:     "core/code",
			nodeType: "core/code",
			raw:      map[string]interface{}{"script": "return 42;", "timeout": 30},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.CodeNodeConfig)
				if !ok {
					t.Fatalf("expected *CodeNodeConfig, got %T", result)
				}
				if cfg.Script != "return 42;" {
					t.Fatalf("expected script=return 42;, got %s", cfg.Script)
				}
				if cfg.Timeout != 30 {
					t.Fatalf("expected timeout=30, got %d", cfg.Timeout)
				}
			},
		},
		{
			name:     "core/delay",
			nodeType: "core/delay",
			raw:      map[string]interface{}{"duration": 10, "unit": "s"},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.DelayNodeConfig)
				if !ok {
					t.Fatalf("expected *DelayNodeConfig, got %T", result)
				}
				if cfg.Duration != 10 {
					t.Fatalf("expected duration=10, got %d", cfg.Duration)
				}
				if cfg.Unit != "s" {
					t.Fatalf("expected unit=s, got %s", cfg.Unit)
				}
			},
		},
		{
			name:     "core/wait_signal",
			nodeType: "core/wait_signal",
			raw: map[string]interface{}{
				"signalName": "approval",
				"mappings": []interface{}{
					map[string]interface{}{
						"paramName": "approver",
						"value": map[string]interface{}{
							"type":  "event",
							"value": "userId",
						},
					},
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.WaitSignalNodeConfig)
				if !ok {
					t.Fatalf("expected *WaitSignalNodeConfig, got %T", result)
				}
				if cfg.SignalName != "approval" {
					t.Fatalf("expected signalName=approval, got %s", cfg.SignalName)
				}
				if len(cfg.Mappings) != 1 {
					t.Fatalf("expected 1 mapping, got %d", len(cfg.Mappings))
				}
				if cfg.Mappings[0].ParamName != "approver" {
					t.Fatalf("expected paramName=approver, got %s", cfg.Mappings[0].ParamName)
				}
			},
		},
		{
			name:     "core/wait_for",
			nodeType: "core/wait_for",
			raw: map[string]interface{}{
				"field":    "status",
				"operator": "==",
				"compareTo": map[string]interface{}{
					"type":  "literal",
					"value": "done",
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.WaitForNodeConfig)
				if !ok {
					t.Fatalf("expected *WaitForNodeConfig, got %T", result)
				}
				if cfg.Field != "status" {
					t.Fatalf("expected field=status, got %s", cfg.Field)
				}
				if cfg.Operator != "==" {
					t.Fatalf("expected operator===, got %s", cfg.Operator)
				}
				if cfg.CompareTo.Value != "done" {
					t.Fatalf("expected compareTo value=done, got %s", cfg.CompareTo.Value)
				}
			},
		},
		{
			name:     "core/subworkflow",
			nodeType: "core/subworkflow",
			raw: map[string]interface{}{
				"workflowId":    "wf-123",
				"workflowName":  "Child",
				"executionMode": "sync",
				"timeout": map[string]interface{}{
					"duration": 60,
					"unit":     "s",
				},
				"inputMappings": []interface{}{
					map[string]interface{}{
						"childParamName": "param1",
						"value": map[string]interface{}{
							"type":  "state",
							"value": "myField",
						},
					},
				},
				"outputMappings": []interface{}{
					map[string]interface{}{
						"outputName": "result",
						"stateField": "childResult",
					},
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.SubworkflowNodeConfig)
				if !ok {
					t.Fatalf("expected *SubworkflowNodeConfig, got %T", result)
				}
				if cfg.WorkflowID != "wf-123" {
					t.Fatalf("expected workflowId=wf-123, got %s", cfg.WorkflowID)
				}
				if cfg.ExecutionMode != "sync" {
					t.Fatalf("expected executionMode=sync, got %s", cfg.ExecutionMode)
				}
				if cfg.Timeout.Duration != 60 {
					t.Fatalf("expected timeout.duration=60, got %d", cfg.Timeout.Duration)
				}
				if len(cfg.InputMappings) != 1 {
					t.Fatalf("expected 1 input mapping, got %d", len(cfg.InputMappings))
				}
				if cfg.InputMappings[0].ChildParamName != "param1" {
					t.Fatalf("expected childParamName=param1, got %s", cfg.InputMappings[0].ChildParamName)
				}
				if len(cfg.OutputMappings) != 1 {
					t.Fatalf("expected 1 output mapping, got %d", len(cfg.OutputMappings))
				}
				if cfg.OutputMappings[0].OutputName != "result" {
					t.Fatalf("expected outputName=result, got %s", cfg.OutputMappings[0].OutputName)
				}
			},
		},
		{
			name:     "core/end",
			nodeType: "core/end",
			raw: map[string]interface{}{
				"terminateWithError": true,
				"errorCode":          "ERR_001",
				"errorMessage": map[string]interface{}{
					"type":  "literal",
					"value": "Something failed",
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.EndNodeConfig)
				if !ok {
					t.Fatalf("expected *EndNodeConfig, got %T", result)
				}
				if !cfg.TerminateWithError {
					t.Fatal("expected terminateWithError=true")
				}
				if cfg.ErrorCode != "ERR_001" {
					t.Fatalf("expected errorCode=ERR_001, got %s", cfg.ErrorCode)
				}
				if cfg.ErrorMessage.Value != "Something failed" {
					t.Fatalf("expected errorMessage value=Something failed, got %s", cfg.ErrorMessage.Value)
				}
			},
		},
		{
			name:     "core/trigger_event",
			nodeType: "core/trigger_event",
			raw: map[string]interface{}{
				"eventType": "user.created",
				"payloadMapping": []interface{}{
					map[string]interface{}{
						"key": "name",
						"value": map[string]interface{}{
							"type":  "state",
							"value": "userName",
						},
					},
				},
			},
			verify: func(t *testing.T, result interface{}) {
				cfg, ok := result.(*entities.TriggerEventNodeConfig)
				if !ok {
					t.Fatalf("expected *TriggerEventNodeConfig, got %T", result)
				}
				if cfg.EventType != "user.created" {
					t.Fatalf("expected eventType=user.created, got %s", cfg.EventType)
				}
				if len(cfg.PayloadMapping) != 1 {
					t.Fatalf("expected 1 payload mapping, got %d", len(cfg.PayloadMapping))
				}
				if cfg.PayloadMapping[0].Key != "name" {
					t.Fatalf("expected key=name, got %s", cfg.PayloadMapping[0].Key)
				}
				if cfg.PayloadMapping[0].Value.Type != defEntities.FieldValueState {
					t.Fatalf("expected value type=state, got %s", cfg.PayloadMapping[0].Value.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseNodeConfig(tt.nodeType, tt.raw)
			tt.verify(t, result)
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("model.MapGetString", func(t *testing.T) {
		m := map[string]interface{}{"key": "value", "num": 42}
		if model.MapGetString(m, "key") != "value" {
			t.Fatal("expected value")
		}
		if model.MapGetString(m, "num") != "" {
			t.Fatal("expected empty for non-string")
		}
		if model.MapGetString(m, "missing") != "" {
			t.Fatal("expected empty for missing")
		}
	})

	t.Run("model.MapGetInt", func(t *testing.T) {
		m := map[string]interface{}{"int": 42, "float": float64(3.14), "str": "hi"}
		if model.MapGetInt(m, "int") != 42 {
			t.Fatal("expected 42")
		}
		if model.MapGetInt(m, "float") != 3 {
			t.Fatal("expected 3")
		}
		if model.MapGetInt(m, "str") != 0 {
			t.Fatal("expected 0 for non-numeric")
		}
	})

	t.Run("model.MapGetBool", func(t *testing.T) {
		m := map[string]interface{}{"yes": true, "no": false, "str": "hi"}
		if !model.MapGetBool(m, "yes") {
			t.Fatal("expected true")
		}
		if model.MapGetBool(m, "no") {
			t.Fatal("expected false")
		}
		if model.MapGetBool(m, "str") {
			t.Fatal("expected false for non-bool")
		}
	})

	t.Run("model.MapGetMap", func(t *testing.T) {
		inner := map[string]interface{}{"a": 1}
		m := map[string]interface{}{"nested": inner, "str": "hi"}
		if model.MapGetMap(m, "nested")["a"] != 1 {
			t.Fatal("expected nested map")
		}
		if model.MapGetMap(m, "str") != nil {
			t.Fatal("expected nil for non-map")
		}
	})

	t.Run("model.MapGetSlice", func(t *testing.T) {
		m := map[string]interface{}{"arr": []interface{}{1, 2}, "str": "hi"}
		if len(model.MapGetSlice(m, "arr")) != 2 {
			t.Fatal("expected 2 elements")
		}
		if model.MapGetSlice(m, "str") != nil {
			t.Fatal("expected nil for non-slice")
		}
	})

	t.Run("model.MapGetStringSlice", func(t *testing.T) {
		m := map[string]interface{}{"tags": []interface{}{"a", "b"}}
		ss := model.MapGetStringSlice(m, "tags")
		if len(ss) != 2 || ss[0] != "a" || ss[1] != "b" {
			t.Fatalf("expected [a b], got %v", ss)
		}
	})
}
