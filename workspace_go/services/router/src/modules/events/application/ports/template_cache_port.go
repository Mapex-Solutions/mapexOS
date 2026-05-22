package ports

import (
	"context"
)

// CachedTemplate contains template metadata needed for event enrichment.
//
// Minimal view of AssetTemplate — only metadata needed for event enrichment.
// The full template (with DynamicFields, scripts) is deserialized but we only
// read Name and Description. JSON tags match the MinIO payload format.
type CachedTemplate struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TemplateCachePort defines the contract for cached access to template metadata
// used for event enrichment.
//
// This port enables Hexagonal Architecture by decoupling the application layer
// from the concrete template cache adapter. Infrastructure implementations
// (e.g., TieredCache-backed) satisfy this port.
//
// Methods:
//   - GetTemplate: Resolves template metadata (name, description) by orgId + templateId.
//   - Invalidate: Removes a template from local cache (L0 + L1) on external updates.
type TemplateCachePort interface {
	// GetTemplate retrieves a template by orgId + templateId from cache.
	// Cache key format: {templateOrgId}/{templateId}
	// Falls back to "mapexos_public" if templateOrgId is empty (system template default).
	GetTemplate(ctx context.Context, templateOrgId, templateId string) (*CachedTemplate, error)

	// Invalidate removes a template from local cache (L0 + L1).
	// Called by the template_invalidate FANOUT consumer when a template is updated.
	Invalidate(ctx context.Context, key string) error
}
