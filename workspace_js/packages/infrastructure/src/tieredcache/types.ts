/**
 * TieredCache Types
 * Same structure as workspace_go/packages/infrastructure/tieredcache/types.go
 */

/**
 * Config holds the TieredCache configuration.
 */
export interface Config {
	// L0 (RAM) Configuration
	EnableL0: boolean; // Enable L0 RAM cache (default: true)
	L0MaxSize: number; // Maximum size in bytes for L0 cache
	L0MaxItems: number; // Maximum number of items in L0
	L0DefaultTTL: number; // Default TTL for L0 entries (ms)

	// L1 (Disk) Configuration
	EnableL1: boolean; // Enable L1 disk cache (default: true)
	L1Dir: string; // Directory for L1 disk cache
	L1MaxSize: number; // Maximum size in bytes for L1 cache
	L1DefaultTTL: number; // Default TTL for L1 entries (ms, 0 = no TTL)

	// L2 (Remote) Configuration
	EnableL2?: boolean; // Enable L2 remote loader (MinIO/S3)
	L2Loader?: LocalCacheLoader; // Loader function for L2 (required if EnableL2 is true)

	// Fallback Configuration (optional)
	// Used when L2 (MinIO) misses - calls HTTP API to fetch from source
	FallbackBaseURL?: string; // Base URL of source service (e.g., "http://assets-service:5001")
	FallbackAPIKey?: string; // API Key for authentication
	FallbackEndpoint?: string; // Endpoint path (e.g., "/internal/assets")
	FallbackTimeout?: number; // Request timeout in ms (default: 5000)
	FallbackKeyTransformer?: (key: string) => string; // Transforms cache key to URL path for fallback (optional)

	// General Configuration
	KeyPrefix: string; // Prefix for all cache keys
	EnableMetrics: boolean; // Enable detailed metrics collection
}

/**
 * CacheStats holds atomic cache statistics.
 */
export interface CacheStats {
	L0Hits: number;
	L0Misses: number;
	L0Size: number;

	L1Hits: number;
	L1Misses: number;
	L1Size: number;

	L2Hits: number;
	L2Misses: number;

	FallbackHits: number; // HTTP fallback hits when L2 misses
	FallbackMisses: number; // HTTP fallback misses

	// Lazy cleanup statistics (on-read expiration removal)
	L1LazyExpired: number; // Files removed during read (TTL expired)
}

/**
 * CacheEntry represents a cached item with metadata.
 */
export interface CacheEntry {
	Data?: Buffer;
	ExpiresAt: Date;
	CreatedAt: Date;
	Size: number;
}

/**
 * CacheTier represents which cache tier served the data.
 */
export enum CacheTier {
	// TierMiss indicates no cache hit
	Miss = -1,
	// TierL0 indicates hit from RAM cache
	L0 = 0,
	// TierL1 indicates hit from disk cache
	L1 = 1,
	// TierL2 indicates hit from remote storage (MinIO)
	L2 = 2,
	// TierFallback indicates hit from HTTP fallback API
	Fallback = 3,
}

/**
 * Returns a human-readable tier name.
 */
export function tierToString(tier: CacheTier): string {
	switch (tier) {
		case CacheTier.L0:
			return 'L0-RAM';
		case CacheTier.L1:
			return 'L1-Disk';
		case CacheTier.L2:
			return 'L2-Remote';
		case CacheTier.Fallback:
			return 'Fallback-HTTP';
		default:
			return 'MISS';
	}
}

/**
 * LocalCacheLoader is a function that loads data from L2 (source of truth).
 * Used when data is not found in L0 or L1.
 */
export type LocalCacheLoader = (key: string) => Promise<Buffer | null>;

/**
 * L0Entry represents an entry in the L0 RAM cache.
 */
export interface L0Entry {
	Data: Buffer;
	ExpiresAt: number; // Unix timestamp (ms)
	Size: number;
}
