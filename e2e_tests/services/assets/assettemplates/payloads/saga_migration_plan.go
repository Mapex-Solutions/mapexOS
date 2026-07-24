package payloads

import (
	"fmt"
	"time"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"
)

// MigrationPlanBuilder assembles a template-migration create request. It is pure —
// it depends only on the runID (for isolation) and the values passed in, never on
// the saga bag — so a step can build the body inside its Do and keep the payload
// reusable and testable in isolation.
type MigrationPlanBuilder struct {
	spec contracts.MigrationPlanCreateRequest
}

// NewMigrationPlan builds an immediate (scheduleAt nil) plan migrating the given
// assets from one template to another. The name is runID-stamped so parallel
// journeys never collide. The asset id slice is copied so the caller's slice is
// never aliased or mutated by the With* helpers.
func NewMigrationPlan(runID, fromTemplateID, toTemplateID string, assetIDs []string) *MigrationPlanBuilder {
	ids := make([]string, len(assetIDs))
	copy(ids, assetIDs)
	return &MigrationPlanBuilder{
		spec: contracts.MigrationPlanCreateRequest{
			Name:           fmt.Sprintf("saga-migration-%s", runID),
			FromTemplateID: fromTemplateID,
			ToTemplateID:   toTemplateID,
			AssetIDs:       ids,
		},
	}
}

// WithScheduleAt sets a future run time; nil (the default) runs immediately.
func (b *MigrationPlanBuilder) WithScheduleAt(t time.Time) *MigrationPlanBuilder {
	b.spec.ScheduleAt = &t
	return b
}

// WithBatchSize caps how many assets migrate per batch.
func (b *MigrationPlanBuilder) WithBatchSize(n int) *MigrationPlanBuilder {
	b.spec.BatchSize = &n
	return b
}

// WithExtraAssetID appends one more asset id — used to add a fabricated,
// non-existent asset that exercises the per-asset failure path.
func (b *MigrationPlanBuilder) WithExtraAssetID(id string) *MigrationPlanBuilder {
	b.spec.AssetIDs = append(b.spec.AssetIDs, id)
	return b
}

// Build returns the assembled request value.
func (b *MigrationPlanBuilder) Build() contracts.MigrationPlanCreateRequest {
	return b.spec
}
