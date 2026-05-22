package cache

import (
	"context"
	"encoding/json"
	"fmt"

	defPorts "workflow/src/modules/definitions/application/ports"
	defRepos "workflow/src/modules/definitions/domain/repositories"
	"workflow/src/modules/runtime/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/*
 * DEFINITION LOADER
 * Wraps TieredCache + MongoDB DefinitionRepository for cached WorkflowDefinition access.
 *
 * Cache key: {orgId}/{definitionId} (matches MinIO path in L2)
 * L0 (RAM) → L1 (Disk) → L2 (MinIO) → fallback MongoDB
 * TTLs are managed by TieredCache bootstrap config (L0DefaultTTL / L1DefaultTTL).
 *
 * HandleExecution (newInstance): no orgId → MongoDB direct → warm cache
 * HandleResume: has orgId from instance → cache hit expected
 */

// Compile-time check
var _ ports.DefinitionLoaderPort = (*DefinitionLoader)(nil)

// New creates a DefinitionLoader with TieredCache and MongoDB fallback.
func New(cache common.TieredCache, repo defRepos.DefinitionRepository) ports.DefinitionLoaderPort {
	return &DefinitionLoader{cache: cache, repo: repo}
}

// GetDefinition retrieves a WorkflowDefinition by ID.
// When orgId is provided, attempts cache lookup first (key: {orgId}/{defId}).
// When orgId is nil, goes directly to MongoDB and warms the cache with the result.
func (l *DefinitionLoader) GetDefinition(ctx context.Context, defId string, orgId *model.ObjectId) (*defPorts.WorkflowDefinition, error) {
	// If orgId is available, try cache first
	if orgId != nil {
		key := orgId.Hex() + "/" + defId
		data, tier, err := l.cache.Get(ctx, key)
		if err == nil && len(data) > 0 {
			var def defPorts.WorkflowDefinition
			if unmarshalErr := json.Unmarshal(data, &def); unmarshalErr != nil {
				logger.Warn(fmt.Sprintf("[INFRA:DefinitionLoader] Cache hit (L%d) for %s but unmarshal failed: %s — falling back to MongoDB", tier, defId, unmarshalErr))
			} else {
				logger.Debug(fmt.Sprintf("[INFRA:DefinitionLoader] Cache hit (L%d) for %s", tier, defId))
				return &def, nil
			}
		}
	}

	// Cache miss or no orgId → fallback to MongoDB
	def, err := l.repo.FindById(ctx, &defId)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:DefinitionLoader] MongoDB FindById failed: %w", err)
	}
	if def == nil {
		return nil, nil
	}

	// Warm cache for future lookups
	l.warmCache(ctx, def)

	return def, nil
}

// warmCache serializes the definition and stores it in TieredCache.
func (l *DefinitionLoader) warmCache(ctx context.Context, def *defPorts.WorkflowDefinition) {
	if def.OrgID == nil {
		return
	}

	key := def.OrgID.Hex() + "/" + def.ID.Hex()
	data, err := json.Marshal(def)
	if err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:DefinitionLoader] Failed to marshal definition %s: %s", def.ID.Hex(), err))
		return
	}

	if err := l.cache.Set(ctx, key, data, 0); err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:DefinitionLoader] Failed to cache definition %s: %s", def.ID.Hex(), err))
	}
}
