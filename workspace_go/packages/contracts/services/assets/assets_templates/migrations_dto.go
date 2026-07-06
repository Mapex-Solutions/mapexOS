package assetstemplate

import "time"

// MigrationPlanCreateRequest creates a plan migrating the selected assets from
// their source template to the target template. A nil or past ScheduleAt runs
// the plan immediately.
type MigrationPlanCreateRequest struct {
	Name           string     `json:"name" validate:"required"`
	Description    string     `json:"description,omitempty"`
	FromTemplateID string     `json:"fromTemplateId" validate:"required,mongoid"`
	ToTemplateID   string     `json:"toTemplateId" validate:"required,mongoid"`
	AssetIDs       []string   `json:"assetIds" validate:"required,min=1,dive,mongoid"`
	ScheduleAt     *time.Time `json:"scheduleAt,omitempty"`
	BatchSize      *int       `json:"batchSize,omitempty" validate:"omitempty,min=1"`
}

// MigrationPlanUpdateRequest is the partial-update body for a plan (only
// provided fields are applied).
type MigrationPlanUpdateRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	ScheduleAt  *time.Time `json:"scheduleAt,omitempty"`
}

// MigrationPlanResponse is the read model for a plan in list/detail responses
// (ids as strings so the shared mapper populates it field-for-field).
type MigrationPlanResponse struct {
	ID               string    `json:"id"`
	OrgID            string    `json:"orgId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	FromTemplateID   string    `json:"fromTemplateId"`
	FromTemplateName string    `json:"fromTemplateName"`
	ToTemplateID     string    `json:"toTemplateId"`
	ToTemplateName   string    `json:"toTemplateName"`
	ScheduleAt       time.Time `json:"scheduleAt"`
	BatchSize      int       `json:"batchSize"`
	Status         string    `json:"status"`
	Total          int       `json:"total"`
	Migrated       int       `json:"migrated"`
	Failed         int       `json:"failed"`
	ProgressPct    int       `json:"progressPct"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
}

// MigrationExecutionResponse is the read model for a single per-asset migration
// execution.
type MigrationExecutionResponse struct {
	ID       string    `json:"id"`
	PlanID   string    `json:"planId"`
	AssetID  string    `json:"assetId"`
	OrgID    string    `json:"orgId"`
	Status   string    `json:"status"`
	Error    string    `json:"error,omitempty"`
	Attempts int       `json:"attempts"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}
