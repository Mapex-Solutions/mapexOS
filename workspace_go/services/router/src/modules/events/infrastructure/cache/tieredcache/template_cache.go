package tieredcache

import (
	"context"
	"encoding/json"
	"fmt"

	"router/src/modules/events/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * Template Cache for Event Enrichment
 *
 * Thin domain adapter over common.TieredCache (injected via DIG).
 * Infrastructure setup (MinIO client, L0/L1/L2 config) lives in bootstrap/cache.go.
 *
 * Used by EventService.buildEnrichedEvent() to resolve templateName and
 * templateDescription from AssetTemplateID + AssetTemplateOrgID.
 *
 * Architecture (configured in bootstrap/cache.go):
 *   - L0 (RAM): ~50µs latency, in-memory ristretto cache
 *   - L1 (Disk): ~500µs latency, local NVMe/SSD
 *   - L2 (MinIO): ~5-50ms latency, source of truth
 *   - Fallback (HTTP): Assets Service /internal/templates/:id
 */

// Compile-time check to ensure TemplateCache implements TemplateCachePort interface.
var _ ports.TemplateCachePort = (*TemplateCache)(nil)

// New creates a new TemplateCache wrapping the given TieredCache.
// The TieredCache is fully configured (L0/L1/L2/Fallback) and injected via DIG.
func New(cache common.TieredCache) ports.TemplateCachePort {
	return &TemplateCache{cache: cache}
}

// GetTemplate retrieves a template by orgId + templateId from cache.
// Cache key format: {templateOrgId}/{templateId}
// Falls back to "mapexos_public" if templateOrgId is empty (system template default).
func (tc *TemplateCache) GetTemplate(ctx context.Context, templateOrgId, templateId string) (*ports.CachedTemplate, error) {
	if templateId == "" {
		return nil, fmt.Errorf("templateId cannot be empty")
	}

	if templateOrgId == "" {
		templateOrgId = "mapexos_public"
	}
	cacheKey := templateOrgId + "/" + templateId

	data, _, err := tc.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s: %w", cacheKey, err)
	}

	var template ports.CachedTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template %s: %w", cacheKey, err)
	}

	logger.Debug(fmt.Sprintf("[INFRA:TemplateCache] Resolved key=%s name=%s", cacheKey, template.Name))
	return &template, nil
}

// Invalidate removes a template from local cache (L0 + L1).
// Called by the template_invalidate FANOUT consumer when a template is updated.
func (tc *TemplateCache) Invalidate(ctx context.Context, key string) error {
	return tc.cache.Invalidate(ctx, key)
}
