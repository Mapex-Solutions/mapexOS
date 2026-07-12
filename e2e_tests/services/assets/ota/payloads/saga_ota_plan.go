package payloads

import (
	"fmt"
	"time"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
)

// planRatePerMinute paces the rollout high enough that the reconciler dispatches
// the whole (single-asset) plan on its first scan — the e2e wants determinism,
// not throttling.
const planRatePerMinute = 1000

// planMaxTime bounds the plan; short so a straggler / negative path closes
// without a long wait. The happy path finishes (early-close) well before it.
const planMaxTime = 2 * time.Minute

// OTAPlanBuilder wraps the OTAPlanCreateRequest contract with the determinism
// defaults saga journeys need.
type OTAPlanBuilder struct {
	spec otaDtos.OTAPlanCreateRequest
}

// NewOTAPlan builds a plan-create request: migrate the given assets from the
// source template to the firmware's target template, starting now with a high
// rate and a short max time. abort is disabled so a single expected failure in a
// negative test never trips it accidentally.
func NewOTAPlan(runID, firmwareID, sourceTemplateID string, assetIDs []string) *OTAPlanBuilder {
	now := time.Now().UTC()
	return &OTAPlanBuilder{
		spec: otaDtos.OTAPlanCreateRequest{
			FirmwareID:       firmwareID,
			SourceTemplateID: sourceTemplateID,
			AssetIDs:         assetIDs,
			Name:             fmt.Sprintf("ota-%s", runID),
			StartAt:          now,
			MaxTime:          now.Add(planMaxTime),
			RolloutConfig: otaDtos.OTARolloutConfigDTO{
				RatePerMinute:     planRatePerMinute,
				AbortThresholdPct: 0,
				AbortMinExecuted:  0,
			},
		},
	}
}

// WithMaxTime overrides the plan's max time (StartAt + d). Used by the negative
// path (silent device → TIMED_OUT) follow-up.
func (b *OTAPlanBuilder) WithMaxTime(d time.Duration) *OTAPlanBuilder {
	b.spec.MaxTime = b.spec.StartAt.Add(d)
	return b
}

// Build returns the contracts payload ready for POST /api/v1/ota/plans.
func (b *OTAPlanBuilder) Build() otaDtos.OTAPlanCreateRequest {
	return b.spec
}
