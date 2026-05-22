import type { TieredCacheClient } from '@mapexos/infrastructure';
import type { ScriptSourceCachePort } from '../../application/ports/script_source_cache_port';

/**
 * TieredScriptSourceCacheAdapter wraps TieredCacheClient for workflow script source access.
 *
 * Used by WorkflowScriptService to fetch script source (JS code) for code nodes
 * and by DefinitionInvalidateConsumer for cache invalidation.
 */
export class TieredScriptSourceCacheAdapter implements ScriptSourceCachePort {
	constructor(private readonly cache: TieredCacheClient) {}

	/**
	 * Fetches script source string from TieredCache.
	 * Returns null on cache miss or on error (logged by caller).
	 *
	 * @param key - Cache key: {orgId}/{workflowId}/scripts/{nodeId}
	 * @returns Script source string or null
	 */
	async get(key: string): Promise<string | null> {
		const result = await this.cache.Get(key);
		if (!result || !result.data || result.data.length === 0) {
			return null;
		}

		return typeof result.data === 'string'
			? result.data
			: Buffer.from(result.data).toString('utf8');
	}

	/**
	 * Invalidates script source cache entry (L0 + L1).
	 *
	 * @param key - Cache key
	 */
	invalidate(key: string): void {
		this.cache.Invalidate(key);
	}
}
