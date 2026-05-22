package mocks

import (
	ctx "context"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"
	"router/src/modules/routegroups/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// Compile-time interface check
var _ ports.RouteGroupServicePort = (*MockRouteGroupService)(nil)

// MockRouteGroupService is a mock implementation of RouteGroupServicePort.
type MockRouteGroupService struct {
	// Return values for each method
	GetByIdResponse *dtos.RouteGroupResponse
	GetByIdError    error

	CreateResponse *dtos.RouteGroupResponse
	CreateError    error

	UpdateResponse *dtos.RouteGroupResponse
	UpdateError    error

	DeleteResponse map[string]bool
	DeleteError    error

	GetListResponse *model.PaginatedResult[dtos.RouteGroupResponse]
	GetListError    error

	GetByIdsResponse []dtos.RouteGroupResponse
	GetByIdsError    error

	CountResponse int64
	CountError    error

	// Call tracking
	GetByIdCalls   []string
	CreateCalls    []dtos.RouteGroupCreateDTO
	UpdateCalls    []UpdateCallRecord
	DeleteCalls    []string
	GetListCalls   []dtos.RouteGroupQueryDTO
	GetByIdsCalls  [][]string
	CountCalls     int
}

// UpdateCallRecord tracks calls to UpdateRouteGroupById method
type UpdateCallRecord struct {
	RouteGroupId string
	DTO          dtos.RouteGroupUpdateDTO
}

// NewMockRouteGroupService creates a new MockRouteGroupService.
func NewMockRouteGroupService() *MockRouteGroupService {
	return &MockRouteGroupService{
		GetByIdCalls:  make([]string, 0),
		CreateCalls:   make([]dtos.RouteGroupCreateDTO, 0),
		UpdateCalls:   make([]UpdateCallRecord, 0),
		DeleteCalls:   make([]string, 0),
		GetListCalls:  make([]dtos.RouteGroupQueryDTO, 0),
		GetByIdsCalls: make([][]string, 0),
	}
}

// CreateRouteGroup creates a new route group from the provided DTO.
func (m *MockRouteGroupService) CreateRouteGroup(
	context ctx.Context,
	requestContext *reqCtx.RequestContext,
	dto *dtos.RouteGroupCreateDTO,
) (*dtos.RouteGroupResponse, error) {
	if dto != nil {
		m.CreateCalls = append(m.CreateCalls, *dto)
	}

	if m.CreateError != nil {
		return nil, m.CreateError
	}

	return m.CreateResponse, nil
}

// GetRouteGroupById retrieves a route group by its unique identifier.
func (m *MockRouteGroupService) GetRouteGroupById(
	context ctx.Context,
	routeGroupId *string,
) (*dtos.RouteGroupResponse, error) {
	if routeGroupId != nil {
		m.GetByIdCalls = append(m.GetByIdCalls, *routeGroupId)
	}

	if m.GetByIdError != nil {
		return nil, m.GetByIdError
	}

	return m.GetByIdResponse, nil
}

// UpdateRouteGroupById updates an existing route group.
func (m *MockRouteGroupService) UpdateRouteGroupById(
	context ctx.Context,
	routeGroupId *string,
	dto *dtos.RouteGroupUpdateDTO,
) (*dtos.RouteGroupResponse, error) {
	if routeGroupId != nil && dto != nil {
		m.UpdateCalls = append(m.UpdateCalls, UpdateCallRecord{
			RouteGroupId: *routeGroupId,
			DTO:          *dto,
		})
	}

	if m.UpdateError != nil {
		return nil, m.UpdateError
	}

	return m.UpdateResponse, nil
}

// GetRouteGroups retrieves a paginated and filtered list of route groups.
func (m *MockRouteGroupService) GetRouteGroups(
	context ctx.Context,
	requestContext *reqCtx.RequestContext,
	query *dtos.RouteGroupQueryDTO,
) (*model.PaginatedResult[dtos.RouteGroupResponse], error) {
	if query != nil {
		m.GetListCalls = append(m.GetListCalls, *query)
	}

	if m.GetListError != nil {
		return nil, m.GetListError
	}

	return m.GetListResponse, nil
}

// DeleteRouteGroupById removes a route group by its unique identifier.
func (m *MockRouteGroupService) DeleteRouteGroupById(
	context ctx.Context,
	routeGroupId *string,
) (map[string]bool, error) {
	if routeGroupId != nil {
		m.DeleteCalls = append(m.DeleteCalls, *routeGroupId)
	}

	if m.DeleteError != nil {
		return nil, m.DeleteError
	}

	return m.DeleteResponse, nil
}

// GetRouteGroupsByIds retrieves multiple route groups by their IDs.
func (m *MockRouteGroupService) GetRouteGroupsByIds(
	context ctx.Context,
	ids []string,
) ([]dtos.RouteGroupResponse, error) {
	m.GetByIdsCalls = append(m.GetByIdsCalls, ids)

	if m.GetByIdsError != nil {
		return nil, m.GetByIdsError
	}

	return m.GetByIdsResponse, nil
}

// CountRouteGroups returns the total count of route groups for the given org context.
func (m *MockRouteGroupService) CountRouteGroups(
	context ctx.Context,
	requestContext *reqCtx.RequestContext,
) (int64, error) {
	m.CountCalls++

	if m.CountError != nil {
		return 0, m.CountError
	}

	return m.CountResponse, nil
}

func (m *MockRouteGroupService) GetRouteGroupEntityById(
	context ctx.Context,
	routeGroupId *string,
) (*entities.RouteGroup, error) {
	if routeGroupId != nil {
		m.GetByIdCalls = append(m.GetByIdCalls, *routeGroupId)
	}
	if m.GetByIdError != nil {
		return nil, m.GetByIdError
	}
	if m.GetByIdResponse == nil {
		return nil, nil
	}

	entity := &entities.RouteGroup{}
	if m.GetByIdResponse.ID != nil {
		entity.ID = *m.GetByIdResponse.ID
	}
	if m.GetByIdResponse.Name != nil {
		entity.Name = *m.GetByIdResponse.Name
	}
	if m.GetByIdResponse.Description != nil {
		entity.Description = *m.GetByIdResponse.Description
	}
	if m.GetByIdResponse.Routers != nil {
		dtoRouters := *m.GetByIdResponse.Routers
		entity.Routers = make([]entities.Router, len(dtoRouters))
		for i, r := range dtoRouters {
			entity.Routers[i] = entities.Router{Kind: r.Kind}
			if r.Match != nil {
				rules := []entities.MatchRule{}
				if r.Match.Rules != nil {
					for _, rule := range *r.Match.Rules {
						rules = append(rules, entities.MatchRule{
							Field:    rule.Field,
							Operator: rule.Operator,
							Value:    rule.Value,
						})
					}
				}
				entity.Routers[i].Match = &entities.MatchConfig{
					Policy: r.Match.Policy,
					Rules:  rules,
				}
			}
		}
	}
	return entity, nil
}
