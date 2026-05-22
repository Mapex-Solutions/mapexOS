package entities

import (
	defPorts "workflow/src/modules/definitions/application/ports"
)

// Node config structs are parsed from node Config map[string]interface{} by GraphBuilder.
// Each executor receives its typed config via NodeExecutionContext.ParsedConfig.
// Field names mirror the Config map keys from the frontend DSL — actual key matching
// is performed manually in services/config_parsing.go (no json.Unmarshal involved).

// ConditionNodeConfig holds parsed configuration for a condition node.
type ConditionNodeConfig struct {
	Condition           defPorts.ConditionGroup
	SelectedTemplateIds []string
}

// SwitchNodeConfig holds parsed configuration for a switch node.
type SwitchNodeConfig struct {
	Cases               []defPorts.SwitchCase
	MatchMode           string
	SelectedTemplateIds []string
}

// FanoutNodeConfig holds parsed configuration for a fanout (fork) node.
type FanoutNodeConfig struct {
	Branches int
	Mode     string // "waitAll" (default) | "firstCompleted"
}

// MergeNodeConfig holds parsed configuration for a merge (join) node.
type MergeNodeConfig struct {
	Branches int
	Strategy string
}

// SequenceNodeConfig holds parsed configuration for a sequence node.
type SequenceNodeConfig struct {
	Steps int
}

// LoopNodeConfig holds parsed configuration for a loop node.
type LoopNodeConfig struct {
	Source defPorts.FieldValue
}

// GotoNodeConfig holds parsed configuration for a goto portal node.
// Role is "sender" or "receiver"; PairLabel links matching pairs.
type GotoNodeConfig struct {
	Role      string
	PairLabel string
	PairColor string
}

// SetStateNodeConfig holds parsed configuration for a set_state node.
type SetStateNodeConfig struct {
	Operation           string
	TargetField         string
	ValueSource         defPorts.FieldValue
	SelectedTemplateIds []string
}

// LogNodeConfig holds parsed configuration for a log node.
type LogNodeConfig struct {
	Message string
	Level   string
}

// CodeNodeConfig holds parsed configuration for a code execution node.
type CodeNodeConfig struct {
	Script  string
	Timeout int
}

// DelayNodeConfig holds parsed configuration for a delay node.
type DelayNodeConfig struct {
	Duration int
	Unit     string
}

// WaitSignalNodeConfig holds parsed configuration for a wait_signal node.
// Timeout is now managed at node level (node.timeout), not inside config.
type WaitSignalNodeConfig struct {
	SignalName string
	Mappings   []SignalMapping
}

// SignalMapping maps a signal parameter name to a resolved field value.
type SignalMapping struct {
	ParamName string
	Value     defPorts.FieldValue
}

// WaitForNodeConfig holds parsed configuration for a wait_for (polling condition) node.
// Timeout is now managed at node level (node.timeout), not inside config.
type WaitForNodeConfig struct {
	Field     string
	Operator  string
	CompareTo defPorts.FieldValue
}

// SubworkflowNodeConfig holds parsed configuration for a subworkflow node.
type SubworkflowNodeConfig struct {
	WorkflowID     string
	WorkflowName   string
	ExecutionMode  string
	Timeout        TimeoutConfig
	InputMappings  []InputMapping
	OutputMappings []OutputMapping
}

// TimeoutConfig specifies a duration and unit for async timeouts.
// When EnableOutput is true, timeout expiry routes to a "timeout" output handle
// instead of failing the execution with TIMEOUT_EXCEEDED.
type TimeoutConfig struct {
	Duration     int    `bson:"duration"`
	Unit         string `bson:"unit"`
	EnableOutput bool   `bson:"enableOutput"`
}

// InputMapping maps a child workflow parameter to a resolved field value.
type InputMapping struct {
	ChildParamName string
	Value          defPorts.FieldValue
}

// OutputMapping maps a child workflow output to a parent state field.
type OutputMapping struct {
	OutputName string
	StateField string
}

// TriggerEventNodeConfig holds parsed configuration for a trigger_event node.
type TriggerEventNodeConfig struct {
	EventType      string
	PayloadMapping []TriggerPayloadField
}

// TriggerPayloadField maps a key to a resolved field value in the trigger event payload.
type TriggerPayloadField struct {
	Key   string
	Value defPorts.FieldValue
}

// EndNodeConfig holds parsed configuration for an end node.
type EndNodeConfig struct {
	TerminateWithError bool
	ErrorCode          string
	ErrorMessage       defPorts.FieldValue
}

// PluginNodeConfig holds parsed configuration for a marketplace plugin node.
// Extracted by parsePluginNodeConfig — preserves the raw config for template resolution at dispatch time.
type PluginNodeConfig struct {
	Operation    string
	CredentialID string
	RawConfig    map[string]interface{}
}
