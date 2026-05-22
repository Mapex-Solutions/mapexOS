package services

import (
	"encoding/json"
	"errors"
	"testing"

	"router/src/bootstrap"
	"router/src/modules/events/application/di"
	templateCache "router/src/modules/events/infrastructure/cache/tieredcache"
	rgMocks "router/src/modules/routegroups/application/services/mocks"
	sharedMocks "router/src/shared/mocks"

	"github.com/prometheus/client_golang/prometheus"
	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	routegroups "github.com/Mapex-Solutions/MapexOS/contracts/services/router/routegroups"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
)

// createTestMetrics creates a real RouterMetrics instance for testing.
// Prometheus metrics work fine in tests — no mocking needed.
func createTestMetrics() *bootstrap.RouterMetrics {
	reg := metrics.NewRegistry("router_test")

	return &bootstrap.RouterMetrics{
		Registry: reg,

		EventsProcessed: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "event", Name: "processed_total", Help: "test",
		}, []string{"status"}),
		EventProcessingDuration: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "event", Name: "processing_duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}),
		EventsBatchSize: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "event", Name: "batch_size", Help: "test",
			Buckets: prometheus.DefBuckets,
		}),
		MessagesTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "message", Name: "total", Help: "test",
		}, []string{"result"}),
		AssetCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "asset", Name: "cache_total", Help: "test",
		}, []string{"tier"}),
		AssetCacheDuration: reg.NewHistogramVec(metrics.HistogramOpts{
			Subsystem: "asset", Name: "cache_duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}, []string{"tier"}),
		CacheInvalidationsTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "cache", Name: "invalidations_total", Help: "test",
		}, []string{"status"}),
		MatchEvaluationsTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "match", Name: "evaluations_total", Help: "test",
		}, []string{"result"}),
		MatchRulesEvaluatedTotal: reg.NewCounter(metrics.CounterOpts{
			Subsystem: "match", Name: "rules_evaluated_total", Help: "test",
		}),
		EventsPublished: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "event", Name: "published_total", Help: "test",
		}, []string{"kind", "status"}),
		PublishDuration: reg.NewHistogramVec(metrics.HistogramOpts{
			Subsystem: "publish", Name: "duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}, []string{"kind"}),
		RouteGroupOperations: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "routegroup", Name: "operations_total", Help: "test",
		}, []string{"operation", "status"}),
		RouteGroupOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
			Subsystem: "routegroup", Name: "operation_duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}, []string{"operation"}),
		RouteGroupListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "routegroup", Name: "list_results_count", Help: "test",
			Buckets: []float64{0, 1, 5, 10, 25, 50, 100, 250},
		}),
		RouteGroupCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "routegroup", Name: "cache_total", Help: "test",
		}, []string{"result"}),
	}
}

// createTestDependencies creates mock dependencies for testing.
func createTestDependencies() (*sharedMocks.MockTieredCache, *rgMocks.MockRouteGroupService, *sharedMocks.MockNatsBus, di.EventServiceDependenciesInjection) {
	assetCache := sharedMocks.NewMockTieredCache()
	templateCacheRaw := sharedMocks.NewMockTieredCache()
	routeGroupService := rgMocks.NewMockRouteGroupService()
	natsBus := sharedMocks.NewMockNatsBus()
	testMetrics := createTestMetrics()

	deps := di.EventServiceDependenciesInjection{
		AssetCache:        assetCache,
		TemplateCache:     templateCache.New(templateCacheRaw),
		RouteGroupService: routeGroupService,
		NatsBus:           natsBus,
		Metrics:           testMetrics,
	}

	return assetCache, routeGroupService, natsBus, deps
}

// createTestAsset creates a mock asset for testing.
func createTestAsset(orgId, uuid string) *assetsContract.AssetReadModel {
	return &assetsContract.AssetReadModel{
		ID:            "asset-mongo-id",
		UUID:          uuid,
		OrgId:         orgId,
		PathKey:       "org1/suborg1",
		Enabled:       true,
		Name:          "Test Asset",
		Description:   "Test Description",
		RouteGroupIds: []string{"routegroup-1"},
	}
}

// createTestRouteGroup creates a mock route group for testing.
func createTestRouteGroup(id string, kind string, hasMatch bool) *routegroups.RouteGroupResponse {
	name := "Test Route Group"
	description := "Test Description"
	objectId := common.ObjectID{}

	routers := []routegroups.Router{
		{
			Kind: kind,
		},
	}

	if hasMatch {
		policy := "all"
		rules := []routegroups.MatchRule{
			{
				Field:    "temperature",
				Operator: "gt",
				Value:    float64(25),
			},
		}
		routers[0].Match = &routegroups.MatchConfig{
			Policy: policy,
			Rules:  &rules,
		}
	}

	return &routegroups.RouteGroupResponse{
		ID:          &objectId,
		Name:        &name,
		Description: &description,
		Routers:     &routers,
	}
}

/**
 * TEST: ProcessEvent (Legacy V1)
 */

func TestProcessEvent_Success(t *testing.T) {
	assetCache, routeGroupService, natsBus, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	routeGroup := createTestRouteGroup("routegroup-1", "save_event", false)
	routeGroupService.GetByIdResponse = routeGroup

	// Create service and process event
	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify publish was called (save_event + history)
	if len(natsBus.PublishCalls) < 1 {
		t.Errorf("Expected at least 1 publish call, got: %d", len(natsBus.PublishCalls))
	}
}

func TestProcessEvent_MissingOrgId(t *testing.T) {
	_, _, _, deps := createTestDependencies()
	service := New(deps)

	eventData := map[string]interface{}{
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err == nil {
		t.Error("Expected error for missing orgId")
	}
}

func TestProcessEvent_MissingAssetId(t *testing.T) {
	_, _, _, deps := createTestDependencies()
	service := New(deps)

	eventData := map[string]interface{}{
		"orgId": "org-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err == nil {
		t.Error("Expected error for missing assetId")
	}
}

func TestProcessEvent_MissingEvent(t *testing.T) {
	_, _, _, deps := createTestDependencies()
	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err == nil {
		t.Error("Expected error for missing event")
	}
}

func TestProcessEvent_InvalidJSON(t *testing.T) {
	_, _, _, deps := createTestDependencies()
	service := New(deps)

	err := service.ProcessEvent([]byte("invalid json"), 0, nil)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestProcessEvent_AssetNotFound(t *testing.T) {
	assetCache, _, _, deps := createTestDependencies()
	assetCache.GetError = errors.New("asset not found")

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err == nil {
		t.Error("Expected error for asset not found")
	}
}

/**
 * TEST: ProcessEventBatch
 */

func TestProcessEventBatch_EmptyMessages(t *testing.T) {
	_, _, _, deps := createTestDependencies()
	service := New(deps)

	err := service.ProcessEventBatch([]*natsModel.Message{})

	if err != nil {
		t.Errorf("Expected no error for empty batch, got: %v", err)
	}
}

/**
 * TEST: ProcessAssetInvalidateBatch
 */

func TestProcessAssetInvalidateBatch_Success(t *testing.T) {
	assetCache, _, _, deps := createTestDependencies()

	// Pre-populate cache
	assetCache.Data["org-1/asset-uuid-1"] = []byte("{}")

	service := New(deps)

	invalidateData := map[string]interface{}{
		"orgId":     "org-1",
		"assetUUID": "asset-uuid-1",
	}
	invalidateJSON, _ := json.Marshal(invalidateData)

	msg := &natsModel.Message{
		Data: invalidateJSON,
	}

	service.ProcessAssetInvalidateBatch([]*natsModel.Message{msg})

	// Verify cache invalidation was called
	if len(assetCache.InvalidateCalls) != 1 {
		t.Errorf("Expected 1 invalidate call, got: %d", len(assetCache.InvalidateCalls))
	}

	if assetCache.InvalidateCalls[0] != "org-1/asset-uuid-1" {
		t.Errorf("Expected invalidate key 'org-1/asset-uuid-1', got: %s", assetCache.InvalidateCalls[0])
	}
}

func TestProcessAssetInvalidateBatch_EmptyMessages(t *testing.T) {
	assetCache, _, _, deps := createTestDependencies()
	service := New(deps)

	service.ProcessAssetInvalidateBatch([]*natsModel.Message{})

	// Verify no invalidation was called
	if len(assetCache.InvalidateCalls) != 0 {
		t.Errorf("Expected 0 invalidate calls, got: %d", len(assetCache.InvalidateCalls))
	}
}

/**
 * TEST: Match Evaluator Integration
 */

func TestProcessEvent_MatchConditionPasses(t *testing.T) {
	assetCache, routeGroupService, natsBus, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	// RouteGroup with match condition: temperature > 25
	routeGroup := createTestRouteGroup("routegroup-1", "save_event", true)
	routeGroupService.GetByIdResponse = routeGroup

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0, // > 25, should match
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify publish was called (event matched)
	// Should have at least 2 calls: save_event publish + history publish
	if len(natsBus.PublishCalls) < 2 {
		t.Errorf("Expected at least 2 publish calls (event + history), got: %d", len(natsBus.PublishCalls))
	}
}

func TestProcessEvent_MatchConditionFails(t *testing.T) {
	assetCache, routeGroupService, natsBus, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	// RouteGroup with match condition: temperature > 25
	routeGroup := createTestRouteGroup("routegroup-1", "save_event", true)
	routeGroupService.GetByIdResponse = routeGroup

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 20.0, // < 25, should NOT match
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify only history was published (event did NOT match)
	// Should have only 1 call: history publish
	if len(natsBus.PublishCalls) != 1 {
		t.Errorf("Expected 1 publish call (history only), got: %d", len(natsBus.PublishCalls))
	}
}

func TestProcessEvent_NoMatchConfig(t *testing.T) {
	assetCache, routeGroupService, natsBus, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	// RouteGroup WITHOUT match config (should always match)
	routeGroup := createTestRouteGroup("routegroup-1", "save_event", false)
	routeGroupService.GetByIdResponse = routeGroup

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 10.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify publish was called (no match = always allowed)
	if len(natsBus.PublishCalls) < 2 {
		t.Errorf("Expected at least 2 publish calls, got: %d", len(natsBus.PublishCalls))
	}
}

/**
 * TEST: RouteGroup Not Found
 */

func TestProcessEvent_RouteGroupNotFound(t *testing.T) {
	assetCache, routeGroupService, _, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	// RouteGroup not found
	routeGroupService.GetByIdError = errors.New("route group not found")

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	// Should not return error, just log warning
	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error (warning only), got: %v", err)
	}
}

/**
 * TEST: Publish Error
 */

func TestProcessEvent_PublishError(t *testing.T) {
	assetCache, routeGroupService, natsBus, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	routeGroup := createTestRouteGroup("routegroup-1", "save_event", false)
	routeGroupService.GetByIdResponse = routeGroup

	// Simulate publish error
	natsBus.PublishError = errors.New("publish failed")

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	// Should not return error, just log error
	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error (error is logged), got: %v", err)
	}
}

/**
 * TEST: Empty Routers
 */

func TestProcessEvent_EmptyRouters(t *testing.T) {
	assetCache, routeGroupService, natsBus, deps := createTestDependencies()

	// Setup mock data
	asset := createTestAsset("org-1", "asset-uuid-1")
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	// RouteGroup with no routers
	name := "Empty Route Group"
	routeGroup := &routegroups.RouteGroupResponse{
		Name:    &name,
		Routers: nil,
	}
	routeGroupService.GetByIdResponse = routeGroup

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify no publish was called (no routers)
	if len(natsBus.PublishCalls) != 0 {
		t.Errorf("Expected 0 publish calls, got: %d", len(natsBus.PublishCalls))
	}
}

/**
 * TEST: Multiple RouteGroups
 */

func TestProcessEvent_MultipleRouteGroups(t *testing.T) {
	assetCache, routeGroupService, _, deps := createTestDependencies()

	// Asset with multiple route groups
	asset := &assetsContract.AssetReadModel{
		ID:            "asset-mongo-id",
		UUID:          "asset-uuid-1",
		OrgId:         "org-1",
		PathKey:       "org1/suborg1",
		Enabled:       true,
		Name:          "Test Asset",
		RouteGroupIds: []string{"routegroup-1", "routegroup-2"},
	}
	assetJSON, _ := json.Marshal(asset)
	assetCache.Data["org-1/asset-uuid-1"] = assetJSON

	routeGroup := createTestRouteGroup("routegroup-1", "save_event", false)
	routeGroupService.GetByIdResponse = routeGroup

	service := New(deps)

	eventData := map[string]interface{}{
		"orgId":   "org-1",
		"assetId": "asset-uuid-1",
		"event": map[string]interface{}{
			"temperature": 30.0,
		},
	}
	eventJSON, _ := json.Marshal(eventData)

	err := service.ProcessEvent(eventJSON, 0, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify GetRouteGroupById was called twice (for both route groups)
	if len(routeGroupService.GetByIdCalls) != 2 {
		t.Errorf("Expected 2 GetRouteGroupById calls, got: %d", len(routeGroupService.GetByIdCalls))
	}
}
