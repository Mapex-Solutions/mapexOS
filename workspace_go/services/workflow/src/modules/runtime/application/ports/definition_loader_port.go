package ports

import (
	"context"

	defPorts "workflow/src/modules/definitions/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// DefinitionLoaderPort provides cached access to WorkflowDefinition entities.
// Uses TieredCache (L0 RAM → L1 Disk → L2 MinIO) for fast lookups.
// Falls back to MongoDB when cache misses or orgId is unknown.
type DefinitionLoaderPort interface {
	// GetDefinition retrieves a WorkflowDefinition by ID.
	// When orgId is provided, attempts cache lookup first (key: {orgId}/{defId}).
	// When orgId is nil, goes directly to MongoDB and warms the cache with the result.
	GetDefinition(ctx context.Context, defId string, orgId *model.ObjectId) (*defPorts.WorkflowDefinition, error)
}
