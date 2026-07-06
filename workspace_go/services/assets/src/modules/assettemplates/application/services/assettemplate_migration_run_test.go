package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/domain/entities"
	"assets/src/modules/assettemplates/domain/repositories"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * Inline port mocks for the migration run path. Only the three ports the run
 * path touches are faked (MigrationPlanRepo, MigrationExecutionRepo,
 * TemplateSwitcher); the remaining DI fields stay zero and are never accessed.
 * Stdlib + inline mocks only (no testify), per the unit-test standard.
 */

// planUpdate records one FindByIdAndUpdate call on the plan repo mock.
type planUpdate struct {
	id      string
	payload map[string]any
}

// mockMigrationPlanRepo returns a configured plan from FindById and records
// every FindByIdAndUpdate payload.
type mockMigrationPlanRepo struct {
	plan        *entities.MigrationPlan
	findErr     error
	updateCalls []planUpdate
}

var _ repositories.MigrationPlanRepository = (*mockMigrationPlanRepo)(nil)

func (m *mockMigrationPlanRepo) Create(_ context.Context, p *entities.MigrationPlan) (*entities.MigrationPlan, error) {
	return p, nil
}

func (m *mockMigrationPlanRepo) FindById(_ context.Context, _ *string) (*entities.MigrationPlan, error) {
	return m.plan, m.findErr
}

func (m *mockMigrationPlanRepo) FindByIdAndUpdate(_ context.Context, id *string, payload map[string]any) (*entities.MigrationPlan, error) {
	idStr := ""
	if id != nil {
		idStr = *id
	}
	m.updateCalls = append(m.updateCalls, planUpdate{id: idStr, payload: payload})
	return m.plan, nil
}

func (m *mockMigrationPlanRepo) FindWithFilters(_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[entities.MigrationPlan], error) {
	return &model.PaginatedResult[entities.MigrationPlan]{}, nil
}

func (m *mockMigrationPlanRepo) IncrementCounter(_ context.Context, _ string, _ string, _ int) error {
	return nil
}

// execUpdate records one UpdateByPlanAndAsset call on the execution repo mock.
type execUpdate struct {
	planID  model.ObjectId
	assetID model.ObjectId
	fields  map[string]any
}

// mockMigrationExecutionRepo records UpdateByPlanAndAsset calls and returns
// configured counts per status from CountByPlanAndStatus.
type mockMigrationExecutionRepo struct {
	migratedCount int64
	failedCount   int64
	updateCalls   []execUpdate
}

var _ repositories.MigrationExecutionRepository = (*mockMigrationExecutionRepo)(nil)

func (m *mockMigrationExecutionRepo) BulkInsert(_ context.Context, execs []*entities.MigrationExecution) (int, error) {
	return len(execs), nil
}

func (m *mockMigrationExecutionRepo) FindByPlan(_ context.Context, _ *string, _ model.Map, _ *model.PaginationOpts) (*model.PaginatedResult[entities.MigrationExecution], error) {
	return &model.PaginatedResult[entities.MigrationExecution]{}, nil
}

func (m *mockMigrationExecutionRepo) FindByIdAndUpdate(_ context.Context, _ *string, _ map[string]any) (*entities.MigrationExecution, error) {
	return &entities.MigrationExecution{}, nil
}

func (m *mockMigrationExecutionRepo) CountByPlanAndStatus(_ context.Context, _ string, status entities.ExecStatus) (int64, error) {
	switch status {
	case entities.ExecMigrated:
		return m.migratedCount, nil
	case entities.ExecFailed:
		return m.failedCount, nil
	default:
		return 0, nil
	}
}

func (m *mockMigrationExecutionRepo) UpdateByPlanAndAsset(_ context.Context, planID, assetID model.ObjectId, fields map[string]any) error {
	m.updateCalls = append(m.updateCalls, execUpdate{planID: planID, assetID: assetID, fields: fields})
	return nil
}

func (m *mockMigrationExecutionRepo) CancelPending(_ context.Context, _ model.ObjectId) (int64, error) {
	return 0, nil
}

// mockTemplateSwitcher counts SwitchTemplate calls and fails for a configured asset id.
type mockTemplateSwitcher struct {
	failAssetID string
	switched    []string
}

var _ ports.TemplateSwitcherPort = (*mockTemplateSwitcher)(nil)

func (m *mockTemplateSwitcher) SwitchTemplate(_ context.Context, assetID, _ string) error {
	m.switched = append(m.switched, assetID)
	if m.failAssetID != "" && assetID == m.failAssetID {
		return errors.New("switch failed")
	}
	return nil
}

// newMigrationTestService wires a service with only the migration-run mocks; the
// remaining DI fields stay zero because the run path never reaches them.
func newMigrationTestService(planRepo *mockMigrationPlanRepo, execRepo *mockMigrationExecutionRepo, switcher *mockTemplateSwitcher) ports.AssetTemplateServicePort {
	return New(di.AssetTemplateServiceDependenciesInjection{
		MigrationPlanRepo:      planRepo,
		MigrationExecutionRepo: execRepo,
		TemplateSwitcher:       switcher,
	})
}

// buildMigrationPlan builds a plan with n assets for the run-path tests.
func buildMigrationPlan(status entities.PlanStatus, scheduleAt time.Time, batchSize, assetCount int) *entities.MigrationPlan {
	assetIDs := make([]model.ObjectId, assetCount)
	for i := range assetIDs {
		assetIDs[i] = model.NewObjectID()
	}
	return &entities.MigrationPlan{
		ID:           model.NewObjectID(),
		ToTemplateID: model.NewObjectID(),
		AssetIDs:     assetIDs,
		ScheduleAt:   scheduleAt,
		BatchSize:    batchSize,
		Status:       status,
	}
}

// finalUpdate returns the last recorded plan update, which finalization writes.
func finalUpdate(t *testing.T, planRepo *mockMigrationPlanRepo) map[string]any {
	t.Helper()
	if len(planRepo.updateCalls) == 0 {
		t.Fatal("expected at least one plan update, got none")
	}
	return planRepo.updateCalls[len(planRepo.updateCalls)-1].payload
}

func TestRunMigrationPlan_HappyPath(t *testing.T) {
	plan := buildMigrationPlan(entities.PlanScheduled, time.Now().Add(-time.Minute), 2, 5)
	planRepo := &mockMigrationPlanRepo{plan: plan}
	execRepo := &mockMigrationExecutionRepo{migratedCount: 5, failedCount: 0}
	switcher := &mockTemplateSwitcher{}
	service := newMigrationTestService(planRepo, execRepo, switcher)

	if err := service.RunMigrationPlan(context.Background(), plan.ID.Hex()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(switcher.switched) != 5 {
		t.Fatalf("expected SwitchTemplate called 5 times, got %d", len(switcher.switched))
	}
	// running transition, then finalize -> two plan updates
	if len(planRepo.updateCalls) != 2 {
		t.Fatalf("expected 2 plan updates (running + finalize), got %d", len(planRepo.updateCalls))
	}
	if got := planRepo.updateCalls[0].payload["status"]; got != string(entities.PlanRunning) {
		t.Fatalf("expected first update to running, got %v", got)
	}
	final := finalUpdate(t, planRepo)
	if got := final["status"]; got != string(entities.PlanComplete) {
		t.Fatalf("expected final status complete, got %v", got)
	}
	if got := final["counters.migrated"]; got != 5 {
		t.Fatalf("expected counters.migrated=5, got %v", got)
	}
	if got := final["counters.failed"]; got != 0 {
		t.Fatalf("expected counters.failed=0, got %v", got)
	}
}

func TestRunMigrationPlan_ContinuesOnFailure(t *testing.T) {
	plan := buildMigrationPlan(entities.PlanScheduled, time.Now().Add(-time.Minute), 2, 5)
	failAsset := plan.AssetIDs[2]
	planRepo := &mockMigrationPlanRepo{plan: plan}
	execRepo := &mockMigrationExecutionRepo{migratedCount: 4, failedCount: 1}
	switcher := &mockTemplateSwitcher{failAssetID: failAsset.Hex()}
	service := newMigrationTestService(planRepo, execRepo, switcher)

	if err := service.RunMigrationPlan(context.Background(), plan.ID.Hex()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// all 5 assets are still attempted despite the single failure
	if len(switcher.switched) != 5 {
		t.Fatalf("expected SwitchTemplate called 5 times, got %d", len(switcher.switched))
	}
	if len(execRepo.updateCalls) != 5 {
		t.Fatalf("expected 5 execution updates, got %d", len(execRepo.updateCalls))
	}
	// the failed asset's execution carries the failed status; every other asset migrated
	for _, call := range execRepo.updateCalls {
		wantStatus := string(entities.ExecMigrated)
		if call.assetID == failAsset {
			wantStatus = string(entities.ExecFailed)
		}
		if got := call.fields["status"]; got != wantStatus {
			t.Fatalf("asset %s: expected status %q, got %v", call.assetID.Hex(), wantStatus, got)
		}
	}
	final := finalUpdate(t, planRepo)
	if got := final["status"]; got != string(entities.PlanCompletedWithErrors) {
		t.Fatalf("expected final status completed_with_errors, got %v", got)
	}
	if got := final["counters.migrated"]; got != 4 {
		t.Fatalf("expected counters.migrated=4, got %v", got)
	}
	if got := final["counters.failed"]; got != 1 {
		t.Fatalf("expected counters.failed=1, got %v", got)
	}
}

func TestRunMigrationPlan_IdempotentNoOps(t *testing.T) {
	tests := []struct {
		name   string
		status entities.PlanStatus
	}{
		{"running is a no-op", entities.PlanRunning},
		{"complete is a no-op", entities.PlanComplete},
		{"cancelled is a no-op", entities.PlanCancelled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := buildMigrationPlan(tt.status, time.Now().Add(-time.Minute), 2, 5)
			planRepo := &mockMigrationPlanRepo{plan: plan}
			execRepo := &mockMigrationExecutionRepo{}
			switcher := &mockTemplateSwitcher{}
			service := newMigrationTestService(planRepo, execRepo, switcher)

			if err := service.RunMigrationPlan(context.Background(), plan.ID.Hex()); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(switcher.switched) != 0 {
				t.Fatalf("expected switcher never called, got %d calls", len(switcher.switched))
			}
			if len(planRepo.updateCalls) != 0 {
				t.Fatalf("expected no plan updates, got %d", len(planRepo.updateCalls))
			}
		})
	}
}

func TestRunMigrationPlan_StaleTimerGuard(t *testing.T) {
	// scheduled but due an hour out -> superseded old timer, must no-op
	plan := buildMigrationPlan(entities.PlanScheduled, time.Now().Add(time.Hour), 2, 5)
	planRepo := &mockMigrationPlanRepo{plan: plan}
	execRepo := &mockMigrationExecutionRepo{}
	switcher := &mockTemplateSwitcher{}
	service := newMigrationTestService(planRepo, execRepo, switcher)

	if err := service.RunMigrationPlan(context.Background(), plan.ID.Hex()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(switcher.switched) != 0 {
		t.Fatalf("expected switcher never called, got %d calls", len(switcher.switched))
	}
	if len(planRepo.updateCalls) != 0 {
		t.Fatalf("expected no plan updates, got %d", len(planRepo.updateCalls))
	}
}
