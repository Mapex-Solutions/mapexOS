package services

import (
	"triggers/src/modules/triggers/application/di"
)

// TriggerService provides methods for managing trigger-related operations.
// It serves as an application service layer that interacts with the
// TriggerRepository to perform domain-level actions on Trigger entities.
//
// This service implements the TriggerServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type TriggerService struct {
	deps di.TriggerServiceDependenciesInjection
}
