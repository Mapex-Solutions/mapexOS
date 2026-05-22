/**
 * ScriptSourceCachePort abstracts workflow script source retrieval from TieredCache.
 *
 * Implemented by TieredScriptSourceCacheAdapter (infrastructure layer).
 * Used by WorkflowScriptService to fetch script source code for code nodes.
 */
export interface ScriptSourceCachePort {
	/**
	 * Fetches script source code by cache key.
	 * Key format: {orgId}/{workflowId}/scripts/{nodeId}
	 *
	 * @param key - Cache key
	 * @returns Script source string or null if not found
	 */
	get(key: string): Promise<string | null>;

	/**
	 * Invalidates a script source cache entry (L0 + L1).
	 * L2 (MinIO) remains as source of truth.
	 *
	 * @param key - Cache key
	 */
	invalidate(key: string): void;
}
