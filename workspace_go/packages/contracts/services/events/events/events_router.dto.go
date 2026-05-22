package dtos

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * Router Events DTOs
 * Shared between Router service (publisher) and Events service (consumer)
 */

// RouterEventDTO is the payload published by Router service to NATS.
// Events service consumes this to store routing execution history.
type RouterEventDTO struct {
	Created        time.Time         `json:"created"`
	EventTrackerId string            `json:"eventTrackerId"`
	ThreadId       string            `json:"threadId"`
	OrgId          string            `json:"orgId"`
	PathKey        string            `json:"pathKey"`
	AssetId        string            `json:"assetId"`
	RouterId       string            `json:"routerId"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Routers        []RouterResultDTO `json:"routers"`
}

// RouterResultDTO represents a single router execution result.
type RouterResultDTO struct {
	Kind       string               `json:"kind"`
	Matched    bool                 `json:"matched"`
	Published  bool                 `json:"published"`
	Conditions []ConditionResultDTO `json:"conditions"`
}

// ConditionResultDTO represents a single condition evaluation result.
type ConditionResultDTO struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
	Passed   bool        `json:"passed"`
}

// EventsRouterQuery represents query parameters for listing router events.
// Uses cursor-based pagination for efficient querying.
type EventsRouterQuery struct {
	query.CursorQueryDTO

	// Filters
	EventTrackerId *string    `query:"eventTrackerId" validate:"omitempty"`
	ThreadId       *string    `query:"threadId" validate:"omitempty"`
	AssetId        *string    `query:"assetId" validate:"omitempty"`
	RouterId       *string    `query:"routerId" validate:"omitempty"`
	Success        *bool      `query:"success" validate:"omitempty"`
	PublishedCount *int       `query:"publishedCount" validate:"omitempty"`
	StartTime      *time.Time `query:"startTime" validate:"omitempty"`
	EndTime        *time.Time `query:"endTime" validate:"omitempty"`
}

// EventsRouterResponse represents a router event response.
type EventsRouterResponse struct {
	Created        time.Time `json:"created"`
	EventTrackerId string    `json:"eventTrackerId,omitempty"`
	ThreadId       string    `json:"threadId"`
	OrgId          string    `json:"orgId"`
	PathKey        string    `json:"pathKey,omitempty"`
	AssetId        string    `json:"assetId"`
	RouterId       string    `json:"routerId"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	TotalRouters   uint8     `json:"totalRouters"`
	MatchedCount   uint8     `json:"matchedCount"`
	PublishedCount uint8     `json:"publishedCount"`
	Event          string    `json:"event"`
	Success        bool      `json:"success"`
	Error          string    `json:"error,omitempty"`
	RetentionDays  uint16    `json:"retentionDays"`
}

// EventsRouterCursorResult represents the cursor-paginated response for router events.
type EventsRouterCursorResult struct {
	Items       []EventsRouterResponse `json:"items"`
	NextCursor  *time.Time             `json:"nextCursor,omitempty"`
	PrevCursor  *time.Time             `json:"prevCursor,omitempty"`
	HasNext     bool                   `json:"hasNext"`
	HasPrevious bool                   `json:"hasPrevious"`
}
