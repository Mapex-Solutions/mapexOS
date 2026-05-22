package services

import (
	"context"
	"testing"
	"time"

	"assets/src/modules/assets/application/di"
	"assets/src/modules/assets/application/dtos"
	healthPorts "assets/src/modules/healthmonitor/application/ports"
)

// Compile-time check — the inline fake must implement the real port.
var _ healthPorts.HealthRepository = (*fakeHealthRepo)(nil)

// fakeHealthRepo records call counts for batch methods so tests can assert the
// zero-N+1 invariant (exactly one GetLastSeenBatch + IsAlertedBatch per list request).
type fakeHealthRepo struct {
	getLastSeenBatchCalls int
	isAlertedBatchCalls   int
	lastSeenBatchResult   map[string]*time.Time
	alertedBatchResult    map[string]bool

	singleLastSeen *time.Time
	singleAlerted  bool
}

func (f *fakeHealthRepo) UpdateLastSeen(ctx context.Context, orgId string, assetUUID string, ts time.Time) error {
	return nil
}

func (f *fakeHealthRepo) ResetMissCounter(ctx context.Context, orgId string, assetUUID string) error {
	return nil
}

func (f *fakeHealthRepo) IsAlerted(ctx context.Context, orgId string, assetUUID string) (bool, error) {
	return f.singleAlerted, nil
}

func (f *fakeHealthRepo) RemoveAlerted(ctx context.Context, orgId string, assetUUID string) (bool, error) {
	return false, nil
}

func (f *fakeHealthRepo) RegisterOrg(ctx context.Context, orgId string) error { return nil }

func (f *fakeHealthRepo) IsKnownOnline(ctx context.Context, orgId string, assetUUID string) (bool, error) {
	return false, nil
}

func (f *fakeHealthRepo) MarkKnownOnline(ctx context.Context, orgId string, assetUUID string) error {
	return nil
}

func (f *fakeHealthRepo) FindStale(ctx context.Context, orgId string, cutoff time.Time, offset int64, limit int64) ([]string, error) {
	return nil, nil
}

func (f *fakeHealthRepo) IncrementMiss(ctx context.Context, orgId string, assetUUID string) (int64, error) {
	return 0, nil
}

func (f *fakeHealthRepo) MarkAlerted(ctx context.Context, orgId string, assetUUID string) error {
	return nil
}

func (f *fakeHealthRepo) GetActiveOrgs(ctx context.Context) ([]string, error) { return nil, nil }

func (f *fakeHealthRepo) GetLastSeen(ctx context.Context, orgId string, assetUUID string) (*time.Time, error) {
	return f.singleLastSeen, nil
}

func (f *fakeHealthRepo) GetLastSeenBatch(ctx context.Context, orgId string, assetUUIDs []string) (map[string]*time.Time, error) {
	f.getLastSeenBatchCalls++
	return f.lastSeenBatchResult, nil
}

func (f *fakeHealthRepo) IsAlertedBatch(ctx context.Context, orgId string, assetUUIDs []string) (map[string]bool, error) {
	f.isAlertedBatchCalls++
	return f.alertedBatchResult, nil
}

func (f *fakeHealthRepo) RemoveAsset(ctx context.Context, orgId string, assetUUID string) error {
	return nil
}

func (f *fakeHealthRepo) SetLastConnectAt(_ context.Context, _ string, _ string, _ time.Time) error {
	return nil
}

func (f *fakeHealthRepo) GetLastConnectAt(_ context.Context, _ string, _ string) (*time.Time, error) {
	return nil, nil
}

// newTestService builds an AssetService with only HealthRepo wired — sufficient for
// the enrichment helpers which don't touch any other dependency.
func newTestService(repo healthPorts.HealthRepository) *AssetService {
	return &AssetService{deps: di.AssetServiceDependenciesInjection{HealthRepo: repo}}
}

func ptrBool(b bool) *bool       { return &b }
func ptrString(s string) *string { return &s }

// TestEnrichHealthStatusBatch_SingleRoundTrip locks in NFR-1: the list enrichment
// MUST call GetLastSeenBatch + IsAlertedBatch exactly once regardless of page size.
func TestEnrichHealthStatusBatch_SingleRoundTrip(t *testing.T) {
	repo := &fakeHealthRepo{
		lastSeenBatchResult: map[string]*time.Time{},
		alertedBatchResult:  map[string]bool{},
	}
	s := newTestService(repo)

	// 5 responses, all with HealthMonitor.Enabled=true and distinct UUIDs.
	responses := make([]dtos.AssetResponse, 5)
	for i := range responses {
		uuid := "uuid-" + string(rune('a'+i))
		responses[i] = dtos.AssetResponse{
			AssetUUID: ptrString(uuid),
			HealthMonitor: &dtos.HealthMonitorConfig{
				Enabled: ptrBool(true),
			},
		}
	}

	s.enrichHealthStatusBatch(context.Background(), responses, "org-1")

	if repo.getLastSeenBatchCalls != 1 {
		t.Fatalf("GetLastSeenBatch: want 1 call, got %d", repo.getLastSeenBatchCalls)
	}
	if repo.isAlertedBatchCalls != 1 {
		t.Fatalf("IsAlertedBatch: want 1 call, got %d", repo.isAlertedBatchCalls)
	}
}

// TestEnrichHealthStatusBatch_EmptyListNoRedisCall — degenerate case: zero items
// means zero Redis round-trips.
func TestEnrichHealthStatusBatch_EmptyListNoRedisCall(t *testing.T) {
	repo := &fakeHealthRepo{}
	s := newTestService(repo)

	s.enrichHealthStatusBatch(context.Background(), []dtos.AssetResponse{}, "org-1")

	if repo.getLastSeenBatchCalls != 0 {
		t.Fatalf("GetLastSeenBatch: want 0 calls for empty list, got %d", repo.getLastSeenBatchCalls)
	}
}

// TestEnrichHealthStatusBatch_AllDisabledNoRedisCall — when no asset has
// HealthMonitor.Enabled=true, the helper must short-circuit before any Redis call.
func TestEnrichHealthStatusBatch_AllDisabledNoRedisCall(t *testing.T) {
	repo := &fakeHealthRepo{}
	s := newTestService(repo)

	responses := []dtos.AssetResponse{
		{AssetUUID: ptrString("uuid-a"), HealthMonitor: &dtos.HealthMonitorConfig{Enabled: ptrBool(false)}},
		{AssetUUID: ptrString("uuid-b"), HealthMonitor: nil},
	}

	s.enrichHealthStatusBatch(context.Background(), responses, "org-1")

	if repo.getLastSeenBatchCalls != 0 {
		t.Fatalf("GetLastSeenBatch: want 0 calls when all disabled, got %d", repo.getLastSeenBatchCalls)
	}
}
