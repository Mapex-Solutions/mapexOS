package services

import (
	"mapexIam/src/modules/roles/application/di"
)

// RoleService provides methods for managing role-related operations.
//
// Architecture Pattern: Dependency Injection
//   - Uses RoleServiceDependenciesInjection struct to aggregate all dependencies
//   - Follows same pattern as AuthService for consistency
type RoleService struct {
	deps di.RoleServiceDependenciesInjection
}
