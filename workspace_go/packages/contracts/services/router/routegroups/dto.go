package routegroups

import (
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type RouteGroupId struct {
	RouteGroupId string `params:"routeGroupId" validate:"required,mongoid"`
}

type LakeHouseData struct {
	LakeHouseId string                 `json:"lakeHouseId" validate:"required,mongoid"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type NotificationData struct {
	NotificationId string                 `json:"notificationId" validate:"required,mongoid"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type TriggerData struct {
	TriggerId string                 `json:"triggerId" validate:"required,mongoid"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type SaveEventData struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// WorkflowData defines how events are delivered to the Workflow Service.
// Router publishes to WORKFLOW-EXECUTION stream with a fixed subject.
// Routing is done by the "mode" field in the payload — never by subject.
//
// Three modes:
//
//   - "newInstance": Creates a new execution from an instance config.
//     data: { instanceId (required), workflowUUID (optional — if omitted, runtime generates GUIDv4) }
//
//   - "signal": Delivers a signal to a waiting execution.
//     data: { workflowUUID (required), signalName (required), signalData (optional) }
//
//   - "signalOrStart": Tries signal first, falls back to newInstance.
//     data: { instanceId, workflowUUID, signalName, signalData }
type WorkflowData struct {
	// Delivery mode: "newInstance", "signal", "signalOrStart"
	Mode string `json:"mode" validate:"required,oneof=newInstance signal signalOrStart"`

	// Mode-specific configuration. Schema defined per mode.
	Data map[string]interface{} `json:"data" validate:"required"`

	// Additional metadata merged into the workflow execution message.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Transform validates WorkflowData based on the selected mode.
func (w *WorkflowData) Transform() error {
	if w.Data == nil {
		return fmt.Errorf("field 'data' is required")
	}

	switch w.Mode {
	case "newInstance":
		if _, ok := w.Data["instanceId"].(string); !ok || w.Data["instanceId"] == "" {
			return fmt.Errorf("data.instanceId is required when mode is 'newInstance'")
		}

	case "signal":
		if _, ok := w.Data["workflowUUID"].(string); !ok || w.Data["workflowUUID"] == "" {
			return fmt.Errorf("data.workflowUUID is required when mode is 'signal'")
		}
		if _, ok := w.Data["signalName"].(string); !ok || w.Data["signalName"] == "" {
			return fmt.Errorf("data.signalName is required when mode is 'signal'")
		}

	case "signalOrStart":
		if _, ok := w.Data["instanceId"].(string); !ok || w.Data["instanceId"] == "" {
			return fmt.Errorf("data.instanceId is required when mode is 'signalOrStart'")
		}
		if _, ok := w.Data["workflowUUID"].(string); !ok || w.Data["workflowUUID"] == "" {
			return fmt.Errorf("data.workflowUUID is required when mode is 'signalOrStart'")
		}
		if _, ok := w.Data["signalName"].(string); !ok || w.Data["signalName"] == "" {
			return fmt.Errorf("data.signalName is required when mode is 'signalOrStart'")
		}
	}

	return nil
}

// MatchRule defines a single conditional rule for event routing
// Supports multiple operators for flexible event filtering
type MatchRule struct {
	Field    string      `json:"field" validate:"required,min=1"`                                // JSON path to the field (e.g., "payload.temperature", "metadata.deviceType")
	Operator string      `json:"operator" validate:"required,oneof=eq neq gt gte lt lte in nin"` // Comparison operator
	Value    interface{} `json:"value" validate:"required"`                                      // Value to compare against
}

// MatchConfig defines the matching strategy for conditional routing
// Allows AND/OR logic for multiple rules
type MatchConfig struct {
	Policy string       `json:"policy" validate:"required,oneof=all any"` // "all" = AND logic (all rules must pass), "any" = OR logic (at least one rule must pass)
	Rules  *[]MatchRule `json:"rules" validate:"required,min=1,dive"`     // Array of matching rules
}

type Router struct {
	Kind         string            `json:"kind" validate:"required,oneof=lake_house notification trigger save_event workflow"`
	Match        *MatchConfig      `json:"match,omitempty"` // Optional conditional routing rules
	LakeHouse    *LakeHouseData    `json:"lakeHouse,omitempty"`
	Notification *NotificationData `json:"notification,omitempty"`
	Trigger      *TriggerData      `json:"trigger,omitempty"`
	SaveEvent    *SaveEventData    `json:"saveEvent,omitempty"`
	Workflow     *WorkflowData     `json:"workflow,omitempty"`
}

type RouteGroupCreate struct {
	Version     string    `json:"version" validate:"required,semver"`
	Name        string    `json:"name" validate:"required,min=1"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1"`
	Enabled     *bool     `json:"enabled,omitempty"`
	IsSystem    *bool     `json:"isSystem,omitempty"`
	IsTemplate  *bool     `json:"isTemplate,omitempty"`
	Routers     *[]Router `json:"routers,omitempty" validate:"omitempty,min=1"`

	// Multi-tenant fields (populated automatically by coverage middleware)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type RouteGroupUpdate struct {
	Version     *string   `json:"version,omitempty" validate:"omitempty,semver"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=1"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1"`
	Enabled     *bool     `json:"enabled,omitempty"`
	OrgId       *string   `json:"orgId,omitempty" validate:"omitempty,min=3"`
	Routers     *[]Router `json:"routers,omitempty" validate:"omitempty,min=1,dive"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type RouteGroupResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Version     *string          `json:"version,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	IsSystem    *bool            `json:"isSystem,omitempty"`
	IsTemplate  *bool            `json:"isTemplate,omitempty"`

	// Multi-tenant fields
	OrgId   *model.ObjectId `json:"orgId,omitempty"`
	PathKey *string         `json:"pathKey,omitempty"`

	Routers *[]Router `json:"routers,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

func (r *RouteGroupResponse) SetCreated(t *common.NullTime) { r.Created = t }
func (r *RouteGroupResponse) SetUpdated(t *common.NullTime) { r.Updated = t }

// RouteGroupQuery represents query parameters for listing route groups.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
//
// Standard fields (from BaseQueryDTO):
//   - Projection: comma-separated fields to return
//   - Page: page number (default: 1)
//   - PerPage: items per page (default: 20)
//   - Sort: sort order (default: "created:desc")
//   - IncludeChildren: include child orgs hierarchically (default: false)
//
// Module-specific filters:
//   - Name: filter by route group name (partial match)
//   - Enabled: filter by enabled status (true/false)
//   - Version: filter by version (exact match)
//   - Kinds: filter by router kinds — returns only groups whose every router.kind is in this set (strict)
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type RouteGroupQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name       *string `query:"name" validate:"omitempty,max=100"`
	Enabled    *bool   `query:"enabled" validate:"omitempty"`
	IsTemplate *bool   `query:"isTemplate" validate:"omitempty"`
	IsSystem   *bool   `query:"isSystem" validate:"omitempty"`
	Version    *string `query:"version" validate:"omitempty,semver"`

	// Kinds optionally restricts the result set to RouteGroups whose every
	// router kind is contained in the provided set. Strict semantic: a group
	// with mixed kinds (e.g., [trigger, save_event]) is excluded when kinds=
	// [trigger, workflow]. Used by the asset wizard's Health step to surface
	// only RouteGroups acceptable to validateHealthMonitorConfig.
	// Wire format: repeated query key — `?kinds=trigger&kinds=workflow`.
	Kinds []string `query:"kinds" validate:"omitempty,dive,oneof=lake_house notification trigger save_event workflow"`
}

/** INTERNAL API DTOs - For inter-service communication via API Key **/

// RouteGroupInternalIdsQuery is used for internal API calls to fetch multiple route groups by IDs.
// IDs are comma-separated in the query string (e.g., "id1,id2,id3").
// Supports projection to optimize fields returned.
type RouteGroupInternalIdsQuery struct {
	Ids        string  `query:"ids" validate:"required"`              // Comma-separated RouteGroup IDs (e.g., "id1,id2,id3")
	Projection *string `query:"projection" validate:"omitempty"` // Optional: comma-separated fields to return (e.g., "name,enabled,version")
}

/** TRANSFORMATIONS **/

func (m *MatchRule) Transform() error {
	// Validate operator-specific value types
	switch m.Operator {
	case "eq", "neq", "gt", "gte", "lt", "lte":
		// These operators require non-nil scalar values
		if m.Value == nil {
			return fmt.Errorf("operator '%s' requires a non-nil value", m.Operator)
		}

	case "in", "nin":
		// These operators require array values
		if m.Value == nil {
			return fmt.Errorf("operator '%s' requires an array value", m.Operator)
		}
		// Check if value is a slice/array
		switch v := m.Value.(type) {
		case []interface{}, []string, []int, []float64, []bool:
			// Valid array types
		default:
			return fmt.Errorf("operator '%s' requires an array value, got %T", m.Operator, v)
		}
	}

	return nil
}

func (m *MatchConfig) Transform() error {
	// Validate that rules array is not empty
	if m.Rules == nil || len(*m.Rules) == 0 {
		return fmt.Errorf("match rules cannot be empty")
	}

	// Validate each rule
	for i, rule := range *m.Rules {
		if err := rule.Transform(); err != nil {
			return fmt.Errorf("invalid match rule at index %d: %w", i, err)
		}
	}

	return nil
}

func (r *Router) Transform() error {
	// Validate match config if provided
	if r.Match != nil {
		if err := r.Match.Transform(); err != nil {
			return fmt.Errorf("invalid match configuration: %w", err)
		}
	}

	// Validate router kind configuration
	switch r.Kind {
	case "save_event":
		return nil


	case "lake_house":
		if r.LakeHouse == nil {
			return fmt.Errorf("field 'lakeHouse' must be provided when kind is 'lake_house'")
		}

	case "notification":
		if r.Notification == nil {
			return fmt.Errorf("field 'notification' must be provided when kind is 'notification'")
		}

	case "trigger":
		if r.Trigger == nil {
			return fmt.Errorf("field 'trigger' must be provided when kind is 'trigger'")
		}

	case "workflow":
		if r.Workflow == nil {
			return fmt.Errorf("field 'workflow' must be provided when kind is 'workflow'")
		}
		if err := r.Workflow.Transform(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid router kind: %s", r.Kind)
	}

	return nil
}

/** DOCUMENTATION & EXAMPLES **/

// Conditional Routing Examples:
//
// Example 1: Route events to LakeHouse ONLY if temperature >= 25 AND deviceType = "sensor"
// {
//   "kind": "lake_house",
//   "match": {
//     "policy": "all",
//     "rules": [
//       {
//         "field": "payload.temperature",
//         "operator": "gte",
//         "value": 25
//       },
//       {
//         "field": "metadata.deviceType",
//         "operator": "eq",
//         "value": "sensor"
//       }
//     ]
//   },
//   "lakeHouse": {
//     "lakeHouseId": "507f1f77bcf86cd799439011",
//     "metadata": {}
//   }
// }
//
// Example 2: Route events to Notification if status is "critical" OR "error"
// {
//   "kind": "notification",
//   "match": {
//     "policy": "any",
//     "rules": [
//       {
//         "field": "payload.status",
//         "operator": "eq",
//         "value": "critical"
//       },
//       {
//         "field": "payload.status",
//         "operator": "eq",
//         "value": "error"
//       }
//     ]
//   },
//   "notification": {
//     "notificationId": "507f1f77bcf86cd799439012",
//     "metadata": {}
//   }
// }
//
// Example 3: Route events if deviceId is in a specific list
// {
//   "match": {
//     "policy": "all",
//     "rules": [
//       {
//         "field": "metadata.deviceId",
//         "operator": "in",
//         "value": ["dev-001", "dev-002", "dev-003"]
//       }
//     ]
//   },
//     "businessRuleId": "507f1f77bcf86cd799439013",
//     "metadata": {}
//   }
// }
//
// Example 4: Complex AND logic - Route if temp > 30 AND humidity < 40 AND deviceType in allowed list
// {
//   "kind": "save_event",
//   "match": {
//     "policy": "all",
//     "rules": [
//       {
//         "field": "payload.temperature",
//         "operator": "gt",
//         "value": 30
//       },
//       {
//         "field": "payload.humidity",
//         "operator": "lt",
//         "value": 40
//       },
//       {
//         "field": "metadata.deviceType",
//         "operator": "in",
//         "value": ["sensor", "gateway", "actuator"]
//       }
//     ]
//   }
// }
//
// Supported Operators:
//   - eq: Equal (payload.status == "active")
//   - neq: Not Equal (payload.status != "inactive")
//   - gt: Greater Than (payload.temperature > 25)
//   - gte: Greater Than or Equal (payload.temperature >= 25)
//   - lt: Less Than (payload.humidity < 60)
//   - lte: Less Than or Equal (payload.humidity <= 60)
//   - in: Value in Array (payload.status in ["active", "pending"])
//   - nin: Value NOT in Array (payload.status not in ["disabled", "error"])
//
// Policy Types:
//   - all: AND logic - ALL rules must pass for routing to occur
//   - any: OR logic - AT LEAST ONE rule must pass for routing to occur
