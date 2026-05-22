import type { TieredCacheClient } from '@mapexos/infrastructure';
import type { AssetCachePort } from '../../application/ports/asset_cache_port';
import type { TemplateCachePort } from '../../application/ports/template_cache_port';
import type { AssetReadModel, TemplateReadModel } from '../../application/types';

/**
 * AssetCacheAdapter wraps TieredCacheClient for asset data access.
 *
 * Used by consumers for preprocessing (fetch asset before execution)
 * and by AssetInvalidateConsumer for cache invalidation.
 */
export class AssetCacheAdapter implements AssetCachePort {
	constructor(private readonly cache: TieredCacheClient) {}

	/**
	 * Fetches asset read model from TieredCache.
	 *
	 * @param key - Cache key in format orgId/assetUUID
	 * @returns Asset read model or null
	 */
	async get(key: string): Promise<AssetReadModel | null> {
		const result = await this.cache.Get(key);
		if (!result || !result.data) return null;
		return JSON.parse(result.data.toString('utf-8'));
	}

	/**
	 * Invalidates asset cache entry (L0 + L1).
	 *
	 * @param key - Cache key in format orgId/assetUUID
	 */
	invalidate(key: string): void {
		this.cache.Invalidate(key);
	}
}

/**
 * TemplateCacheAdapter wraps TieredCacheClient for template data access.
 *
 * Used by consumers for preprocessing (fetch template scripts before execution)
 * and by TemplateInvalidateConsumer for cache invalidation.
 */
export class TemplateCacheAdapter implements TemplateCachePort {
	constructor(private readonly cache: TieredCacheClient) {}

	/**
	 * Fetches template read model from TieredCache.
	 *
	 * @param key - Cache key in format orgId/templateId
	 * @returns Template read model or null
	 */
	async get(key: string): Promise<TemplateReadModel | null> {
		const result = await this.cache.Get(key);
		if (!result || !result.data) return null;
		return JSON.parse(result.data.toString('utf-8'));
	}

	/**
	 * Invalidates template cache entry (L0 + L1).
	 *
	 * @param key - Cache key in format orgId/templateId
	 */
	invalidate(key: string): void {
		this.cache.Invalidate(key);
	}
}
