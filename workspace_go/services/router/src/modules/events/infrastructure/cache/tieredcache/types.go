package tieredcache

import (
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// TemplateCache provides cached access to template metadata for event enrichment.
// Wraps common.TieredCache with domain-specific key building and unmarshal.
type TemplateCache struct {
	cache common.TieredCache
}
