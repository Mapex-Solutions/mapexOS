/**
 * Cache context for workflow bytecode operations.
 *
 * Key format: {orgId}/{definitionId}/bytecode/{nodeId}
 */
export interface BytecodeCacheContext {
	/** Organization ID */
	orgId: string;
	/** Workflow definition ID */
	workflowId: string;
	/** Node ID within the workflow */
	nodeId: string;
}

/**
 * Port interface for Bytecode Cache (infrastructure contract).
 *
 * Stores compressed V8 bytecode to avoid recompiling scripts on each execution.
 * Bytecode belongs to DEFINITIONS (not instances) — multiple instances of the same
 * workflow reuse the same compiled bytecode.
 *
 * Key format: {orgId}/{definitionId}/bytecode/{nodeId}.bin
 *
 * L2 (MinIO) lifecycle:
 *   - WRITE: js-workflow-executor writes bytecode after first compile
 *   - DELETE: Go workflow service deletes when script source changes
 *   - js-workflow-executor only invalidates L0 (RAM) + L1 (Disk) on FANOUT
 */
export interface BytecodeCachePort {
	/**
	 * Retrieves cached bytecode for a script.
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{definitionId}/bytecode/{nodeId})
	 * @returns Promise resolving to the decompressed bytecode Buffer, or null if not cached
	 */
	get(cacheKey: string): Promise<Buffer | null>;

	/**
	 * Stores compressed bytecode in cache (L0/L1 + L2 MinIO).
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{definitionId}/bytecode/{nodeId})
	 * @param bytecode - The raw bytecode ArrayBuffer to compress and store
	 * @param ttlMs - Time-to-live in milliseconds (default: 3600000 = 1 hour)
	 */
	set(cacheKey: string, bytecode: ArrayBuffer, ttlMs?: number): Promise<void>;

	/**
	 * Invalidates cached bytecode for a workflow node (L0 + L1 only).
	 * L2 (MinIO) cleanup is handled by the Go workflow service.
	 *
	 * @param ctx - The cache context with orgId, workflowId, nodeId
	 */
	invalidate(ctx: BytecodeCacheContext): Promise<void>;

	/**
	 * Invalidates cached bytecode for specific nodes (L0 + L1 only).
	 * Called on FANOUT with granular nodeIds from Go workflow service.
	 *
	 * @param orgId - Organization ID
	 * @param definitionId - Workflow definition ID
	 * @param nodeIds - Node IDs to invalidate
	 */
	invalidateNodes(orgId: string, definitionId: string, nodeIds: string[]): Promise<void>;

	/**
	 * Invalidates ALL cached bytecode for a workflow definition (L0 + L1 only).
	 * Fallback when nodeIds are not available — relies on TTL expiry.
	 *
	 * @param orgId - Organization ID
	 * @param workflowId - Workflow definition ID
	 */
	invalidateWorkflow(orgId: string, workflowId: string): Promise<void>;

	/**
	 * Builds a cache key for a workflow node's bytecode.
	 *
	 * @param ctx - The cache context with orgId, workflowId, nodeId
	 * @returns The cache key string
	 */
	buildCacheKey(ctx: BytecodeCacheContext): string;
}
