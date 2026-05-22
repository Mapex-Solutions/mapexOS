package services

import (
	"router/src/modules/routegroups/application/di"
)

// RouteGroupService provides methods for managing routegroup-related operations.
// It serves as an application service layer that interacts with the
// RouteGroupRepository to perform domain-level actions on RouteGroup entities.
//
// This service implements the RouteGroupServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type RouteGroupService struct {
	deps di.RouteGroupServiceDependenciesInjection
}
