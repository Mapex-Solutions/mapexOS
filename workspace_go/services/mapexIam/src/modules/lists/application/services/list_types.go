package services

import (
	"mapexIam/src/modules/lists/application/di"
)

// ListService provides methods for managing list-related operations.
// It serves as an application service layer that interacts with the
// ListRepository to perform domain-level actions on List entities.
//
// Architecture Pattern: Dependency Injection
//   - Uses ListServiceDependenciesInjection struct to aggregate all dependencies
//   - Follows same pattern as AuthService for consistency
type ListService struct {
	deps di.ListServiceDependenciesInjection
}
