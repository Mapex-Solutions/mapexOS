package executions

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/*
 * EXECUTION STATUS
 */

type ExecutionStatus string

const (
	StatusCreated   ExecutionStatus = "created"
	StatusRunning   ExecutionStatus = "running"
	StatusWaiting   ExecutionStatus = "waiting"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
	StatusCancelled ExecutionStatus = "cancelled"
)

/*
 * DTOs (API Layer) for Workflow Executions
 */

// ExecutionId represents the params DTO for /:executionId
type ExecutionId struct {
	ExecutionId string `params:"executionId" validate:"required"`
}

// ExecutionQuery represents the query parameters for listing executions.
// Status supports comma-separated values (e.g., "running,waiting") for $in queries.
type ExecutionQuery struct {
	query.BaseQueryDTO

	InstanceID   *string `query:"instanceId" validate:"omitempty"`
	DefinitionID *string `query:"definitionId" validate:"omitempty"`
	Status       *string `query:"status" validate:"omitempty"`
}

// ExecutionResponse represents the API response for a workflow execution.
type ExecutionResponse struct {
	ID                *string                `json:"_id,omitempty"`
	WorkflowUUID      *string                `json:"workflowUUID,omitempty"`
	InstanceID        *common.ObjectID       `json:"instanceId,omitempty"`
	DefinitionID      *common.ObjectID       `json:"definitionId,omitempty"`
	WorkflowName      *string                `json:"workflowName,omitempty"`
	InstanceName      *string                `json:"instanceName,omitempty"`
	DefinitionName    *string                `json:"definitionName,omitempty"`
	OrgID             *common.ObjectID       `json:"orgId,omitempty"`
	Version           *int                   `json:"version,omitempty"`
	Status            *ExecutionStatus       `json:"status,omitempty"`
	ActiveNodeIDs     []string               `json:"activeNodeIds,omitempty"`
	State             map[string]interface{} `json:"state,omitempty"`
	EventPayload      map[string]interface{} `json:"eventPayload,omitempty"`
	ExecutionPath     []PathEntryResponse    `json:"executionPath,omitempty"`
	NodeOutputs       map[string]interface{} `json:"nodeOutputs,omitempty"`
	ErrorInfo         *ErrorInfoResponse     `json:"errorInfo,omitempty"`
	TriggerSource     *string                `json:"triggerSource,omitempty"`
	ParentExecutionID *string                `json:"parentExecutionId,omitempty"`
	Depth             *int                   `json:"depth,omitempty"`
	StartedAt         *common.NullTime       `json:"startedAt,omitempty"`
	CompletedAt       *common.NullTime       `json:"completedAt,omitempty"`
	Created           *common.NullTime       `json:"created,omitempty"`
	Updated           *common.NullTime       `json:"updated,omitempty"`
}

func (r *ExecutionResponse) SetCreated(t *common.NullTime) { r.Created = t }
func (r *ExecutionResponse) SetUpdated(t *common.NullTime) { r.Updated = t }

/*
 * SUPPORTING RESPONSE TYPES
 */

type PathEntryResponse struct {
	NodeID       string           `json:"nodeId"`
	NodeType     string           `json:"nodeType"`
	Status       string           `json:"status"`
	EnteredAt    *common.NullTime `json:"enteredAt,omitempty"`
	ExitedAt     *common.NullTime `json:"exitedAt,omitempty"`
	DurationMs   int64            `json:"durationMs"`
	OutputHandle string           `json:"outputHandle,omitempty"`
	Error        *string          `json:"error,omitempty"`
}

type ErrorInfoResponse struct {
	Code       string           `json:"code"`
	Message    string           `json:"message"`
	NodeID     string           `json:"nodeId"`
	NodeType   string           `json:"nodeType"`
	Timestamp  *common.NullTime `json:"timestamp,omitempty"`
	StackTrace string           `json:"stackTrace,omitempty"`
}

// SignalRequest represents the payload for sending a signal to a waiting execution.
type SignalRequest struct {
	SignalName string                 `json:"signalName" validate:"required"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// WorkflowExecutionMessage is the NATS message payload published to the WORKFLOW-EXECUTION stream.
// Routing is done by the "mode" field in the payload — never by subject.
//
// Modes:
//   - "newInstance": Create a new execution from an instance config.
//     data: { instanceId (required), workflowUUID (optional) }
//
//   - "signal": Deliver a signal to a waiting execution.
//     data: { workflowUUID (required), signalName (required), signalData (optional) }
//
//   - "signalOrStart": Try signal first, fall back to newInstance.
//     data: { instanceId, workflowUUID, signalName, signalData }
type WorkflowExecutionMessage struct {
	// Dispatch mode: "newInstance", "signal", "signalOrStart", "subworkflow"
	Mode string `json:"mode" validate:"required,oneof=newInstance signal signalOrStart subworkflow"`

	// Event payload from the Router (original sensor/trigger data).
	Event map[string]interface{} `json:"event,omitempty"`

	// Router metadata (routerId, matchRuleId, timestamp, etc.).
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Mode-specific configuration. Schema defined per mode.
	Data map[string]interface{} `json:"data"`
}

// NewInstanceData is the typed data for mode "newInstance".
type NewInstanceData struct {
	InstanceID   string `json:"instanceId"`
	WorkflowUUID string `json:"workflowUUID,omitempty"`
}

// SignalData is the typed data for mode "signal".
type SignalData struct {
	WorkflowUUID string                 `json:"workflowUUID"`
	SignalName   string                 `json:"signalName"`
	SignalData   map[string]interface{} `json:"signalData,omitempty"`
}

// SignalOrStartData is the typed data for mode "signalOrStart".
type SignalOrStartData struct {
	InstanceID   string                 `json:"instanceId"`
	WorkflowUUID string                 `json:"workflowUUID"`
	SignalName   string                 `json:"signalName"`
	SignalData   map[string]interface{} `json:"signalData,omitempty"`
}
