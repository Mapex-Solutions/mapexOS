package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"triggers/src/modules/triggers/application/di"
	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
)

/**
 * Mock Implementations
 */

// MockTriggerRepository is a mock implementation of TriggerRepository
type MockTriggerRepository struct {
	CreateFunc           func(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error)
	FindByIdFunc         func(ctx context.Context, triggerId *string) (*entities.Trigger, error)
	FindByIdAndUpdateFunc func(ctx context.Context, triggerId *string, payload map[string]any) (*entities.Trigger, error)
	DeleteByIdFunc       func(ctx context.Context, triggerId *string) error
	FindWithFiltersFunc  func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error)
}

func (m *MockTriggerRepository) Create(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, trigger)
	}
	return trigger, nil
}

func (m *MockTriggerRepository) FindById(ctx context.Context, triggerId *string) (*entities.Trigger, error) {
	if m.FindByIdFunc != nil {
		return m.FindByIdFunc(ctx, triggerId)
	}
	return nil, errors.New("not implemented")
}

func (m *MockTriggerRepository) FindByIdAndUpdate(ctx context.Context, triggerId *string, payload map[string]any) (*entities.Trigger, error) {
	if m.FindByIdAndUpdateFunc != nil {
		return m.FindByIdAndUpdateFunc(ctx, triggerId, payload)
	}
	return nil, errors.New("not implemented")
}

func (m *MockTriggerRepository) DeleteById(ctx context.Context, triggerId *string) error {
	if m.DeleteByIdFunc != nil {
		return m.DeleteByIdFunc(ctx, triggerId)
	}
	return nil
}

func (m *MockTriggerRepository) FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
	if m.FindWithFiltersFunc != nil {
		return m.FindWithFiltersFunc(ctx, filters, pagination, projection)
	}
	return &model.PaginatedResult[entities.Trigger]{
		Items:      []entities.Trigger{},
		Pagination: model.Pagination{},
	}, nil
}

func (m *MockTriggerRepository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	return 0, nil
}

// MockCacheRepository is a mock implementation of CacheRepository
type MockCacheRepository struct {
	GetFunc        func(ctx context.Context, key string, dest interface{}) error
	SetFunc        func(ctx context.Context, key string, value interface{}) error
	SetExFunc      func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	DelFunc        func(ctx context.Context, key string) error
	GetOrSetExFunc func(params common.GetOrSetParams) (interface{}, error)
}

func (m *MockCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key, dest)
	}
	return errors.New("cache miss")
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value)
	}
	return nil
}

func (m *MockCacheRepository) SetEx(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if m.SetExFunc != nil {
		return m.SetExFunc(ctx, key, value, ttl)
	}
	return nil
}

func (m *MockCacheRepository) Del(ctx context.Context, key string) error {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, key)
	}
	return nil
}

func (m *MockCacheRepository) GetOrSetEx(params common.GetOrSetParams) (interface{}, error) {
	if m.GetOrSetExFunc != nil {
		return m.GetOrSetExFunc(params)
	}
	// Default behavior: execute callback
	if params.Callback != nil {
		return params.Callback()
	}
	return nil, nil
}

// MockCacheKeyBuilder is a mock implementation of TriggerCacheKeyBuilderPort.
type MockCacheKeyBuilder struct{}

func (m *MockCacheKeyBuilder) TriggerKey(triggerId string) string {
	return "TRIGGER:" + triggerId
}

func (m *MockCacheKeyBuilder) CounterKey(orgId string) string {
	return "counter:triggers:" + orgId
}

/**
 * Helper Functions
 */

func createTestService(triggerRepo *MockTriggerRepository, cacheRepo *MockCacheRepository) *TriggerService {
	return &TriggerService{
		deps: di.TriggerServiceDependenciesInjection{
			TriggerRepository: triggerRepo,
			CacheRepository:   cacheRepo,
			AppCache:          cacheRepo,
			CacheKeyBuilder:   &MockCacheKeyBuilder{},
		},
	}
}

func createTestRequestContext() *reqCtx.RequestContext {
	orgId := "org123"
	return &reqCtx.RequestContext{
		OrgContext: &orgId,
		OrgContextData: &reqCtx.CoverageOrg{
			ID:      "org123",
			Name:    "Test Org",
			Type:    "customer",
			PathKey: "mapex.vendor.customer",
		},
	}
}

func createTestTrigger() *entities.Trigger {
	id, _ := model.ToObjectID("trigger123")
	return &entities.Trigger{
		ID:          id,
		Name:        "Test Trigger",
		Description: typeconv.PtrString("Test Description"),
		TriggerType: "http",
		Category:    "technical",
		Enabled:     true,
		Config: entities.TriggerConfig{
			Http: &entities.HttpConfig{
				Endpoint: "https://api.example.com",
				Method:   "POST",
			},
		},
		IsSystem:   false,
		IsTemplate: false,
		Created:    time.Now(),
		Updated:    time.Now(),
	}
}

/**
 * CreateTrigger Tests
 */

func TestTriggerService_CreateTrigger_Success(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		CreateFunc: func(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error) {
			id, _ := model.ToObjectID("newtrigger123")
			trigger.ID = id
			return trigger, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	dto := &dtos.CreateTriggerDto{
		Name:        "New HTTP Trigger",
		Description: typeconv.PtrString("A test HTTP trigger"),
		TriggerType: "http",
		Category:    "technical",
		Enabled:     true,
		IsSystem:    false,
		IsTemplate:  false,
	}

	result, err := service.CreateTrigger(context.Background(), reqCtx, dto)

	if err != nil {
		t.Fatalf("CreateTrigger() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("CreateTrigger() returned nil result")
	}

	if result.Name == nil || *result.Name != "New HTTP Trigger" {
		t.Errorf("CreateTrigger() name = %v, want 'New HTTP Trigger'", result.Name)
	}
}

func TestTriggerService_CreateTrigger_SystemTrigger(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		CreateFunc: func(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error) {
			// Verify system trigger properties
			if trigger.OrgID != nil {
				t.Error("System trigger should have nil OrgID")
			}
			if trigger.PathKey != "" {
				t.Error("System trigger should have empty PathKey")
			}
			if !trigger.IsSystem {
				t.Error("System trigger should have IsSystem=true")
			}

			id, _ := model.ToObjectID("systemtrigger123")
			trigger.ID = id
			return trigger, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	dto := &dtos.CreateTriggerDto{
		Name:        "System HTTP Trigger",
		TriggerType: "http",
		Category:    "technical",
		Enabled:     true,
		IsSystem:    true,
	}

	result, err := service.CreateTrigger(context.Background(), reqCtx, dto)

	if err != nil {
		t.Fatalf("CreateTrigger() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("CreateTrigger() returned nil result")
	}
}

func TestTriggerService_CreateTrigger_RepositoryError(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		CreateFunc: func(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error) {
			return nil, errors.New("database error")
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	dto := &dtos.CreateTriggerDto{
		Name:        "Test Trigger",
		TriggerType: "http",
		Category:    "technical",
		Enabled:     true,
	}

	_, err := service.CreateTrigger(context.Background(), reqCtx, dto)

	if err == nil {
		t.Fatal("CreateTrigger() expected error, got nil")
	}
}

/**
 * GetTriggerById Tests
 */

func TestTriggerService_GetTriggerById_Success(t *testing.T) {
	testTrigger := createTestTrigger()

	mockRepo := &MockTriggerRepository{
		FindByIdFunc: func(ctx context.Context, triggerId *string) (*entities.Trigger, error) {
			return testTrigger, nil
		},
	}
	mockCache := &MockCacheRepository{
		GetOrSetExFunc: func(params common.GetOrSetParams) (interface{}, error) {
			// Simulate cache miss, callback executed
			result, err := params.Callback()
			if err != nil {
				return nil, err
			}
			// Copy to dest
			if trigger, ok := result.(*entities.Trigger); ok {
				if dest, ok := params.Dest.(*entities.Trigger); ok {
					*dest = *trigger
				}
			}
			return result, nil
		},
	}

	service := createTestService(mockRepo, mockCache)
	triggerId := "trigger123"

	result, err := service.GetTriggerById(context.Background(), &triggerId)

	if err != nil {
		t.Fatalf("GetTriggerById() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetTriggerById() returned nil result")
	}

	if result.Name == nil || *result.Name != "Test Trigger" {
		t.Errorf("GetTriggerById() name = %v, want 'Test Trigger'", result.Name)
	}
}

func TestTriggerService_GetTriggerById_NotFound(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindByIdFunc: func(ctx context.Context, triggerId *string) (*entities.Trigger, error) {
			return nil, errors.New("not found")
		},
	}
	mockCache := &MockCacheRepository{
		GetOrSetExFunc: func(params common.GetOrSetParams) (interface{}, error) {
			return params.Callback()
		},
	}

	service := createTestService(mockRepo, mockCache)
	triggerId := "nonexistent"

	_, err := service.GetTriggerById(context.Background(), &triggerId)

	if err == nil {
		t.Fatal("GetTriggerById() expected error for non-existent trigger, got nil")
	}
}

/**
 * UpdateTriggerById Tests
 */

func TestTriggerService_UpdateTriggerById_Success(t *testing.T) {
	updatedTrigger := createTestTrigger()
	updatedTrigger.Name = "Updated Trigger"

	mockRepo := &MockTriggerRepository{
		FindByIdAndUpdateFunc: func(ctx context.Context, triggerId *string, payload map[string]any) (*entities.Trigger, error) {
			return updatedTrigger, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()
	triggerId := "trigger123"

	newName := "Updated Trigger"
	dto := &dtos.UpdateTriggerDto{
		Name: &newName,
	}

	result, err := service.UpdateTriggerById(context.Background(), reqCtx, &triggerId, dto)

	if err != nil {
		t.Fatalf("UpdateTriggerById() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("UpdateTriggerById() returned nil result")
	}
}

func TestTriggerService_UpdateTriggerById_NotFound(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindByIdAndUpdateFunc: func(ctx context.Context, triggerId *string, payload map[string]any) (*entities.Trigger, error) {
			return nil, errors.New("not found")
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()
	triggerId := "nonexistent"

	newName := "Updated Trigger"
	dto := &dtos.UpdateTriggerDto{
		Name: &newName,
	}

	_, err := service.UpdateTriggerById(context.Background(), reqCtx, &triggerId, dto)

	if err == nil {
		t.Fatal("UpdateTriggerById() expected error for non-existent trigger, got nil")
	}
}

/**
 * DeleteTriggerById Tests
 */

func TestTriggerService_DeleteTriggerById_Success(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		DeleteByIdFunc: func(ctx context.Context, triggerId *string) error {
			return nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	triggerId := "trigger123"

	result, err := service.DeleteTriggerById(context.Background(), &triggerId)

	if err != nil {
		t.Fatalf("DeleteTriggerById() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("DeleteTriggerById() returned nil result")
	}

	if !result["success"] {
		t.Error("DeleteTriggerById() expected success=true")
	}
}

func TestTriggerService_DeleteTriggerById_NotFound(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		DeleteByIdFunc: func(ctx context.Context, triggerId *string) error {
			return errors.New("not found")
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	triggerId := "nonexistent"

	_, err := service.DeleteTriggerById(context.Background(), &triggerId)

	if err == nil {
		t.Fatal("DeleteTriggerById() expected error for non-existent trigger, got nil")
	}
}

/**
 * GetTriggers Tests
 */

func TestTriggerService_GetTriggers_Success(t *testing.T) {
	testTrigger := createTestTrigger()

	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			return &model.PaginatedResult[entities.Trigger]{
				Items: []entities.Trigger{*testTrigger},
				Pagination: model.Pagination{
					Page:       1,
					PerPage:    10,
					TotalItems: 1,
					TotalPages: 1,
				},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	query := &dtos.TriggerQueryDto{}

	result, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("GetTriggers() returned nil result")
	}

	if len(result.Items) != 1 {
		t.Errorf("GetTriggers() items count = %d, want 1", len(result.Items))
	}
}

func TestTriggerService_GetTriggers_WithFilters(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify filters are applied
			if filters["triggerType"] != "http" {
				t.Error("Expected triggerType filter to be 'http'")
			}
			if filters["enabled"] != true {
				t.Error("Expected enabled filter to be true")
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items:      []entities.Trigger{},
				Pagination: model.Pagination{},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	triggerType := "http"
	enabled := true
	query := &dtos.TriggerQueryDto{
		TriggerType: &triggerType,
		Enabled:     &enabled,
	}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}
}

func TestTriggerService_GetTriggers_WithPagination(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify pagination is applied
			if pagination.Page != 2 {
				t.Errorf("Expected page = 2, got %d", pagination.Page)
			}
			if pagination.PerPage != 20 {
				t.Errorf("Expected perPage = 20, got %d", pagination.PerPage)
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items: []entities.Trigger{},
				Pagination: model.Pagination{
					Page:       2,
					PerPage:    20,
					TotalItems: 50,
					TotalPages: 3,
				},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	page := 2
	perPage := 20
	query := &dtos.TriggerQueryDto{
		Page:    &page,
		PerPage: &perPage,
	}

	result, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}

	if result.Pagination.Page != 2 {
		t.Errorf("GetTriggers() page = %d, want 2", result.Pagination.Page)
	}
}

func TestTriggerService_GetTriggers_RepositoryError(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			return nil, errors.New("database error")
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	query := &dtos.TriggerQueryDto{}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err == nil {
		t.Fatal("GetTriggers() expected error, got nil")
	}
}

/**
 * GetTriggers $or Filter Tests
 */

func TestTriggerService_GetTriggers_DefaultIncludesSystemTemplates(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify $or contains isSystem: true by default
			orConditions, ok := filters["$or"].([]model.Map)
			if !ok {
				t.Fatal("Expected $or filter to be present")
			}

			// Check that isSystem: true is in the $or conditions
			hasSystemFilter := false
			for _, cond := range orConditions {
				if val, exists := cond["isSystem"]; exists && val == true {
					hasSystemFilter = true
					break
				}
			}

			if !hasSystemFilter {
				t.Error("Default query should include {isSystem: true} in $or conditions")
			}

			// Verify no isSystem AND filter exists (not excluded)
			if _, exists := filters["isSystem"]; exists {
				t.Error("Default query should NOT have isSystem as an AND filter")
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items:      []entities.Trigger{},
				Pagination: model.Pagination{},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	// Default query: both isSystem and isTemplate are nil
	query := &dtos.TriggerQueryDto{}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}
}

func TestTriggerService_GetTriggers_IsSystemFalseExcludesSystem(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify $or does NOT contain isSystem: true
			orConditions, ok := filters["$or"].([]model.Map)
			if ok {
				for _, cond := range orConditions {
					if val, exists := cond["isSystem"]; exists && val == true {
						t.Error("isSystem=false should NOT include {isSystem: true} in $or conditions")
					}
				}
			}

			// Verify isSystem=false is set as AND filter
			if val, exists := filters["isSystem"]; !exists || val != false {
				t.Error("isSystem=false should set isSystem: false as an AND filter")
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items:      []entities.Trigger{},
				Pagination: model.Pagination{},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	isSystem := false
	query := &dtos.TriggerQueryDto{
		IsSystem: &isSystem,
	}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}
}

func TestTriggerService_GetTriggers_IsSystemTrueIncludesSystem(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify $or contains isSystem: true
			orConditions, ok := filters["$or"].([]model.Map)
			if !ok {
				t.Fatal("Expected $or filter to be present")
			}

			hasSystemFilter := false
			for _, cond := range orConditions {
				if val, exists := cond["isSystem"]; exists && val == true {
					hasSystemFilter = true
					break
				}
			}

			if !hasSystemFilter {
				t.Error("isSystem=true should include {isSystem: true} in $or conditions")
			}

			// Verify no isSystem AND filter (not excluded)
			if _, exists := filters["isSystem"]; exists {
				t.Error("isSystem=true should NOT have isSystem as an AND filter")
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items:      []entities.Trigger{},
				Pagination: model.Pagination{},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	isSystem := true
	query := &dtos.TriggerQueryDto{
		IsSystem: &isSystem,
	}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}
}

func TestTriggerService_GetTriggers_IsTemplateTrueIncludesAncestors(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify $or contains a template ancestor filter
			orConditions, ok := filters["$or"].([]model.Map)
			if !ok {
				t.Fatal("Expected $or filter to be present")
			}

			hasTemplateFilter := false
			for _, cond := range orConditions {
				if _, exists := cond["isTemplate"]; exists {
					hasTemplateFilter = true
					break
				}
			}

			if !hasTemplateFilter {
				t.Error("isTemplate=true should include ancestor template filter in $or conditions")
			}

			// Verify no isTemplate AND filter (not excluded)
			if _, exists := filters["isTemplate"]; exists {
				t.Error("isTemplate=true should NOT have isTemplate as an AND filter")
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items:      []entities.Trigger{},
				Pagination: model.Pagination{},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	isTemplate := true
	query := &dtos.TriggerQueryDto{
		IsTemplate: &isTemplate,
	}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}
}

func TestTriggerService_GetTriggers_IsTemplateFalseExcludesTemplates(t *testing.T) {
	mockRepo := &MockTriggerRepository{
		FindWithFiltersFunc: func(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error) {
			// Verify $or does NOT contain template filter
			orConditions, ok := filters["$or"].([]model.Map)
			if ok {
				for _, cond := range orConditions {
					if _, exists := cond["isTemplate"]; exists {
						t.Error("isTemplate=false should NOT include template filter in $or conditions")
					}
				}
			}

			// Verify isTemplate=false is set as AND filter
			if val, exists := filters["isTemplate"]; !exists || val != false {
				t.Error("isTemplate=false should set isTemplate: false as an AND filter")
			}

			return &model.PaginatedResult[entities.Trigger]{
				Items:      []entities.Trigger{},
				Pagination: model.Pagination{},
			}, nil
		},
	}
	mockCache := &MockCacheRepository{}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	isTemplate := false
	query := &dtos.TriggerQueryDto{
		IsTemplate: &isTemplate,
	}

	_, err := service.GetTriggers(context.Background(), reqCtx, query)

	if err != nil {
		t.Fatalf("GetTriggers() unexpected error: %v", err)
	}
}

/**
 * Cache Behavior Tests
 */

func TestTriggerService_CreateTrigger_CachesProperly(t *testing.T) {
	cacheSetCalled := false
	expectedId := "507f1f77bcf86cd799439011" // Valid MongoDB ObjectID

	mockRepo := &MockTriggerRepository{
		CreateFunc: func(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error) {
			id, _ := model.ToObjectID(expectedId)
			trigger.ID = id
			return trigger, nil
		},
	}
	mockCache := &MockCacheRepository{
		SetExFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
			cacheSetCalled = true
			expectedKey := "TRIGGER:" + expectedId
			if key != expectedKey {
				t.Errorf("Cache key = %s, want '%s'", key, expectedKey)
			}
			return nil
		},
	}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()

	dto := &dtos.CreateTriggerDto{
		Name:        "Test Trigger",
		TriggerType: "http",
		Category:    "technical",
		Enabled:     true,
	}

	_, err := service.CreateTrigger(context.Background(), reqCtx, dto)

	if err != nil {
		t.Fatalf("CreateTrigger() unexpected error: %v", err)
	}

	if !cacheSetCalled {
		t.Error("CreateTrigger() should cache the created trigger")
	}
}

func TestTriggerService_UpdateTriggerById_InvalidatesCache(t *testing.T) {
	cacheDelCalled := false
	updatedTrigger := createTestTrigger()

	mockRepo := &MockTriggerRepository{
		FindByIdAndUpdateFunc: func(ctx context.Context, triggerId *string, payload map[string]any) (*entities.Trigger, error) {
			return updatedTrigger, nil
		},
	}
	mockCache := &MockCacheRepository{
		DelFunc: func(ctx context.Context, key string) error {
			cacheDelCalled = true
			return nil
		},
	}

	service := createTestService(mockRepo, mockCache)
	reqCtx := createTestRequestContext()
	triggerId := "trigger123"

	newName := "Updated Trigger"
	dto := &dtos.UpdateTriggerDto{
		Name: &newName,
	}

	_, err := service.UpdateTriggerById(context.Background(), reqCtx, &triggerId, dto)

	if err != nil {
		t.Fatalf("UpdateTriggerById() unexpected error: %v", err)
	}

	if !cacheDelCalled {
		t.Error("UpdateTriggerById() should invalidate cache")
	}
}

func TestTriggerService_DeleteTriggerById_InvalidatesCache(t *testing.T) {
	cacheDelCalled := false

	mockRepo := &MockTriggerRepository{
		DeleteByIdFunc: func(ctx context.Context, triggerId *string) error {
			return nil
		},
	}
	mockCache := &MockCacheRepository{
		DelFunc: func(ctx context.Context, key string) error {
			cacheDelCalled = true
			return nil
		},
	}

	service := createTestService(mockRepo, mockCache)
	triggerId := "trigger123"

	_, err := service.DeleteTriggerById(context.Background(), &triggerId)

	if err != nil {
		t.Fatalf("DeleteTriggerById() unexpected error: %v", err)
	}

	if !cacheDelCalled {
		t.Error("DeleteTriggerById() should invalidate cache")
	}
}
