package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"events/src/modules/retention/application/constants"
	"events/src/modules/retention/application/di"
	"events/src/modules/retention/application/ports"
	"events/src/modules/retention/domain/entities"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * Test doubles — inline mocks (no testify).
 */

type fakeRetentionRepo struct {
	createErr      error
	findByIdErr    error
	findByTypeErr  error
	upsertErr      error
	deleteByIdErr  error
	filterErr      error
	upsertCalled   int
	upsertWithNil  int // times called with nil orgId
	upsertLastType string
}

func (f *fakeRetentionRepo) Create(ctx context.Context, policy *entities.RetentionPolicy) (*entities.RetentionPolicy, error) {
	return policy, f.createErr
}
func (f *fakeRetentionRepo) FindById(ctx context.Context, policyId *string) (*entities.RetentionPolicy, error) {
	return nil, f.findByIdErr
}
func (f *fakeRetentionRepo) FindByOrgIdAndType(ctx context.Context, orgId *string, retentionType string) (*entities.RetentionPolicy, error) {
	return nil, f.findByTypeErr
}
func (f *fakeRetentionRepo) Upsert(ctx context.Context, orgId *model.ObjectId, retentionType string, policy *entities.RetentionPolicy) (*entities.RetentionPolicy, error) {
	f.upsertCalled++
	if orgId == nil {
		f.upsertWithNil++
	}
	f.upsertLastType = retentionType
	return policy, f.upsertErr
}
func (f *fakeRetentionRepo) DeleteById(ctx context.Context, policyId *string) error {
	return f.deleteByIdErr
}
func (f *fakeRetentionRepo) FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.RetentionPolicy], error) {
	return nil, f.filterErr
}

/**
 * fakeCacheRepo — only need to not crash the Upsert cache invalidation path.
 */
type fakeCacheRepo struct{}

func (f *fakeCacheRepo) Set(ctx context.Context, key string, value interface{}) error {
	return nil
}
func (f *fakeCacheRepo) SetEx(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}
func (f *fakeCacheRepo) Get(ctx context.Context, key string, dest interface{}) error { return nil }
func (f *fakeCacheRepo) Del(ctx context.Context, key string) error                   { return nil }
func (f *fakeCacheRepo) GetOrSetEx(params common.GetOrSetParams) (any, error) {
	return nil, nil
}

/**
 * fakeCHConn — records ALTER TABLE statements; can be configured to fail.
 * Implements ports.ClickHouseConnPort (minimal surface — only Exec).
 */

type fakeCHConn struct {
	execCalls []string
	execErr   error
}

func (f *fakeCHConn) Exec(ctx context.Context, q string, args ...any) error {
	f.execCalls = append(f.execCalls, q)
	return f.execErr
}

// Compile-time sanity — fakeCHConn must satisfy ClickHouseConnPort at test time.
var _ ports.ClickHouseConnPort = (*fakeCHConn)(nil)

func newTestService(repo *fakeRetentionRepo, ch *fakeCHConn) *RetentionService {
	return &RetentionService{
		deps: di.RetentionServiceDependenciesInjection{
			RetentionRepo:  repo,
			CacheRepo:      &fakeCacheRepo{},
			ClickHouseConn: ch,
		},
	}
}

/**
 * Tests — SeedPlatformPolicies
 */

func TestSeedPlatformPolicies_UpsertsAssetStatusHistoryWithNilOrgId(t *testing.T) {
	repo := &fakeRetentionRepo{}
	ch := &fakeCHConn{}
	svc := newTestService(repo, ch)

	if err := svc.SeedPlatformPolicies(context.Background()); err != nil {
		t.Fatalf("SeedPlatformPolicies returned error: %v", err)
	}

	if repo.upsertCalled != 1 {
		t.Fatalf("expected 1 upsert, got %d", repo.upsertCalled)
	}
	if repo.upsertWithNil != 1 {
		t.Fatalf("expected upsert with nil orgId (platform-level), got %d", repo.upsertWithNil)
	}
	if repo.upsertLastType != constants.TableAssetStatusHistory {
		t.Fatalf("expected type=%s, got %s", constants.TableAssetStatusHistory, repo.upsertLastType)
	}
	if len(ch.execCalls) != 1 {
		t.Fatalf("expected 1 ALTER TABLE call, got %d", len(ch.execCalls))
	}
}

func TestSeedPlatformPolicies_PropagatesRepoError(t *testing.T) {
	repo := &fakeRetentionRepo{upsertErr: errors.New("mongo down")}
	svc := newTestService(repo, &fakeCHConn{})

	if err := svc.SeedPlatformPolicies(context.Background()); err == nil {
		t.Fatal("expected error when repo upsert fails")
	}
}

func TestSeedPlatformPolicies_ContinuesWhenCHTTLFails(t *testing.T) {
	// Seed Mongo write succeeds; CH TTL apply fails — seed MUST still succeed
	// (we log the TTL failure but return nil so boot isn't blocked).
	repo := &fakeRetentionRepo{}
	ch := &fakeCHConn{execErr: errors.New("clickhouse down")}
	svc := newTestService(repo, ch)

	if err := svc.SeedPlatformPolicies(context.Background()); err != nil {
		t.Fatalf("seed should still succeed when CH TTL apply fails, got %v", err)
	}
}

/**
 * Tests — ApplyAssetStatusHistoryTTL
 */

func TestApplyAssetStatusHistoryTTL_BuildsCorrectAlterStatement(t *testing.T) {
	ch := &fakeCHConn{}
	svc := newTestService(&fakeRetentionRepo{}, ch)

	if err := svc.ApplyAssetStatusHistoryTTL(context.Background(), 14); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ch.execCalls) != 1 {
		t.Fatalf("expected 1 Exec call, got %d", len(ch.execCalls))
	}
	want := "ALTER TABLE asset_status_history MODIFY TTL created + toIntervalDay(14)"
	if ch.execCalls[0] != want {
		t.Fatalf("ALTER TABLE mismatch:\n got:  %q\n want: %q", ch.execCalls[0], want)
	}
}

func TestApplyAssetStatusHistoryTTL_RejectsOutOfRangeDays(t *testing.T) {
	svc := newTestService(&fakeRetentionRepo{}, &fakeCHConn{})
	tests := []struct {
		name string
		days uint16
	}{
		{"zero", 0},
		{"above max", 91},
		{"way above max", 9999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.ApplyAssetStatusHistoryTTL(context.Background(), tt.days); err == nil {
				t.Fatalf("expected out-of-range error for days=%d", tt.days)
			}
		})
	}
}

func TestApplyAssetStatusHistoryTTL_ErrorsWhenConnNil(t *testing.T) {
	svc := &RetentionService{
		deps: di.RetentionServiceDependenciesInjection{
			RetentionRepo: &fakeRetentionRepo{},
			CacheRepo:     &fakeCacheRepo{},
			// ClickHouseConn intentionally left nil
		},
	}
	if err := svc.ApplyAssetStatusHistoryTTL(context.Background(), 7); err == nil {
		t.Fatal("expected error when ClickHouseConn is not injected")
	}
}

func TestApplyAssetStatusHistoryTTL_PropagatesExecError(t *testing.T) {
	ch := &fakeCHConn{execErr: errors.New("network blip")}
	svc := newTestService(&fakeRetentionRepo{}, ch)

	if err := svc.ApplyAssetStatusHistoryTTL(context.Background(), 30); err == nil {
		t.Fatal("expected Exec error to propagate")
	}
}
