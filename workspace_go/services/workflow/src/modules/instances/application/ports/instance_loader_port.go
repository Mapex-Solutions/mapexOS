package ports

import (
	"context"

	"workflow/src/modules/instances/domain/entities"
)

// InstanceLoaderPort abstracts cached access to workflow instance configs.
// Wraps TieredCache (L0 RAM → L1 Disk) with MongoDB fallback.
type InstanceLoaderPort interface {
	// GetInstance retrieves an instance config by ID.
	// Cache key: "instance:{id}" — L0→L1→MongoDB fallback.
	GetInstance(ctx context.Context, instanceId string) (*entities.WorkflowInstance, error)

	// Invalidate removes an instance from L0+L1 cache.
	// Called after CRUD operations.
	Invalidate(ctx context.Context, instanceId string) error
}
