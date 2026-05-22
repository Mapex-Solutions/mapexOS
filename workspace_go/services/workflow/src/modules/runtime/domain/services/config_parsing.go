package services

import (
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

	defPorts "workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/runtime/domain/entities"
)

// parseNodeConfig dispatches config parsing by node type.
// Returns nil for node types that don't need config (e.g., core/start).
func parseNodeConfig(nodeType string, raw map[string]interface{}) interface{} {
	if raw == nil {
		return nil
	}
	switch nodeType {
	case "core/condition":
		return parseConditionNodeConfig(raw)
	case "core/switch":
		return parseSwitchNodeConfig(raw)
	case "core/fanout":
		return parseFanoutNodeConfig(raw)
	case "core/merge":
		return parseMergeNodeConfig(raw)
	case "core/sequence":
		return parseSequenceNodeConfig(raw)
	case "core/loop":
		return parseLoopNodeConfig(raw)
	case "core/goto":
		return parseGotoNodeConfig(raw)
	case "core/set_state":
		return parseSetStateNodeConfig(raw)
	case "core/log":
		return parseLogNodeConfig(raw)
	case "core/code":
		return parseCodeNodeConfig(raw)
	case "core/delay":
		return parseDelayNodeConfig(raw)
	case "core/wait_signal":
		return parseWaitSignalNodeConfig(raw)
	case "core/wait_for":
		return parseWaitForNodeConfig(raw)
	case "core/subworkflow":
		return parseSubworkflowNodeConfig(raw)
	case "core/end":
		return parseEndNodeConfig(raw)
	case "core/trigger_event":
		return parseTriggerEventNodeConfig(raw)
	default:
		// Non-core node types are plugin nodes — extract minimal fields, preserve raw config
		return parsePluginNodeConfig(raw)
	}
}

// parseConditionNodeConfig parses a condition node config from a raw map.
func parseConditionNodeConfig(m map[string]interface{}) *entities.ConditionNodeConfig {
	return &entities.ConditionNodeConfig{
		Condition:           parseConditionGroup(m),
		SelectedTemplateIds: model.MapGetStringSlice(m, "selectedTemplateIds"),
	}
}

// parseSwitchNodeConfig parses a switch node config from a raw map.
func parseSwitchNodeConfig(m map[string]interface{}) *entities.SwitchNodeConfig {
	cfg := &entities.SwitchNodeConfig{
		MatchMode:           model.MapGetString(m, "matchMode"),
		SelectedTemplateIds: model.MapGetStringSlice(m, "selectedTemplateIds"),
	}
	rawCases := model.MapGetSlice(m, "cases")
	if rawCases != nil {
		cfg.Cases = make([]defPorts.SwitchCase, 0, len(rawCases))
		for _, raw := range rawCases {
			if caseMap, ok := raw.(map[string]interface{}); ok {
				cfg.Cases = append(cfg.Cases, parseSwitchCase(caseMap))
			}
		}
	}
	return cfg
}

// parseFanoutNodeConfig parses a fanout node config from a raw map.
func parseFanoutNodeConfig(m map[string]interface{}) *entities.FanoutNodeConfig {
	mode := model.MapGetString(m, "mode")
	if mode == "" {
		mode = "waitAll"
	}
	return &entities.FanoutNodeConfig{
		Branches: model.MapGetInt(m, "branches"),
		Mode:     mode,
	}
}

// parseMergeNodeConfig parses a merge node config from a raw map.
func parseMergeNodeConfig(m map[string]interface{}) *entities.MergeNodeConfig {
	return &entities.MergeNodeConfig{
		Branches: model.MapGetInt(m, "branches"),
		Strategy: model.MapGetString(m, "strategy"),
	}
}

// parseSequenceNodeConfig parses a sequence node config from a raw map.
func parseSequenceNodeConfig(m map[string]interface{}) *entities.SequenceNodeConfig {
	return &entities.SequenceNodeConfig{
		Steps: model.MapGetInt(m, "steps"),
	}
}

// parseLoopNodeConfig parses a loop node config from a raw map.
func parseLoopNodeConfig(m map[string]interface{}) *entities.LoopNodeConfig {
	return &entities.LoopNodeConfig{
		Source: parseFieldValue(model.MapGetMap(m, "source")),
	}
}

// parseGotoNodeConfig parses a goto node config from a raw map.
func parseGotoNodeConfig(m map[string]interface{}) *entities.GotoNodeConfig {
	return &entities.GotoNodeConfig{
		Role:      model.MapGetString(m, "role"),
		PairLabel: model.MapGetString(m, "pairLabel"),
		PairColor: model.MapGetString(m, "pairColor"),
	}
}

// parseSetStateNodeConfig parses a set_state node config from a raw map.
func parseSetStateNodeConfig(m map[string]interface{}) *entities.SetStateNodeConfig {
	return &entities.SetStateNodeConfig{
		Operation:           model.MapGetString(m, "operation"),
		TargetField:         model.MapGetString(m, "targetField"),
		ValueSource:         parseFieldValue(model.MapGetMap(m, "valueSource")),
		SelectedTemplateIds: model.MapGetStringSlice(m, "selectedTemplateIds"),
	}
}

// parseLogNodeConfig parses a log node config from a raw map.
func parseLogNodeConfig(m map[string]interface{}) *entities.LogNodeConfig {
	return &entities.LogNodeConfig{
		Message: model.MapGetString(m, "message"),
		Level:   model.MapGetString(m, "level"),
	}
}

// parseCodeNodeConfig parses a code node config from a raw map.
func parseCodeNodeConfig(m map[string]interface{}) *entities.CodeNodeConfig {
	return &entities.CodeNodeConfig{
		Script:  model.MapGetString(m, "script"),
		Timeout: model.MapGetInt(m, "timeout"),
	}
}

// parseDelayNodeConfig parses a delay node config from a raw map.
func parseDelayNodeConfig(m map[string]interface{}) *entities.DelayNodeConfig {
	return &entities.DelayNodeConfig{
		Duration: model.MapGetInt(m, "duration"),
		Unit:     model.MapGetString(m, "unit"),
	}
}

// parseWaitSignalNodeConfig parses a wait_signal node config from a raw map.
func parseWaitSignalNodeConfig(m map[string]interface{}) *entities.WaitSignalNodeConfig {
	cfg := &entities.WaitSignalNodeConfig{
		SignalName: model.MapGetString(m, "signalName"),
	}
	rawMappings := model.MapGetSlice(m, "mappings")
	if rawMappings != nil {
		cfg.Mappings = make([]entities.SignalMapping, 0, len(rawMappings))
		for _, raw := range rawMappings {
			if mapItem, ok := raw.(map[string]interface{}); ok {
				cfg.Mappings = append(cfg.Mappings, entities.SignalMapping{
					ParamName: model.MapGetString(mapItem, "paramName"),
					Value:     parseFieldValue(model.MapGetMap(mapItem, "value")),
				})
			}
		}
	}
	return cfg
}

// parseWaitForNodeConfig parses a wait_for node config from a raw map.
func parseWaitForNodeConfig(m map[string]interface{}) *entities.WaitForNodeConfig {
	return &entities.WaitForNodeConfig{
		Field:     model.MapGetString(m, "field"),
		Operator:  model.MapGetString(m, "operator"),
		CompareTo: parseFieldValue(model.MapGetMap(m, "compareTo")),
	}
}

// parseSubworkflowNodeConfig parses a subworkflow node config from a raw map.
func parseSubworkflowNodeConfig(m map[string]interface{}) *entities.SubworkflowNodeConfig {
	cfg := &entities.SubworkflowNodeConfig{
		WorkflowID:    model.MapGetString(m, "workflowId"),
		WorkflowName:  model.MapGetString(m, "workflowName"),
		ExecutionMode: model.MapGetString(m, "executionMode"),
	}
	if timeoutMap := model.MapGetMap(m, "timeout"); timeoutMap != nil {
		cfg.Timeout = entities.TimeoutConfig{
			Duration: model.MapGetInt(timeoutMap, "duration"),
			Unit:     model.MapGetString(timeoutMap, "unit"),
		}
	}
	rawInputs := model.MapGetSlice(m, "inputMappings")
	if rawInputs != nil {
		cfg.InputMappings = make([]entities.InputMapping, 0, len(rawInputs))
		for _, raw := range rawInputs {
			if mapItem, ok := raw.(map[string]interface{}); ok {
				cfg.InputMappings = append(cfg.InputMappings, entities.InputMapping{
					ChildParamName: model.MapGetString(mapItem, "childParamName"),
					Value:          parseFieldValue(model.MapGetMap(mapItem, "value")),
				})
			}
		}
	}
	rawOutputs := model.MapGetSlice(m, "outputMappings")
	if rawOutputs != nil {
		cfg.OutputMappings = make([]entities.OutputMapping, 0, len(rawOutputs))
		for _, raw := range rawOutputs {
			if mapItem, ok := raw.(map[string]interface{}); ok {
				cfg.OutputMappings = append(cfg.OutputMappings, entities.OutputMapping{
					OutputName: model.MapGetString(mapItem, "outputName"),
					StateField: model.MapGetString(mapItem, "stateField"),
				})
			}
		}
	}
	return cfg
}

// parseEndNodeConfig parses an end node config from a raw map.
func parseEndNodeConfig(m map[string]interface{}) *entities.EndNodeConfig {
	return &entities.EndNodeConfig{
		TerminateWithError: model.MapGetBool(m, "terminateWithError"),
		ErrorCode:          model.MapGetString(m, "errorCode"),
		ErrorMessage:       parseFieldValue(model.MapGetMap(m, "errorMessage")),
	}
}

// parseTriggerEventNodeConfig parses a trigger_event node config from a raw map.
func parseTriggerEventNodeConfig(m map[string]interface{}) *entities.TriggerEventNodeConfig {
	cfg := &entities.TriggerEventNodeConfig{
		EventType: model.MapGetString(m, "eventType"),
	}
	rawMapping := model.MapGetSlice(m, "payloadMapping")
	if rawMapping != nil {
		cfg.PayloadMapping = make([]entities.TriggerPayloadField, 0, len(rawMapping))
		for _, raw := range rawMapping {
			if fieldMap, ok := raw.(map[string]interface{}); ok {
				cfg.PayloadMapping = append(cfg.PayloadMapping, entities.TriggerPayloadField{
					Key:   model.MapGetString(fieldMap, "key"),
					Value: parseFieldValue(model.MapGetMap(fieldMap, "value")),
				})
			}
		}
	}
	return cfg
}

// parseFieldValue converts a raw map to a FieldValue struct.
// Accepts both "type" and "source" keys for the field type — the UI uses "source"
// in wait_for.compareTo while all other nodes use "type".
func parseFieldValue(m map[string]interface{}) defPorts.FieldValue {
	if m == nil {
		return defPorts.FieldValue{}
	}
	fieldType := model.MapGetString(m, "type")
	if fieldType == "" {
		fieldType = model.MapGetString(m, "source")
	}
	return defPorts.FieldValue{
		Type:   defPorts.FieldValueType(fieldType),
		Value:  model.MapGetString(m, "value"),
		Mode:   model.MapGetString(m, "mode"),
		NodeID: model.MapGetString(m, "nodeId"),
	}
}

// parseConditionItem converts a raw map to a ConditionItem struct.
func parseConditionItem(m map[string]interface{}) defPorts.ConditionItem {
	return defPorts.ConditionItem{
		ID:       model.MapGetString(m, "id"),
		Name:     model.MapGetString(m, "name"),
		Field:    parseFieldValue(model.MapGetMap(m, "field")),
		Operator: model.MapGetString(m, "operator"),
		Value:    parseFieldValue(model.MapGetMap(m, "value")),
	}
}

// parseConditionGroupItem converts a raw map to a ConditionGroupItem struct.
// Recursively parses nested condition groups.
func parseConditionGroupItem(m map[string]interface{}) defPorts.ConditionGroupItem {
	item := defPorts.ConditionGroupItem{
		Type: model.MapGetString(m, "type"),
	}
	data := model.MapGetMap(m, "data")
	if data == nil {
		return item
	}
	switch item.Type {
	case "condition":
		parsed := parseConditionItem(data)
		item.Data = parsed
	case "group":
		parsed := parseConditionGroup(data)
		item.Data = parsed
	}
	return item
}

// parseConditionGroup converts a raw map to a ConditionGroup struct.
// Recursively parses nested items (conditions and sub-groups).
func parseConditionGroup(m map[string]interface{}) defPorts.ConditionGroup {
	if m == nil {
		return defPorts.ConditionGroup{}
	}
	group := defPorts.ConditionGroup{
		ID:    model.MapGetString(m, "id"),
		Name:  model.MapGetString(m, "name"),
		Logic: defPorts.GroupLogicOperator(model.MapGetString(m, "logic")),
	}
	rawItems := model.MapGetSlice(m, "items")
	if rawItems != nil {
		group.Items = make([]defPorts.ConditionGroupItem, 0, len(rawItems))
		for _, raw := range rawItems {
			if itemMap, ok := raw.(map[string]interface{}); ok {
				group.Items = append(group.Items, parseConditionGroupItem(itemMap))
			}
		}
	}
	return group
}

// parseSwitchCase converts a raw map to a SwitchCase struct.
func parseSwitchCase(m map[string]interface{}) defPorts.SwitchCase {
	return defPorts.SwitchCase{
		ID:        model.MapGetString(m, "id"),
		Name:      model.MapGetString(m, "name"),
		Condition: parseConditionGroup(model.MapGetMap(m, "condition")),
	}
}

// parsePluginNodeConfig extracts minimal fields from a plugin node config.
// Preserves the raw config map for template resolution at dispatch time.
func parsePluginNodeConfig(m map[string]interface{}) *entities.PluginNodeConfig {
	return &entities.PluginNodeConfig{
		Operation:    model.MapGetString(m, "operation"),
		CredentialID: model.MapGetString(m, "credentialId"),
		RawConfig:    m,
	}
}
