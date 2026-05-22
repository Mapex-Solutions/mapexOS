package mocks

import (
	"context"

	"router/src/modules/routegroups/domain/entities"
	"router/src/modules/routegroups/domain/repositories"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// Compile-time check to ensure MockRouteGroupRepository implements RouteGroupRepository.
var _ repositories.RouteGroupRepository = (*MockRouteGroupRepository)(nil)

// MockRouteGroupRepository is a mock implementation of RouteGroupRepository.
type MockRouteGroupRepository struct {
	// Return values
	CreateResponse           *entities.RouteGroup
	CreateError              error
	FindByIdResponse         *entities.RouteGroup
	FindByIdError            error
	FindByIdAndUpdateResponse *entities.RouteGroup
	FindByIdAndUpdateError    error
	DeleteByIdError          error
	FindWithFiltersResponse  *model.PaginatedResult[entities.RouteGroup]
	FindWithFiltersError     error
	CountDocumentsResponse   int64
	CountDocumentsError      error

	// Call tracking
	CreateCalls            []*entities.RouteGroup
	FindByIdCalls          []string
	FindByIdAndUpdateCalls []FindByIdAndUpdateCallRecord
	DeleteByIdCalls        []string
	FindWithFiltersCalls   []FindWithFiltersCallRecord
	CountDocumentsCalls    []model.Map
}

// FindByIdAndUpdateCallRecord tracks calls to FindByIdAndUpdate.
type FindByIdAndUpdateCallRecord struct {
	ID      string
	Payload map[string]any
}

// FindWithFiltersCallRecord tracks calls to FindWithFilters.
type FindWithFiltersCallRecord struct {
	Filters    model.Map
	Pagination *model.PaginationOpts
	Projection model.Map
}

// NewMockRouteGroupRepository creates a new MockRouteGroupRepository.
func NewMockRouteGroupRepository() *MockRouteGroupRepository {
	return &MockRouteGroupRepository{
		CreateCalls:            make([]*entities.RouteGroup, 0),
		FindByIdCalls:          make([]string, 0),
		FindByIdAndUpdateCalls: make([]FindByIdAndUpdateCallRecord, 0),
		DeleteByIdCalls:        make([]string, 0),
		FindWithFiltersCalls:   make([]FindWithFiltersCallRecord, 0),
		CountDocumentsCalls:    make([]model.Map, 0),
	}
}

// Create mocks the repository Create method.
func (m *MockRouteGroupRepository) Create(_ context.Context, entity *entities.RouteGroup) (*entities.RouteGroup, error) {
	m.CreateCalls = append(m.CreateCalls, entity)

	if m.CreateError != nil {
		return nil, m.CreateError
	}

	return m.CreateResponse, nil
}

// FindById mocks the repository FindById method.
func (m *MockRouteGroupRepository) FindById(_ context.Context, id *string) (*entities.RouteGroup, error) {
	if id != nil {
		m.FindByIdCalls = append(m.FindByIdCalls, *id)
	}

	if m.FindByIdError != nil {
		return nil, m.FindByIdError
	}

	return m.FindByIdResponse, nil
}

// FindByIdAndUpdate mocks the repository FindByIdAndUpdate method.
func (m *MockRouteGroupRepository) FindByIdAndUpdate(_ context.Context, id *string, payload map[string]any) (*entities.RouteGroup, error) {
	if id != nil {
		m.FindByIdAndUpdateCalls = append(m.FindByIdAndUpdateCalls, FindByIdAndUpdateCallRecord{
			ID:      *id,
			Payload: payload,
		})
	}

	if m.FindByIdAndUpdateError != nil {
		return nil, m.FindByIdAndUpdateError
	}

	return m.FindByIdAndUpdateResponse, nil
}

// DeleteById mocks the repository DeleteById method.
func (m *MockRouteGroupRepository) DeleteById(_ context.Context, id *string) error {
	if id != nil {
		m.DeleteByIdCalls = append(m.DeleteByIdCalls, *id)
	}

	return m.DeleteByIdError
}

// FindWithFilters mocks the repository FindWithFilters method.
func (m *MockRouteGroupRepository) FindWithFilters(
	_ context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.RouteGroup], error) {
	m.FindWithFiltersCalls = append(m.FindWithFiltersCalls, FindWithFiltersCallRecord{
		Filters:    filters,
		Pagination: pagination,
		Projection: projection,
	})

	if m.FindWithFiltersError != nil {
		return nil, m.FindWithFiltersError
	}

	return m.FindWithFiltersResponse, nil
}

// CountDocuments mocks the repository CountDocuments method.
func (m *MockRouteGroupRepository) CountDocuments(_ context.Context, filters model.Map) (int64, error) {
	m.CountDocumentsCalls = append(m.CountDocumentsCalls, filters)

	if m.CountDocumentsError != nil {
		return 0, m.CountDocumentsError
	}

	return m.CountDocumentsResponse, nil
}
