package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// GroupMember represents a junction entity for the many-to-many relationship
// between groups and users. This replaces the embedded Members array in Group
// for better scalability with large member counts (100K+ tenants).
type GroupMember struct {
	ID      model.ObjectId  `bson:"_id,omitempty"`
	GroupID model.ObjectId  `bson:"groupId"`
	UserID  model.ObjectId  `bson:"userId"`
	OrgID   model.ObjectId  `bson:"orgId"`   // Denormalized for efficient org-scoped queries
	PathKey string          `bson:"pathKey"` // Denormalized base36 for hierarchical range queries ($gt/$lt)
	AddedAt time.Time       `bson:"addedAt"` // When the user was added to the group
	AddedBy *model.ObjectId `bson:"addedBy,omitempty"` // Who added the user (audit trail)

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (gm *GroupMember) GetCreated() time.Time { return gm.Created }
func (gm *GroupMember) GetUpdated() time.Time { return gm.Updated }
