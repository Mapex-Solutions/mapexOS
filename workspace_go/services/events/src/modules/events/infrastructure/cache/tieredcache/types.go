package tieredcache

import (
	"events/src/modules/events/domain/entities"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

/* Domain type aliases */

// Type aliases for domain entities used in this package.
// The canonical types live in domain/entities (DDD principle).
type DynamicField = entities.DynamicField
type CachedTemplate = entities.CachedTemplate

// TemplateCache provides cached access to AssetTemplate data for EVA field resolution.
// Wraps common.TieredCache with domain-specific unmarshal and field helpers.
type TemplateCache struct {
	cache common.TieredCache
}
