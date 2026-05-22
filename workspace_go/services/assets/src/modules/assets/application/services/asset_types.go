package services

import (
	"assets/src/modules/assets/application/di"
)

// AssetService provides methods for managing asset-related operations.
// It serves as an application service layer that interacts with the
// AssetRepository to perform domain-level actions on Asset entities.
//
// This service implements the AssetServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type AssetService struct {
	deps di.AssetServiceDependenciesInjection
}
