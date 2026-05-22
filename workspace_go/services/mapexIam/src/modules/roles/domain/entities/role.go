package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type Role struct {
	ID          model.ObjectId `bson:"_id,omitempty"`
	Name        string         `bson:"name"`
	Description *string        `bson:"description,omitempty"`
	Permissions []string       `bson:"permissions"` // e.g., ["read_users", "write_devices", "read_logs"]

	// Template visibility flags
	IsSystem   bool `bson:"isSystem"`   // true = visible to everyone (MAPEX global templates)
	IsTemplate bool `bson:"isTemplate"` // true = shared template (vendor/customer only)

	// Multi-tenant hierarchical fields
	OrgID   *model.ObjectId `bson:"orgId,omitempty"` // null for system, org for template/local
	PathKey string          `bson:"pathKey"`         // null for system, pathKey for template/local
	Scope   string          `bson:"scope"`           // "global" | "local" - inheritance behavior

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (r *Role) GetCreated() time.Time { return r.Created }
func (r *Role) GetUpdated() time.Time { return r.Updated }
