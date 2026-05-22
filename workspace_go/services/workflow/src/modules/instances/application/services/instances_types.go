package services

import (
	"workflow/src/modules/instances/application/di"
)

// InstancesService handles CRUD operations for workflow instance configs.
type InstancesService struct {
	deps di.InstancesServiceDependenciesInjection
}
