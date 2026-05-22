package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * VALUE TYPES (Domain Layer - MongoDB persistence)
 * Mirrored from contracts but with bson tags for persistence
 */

type FieldValueType string

const (
	FieldValueEvent      FieldValueType = "event"
	FieldValueState      FieldValueType = "state"
	FieldValueVariable   FieldValueType = "variable" // Legacy alias for state
	FieldValueInput      FieldValueType = "input"    // ExternalInputs
	FieldValueLiteral    FieldValueType = "literal"
	FieldValueNodeOutput FieldValueType = "nodeOutput"
)

type FieldValue struct {
	Type   FieldValueType `bson:"type"`
	Value  string         `bson:"value"`
	Mode   string         `bson:"mode,omitempty"`
	NodeID string         `bson:"nodeId,omitempty"`
}

/*
 * NODE STRUCTURES (Domain Layer - MongoDB persistence)
 */

type Position struct {
	X float64 `bson:"x"`
	Y float64 `bson:"y"`
}

// TimeoutConfig specifies async timeout duration and behavior at the node level.
// Lives alongside config (not inside it) — separate concern from node functionality.
type TimeoutConfig struct {
	Duration     int    `bson:"duration"`
	Unit         string `bson:"unit"`
	EnableOutput bool   `bson:"enableOutput"`
}

// ErrorHandlerConfig specifies retry policy at the node level.
// When enabled, retries the node with exponential backoff before following the "error" handle.
// Independent from timeout — they never interact.
type ErrorHandlerConfig struct {
	Enabled           bool    `bson:"enabled"`
	MaxAttempts       int     `bson:"maxAttempts"`
	InitialInterval   int     `bson:"initialInterval"`
	IntervalUnit      string  `bson:"intervalUnit"`
	BackoffMultiplier float64 `bson:"backoffMultiplier"`
}

type WorkflowNode struct {
	ID           string                 `bson:"id"`
	Type         string                 `bson:"type"`
	Label        string                 `bson:"label"`
	Position     Position               `bson:"position"`
	Config       map[string]interface{} `bson:"config"`
	Timeout      *TimeoutConfig         `bson:"timeout,omitempty"`
	ErrorHandler *ErrorHandlerConfig    `bson:"errorHandler,omitempty"`
	ParentNodeID string                 `bson:"parentNodeId"`
}

/*
 * EDGE STRUCTURES (Domain Layer - MongoDB persistence)
 */

type WorkflowEdge struct {
	ID           string  `bson:"id"`
	Source       string  `bson:"source"`
	SourceHandle string  `bson:"sourceHandle"`
	Target       string  `bson:"target"`
	TargetHandle string  `bson:"targetHandle"`
	Label        string  `bson:"label"`
	PathOffsetX  float64 `bson:"pathOffsetX"`
	PathOffsetY  float64 `bson:"pathOffsetY"`
}

/*
 * VARIABLE STRUCTURES (Domain Layer - MongoDB persistence)
 */

type VariableType string

const (
	VarTypeString  VariableType = "string"
	VarTypeNumber  VariableType = "number"
	VarTypeBoolean VariableType = "boolean"
	VarTypeJSON    VariableType = "json"
)

type WorkflowVariable struct {
	Field        string       `bson:"field"`
	Type         VariableType `bson:"type"`
	DefaultValue interface{}  `bson:"defaultValue"`
	Description  string       `bson:"description"`
	Durable      bool         `bson:"durable"`
}

type CaptureField struct {
	Field       string       `bson:"field"`
	Type        VariableType `bson:"type"`
	Description string       `bson:"description"`
}

type ExternalInput struct {
	Field           string       `bson:"field"`
	Label           string       `bson:"label"`
	Icon            string       `bson:"icon"`
	Type            VariableType `bson:"type"`
	Description     string       `bson:"description"`
	DefaultValue    interface{}  `bson:"defaultValue"`
	Required        bool         `bson:"required"`
	AssetTemplateId string       `bson:"assetTemplateId,omitempty"`
	FieldPath       string       `bson:"fieldPath,omitempty"`
}

type ExternalSignal struct {
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

/*
 * CONDITION SYSTEM (Domain Layer - MongoDB persistence)
 */

type GroupLogicOperator string

const (
	LogicAND  GroupLogicOperator = "AND"
	LogicOR   GroupLogicOperator = "OR"
	LogicNAND GroupLogicOperator = "NAND"
	LogicNOR  GroupLogicOperator = "NOR"
)

type ConditionItem struct {
	ID       string     `bson:"id"`
	Name     string     `bson:"name"`
	Field    FieldValue `bson:"field"`
	Operator string     `bson:"operator"`
	Value    FieldValue `bson:"value"`
}

type ConditionGroupItem struct {
	Type string      `bson:"type"`
	Data interface{} `bson:"data"`
}

type ConditionGroup struct {
	ID    string               `bson:"id"`
	Name  string               `bson:"name"`
	Logic GroupLogicOperator   `bson:"logic"`
	Items []ConditionGroupItem `bson:"items"`
}

type SwitchCase struct {
	ID        string         `bson:"id"`
	Name      string         `bson:"name"`
	Condition ConditionGroup `bson:"condition"`
}

/*
 * RETRY POLICY (Domain Layer - MongoDB persistence)
 */

type RetryPolicy struct {
	Enabled            bool     `bson:"enabled"`
	MaxAttempts        int      `bson:"maxAttempts"`
	InitialInterval    string   `bson:"initialInterval"`
	BackoffMultiplier  float64  `bson:"backoffMultiplier"`
	MaxInterval        string   `bson:"maxInterval"`
	NonRetryableErrors []string `bson:"nonRetryableErrors"`
}

/*
 * METADATA (Domain Layer - MongoDB persistence)
 */

type CanvasViewport struct {
	X    float64 `bson:"x"`
	Y    float64 `bson:"y"`
	Zoom float64 `bson:"zoom"`
}

type DefinitionMetadata struct {
	CanvasViewport CanvasViewport `bson:"canvasViewport"`
}

/*
 * DEFINITION STATUS CONSTANTS
 */

type DefinitionStatus string

const (
	StatusValid         DefinitionStatus = "valid"
	StatusPluginMissing DefinitionStatus = "plugin_missing"
	StatusInvalid       DefinitionStatus = "invalid"
)

/*
 * WORKFLOW DEFINITION (Domain Layer - Root Aggregate)
 */

type WorkflowDefinition struct {
	ID                model.ObjectId     `bson:"_id,omitempty"`
	OrgID             *model.ObjectId    `bson:"orgId"`
	Name              string             `bson:"name"`
	Description       string             `bson:"description"`
	Enabled           bool               `bson:"enabled"`
	IsTemplate        bool               `bson:"isTemplate"`
	DefinitionVersion int                `bson:"definitionVersion"`
	Timezone          FieldValue         `bson:"timezone"`
	RetryPolicy       RetryPolicy        `bson:"retryPolicy"`
	States            []WorkflowVariable `bson:"states"`
	CaptureFields     []CaptureField     `bson:"captureFields"`
	ExternalInputs    []ExternalInput    `bson:"externalInputs"`
	ExternalSignals   []ExternalSignal   `bson:"externalSignals"`
	Nodes             []WorkflowNode     `bson:"nodes"`
	Edges             []WorkflowEdge     `bson:"edges"`
	InstalledPlugins  []string           `bson:"installedPlugins"`
	MissingPlugins    []string           `bson:"missingPlugins"`
	Status            string             `bson:"status"`
	Metadata          DefinitionMetadata `bson:"metadata"`
	PathKey           string             `bson:"pathKey"`
	Scope             string             `bson:"scope"`
	Created           time.Time          `bson:"created"`
	Updated           time.Time          `bson:"updated"`
}

func (d *WorkflowDefinition) GetCreated() time.Time { return d.Created }
func (d *WorkflowDefinition) GetUpdated() time.Time { return d.Updated }

/*
 * WORKFLOW DEFINITION UPDATE (Domain Layer - Partial update)
 */

type WorkflowDefinitionUpdate struct {
	Name              *string             `bson:"name,omitempty"`
	Description       *string             `bson:"description,omitempty"`
	Enabled           *bool               `bson:"enabled,omitempty"`
	IsTemplate        *bool               `bson:"isTemplate,omitempty"`
	Timezone          *FieldValue         `bson:"timezone,omitempty"`
	RetryPolicy       *RetryPolicy        `bson:"retryPolicy,omitempty"`
	States            *[]WorkflowVariable `bson:"states,omitempty"`
	CaptureFields     *[]CaptureField     `bson:"captureFields,omitempty"`
	ExternalInputs    *[]ExternalInput    `bson:"externalInputs,omitempty"`
	ExternalSignals   *[]ExternalSignal   `bson:"externalSignals,omitempty"`
	Nodes             *[]WorkflowNode     `bson:"nodes,omitempty"`
	Edges             *[]WorkflowEdge     `bson:"edges,omitempty"`
	InstalledPlugins  *[]string           `bson:"installedPlugins,omitempty"`
	MissingPlugins    *[]string           `bson:"missingPlugins,omitempty"`
	Status            *string             `bson:"status,omitempty"`
	Metadata          *DefinitionMetadata `bson:"metadata,omitempty"`
	Scope             *string             `bson:"scope,omitempty"`
}
