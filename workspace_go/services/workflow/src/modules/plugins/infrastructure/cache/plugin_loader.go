package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/plugins/domain/entities"
	"workflow/src/modules/plugins/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/*
 * PLUGIN LOADER
 * Wraps TieredCache + MongoDB PluginManifestRepository for cached access.
 *
 * Cache key formats:
 *   - Single: "plugin:{pluginId}"
 *   - All enabled: "plugins:all:enabled"
 *
 * TTLs are managed by TieredCache bootstrap config (L0DefaultTTL / L1DefaultTTL).
 * L0 (RAM) → L1 (Disk) → MongoDB fallback (no L2/MinIO — MongoDB is SoT)
 */

// Compile-time check
var _ ports.PluginLoaderPort = (*PluginLoader)(nil)

// New creates a PluginLoader with TieredCache and MongoDB fallback.
func New(cache common.TieredCache, repo repositories.PluginManifestRepository) ports.PluginLoaderPort {
	return &PluginLoader{cache: cache, repo: repo}
}

// GetManifest retrieves a single manifest by pluginId.
// Tries L0→L1 cache first, falls back to MongoDB on miss.
func (l *PluginLoader) GetManifest(ctx context.Context, pluginId string) (*entities.PluginManifest, error) {
	key := keyPrefix + pluginId

	// Try cache
	data, tier, err := l.cache.Get(ctx, key)
	if err == nil && len(data) > 0 {
		var manifest entities.PluginManifest
		if unmarshalErr := json.Unmarshal(data, &manifest); unmarshalErr != nil {
			logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Cache hit (L%d) for %s but unmarshal failed: %s — falling back to MongoDB", tier, pluginId, unmarshalErr))
		} else {
			logger.Debug(fmt.Sprintf("[INFRA:PluginLoader] Cache hit (L%d) for %s", tier, pluginId))
			return &manifest, nil
		}
	}

	// Cache miss → MongoDB
	manifest, err := l.repo.FindByPluginId(ctx, pluginId)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:PluginLoader] MongoDB FindByPluginId failed: %w", err)
	}
	if manifest == nil {
		return nil, nil
	}

	// Warm cache
	l.warmSingle(ctx, manifest)

	return manifest, nil
}

// GetAllEnabled retrieves all enabled manifests (for editor boot).
func (l *PluginLoader) GetAllEnabled(ctx context.Context) ([]entities.PluginManifest, error) {
	// Try cache
	data, tier, err := l.cache.Get(ctx, allEnabledKey)
	if err == nil && len(data) > 0 {
		var manifests []entities.PluginManifest
		if unmarshalErr := json.Unmarshal(data, &manifests); unmarshalErr != nil {
			logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Cache hit (L%d) for all:enabled but unmarshal failed: %s", tier, unmarshalErr))
		} else {
			logger.Debug(fmt.Sprintf("[INFRA:PluginLoader] Cache hit (L%d) for all:enabled (%d manifests)", tier, len(manifests)))
			return manifests, nil
		}
	}

	// Cache miss → MongoDB
	filters := model.Map{"enabled": true}
	result, err := l.repo.FindWithFilters(ctx, filters, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:PluginLoader] MongoDB FindWithFilters failed: %w", err)
	}

	manifests := result.Items

	// Warm cache
	l.warmAll(ctx, manifests)

	return manifests, nil
}

// Invalidate removes a single pluginId from cache.
func (l *PluginLoader) Invalidate(ctx context.Context, pluginId string) error {
	key := keyPrefix + pluginId
	if err := l.cache.Delete(ctx, key); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Failed to invalidate cache for %s: %s", pluginId, err))
		return err
	}
	logger.Debug(fmt.Sprintf("[INFRA:PluginLoader] Cache invalidated for %s", pluginId))
	return nil
}

// InvalidateAll removes the "all enabled" cache entry.
func (l *PluginLoader) InvalidateAll(ctx context.Context) error {
	if err := l.cache.Delete(ctx, allEnabledKey); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Failed to invalidate all:enabled cache: %s", err))
		return err
	}
	logger.Debug("[INFRA:PluginLoader] Cache invalidated for all:enabled")
	return nil
}

// warmSingle caches a single manifest. TTL=0 uses bootstrap defaults.
func (l *PluginLoader) warmSingle(ctx context.Context, manifest *entities.PluginManifest) {
	key := keyPrefix + manifest.PluginID
	data, err := json.Marshal(manifest)
	if err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Failed to marshal manifest %s: %s", manifest.PluginID, err))
		return
	}
	if err := l.cache.Set(ctx, key, data, 0); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Failed to cache manifest %s: %s", manifest.PluginID, err))
	}
}

// warmAll caches the "all enabled" list. TTL=0 uses bootstrap defaults.
func (l *PluginLoader) warmAll(ctx context.Context, manifests []entities.PluginManifest) {
	data, err := json.Marshal(manifests)
	if err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Failed to marshal all:enabled: %s", err))
		return
	}
	if err := l.cache.Set(ctx, allEnabledKey, data, 0); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:PluginLoader] Failed to cache all:enabled: %s", err))
	}
}
