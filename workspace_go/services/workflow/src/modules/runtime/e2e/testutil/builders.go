package dagwalker

import (
	"fmt"

	defEntities "workflow/src/modules/definitions/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * DefinitionBuilder
 * Provides a fluent API to construct WorkflowDefinition for tests.
 */
type DefinitionBuilder struct {
	def defEntities.WorkflowDefinition
}

// NewDefinition creates a new builder with required defaults.
func NewDefinition(name string) *DefinitionBuilder {
	return &DefinitionBuilder{
		def: defEntities.WorkflowDefinition{
			ID:                model.NewObjectID(),
			OrgID:             ptrObjID(model.NewObjectID()),
			Name:              name,
			Enabled:           true,
			Status:            "valid",
			DefinitionVersion: 1,
			Nodes:             []defEntities.WorkflowNode{},
			Edges:             []defEntities.WorkflowEdge{},
			PathKey:           "000001",
		},
	}
}

func (b *DefinitionBuilder) WithState(field string, varType string, defaultValue interface{}) *DefinitionBuilder {
	b.def.States = append(b.def.States, defEntities.WorkflowVariable{
		Field:        field,
		Type:         defEntities.VariableType(varType),
		DefaultValue: defaultValue,
	})
	return b
}

func (b *DefinitionBuilder) WithExternalInput(field string, defaultValue interface{}) *DefinitionBuilder {
	b.def.ExternalInputs = append(b.def.ExternalInputs, defEntities.ExternalInput{
		Field:        field,
		Type:         "string",
		DefaultValue: defaultValue,
	})
	return b
}

func (b *DefinitionBuilder) AddNode(id, nodeType string, config map[string]interface{}) *DefinitionBuilder {
	b.def.Nodes = append(b.def.Nodes, defEntities.WorkflowNode{
		ID:     id,
		Type:   nodeType,
		Label:  id,
		Config: config,
	})
	return b
}

func (b *DefinitionBuilder) AddNodeWithLabel(id, nodeType, label string, config map[string]interface{}) *DefinitionBuilder {
	b.def.Nodes = append(b.def.Nodes, defEntities.WorkflowNode{
		ID:     id,
		Type:   nodeType,
		Label:  label,
		Config: config,
	})
	return b
}

func (b *DefinitionBuilder) AddEdge(source, sourceHandle, target string) *DefinitionBuilder {
	b.def.Edges = append(b.def.Edges, defEntities.WorkflowEdge{
		ID:           fmt.Sprintf("e_%s_%s_%s", source, sourceHandle, target),
		Source:       source,
		SourceHandle: sourceHandle,
		Target:       target,
		TargetHandle: "in",
	})
	return b
}

func (b *DefinitionBuilder) Build() *defEntities.WorkflowDefinition {
	return &b.def
}

/**
 * Node Convenience Functions
 * Each function returns (id, type, config) for use with AddNode.
 */

// StartNode returns (id, type, config) for a core/start node.
func StartNode(id string) (string, string, map[string]interface{}) {
	return id, "core/start", map[string]interface{}{}
}

// EndNode returns (id, type, config) for a core/end node (success).
func EndNode(id string) (string, string, map[string]interface{}) {
	return id, "core/end", map[string]interface{}{
		"terminateWithError": false,
	}
}

// EndNodeWithError returns (id, type, config) for a core/end node that terminates with error (literal message).
func EndNodeWithError(id, errorCode, errorMessage string) (string, string, map[string]interface{}) {
	return id, "core/end", map[string]interface{}{
		"terminateWithError": true,
		"errorCode":          errorCode,
		"errorMessage": map[string]interface{}{
			"type":  "literal",
			"value": errorMessage,
		},
	}
}

// EndNodeWithErrorSource returns (id, type, config) for a core/end node that terminates
// with error using a dynamic FieldValue source for the message.
func EndNodeWithErrorSource(id, errorCode string, errorMessage map[string]interface{}) (string, string, map[string]interface{}) {
	return id, "core/end", map[string]interface{}{
		"terminateWithError": true,
		"errorCode":          errorCode,
		"errorMessage":       errorMessage,
	}
}

// SetStateNode returns (id, type, config) for a core/set_state node.
func SetStateNode(id, operation, targetField string, valueSource map[string]interface{}) (string, string, map[string]interface{}) {
	return id, "core/set_state", map[string]interface{}{
		"operation":   operation,
		"targetField": targetField,
		"valueSource": valueSource,
	}
}

// LogNode returns (id, type, config) for a core/log node.
func LogNode(id, message, level string) (string, string, map[string]interface{}) {
	return id, "core/log", map[string]interface{}{
		"message": message,
		"level":   level,
	}
}

// CodeNode returns (id, type, config) for a core/code node.
func CodeNode(id, script string, timeout int) (string, string, map[string]interface{}) {
	return id, "core/code", map[string]interface{}{
		"script":  script,
		"timeout": timeout,
	}
}

// DelayNode returns (id, type, config) for a core/delay node.
func DelayNode(id string, duration int, unit string) (string, string, map[string]interface{}) {
	return id, "core/delay", map[string]interface{}{
		"duration": duration,
		"unit":     unit,
	}
}

// WaitSignalNode returns (id, type, config) for a core/wait_signal node.
func WaitSignalNode(id, signalName string) (string, string, map[string]interface{}) {
	return id, "core/wait_signal", map[string]interface{}{
		"signalName": signalName,
	}
}

// FanoutNode returns (id, type, config) for a core/fanout node.
func FanoutNode(id string, branches int, mode string) (string, string, map[string]interface{}) {
	cfg := map[string]interface{}{
		"branches": branches,
	}
	if mode != "" {
		cfg["mode"] = mode
	}
	return id, "core/fanout", cfg
}

// MergeNode returns (id, type, config) for a core/merge node.
func MergeNode(id string, branches int) (string, string, map[string]interface{}) {
	return id, "core/merge", map[string]interface{}{
		"branches": branches,
	}
}

// LoopNode returns (id, type, config) for a core/loop node.
func LoopNode(id string, source map[string]interface{}) (string, string, map[string]interface{}) {
	return id, "core/loop", map[string]interface{}{
		"source": source,
	}
}

// GotoSender returns (id, type, config) for a core/goto sender node.
func GotoSender(id, pairLabel string) (string, string, map[string]interface{}) {
	return id, "core/goto", map[string]interface{}{
		"role":      "sender",
		"pairLabel": pairLabel,
	}
}

// GotoReceiver returns (id, type, config) for a core/goto receiver node.
func GotoReceiver(id, pairLabel string) (string, string, map[string]interface{}) {
	return id, "core/goto", map[string]interface{}{
		"role":      "receiver",
		"pairLabel": pairLabel,
	}
}

// SubworkflowNode returns (id, type, config) for a core/subworkflow node.
func SubworkflowNode(id, workflowId string) (string, string, map[string]interface{}) {
	return id, "core/subworkflow", map[string]interface{}{
		"workflowId": workflowId,
	}
}

// SequenceNode returns (id, type, config) for a core/sequence node.
func SequenceNode(id string, steps int) (string, string, map[string]interface{}) {
	return id, "core/sequence", map[string]interface{}{
		"steps": steps,
	}
}

// ConditionItem represents a single condition for test builders.
type ConditionItem struct {
	Field    map[string]interface{}
	Operator string
	Value    map[string]interface{}
}

// ConditionNode returns (id, type, config) for a core/condition node.
// The parser calls parseConditionGroup(configMap) directly, so logic/items are top-level keys.
func ConditionNode(id, logic string, items []ConditionItem) (string, string, map[string]interface{}) {
	groupItems := make([]interface{}, len(items))
	for i, item := range items {
		groupItems[i] = map[string]interface{}{
			"type": "condition",
			"data": map[string]interface{}{
				"field":    item.Field,
				"operator": item.Operator,
				"value":    item.Value,
			},
		}
	}
	return id, "core/condition", map[string]interface{}{
		"logic": logic,
		"items": groupItems,
	}
}

// SwitchCaseItem represents a single switch case for test builders.
type SwitchCaseItem struct {
	ID         string
	Logic      string
	Conditions []ConditionItem
}

// SwitchNode returns (id, type, config) for a core/switch node.
func SwitchNode(id, matchMode string, cases []SwitchCaseItem) (string, string, map[string]interface{}) {
	caseMaps := make([]interface{}, 0, len(cases))
	for _, c := range cases {
		items := make([]interface{}, len(c.Conditions))
		for i, cond := range c.Conditions {
			items[i] = map[string]interface{}{
				"type": "condition",
				"data": map[string]interface{}{
					"field":    cond.Field,
					"operator": cond.Operator,
					"value":    cond.Value,
				},
			}
		}
		caseMaps = append(caseMaps, map[string]interface{}{
			"id": c.ID,
			"condition": map[string]interface{}{
				"logic": c.Logic,
				"items": items,
			},
		})
	}
	return id, "core/switch", map[string]interface{}{
		"matchMode": matchMode,
		"cases":     caseMaps,
	}
}

// WaitForNode returns (id, type, config) for a core/wait_for node.
// Uses the real UI format: compareTo has "source" key (not "type").
func WaitForNode(id, field, operator string, compareTo map[string]interface{}) (string, string, map[string]interface{}) {
	return id, "core/wait_for", map[string]interface{}{
		"field":     field,
		"operator":  operator,
		"compareTo": compareTo,
		"interval":  "30s",
	}
}

// WaitForLiteral is a shorthand for WaitForNode with a literal compareTo value.
func WaitForLiteral(id, field, operator, value string) (string, string, map[string]interface{}) {
	return WaitForNode(id, field, operator, map[string]interface{}{
		"source": "literal",
		"value":  value,
	})
}

/**
 * FieldValue Helpers
 * Shorthand constructors for FieldValue maps used in node configs.
 */

// Literal creates a literal FieldValue.
func Literal(value string) map[string]interface{} {
	return map[string]interface{}{"type": "literal", "value": value}
}

// FromState creates a state FieldValue.
func FromState(field string) map[string]interface{} {
	return map[string]interface{}{"type": "state", "value": field}
}

// FromEvent creates an event FieldValue.
func FromEvent(path string) map[string]interface{} {
	return map[string]interface{}{"type": "event", "value": path}
}

// FromInput creates an input FieldValue.
func FromInput(field string) map[string]interface{} {
	return map[string]interface{}{"type": "input", "value": field}
}

// FromNodeOutput creates a nodeOutput FieldValue.
func FromNodeOutput(nodeId, path string) map[string]interface{} {
	return map[string]interface{}{"type": "nodeOutput", "value": path, "nodeId": nodeId}
}

func ptrObjID(id model.ObjectId) *model.ObjectId {
	return &id
}

