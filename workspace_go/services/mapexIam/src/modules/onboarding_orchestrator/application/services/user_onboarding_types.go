package services

import (
	"mapexIam/src/modules/onboarding_orchestrator/application/di"
)

// UserOnboardingService is an Application Service that orchestrates user creation with memberships.
// This is a Use Case that coordinates multiple Domain Services (User, Membership) to achieve
// a complex business operation atomically.
//
// Architecture Pattern: Dependency Injection
//   - Uses UserOnboardingServiceDependenciesInjection struct to aggregate all dependencies
//   - Orchestrates multiple domain services following DDD principles
type UserOnboardingService struct {
	deps di.UserOnboardingServiceDependenciesInjection
}
