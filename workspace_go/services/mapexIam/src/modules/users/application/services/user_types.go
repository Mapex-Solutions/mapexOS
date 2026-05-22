package services

import (
	"mapexIam/src/modules/users/application/di"
)

// UserService provides methods for managing user-related operations.
// It serves as an application service layer that interacts with the
// UserRepository to perform domain-level actions on User entities.
//
// Architecture Pattern: Dependency Injection
//   - Uses UserServiceDependenciesInjection struct to aggregate all dependencies
//   - Follows same pattern as AuthService for consistency
type UserService struct {
	deps di.UserServiceDependenciesInjection
}
