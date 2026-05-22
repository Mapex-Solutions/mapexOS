package dtos

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * Trigger Events DTOs
 * Shared between Triggers service (publisher) and Events service (consumer)
 */

// TriggerEventDTO is the payload published by Triggers service to NATS.
// Events service consumes this to store trigger execution history.
type TriggerEventDTO struct {
	Created        time.Time `json:"created"`
	EventTrackerId string    `json:"eventTrackerId"`
	OrgId          string    `json:"orgId"`
	PathKey        string    `json:"pathKey"`
	TriggerId      string    `json:"triggerId"`
	TriggerName    string    `json:"triggerName"`
	TriggerType    string    `json:"triggerType"`  // http, mqtt, rabbitmq, nats, websocket, email, teams, slack
	Category       string    `json:"category"`     // technical, communication
	Source         string    `json:"source"`       // router
	Success        bool      `json:"success"`
	DurationMs     int64     `json:"durationMs"`
	Error          string    `json:"error,omitempty"`
	RequestData    string    `json:"requestData,omitempty"`  // JSON of resolved config sent to trigger
	ResponseData   string    `json:"responseData,omitempty"` // JSON of response/result from trigger
}

// EventsTriggerQuery represents query parameters for listing trigger events.
// Uses cursor-based pagination for efficient querying.
type EventsTriggerQuery struct {
	query.CursorQueryDTO

	// Filters
	EventTrackerId *string    `query:"eventTrackerId" validate:"omitempty"`
	TriggerId      *string    `query:"triggerId" validate:"omitempty"`
	TriggerType    *string    `query:"triggerType" validate:"omitempty"`
	Category       *string    `query:"category" validate:"omitempty"`
	Source         *string    `query:"source" validate:"omitempty"`
	Success        *bool      `query:"success" validate:"omitempty"`
	StartTime      *time.Time `query:"startTime" validate:"omitempty"`
	EndTime        *time.Time `query:"endTime" validate:"omitempty"`
}

// EventsTriggerResponse represents a trigger event response.
type EventsTriggerResponse struct {
	Created        time.Time `json:"created"`
	EventTrackerId string    `json:"eventTrackerId,omitempty"`
	OrgId          string    `json:"orgId"`
	PathKey        string    `json:"pathKey,omitempty"`
	TriggerId      string    `json:"triggerId"`
	TriggerName    string    `json:"triggerName"`
	TriggerType    string    `json:"triggerType"`
	Category       string    `json:"category"`
	Source         string    `json:"source"`
	Success        bool      `json:"success"`
	DurationMs     int64     `json:"durationMs"`
	Error          string    `json:"error,omitempty"`
	RequestData    string    `json:"requestData,omitempty"`
	ResponseData   string    `json:"responseData,omitempty"`

	RetentionDays uint16 `json:"retentionDays"`
}

// EventsTriggerCursorResult represents the cursor-paginated response for trigger events.
type EventsTriggerCursorResult struct {
	Items       []EventsTriggerResponse `json:"items"`
	NextCursor  *time.Time              `json:"nextCursor,omitempty"`
	PrevCursor  *time.Time              `json:"prevCursor,omitempty"`
	HasNext     bool                    `json:"hasNext"`
	HasPrevious bool                    `json:"hasPrevious"`
}
