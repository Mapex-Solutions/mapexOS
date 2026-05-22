package dtos

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * Event Store DTOs
 * Shared between Router service (publisher) and Events service (consumer)
 * For processed events with EVA (Entity-Value-Attribute) field mapping
 */

// EventStoreDTO is the payload published to NATS for event storage.
// Any service can publish: Router (assets), Workflow, etc.
// The publisher defines the source and provides the template context for EVA resolution.
type EventStoreDTO struct {
	// Created is the event occurrence time
	Created time.Time `json:"created"`

	// EventTrackerId is the UUID for end-to-end event tracking across services
	EventTrackerId string `json:"eventTrackerId"`

	// ThreadId groups related events for correlation.
	// For router/assets: assetId. For rule engine: businessRuleId.
	ThreadId string `json:"threadId"`

	// Source identifies who published the event (e.g., "asset", "rule")
	// Used by the detail endpoint to determine where to fetch DynamicFields for EVA resolution.
	Source string `json:"source" validate:"required"`

	// AssetId is the MongoDB ObjectId of the asset
	AssetId string `json:"assetId" validate:"required"`

	// AssetUUID is the device/gateway identifier (devEUI, deviceId, etc.)
	AssetUUID string `json:"assetUUID"`

	// AssetTemplateId is the template ID for EVA field resolution (AssetTemplate, BusinessRule, etc.)
	// Used to fetch DynamicFields from cache for fieldId mapping.
	AssetTemplateId string `json:"assetTemplateId"`

	// AssetTemplateOrgId is the org that owns the template.
	// "mapexos_public" for system templates, or the org's ID for private templates.
	// Used as cache key prefix: {templateOrgId}/{templateId}
	AssetTemplateOrgId string `json:"assetTemplateOrgId"`

	// Human-readable names (denormalized at write-time by Router)
	AssetName           string `json:"assetName,omitempty"`
	AssetDescription    string `json:"assetDescription,omitempty"`
	TemplateName        string `json:"templateName,omitempty"`
	TemplateDescription string `json:"templateDescription,omitempty"`

	// OrgId is the organization ID for multi-tenancy
	OrgId string `json:"orgId" validate:"required"`

	// PathKey is the asset's hierarchical path for organizational queries
	PathKey string `json:"pathKey"`

	// Event is the actual event data containing dynamic fields.
	// Fields are mapped to EVA columns based on DynamicFields from the template.
	Event map[string]interface{} `json:"event"`

	// Metadata contains additional context from the publisher
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EvaFilterOperator defines valid comparison operators for EVA field filtering.
// These are safe string values sent from the frontend and mapped to ClickHouse operators.
type EvaFilterOperator string

const (
	EvaOpEqual        EvaFilterOperator = "eq"
	EvaOpNotEqual     EvaFilterOperator = "neq"
	EvaOpGreater      EvaFilterOperator = "gt"
	EvaOpGreaterEqual EvaFilterOperator = "gte"
	EvaOpLess         EvaFilterOperator = "lt"
	EvaOpLessEqual    EvaFilterOperator = "lte"
	EvaOpBetween      EvaFilterOperator = "between"
	EvaOpLike         EvaFilterOperator = "like"
)

// EvaFilter represents a single filter condition on an EVA MAP column field.
// The frontend resolves the fieldId from the AssetTemplate.DynamicFields array
// and sends it along with the bucket type, operator, and value(s).
//
// Example ClickHouse SQL generated: eva_number[42] >= 10.5
type EvaFilter struct {
	// FieldId is the uint16 key in the EVA MAP column (from DynamicField.FieldId property)
	// Auto-increment starting at 1, so min=1 is correct
	FieldId uint16 `json:"fieldId" validate:"min=1,max=65535"`

	// Bucket is the EVA type bucket: "number", "string", "bool", "date"
	Bucket string `json:"bucket" validate:"required,oneof=number string bool date"`

	// Operator is the comparison operator (eq, neq, gt, gte, lt, lte, between, like)
	Operator EvaFilterOperator `json:"operator" validate:"required,oneof=eq neq gt gte lt lte between like"`

	// Value is the primary comparison value (sent as string, cast per bucket type)
	Value string `json:"value" validate:"required"`

	// EndValue is the range end value, used only with "between" operator
	EndValue string `json:"endValue,omitempty"`
}

// EventsStoreQuery represents query parameters for listing processed events.
// Uses cursor-based pagination for efficient querying of large datasets.
// Supports both indexed column filters and EVA MAP column filters.
type EventsStoreQuery struct {
	query.CursorQueryDTO

	// Indexed column filters
	EventTrackerId  *string    `query:"eventTrackerId" validate:"omitempty"`
	ThreadId        *string    `query:"threadId" validate:"omitempty"`
	AssetId         *string    `query:"assetId" validate:"omitempty"`
	AssetTemplateId *string    `query:"assetTemplateId" validate:"omitempty"`
	EventType       *string    `query:"eventType" validate:"omitempty"`
	Source          *string    `query:"source" validate:"omitempty"`
	StartTime       *time.Time `query:"startTime" validate:"omitempty"`
	EndTime         *time.Time `query:"endTime" validate:"omitempty"`

	// EVA dynamic field filters (sent via POST body as JSON array)
	EvaFilters []EvaFilter `query:"evaFilters" validate:"omitempty,dive"`
}

// EventsStoreResponse represents a processed event in the list response.
// EVA fields are NOT included in the list — they are resolved in the detail view
// where the frontend fetches the template to translate fieldIds to field names.
type EventsStoreResponse struct {
	Created             time.Time `json:"created"`
	EventTrackerId      string    `json:"eventTrackerId,omitempty"`
	ThreadId            string    `json:"threadId,omitempty"`
	AssetId             string    `json:"assetId"`
	AssetName           string    `json:"assetName,omitempty"`
	AssetDescription    string    `json:"assetDescription,omitempty"`
	TemplateName        string    `json:"templateName,omitempty"`
	TemplateDescription string    `json:"templateDescription,omitempty"`
	OrgId               string    `json:"orgId"`
	PathKey             string    `json:"pathKey,omitempty"`
	Source              string    `json:"source,omitempty"`
	Payload             string    `json:"payload"`
	Metadata            string    `json:"metadata,omitempty"`
}

// EventsStoreDetailResponse represents the detail view of a single processed event.
// Includes resolved EVA field names as advancedSearch map (fieldId→fieldName resolved at read time
// from the appropriate template based on Source: "asset"→AssetTemplate, "rule"→BusinessRule).
type EventsStoreDetailResponse struct {
	Created             time.Time              `json:"created"`
	EventTrackerId      string                 `json:"eventTrackerId,omitempty"`
	Source              string                 `json:"source"`
	AssetId             string                 `json:"assetId"`
	AssetName           string                 `json:"assetName,omitempty"`
	AssetDescription    string                 `json:"assetDescription,omitempty"`
	TemplateName        string                 `json:"templateName,omitempty"`
	TemplateDescription string                 `json:"templateDescription,omitempty"`
	OrgId               string                 `json:"orgId"`
	PathKey             string                 `json:"pathKey,omitempty"`
	Payload             string                 `json:"payload"`
	Metadata            string                 `json:"metadata,omitempty"`
	AdvancedSearch      map[string]interface{} `json:"advancedSearch,omitempty"`
}

// EventsStoreCursorResult represents the cursor-paginated response for processed events.
type EventsStoreCursorResult struct {
	Items       []EventsStoreResponse `json:"items"`
	NextCursor  *time.Time            `json:"nextCursor,omitempty"`
	PrevCursor  *time.Time            `json:"prevCursor,omitempty"`
	HasNext     bool                  `json:"hasNext"`
	HasPrevious bool                  `json:"hasPrevious"`
}
