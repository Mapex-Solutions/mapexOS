package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// RetentionPolicy defines the retention configuration for one type per organization.
// Each organization has one document per retention type (8 documents per org).
// Unique constraint: {orgId, type} compound index.
type RetentionPolicy struct {
	ID model.ObjectId `bson:"_id,omitempty"`

	// Human-readable name for the policy (e.g., "Events Raw")
	Name string `bson:"name"`

	// Type identifies the ClickHouse table (e.g., "events", "eventsRaw")
	Type string `bson:"type"`

	// RetentionDays is the number of days to retain data
	RetentionDays uint16 `bson:"retentionDays"`

	// Multi-tenant fields
	OrgId   *model.ObjectId `bson:"orgId,omitempty"`
	PathKey string          `bson:"pathKey"`

	// Enabled controls whether this retention policy is active
	Enabled bool `bson:"enabled"`

	// Timestamps
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (rp *RetentionPolicy) GetCreated() time.Time { return rp.Created }
func (rp *RetentionPolicy) GetUpdated() time.Time { return rp.Updated }
