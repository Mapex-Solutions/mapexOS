package services

import (
	ctx "context"

	"mapexIam/src/modules/lists/application/di"
	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/application/ports"
	"mapexIam/src/modules/lists/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// Compile-time check to ensure ListService implements ListServicePort
var _ ports.ListServicePort = (*ListService)(nil)

// New creates and returns a new instance of ListService.
func New(deps di.ListServiceDependenciesInjection) ports.ListServicePort {
	return &ListService{deps: deps}
}

// CreateList orchestrates list creation:
// resolve scope (system / template-ancestor / org-local) and apply org
// context -> map DTO to entity, defaulting Scope to "local" -> persist
// -> map back to response DTO.
func (s *ListService) CreateList(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.ListCreateDTO) (*dtos.ListResponse, error) {
	if err := s.applyListScope(requestContext, dto); err != nil {
		return nil, err
	}
	listEntity, _ := mapper.DtoToEntity[dtos.ListCreateDTO, entities.List](dto)
	if listEntity.Scope == "" && !listEntity.IsSystem {
		listEntity.Scope = "local"
	}
	created, err := s.deps.Repo.Create(c, listEntity)
	if err != nil {
		return nil, err
	}
	resp, _ := mapper.EntityToDto[entities.List, dtos.ListResponse](created)
	return resp, nil
}

// GetListById fetches a single list by id. Returns 404 when the id is
// unknown.
func (s *ListService) GetListById(c ctx.Context, listId *string) (*dtos.ListResponse, error) {
	listEntity, err := s.deps.Repo.FindById(c, listId)
	if err != nil {
		return nil, err
	}
	if listEntity == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"List not found"}}
	}
	resp, _ := mapper.EntityToDto[entities.List, dtos.ListResponse](listEntity)
	return resp, nil
}

// UpdateListById orchestrates a partial update: load the prior entity
// (404 on miss) -> apply the patch in Mongo -> when the name flipped on
// a classification list, publish a NATS event so consuming services can
// re-denormalize their cached names.
func (s *ListService) UpdateListById(c ctx.Context, listId *string, dto *dtos.ListUpdateDTO) (*dtos.ListResponse, error) {
	currentList, err := s.deps.Repo.FindById(c, listId)
	if err != nil || currentList == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"List not found"}}
	}

	fields, _ := mapper.DtoToMap(dto)
	updated, _ := s.deps.Repo.FindByIdAndUpdate(c, listId, fields)
	if updated.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"List not found"}}
	}

	if dto.Name != nil && *dto.Name != currentList.Name {
		s.publishListNameUpdatedEvent(c, updated)
	}

	resp, _ := mapper.EntityToDto[entities.List, dtos.ListResponse](updated)
	return resp, nil
}

// DeleteListById orchestrates list removal: delete the row, mapping the
// repository's "document not found" sentinel to 404.
func (s *ListService) DeleteListById(c ctx.Context, listId *string) (map[string]bool, error) {
	if err := s.deps.Repo.DeleteById(c, listId); err != nil {
		if err.Error() == "document not found" {
			return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"List not found"}}
		}
		return nil, err
	}
	return map[string]bool{"success": true}, nil
}

// GetListByEmail looks up a list by its email address.
func (s *ListService) GetListByEmail(c ctx.Context, email *string) (*entities.List, error) {
	return s.deps.Repo.FindByEmail(c, email)
}

// GetLists orchestrates the paginated list:
// build the org/system/ancestor filter union -> apply optional type /
// name / enabled / parentId filters -> resolve pagination + projection
// -> delegate to repository -> map entities to response DTOs.
func (s *ListService) GetLists(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.ListQueryDTO) (*model.PaginatedResult[dtos.ListResponse], error) {
	filters := s.buildListListFilters(requestContext, query)
	pagination := &model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}
	projection := s.buildListProjection(query)

	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]dtos.ListResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.List, dtos.ListResponse](&entity)
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.ListResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}
