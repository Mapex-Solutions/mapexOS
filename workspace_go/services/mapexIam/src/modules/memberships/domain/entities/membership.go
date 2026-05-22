package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type Membership struct {
	ID           model.ObjectId  `bson:"_id,omitempty"`

	// Assignee (who receives access)
	AssigneeType string          `bson:"assigneeType"` // "user" | "group"
	AssigneeID   model.ObjectId  `bson:"assigneeId"`   // User or Group ID

	// Organization
	OrgID        *model.ObjectId `bson:"orgId,omitempty"`
	OrgPathKey   string          `bson:"orgPathKey"` // Denormalized from Organization for range queries

	// Denormalized tenant anchor (optional - only for customer organizations, not vendors)
	// Vendor organizations don't have a customerId (they ARE the top-level tenant)
	// Customer organizations have customerId pointing to themselves or parent customer
	CustomerID   *model.ObjectId `bson:"customerId,omitempty"`

	// Permissions
	RoleIds      []model.ObjectId `bson:"roleIds"` // Array of Role IDs (references to Role entities)
	Scope        string           `bson:"scope"`   // "local" | "recursive"

	// Status
	Enabled      bool            `bson:"enabled"`
	Created      time.Time       `bson:"created"`
	Updated      time.Time       `bson:"updated"`
}

func (m *Membership) GetCreated() time.Time { return m.Created }
func (m *Membership) GetUpdated() time.Time { return m.Updated }
