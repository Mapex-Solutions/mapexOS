package cache

import (
	defRepos "workflow/src/modules/definitions/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

type DefinitionLoader struct {
	cache common.TieredCache
	repo  defRepos.DefinitionRepository
}
