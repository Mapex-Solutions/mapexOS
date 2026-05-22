package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type Group struct {
	ID          model.ObjectId `bson:"_id,omitempty"`
	Name        string         `bson:"name"`
	Description *string        `bson:"description,omitempty"`
	// Members field removed - now stored in group_members junction collection
	// Use GroupMemberRepository for member operations (scalable for 100K+ tenants)
	Enabled bool `bson:"enabled"`

	// Multi-tenant hierarchical fields
	OrgID   *model.ObjectId `bson:"orgId,omitempty"`
	PathKey string          `bson:"pathKey"`
	Scope   string          `bson:"scope"` // "global" | "local" - inheritance behavior

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (g *Group) GetCreated() time.Time { return g.Created }
func (g *Group) GetUpdated() time.Time { return g.Updated }
