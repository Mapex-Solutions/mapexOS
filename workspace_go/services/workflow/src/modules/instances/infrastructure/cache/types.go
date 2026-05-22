package cache

import (
	"workflow/src/modules/instances/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// InstanceLoader provides cached access to workflow instance configs.
type InstanceLoader struct {
	cache common.TieredCache
	repo  repositories.InstanceRepository
}
