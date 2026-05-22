package retention

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// RetentionPolicyUpsert is the body DTO for PUT (upsert) operations.
// The orgId and pathKey are populated from RequestContext (coverage middleware).
type RetentionPolicyUpsert struct {
	// Type identifies the ClickHouse table (e.g., "events", "eventsRaw")
	Type string `json:"type" validate:"required,min=1"`

	// Name is a human-readable label for the policy
	Name string `json:"name" validate:"required,min=1"`

	// RetentionDays is the number of days to retain data
	RetentionDays uint16 `json:"retentionDays" validate:"required,min=1"`

	// Enabled controls whether this retention policy is active
	Enabled *bool `json:"enabled,omitempty"`

	// Multi-tenant fields (populated from RequestContext, not from body)
	OrgId   *model.ObjectId `json:"-"`
	PathKey *string         `json:"-"`
}

// Transform performs any necessary transformations on the DTO.
func (dto *RetentionPolicyUpsert) Transform() error {
	return nil
}

// RetentionPolicyResponse is the API response DTO for retention policies.
type RetentionPolicyResponse struct {
	ID            *string `json:"id,omitempty"`
	Name          *string `json:"name,omitempty"`
	Type          *string `json:"type,omitempty"`
	RetentionDays *uint16 `json:"retentionDays,omitempty"`
	OrgId         *string `json:"orgId,omitempty"`
	PathKey       *string `json:"pathKey,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
	Created       *string `json:"created,omitempty"`
	Updated       *string `json:"updated,omitempty"`
}

// RetentionPolicyQuery is the query DTO for listing retention policies.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
type RetentionPolicyQuery struct {
	query.BaseQueryDTO

	// Type filters by retention policy type (comma-separated for $in, e.g., "events,eventsRaw")
	Type *string `query:"type" validate:"omitempty"`
}

// RetentionPolicyParams is the route params DTO for individual retention policy operations.
type RetentionPolicyParams struct {
	RetentionPolicyId string `params:"retentionPolicyId" validate:"required,mongoid"`
}
