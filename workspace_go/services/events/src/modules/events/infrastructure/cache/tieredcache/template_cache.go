package tieredcache

import (
	"context"
	"encoding/json"
	"fmt"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/*
 * Template Cache for EVA Field Resolution
 *
 * This cache is a thin domain adapter over common.TieredCache (injected via DIG).
 * Infrastructure setup (S3 client, L0/L1/L2 config) lives in main.go.
 *
 * It stores AssetTemplate data for resolving EVA field mappings:
 * field_name → fieldId for efficient ClickHouse MAP<UInt16, Type> storage.
 *
 * Architecture (configured in main.go):
 *   - L0 (RAM): ~50µs latency, in-memory ristretto cache
 *   - L1 (Disk): ~500µs latency, local NVMe/SSD
 *   - L2 (S3/MinIO): ~5-50ms latency, source of truth
 *   - Fallback (HTTP): Calls Assets Service /internal/templates/:id
 */

// New creates a new TemplateCache wrapping the given TieredCache.
// The TieredCache is fully configured (L0/L1/L2/Fallback) and injected via DIG.
func New(cache common.TieredCache) *TemplateCache {
	return &TemplateCache{cache: cache}
}

// GetTemplate retrieves a template by orgId + templateId from cache.
// Returns the cached template or fetches from Assets Service if not cached.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - templateOrgId: The org that owns the template ("mapexos_public" for system templates)
//   - templateId: The AssetTemplate ID (MongoDB ObjectId as string)
//
// Cache key format: {templateOrgId}/{templateId}
// S3 path: {templateOrgId}/{templateId}.json
//
// Returns:
//   - *CachedTemplate: The template with DynamicFields for EVA mapping
//   - error: If template not found or fetch failed
func (tc *TemplateCache) GetTemplate(ctx context.Context, templateOrgId string, templateId string) (*CachedTemplate, error) {
	if templateId == "" {
		return nil, fmt.Errorf("templateId cannot be empty")
	}

	// Build cache key: {templateOrgId}/{templateId}
	// Falls back to "mapexos_public" if templateOrgId is empty (system template default)
	if templateOrgId == "" {
		templateOrgId = "mapexos_public"
	}
	cacheKey := templateOrgId + "/" + templateId

	// Get from cache (L0 → L1 → L2 → Fallback HTTP)
	data, tier, err := tc.cache.Get(ctx, cacheKey)
	if err != nil {
		stats := tc.cache.Stats()
		logger.Warn(fmt.Sprintf(
			"[CACHE:Template] MISS key=%s err=%v | stats: L0(hits=%d miss=%d) L1(hits=%d miss=%d) L2(hits=%d miss=%d) Fallback(hits=%d miss=%d)",
			cacheKey, err,
			stats.L0Hits, stats.L0Misses,
			stats.L1Hits, stats.L1Misses,
			stats.L2Hits, stats.L2Misses,
			stats.FallbackHits, stats.FallbackMisses,
		))
		return nil, fmt.Errorf("template not found: %s: %w", cacheKey, err)
	}

	// Log which tier served the data
	tierName := tierToString(tier)
	logger.Info(fmt.Sprintf("[CACHE:Template] HIT key=%s tier=%s size=%d bytes", cacheKey, tierName, len(data)))

	// Unmarshal JSON to CachedTemplate
	var template CachedTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		logger.Error(err, fmt.Sprintf("[CACHE:Template] Unmarshal failed key=%s dataLen=%d", cacheKey, len(data)))
		return nil, fmt.Errorf("failed to unmarshal template %s: %w", cacheKey, err)
	}

	logger.Info(fmt.Sprintf("[CACHE:Template] Resolved key=%s fields=%d", cacheKey, len(template.DynamicFields)))
	return &template, nil
}

// tierToString converts a tier integer to a human-readable string.
func tierToString(tier int) string {
	switch tier {
	case 0:
		return "L0(RAM)"
	case 1:
		return "L1(Disk)"
	case 2:
		return "L2(S3)"
	case 3:
		return "Fallback(HTTP)"
	default:
		return fmt.Sprintf("Unknown(%d)", tier)
	}
}

// GetActiveFields returns only active DynamicFields (status=1) for a template.
// This is useful when processing new events - deprecated fields should not be populated.
func (tc *TemplateCache) GetActiveFields(ctx context.Context, templateOrgId string, templateId string) ([]DynamicField, error) {
	template, err := tc.GetTemplate(ctx, templateOrgId, templateId)
	if err != nil {
		return nil, err
	}

	var activeFields []DynamicField
	for _, field := range template.DynamicFields {
		if field.Status == 1 {
			activeFields = append(activeFields, field)
		}
	}

	return activeFields, nil
}

// BuildFieldIndex returns a map of field_name → fieldId for quick lookup.
// Only includes active fields (status=1).
func (tc *TemplateCache) BuildFieldIndex(ctx context.Context, templateOrgId string, templateId string) (map[string]uint16, error) {
	fields, err := tc.GetActiveFields(ctx, templateOrgId, templateId)
	if err != nil {
		return nil, err
	}

	index := make(map[string]uint16, len(fields))
	for _, field := range fields {
		index[field.Field] = field.FieldId
	}

	return index, nil
}

// Invalidate removes a template from local cache (L0 + L1).
// Called when template is updated via NATS broadcast.
func (tc *TemplateCache) Invalidate(ctx context.Context, templateId string) error {
	return tc.cache.Invalidate(ctx, templateId)
}

// Stats returns cache statistics for monitoring.
func (tc *TemplateCache) Stats() map[string]uint64 {
	stats := tc.cache.Stats()
	return map[string]uint64{
		"l0_hits":         stats.L0Hits,
		"l0_misses":       stats.L0Misses,
		"l1_hits":         stats.L1Hits,
		"l1_misses":       stats.L1Misses,
		"l2_hits":         stats.L2Hits,
		"l2_misses":       stats.L2Misses,
		"fallback_hits":   stats.FallbackHits,
		"fallback_misses": stats.FallbackMisses,
	}
}
