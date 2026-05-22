package dtos

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * Business Rule Events DTOs
 */

// Events service consumes this to store business rule execution history.
type BusinessRuleEventDTO struct {
	Created                 time.Time `json:"created"`
	EventTrackerId          string    `json:"eventTrackerId"`
	ThreadId                string    `json:"threadId"`
	OrgId                   string    `json:"orgId"`
	PathKey                 string    `json:"pathKey"`
	RuleId                  string    `json:"ruleId"`
	BusinessRuleId          string    `json:"businessRuleId"`
	BusinessRuleName        string    `json:"businessRuleName"`
	BusinessRuleDescription string    `json:"businessRuleDescription"`

	// Execution result
	Matched    bool  `json:"matched"`
	DurationMs int64 `json:"durationMs"`

	// Evaluation metrics
	ConditionsEvaluated int `json:"conditionsEvaluated"`
	ConditionsMatched   int `json:"conditionsMatched"`
	GroupsEvaluated     int `json:"groupsEvaluated"`
	MaxDepthReached     int `json:"maxDepthReached"`

	// State data (JSON encoded)
	FinalState   map[string]interface{} `json:"finalState,omitempty"`
	StateChanges map[string]interface{} `json:"stateChanges,omitempty"`

	// Detailed logs (JSON encoded)
	EvaluationTree interface{} `json:"evaluationTree,omitempty"`
	ConditionLogs  interface{} `json:"conditionLogs,omitempty"`

	// Actions (JSON encoded)
	ActionsToDispatch interface{} `json:"actionsToDispatch,omitempty"`
}

// EventsBusinessRuleQuery represents query parameters for listing business rule events.
// Uses cursor-based pagination for efficient querying.
type EventsBusinessRuleQuery struct {
	query.CursorQueryDTO

	// Filters
	EventTrackerId *string    `query:"eventTrackerId" validate:"omitempty"`
	ThreadId       *string    `query:"threadId" validate:"omitempty"`
	RuleId         *string    `query:"ruleId" validate:"omitempty"`
	BusinessRuleId *string    `query:"businessRuleId" validate:"omitempty"`
	Matched        *bool      `query:"matched" validate:"omitempty"`
	StartTime      *time.Time `query:"startTime" validate:"omitempty"`
	EndTime        *time.Time `query:"endTime" validate:"omitempty"`
}

// EventsBusinessRuleResponse represents a business rule event response.
type EventsBusinessRuleResponse struct {
	Created                 time.Time `json:"created"`
	EventTrackerId          string    `json:"eventTrackerId,omitempty"`
	ThreadId                string    `json:"threadId"`
	OrgId                   string    `json:"orgId"`
	PathKey                 string    `json:"pathKey,omitempty"`
	RuleId                  string    `json:"ruleId"`
	BusinessRuleId          string    `json:"businessRuleId"`
	BusinessRuleName        string    `json:"businessRuleName"`
	BusinessRuleDescription string    `json:"businessRuleDescription,omitempty"`

	// Execution result
	Matched    bool  `json:"matched"`
	DurationMs int64 `json:"durationMs"`

	// Evaluation metrics
	ConditionsEvaluated int `json:"conditionsEvaluated"`
	ConditionsMatched   int `json:"conditionsMatched"`
	GroupsEvaluated     int `json:"groupsEvaluated"`
	MaxDepthReached     int `json:"maxDepthReached"`

	// State data (JSON strings)
	FinalState   string `json:"finalState,omitempty"`
	StateChanges string `json:"stateChanges,omitempty"`

	// Detailed logs (JSON strings)
	EvaluationTree string `json:"evaluationTree,omitempty"`
	ConditionLogs  string `json:"conditionLogs,omitempty"`

	// Actions (JSON string)
	ActionsToDispatch string `json:"actionsToDispatch,omitempty"`

	RetentionDays uint16 `json:"retentionDays"`
}

// EventsBusinessRuleCursorResult represents the cursor-paginated response for business rule events.
type EventsBusinessRuleCursorResult struct {
	Items       []EventsBusinessRuleResponse `json:"items"`
	NextCursor  *time.Time                   `json:"nextCursor,omitempty"`
	PrevCursor  *time.Time                   `json:"prevCursor,omitempty"`
	HasNext     bool                         `json:"hasNext"`
	HasPrevious bool                         `json:"hasPrevious"`
}
