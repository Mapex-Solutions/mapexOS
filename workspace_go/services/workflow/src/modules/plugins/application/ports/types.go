package ports

import (
	"workflow/src/modules/plugins/domain/entities"
	"workflow/src/modules/plugins/domain/repositories"

	pluginsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/plugins"
)

// Port-level type aliases — expose domain entities through the port boundary.
// Other modules import these types from ports, NEVER from domain/entities directly.

type PluginManifest = entities.PluginManifest
type NodeTypeManifest = entities.NodeTypeManifest
type ActionDef = entities.ActionDef
type NodeHooks = entities.NodeHooks
type CredentialDef = entities.CredentialDef
type FetchOptionsDef = entities.FetchOptionsDef

// Repository alias — other modules access via this port type
type PluginManifestRepository = repositories.PluginManifestRepository

// PluginInvalidatePayload is the FANOUT message payload for plugin cache invalidation.
// Type alias to packages/contracts/services/workflow/plugins.PluginInvalidatePayload (canonical contract live in packages/contracts).
type PluginInvalidatePayload = pluginsContract.PluginInvalidatePayload
