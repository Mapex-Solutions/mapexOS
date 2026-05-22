import type { AssetReadModel } from '../types';

/**
 * AssetCachePort abstracts asset data retrieval from TieredCache.
 *
 * Implemented by AssetCacheAdapter (infrastructure layer).
 * Used by consumers for preprocessing (fetch asset scripts before execution)
 * and by adapters for cache invalidation (FANOUT events).
 */
export interface AssetCachePort {
	/**
	 * Fetches asset read model by cache key.
	 * Key format: {orgId}/{assetUUID}
	 *
	 * @param key - Cache key in format orgId/assetUUID
	 * @returns Asset read model or null if not found
	 */
	get(key: string): Promise<AssetReadModel | null>;

	/**
	 * Invalidates an asset cache entry (L0 + L1).
	 * L2 (MinIO) remains as source of truth.
	 *
	 * @param key - Cache key in format orgId/assetUUID
	 */
	invalidate(key: string): void;
}
