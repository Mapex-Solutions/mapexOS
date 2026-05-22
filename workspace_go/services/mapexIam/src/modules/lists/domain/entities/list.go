package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type List struct {
	ID model.ObjectId `bson:"_id,omitempty"`

	Type    string `bson:"type"`
	Name    string `bson:"name"`
	Value   string `bson:"value"`
	Enabled bool   `bson:"enabled"`

	// Hierarchical reference - links to parent list item (e.g., manufacturer for asset types)
	ParentId *model.ObjectId `bson:"parentId,omitempty"`

	// Flexible metadata for type-specific fields
	Metadata map[string]interface{} `bson:"metadata,omitempty"`

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

func (u *List) GetCreated() time.Time { return u.Created }
func (u *List) GetUpdated() time.Time { return u.Updated }

// ListUpdateDTO is used for PATCH/UPDATE operations.
// Every field is optional (pointers), so nil means "ignore".
type ListUpdateDTO struct {
	Name     *string                `bson:"name,omitempty"`
	Value    *string                `bson:"value,omitempty"`
	Enabled  *bool                  `bson:"enabled,omitempty"`
	ParentId *model.ObjectId        `bson:"parentId,omitempty"`
	Metadata map[string]interface{} `bson:"metadata,omitempty"`

	Created *time.Time `bson:"created"`
	Updated time.Time  `bson:"updated"`
}

func (u *ListUpdateDTO) GetCreated() *time.Time { return u.Created }
func (u *ListUpdateDTO) GetUpdated() time.Time  { return u.Updated }
