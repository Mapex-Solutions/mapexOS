/**
 * TieredCache Methods
 * Same structure as workspace_go/packages/infrastructure/tieredcache/methods.go
 */

import type { Config, CacheStats, LocalCacheLoader, L0Entry } from './types';
import { CacheTier } from './types';
import { ErrCacheMiss, ErrEmptyKey, ErrNilData, ErrNoLoader } from './errors';
import {
	prefixKey,
	writeL1,
	readL1,
	deleteL1Files,
	isL0Enabled,
	isL1Enabled,
} from './internals';
import { HTTPClient } from '../httpclient';
import { getInfraLogger } from '../logger';

/**
 * TieredCacheClient implements a multi-tier cache system.
 *
 * Architecture:
 *   - L0 (RAM): Map-based in-memory cache (~50µs latency)
 *   - L1 (Disk): Local NVMe/SSD storage (~500µs latency)
 *   - L2 (Remote): MinIO/S3 source of truth (~5-50ms latency)
 *   - Fallback (HTTP): API call to source service when L2 misses (~10-100ms latency)
 */
export class TieredCacheClient {
	// L0 - In-memory cache (Map)
	private l0Cache: Map<string, L0Entry> | null = null;

	// L1 - Disk cache directory
	private l1Dir: string = '';

	// L2 - Remote loader function (MinIO)
	private l2Loader: LocalCacheLoader | null = null;

	// Fallback - HTTP client for L2 miss recovery
	private fallbackClient: HTTPClient | null = null;
	private fallbackEndpoint: string = '';
	private fallbackKeyTransformer: ((key: string) => string) | null = null;

	// Configuration
	private config: Config;

	// Statistics
	private stats: CacheStats = {
		L0Hits: 0,
		L0Misses: 0,
		L0Size: 0,
		L1Hits: 0,
		L1Misses: 0,
		L1Size: 0,
		L2Hits: 0,
		L2Misses: 0,
		FallbackHits: 0,
		FallbackMisses: 0,
		L1LazyExpired: 0,
	};

	// Key prefix for all cache operations
	private keyPrefix: string;

	constructor(
		config: Config,
		l0Cache: Map<string, L0Entry> | null,
		l1Dir: string,
		l2Loader: LocalCacheLoader | null = null,
		fallbackClient: HTTPClient | null = null,
		fallbackEndpoint: string = '',
		fallbackKeyTransformer: ((key: string) => string) | null = null,
	) {
		this.config = config;
		this.l0Cache = l0Cache;
		this.l1Dir = l1Dir;
		this.l2Loader = l2Loader;
		this.keyPrefix = config.KeyPrefix;
		this.fallbackClient = fallbackClient;
		this.fallbackEndpoint = fallbackEndpoint;
		this.fallbackKeyTransformer = fallbackKeyTransformer;
	}

	/**
	 * Get retrieves a value from cache following the tier hierarchy: L0 → L1 → L2.
	 *
	 * Returns:
	 *   - data: the cached bytes (null if not found)
	 *   - tier: which tier served the data (0=L0, 1=L1, 2=L2, -1=miss)
	 *   - error: throws ErrCacheMiss if not found in any tier
	 */
	async Get(key: string): Promise<{ data: Buffer; tier: CacheTier }> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		// L0 (RAM) lookup - only if enabled
		if (isL0Enabled(this.config, this.l0Cache)) {
			const result = this.GetFromL0(prefixed);
			if (result.found) {
				this.stats.L0Hits++;
				return { data: result.data!, tier: CacheTier.L0 };
			}
			this.stats.L0Misses++;
		}

		// L1 (Disk) lookup
		if (isL1Enabled(this.config, this.l1Dir)) {
			try {
				const data = this.GetFromL1(prefixed);
				this.stats.L1Hits++;
				// Promote to L0 if enabled
				this.setL0(prefixed, data, this.config.L0DefaultTTL);
				return { data, tier: CacheTier.L1 };
			} catch {
				this.stats.L1Misses++;
			}
		}

		// L2 (Remote) lookup via loader
		if (this.l2Loader) {
			const data = await this.l2Loader(key); // Use original key for L2
			if (data) {
				this.stats.L2Hits++;
				// Promote to L0 (if enabled) and L1 (if enabled)
				this.setL0(prefixed, data, this.config.L0DefaultTTL);
				if (isL1Enabled(this.config, this.l1Dir)) {
					try {
						writeL1(this.l1Dir, prefixed, data, this.config.L1DefaultTTL, this.stats);
					} catch {
						// Log but don't fail
					}
				}
				return { data, tier: CacheTier.L2 };
			}
			this.stats.L2Misses++;
		}

		// Fallback HTTP call when L2 misses
		// This calls the source service to fetch from MongoDB and repopulate L2
		if (this.fallbackClient) {
			try {
				const data = await this.fetchFromFallback(key);
				if (data) {
					this.stats.FallbackHits++;
					// L2 was repopulated by the internal endpoint
					// Populate L0/L1 locally
					this.setL0(prefixed, data, this.config.L0DefaultTTL);
					if (isL1Enabled(this.config, this.l1Dir)) {
						try {
							writeL1(this.l1Dir, prefixed, data, this.config.L1DefaultTTL, this.stats);
						} catch {
							// Log but don't fail
						}
					}
					return { data, tier: CacheTier.Fallback };
				}
			} catch {
				this.stats.FallbackMisses++;
			}
		}

		throw ErrCacheMiss;
	}

	/**
	 * Set stores a value in L0 (RAM) and/or L1 (Disk) based on configuration.
	 *
	 * Does NOT write to L2 - that's the source of truth managed separately.
	 */
	Set(key: string, value: Buffer, ttlMs?: number): void {
		if (!key) {
			throw ErrEmptyKey;
		}
		if (!value) {
			throw ErrNilData;
		}

		const prefixed = prefixKey(this.keyPrefix, key);
		const l0Ttl = ttlMs ?? this.config.L0DefaultTTL;
		const l1Ttl = ttlMs ?? this.config.L1DefaultTTL;

		// Set in L0 (RAM) if enabled
		this.setL0(prefixed, value, l0Ttl);

		// Set in L1 (Disk) if enabled
		if (isL1Enabled(this.config, this.l1Dir)) {
			try {
				writeL1(this.l1Dir, prefixed, value, l1Ttl, this.stats);
			} catch {
				// Log but don't fail - L0 may still be cached
			}
		}
	}

	/**
	 * Delete removes a value from all cache tiers (L0, L1).
	 *
	 * Does NOT delete from L2 - use the MinIO client directly for that.
	 */
	Delete(key: string): void {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		// Delete from L0 if enabled
		if (isL0Enabled(this.config, this.l0Cache)) {
			this.l0Cache!.delete(prefixed);
		}

		// Delete from L1 if enabled
		if (isL1Enabled(this.config, this.l1Dir)) {
			deleteL1Files(this.l1Dir, prefixed, this.stats);
		}
	}

	/**
	 * Invalidate removes a value from local cache only (L0 + L1).
	 * Alias for Delete - used for semantic clarity in invalidation scenarios.
	 */
	Invalidate(key: string): void {
		return this.Delete(key);
	}

	/**
	 * Stats returns current cache statistics.
	 */
	Stats(): CacheStats {
		return { ...this.stats };
	}

	/**
	 * GetFromL0 retrieves directly from RAM cache.
	 * Returns the data and whether it was found.
	 */
	GetFromL0(key: string): { data: Buffer | null; found: boolean } {
		if (!isL0Enabled(this.config, this.l0Cache)) {
			return { data: null, found: false };
		}

		const entry = this.l0Cache!.get(key);
		if (!entry) {
			return { data: null, found: false };
		}

		// Check TTL
		if (Date.now() > entry.ExpiresAt) {
			this.l0Cache!.delete(key);
			return { data: null, found: false };
		}

		return { data: entry.Data, found: true };
	}

	/**
	 * GetFromL1 retrieves directly from disk cache.
	 */
	GetFromL1(key: string): Buffer {
		if (!isL1Enabled(this.config, this.l1Dir)) {
			throw ErrCacheMiss;
		}
		return readL1(this.l1Dir, key, this.stats);
	}

	/**
	 * Warmup preloads keys into L0/L1 from L2.
	 * Useful for pre-populating cache on service startup.
	 */
	async Warmup(keys: string[]): Promise<void> {
		if (!this.l2Loader) {
			throw ErrNoLoader;
		}

		for (const key of keys) {
			try {
				await this.Get(key);
			} catch {
				// Ignore individual failures during warmup
			}
		}
	}

	/**
	 * setL0 stores a value in L0 (RAM) cache with TTL.
	 * Does nothing if L0 is not enabled.
	 */
	private setL0(key: string, value: Buffer, ttlMs: number): void {
		if (!isL0Enabled(this.config, this.l0Cache)) {
			return;
		}

		// Simple LRU eviction: remove oldest entries if at capacity
		while (this.l0Cache!.size >= this.config.L0MaxItems) {
			const firstKey = this.l0Cache!.keys().next().value;
			if (firstKey) {
				this.l0Cache!.delete(firstKey);
			} else {
				break;
			}
		}

		const entry: L0Entry = {
			Data: value,
			ExpiresAt: Date.now() + ttlMs,
			Size: value.length,
		};
		this.l0Cache!.set(key, entry);
	}

	/**
	 * GetOrLoad retrieves from cache or loads from L2 if not cached.
	 */
	async GetOrLoad(key: string, loader?: LocalCacheLoader): Promise<Buffer> {
		try {
			const result = await this.Get(key);
			return result.data;
		} catch {
			// Cache miss - use provided loader or default
			const l = loader ?? this.l2Loader;
			if (!l) {
				throw ErrNoLoader;
			}

			// Load from L2
			const data = await l(key);
			if (!data) {
				throw ErrCacheMiss;
			}

			// Cache the loaded data
			this.Set(key, data);
			this.stats.L2Hits++;

			return data;
		}
	}

	/**
	 * fetchFromFallback calls the HTTP fallback endpoint to fetch data from source.
	 *
	 * Uses FallbackKeyTransformer (if configured) to convert the cache key into a URL path.
	 * Falls back to last-segment extraction for backward compatibility.
	 * The internal endpoint fetches from MongoDB and repopulates L2 (MinIO) before returning.
	 *
	 * @param key - The cache key (e.g., "orgId123/defId456/scripts/nodeId")
	 * @returns The fetched data as Buffer
	 */
	private async fetchFromFallback(key: string): Promise<Buffer | null> {
		if (!this.fallbackClient) {
			return null;
		}

		// Use transformer if provided, otherwise extract last segment (backward compatible)
		const urlPath = this.fallbackKeyTransformer
			? this.fallbackKeyTransformer(key)
			: this.defaultKeyToUrlPath(key);
		const endpoint = `${this.fallbackEndpoint}${urlPath}`;

		try {
			const response = await this.fallbackClient.get<{ success?: boolean; data?: unknown }>(endpoint);

			// HTTPClient.get<T> returns the parsed body directly: { success: true, data: <payload> }
			// Extract the .data field from the API response wrapper
			const body = response as { success?: boolean; data?: unknown };
			if (!body || body.data === undefined || body.data === null) {
				getInfraLogger().error({ key, endpoint, body }, '[INFRA:CACHE] Fallback response missing data field');
				return null;
			}

			const payload = body.data;

			// Convert to Buffer — if string, store as-is; if object, stringify
			const serialized = typeof payload === 'string'
				? payload
				: JSON.stringify(payload);

			getInfraLogger().debug({ key, payloadLength: serialized.length }, '[INFRA:CACHE] Fallback hit');
			return Buffer.from(serialized);
		} catch (error) {
			getInfraLogger().debug({ key }, '[INFRA:CACHE] Fallback miss');
			throw error;
		}
	}

	/**
	 * defaultKeyToUrlPath extracts the last segment of a key to build a URL path.
	 * This is the original behavior — used when no FallbackKeyTransformer is configured.
	 *
	 * @param key - The cache key (e.g., "orgId123/assetUUID456")
	 * @returns URL path segment (e.g., "/assetUUID456")
	 */
	private defaultKeyToUrlPath(key: string): string {
		const lastSlashIndex = key.lastIndexOf('/');
		const resourceId = lastSlashIndex !== -1 ? key.substring(lastSlashIndex + 1) : key;
		return `/${resourceId}`;
	}

	/**
	 * Close closes the cache and releases all resources.
	 */
	Close(): void {
		if (this.l0Cache) {
			this.l0Cache.clear();
		}
		getInfraLogger().info('[INFRA:CACHE] Closed');
	}
}
