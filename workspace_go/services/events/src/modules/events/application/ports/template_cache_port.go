package ports

import (
	"context"

	"events/src/modules/events/domain/entities"
)

// TemplateCachePort defines the outbound port for accessing AssetTemplate data.
// The Application Service uses this interface to fetch templates for EVA field resolution,
// without knowing the caching strategy (L0/L1/L2/Fallback HTTP).
//
// Cache key format: {templateOrgId}/{templateId}
// S3 path: {templateOrgId}/{templateId}.json
type TemplateCachePort interface {
	GetTemplate(ctx context.Context, templateOrgId string, templateId string) (*entities.CachedTemplate, error)

	// Invalidate clears the local L0+L1 cache entry for the given key.
	// Key format: {orgId}/{templateId}. L2 (MinIO source of truth) is not affected.
	// Called by the FANOUT consumer when the assets service publishes
	// mapexos.fanout.template.invalidate.
	Invalidate(ctx context.Context, key string) error
}
