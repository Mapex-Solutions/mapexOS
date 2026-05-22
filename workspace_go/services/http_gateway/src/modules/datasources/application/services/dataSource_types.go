package services

import (
	"http_gateway/src/modules/datasources/application/di"
)

// DataSourceService provides methods for managing dataSource-related operations.
// It serves as an application service layer that interacts with the
// DataSourceRepository to perform domain-level actions on DataSource entities.
//
// This service implements the DataSourceServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type DataSourceService struct {
	deps di.DataSourceServiceDependenciesInjection
}
