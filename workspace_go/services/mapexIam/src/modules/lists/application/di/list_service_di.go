package di

import (
	"mapexIam/src/modules/lists/application/ports"
	"mapexIam/src/modules/lists/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// ListServiceDependenciesInjection defines the dependency injection container for ListService.
// This struct aggregates all dependencies required by ListService using dig.In.
//
// Architecture Pattern: Dependency Injection with Uber Dig
//   - dig.In: Instructs Dig to inject all fields automatically
//   - Provides clean constructor signature (single parameter instead of multiple)
//   - Scalable: Adding new dependencies doesn't change constructor signature
//
// Dependencies:
//   - Repo: Repository for list persistence operations
//   - NatsBus: For publishing list name update events
//
// Used by:
//   - ListService (as single dependency container)
type ListServiceDependenciesInjection struct {
	dig.In

	// Same domain repository
	Repo repositories.ListRepository

	// NatsBus provides NATS messaging capabilities (uses Publisher interface for Hexagonal Architecture)
	NatsBus natsModel.Publisher
}

// ListServiceDI is used to inject ListService into other services.
// This struct is used by consumers of ListService (e.g., other modules).
type ListServiceDI struct {
	dig.In

	ListService ports.ListServicePort
}
