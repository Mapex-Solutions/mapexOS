package ports

import (
	"context"

	"workflow/src/modules/plugins/domain/entities"
)

// PluginLoaderPort abstracts cached access to plugin manifests.
// Wraps TieredCache (L0 RAM → L1 Disk) with MongoDB fallback.
type PluginLoaderPort interface {
	// GetManifest retrieves a single manifest by pluginId.
	// Cache key: "plugin:{pluginId}" — L0→L1→MongoDB fallback.
	GetManifest(ctx context.Context, pluginId string) (*entities.PluginManifest, error)

	// GetAllEnabled retrieves all enabled manifests (for editor boot).
	// Cache key: "plugins:all:enabled" — shorter TTL.
	GetAllEnabled(ctx context.Context) ([]entities.PluginManifest, error)

	// Invalidate removes a pluginId from L0+L1 cache.
	// Called after CRUD operations and on NATS fanout receive.
	Invalidate(ctx context.Context, pluginId string) error

	// InvalidateAll removes the "all enabled" cache entry.
	// Called on any CRUD operation that changes enabled state.
	InvalidateAll(ctx context.Context) error
}
