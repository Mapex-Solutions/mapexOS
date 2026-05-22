package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"workflow/src/modules/instances/application/ports"
	"workflow/src/modules/instances/domain/entities"
	"workflow/src/modules/instances/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/*
 * INSTANCE LOADER
 * Wraps TieredCache + MongoDB InstanceRepository for cached access.
 *
 * Cache key format:
 *   - Single: "instance:{id}"
 *
 * TTLs are managed by TieredCache bootstrap config (L0DefaultTTL / L1DefaultTTL).
 * L0 (RAM) → L1 (Disk) → MongoDB fallback (no L2/MinIO — MongoDB is SoT)
 */

// Compile-time check
var _ ports.InstanceLoaderPort = (*InstanceLoader)(nil)

// New creates an InstanceLoader with TieredCache and MongoDB fallback.
func New(cache common.TieredCache, repo repositories.InstanceRepository) ports.InstanceLoaderPort {
	return &InstanceLoader{cache: cache, repo: repo}
}

// GetInstance retrieves a single instance config by ID.
// Tries L0→L1 cache first, falls back to MongoDB on miss.
func (l *InstanceLoader) GetInstance(ctx context.Context, instanceId string) (*entities.WorkflowInstance, error) {
	key := keyPrefix + instanceId

	// Try cache
	data, tier, err := l.cache.Get(ctx, key)
	if err == nil && len(data) > 0 {
		var instance entities.WorkflowInstance
		if unmarshalErr := json.Unmarshal(data, &instance); unmarshalErr != nil {
			logger.Warn(fmt.Sprintf("[INFRA:InstanceLoader] Cache hit (L%d) for %s but unmarshal failed: %s — falling back to MongoDB", tier, instanceId, unmarshalErr))
		} else {
			logger.Debug(fmt.Sprintf("[INFRA:InstanceLoader] Cache hit (L%d) for %s", tier, instanceId))
			return &instance, nil
		}
	}

	// Cache miss → MongoDB
	instance, err := l.repo.FindById(ctx, &instanceId)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:InstanceLoader] MongoDB FindById failed: %w", err)
	}
	if instance == nil {
		return nil, nil
	}

	// Warm cache
	l.warm(ctx, instanceId, instance)

	return instance, nil
}

// Invalidate removes an instance from L0+L1 cache.
func (l *InstanceLoader) Invalidate(ctx context.Context, instanceId string) error {
	key := keyPrefix + instanceId
	if err := l.cache.Delete(ctx, key); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:InstanceLoader] Failed to invalidate cache for %s: %s", instanceId, err))
		return err
	}
	logger.Debug(fmt.Sprintf("[INFRA:InstanceLoader] Cache invalidated for %s", instanceId))
	return nil
}

// warm caches a single instance. TTL=0 uses bootstrap defaults.
func (l *InstanceLoader) warm(ctx context.Context, instanceId string, instance *entities.WorkflowInstance) {
	key := keyPrefix + instanceId
	data, err := json.Marshal(instance)
	if err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:InstanceLoader] Failed to marshal instance %s: %s", instanceId, err))
		return
	}
	if err := l.cache.Set(ctx, key, data, 0); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:InstanceLoader] Failed to cache instance %s: %s", instanceId, err))
	}
}
