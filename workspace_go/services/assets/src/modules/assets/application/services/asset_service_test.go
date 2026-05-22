package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"assets/src/bootstrap"
	"assets/src/modules/assets/application/di"
	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/domain/entities"
	redisCache "assets/src/modules/assets/infrastructure/cache/redis"
	templatePorts "assets/src/modules/assettemplates/application/ports"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
	"github.com/prometheus/client_golang/prometheus"
)

// errTest is the sentinel error used by failure-path tests.
var errTest = errors.New("test sentinel error")

/*
 * Inline port mocks. Tests use stdlib + inline mocks of port
 * interfaces so each test wires only the methods it exercises;
 * unset methods return zero values. fakeRouteGroupPort + boolPtr
 * live in asset_helpers_test.go and are reused from this file (same
 * package).
 */

type fakeAssetRepo struct {
	createFn                          func(ctx context.Context, a *entities.Asset) (*entities.Asset, error)
	findByIdFn                        func(ctx context.Context, id *string) (*entities.Asset, error)
	findByAssetUUIDFn                 func(ctx context.Context, uuid *string) (*entities.Asset, error)
	findByMqttUsernameFn              func(ctx context.Context, username string) (*entities.Asset, error)
	findByIdAndUpdateFn               func(ctx context.Context, id *string, payload map[string]any) (*entities.Asset, error)
	deleteByIdFn                      func(ctx context.Context, id *string) error
	findWithFiltersFn                 func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Asset], error)
	findWithFiltersAndTemplateFn      func(ctx context.Context, assetFilters model.Map, templateFilters model.Map, pagination *model.PaginationOpts, sort model.Map) (*model.PaginatedResult[entities.AssetWithTemplate], error)
	countDocumentsFn                  func(ctx context.Context, filters model.Map) (int64, error)
	updateHealthStatusWithChangedAtFn func(ctx context.Context, uuid *string, status string, changedAt time.Time) error
}

func (r *fakeAssetRepo) Create(ctx context.Context, a *entities.Asset) (*entities.Asset, error) {
	if r.createFn != nil {
		return r.createFn(ctx, a)
	}
	return nil, nil
}
func (r *fakeAssetRepo) FindById(ctx context.Context, id *string) (*entities.Asset, error) {
	if r.findByIdFn != nil {
		return r.findByIdFn(ctx, id)
	}
	return nil, nil
}
func (r *fakeAssetRepo) FindByAssetUUID(ctx context.Context, uuid *string) (*entities.Asset, error) {
	if r.findByAssetUUIDFn != nil {
		return r.findByAssetUUIDFn(ctx, uuid)
	}
	return nil, nil
}
func (r *fakeAssetRepo) FindByMqttUsername(ctx context.Context, username string) (*entities.Asset, error) {
	if r.findByMqttUsernameFn != nil {
		return r.findByMqttUsernameFn(ctx, username)
	}
	return nil, nil
}
func (r *fakeAssetRepo) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.Asset, error) {
	if r.findByIdAndUpdateFn != nil {
		return r.findByIdAndUpdateFn(ctx, id, payload)
	}
	return &entities.Asset{}, nil
}
func (r *fakeAssetRepo) DeleteById(ctx context.Context, id *string) error {
	if r.deleteByIdFn != nil {
		return r.deleteByIdFn(ctx, id)
	}
	return nil
}
func (r *fakeAssetRepo) FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Asset], error) {
	if r.findWithFiltersFn != nil {
		return r.findWithFiltersFn(ctx, filters, pagination, projection)
	}
	return &model.PaginatedResult[entities.Asset]{}, nil
}
func (r *fakeAssetRepo) FindWithFiltersAndTemplate(ctx context.Context, assetFilters model.Map, templateFilters model.Map, pagination *model.PaginationOpts, sort model.Map) (*model.PaginatedResult[entities.AssetWithTemplate], error) {
	if r.findWithFiltersAndTemplateFn != nil {
		return r.findWithFiltersAndTemplateFn(ctx, assetFilters, templateFilters, pagination, sort)
	}
	return &model.PaginatedResult[entities.AssetWithTemplate]{}, nil
}
func (r *fakeAssetRepo) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	if r.countDocumentsFn != nil {
		return r.countDocumentsFn(ctx, filters)
	}
	return 0, nil
}
func (r *fakeAssetRepo) UpdateHealthStatusWithChangedAt(ctx context.Context, uuid *string, status string, changedAt time.Time) error {
	if r.updateHealthStatusWithChangedAtFn != nil {
		return r.updateHealthStatusWithChangedAtFn(ctx, uuid, status, changedAt)
	}
	return nil
}

type fakeAssetTemplateRepo struct {
	createFn            func(ctx context.Context, t *templatePorts.Assettemplate) (*templatePorts.Assettemplate, error)
	findByIdFn          func(ctx context.Context, id *string) (*templatePorts.Assettemplate, error)
	findByIdAndUpdateFn func(ctx context.Context, id *string, payload map[string]any) (*templatePorts.Assettemplate, error)
	deleteByIdFn        func(ctx context.Context, id *string) error
	findWithFiltersFn   func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[templatePorts.Assettemplate], error)
	updateManyFn        func(ctx context.Context, filter model.Map, update model.Map) (int64, error)
	countDocumentsFn    func(ctx context.Context, filters model.Map) (int64, error)
}

func (r *fakeAssetTemplateRepo) Create(ctx context.Context, t *templatePorts.Assettemplate) (*templatePorts.Assettemplate, error) {
	if r.createFn != nil {
		return r.createFn(ctx, t)
	}
	return nil, nil
}
func (r *fakeAssetTemplateRepo) FindById(ctx context.Context, id *string) (*templatePorts.Assettemplate, error) {
	if r.findByIdFn != nil {
		return r.findByIdFn(ctx, id)
	}
	return nil, nil
}
func (r *fakeAssetTemplateRepo) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*templatePorts.Assettemplate, error) {
	if r.findByIdAndUpdateFn != nil {
		return r.findByIdAndUpdateFn(ctx, id, payload)
	}
	return &templatePorts.Assettemplate{}, nil
}
func (r *fakeAssetTemplateRepo) DeleteById(ctx context.Context, id *string) error {
	if r.deleteByIdFn != nil {
		return r.deleteByIdFn(ctx, id)
	}
	return nil
}
func (r *fakeAssetTemplateRepo) FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[templatePorts.Assettemplate], error) {
	if r.findWithFiltersFn != nil {
		return r.findWithFiltersFn(ctx, filters, pagination, projection)
	}
	return &model.PaginatedResult[templatePorts.Assettemplate]{}, nil
}
func (r *fakeAssetTemplateRepo) UpdateMany(ctx context.Context, filter model.Map, update model.Map) (int64, error) {
	if r.updateManyFn != nil {
		return r.updateManyFn(ctx, filter, update)
	}
	return 0, nil
}
func (r *fakeAssetTemplateRepo) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	if r.countDocumentsFn != nil {
		return r.countDocumentsFn(ctx, filters)
	}
	return 0, nil
}

type fakeAssetStoragePort struct {
	writeAssetFn       func(ctx context.Context, a *entities.Asset, templateOrgId string) error
	deleteAssetFn      func(ctx context.Context, orgId string, assetUUID string) error
	writeAssetAuthFn   func(ctx context.Context, projection assetsAuthContract.AuthProjection) error
	deleteAssetAuthFn  func(ctx context.Context, assetUUID string) error
}

func (s *fakeAssetStoragePort) WriteAsset(ctx context.Context, a *entities.Asset, templateOrgId string) error {
	if s.writeAssetFn != nil {
		return s.writeAssetFn(ctx, a, templateOrgId)
	}
	return nil
}
func (s *fakeAssetStoragePort) DeleteAsset(ctx context.Context, orgId string, assetUUID string) error {
	if s.deleteAssetFn != nil {
		return s.deleteAssetFn(ctx, orgId, assetUUID)
	}
	return nil
}
func (s *fakeAssetStoragePort) WriteAssetAuth(ctx context.Context, projection assetsAuthContract.AuthProjection) error {
	if s.writeAssetAuthFn != nil {
		return s.writeAssetAuthFn(ctx, projection)
	}
	return nil
}
func (s *fakeAssetStoragePort) DeleteAssetAuth(ctx context.Context, assetUUID string) error {
	if s.deleteAssetAuthFn != nil {
		return s.deleteAssetAuthFn(ctx, assetUUID)
	}
	return nil
}

type fakeFanout struct {
	publishFanoutFn func(ctx context.Context, subject string, data []byte) error
}

func (f *fakeFanout) PublishFanout(ctx context.Context, subject string, data []byte) error {
	if f.publishFanoutFn != nil {
		return f.publishFanoutFn(ctx, subject, data)
	}
	return nil
}
func (f *fakeFanout) SubscribeFanout(stream, serviceName, subject string, handler natsModel.FanoutHandler) (*natsModel.FanoutSubscription, error) {
	return nil, nil
}
func (f *fakeFanout) EnsureFanoutStream(config natsModel.FanoutStreamConfig) error {
	return nil
}

// fakeAppCache satisfies common.AppCache (Cache + CacheGetOrSetEx) with no-ops.
// Get returns errTest by default so the cache-aside lookups behave as misses.
type fakeAppCache struct {
	delFn func(ctx context.Context, key string) error
}

func (c *fakeAppCache) Set(ctx context.Context, key string, value interface{}) error {
	return nil
}
func (c *fakeAppCache) SetEx(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}
func (c *fakeAppCache) Get(ctx context.Context, key string, dest interface{}) error {
	return errTest
}
func (c *fakeAppCache) Del(ctx context.Context, key string) error {
	if c.delFn != nil {
		return c.delFn(ctx, key)
	}
	return nil
}
func (c *fakeAppCache) GetOrSetEx(p common.GetOrSetParams) (any, error) { return nil, nil }

// fakeHealthRepo (the no-op variant satisfying healthPorts.HealthRepository)
// is declared in asset_handler_enrichment_test.go and shared across test files.

type fakeHealthLifecycle struct {
	clearAssetStateFn func(ctx context.Context, orgId string, assetUUID string) error
}

func (h *fakeHealthLifecycle) ClearAssetState(ctx context.Context, orgId string, assetUUID string) error {
	if h.clearAssetStateFn != nil {
		return h.clearAssetStateFn(ctx, orgId, assetUUID)
	}
	return nil
}

/*
 * TEST FIXTURES
 */

// createTestMetrics creates a real AssetsMetrics with a throwaway registry so
// instrumented methods can run without nil-pointer panics.
func createTestMetrics() *bootstrap.AssetsMetrics {
	reg := metrics.NewRegistry("assets_test")

	return &bootstrap.AssetsMetrics{
		Registry: reg,
		AssetOperations: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "asset", Name: "operations_total", Help: "test",
		}, []string{"operation", "status"}),
		AssetOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
			Subsystem: "asset", Name: "operation_duration_seconds", Help: "test", Buckets: prometheus.DefBuckets,
		}, []string{"operation"}),
		AssetListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "asset", Name: "list_results_count", Help: "test", Buckets: []float64{0, 1, 5, 10, 25, 50, 100, 250},
		}),
		AssetCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "asset", Name: "cache_total", Help: "test",
		}, []string{"result"}),
		TemplateOperations: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "template", Name: "operations_total", Help: "test",
		}, []string{"operation", "status"}),
		TemplateOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
			Subsystem: "template", Name: "operation_duration_seconds", Help: "test", Buckets: prometheus.DefBuckets,
		}, []string{"operation"}),
		TemplateListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "template", Name: "list_results_count", Help: "test", Buckets: []float64{0, 1, 5, 10, 25, 50, 100, 250},
		}),
		TemplateCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "template", Name: "cache_total", Help: "test",
		}, []string{"result"}),
		AuthCalloutTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "auth", Name: "callout_total", Help: "test",
		}, []string{"result"}),
		AuthCalloutDuration: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "auth", Name: "callout_duration_seconds", Help: "test", Buckets: prometheus.DefBuckets,
		}),
		AuthCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "auth", Name: "cache_total", Help: "test",
		}, []string{"result"}),
	}
}

type testServiceHandles struct {
	service         *AssetService
	repo            *fakeAssetRepo
	templateRepo    *fakeAssetTemplateRepo
	storage         *fakeAssetStoragePort
	routeGroup      *fakeRouteGroupPort
	cache           *fakeAppCache
	fanout          *fakeFanout
	healthRepo      *fakeHealthRepo
	healthLifecycle *fakeHealthLifecycle
}

// newAssetTestService wires AssetService with inline port mocks.
func newAssetTestService() *testServiceHandles {
	repo := &fakeAssetRepo{}
	tmpl := &fakeAssetTemplateRepo{}
	storage := &fakeAssetStoragePort{}
	rg := &fakeRouteGroupPort{}
	cache := &fakeAppCache{}
	fan := &fakeFanout{}
	hRepo := &fakeHealthRepo{}
	hLifecycle := &fakeHealthLifecycle{}

	deps := di.AssetServiceDependenciesInjection{
		AssetRepo:         repo,
		AssetTemplateRepo: tmpl,
		AppCache:          cache,
		NatsBus:           fan,
		RouteGroupPort:    rg,
		AssetStoragePort:  storage,
		CacheKeyBuilder:   redisCache.NewCacheKeyBuilderAdapter(),
		Metrics:           createTestMetrics(),
		HealthRepo:        hRepo,
		HealthLifecycle:   hLifecycle,
	}

	return &testServiceHandles{
		service:         &AssetService{deps: deps},
		repo:            repo,
		templateRepo:    tmpl,
		storage:         storage,
		routeGroup:      rg,
		cache:           cache,
		fanout:          fan,
		healthRepo:      hRepo,
		healthLifecycle: hLifecycle,
	}
}

/*
 * TESTS
 */

func TestCreateAsset(t *testing.T) {
	t.Run("should create asset successfully", func(t *testing.T) {
		h := newAssetTestService()

		ctx := context.Background()
		orgId := "68f5bbce1aef22967c3ebb30"
		pathKey := "P1234"
		orgObjectId, _ := model.ToObjectID(orgId)
		templateId := "691bb4071e717d77a2430b46"

		requestContext := &reqCtx.RequestContext{
			OrgContext: &orgId,
			OrgContextData: &reqCtx.CoverageOrg{
				PathKey: pathKey,
			},
		}

		dto := &dtos.AssetCreateDTO{
			Name:            "Test Asset",
			Enabled:         true,
			AssetUUID:       "device-uuid-123",
			AssetTemplateID: templateId,
			RouteGroupIds:   []string{"route-group-1"},
			Protocol: assetsContract.ProtocolType{
				Type: "mqtt",
				Mqtt: &assetsContract.MqttConfig{
					ClientId: "client-1",
					Username: "user-1",
					Password: "test-password-12345678",
				},
			},
		}

		createdAsset := &entities.Asset{
			ID:        orgObjectId,
			Name:      "Test Asset",
			Enabled:   true,
			AssetUUID: "device-uuid-123",
			OrgID:     orgObjectId,
			PathKey:   pathKey,
			RouteGroupIds: []string{"route-group-1"},
			Protocol: entities.ProtocolType{
				Type: "mqtt",
				Mqtt: &entities.MqttConfig{
					ClientId: "client-1",
					Username: "user-1",
				},
			},
			Created: time.Now(),
			Updated: time.Now(),
		}

		var createCalls int
		h.repo.createFn = func(_ context.Context, _ *entities.Asset) (*entities.Asset, error) {
			createCalls++
			return createdAsset, nil
		}

		result, err := h.service.CreateAsset(ctx, requestContext, dto)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Name == nil || *result.Name != "Test Asset" {
			t.Fatalf("expected Name='Test Asset', got %#v", result)
		}
		if createCalls != 1 {
			t.Fatalf("expected Create to be called exactly once, got %d", createCalls)
		}
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		h := newAssetTestService()

		orgId := "68f5bbce1aef22967c3ebb30"
		requestContext := &reqCtx.RequestContext{
			OrgContext: &orgId,
			OrgContextData: &reqCtx.CoverageOrg{
				PathKey: "P1234",
			},
		}
		dto := &dtos.AssetCreateDTO{
			Name:            "Test Asset",
			Enabled:         true,
			AssetUUID:       "device-uuid-123",
			AssetTemplateID: "691bb4071e717d77a2430b46",
			RouteGroupIds:   []string{"route-group-1"},
		}

		h.repo.createFn = func(_ context.Context, _ *entities.Asset) (*entities.Asset, error) {
			return nil, errTest
		}

		result, err := h.service.CreateAsset(context.Background(), requestContext, dto)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestGetAssetById(t *testing.T) {
	t.Run("should get asset by id successfully", func(t *testing.T) {
		h := newAssetTestService()

		ctx := context.Background()
		assetId := "68f5bbce1aef22967c3ebb30"
		templateId, _ := model.ToObjectID("691bb4071e717d77a2430b46")
		orgId, _ := model.ToObjectID("68f5bbce1aef22967c3ebb30")

		existingAsset := &entities.Asset{
			ID:              orgId,
			Name:            "Test Asset",
			Enabled:         true,
			AssetUUID:       "device-uuid-123",
			OrgID:           orgId,
			AssetTemplateID: templateId,
			RouteGroupIds:   []string{"route-group-1", "route-group-2"},
		}

		template := &templatePorts.Assettemplate{
			ID:               templateId,
			Name:             "Test Template",
			ManufacturerName: typeconv.PtrString("Test Manufacturer"),
			ModelName:        typeconv.PtrString("Test Model"),
			CategoryName:     typeconv.PtrString("Test Category"),
			Version:          typeconv.PtrString("1.0"),
			AssetIDPath:      "$.deviceId",
		}

		h.repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return existingAsset, nil
		}
		h.templateRepo.findByIdFn = func(_ context.Context, _ *string) (*templatePorts.Assettemplate, error) {
			return template, nil
		}
		h.routeGroup.namesByIdsFn = func(_ context.Context, ids []string) ([]string, error) {
			return []string{"Route 1", "Route 2"}, nil
		}

		result, err := h.service.GetAssetById(ctx, &assetId)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Name == nil || *result.Name != "Test Asset" {
			t.Fatalf("expected Name='Test Asset', got %#v", result)
		}
		if result.ManufacturerName == nil || *result.ManufacturerName != "Test Manufacturer" {
			t.Fatalf("expected ManufacturerName='Test Manufacturer', got %#v", result.ManufacturerName)
		}
		if result.RouteGroupNames == nil || len(*result.RouteGroupNames) != 2 {
			t.Fatalf("expected 2 RouteGroupNames, got %#v", result.RouteGroupNames)
		}
	})

	t.Run("should return error when asset not found", func(t *testing.T) {
		h := newAssetTestService()
		h.repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return nil, nil
		}

		assetId := "nonexistent-id"
		result, err := h.service.GetAssetById(context.Background(), &assetId)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestUpdateAssetById(t *testing.T) {
	t.Run("should return error when asset not found on initial fetch", func(t *testing.T) {
		h := newAssetTestService()
		h.repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return nil, nil
		}

		assetId := "nonexistent-id"
		newName := "Updated Asset"
		dto := &dtos.AssetUpdateDTO{Name: &newName}

		result, err := h.service.UpdateAssetById(context.Background(), &assetId, dto)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestDeleteAssetById(t *testing.T) {
	t.Run("should return error when delete fails", func(t *testing.T) {
		h := newAssetTestService()

		assetId := "68f5bbce1aef22967c3ebb30"
		orgId, _ := model.ToObjectID("68f5bbce1aef22967c3ebb30")

		h.repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return &entities.Asset{ID: orgId, AssetUUID: "device-uuid-123", OrgID: orgId}, nil
		}
		h.repo.deleteByIdFn = func(_ context.Context, _ *string) error { return errTest }

		result, err := h.service.DeleteAssetById(context.Background(), &assetId)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestGetAssetByMqttUsername(t *testing.T) {
	t.Run("should get asset by mqtt username successfully", func(t *testing.T) {
		h := newAssetTestService()

		ctx := context.Background()
		username := "mqtt-user-123"
		orgId, _ := model.ToObjectID("68f5bbce1aef22967c3ebb30")

		existingAsset := &entities.Asset{
			ID:        orgId,
			Name:      "Test Asset",
			Enabled:   true,
			AssetUUID: "device-uuid-123",
			OrgID:     orgId,
			Protocol: entities.ProtocolType{
				Type: "mqtt",
				Mqtt: &entities.MqttConfig{Username: username},
			},
		}
		h.repo.findByMqttUsernameFn = func(_ context.Context, _ string) (*entities.Asset, error) {
			return existingAsset, nil
		}

		result, err := h.service.GetAssetByMqttUsername(ctx, username)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Name == nil || *result.Name != "Test Asset" {
			t.Fatalf("expected Name='Test Asset', got %#v", result)
		}
	})

	t.Run("should return error when asset not found", func(t *testing.T) {
		h := newAssetTestService()
		h.repo.findByMqttUsernameFn = func(_ context.Context, _ string) (*entities.Asset, error) {
			return nil, nil
		}

		result, err := h.service.GetAssetByMqttUsername(context.Background(), "nonexistent-user")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestGetAssetReadModelByUUID(t *testing.T) {
	t.Run("should get asset read model by UUID successfully", func(t *testing.T) {
		h := newAssetTestService()

		ctx := context.Background()
		assetUUID := "device-uuid-123"
		orgId, _ := model.ToObjectID("68f5bbce1aef22967c3ebb30")
		templateId, _ := model.ToObjectID("691bb4071e717d77a2430b46")

		existingAsset := &entities.Asset{
			ID:              orgId,
			Name:            "Test Asset",
			Enabled:         true,
			AssetUUID:       assetUUID,
			OrgID:           orgId,
			AssetTemplateID: templateId,
			PathKey:         "P1234",
			Protocol: entities.ProtocolType{
				Type: "mqtt",
				Mqtt: &entities.MqttConfig{ClientId: "client-1", Username: "user-1"},
			},
		}
		template := &templatePorts.Assettemplate{ID: templateId, IsSystem: true}

		h.repo.findByAssetUUIDFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return existingAsset, nil
		}
		h.templateRepo.findByIdFn = func(_ context.Context, _ *string) (*templatePorts.Assettemplate, error) {
			return template, nil
		}

		var capturedTemplateOrgId string
		h.storage.writeAssetFn = func(_ context.Context, _ *entities.Asset, templateOrgId string) error {
			capturedTemplateOrgId = templateOrgId
			return nil
		}

		result, err := h.service.GetAssetReadModelByUUID(ctx, assetUUID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.UUID != assetUUID {
			t.Fatalf("expected UUID=%q, got %q", assetUUID, result.UUID)
		}
		if result.OrgId != orgId.Hex() {
			t.Fatalf("expected OrgId=%q, got %q", orgId.Hex(), result.OrgId)
		}
		if capturedTemplateOrgId != "mapexos_public" {
			t.Fatalf("expected templateOrgId='mapexos_public' for system template, got %q", capturedTemplateOrgId)
		}
	})

	t.Run("should return error when asset not found", func(t *testing.T) {
		h := newAssetTestService()
		h.repo.findByAssetUUIDFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return nil, nil
		}

		result, err := h.service.GetAssetReadModelByUUID(context.Background(), "nonexistent-uuid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestGetTemplateOrgId(t *testing.T) {
	t.Run("should return mapexos_public for system template", func(t *testing.T) {
		h := newAssetTestService()

		ctx := context.Background()
		assetUUID := "device-uuid-123"
		orgId, _ := model.ToObjectID("68f5bbce1aef22967c3ebb30")
		templateId, _ := model.ToObjectID("691bb4071e717d77a2430b46")

		existingAsset := &entities.Asset{
			ID:              orgId,
			AssetUUID:       assetUUID,
			OrgID:           orgId,
			AssetTemplateID: templateId,
		}
		systemTemplate := &templatePorts.Assettemplate{ID: templateId, IsSystem: true}

		h.repo.findByAssetUUIDFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return existingAsset, nil
		}
		h.templateRepo.findByIdFn = func(_ context.Context, _ *string) (*templatePorts.Assettemplate, error) {
			return systemTemplate, nil
		}

		var capturedTemplateOrgId string
		h.storage.writeAssetFn = func(_ context.Context, _ *entities.Asset, templateOrgId string) error {
			capturedTemplateOrgId = templateOrgId
			return nil
		}

		_, err := h.service.GetAssetReadModelByUUID(ctx, assetUUID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedTemplateOrgId != "mapexos_public" {
			t.Fatalf("expected templateOrgId='mapexos_public', got %q", capturedTemplateOrgId)
		}
	})

	t.Run("should return asset orgId for non-system template", func(t *testing.T) {
		h := newAssetTestService()

		ctx := context.Background()
		assetUUID := "device-uuid-123"
		orgId, _ := model.ToObjectID("68f5bbce1aef22967c3ebb30")
		templateId, _ := model.ToObjectID("691bb4071e717d77a2430b46")

		existingAsset := &entities.Asset{
			ID:              orgId,
			AssetUUID:       assetUUID,
			OrgID:           orgId,
			AssetTemplateID: templateId,
		}
		privateTemplate := &templatePorts.Assettemplate{
			ID:       templateId,
			IsSystem: false,
			OrgID:    &orgId,
		}

		h.repo.findByAssetUUIDFn = func(_ context.Context, _ *string) (*entities.Asset, error) {
			return existingAsset, nil
		}
		h.templateRepo.findByIdFn = func(_ context.Context, _ *string) (*templatePorts.Assettemplate, error) {
			return privateTemplate, nil
		}

		var capturedTemplateOrgId string
		h.storage.writeAssetFn = func(_ context.Context, _ *entities.Asset, templateOrgId string) error {
			capturedTemplateOrgId = templateOrgId
			return nil
		}

		_, err := h.service.GetAssetReadModelByUUID(ctx, assetUUID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedTemplateOrgId != orgId.Hex() {
			t.Fatalf("expected templateOrgId=%q, got %q", orgId.Hex(), capturedTemplateOrgId)
		}
	})
}
