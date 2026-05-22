package types

import "time"

/**
 * CACHE INVALIDATION EVENTS
 * Types for FANOUT cache invalidation messages.
 */

// AssetInvalidateEvent is the payload for FANOUT cache invalidation.
// Key format: {orgId}/{assetUUID}
type AssetInvalidateEvent struct {
	OrgId     string `json:"orgId"`
	AssetUUID string `json:"assetUUID"`
}

// TemplateInvalidateEvent is the payload for FANOUT template cache invalidation.
// Key format: {orgId}/{templateId}
type TemplateInvalidateEvent struct {
	OrgId      string `json:"orgId"`
	TemplateId string `json:"templateId"`
}

/**
 * ROUTING HISTORY EVENTS
 * Types for routing execution history published to Events Service.
 */

// RouterHistoryEvent is published to NATS for routing execution history.
type RouterHistoryEvent struct {
	Created        time.Time            `json:"created"`
	EventTrackerId string               `json:"eventTrackerId"` // UUID for end-to-end event tracking across services
	ThreadId       string               `json:"threadId"`
	OrgId          string               `json:"orgId"`
	PathKey        string               `json:"pathKey"`
	AssetId        string               `json:"assetId"`
	RouterId       string               `json:"routerId"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Routers        []RouterResultRecord `json:"routers"`
}

// RouterResultRecord holds a single router execution result.
type RouterResultRecord struct {
	Kind       string                  `json:"kind"`
	Matched    bool                    `json:"matched"`
	Published  bool                    `json:"published"`
	Conditions []ConditionResultRecord `json:"conditions"`
}

// ConditionResultRecord holds a single condition evaluation result.
type ConditionResultRecord struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
	Passed   bool        `json:"passed"`
}
