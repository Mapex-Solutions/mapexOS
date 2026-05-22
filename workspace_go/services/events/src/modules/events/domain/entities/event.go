package entities

import (
	"time"
)

/**
 * Event represents an event stored in ClickHouse for analytics and history.
 * ClickHouse is optimized for time-series data and OLAP queries.
 *
 * EVA (Entity-Value-Attribute) Storage:
 * Uses MAP<UInt16, Type> instead of Array(Tuple(String, Type)) for 3-4x faster queries.
 * FieldId (uint16) comes from AssetTemplate.DynamicFields.
 *
 * Example:
 *   EvaNumber = map[uint16]float64{1: 25.5, 2: 60.0}  // 1=temperature, 2=humidity
 *   EvaString = map[uint16]string{3: "ok"}             // 3=status
 */
type Event struct {
	// Created is the event occurrence time (indexed for efficient queries)
	Created time.Time `ch:"created"`

	// EventTrackerId is the UUID for end-to-end event tracking across services
	EventTrackerId string `ch:"event_tracker_id"`

	// ThreadId is the internal processing thread ID
	ThreadId string `ch:"thread_id"`

	// AssetId is the device/gateway UUID that generated the event
	AssetId string `ch:"asset_id"`

	// AssetTemplateId is the template ID for EVA field resolution (AssetTemplate, BusinessRule, etc.)
	AssetTemplateId string `ch:"asset_template_id"`

	// AssetTemplateOrgId is the template owner org for cache lookup ("mapexos_public" or orgId)
	AssetTemplateOrgId string `ch:"asset_template_org_id"`

	// Human-readable names (denormalized at write-time by Router)
	AssetName           string `ch:"asset_name"`
	AssetDescription    string `ch:"asset_description"`
	TemplateName        string `ch:"template_name"`
	TemplateDescription string `ch:"template_description"`

	// OrgId is the organization identifier for multi-tenancy
	OrgId string `ch:"org_id"`

	// PathKey is the asset's hierarchical path for organizational queries
	PathKey string `ch:"path_key"`

	// EventType categorizes the event (e.g., "telemetry", "alarm", "command")
	EventType string `ch:"event_type"`

	// Source identifies the originating service/module
	Source string `ch:"source"`

	// Payload contains the event data as JSON string (ClickHouse supports JSON queries)
	Payload string `ch:"payload"`

	// Metadata contains additional context as JSON string
	Metadata string `ch:"metadata"`

	// EVA Fields: MAP<UInt16, Type>
	// Key = fieldId from AssetTemplate.DynamicFields, Value = typed value (not string!)

	// EvaNumber stores numeric values (temperature, pressure, count, etc.)
	// Key is fieldId from AssetTemplate.DynamicFields
	EvaNumber map[uint16]float64 `ch:"eva_number"`

	// EvaString stores text values (status, location, device_type, etc.)
	// Key is fieldId from AssetTemplate.DynamicFields
	EvaString map[uint16]string `ch:"eva_string"`

	// EvaBool stores boolean values (alarm, online, active, etc.)
	// Key is fieldId from AssetTemplate.DynamicFields
	// Value: 1=true, 0=false
	EvaBool map[uint16]uint8 `ch:"eva_bool"`

	// EvaDate stores datetime values (last_update, alarm_time, etc.)
	// Key is fieldId from AssetTemplate.DynamicFields
	EvaDate map[uint16]time.Time `ch:"eva_date"`

	// RetentionDays specifies how long this event should be kept (dynamic TTL per organization)
	// Fetched from MapexOS retention policies and cached for 24 hours
	RetentionDays uint16 `ch:"retention_days"`
}

// NewEvent creates a new Event with initialized EVA maps.
func NewEvent() *Event {
	return &Event{
		EvaNumber: make(map[uint16]float64),
		EvaString: make(map[uint16]string),
		EvaBool:   make(map[uint16]uint8),
		EvaDate:   make(map[uint16]time.Time),
	}
}

// SetEvaNumber sets a numeric EVA field by fieldId.
func (e *Event) SetEvaNumber(fieldId uint16, value float64) {
	if e.EvaNumber == nil {
		e.EvaNumber = make(map[uint16]float64)
	}
	e.EvaNumber[fieldId] = value
}

// SetEvaString sets a string EVA field by fieldId.
func (e *Event) SetEvaString(fieldId uint16, value string) {
	if e.EvaString == nil {
		e.EvaString = make(map[uint16]string)
	}
	e.EvaString[fieldId] = value
}

// SetEvaBool sets a boolean EVA field by fieldId.
func (e *Event) SetEvaBool(fieldId uint16, value bool) {
	if e.EvaBool == nil {
		e.EvaBool = make(map[uint16]uint8)
	}
	if value {
		e.EvaBool[fieldId] = 1
	} else {
		e.EvaBool[fieldId] = 0
	}
}

// SetEvaDate sets a datetime EVA field by fieldId.
func (e *Event) SetEvaDate(fieldId uint16, value time.Time) {
	if e.EvaDate == nil {
		e.EvaDate = make(map[uint16]time.Time)
	}
	e.EvaDate[fieldId] = value
}
