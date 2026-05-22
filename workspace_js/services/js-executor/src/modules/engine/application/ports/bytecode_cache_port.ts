/**
 * Cache key context for bytecode operations.
 *
 * Contains the template information needed to build cache keys.
 * Bytecode belongs to TEMPLATES (not assets), so keys are based on templateId.
 *
 * Key format: {orgId}/{templateId}/{scriptType}.bin
 * Example: mapexos_public/507f1f77bcf86cd799439011/DECODE.bin (public template)
 * Example: 607f1f77bcf86cd799439022/507f1f77bcf86cd799439011/DECODE.bin (org template)
 */
export interface BytecodeCacheContext {
	/** Template ID (MongoDB ObjectId as hex string) */
	templateId: string;
	/**
	 * Organization ID for the template.
	 * If template.isSystem=true, this should be "mapexos_public"
	 */
	templateOrgId: string;
}

/**
 * Port interface for Bytecode Cache (infrastructure contract).
 *
 * This interface follows the Hexagonal Architecture pattern, defining the contract
 * for bytecode storage operations. The implementation can use Redis, file system,
 * or any other storage mechanism.
 *
 * @remarks
 * The cache stores compressed V8 bytecode to avoid recompiling scripts on each execution.
 * This significantly improves performance for frequently executed scripts.
 *
 * Bytecode belongs to TEMPLATES, not assets. Multiple assets using the same template
 * will share the cached bytecode, saving CPU by compiling only once.
 *
 * Key format: {orgId}/{templateId}/{scriptType}.bin
 *
 * This port only handles STORAGE - the actual compilation logic is in the Domain layer.
 */
export interface BytecodeCachePort {
	/**
	 * Retrieves cached bytecode for a script.
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{templateId}/{scriptType})
	 * @returns Promise resolving to the decompressed bytecode Buffer, or null if not cached
	 */
	get(cacheKey: string): Promise<Buffer | null>;

	/**
	 * Stores compressed bytecode in cache.
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{templateId}/{scriptType})
	 * @param bytecode - The raw bytecode ArrayBuffer to compress and store
	 * @param ttlSeconds - Time-to-live in seconds (default: 3600)
	 */
	set(cacheKey: string, bytecode: ArrayBuffer, ttlSeconds?: number): Promise<void>;

	/**
	 * Invalidates cached bytecode for a specific template.
	 *
	 * @param ctx - The cache context containing templateId and templateOrgId
	 * @param scriptName - Optional specific script to invalidate. If not provided, all scripts for the template are invalidated.
	 */
	invalidate(ctx: BytecodeCacheContext, scriptName?: string): Promise<void>;

	/**
	 * Builds a cache key for a script.
	 *
	 * Key format: {orgId}/{templateId}/{scriptType}
	 *
	 * @param ctx - The cache context containing templateId and templateOrgId
	 * @param scriptName - The script name (payloadDecode, payloadValidation, payloadTransform)
	 * @returns The cache key string
	 */
	buildCacheKey(ctx: BytecodeCacheContext, scriptName: string): string;
}
