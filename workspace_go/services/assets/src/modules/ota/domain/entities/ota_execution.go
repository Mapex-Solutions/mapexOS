package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// ExecutionState is the device-reported per-asset OTA execution state. The Asset
// MS owns QUEUED/INITIATED; the device reports everything from DOWNLOADING on.
type ExecutionState string

const (
	ExecQueued      ExecutionState = "QUEUED"
	ExecInitiated   ExecutionState = "INITIATED"
	ExecDownloading ExecutionState = "DOWNLOADING"
	ExecDownloaded  ExecutionState = "DOWNLOADED"
	ExecVerified    ExecutionState = "VERIFIED"
	ExecUpdating    ExecutionState = "UPDATING"
	ExecUpdated     ExecutionState = "UPDATED"
	ExecFailed      ExecutionState = "FAILED"
	ExecTimedOut    ExecutionState = "TIMED_OUT"
)

// OTAExecution is one plan instance on one asset. It is its own document so a
// plan over thousands of assets does not bloat the plan aggregate, and each
// device's status update writes only its own small document.
type OTAExecution struct {
	ID         model.ObjectId `bson:"_id,omitempty"`
	PlanID     model.ObjectId `bson:"planId"`
	AssetID    model.ObjectId `bson:"assetId"`
	OrgID      model.ObjectId `bson:"orgId"`
	AssetUUID  string         `bson:"assetUUID,omitempty"`
	Protocol   string         `bson:"protocol"`
	State      ExecutionState `bson:"state"`
	Percentage int32          `bson:"percentage"`
	Error      string         `bson:"error,omitempty"`
	Attempts   int            `bson:"attempts"`
	QueuedAt   time.Time      `bson:"queuedAt"`
	StartedAt  time.Time      `bson:"startedAt,omitempty"`
	Updated    time.Time      `bson:"updated"`
}

func (e *OTAExecution) GetCreated() time.Time { return e.QueuedAt }

// IsTerminal reports whether the execution has reached a final state.
func (e *OTAExecution) IsTerminal() bool {
	return e.State == ExecUpdated || e.State == ExecFailed || e.State == ExecTimedOut
}
