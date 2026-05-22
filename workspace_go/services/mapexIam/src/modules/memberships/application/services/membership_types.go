package services

import (
	"mapexIam/src/modules/memberships/application/di"
)

// MembershipService provides methods for managing membership-related operations.
//
// Architecture Pattern: Dependency Injection
//   - Uses MembershipServiceDependenciesInjection struct to aggregate all dependencies
//   - Follows same pattern as AuthService for consistency
type MembershipService struct {
	deps di.MembershipServiceDependenciesInjection
}
