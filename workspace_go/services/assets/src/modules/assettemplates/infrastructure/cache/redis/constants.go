// Package redis holds Redis-specific constants and key-builder adapters for
// the assettemplates module's application cache (Redis DB 0 - AppCache).
//
// Scope: infrastructure/cache/redis — key prefixes and Redis-technology
// details that are intentionally not part of the domain, the cross-service
// contract, or the application-behavior layer. Cache TTLs live in
// application/constants (application behavior, technology-agnostic).
//
// Redis DB 0 roles covered:
//   - Per-org counter cache for asset template totals
package redis

/**
 * COUNTER CACHE KEY FORMAT
 *
 * Per-org asset template totals are cached to avoid MongoDB CountDocuments
 * on every list request. Keys are invalidated on create/delete.
 */
const (
	// CounterCacheKeyPrefix is the Redis key prefix for asset template
	// counter cache. Format: counter:asset_templates:{orgId}
	CounterCacheKeyPrefix = "counter:asset_templates:"
)
