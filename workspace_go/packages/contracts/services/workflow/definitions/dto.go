package definitions

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * VALUE TYPES (API Layer)
 * These types mirror the domain entities but with json tags for API serialization.
 */

type FieldValueType string

const (
	FieldValueEvent      FieldValueType = "event"
	FieldValueState      FieldValueType = "state"
	FieldValueVariable   FieldValueType = "variable"
	FieldValueLiteral    FieldValueType = "literal"
	FieldValueNodeOutput FieldValueType = "nodeOutput"
	FieldValueEngine     FieldValueType = "engine"
)

type FieldValue struct {
	Type   FieldValueType `json:"type"`
	Value  string         `json:"value"`
	Mode   string         `json:"mode,omitempty"`
	NodeID string         `json:"nodeId,omitempty"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeTimeoutConfig specifies async timeout at the node level.
type NodeTimeoutConfig struct {
	Duration     int    `json:"duration"`
	Unit         string `json:"unit"`
	EnableOutput bool   `json:"enableOutput"`
}

// NodeErrorHandlerConfig specifies retry policy at the node level.
// When enabled, retries the node with exponential backoff before following the "error" handle.
type NodeErrorHandlerConfig struct {
	Enabled           bool    `json:"enabled"`
	MaxAttempts       int     `json:"maxAttempts"`
	InitialInterval   int     `json:"initialInterval"`
	IntervalUnit      string  `json:"intervalUnit"`
	BackoffMultiplier float64 `json:"backoffMultiplier"`
}

type WorkflowNode struct {
	ID           string                  `json:"id"`
	Type         string                  `json:"type"`
	Label        string                  `json:"label"`
	Position     Position                `json:"position"`
	Config       map[string]interface{}  `json:"config"`
	Timeout      *NodeTimeoutConfig      `json:"timeout,omitempty"`
	ErrorHandler *NodeErrorHandlerConfig `json:"errorHandler,omitempty"`
	ParentNodeID string                  `json:"parentNodeId"`
}

type WorkflowEdge struct {
	ID           string  `json:"id"`
	Source       string  `json:"source"`
	SourceHandle string  `json:"sourceHandle"`
	Target       string  `json:"target"`
	TargetHandle string  `json:"targetHandle"`
	Label        string  `json:"label"`
	PathOffsetX  float64 `json:"pathOffsetX"`
	PathOffsetY  float64 `json:"pathOffsetY"`
}

type VariableType string

const (
	VarTypeString  VariableType = "string"
	VarTypeNumber  VariableType = "number"
	VarTypeBoolean VariableType = "boolean"
	VarTypeJSON    VariableType = "json"
)

type WorkflowVariable struct {
	Field        string       `json:"field"`
	Type         VariableType `json:"type"`
	DefaultValue interface{}  `json:"defaultValue"`
	Description  string       `json:"description"`
	Durable      bool         `json:"durable"`
}

type CaptureField struct {
	Field       string       `json:"field"`
	Type        VariableType `json:"type"`
	Description string       `json:"description"`
}

type ExternalInput struct {
	Field           string       `json:"field"`
	Label           string       `json:"label"`
	Icon            string       `json:"icon"`
	Type            VariableType `json:"type"`
	Description     string       `json:"description"`
	DefaultValue    interface{}  `json:"defaultValue"`
	Required        bool         `json:"required"`
	AssetTemplateId string       `json:"assetTemplateId,omitempty"`
	FieldPath       string       `json:"fieldPath,omitempty"`
}

type ExternalSignal struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupLogicOperator string

const (
	LogicAND  GroupLogicOperator = "AND"
	LogicOR   GroupLogicOperator = "OR"
	LogicNAND GroupLogicOperator = "NAND"
	LogicNOR  GroupLogicOperator = "NOR"
)

type ConditionItem struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Field    FieldValue `json:"field"`
	Operator string     `json:"operator"`
	Value    FieldValue `json:"value"`
}

type ConditionGroupItem struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ConditionGroup struct {
	ID    string               `json:"id"`
	Name  string               `json:"name"`
	Logic GroupLogicOperator   `json:"logic"`
	Items []ConditionGroupItem `json:"items"`
}

type SwitchCase struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Condition ConditionGroup `json:"condition"`
}

type RetryPolicy struct {
	Enabled            bool     `json:"enabled"`
	MaxAttempts        int      `json:"maxAttempts"`
	InitialInterval    string   `json:"initialInterval"`
	BackoffMultiplier  float64  `json:"backoffMultiplier"`
	MaxInterval        string   `json:"maxInterval"`
	NonRetryableErrors []string `json:"nonRetryableErrors"`
}

type CanvasViewport struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}

type DefinitionMetadata struct {
	CanvasViewport CanvasViewport `json:"canvasViewport"`
}

/**
 * DTOs (API Layer)
 */

// DefinitionId represents the params DTO for /:workflowId
type DefinitionId struct {
	WorkflowId string `params:"workflowId" validate:"required"`
}

// DefinitionCreate represents the payload to create a workflow definition.
type DefinitionCreate struct {
	Name          string             `json:"name"          validate:"required,min=1,max=255"`
	Description   string             `json:"description"   validate:"max=1000"`
	Enabled       bool               `json:"enabled"`
	IsTemplate    bool               `json:"isTemplate"`
	Timezone      FieldValue         `json:"timezone"`
	RetryPolicy   RetryPolicy        `json:"retryPolicy"`
	States           []WorkflowVariable `json:"states"`
	CaptureFields    []CaptureField     `json:"captureFields"`
	ExternalInputs   []ExternalInput    `json:"externalInputs"`
	ExternalSignals  []ExternalSignal   `json:"externalSignals"`
	Nodes            []WorkflowNode     `json:"nodes"             validate:"required,min=1"`
	Edges            []WorkflowEdge     `json:"edges"`
	InstalledPlugins []string           `json:"installedPlugins"`
	Metadata         DefinitionMetadata `json:"metadata"`
	Scope            string             `json:"scope"`

	// Multi-tenant fields (populated by coverage middleware / service)
	OrgID   *model.ObjectId  `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string          `json:"pathKey,omitempty" validate:"omitempty"`
	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

// DefinitionUpdate represents the payload to update a workflow definition.
type DefinitionUpdate struct {
	Name          *string              `json:"name,omitempty"          validate:"omitempty,min=1,max=255"`
	Description   *string              `json:"description,omitempty"   validate:"omitempty,max=1000"`
	Enabled       *bool                `json:"enabled,omitempty"`
	IsTemplate    *bool                `json:"isTemplate,omitempty"`
	Timezone      *FieldValue          `json:"timezone,omitempty"`
	RetryPolicy   *RetryPolicy         `json:"retryPolicy,omitempty"`
	States           *[]WorkflowVariable  `json:"states,omitempty"`
	CaptureFields    *[]CaptureField      `json:"captureFields,omitempty"`
	ExternalInputs   *[]ExternalInput     `json:"externalInputs,omitempty"`
	ExternalSignals  *[]ExternalSignal    `json:"externalSignals,omitempty"`
	Nodes            *[]WorkflowNode      `json:"nodes,omitempty"`
	Edges            *[]WorkflowEdge     `json:"edges,omitempty"`
	InstalledPlugins *[]string           `json:"installedPlugins,omitempty"`
	Metadata         *DefinitionMetadata `json:"metadata,omitempty"`
	Scope            *string             `json:"scope,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

// DefinitionQuery represents the query parameters for listing definitions.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
type DefinitionQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name              *string `query:"name" validate:"omitempty,max=100"`
	Enabled           *bool   `query:"enabled" validate:"omitempty"`
	Status            *string `query:"status" validate:"omitempty,oneof=valid plugin_missing invalid"`
	IsTemplate        *bool   `query:"isTemplate" validate:"omitempty"`
	DefinitionVersion *int    `query:"definitionVersion" validate:"omitempty,min=1"`
}

// DefinitionResponse represents the API response for a workflow definition.
type DefinitionResponse struct {
	ID                *common.ObjectID    `json:"_id,omitempty"`
	OrgID             *common.ObjectID    `json:"orgId,omitempty"`
	Name              *string             `json:"name,omitempty"`
	Description       *string             `json:"description,omitempty"`
	Enabled           *bool               `json:"enabled,omitempty"`
	IsTemplate        *bool               `json:"isTemplate,omitempty"`
	DefinitionVersion *int                `json:"definitionVersion,omitempty"`
	Timezone          *FieldValue         `json:"timezone,omitempty"`
	RetryPolicy       *RetryPolicy        `json:"retryPolicy,omitempty"`
	States            []WorkflowVariable  `json:"states,omitempty"`
	CaptureFields     []CaptureField      `json:"captureFields,omitempty"`
	ExternalInputs    []ExternalInput     `json:"externalInputs,omitempty"`
	ExternalSignals   []ExternalSignal    `json:"externalSignals,omitempty"`
	Nodes             []WorkflowNode      `json:"nodes,omitempty"`
	Edges             []WorkflowEdge      `json:"edges,omitempty"`
	InstalledPlugins  []string            `json:"installedPlugins,omitempty"`
	MissingPlugins    []string            `json:"missingPlugins,omitempty"`
	Status            *string             `json:"status,omitempty"`
	Metadata          *DefinitionMetadata `json:"metadata,omitempty"`
	PathKey           *string             `json:"pathKey,omitempty"`
	Scope             *string             `json:"scope,omitempty"`
	Created           *common.NullTime    `json:"created,omitempty"`
	Updated           *common.NullTime    `json:"updated,omitempty"`
}

func (d *DefinitionResponse) SetCreated(t *common.NullTime) { d.Created = t }
func (d *DefinitionResponse) SetUpdated(t *common.NullTime) { d.Updated = t }
