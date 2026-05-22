package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"assets/src/bootstrap"
	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"
	redisCache "assets/src/modules/assettemplates/infrastructure/cache/redis"

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
 * interfaces (no testify, no shared mock packages). Each fake exposes
 * function fields for the methods the tests configure; unset methods
 * return zero values so tests only wire what they need.
 */

type fakeRepo struct {
	createFn            func(ctx context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error)
	findByIdFn          func(ctx context.Context, id *string) (*entities.Assettemplate, error)
	findByIdAndUpdateFn func(ctx context.Context, id *string, payload map[string]any) (*entities.Assettemplate, error)
	deleteByIdFn        func(ctx context.Context, id *string) error
	findWithFiltersFn   func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Assettemplate], error)
	updateManyFn        func(ctx context.Context, filter model.Map, update model.Map) (int64, error)
	countDocumentsFn    func(ctx context.Context, filters model.Map) (int64, error)

	findByIdCalls int
}

func (r *fakeRepo) Create(ctx context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
	if r.createFn != nil {
		return r.createFn(ctx, e)
	}
	return nil, nil
}

func (r *fakeRepo) FindById(ctx context.Context, id *string) (*entities.Assettemplate, error) {
	r.findByIdCalls++
	if r.findByIdFn != nil {
		return r.findByIdFn(ctx, id)
	}
	return nil, nil
}

func (r *fakeRepo) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.Assettemplate, error) {
	if r.findByIdAndUpdateFn != nil {
		return r.findByIdAndUpdateFn(ctx, id, payload)
	}
	return &entities.Assettemplate{}, nil
}

func (r *fakeRepo) DeleteById(ctx context.Context, id *string) error {
	if r.deleteByIdFn != nil {
		return r.deleteByIdFn(ctx, id)
	}
	return nil
}

func (r *fakeRepo) FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Assettemplate], error) {
	if r.findWithFiltersFn != nil {
		return r.findWithFiltersFn(ctx, filters, pagination, projection)
	}
	return &model.PaginatedResult[entities.Assettemplate]{}, nil
}

func (r *fakeRepo) UpdateMany(ctx context.Context, filter model.Map, update model.Map) (int64, error) {
	if r.updateManyFn != nil {
		return r.updateManyFn(ctx, filter, update)
	}
	return 0, nil
}

func (r *fakeRepo) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	if r.countDocumentsFn != nil {
		return r.countDocumentsFn(ctx, filters)
	}
	return 0, nil
}

type fakeStorage struct {
	writeScriptsFn  func(ctx context.Context, t *entities.Assettemplate) error
	deleteScriptsFn func(ctx context.Context, orgId, templateId string) error
}

func (s *fakeStorage) WriteScripts(ctx context.Context, t *entities.Assettemplate) error {
	if s.writeScriptsFn != nil {
		return s.writeScriptsFn(ctx, t)
	}
	return nil
}

func (s *fakeStorage) DeleteScripts(ctx context.Context, orgId, templateId string) error {
	if s.deleteScriptsFn != nil {
		return s.deleteScriptsFn(ctx, orgId, templateId)
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
// Production code only invokes Del/Get/SetEx in counter-cache paths; tests
// don't exercise those, so default no-ops are sufficient.
type fakeAppCache struct{}

func (c *fakeAppCache) Set(ctx context.Context, key string, value interface{}) error {
	return nil
}
func (c *fakeAppCache) SetEx(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}
func (c *fakeAppCache) Get(ctx context.Context, key string, dest interface{}) error {
	return errTest
}
func (c *fakeAppCache) Del(ctx context.Context, key string) error       { return nil }
func (c *fakeAppCache) GetOrSetEx(p common.GetOrSetParams) (any, error) { return nil, nil }

/*
 * TEST FIXTURES
 */

// createTestMetrics creates a real AssetsMetrics with a throwaway registry so
// instrumented methods can run without nil-pointer panics.
func createTestMetrics() *bootstrap.AssetsMetrics {
	reg := metrics.NewRegistry("assets_tmpl_test")

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

// newTestService wires a service with inline port mocks. Returned mocks let
// each subtest configure only the call paths it cares about.
func newTestService() (*AssetTemplateService, *fakeRepo, *fakeStorage, *fakeFanout) {
	repo := &fakeRepo{}
	storage := &fakeStorage{}
	fanout := &fakeFanout{}

	deps := di.AssetTemplateServiceDependenciesInjection{
		AssetTemplateRepo:   repo,
		AppCache:            &fakeAppCache{},
		NatsBus:             fanout,
		TemplateStoragePort: storage,
		TieredCache:         nil, // never accessed by production code under test
		CacheKeyBuilder:     redisCache.NewCacheKeyBuilderAdapter(),
		Metrics:             createTestMetrics(),
	}

	return &AssetTemplateService{deps: deps}, repo, storage, fanout
}

/*
 * TESTS
 */

func TestCreateAssetTemplate(t *testing.T) {
	t.Run("should return error when repository fails", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.createFn = func(_ context.Context, _ *entities.Assettemplate) (*entities.Assettemplate, error) {
			return nil, errTest
		}

		validator := "function validate() {}"
		dto := &dtos.AssetTemplateCreateDTO{
			Name:             "Test Template",
			Enabled:          true,
			IsSystem:         true,
			ScriptValidator:  &validator,
			ScriptConversion: "function convert() {}",
		}

		result, err := service.CreateAssetTemplate(context.Background(), &reqCtx.RequestContext{}, dto)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestGetAssetTemplateById(t *testing.T) {
	t.Run("should get template by id successfully", func(t *testing.T) {
		service, repo, _, _ := newTestService()

		ctx := context.Background()
		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)
		name := "Test Template"

		existing := &entities.Assettemplate{
			ID:               templateObjectId,
			Name:             name,
			Enabled:          true,
			IsSystem:         true,
			ScriptValidator:  "function validate() {}",
			ScriptConversion: "function convert() {}",
		}
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return existing, nil
		}

		result, err := service.GetAssetTemplateById(ctx, &templateId)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Name == nil || *result.Name != name {
			t.Fatalf("expected name=%q, got %#v", name, result)
		}
	})

	t.Run("should return error when template not found", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return nil, nil
		}

		templateId := "nonexistent-id"
		result, err := service.GetAssetTemplateById(context.Background(), &templateId)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestUpdateAssetTemplateById(t *testing.T) {
	t.Run("should return error when template not found", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return nil, nil
		}

		templateId := "nonexistent-id"
		newName := "Updated Template"
		dto := &dtos.AssetTemplateUpdateDTO{Name: &newName}

		result, err := service.UpdateAssetTemplateById(context.Background(), &templateId, dto)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestDeleteAssetTemplateById(t *testing.T) {
	t.Run("should return error when template not found", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return nil, nil
		}

		templateId := "nonexistent-id"
		result, err := service.DeleteAssetTemplateById(context.Background(), &templateId)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)

		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return &entities.Assettemplate{ID: templateObjectId, Name: "Test Template", IsSystem: true}, nil
		}
		repo.deleteByIdFn = func(_ context.Context, _ *string) error { return errTest }

		result, err := service.DeleteAssetTemplateById(context.Background(), &templateId)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestGetTemplateByIdForCacheFallback(t *testing.T) {
	t.Run("should fetch template and repopulate L2 cache", func(t *testing.T) {
		service, repo, storage, _ := newTestService()
		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)
		name := "Test Template"

		existing := &entities.Assettemplate{
			ID:               templateObjectId,
			Name:             name,
			Enabled:          true,
			IsSystem:         true,
			ScriptValidator:  "function validate() {}",
			ScriptConversion: "function convert() {}",
		}
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return existing, nil
		}

		writeScriptsCalls := 0
		storage.writeScriptsFn = func(_ context.Context, t *entities.Assettemplate) error {
			writeScriptsCalls++
			if t != existing {
				return errors.New("expected the same template instance to be passed")
			}
			return nil
		}

		result, err := service.GetTemplateByIdForCacheFallback(context.Background(), templateId)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Name == nil || *result.Name != name {
			t.Fatalf("expected name=%q, got %#v", name, result)
		}
		if writeScriptsCalls != 1 {
			t.Fatalf("expected WriteScripts to be called exactly once, got %d", writeScriptsCalls)
		}
	})

	t.Run("should return error when template not found", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return nil, nil
		}

		result, err := service.GetTemplateByIdForCacheFallback(context.Background(), "nonexistent-id")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

func TestUpdateManufacturerName(t *testing.T) {
	t.Run("should update manufacturer name for all templates", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.updateManyFn = func(_ context.Context, _ model.Map, _ model.Map) (int64, error) {
			return 5, nil
		}

		err := service.UpdateManufacturerName(context.Background(), "68f5bbce1aef22967c3ebb30", "New Manufacturer Name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should return error for invalid manufacturer id", func(t *testing.T) {
		service, _, _, _ := newTestService()

		err := service.UpdateManufacturerName(context.Background(), "invalid-id", "New Manufacturer Name")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUpdateModelName(t *testing.T) {
	t.Run("should update model name for all templates", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.updateManyFn = func(_ context.Context, _ model.Map, _ model.Map) (int64, error) {
			return 3, nil
		}

		err := service.UpdateModelName(context.Background(), "68f5bbce1aef22967c3ebb30", "New Model Name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUpdateCategoryName(t *testing.T) {
	t.Run("should update category name for all templates", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.updateManyFn = func(_ context.Context, _ model.Map, _ model.Map) (int64, error) {
			return 10, nil
		}

		err := service.UpdateCategoryName(context.Background(), "68f5bbce1aef22967c3ebb30", "New Category Name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGetAvailableFields(t *testing.T) {
	t.Run("should get available fields successfully", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)

		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return &entities.Assettemplate{
				ID:              templateObjectId,
				Name:            "Test Template",
				AvailableFields: []string{"temperature", "humidity", "pressure"},
			}, nil
		}

		result, err := service.GetAvailableFields(context.Background(), &templateId, &reqCtx.RequestContext{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		raw, ok := result["availableFields"]
		if !ok {
			t.Fatal("expected key 'availableFields' in result")
		}
		fields, ok := raw.([]string)
		if !ok {
			t.Fatalf("expected []string, got %T", raw)
		}
		if len(fields) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(fields))
		}
		if !containsString(fields, "temperature") {
			t.Fatal("expected 'temperature' in fields")
		}
	})

	t.Run("should return error when template not found", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return nil, nil
		}

		templateId := "nonexistent-id"
		result, err := service.GetAvailableFields(context.Background(), &templateId, &reqCtx.RequestContext{})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %#v", result)
		}
	})
}

/*
 * DYNAMIC FIELDS EVA PATTERN TESTS — pure helper tests, no port interactions.
 */

func TestProcessDynamicFieldsUpdate(t *testing.T) {
	t.Run("should preserve existing FieldIds for same field names", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
				{FieldId: 2, Field: "humidity", Value: "data.humidity", Type: "number", Status: 1},
			},
			NextFieldId: 3,
		}
		incoming := []dtos.DynamicField{
			{Field: "temperature", Value: "data.temp_updated", Type: "number"},
			{Field: "humidity", Value: "data.humidity", Type: "number"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(result.Fields))
		}
		if result.NextFieldId != 3 {
			t.Fatalf("expected NextFieldId=3, got %d", result.NextFieldId)
		}
		got := indexFields(result.Fields)
		assertField(t, got, "temperature", 1, 1, "data.temp_updated")
		assertField(t, got, "humidity", 2, 1, "data.humidity")
	})

	t.Run("should assign new FieldIds to new fields", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
			},
			NextFieldId: 2,
		}
		incoming := []dtos.DynamicField{
			{Field: "temperature", Value: "data.temp", Type: "number"},
			{Field: "pressure", Value: "data.pressure", Type: "number"},
			{Field: "windSpeed", Value: "data.windSpeed", Type: "number"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if result.NextFieldId != 4 {
			t.Fatalf("expected NextFieldId=4, got %d", result.NextFieldId)
		}
		got := indexFields(result.Fields)
		assertFieldId(t, got, "temperature", 1, 1)
		assertFieldId(t, got, "pressure", 2, 1)
		assertFieldId(t, got, "windSpeed", 3, 1)
	})

	t.Run("should mark removed fields as deprecated (Status=0)", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
				{FieldId: 2, Field: "humidity", Value: "data.humidity", Type: "number", Status: 1},
				{FieldId: 3, Field: "pressure", Value: "data.pressure", Type: "number", Status: 1},
			},
			NextFieldId: 4,
		}
		incoming := []dtos.DynamicField{
			{Field: "temperature", Value: "data.temp", Type: "number"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 3 {
			t.Fatalf("expected 3 fields (1 active + 2 deprecated), got %d", len(result.Fields))
		}
		if result.NextFieldId != 4 {
			t.Fatalf("expected NextFieldId=4, got %d", result.NextFieldId)
		}
		got := indexFields(result.Fields)
		assertFieldId(t, got, "temperature", 1, 1)
		assertFieldId(t, got, "humidity", 2, 0)
		assertFieldId(t, got, "pressure", 3, 0)
	})

	t.Run("should handle mixed scenario: add, keep, and remove fields", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
				{FieldId: 2, Field: "humidity", Value: "data.humidity", Type: "number", Status: 1},
				{FieldId: 3, Field: "oldField", Value: "data.old", Type: "string", Status: 1},
			},
			NextFieldId: 4,
		}
		incoming := []dtos.DynamicField{
			{Field: "temperature", Value: "data.temp", Type: "number"},
			{Field: "humidity", Value: "data.humidity_new", Type: "number"},
			{Field: "newField1", Value: "data.new1", Type: "string"},
			{Field: "newField2", Value: "data.new2", Type: "bool"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 5 {
			t.Fatalf("expected 5 fields, got %d", len(result.Fields))
		}
		if result.NextFieldId != 6 {
			t.Fatalf("expected NextFieldId=6, got %d", result.NextFieldId)
		}
		got := indexFields(result.Fields)
		assertField(t, got, "temperature", 1, 1, "data.temp")
		assertField(t, got, "humidity", 2, 1, "data.humidity_new")
		assertFieldId(t, got, "oldField", 3, 0)
		assertFieldId(t, got, "newField1", 4, 1)
		assertFieldId(t, got, "newField2", 5, 1)
	})

	t.Run("should handle empty existing fields", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{},
			NextFieldId:   0,
		}
		incoming := []dtos.DynamicField{
			{Field: "temperature", Value: "data.temp", Type: "number"},
			{Field: "humidity", Value: "data.humidity", Type: "number"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(result.Fields))
		}
		if result.NextFieldId != 3 {
			t.Fatalf("expected NextFieldId=3, got %d", result.NextFieldId)
		}
		got := indexFields(result.Fields)
		assertFieldId(t, got, "temperature", 1, 1)
		assertFieldId(t, got, "humidity", 2, 1)
	})

	t.Run("should handle removing all fields", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
				{FieldId: 2, Field: "humidity", Value: "data.humidity", Type: "number", Status: 1},
			},
			NextFieldId: 3,
		}
		incoming := []dtos.DynamicField{}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 2 {
			t.Fatalf("expected 2 deprecated fields, got %d", len(result.Fields))
		}
		if result.NextFieldId != 3 {
			t.Fatalf("expected NextFieldId=3 unchanged, got %d", result.NextFieldId)
		}
		for _, f := range result.Fields {
			if f.Status != 0 {
				t.Fatalf("expected field %s status=0, got %d", f.Field, f.Status)
			}
		}
	})

	t.Run("should not re-deprecate already deprecated fields", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
				{FieldId: 2, Field: "oldField", Value: "data.old", Type: "string", Status: 0},
			},
			NextFieldId: 3,
		}
		incoming := []dtos.DynamicField{
			{Field: "temperature", Value: "data.temp", Type: "number"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 1 {
			t.Fatalf("expected 1 field (already-deprecated must not reappear), got %d", len(result.Fields))
		}
		if result.Fields[0].Field != "temperature" || result.Fields[0].Status != 1 {
			t.Fatalf("expected active temperature, got %#v", result.Fields[0])
		}
	})

	t.Run("should preserve geo field paths", func(t *testing.T) {
		service, _, _, _ := newTestService()

		existing := &entities.Assettemplate{
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "location", Type: "geo", Status: 1, LatitudePath: "gps.lat", LongitudePath: "gps.lng"},
			},
			NextFieldId: 2,
		}
		incoming := []dtos.DynamicField{
			{Field: "location", Type: "geo", LatitudePath: "gps.latitude", LongitudePath: "gps.longitude"},
		}

		result := service.processDynamicFieldsUpdate(existing, incoming)

		if len(result.Fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(result.Fields))
		}
		got := result.Fields[0]
		if got.FieldId != 1 || got.LatitudePath != "gps.latitude" || got.LongitudePath != "gps.longitude" {
			t.Fatalf("expected preserved FieldId with updated paths, got %#v", got)
		}
	})
}

func TestCreateAssetTemplate_DynamicFields(t *testing.T) {
	t.Run("should assign sequential FieldIds starting from 1", func(t *testing.T) {
		service, repo, _, _ := newTestService()

		dto := &dtos.AssetTemplateCreateDTO{
			Name:             "Test Template",
			Enabled:          true,
			IsSystem:         true,
			ScriptValidator:  typeconv.PtrString("function validate() {}"),
			ScriptConversion: "function convert() {}",
			DynamicFields: []dtos.DynamicField{
				{Field: "temperature", Value: "data.temp", Type: "number"},
				{Field: "humidity", Value: "data.humidity", Type: "number"},
				{Field: "status", Value: "data.status", Type: "string"},
			},
		}

		var captured *entities.Assettemplate
		createdId, _ := model.ToObjectID("691bb4071e717d77a2430b47")
		repo.createFn = func(_ context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
			captured = e
			return &entities.Assettemplate{ID: createdId, Name: "Test Template", Enabled: true, IsSystem: true}, nil
		}

		_, _ = service.CreateAssetTemplate(context.Background(), &reqCtx.RequestContext{}, dto)

		if captured == nil {
			t.Fatal("expected Create to be called")
		}
		if len(captured.DynamicFields) != 3 {
			t.Fatalf("expected 3 dynamic fields, got %d", len(captured.DynamicFields))
		}
		expected := []uint16{1, 2, 3}
		for i, want := range expected {
			if captured.DynamicFields[i].FieldId != want {
				t.Fatalf("DynamicFields[%d].FieldId: want %d, got %d", i, want, captured.DynamicFields[i].FieldId)
			}
			if captured.DynamicFields[i].Status != 1 {
				t.Fatalf("DynamicFields[%d].Status (%s): want 1, got %d", i, captured.DynamicFields[i].Field, captured.DynamicFields[i].Status)
			}
		}
		if captured.NextFieldId != 4 {
			t.Fatalf("expected NextFieldId=4, got %d", captured.NextFieldId)
		}
	})

	t.Run("should handle empty DynamicFields", func(t *testing.T) {
		service, repo, _, _ := newTestService()

		dto := &dtos.AssetTemplateCreateDTO{
			Name:             "Test Template",
			Enabled:          true,
			IsSystem:         true,
			ScriptValidator:  typeconv.PtrString("function validate() {}"),
			ScriptConversion: "function convert() {}",
			DynamicFields:    []dtos.DynamicField{},
		}

		var captured *entities.Assettemplate
		createdId, _ := model.ToObjectID("691bb4071e717d77a2430b48")
		repo.createFn = func(_ context.Context, e *entities.Assettemplate) (*entities.Assettemplate, error) {
			captured = e
			return &entities.Assettemplate{ID: createdId, Name: "Test Template", Enabled: true, IsSystem: true}, nil
		}

		_, _ = service.CreateAssetTemplate(context.Background(), &reqCtx.RequestContext{}, dto)

		if captured == nil {
			t.Fatal("expected Create to be called")
		}
		if len(captured.DynamicFields) != 0 {
			t.Fatalf("expected 0 dynamic fields, got %d", len(captured.DynamicFields))
		}
		if captured.NextFieldId != 0 {
			t.Fatalf("expected NextFieldId=0 when no fields, got %d", captured.NextFieldId)
		}
	})
}

func TestUpdateAssetTemplateById_DynamicFields(t *testing.T) {
	t.Run("should preserve existing FieldIds and assign new ones during update", func(t *testing.T) {
		service, repo, _, _ := newTestService()

		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)

		existing := &entities.Assettemplate{
			ID:       templateObjectId,
			Name:     "Test Template",
			Enabled:  true,
			IsSystem: true,
			DynamicFields: []entities.DynamicField{
				{FieldId: 1, Field: "temperature", Value: "data.temp", Type: "number", Status: 1},
				{FieldId: 2, Field: "humidity", Value: "data.humidity", Type: "number", Status: 1},
			},
			NextFieldId: 3,
		}

		dto := &dtos.AssetTemplateUpdateDTO{
			DynamicFields: []dtos.DynamicField{
				{Field: "temperature", Value: "data.temp", Type: "number"},
				{Field: "pressure", Value: "data.pressure", Type: "number"},
			},
		}

		var capturedUpdate map[string]any
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return existing, nil
		}
		repo.findByIdAndUpdateFn = func(_ context.Context, _ *string, payload map[string]any) (*entities.Assettemplate, error) {
			capturedUpdate = payload
			return &entities.Assettemplate{ID: templateObjectId, Name: "Test Template", Enabled: true, IsSystem: true}, nil
		}

		_, _ = service.UpdateAssetTemplateById(context.Background(), &templateId, dto)

		if capturedUpdate == nil {
			t.Fatal("expected FindByIdAndUpdate to be called")
		}
		dynamicFields, ok := capturedUpdate["dynamicFields"].([]entities.DynamicField)
		if !ok {
			t.Fatalf("expected []entities.DynamicField in update payload, got %T", capturedUpdate["dynamicFields"])
		}
		if len(dynamicFields) != 3 {
			t.Fatalf("expected 3 fields (2 incoming + 1 deprecated), got %d", len(dynamicFields))
		}
		got := indexFields(dynamicFields)
		assertFieldId(t, got, "temperature", 1, 1)
		assertFieldId(t, got, "humidity", 2, 0)
		assertFieldId(t, got, "pressure", 3, 1)

		nextFieldId, ok := capturedUpdate["nextFieldId"].(uint16)
		if !ok {
			t.Fatalf("expected nextFieldId of type uint16 in payload, got %T", capturedUpdate["nextFieldId"])
		}
		if nextFieldId != 4 {
			t.Fatalf("expected NextFieldId=4, got %d", nextFieldId)
		}
	})
}

func TestFetchTemplateById(t *testing.T) {
	t.Run("should reuse fetchTemplateById in GetAssetTemplateById", func(t *testing.T) {
		service, repo, _, _ := newTestService()
		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)

		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return &entities.Assettemplate{ID: templateObjectId, Name: "Test Template"}, nil
		}

		result, err := service.GetAssetTemplateById(context.Background(), &templateId)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if repo.findByIdCalls != 1 {
			t.Fatalf("expected FindById to be called exactly once, got %d", repo.findByIdCalls)
		}
	})

	t.Run("should reuse fetchTemplateById in GetTemplateByIdForCacheFallback", func(t *testing.T) {
		service, repo, storage, _ := newTestService()
		templateId := "691bb4071e717d77a2430b46"
		templateObjectId, _ := model.ToObjectID(templateId)

		existing := &entities.Assettemplate{ID: templateObjectId, Name: "Test Template"}
		repo.findByIdFn = func(_ context.Context, _ *string) (*entities.Assettemplate, error) {
			return existing, nil
		}
		storage.writeScriptsFn = func(_ context.Context, _ *entities.Assettemplate) error { return nil }

		result, err := service.GetTemplateByIdForCacheFallback(context.Background(), templateId)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if repo.findByIdCalls != 1 {
			t.Fatalf("expected FindById to be called exactly once, got %d", repo.findByIdCalls)
		}
	})
}

/*
 * TEST HELPERS
 */

func indexFields(fields []entities.DynamicField) map[string]entities.DynamicField {
	out := make(map[string]entities.DynamicField, len(fields))
	for _, f := range fields {
		out[f.Field] = f
	}
	return out
}

func assertFieldId(t *testing.T, fields map[string]entities.DynamicField, name string, fieldId uint16, status uint8) {
	t.Helper()
	got, ok := fields[name]
	if !ok {
		t.Fatalf("expected field %q to be present", name)
	}
	if got.FieldId != fieldId {
		t.Fatalf("field %q: want FieldId=%d, got %d", name, fieldId, got.FieldId)
	}
	if got.Status != status {
		t.Fatalf("field %q: want Status=%d, got %d", name, status, got.Status)
	}
}

func assertField(t *testing.T, fields map[string]entities.DynamicField, name string, fieldId uint16, status uint8, value string) {
	t.Helper()
	assertFieldId(t, fields, name, fieldId, status)
	if fields[name].Value != value {
		t.Fatalf("field %q: want Value=%q, got %q", name, value, fields[name].Value)
	}
}

func containsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
