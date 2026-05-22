package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type LakeHouseData struct {
	LakeHouseId string                 `bson:"lakeHouseId"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty"`
}

type NotificationData struct {
	NotificationId string                 `bson:"notificationId"`
	Metadata       map[string]interface{} `bson:"metadata,omitempty"`
}

type TriggerData struct {
	TriggerId string                 `bson:"triggerId"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty"`
}

type SaveEventData struct {
	Metadata map[string]interface{} `bson:"metadata,omitempty"`
}

type WorkflowData struct {
	Mode     string                 `bson:"mode"`
	Data     map[string]interface{} `bson:"data"`
	Metadata map[string]interface{} `bson:"metadata,omitempty"`
}

// MatchRule defines a single conditional rule for event routing
type MatchRule struct {
	Field    string      `bson:"field"`    // JSON path to the field (e.g., "payload.temperature")
	Operator string      `bson:"operator"` // Comparison operator (eq, neq, gt, gte, lt, lte, in, nin)
	Value    interface{} `bson:"value"`    // Value to compare against
}

// MatchConfig defines the matching strategy for conditional routing
type MatchConfig struct {
	Policy string      `bson:"policy"` // "all" = AND logic, "any" = OR logic
	Rules  []MatchRule `bson:"rules"`  // Array of matching rules
}

type Router struct {
	Kind         string            `bson:"kind"`
	Match        *MatchConfig      `bson:"match,omitempty"`
	LakeHouse    *LakeHouseData    `bson:"lakeHouse,omitempty"`
	Notification *NotificationData `bson:"notification,omitempty"`
	Trigger      *TriggerData      `bson:"trigger,omitempty"`
	SaveEvent    *SaveEventData    `bson:"saveEvent,omitempty"`
	Workflow     *WorkflowData     `bson:"workflow,omitempty"`
}

type RouteGroup struct {
	ID model.ObjectId `bson:"_id,omitempty"`

	Version     string   `bson:"version"`
	Name        string   `bson:"name"`
	Description string   `bson:"description"`
	Enabled     bool     `bson:"enabled"`
	IsSystem    bool     `bson:"isSystem"`   // true = visible to everyone (MAPEX global templates)
	IsTemplate  bool     `bson:"isTemplate"` // true = shared template (vendor/customer only)

	// Multi-tenant fields
	OrgId   *model.ObjectId `bson:"orgId,omitempty"` // null for system, org for template/local
	PathKey string          `bson:"pathKey"`         // null for system, pathKey for template/local

	Routers []Router `bson:"routers"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (rg *RouteGroup) GetCreated() time.Time { return rg.Created }
func (rg *RouteGroup) GetUpdated() time.Time { return rg.Updated }

// RouteGroupUpdateDTO is used for PATCH/UPDATE operations.
// Every field is optional (pointers), so nil means "ignore".
type RouteGroupUpdateDTO struct {
	Version     *string   `bson:"version,omitempty"`
	Name        *string   `bson:"name,omitempty"`
	Description *string   `bson:"description,omitempty"`
	Enabled     *bool     `bson:"enabled,omitempty"`

	// Multi-tenant fields (usually not updated, but available)
	OrgId   *model.ObjectId `bson:"orgId,omitempty"`
	PathKey *string         `bson:"pathKey,omitempty"`

	Routers *[]Router `bson:"routers,omitempty"`

	Created *time.Time `bson:"created"`
	Updated time.Time  `bson:"updated"`
}

func (rg *RouteGroupUpdateDTO) GetCreated() *time.Time { return rg.Created }
func (rg *RouteGroupUpdateDTO) GetUpdated() time.Time  { return rg.Updated }
