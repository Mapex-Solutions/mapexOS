package services

import (
	"events/src/modules/retention/application/di"
)

// RetentionService provides methods for managing retention policy operations.
// It serves as an application service layer that interacts with the
// RetentionRepository to perform domain-level actions on RetentionPolicy entities.
//
// This service implements the RetentionServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type RetentionService struct {
	deps di.RetentionServiceDependenciesInjection
}
