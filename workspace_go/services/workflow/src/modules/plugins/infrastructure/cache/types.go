package cache

import (
	"workflow/src/modules/plugins/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// PluginLoader provides cached access to plugin manifests.
type PluginLoader struct {
	cache common.TieredCache
	repo  repositories.PluginManifestRepository
}
