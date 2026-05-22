import type { TemplateReadModel } from '../types';

/**
 * TemplateCachePort abstracts template data retrieval from TieredCache.
 *
 * Implemented by TemplateCacheAdapter (infrastructure layer).
 * Used by consumers for preprocessing (fetch template scripts before execution)
 * and by adapters for cache invalidation (FANOUT events).
 */
export interface TemplateCachePort {
	/**
	 * Fetches template read model by cache key.
	 * Key format: {orgId}/{templateId}
	 *
	 * @param key - Cache key in format orgId/templateId
	 * @returns Template read model or null if not found
	 */
	get(key: string): Promise<TemplateReadModel | null>;

	/**
	 * Invalidates a template cache entry (L0 + L1).
	 * L2 (MinIO) remains as source of truth.
	 *
	 * @param key - Cache key in format orgId/templateId
	 */
	invalidate(key: string): void;
}
