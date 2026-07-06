package services

import (
	ctx "context"
	"fmt"
	"time"

	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// migrationStaleTimerTolerance is how far in the future ScheduleAt may sit before
// a firing is treated as a stale timer, absorbing clock skew.
const migrationStaleTimerTolerance = 5 * time.Second

// guardMigrationRun loads the plan and decides whether this timer firing should
// run it. A nil/zero plan, a non-runnable status, or a stale timer all make the
// firing a no-op (returns false).
func (s *AssetTemplateService) guardMigrationRun(c ctx.Context, planId string) (*entities.MigrationPlan, bool) {
	plan, err := s.deps.MigrationPlanRepo.FindById(c, &planId)
	if err != nil || plan == nil || plan.ID.IsZero() {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplateMigration] run skipped, plan not found planId=%s", planId))
		return nil, false
	}
	// Only a still-pending/scheduled plan is runnable; running/terminal means a
	// redelivery or a plan that was cancelled after this timer was armed.
	if plan.Status != entities.PlanPending && plan.Status != entities.PlanScheduled {
		return nil, false
	}
	// Stale-timer guard: editing a plan arms a NEW timer but cannot cancel the
	// OLD one (ScheduleManager has no per-MsgId cancel), so a firing whose plan
	// now says "run later" is that superseded old timer and must be ignored. The
	// tolerance keeps a legitimately-due timer from being dropped by clock skew.
	if time.Until(plan.ScheduleAt) > migrationStaleTimerTolerance {
		return nil, false
	}
	return plan, true
}

// markMigrationPlanRunning transitions the plan to running, guarding the move
// with CanTransition and reflecting the new status on the in-memory plan.
func (s *AssetTemplateService) markMigrationPlanRunning(c ctx.Context, plan *entities.MigrationPlan) error {
	if !entities.CanTransition(plan.Status, entities.PlanRunning) {
		return nil
	}
	id := plan.ID.Hex()
	if _, err := s.deps.MigrationPlanRepo.FindByIdAndUpdate(c, &id, map[string]any{
		"status":  string(entities.PlanRunning),
		"updated": time.Now(),
	}); err != nil {
		return err
	}
	plan.Status = entities.PlanRunning
	return nil
}

// runMigrationBatches walks the plan's assets in BatchSize batches and switches
// each asset's template, recording every per-asset outcome. A single asset's
// failure never aborts the run.
func (s *AssetTemplateService) runMigrationBatches(c ctx.Context, plan *entities.MigrationPlan) {
	batchSize := plan.BatchSize
	if batchSize <= 0 {
		batchSize = defaultMigrationBatchSize
	}
	target := plan.ToTemplateID.Hex()
	for start := 0; start < len(plan.AssetIDs); start += batchSize {
		end := start + batchSize
		if end > len(plan.AssetIDs) {
			end = len(plan.AssetIDs)
		}
		for _, assetID := range plan.AssetIDs[start:end] {
			s.migrateOneAsset(c, plan, assetID, target)
		}
	}
}

// migrateOneAsset switches one asset to the target template and records the
// outcome on its execution.
//
// The per-asset bookkeeping is NOT wrapped in a Mongo transaction: the writes
// interleave with the external SwitchTemplate call (which drives its own
// side effects), and the application layer holds only repository ports — it has
// no transaction-capable handle and injecting the Mongo manager would break the
// inward dependency rule. Finalization instead recomputes counters from the
// executions collection, so it self-heals after a mid-run crash.
func (s *AssetTemplateService) migrateOneAsset(c ctx.Context, plan *entities.MigrationPlan, assetID model.ObjectId, target string) {
	now := time.Now()
	if err := s.deps.TemplateSwitcher.SwitchTemplate(c, assetID.Hex(), target); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplateMigration] asset switch failed planId=%s assetId=%s: %v", plan.ID.Hex(), assetID.Hex(), err))
		// attempts is 1 here: the run-once plan guard means each asset is
		// attempted exactly once per plan, so the initial 0 becomes 0+1.
		_ = s.deps.MigrationExecutionRepo.UpdateByPlanAndAsset(c, plan.ID, assetID, map[string]any{
			"status":   string(entities.ExecFailed),
			"error":    err.Error(),
			"attempts": 1,
			"updated":  now,
		})
		return
	}
	_ = s.deps.MigrationExecutionRepo.UpdateByPlanAndAsset(c, plan.ID, assetID, map[string]any{
		"status":  string(entities.ExecMigrated),
		"updated": now,
	})
}

// finalizeMigrationPlan recomputes authoritative counters straight from the
// executions collection (self-healing regardless of mid-run crashes) and writes
// the terminal status: complete when nothing failed, else completed_with_errors.
func (s *AssetTemplateService) finalizeMigrationPlan(c ctx.Context, plan *entities.MigrationPlan) error {
	planId := plan.ID.Hex()
	migrated, _ := s.deps.MigrationExecutionRepo.CountByPlanAndStatus(c, planId, entities.ExecMigrated)
	failed, _ := s.deps.MigrationExecutionRepo.CountByPlanAndStatus(c, planId, entities.ExecFailed)

	final := entities.PlanComplete
	if failed > 0 {
		final = entities.PlanCompletedWithErrors
	}
	if !entities.CanTransition(plan.Status, final) {
		return nil
	}
	_, err := s.deps.MigrationPlanRepo.FindByIdAndUpdate(c, &planId, map[string]any{
		"status":            string(final),
		"counters.migrated": int(migrated),
		"counters.failed":   int(failed),
		"updated":           time.Now(),
	})
	logger.Info(fmt.Sprintf("[SERVICE:AssetTemplateMigration] plan finished planId=%s migrated=%d failed=%d", planId, migrated, failed))
	return err
}
