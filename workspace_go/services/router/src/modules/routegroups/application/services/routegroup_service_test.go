package services

import (
	ctx "context"
	"errors"
	"testing"
	"time"

	"router/src/bootstrap"
	"router/src/modules/routegroups/application/di"
	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/domain/entities"
	"router/src/modules/routegroups/application/services/mocks"
	sharedMocks "router/src/shared/mocks"

	"github.com/prometheus/client_golang/prometheus"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
)

// createTestMetrics creates a real RouterMetrics instance for testing.
func createTestMetrics() *bootstrap.RouterMetrics {
	reg := metrics.NewRegistry("routegroup_test")

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

// createServiceDependencies creates mock dependencies for testing RouteGroupService.
func createServiceDependencies() (*mocks.MockRouteGroupRepository, *mocks.MockCacheRepository, *sharedMocks.MockAppCache, di.RouteGroupServiceDependenciesInjection) {
	repo := mocks.NewMockRouteGroupRepository()
	cache := mocks.NewMockCacheRepository()
	appCache := sharedMocks.NewMockAppCache()

	deps := di.RouteGroupServiceDependenciesInjection{
		RouteGroupRepo: repo,
		CacheRepo:      cache,
		AppCache:       appCache,
		Metrics:        createTestMetrics(),
	}

	return repo, cache, appCache, deps
}

// createTestRouteGroupEntity creates a RouteGroup entity for testing.
func createTestRouteGroupEntity(hexId string) *entities.RouteGroup {
	objectId, _ := model.ToObjectID(hexId)
	return &entities.RouteGroup{
		ID:          objectId,
		Version:     "1.0.0",
		Name:        "Test Route Group",
		Description: "A test route group",
		Enabled:     true,
		IsSystem:    false,
		IsTemplate:  false,
		PathKey:     "000001/0001",
		Routers:     []entities.Router{},
		Created:     time.Now(),
		Updated:     time.Now(),
	}
}

// createTestRequestContext creates a RequestContext for testing.
func createTestRequestContext(orgId, pathKey string) *reqCtx.RequestContext {
	return &reqCtx.RequestContext{
		ScopedOrgIds: []string{orgId},
		OrgContext:   &orgId,
		OrgContextData: &reqCtx.CoverageOrg{
			ID:      orgId,
			Name:    "Test Org",
			Type:    "customer",
			PathKey: pathKey,
		},
		UserId: "user-1",
	}
}

/**
 * TEST: CreateRouteGroup
 */

func TestCreateRouteGroup_SystemResource(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	entity.IsSystem = true
	repo.CreateResponse = entity

	service := New(deps)

	isSystem := true
	reqContext := createTestRequestContext("org-1", "000001/0001")

	result, err := service.CreateRouteGroup(ctx.Background(), reqContext, &dtos.RouteGroupCreateDTO{
		Version:  "1.0.0",
		Name:     "System Route Group",
		IsSystem: &isSystem,
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Verify repo was called
	if len(repo.CreateCalls) != 1 {
		t.Errorf("Expected 1 create call, got: %d", len(repo.CreateCalls))
	}
}

func TestCreateRouteGroup_LocalResource(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	repo.CreateResponse = entity

	service := New(deps)
	reqContext := createTestRequestContext("507f1f77bcf86cd799439022", "000001/0001")

	result, err := service.CreateRouteGroup(ctx.Background(), reqContext, &dtos.RouteGroupCreateDTO{
		Version: "1.0.0",
		Name:    "Local Route Group",
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(repo.CreateCalls) != 1 {
		t.Errorf("Expected 1 create call, got: %d", len(repo.CreateCalls))
	}
}

func TestCreateRouteGroup_RepoError(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	repo.CreateError = errors.New("database error")

	service := New(deps)
	reqContext := createTestRequestContext("507f1f77bcf86cd799439022", "000001/0001")

	isSystem := true
	_, err := service.CreateRouteGroup(ctx.Background(), reqContext, &dtos.RouteGroupCreateDTO{
		Version:  "1.0.0",
		Name:     "Failed Route Group",
		IsSystem: &isSystem,
	})

	if err == nil {
		t.Error("Expected error for repo failure")
	}
}

/**
 * TEST: GetRouteGroupById
 */

func TestGetRouteGroupById_CacheMiss_DbHit(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	repo.FindByIdResponse = entity

	service := New(deps)

	result, err := service.GetRouteGroupById(ctx.Background(), &entityId)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetRouteGroupById_NotFound(t *testing.T) {
	repo, cache, _, deps := createServiceDependencies()

	// Return nil entity (not found in DB)
	repo.FindByIdResponse = nil
	cache.GetOrSetExError = errors.New("not found")

	service := New(deps)

	id := "507f1f77bcf86cd799439099"
	_, err := service.GetRouteGroupById(ctx.Background(), &id)

	if err == nil {
		t.Error("Expected error for not found route group")
	}
}

func TestGetRouteGroupById_CacheCallbackError(t *testing.T) {
	_, cache, _, deps := createServiceDependencies()

	cache.GetOrSetExError = errors.New("cache callback error")

	service := New(deps)

	id := "507f1f77bcf86cd799439011"
	_, err := service.GetRouteGroupById(ctx.Background(), &id)

	if err == nil {
		t.Error("Expected error when cache callback fails")
	}
}

/**
 * TEST: UpdateRouteGroupById
 */

func TestUpdateRouteGroupById_Success(t *testing.T) {
	repo, cache, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	repo.FindByIdAndUpdateResponse = entity

	service := New(deps)

	newName := "Updated Name"
	result, err := service.UpdateRouteGroupById(ctx.Background(), &entityId, &dtos.RouteGroupUpdateDTO{
		Name: &newName,
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Verify repo was called
	if len(repo.FindByIdAndUpdateCalls) != 1 {
		t.Errorf("Expected 1 FindByIdAndUpdate call, got: %d", len(repo.FindByIdAndUpdateCalls))
	}

	// Verify cache was updated
	if len(cache.SetExCalls) != 1 {
		t.Errorf("Expected 1 SetEx call, got: %d", len(cache.SetExCalls))
	}
}

func TestUpdateRouteGroupById_NotFound(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	// Return entity with zero ID (not found)
	repo.FindByIdAndUpdateResponse = &entities.RouteGroup{}

	service := New(deps)

	id := "507f1f77bcf86cd799439099"
	newName := "Updated Name"
	_, err := service.UpdateRouteGroupById(ctx.Background(), &id, &dtos.RouteGroupUpdateDTO{
		Name: &newName,
	})

	if err == nil {
		t.Error("Expected error for not found route group")
	}
}

/**
 * TEST: DeleteRouteGroupById
 */

func TestDeleteRouteGroupById_Success(t *testing.T) {
	repo, cache, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	repo.FindByIdResponse = entity

	service := New(deps)

	result, err := service.DeleteRouteGroupById(ctx.Background(), &entityId)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result["success"] {
		t.Error("Expected success: true")
	}

	// Verify FindById was called (pre-check)
	if len(repo.FindByIdCalls) != 1 {
		t.Errorf("Expected 1 FindById call, got: %d", len(repo.FindByIdCalls))
	}

	// Verify DeleteById was called
	if len(repo.DeleteByIdCalls) != 1 {
		t.Errorf("Expected 1 DeleteById call, got: %d", len(repo.DeleteByIdCalls))
	}

	// Verify cache was deleted
	if len(cache.DelCalls) != 1 {
		t.Errorf("Expected 1 Del call, got: %d", len(cache.DelCalls))
	}
}

func TestDeleteRouteGroupById_NotFound(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	// FindById returns nil (not found)
	repo.FindByIdResponse = nil
	repo.FindByIdError = errors.New("document not found")

	service := New(deps)

	id := "507f1f77bcf86cd799439099"
	_, err := service.DeleteRouteGroupById(ctx.Background(), &id)

	if err == nil {
		t.Error("Expected error for not found route group")
	}
}

/**
 * TEST: GetRouteGroupsByIds
 */

func TestGetRouteGroupsByIds_AllFound(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	repo.FindByIdResponse = entity

	service := New(deps)

	ids := []string{entityId, entityId}
	results, err := service.GetRouteGroupsByIds(ctx.Background(), ids)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got: %d", len(results))
	}
}

func TestGetRouteGroupsByIds_PartialNotFound(t *testing.T) {
	_, cache, _, deps := createServiceDependencies()

	// All cache misses return error (simulating not found)
	cache.GetOrSetExError = errors.New("not found")

	service := New(deps)

	ids := []string{"507f1f77bcf86cd799439011", "507f1f77bcf86cd799439099"}
	results, err := service.GetRouteGroupsByIds(ctx.Background(), ids)

	if err != nil {
		t.Errorf("Expected no error (skip not found), got: %v", err)
	}

	// All should be skipped since cache returns error
	if len(results) != 0 {
		t.Errorf("Expected 0 results (all not found), got: %d", len(results))
	}
}

func TestGetRouteGroupsByIds_EmptyIds(t *testing.T) {
	_, _, _, deps := createServiceDependencies()

	service := New(deps)

	results, err := service.GetRouteGroupsByIds(ctx.Background(), []string{})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty ids, got: %d", len(results))
	}
}

/**
 * TEST: GetRouteGroups
 */

func TestGetRouteGroups_BasicQuery(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)

	repo.FindWithFiltersResponse = &model.PaginatedResult[entities.RouteGroup]{
		Items: []entities.RouteGroup{*entity},
		Pagination: model.Pagination{
			Page:       1,
			PerPage:    20,
			TotalItems: 1,
			TotalPages: 1,
		},
	}

	service := New(deps)
	reqContext := createTestRequestContext("org-1", "000001/0001")

	result, err := service.GetRouteGroups(ctx.Background(), reqContext, &dtos.RouteGroupQueryDTO{})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got: %d", len(result.Items))
	}

	// Verify FindWithFilters was called
	if len(repo.FindWithFiltersCalls) != 1 {
		t.Errorf("Expected 1 FindWithFilters call, got: %d", len(repo.FindWithFiltersCalls))
	}
}

func TestGetRouteGroups_WithNameFilter(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	repo.FindWithFiltersResponse = &model.PaginatedResult[entities.RouteGroup]{
		Items: []entities.RouteGroup{},
		Pagination: model.Pagination{
			Page:       1,
			PerPage:    20,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	service := New(deps)
	reqContext := createTestRequestContext("org-1", "000001/0001")

	nameFilter := "test"
	result, err := service.GetRouteGroups(ctx.Background(), reqContext, &dtos.RouteGroupQueryDTO{
		Name: &nameFilter,
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Verify filters contain name regex
	if len(repo.FindWithFiltersCalls) != 1 {
		t.Errorf("Expected 1 FindWithFilters call, got: %d", len(repo.FindWithFiltersCalls))
	}

	filters := repo.FindWithFiltersCalls[0].Filters
	nameFilterMap, ok := filters["name"].(model.Map)
	if !ok {
		t.Error("Expected name filter to be a Map with $regex")
	} else {
		if nameFilterMap["$regex"] != "test" {
			t.Errorf("Expected name regex 'test', got: %v", nameFilterMap["$regex"])
		}
	}
}

func TestGetRouteGroups_RepoError(t *testing.T) {
	repo, _, _, deps := createServiceDependencies()

	repo.FindWithFiltersError = errors.New("database error")

	service := New(deps)
	reqContext := createTestRequestContext("org-1", "000001/0001")

	_, err := service.GetRouteGroups(ctx.Background(), reqContext, &dtos.RouteGroupQueryDTO{})

	if err == nil {
		t.Error("Expected error for repo failure")
	}
}

/**
 * TEST: CountRouteGroups
 */

func TestCountRouteGroups_CacheHit(t *testing.T) {
	_, _, appCache, deps := createServiceDependencies()

	// Pre-populate cache with count value
	orgId := "org-1"
	cacheKey := "counter:route_groups:" + orgId
	appCache.Store[cacheKey] = int64(42)

	service := New(deps)
	reqContext := createTestRequestContext(orgId, "000001/0001")

	count, err := service.CountRouteGroups(ctx.Background(), reqContext)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if count != 42 {
		t.Errorf("Expected count 42, got: %d", count)
	}

	// Verify Get was called
	if len(appCache.GetCalls) != 1 {
		t.Errorf("Expected 1 Get call, got: %d", len(appCache.GetCalls))
	}
}

func TestCountRouteGroups_CacheMiss_DbSuccess(t *testing.T) {
	repo, _, appCache, deps := createServiceDependencies()

	// Cache miss (store is empty) + DB returns count
	appCache.GetError = errors.New("cache miss")
	repo.CountDocumentsResponse = 15

	service := New(deps)
	reqContext := createTestRequestContext("507f1f77bcf86cd799439022", "000001/0001")

	count, err := service.CountRouteGroups(ctx.Background(), reqContext)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if count != 15 {
		t.Errorf("Expected count 15, got: %d", count)
	}

	// Verify CountDocuments was called
	if len(repo.CountDocumentsCalls) != 1 {
		t.Errorf("Expected 1 CountDocuments call, got: %d", len(repo.CountDocumentsCalls))
	}

	// Verify SetEx was called to cache the result
	if len(appCache.SetExCalls) != 1 {
		t.Errorf("Expected 1 SetEx call, got: %d", len(appCache.SetExCalls))
	}
}

func TestCountRouteGroups_CacheMiss_DbError(t *testing.T) {
	repo, _, appCache, deps := createServiceDependencies()

	appCache.GetError = errors.New("cache miss")
	repo.CountDocumentsError = errors.New("database error")

	service := New(deps)
	reqContext := createTestRequestContext("org-1", "000001/0001")

	_, err := service.CountRouteGroups(ctx.Background(), reqContext)

	if err == nil {
		t.Error("Expected error for DB failure")
	}

	// Verify SetEx was NOT called (DB failed)
	if len(appCache.SetExCalls) != 0 {
		t.Errorf("Expected 0 SetEx calls, got: %d", len(appCache.SetExCalls))
	}
}

func TestCountRouteGroups_NilOrgContext(t *testing.T) {
	repo, _, appCache, deps := createServiceDependencies()

	appCache.GetError = errors.New("cache miss")
	repo.CountDocumentsResponse = 7

	service := New(deps)

	// Create request context with nil OrgContext
	reqContext := &reqCtx.RequestContext{
		ScopedOrgIds: []string{},
		OrgContext:   nil,
		UserId:       "user-1",
	}

	count, err := service.CountRouteGroups(ctx.Background(), reqContext)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if count != 7 {
		t.Errorf("Expected count 7, got: %d", count)
	}

	// Verify cache key used empty orgId
	if len(appCache.GetCalls) != 1 {
		t.Errorf("Expected 1 Get call, got: %d", len(appCache.GetCalls))
	}

	expectedKey := "counter:route_groups:"
	if appCache.GetCalls[0] != expectedKey {
		t.Errorf("Expected cache key %q, got: %q", expectedKey, appCache.GetCalls[0])
	}
}

func TestCountRouteGroups_CacheInvalidation_OnCreate(t *testing.T) {
	repo, _, appCache, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)
	entity.IsSystem = true
	repo.CreateResponse = entity

	service := New(deps)

	isSystem := true
	reqContext := createTestRequestContext("org-1", "000001/0001")

	_, err := service.CreateRouteGroup(ctx.Background(), reqContext, &dtos.RouteGroupCreateDTO{
		Version:  "1.0.0",
		Name:     "System Route Group",
		IsSystem: &isSystem,
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify AppCache.Del was called with counter key
	expectedCounterKey := "counter:route_groups:org-1"
	found := false
	for _, key := range appCache.DelCalls {
		if key == expectedCounterKey {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected AppCache.Del called with key %q, got calls: %v", expectedCounterKey, appCache.DelCalls)
	}
}

func TestCountRouteGroups_CacheInvalidation_OnDelete(t *testing.T) {
	repo, _, appCache, deps := createServiceDependencies()

	entityId := "507f1f77bcf86cd799439011"
	entity := createTestRouteGroupEntity(entityId)

	// Set OrgId on the entity so deletion can invalidate counter cache
	orgHex := "507f1f77bcf86cd799439022"
	orgObjectId, _ := model.ToObjectID(orgHex)
	entity.OrgId = &orgObjectId
	repo.FindByIdResponse = entity

	service := New(deps)

	result, err := service.DeleteRouteGroupById(ctx.Background(), &entityId)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil || !result["success"] {
		t.Error("Expected success: true")
	}

	// Verify AppCache.Del was called with counter key for the entity's orgId
	expectedCounterKey := "counter:route_groups:" + orgHex
	found := false
	for _, key := range appCache.DelCalls {
		if key == expectedCounterKey {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected AppCache.Del called with key %q, got calls: %v", expectedCounterKey, appCache.DelCalls)
	}
}

/**
 * TEST: GetRouteGroups — Kinds filter (strict semantic)
 *
 * Asserts the new `Kinds` filter wired in routegroup_handler_crud.go (T2):
 *   - When len(query.Kinds) > 0, the emitted filter map MUST contain:
 *       filters["routers.0"] = model.Map{"$exists": true}
 *       filters["routers"]   = model.Map{"$not": model.Map{"$elemMatch":
 *                              model.Map{"kind": model.Map{"$nin": query.Kinds}}}}
 *   - When Kinds is empty/nil, neither key is emitted (backward-compat).
 *
 * Strict semantic outcomes (A–F) are verified via a pure simulator that
 * mirrors the MongoDB matcher (`routers.0 $exists` AND no router with
 * kind $nin Kinds). This avoids spinning up a real Mongo container while
 * still asserting the inclusion/exclusion behavior of the 5 fixtures.
 */

// simulateKindsFilter mirrors the MongoDB strict-Kinds filter semantic:
//
//	(routers.0 $exists) AND NOT $elemMatch{kind $nin Kinds}
//
// i.e. there is at least one router AND every router.kind is in Kinds.
// When Kinds is empty/nil, the filter is not applied → all docs match.
func simulateKindsFilter(routers []entities.Router, kinds []string) bool {
	if len(kinds) == 0 {
		return true // backward-compat: no filter applied
	}
	if len(routers) == 0 {
		return false // routers.0 $exists guard
	}
	allowed := make(map[string]struct{}, len(kinds))
	for _, k := range kinds {
		allowed[k] = struct{}{}
	}
	for _, r := range routers {
		if _, ok := allowed[r.Kind]; !ok {
			return false // $elemMatch{kind $nin Kinds} would match → group excluded
		}
	}
	return true
}

func TestGetRouteGroups_KindsFilter_StrictSemantic(t *testing.T) {
	// 5 fixtures covering inclusion/exclusion under Kinds=["trigger","workflow"].
	type fixture struct {
		label   string
		routers []entities.Router
	}
	fixtures := []fixture{
		// A — INCLUDED: all kinds inside the set.
		{"A_all_in_set", []entities.Router{{Kind: "trigger"}, {Kind: "workflow"}}},
		// B — EXCLUDED: mixed (one outside).
		{"B_mixed_one_outside", []entities.Router{{Kind: "trigger"}, {Kind: "save_event"}}},
		// C — INCLUDED: single kind in set.
		{"C_single_in_set", []entities.Router{{Kind: "trigger"}}},
		// D — EXCLUDED: single kind outside set.
		{"D_single_outside_set", []entities.Router{{Kind: "save_event"}}},
		// E — EXCLUDED: empty routers (routers.0 guard).
		{"E_empty_routers", []entities.Router{}},
	}

	kinds := []string{"trigger", "workflow"}
	expectedIncluded := map[string]bool{
		"A_all_in_set":         true,
		"B_mixed_one_outside":  false,
		"C_single_in_set":      true,
		"D_single_outside_set": false,
		"E_empty_routers":      false,
	}

	for _, f := range fixtures {
		got := simulateKindsFilter(f.routers, kinds)
		want := expectedIncluded[f.label]
		if got != want {
			t.Errorf("[%s] strict-Kinds match: got=%v, want=%v (routers=%v, kinds=%v)",
				f.label, got, want, f.routers, kinds)
		}
	}

	// F — Backward-compat: Kinds empty/nil ⇒ all 5 fixtures match.
	for _, f := range fixtures {
		if !simulateKindsFilter(f.routers, nil) {
			t.Errorf("[F_backcompat_nil] expected fixture %q to match when Kinds=nil", f.label)
		}
		if !simulateKindsFilter(f.routers, []string{}) {
			t.Errorf("[F_backcompat_empty] expected fixture %q to match when Kinds=[]", f.label)
		}
	}

	// Now assert the SERVICE actually emits the expected filter map shape
	// when Kinds is provided (Option A — filter map assertion via the repo mock).
	repo, _, _, deps := createServiceDependencies()
	repo.FindWithFiltersResponse = &model.PaginatedResult[entities.RouteGroup]{
		Items: []entities.RouteGroup{},
		Pagination: model.Pagination{
			Page:       1,
			PerPage:    20,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	service := New(deps)
	reqContext := createTestRequestContext("org-1", "000001/0001")

	_, err := service.GetRouteGroups(ctx.Background(), reqContext, &dtos.RouteGroupQueryDTO{
		Kinds: kinds,
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(repo.FindWithFiltersCalls) != 1 {
		t.Fatalf("Expected 1 FindWithFilters call, got: %d", len(repo.FindWithFiltersCalls))
	}
	filters := repo.FindWithFiltersCalls[0].Filters

	// 1) routers.0 $exists guard.
	routers0, ok := filters["routers.0"].(model.Map)
	if !ok {
		t.Fatalf("Expected filters[\"routers.0\"] to be a model.Map, got: %T (%v)", filters["routers.0"], filters["routers.0"])
	}
	if exists, _ := routers0["$exists"].(bool); !exists {
		t.Errorf("Expected filters[\"routers.0\"][\"$exists\"] = true, got: %v", routers0["$exists"])
	}

	// 2) routers $not $elemMatch {kind $nin Kinds}.
	routersFilter, ok := filters["routers"].(model.Map)
	if !ok {
		t.Fatalf("Expected filters[\"routers\"] to be a model.Map, got: %T (%v)", filters["routers"], filters["routers"])
	}
	notClause, ok := routersFilter["$not"].(model.Map)
	if !ok {
		t.Fatalf("Expected filters[\"routers\"][\"$not\"] to be a model.Map, got: %T", routersFilter["$not"])
	}
	elemMatch, ok := notClause["$elemMatch"].(model.Map)
	if !ok {
		t.Fatalf("Expected $not.$elemMatch to be a model.Map, got: %T", notClause["$elemMatch"])
	}
	kindClause, ok := elemMatch["kind"].(model.Map)
	if !ok {
		t.Fatalf("Expected $elemMatch.kind to be a model.Map, got: %T", elemMatch["kind"])
	}
	ninRaw, ok := kindClause["$nin"]
	if !ok {
		t.Fatalf("Expected kind.$nin to be present, got: %v", kindClause)
	}
	ninKinds, ok := ninRaw.([]string)
	if !ok {
		t.Fatalf("Expected kind.$nin to be []string, got: %T (%v)", ninRaw, ninRaw)
	}
	if len(ninKinds) != len(kinds) {
		t.Fatalf("Expected kind.$nin len=%d, got: %d", len(kinds), len(ninKinds))
	}
	for i, k := range kinds {
		if ninKinds[i] != k {
			t.Errorf("kind.$nin[%d]: got=%q, want=%q", i, ninKinds[i], k)
		}
	}

	// Backward-compat: when Kinds is empty/nil, neither key is emitted.
	repoEmpty, _, _, depsEmpty := createServiceDependencies()
	repoEmpty.FindWithFiltersResponse = &model.PaginatedResult[entities.RouteGroup]{
		Items: []entities.RouteGroup{},
		Pagination: model.Pagination{Page: 1, PerPage: 20, TotalItems: 0, TotalPages: 0},
	}
	serviceEmpty := New(depsEmpty)
	_, err = serviceEmpty.GetRouteGroups(ctx.Background(), reqContext, &dtos.RouteGroupQueryDTO{})
	if err != nil {
		t.Fatalf("Expected no error (empty Kinds), got: %v", err)
	}
	if len(repoEmpty.FindWithFiltersCalls) != 1 {
		t.Fatalf("Expected 1 FindWithFilters call (empty Kinds), got: %d", len(repoEmpty.FindWithFiltersCalls))
	}
	emptyFilters := repoEmpty.FindWithFiltersCalls[0].Filters
	if _, present := emptyFilters["routers.0"]; present {
		t.Errorf("Expected no filters[\"routers.0\"] when Kinds is empty, got: %v", emptyFilters["routers.0"])
	}
	if _, present := emptyFilters["routers"]; present {
		t.Errorf("Expected no filters[\"routers\"] when Kinds is empty, got: %v", emptyFilters["routers"])
	}
}
