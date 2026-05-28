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
//
// StringToObjectId is enabled so the DTO's `*string` ParentId is converted
// to `*model.ObjectId` on the entity (the DTO uses string so the validator's
// `mongoid` tag can inspect the underlying hex).
func (s *ListService) CreateList(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.ListCreateDTO) (*dtos.ListResponse, error) {
	if err := s.applyListScope(requestContext, dto); err != nil {
		return nil, err
	}
	listEntity, _ := mapper.DtoToEntityWithOptions[dtos.ListCreateDTO, entities.List](dto, mapper.MapperOptions{StringToObjectId: true})
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
// unknown. The response is enriched with ParentName / ParentType so the
// UI can show the parent hierarchy without an extra round-trip.
func (s *ListService) GetListById(c ctx.Context, listId *string) (*dtos.ListResponse, error) {
	listEntity, err := s.deps.Repo.FindById(c, listId)
	if err != nil {
		return nil, err
	}
	if listEntity == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"List not found"}}
	}
	resp, _ := mapper.EntityToDto[entities.List, dtos.ListResponse](listEntity)
	s.populateParentFields(c, []*dtos.ListResponse{resp}, []*entities.List{listEntity})
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
	// parentId travels through the contract DTO as a string (so the validator
	// can inspect it) — convert it to ObjectID before letting Mongo $set the
	// document, otherwise the field would be persisted as a string and break
	// downstream lookups.
	if pid, ok := fields["parentId"].(string); ok && pid != "" {
		if oid, err := model.ToObjectID(pid); err == nil {
			fields["parentId"] = oid
		}
	}
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
	dtoPtrs := make([]*dtos.ListResponse, len(result.Items))
	entityPtrs := make([]*entities.List, len(result.Items))
	for i := range result.Items {
		entity := result.Items[i]
		dto, _ := mapper.EntityToDto[entities.List, dtos.ListResponse](&entity)
		dtoItems[i] = *dto
		dtoPtrs[i] = &dtoItems[i]
		entityPtrs[i] = &result.Items[i]
	}
	s.populateParentFields(c, dtoPtrs, entityPtrs)
	return &model.PaginatedResult[dtos.ListResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}

// populateParentFields resolves ParentName / ParentType on the supplied
// response DTOs by issuing one batched FindByIds against the parents.
//
// The DTO slice and the entity slice MUST line up index-by-index: the DTOs
// are mutated in place and the entities are read for their ParentId.
// Top-level items (ParentId == nil) are left untouched.
func (s *ListService) populateParentFields(c ctx.Context, responses []*dtos.ListResponse, items []*entities.List) {
	parentIDStrs := make([]string, 0, len(items))
	for _, item := range items {
		if item == nil || item.ParentId == nil || item.ParentId.IsZero() {
			continue
		}
		parentIDStrs = append(parentIDStrs, item.ParentId.Hex())
	}
	if len(parentIDStrs) == 0 {
		return
	}

	parents, err := s.deps.Repo.FindByIds(c, parentIDStrs)
	if err != nil || len(parents) == 0 {
		return
	}

	parentByID := make(map[string]*entities.List, len(parents))
	for _, p := range parents {
		if p == nil {
			continue
		}
		parentByID[p.ID.Hex()] = p
	}

	for i, item := range items {
		if item == nil || item.ParentId == nil || item.ParentId.IsZero() {
			continue
		}
		parent, ok := parentByID[item.ParentId.Hex()]
		if !ok {
			continue
		}
		name := parent.Name
		parentType := parent.Type
		responses[i].ParentName = &name
		responses[i].ParentType = &parentType
	}
}
