package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// PlanStatus is the lifecycle state of a template migration plan.
type PlanStatus string

const (
	PlanPending             PlanStatus = "pending"
	PlanScheduled           PlanStatus = "scheduled"
	PlanRunning             PlanStatus = "running"
	PlanComplete            PlanStatus = "complete"
	PlanCompletedWithErrors PlanStatus = "completed_with_errors"
	PlanCancelled           PlanStatus = "cancelled"
)

// PlanCounters is the per-plan aggregate of migration outcomes.
type PlanCounters struct {
	Total    int `bson:"total"`
	Migrated int `bson:"migrated"`
	Failed   int `bson:"failed"`
}

// MigrationPlan moves the selected assets from one template to another. AssetIDs
// is the snapshot of assets selected when the plan was created.
type MigrationPlan struct {
	ID             model.ObjectId   `bson:"_id,omitempty"`
	OrgID          model.ObjectId   `bson:"orgId"`
	Name           string           `bson:"name"`
	Description    string           `bson:"description,omitempty"`
	FromTemplateID   model.ObjectId `bson:"fromTemplateId"`
	FromTemplateName string         `bson:"fromTemplateName,omitempty"`
	ToTemplateID     model.ObjectId `bson:"toTemplateId"`
	ToTemplateName   string         `bson:"toTemplateName,omitempty"`
	AssetIDs         []model.ObjectId `bson:"assetIds"`
	ScheduleAt     time.Time        `bson:"scheduleAt"`
	BatchSize      int              `bson:"batchSize"`
	Status         PlanStatus       `bson:"status"`
	Counters       PlanCounters     `bson:"counters"`
	Created        time.Time        `bson:"created"`
	Updated        time.Time        `bson:"updated"`
}

func (p *MigrationPlan) GetCreated() time.Time { return p.Created }

// IsEditable reports whether the plan can still be modified.
func (p *MigrationPlan) IsEditable() bool {
	return p.Status == PlanPending || p.Status == PlanScheduled
}

// IsTerminal reports whether the plan has reached a final state.
func (p *MigrationPlan) IsTerminal() bool {
	return p.Status == PlanComplete || p.Status == PlanCompletedWithErrors || p.Status == PlanCancelled
}

// CanTransition reports whether a plan may move from one status to another.
func CanTransition(from, to PlanStatus) bool {
	switch from {
	case PlanPending:
		return to == PlanRunning || to == PlanCancelled
	case PlanScheduled:
		return to == PlanRunning || to == PlanCancelled
	case PlanRunning:
		return to == PlanComplete || to == PlanCompletedWithErrors
	default:
		return false
	}
}
