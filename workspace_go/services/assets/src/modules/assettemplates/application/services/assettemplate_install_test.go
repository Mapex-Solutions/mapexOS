package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/domain/entities"
	redisCache "assets/src/modules/assettemplates/infrastructure/cache/redis"

	v1lists "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/lists"
	commonPorts "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// --- inline fakes for the install flow ---

type fakeMarketplaceClient struct {
	fetchBundleFn func(ctx context.Context, vendor, slug string) (*ports.MarketplaceBundleFetch, error)
}

func (m *fakeMarketplaceClient) FetchBundle(ctx context.Context, vendor, slug string) (*ports.MarketplaceBundleFetch, error) {
	return m.fetchBundleFn(ctx, vendor, slug)
}

type fakeListsClient struct {
	resolveFn func(req v1lists.ListResolveRequest) (string, error)
	requests  []v1lists.ListResolveRequest
}

func (l *fakeListsClient) Resolve(_ context.Context, req v1lists.ListResolveRequest) (string, error) {
	l.requests = append(l.requests, req)
	if l.resolveFn != nil {
		return l.resolveFn(req)
	}
	return "507f1f77bcf86cd799439011", nil
}

// fakeTieredCache is a no-op stub of the tiered cache; the install tests care
// about the link records and resolve calls, not the cached bytes.
type fakeTieredCache struct{ setCalls int }

func (f *fakeTieredCache) Get(context.Context, string) ([]byte, int, error)   { return nil, -1, nil }
func (f *fakeTieredCache) Set(context.Context, string, []byte, time.Duration) error {
	f.setCalls++
	return nil
}
func (f *fakeTieredCache) Delete(context.Context, string) error          { return nil }
func (f *fakeTieredCache) Invalidate(context.Context, string) error      { return nil }
func (f *fakeTieredCache) Stats() commonPorts.LocalCacheStats            { return commonPorts.LocalCacheStats{} }
func (f *fakeTieredCache) GetFromL0(string) ([]byte, bool)               { return nil, false }
func (f *fakeTieredCache) GetFromL1(context.Context, string) ([]byte, error) { return nil, nil }
func (f *fakeTieredCache) Warmup(context.Context, []string) error        { return nil }

// --- helpers ---

func newTestServiceWithInstall(repo *fakeRepo, mc *fakeMarketplaceClient, lc *fakeListsClient, tc *fakeTieredCache) *AssetTemplateService {
	deps := di.AssetTemplateServiceDependenciesInjection{
		AssetTemplateRepo:   repo,
		FieldVocabularyRepo: &fakeVocabRepo{},
		AppCache:            &fakeAppCache{},
		NatsBus:             &fakeFanout{},
		TemplateStoragePort: &fakeStorage{},
		TieredCache:         tc,
		CacheKeyBuilder:     redisCache.NewCacheKeyBuilderAdapter(),
		Metrics:             createTestMetrics(),
		MarketplaceClient:   mc,
		ListsClient:         lc,
	}
	return &AssetTemplateService{deps: deps}
}

func orgRequestContext(orgId string) *reqCtx.RequestContext {
	oid := orgId
	return &reqCtx.RequestContext{OrgContext: &oid}
}

// verifiedFetch builds a fetch whose DeclaredSha256 matches its RawBytes, so it
// passes the checksum hard-verify.
func verifiedFetch(guid string) *ports.MarketplaceBundleFetch {
	raw := []byte(`{"bundle":"` + guid + `"}`)
	sum := sha256.Sum256(raw)
	return &ports.MarketplaceBundleFetch{
		Bundle: dtos.MarketplaceBundle{
			Name:             map[string]string{"en-US": "Temperature"},
			CategorySlug:     "iot-platforms",
			CategoryName:     "IoT Platforms",
			ManufacturerSlug: "disruptive_technologies",
			ManufacturerName: "Disruptive Technologies",
			ModelSlug:        "temperature",
			ModelName:        "Temperature",
			Version:          "1",
		},
		RawBytes:        raw,
		MarketplaceGuid: guid,
		DeclaredSha256:  hex.EncodeToString(sum[:]),
	}
}

func objectIDHex(t *testing.T, hexStr string) model.ObjectId {
	oid, err := model.ToObjectID(hexStr)
	if err != nil {
		t.Fatalf("bad object id %q: %v", hexStr, err)
	}
	return oid
}

// --- tests (FR-7) ---

func TestInstallFromMarketplace_Idempotent(t *testing.T) {
	repo := &fakeRepo{}
	createCount, findCount := 0, 0
	repo.findByMarketplaceGuidAndOrgFn = func(_ context.Context, _ string, _ model.ObjectId) (*entities.Assettemplate, error) {
		findCount++
		if findCount == 1 {
			return nil, nil // first install: not installed yet
		}
		return &entities.Assettemplate{ID: objectIDHex(t, "507f1f77bcf86cd799439021")}, nil // already installed
	}
	repo.createFn = func(_ context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
		createCount++
		e.ID = objectIDHex(t, "507f1f77bcf86cd799439021")
		return e, nil
	}
	mc := &fakeMarketplaceClient{fetchBundleFn: func(_ context.Context, _, _ string) (*ports.MarketplaceBundleFetch, error) {
		return verifiedFetch("guid-1"), nil
	}}
	tc := &fakeTieredCache{}
	svc := newTestServiceWithInstall(repo, mc, &fakeListsClient{}, tc)
	rc := orgRequestContext("507f1f77bcf86cd799439011")

	if _, err := svc.InstallFromMarketplace(context.Background(), rc, "vendor", "slug", false); err != nil {
		t.Fatalf("install 1: %v", err)
	}
	if _, err := svc.InstallFromMarketplace(context.Background(), rc, "vendor", "slug", false); err != nil {
		t.Fatalf("install 2: %v", err)
	}
	if createCount != 1 {
		t.Fatalf("expected exactly 1 link create across two installs, got %d", createCount)
	}
}

func TestInstallFromMarketplace_SharedAcrossOrgs(t *testing.T) {
	repo := &fakeRepo{}
	createCount := 0
	var createdOrgs []string
	repo.findByMarketplaceGuidAndOrgFn = func(_ context.Context, _ string, _ model.ObjectId) (*entities.Assettemplate, error) {
		return nil, nil // each org is a fresh install
	}
	repo.createFn = func(_ context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
		createCount++
		createdOrgs = append(createdOrgs, e.OrgID.Hex())
		e.ID = objectIDHex(t, "507f1f77bcf86cd799439021")
		return e, nil
	}
	mc := &fakeMarketplaceClient{fetchBundleFn: func(_ context.Context, _, _ string) (*ports.MarketplaceBundleFetch, error) {
		return verifiedFetch("guid-shared"), nil
	}}
	tc := &fakeTieredCache{}
	svc := newTestServiceWithInstall(repo, mc, &fakeListsClient{}, tc)

	if _, err := svc.InstallFromMarketplace(context.Background(), orgRequestContext("507f1f77bcf86cd799439011"), "v", "s", false); err != nil {
		t.Fatalf("install org A: %v", err)
	}
	if _, err := svc.InstallFromMarketplace(context.Background(), orgRequestContext("507f1f77bcf86cd799439012"), "v", "s", false); err != nil {
		t.Fatalf("install org B: %v", err)
	}
	if createCount != 2 {
		t.Fatalf("expected 2 per-org links, got %d", createCount)
	}
	if len(createdOrgs) != 2 || createdOrgs[0] == createdOrgs[1] {
		t.Fatalf("expected two distinct org links, got %v", createdOrgs)
	}
	// One shared content cached per guid (idempotent Set), not one per org content doc.
	if tc.setCalls == 0 {
		t.Fatal("expected the shared content to be cached")
	}
}

func TestUninstallFromMarketplace_OnlyRemovesLink(t *testing.T) {
	repo := &fakeRepo{}
	linkID := objectIDHex(t, "507f1f77bcf86cd799439021")
	deleteCount := 0
	var deletedID string
	repo.findByMarketplaceGuidAndOrgFn = func(_ context.Context, _ string, _ model.ObjectId) (*entities.Assettemplate, error) {
		return &entities.Assettemplate{ID: linkID}, nil
	}
	repo.deleteByIdFn = func(_ context.Context, id *string) error {
		deleteCount++
		deletedID = *id
		return nil
	}
	mc := &fakeMarketplaceClient{fetchBundleFn: func(_ context.Context, _, _ string) (*ports.MarketplaceBundleFetch, error) {
		return verifiedFetch("guid-1"), nil
	}}
	svc := newTestServiceWithInstall(repo, mc, &fakeListsClient{}, &fakeTieredCache{})

	if err := svc.UninstallFromMarketplace(context.Background(), orgRequestContext("507f1f77bcf86cd799439011"), "v", "s"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if deleteCount != 1 || deletedID != linkID.Hex() {
		t.Fatalf("expected the caller's link (id=%s) deleted exactly once, got count=%d id=%s", linkID.Hex(), deleteCount, deletedID)
	}
}

func TestInstallFromMarketplace_ChecksumMismatch(t *testing.T) {
	repo := &fakeRepo{}
	createCount := 0
	repo.createFn = func(_ context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
		createCount++
		return e, nil
	}
	mc := &fakeMarketplaceClient{fetchBundleFn: func(_ context.Context, _, _ string) (*ports.MarketplaceBundleFetch, error) {
		return &ports.MarketplaceBundleFetch{
			RawBytes:        []byte("tampered"),
			DeclaredSha256:  "0000000000000000000000000000000000000000000000000000000000000000",
			MarketplaceGuid: "guid-1",
		}, nil
	}}
	lc := &fakeListsClient{}
	tc := &fakeTieredCache{}
	svc := newTestServiceWithInstall(repo, mc, lc, tc)

	_, err := svc.InstallFromMarketplace(context.Background(), orgRequestContext("507f1f77bcf86cd799439011"), "v", "s", false)
	if err == nil {
		t.Fatal("expected a checksum-mismatch error, got nil")
	}
	// Verify runs BEFORE resolve/cache/create: nothing is created or resolved.
	if createCount != 0 {
		t.Fatalf("expected no link created on checksum mismatch, got %d", createCount)
	}
	if len(lc.requests) != 0 {
		t.Fatalf("expected no classification resolve on checksum mismatch, got %d", len(lc.requests))
	}
	if tc.setCalls != 0 {
		t.Fatalf("expected no content cached on checksum mismatch, got %d", tc.setCalls)
	}
}

func TestInstallFromMarketplace_ClassificationOrgScoped(t *testing.T) {
	repo := &fakeRepo{}
	repo.findByMarketplaceGuidAndOrgFn = func(_ context.Context, _ string, _ model.ObjectId) (*entities.Assettemplate, error) {
		return nil, nil
	}
	repo.createFn = func(_ context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
		e.ID = objectIDHex(t, "507f1f77bcf86cd799439021")
		return e, nil
	}
	// Distinct id per list type so the model's parentId can be checked against the
	// manufacturer's resolved id.
	manufacturerID := "507f1f77bcf86cd799439031"
	lc := &fakeListsClient{resolveFn: func(req v1lists.ListResolveRequest) (string, error) {
		switch req.Type {
		case "asset_manufacturer":
			return manufacturerID, nil
		case "asset_category":
			return "507f1f77bcf86cd799439032", nil
		default:
			return "507f1f77bcf86cd799439033", nil
		}
	}}
	mc := &fakeMarketplaceClient{fetchBundleFn: func(_ context.Context, _, _ string) (*ports.MarketplaceBundleFetch, error) {
		return verifiedFetch("guid-1"), nil
	}}
	svc := newTestServiceWithInstall(repo, mc, lc, &fakeTieredCache{})

	if _, err := svc.InstallFromMarketplace(context.Background(), orgRequestContext("507f1f77bcf86cd799439011"), "v", "s", false); err != nil {
		t.Fatalf("install: %v", err)
	}
	if len(lc.requests) != 3 {
		t.Fatalf("expected 3 classification resolves (manufacturer, category, model), got %d", len(lc.requests))
	}
	for _, r := range lc.requests {
		if r.OrgId != "507f1f77bcf86cd799439011" {
			t.Fatalf("classification resolved outside the caller org: %+v", r)
		}
	}
	if lc.requests[0].Type != "asset_manufacturer" {
		t.Fatalf("manufacturer must resolve first, got %s", lc.requests[0].Type)
	}
	model := lc.requests[2]
	if model.Type != "asset_model" || model.ParentId == nil || *model.ParentId != manufacturerID {
		t.Fatalf("model must be parented on the resolved manufacturer id %s, got %+v", manufacturerID, model)
	}
}
