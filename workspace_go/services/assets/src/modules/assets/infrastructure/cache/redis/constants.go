// Package redis holds Redis-specific constants and key-builder adapters for
// the assets module's application cache (Redis DB 0 - AppCache).
//
// Scope: infrastructure/cache/redis — key prefixes and Redis-technology
// details that are intentionally not part of the domain, the cross-service
// contract, or the application-behavior layer. Cache TTLs live in
// application/constants (application behavior, technology-agnostic).
//
// Redis DB 0 role covered:
//   - Per-org counter cache for asset totals
//
// MQTT auth used to live here as a credential cache; that path was
// removed because the mapex-mqtt-broker plugin owns its own TieredCache
// (Pebble L1 + MinIO L2 + HTTP L3) and decides CONNECTs locally off the
// AssetReadModel.
package redis

/**
 * COUNTER CACHE KEY FORMAT
 *
 * Per-org asset totals are cached to avoid MongoDB CountDocuments on every
 * list request. Keys are invalidated on create/delete.
 */
const (
	// CounterCacheKeyPrefix is the Redis key prefix for asset counter cache.
	// Format: counter:assets:{orgId}
	CounterCacheKeyPrefix = "counter:assets:"
)
