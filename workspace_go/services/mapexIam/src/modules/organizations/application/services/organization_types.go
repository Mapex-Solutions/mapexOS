package services

import (
	"mapexIam/src/modules/organizations/application/di"
)

// OrganizationService provides methods for managing organization-related operations.
type OrganizationService struct {
	deps di.OrganizationServiceDependenciesInjection
}
