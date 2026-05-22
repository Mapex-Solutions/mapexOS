package di

import (
	groupPorts "mapexIam/src/modules/groups/application/ports"
	membershipPorts "mapexIam/src/modules/memberships/application/ports"
	onboardingPorts "mapexIam/src/modules/onboarding_orchestrator/application/ports"
	orgPorts "mapexIam/src/modules/organizations/application/ports"
	userPorts "mapexIam/src/modules/users/application/ports"

	"go.uber.org/dig"
)

// UserOnboardingServiceDependenciesInjection defines the dependency injection container for UserOnboardingService.
// This struct aggregates all dependencies required by UserOnboardingService using dig.In.
//
// Architecture Pattern: Dependency Injection with Uber Dig
//   - dig.In: Instructs Dig to inject all fields automatically
//   - Provides clean constructor signature (single parameter instead of multiple)
//   - Scalable: Adding new dependencies doesn't change constructor signature
//
// DDD/Hexagonal Architecture:
//   - Application Service: Orchestrates multiple domain services
//   - Uses service ports (NOT repositories) to maintain bounded context
//   - All cross-domain access via service ports following DDD principles
//   - Infrastructure access (MongoDB transactions) via a dedicated port to
//     keep the concrete *mongoManager.MongoManager driver out of the
//     application DI layer.
//
// Used by:
//   - UserOnboardingService (as single dependency container)
type UserOnboardingServiceDependenciesInjection struct {
	dig.In

	// Infrastructure port — provides transaction support via RunTransaction()
	MongoManager onboardingPorts.MongoManagerPort

	// Domain services (using ports for DDD bounded context)
	UserService       userPorts.UserServicePort
	MembershipService membershipPorts.MembershipServicePort
	OrgService        orgPorts.OrganizationServicePort
	GroupService      groupPorts.GroupServicePort
}
