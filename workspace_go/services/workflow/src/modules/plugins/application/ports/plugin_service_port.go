package ports

import (
	ctx "context"

	"workflow/src/modules/plugins/application/dtos"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// PluginServicePort defines the contract for plugin manifest business operations.
type PluginServicePort interface {
	// CreatePlugin creates a new plugin manifest.
	CreatePlugin(ctx ctx.Context, requestContext *reqCtx.RequestContext, entity *dtos.PluginManifestResponse) (*dtos.PluginManifestResponse, error)

	// GetPluginById retrieves a plugin manifest by its MongoDB ObjectId.
	GetPluginById(ctx ctx.Context, pluginId *string) (*dtos.PluginManifestResponse, error)

	// GetPluginByPluginId retrieves a plugin manifest by its unique pluginId string.
	// Uses TieredCache (L0→L1→MongoDB fallback) via PluginLoader.
	GetPluginByPluginId(ctx ctx.Context, pluginId string) (*dtos.PluginManifestResponse, error)

	// UpdatePluginById updates an existing plugin manifest.
	UpdatePluginById(ctx ctx.Context, pluginId *string, dto *dtos.PluginManifestUpdate) (*dtos.PluginManifestResponse, error)

	// DeletePluginById removes a plugin manifest.
	DeletePluginById(ctx ctx.Context, pluginId *string) (map[string]bool, error)

	// GetPlugins retrieves a paginated and filtered list of plugin manifests.
	// Respects multi-tenant visibility: system + template + local (org).
	GetPlugins(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.PluginQueryDTO) (*model.PaginatedResult[dtos.PluginManifestResponse], error)

	// GetEnabledPlugins retrieves all enabled plugins for the editor boot sequence.
	// Uses TieredCache (L0→L1→MongoDB fallback) via PluginLoader.
	GetEnabledPlugins(ctx ctx.Context) ([]dtos.PluginManifestResponse, error)

	// HandleFanoutEvent processes a FANOUT cache-invalidation message broadcast across pods.
	// Decodes the payload and invalidates the local TieredCache (L0/L1) for the affected plugin.
	HandleFanoutEvent(msg *natsModel.Message)
}
