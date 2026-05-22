package services

import (
	ctx "context"

	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	utilsPassword "github.com/Mapex-Solutions/mapexGoKit/utils/bcrypt/password"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// ensureEmailIsUnique rejects creates that would collide with an existing
// account. Defense in depth — the email index also enforces uniqueness.
func (s *UserService) ensureEmailIsUnique(c ctx.Context, email string) error {
	if email == "" {
		return nil
	}
	existing, _ := s.deps.Repo.FindByEmail(c, &email)
	if existing != nil {
		return &customErrors.ServerCustomError{
			Code:   status.CONFLICT,
			Errors: []string{"User with this email already exists"},
		}
	}
	return nil
}

// buildUserCreateEntity maps the create DTO into a domain entity, applies
// the V1 defaults (startTour=true, internal AuthProvider), and hashes the
// password when one was supplied.
func (s *UserService) buildUserCreateEntity(dto *dtos.UserCreateDTO) (*entities.User, error) {
	userEntity, _ := mapper.DtoToEntity[dtos.UserCreateDTO, entities.User](dto)
	userEntity.StartTour = true
	userEntity.AuthProvider = entities.AuthProvider{
		Type:     "internal",
		Metadata: make(map[string]interface{}),
	}
	if userEntity.Password != nil && *userEntity.Password != "" {
		hashed, err := utilsPassword.HashPassword(*userEntity.Password)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   status.INTERNAL_SERVER_ERROR,
				Errors: []string{"Failed to hash password"},
			}
		}
		userEntity.Password = &hashed
	}
	return userEntity, nil
}

// invalidateUserCounterCache drops the global counter cache after a
// create or delete (Users are not org-scoped so the key uses "").
func (s *UserService) invalidateUserCounterCache(c ctx.Context) {
	counterKey := s.deps.CounterCache.BuildKey("")
	_ = s.deps.AppCache.Del(c, counterKey)
}

// buildUserResponse converts a User entity to its response DTO.
func (s *UserService) buildUserResponse(user *entities.User) *dtos.UserResponse {
	resp, _ := mapper.EntityToDto[entities.User, dtos.UserResponse](user)
	return resp
}

// buildUserUpdateFields converts the update DTO into a $set map and hashes
// the password when one was supplied.
func (s *UserService) buildUserUpdateFields(dto *dtos.UserUpdateDTO) (map[string]interface{}, error) {
	fields, _ := mapper.DtoToMap(dto)
	if dto.Password != nil && *dto.Password != "" {
		hashed, err := utilsPassword.HashPassword(*dto.Password)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   status.INTERNAL_SERVER_ERROR,
				Errors: []string{"Failed to hash password"},
			}
		}
		fields["password"] = hashed
	}
	return fields, nil
}

// applyUserUpdate runs the partial update against the repository and
// translates a missing document into the canonical 404 contract error.
func (s *UserService) applyUserUpdate(c ctx.Context, userId *string, fields map[string]interface{}) (*entities.User, error) {
	user, _ := s.deps.Repo.FindByIdAndUpdate(c, userId, fields)
	if user.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"User not found"}}
	}
	return user, nil
}

// deleteUserFromRepo removes the document and translates the driver's
// "document not found" string into the canonical 404 contract error.
func (s *UserService) deleteUserFromRepo(c ctx.Context, userId *string) error {
	if err := s.deps.Repo.DeleteById(c, userId); err != nil {
		if err.Error() == "document not found" {
			return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"User not found"}}
		}
		return err
	}
	return nil
}
