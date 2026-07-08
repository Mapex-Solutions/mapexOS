package constants

import "time"

/**
 * CACHE LIFETIME (application-level behavior)
 *
 * This TTL expresses "how long the application allows cached values to
 * remain authoritative" — it is application behavior, not a property of
 * any specific cache technology (any KV cache with TTL semantics can
 * enforce it). Redis key construction (prefix, key format) is an
 * infrastructure concern and lives in infrastructure/cache/redis.
 *
 * Cross-service FANOUT subjects/streams previously declared here now live
 * in packages/contracts/services/assets/assettemplates/constants.go for
 * cross-service contract reciprocity.
 */
const (
	// CounterCacheTTL is the lifetime for per-org asset template counter
	// cache entries. 6 hours.
	CounterCacheTTL = 6 * time.Hour

	// AssetTemplateContentTTL bounds how long a marketplace template's shared,
	// immutable content stays cached before a read re-fetches it from the
	// marketplace. Long — the content never changes for a given marketplaceGuid.
	AssetTemplateContentTTL = 24 * time.Hour
)

// AssetTemplateContentKeyPrefix namespaces the shared marketplace template
// content (the immutable bundle body) in the tiered cache, keyed by
// marketplaceGuid so every per-org install link hydrates one shared copy instead
// of duplicating the heavy body per organization.
const AssetTemplateContentKeyPrefix = "asset_template_content:"
