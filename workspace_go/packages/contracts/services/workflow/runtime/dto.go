package runtime

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/*
 * INSTANCE STATUS
 */

type InstanceStatus string

const (
	StatusCreated   InstanceStatus = "created"
	StatusRunning   InstanceStatus = "running"
	StatusWaiting   InstanceStatus = "waiting"
	StatusCompleted InstanceStatus = "completed"
	StatusFailed    InstanceStatus = "failed"
	StatusCancelled InstanceStatus = "cancelled"
)

/*
 * DTOs (API Layer)
 */

// InstanceId represents the params DTO for /:instanceId
type InstanceId struct {
	InstanceId string `params:"instanceId" validate:"required"`
}

// InstanceQuery represents the query parameters for listing instances.
type InstanceQuery struct {
	query.BaseQueryDTO

	WorkflowID *string         `query:"workflowId" validate:"omitempty"`
	Status     *InstanceStatus `query:"status" validate:"omitempty"`
}

// InstanceResponse represents the API response for a workflow instance.
type InstanceResponse struct {
	ID                *common.ObjectID       `json:"_id,omitempty"`
	WorkflowID        *common.ObjectID       `json:"workflowId,omitempty"`
	WorkflowName      *string                `json:"workflowName,omitempty"`
	OrgID             *common.ObjectID       `json:"orgId,omitempty"`
	DefinitionVersion *int                   `json:"definitionVersion,omitempty"`
	Version           *int                   `json:"version,omitempty"`
	Status            *InstanceStatus        `json:"status,omitempty"`
	CurrentNodeID     *string                `json:"currentNodeId,omitempty"`
	State             map[string]interface{} `json:"state,omitempty"`
	EventPayload      map[string]interface{} `json:"eventPayload,omitempty"`
	ExecutionPath     []PathEntryResponse    `json:"executionPath,omitempty"`
	NodeOutputs       map[string]interface{} `json:"nodeOutputs,omitempty"`
	ErrorInfo         *ErrorInfoResponse     `json:"errorInfo,omitempty"`
	ParentInstanceID  *common.ObjectID       `json:"parentInstanceId,omitempty"`
	Depth             *int                   `json:"depth,omitempty"`
	StartedAt         *common.NullTime       `json:"startedAt,omitempty"`
	CompletedAt       *common.NullTime       `json:"completedAt,omitempty"`
	Created           *common.NullTime       `json:"created,omitempty"`
	Updated           *common.NullTime       `json:"updated,omitempty"`
}

func (r *InstanceResponse) SetCreated(t *common.NullTime) { r.Created = t }
func (r *InstanceResponse) SetUpdated(t *common.NullTime) { r.Updated = t }

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

// SignalRequest represents the payload for sending a signal to a waiting workflow instance.
// Used by the HTTP API (POST /instances/:id/signal) — always targets a specific instance.
type SignalRequest struct {
	SignalName string                 `json:"signalName" validate:"required"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// Deprecated: Use executions.WorkflowExecutionMessage instead.
// WorkflowDeliverMessage is kept for backward compatibility.
type WorkflowDeliverMessage = WorkflowExecutionMessage

// WorkflowExecutionMessage is the NATS message payload published to the WORKFLOW-EXECUTION stream.
type WorkflowExecutionMessage struct {
	Mode     string                 `json:"mode" validate:"required,oneof=newInstance signal signalOrStart"`
	Event    map[string]interface{} `json:"event,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Data     map[string]interface{} `json:"data"`
}

// CodeExecutionRequest is published to WORKFLOW-JS-CODE stream when a code node
// suspends with waitType "callback". The JS Workflow Executor consumes this and runs the script.
type CodeExecutionRequest struct {
	OrgID           string                 `json:"orgId"`
	PathKey         string                 `json:"pathKey"`
	WorkflowID      string                 `json:"workflowId"`
	NodeID          string                 `json:"nodeId"`
	InstanceID      string                 `json:"instanceId"`
	CallbackSubject string                 `json:"callbackSubject"`
	ExecutionToken  string                 `json:"executionToken,omitempty"`
	Timeout         int                    `json:"timeout"`
	EventPayload    map[string]interface{} `json:"eventPayload"`
	State           map[string]interface{} `json:"state"`
	Inputs          map[string]interface{} `json:"inputs"`
	NodeOutputs     map[string]interface{} `json:"nodeOutputs"`
}

// WorkflowTriggerRequest is published to the TRIGGERS stream (subject: trigger.WORKFLOW.execute).
// Unified message for all workflow-to-triggers communication.
//
// mode "trigger": Workflow uses a registered trigger entity (triggerId in data).
//   The Triggers Service fetches config by triggerId and executes.
//
// mode "plugin": Workflow uses a plugin action (fully resolved pipeline in data).
//   The Triggers Service executes the pipeline: before hook → operation → after hook.
//
// Both modes include callbackSubject for resume after execution.
type WorkflowTriggerRequest struct {
	Mode            string                 `json:"mode"` // "trigger" or "plugin"
	OrgID           string                 `json:"orgId"`
	PathKey         string                 `json:"pathKey"`
	WorkflowID      string                 `json:"workflowId"`
	InstanceID      string                 `json:"instanceId"`
	NodeID          string                 `json:"nodeId"`
	CallbackSubject string                 `json:"callbackSubject"`
	ExecutionToken  string                 `json:"executionToken,omitempty"`
	Data            map[string]interface{} `json:"data"` // Mode-specific payload
}
