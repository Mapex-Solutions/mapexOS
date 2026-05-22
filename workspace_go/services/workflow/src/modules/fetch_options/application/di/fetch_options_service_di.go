package di

import (
	pluginPorts "workflow/src/modules/plugins/application/ports"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	"go.uber.org/dig"
)

// FetchOptionsServiceDI aggregates all dependencies required by the FetchOptionsService.
type FetchOptionsServiceDI struct {
	dig.In

	Vault      runtimePorts.VaultPort
	PluginRepo pluginPorts.PluginManifestRepository
}
