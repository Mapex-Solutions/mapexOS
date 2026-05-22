package di

import (
	"mapexIam/src/modules/organizations/application/ports"
	"mapexIam/src/modules/organizations/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// OrganizationServiceDependenciesInjection aggregates the dependencies that
// OrganizationService needs so the constructor stays a single-parameter
// signature even as the dependency list grows.
type OrganizationServiceDependenciesInjection struct {
	dig.In

	// Same domain repository
	Repo repositories.OrganizationRepository

	// Event bus for domain event publishing (uses Publisher interface for Hexagonal Architecture)
	NatsBus natsModel.Publisher
}

// OrganizationServiceDI is used to inject OrganizationService into other services.
// This struct is used by consumers of OrganizationService (e.g., other modules).
type OrganizationServiceDI struct {
	dig.In

	OrganizationService ports.OrganizationServicePort
}
