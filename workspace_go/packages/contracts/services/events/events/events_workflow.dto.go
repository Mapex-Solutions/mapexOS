package dtos

import (
	"encoding/json"
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * Workflow Events DTOs
 * Shared between Workflow service (publisher/Archiver) and Events service (consumer)
 */

// WorkflowEventDTO is the payload published by Workflow Archiver to NATS.
// Events service consumes this to store workflow execution history in ClickHouse.
type WorkflowEventDTO struct {
	Created           time.Time `json:"created"`
	Finished          time.Time `json:"finished"`
	ExecutionId       string    `json:"executionId"`
	EventTrackerId    string    `json:"eventTrackerId,omitempty"`
	OrgId             string    `json:"orgId"`
	PathKey           string    `json:"pathKey"`
	WorkflowUUID      string    `json:"workflowUUID,omitempty"`
	InstanceId        string    `json:"instanceId"`
	DefinitionId      string    `json:"definitionId"`
	WorkflowName      string    `json:"workflowName"`
	InstanceName      string    `json:"instanceName"`
	DefinitionName    string    `json:"definitionName"`
	Status            string    `json:"status"`
	Success           bool      `json:"success"`
	DurationMs        int64     `json:"durationMs"`
	ErrorMessage      string    `json:"errorMessage,omitempty"`
	ExecutionPath     string    `json:"executionPath,omitempty"`
	NodeOutputs       string    `json:"nodeOutputs,omitempty"`
	ErrorInfo         string    `json:"errorInfo,omitempty"`
	EventPayload      string    `json:"eventPayload,omitempty"`
	TriggerSource     string    `json:"triggerSource,omitempty"`
	ParentExecutionId string    `json:"parentExecutionId,omitempty"`
	Depth             uint8     `json:"depth"`
	RetentionDays     int16     `json:"retentionDays"`
	State             string    `json:"state,omitempty"`
	ExternalInputs    string    `json:"externalInputs,omitempty"`
}

// EventsWorkflowExecutionIdParam represents the URL param for getting a single workflow event by executionId.
type EventsWorkflowExecutionIdParam struct {
	ExecutionId string `params:"executionId" validate:"required"`
}

// EventsWorkflowQuery represents query parameters for listing workflow events.
// Uses cursor-based pagination for efficient querying.
type EventsWorkflowQuery struct {
	query.CursorQueryDTO

	// Filters
	EventTrackerId *string    `query:"eventTrackerId" validate:"omitempty"`
	WorkflowUUID   *string    `query:"workflowUUID" validate:"omitempty"`
	InstanceId     *string    `query:"instanceId" validate:"omitempty"`
	DefinitionId   *string    `query:"definitionId" validate:"omitempty"`
	Status         *string    `query:"status" validate:"omitempty"`
	Success        *bool      `query:"success" validate:"omitempty"`
	StartTime      *time.Time `query:"startTime" validate:"omitempty"`
	EndTime        *time.Time `query:"endTime" validate:"omitempty"`
}

// EventsWorkflowResponse represents a workflow event response.
// ExecutionPath, NodeOutputs, ErrorInfo, EventPayload are json.RawMessage
// so they serialize as JSON objects (not escaped strings) in the API response.
type EventsWorkflowResponse struct {
	Created           time.Time       `json:"created"`
	Finished          time.Time       `json:"finished"`
	ExecutionId       string          `json:"executionId"`
	EventTrackerId    string          `json:"eventTrackerId,omitempty"`
	OrgId             string          `json:"orgId"`
	PathKey           string          `json:"pathKey,omitempty"`
	WorkflowUUID      string          `json:"workflowUUID,omitempty"`
	InstanceId        string          `json:"instanceId"`
	DefinitionId      string          `json:"definitionId"`
	WorkflowName      string          `json:"workflowName"`
	InstanceName      string          `json:"instanceName"`
	DefinitionName    string          `json:"definitionName"`
	Status            string          `json:"status"`
	Success           bool            `json:"success"`
	DurationMs        int64           `json:"durationMs"`
	ErrorMessage      string          `json:"errorMessage,omitempty"`
	ExecutionPath     json.RawMessage `json:"executionPath,omitempty"`
	NodeOutputs       json.RawMessage `json:"nodeOutputs,omitempty"`
	ErrorInfo         json.RawMessage `json:"errorInfo,omitempty"`
	EventPayload      json.RawMessage `json:"eventPayload,omitempty"`
	TriggerSource     string          `json:"triggerSource,omitempty"`
	ParentExecutionId string          `json:"parentExecutionId,omitempty"`
	Depth             uint8           `json:"depth"`
	RetentionDays     int16           `json:"retentionDays"`
	State             json.RawMessage `json:"state,omitempty"`
	ExternalInputs    json.RawMessage `json:"externalInputs,omitempty"`
}

// EventsWorkflowCursorResult represents the cursor-paginated response for workflow events.
type EventsWorkflowCursorResult struct {
	Items       []EventsWorkflowResponse `json:"items"`
	NextCursor  *time.Time               `json:"nextCursor,omitempty"`
	PrevCursor  *time.Time               `json:"prevCursor,omitempty"`
	HasNext     bool                     `json:"hasNext"`
	HasPrevious bool                     `json:"hasPrevious"`
}
