package services

import (
	ctx "context"

	"mapexIam/src/modules/roles/application/di"
	"mapexIam/src/modules/roles/application/dtos"
	"mapexIam/src/modules/roles/application/ports"
	"mapexIam/src/modules/roles/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// Compile-time check to ensure RoleService implements RoleServicePort
var _ ports.RoleServicePort = (*RoleService)(nil)

// New creates and returns a new instance of RoleService.
func New(deps di.RoleServiceDependenciesInjection) ports.RoleServicePort {
	return &RoleService{deps: deps}
}

// CreateRole orchestrates role creation:
// resolve scope (system / template-ancestor / org-local) and apply org
// context -> map DTO to entity and assign multi-tenant fields ->
// persist -> map back to response DTO.
func (s *RoleService) CreateRole(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateRoleDto) (*dtos.RoleResponse, error) {
	if err := s.applyRoleScope(requestContext, dto); err != nil {
		return nil, err
	}
	roleEntity, _ := mapper.DtoToEntity[dtos.CreateRoleDto, entities.Role](dto)
	if err := assignRoleTenantFields(roleEntity, dto); err != nil {
		return nil, err
	}
	created, err := s.deps.Repo.Create(c, roleEntity)
	if err != nil {
		return nil, err
	}
	resp, _ := mapper.EntityToDto[entities.Role, dtos.RoleResponse](created)
	return resp, nil
}

// GetRoleById fetches a single role by id. Returns 404 when the id is
// unknown.
func (s *RoleService) GetRoleById(c ctx.Context, roleId *string) (*dtos.RoleResponse, error) {
	roleEntity, err := s.deps.Repo.FindById(c, roleId)
	if err != nil {
		return nil, err
	}
	if roleEntity == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Role not found"}}
	}
	resp, _ := mapper.EntityToDto[entities.Role, dtos.RoleResponse](roleEntity)
	return resp, nil
}

// UpdateRoleById orchestrates a partial role update: load the prior
// entity (404 on miss) -> apply the patch in Mongo -> publish a NATS
// cache-invalidation event when permissions changed -> map to response.
func (s *RoleService) UpdateRoleById(c ctx.Context, roleId *string, dto *dtos.UpdateRoleDto) (*dtos.RoleResponse, error) {
	oldRole, err := s.deps.Repo.FindById(c, roleId)
	if err != nil {
		return nil, err
	}
	if oldRole == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Role not found"}}
	}

	fields, _ := mapper.DtoToMap(dto)
	updated, _ := s.deps.Repo.FindByIdAndUpdate(c, roleId, fields)
	if updated.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Role not found"}}
	}

	if dto.Permissions != nil && len(*dto.Permissions) > 0 {
		go s.publishRolePermissionsChanged(*roleId, oldRole.Permissions, *dto.Permissions)
	}

	resp, _ := mapper.EntityToDto[entities.Role, dtos.RoleResponse](updated)
	return resp, nil
}

// DeleteRoleById orchestrates role deletion: delete the row (404 on
// miss) -> publish a NATS cache-invalidation event so consumers drop
// every cached principal that referenced this role.
func (s *RoleService) DeleteRoleById(c ctx.Context, roleId *string) (map[string]bool, error) {
	if err := s.deps.Repo.DeleteById(c, roleId); err != nil {
		if err.Error() == "document not found" {
			return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Role not found"}}
		}
		return nil, err
	}

	go s.publishRoleDeleted(*roleId)
	return map[string]bool{"success": true}, nil
}

// GetRoles orchestrates the paginated list:
// build the org/system/ancestor filter union -> apply optional name /
// scope / permission filters -> resolve pagination + projection ->
// delegate to repository -> map entities to DTOs.
func (s *RoleService) GetRoles(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.RoleQueryDto) (*model.PaginatedResult[dtos.RoleResponse], error) {
	filters := s.buildRoleListFilters(requestContext, query)
	pagination := &model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}
	projection := s.buildRoleProjection(query)

	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]dtos.RoleResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.Role, dtos.RoleResponse](&entity)
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.RoleResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}

// assignRoleTenantFields fills the multi-tenant fields on the freshly
// mapped entity. System roles get cleared org/path; org-local roles get
// scope="local"; missing OrgID for non-system is rejected as 400.
func assignRoleTenantFields(roleEntity *entities.Role, dto *dtos.CreateRoleDto) error {
	if dto.IsSystem {
		roleEntity.PathKey = ""
		roleEntity.OrgID = nil
		roleEntity.Scope = ""
		return nil
	}
	if dto.OrgID != nil {
		roleEntity.Scope = "local"
		return nil
	}
	return &customErrors.ServerCustomError{
		Code:   httpStatus.BAD_REQUEST,
		Errors: []string{"OrgID is required for non-system roles"},
	}
}

