package services

import (
	ctx "context"

	"mapexIam/src/modules/users/application/di"
	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/application/ports"
	"mapexIam/src/modules/users/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// Compile-time check to ensure UserService implements UserServicePort.
var _ ports.UserServicePort = (*UserService)(nil)

// New creates and returns a new instance of UserService.
func New(deps di.UserServiceDependenciesInjection) ports.UserServicePort {
	return &UserService{deps: deps}
}

// CreateUser orchestrates user creation: enforce email uniqueness ->
// build entity (default startTour, internal AuthProvider) -> hash password
// -> persist -> invalidate counter cache -> return response DTO.
func (s *UserService) CreateUser(c ctx.Context, dto *dtos.UserCreateDTO) (*dtos.UserResponse, error) {
	if err := s.ensureEmailIsUnique(c, dto.Email); err != nil {
		return nil, err
	}
	userEntity, err := s.buildUserCreateEntity(dto)
	if err != nil {
		return nil, err
	}
	created, err := s.deps.Repo.Create(c, userEntity)
	if err != nil {
		return nil, err
	}
	s.invalidateUserCounterCache(c)
	return s.buildUserResponse(created), nil
}

// GetUserById fetches a user by id, enriches the response with groups +
// memberships, and returns the DTO. 404 when the id is unknown.
func (s *UserService) GetUserById(c ctx.Context, userId *string) (*dtos.UserResponse, error) {
	user, err := s.fetchUserOr404(c, userId)
	if err != nil {
		return nil, err
	}
	response := s.buildUserResponse(user)
	response.Groups, response.GroupsCount = s.getUserGroups(c, *userId)
	response.Memberships = s.getUserMemberships(c, *userId)
	return response, nil
}

// UpdateUserById applies a partial user update, hashing the password when
// the patch carries one. 404 when no document matched the id.
func (s *UserService) UpdateUserById(c ctx.Context, userId *string, dto *dtos.UserUpdateDTO) (*dtos.UserResponse, error) {
	fields, err := s.buildUserUpdateFields(dto)
	if err != nil {
		return nil, err
	}
	updated, err := s.applyUserUpdate(c, userId, fields)
	if err != nil {
		return nil, err
	}
	return s.buildUserResponse(updated), nil
}

// DeleteUserById removes a user by id and invalidates the global counter
// cache. 404 when the id is unknown.
func (s *UserService) DeleteUserById(c ctx.Context, userId *string) (map[string]bool, error) {
	if err := s.deleteUserFromRepo(c, userId); err != nil {
		return nil, err
	}
	s.invalidateUserCounterCache(c)
	return map[string]bool{"success": true}, nil
}

// GetUserByEmail returns the user entity matched by email (used by auth
// for credential validation). Forwards repo errors as-is.
func (s *UserService) GetUserByEmail(c ctx.Context, email *string) (*entities.User, error) {
	user, err := s.deps.Repo.FindByEmail(c, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers returns the paginated, filtered user list scoped to the
// caller's org context (resolved via direct + group memberships).
func (s *UserService) GetUsers(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.UserQueryDto) (*model.PaginatedResult[dtos.UserResponse], error) {
	filters, empty, err := s.buildUserListFilters(c, requestContext, query)
	if err != nil {
		return nil, err
	}
	if empty {
		return s.emptyUserListResult(query), nil
	}
	pagination := s.buildUserListPagination(query)
	projection := s.buildUserListProjection(query)
	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}
	return s.mapUserListResultWithGroupCounts(c, result), nil
}

// CountUsers returns the per-org user count with cache-aside semantics
// (6h TTL); the counter cache is invalidated on every Create/Delete.
func (s *UserService) CountUsers(c ctx.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	orgId := s.userOrgIdFromContext(requestContext)
	cacheKey := s.deps.CounterCache.BuildKey(orgId)
	if count, ok := s.tryCachedUserCount(c, cacheKey); ok {
		return count, nil
	}
	count, err := s.countUsersFromRepo(c, requestContext)
	if err != nil {
		return 0, err
	}
	s.cacheUserCount(c, cacheKey, count)
	return count, nil
}

// fetchUserOr404 loads a user by id and returns NOT_FOUND when the repo
// returns nil. Used by every read/update orchestration.
func (s *UserService) fetchUserOr404(c ctx.Context, userId *string) (*entities.User, error) {
	user, err := s.deps.Repo.FindById(c, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"User not found"}}
	}
	return user, nil
}
