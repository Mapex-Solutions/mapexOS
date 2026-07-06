package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// PlanStatus is the lifecycle state of an OTA plan.
type PlanStatus string

const (
	PlanDraft      PlanStatus = "DRAFT"
	PlanScheduled  PlanStatus = "SCHEDULED"
	PlanInProgress PlanStatus = "IN_PROGRESS"
	PlanCompleted  PlanStatus = "COMPLETED"
	PlanClosed     PlanStatus = "CLOSED"
	PlanCanceled   PlanStatus = "CANCELED"
	PlanAborted    PlanStatus = "ABORTED"
)

// OTARolloutConfig controls the pacing and auto-abort of a plan rollout.
type OTARolloutConfig struct {
	RatePerMinute     int `bson:"ratePerMinute"`
	AbortThresholdPct int `bson:"abortThresholdPct"`
	AbortMinExecuted  int `bson:"abortMinExecuted"`
}

// PlanCounters is the per-plan aggregate. The plan document holds only counters;
// the per-device detail lives in OTAExecution documents.
type PlanCounters struct {
	Total     int `bson:"total"`
	Succeeded int `bson:"succeeded"`
	Failed    int `bson:"failed"`
	TimedOut  int `bson:"timedOut"`
}

// OTAPlan is the orchestrating aggregate: a reconciliation engine bounded by two
// timers (StartAt, MaxTime) that migrates the selected assets from the source
// template to the target template carried by the firmware.
type OTAPlan struct {
	ID               model.ObjectId   `bson:"_id,omitempty"`
	OrgID            model.ObjectId   `bson:"orgId"`
	Name             string           `bson:"name"`
	Description      string           `bson:"description,omitempty"`
	FirmwareID       model.ObjectId   `bson:"firmwareId"`
	SourceTemplateID model.ObjectId   `bson:"sourceTemplateId"`
	TargetTemplateID model.ObjectId   `bson:"targetTemplateId"`
	StartAt          time.Time        `bson:"startAt"`
	MaxTime          time.Time        `bson:"maxTime"`
	RolloutConfig    OTARolloutConfig `bson:"rolloutConfig"`
	Status           PlanStatus       `bson:"status"`
	Counters         PlanCounters     `bson:"counters"`
	Created          time.Time        `bson:"created"`
	Updated          time.Time        `bson:"updated"`
}

func (p *OTAPlan) GetCreated() time.Time { return p.Created }
