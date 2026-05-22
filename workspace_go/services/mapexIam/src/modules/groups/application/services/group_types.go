package services

import (
	"mapexIam/src/modules/groups/application/di"
)

// GroupService provides methods for managing group-related operations.
//
// Architecture Pattern: Dependency Injection
//   - Uses GroupServiceDependenciesInjection struct to aggregate all dependencies
//   - Follows same pattern as AuthService for consistency
type GroupService struct {
	deps di.GroupServiceDependenciesInjection
}
