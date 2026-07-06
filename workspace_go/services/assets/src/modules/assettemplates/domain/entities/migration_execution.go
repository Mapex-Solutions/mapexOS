package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// ExecStatus is the per-asset outcome of a template migration.
type ExecStatus string

const (
	ExecPending   ExecStatus = "pending"
	ExecMigrated  ExecStatus = "migrated"
	ExecFailed    ExecStatus = "failed"
	ExecCancelled ExecStatus = "cancelled"
)

// MigrationExecution is one migration plan applied to one asset. It is its own
// document so a plan over thousands of assets does not bloat the plan aggregate.
type MigrationExecution struct {
	ID       model.ObjectId `bson:"_id,omitempty"`
	PlanID   model.ObjectId `bson:"planId"`
	AssetID  model.ObjectId `bson:"assetId"`
	OrgID    model.ObjectId `bson:"orgId"`
	Status   ExecStatus     `bson:"status"`
	Error    string         `bson:"error,omitempty"`
	Attempts int            `bson:"attempts"`
	Created  time.Time      `bson:"created"`
	Updated  time.Time      `bson:"updated"`
}

func (e *MigrationExecution) GetCreated() time.Time { return e.Created }
