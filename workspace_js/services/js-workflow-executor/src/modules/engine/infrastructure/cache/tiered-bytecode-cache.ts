import type { BytecodeCachePort, BytecodeCacheContext } from '@modules/engine/application/ports';
import type { Logger } from '@mapexos/microservices';
import type { TieredCacheClient, MinIOClient } from '@mapexos/infrastructure';

import { compressBytecode, decompressBytecode } from '@shared/utils';

/**
 * TieredCache implementation of BytecodeCachePort for workflow scripts.
 *
 * Cache hierarchy:
 *   - L0 (RAM): Ultra-fast in-memory cache (~50µs)
 *   - L1 (Disk): NVMe/SSD storage (~500µs)
 *   - L2 (MinIO): Source of truth for horizontal scaling
 *
 * Key format: {orgId}/{definitionId}/bytecode/{nodeId}.bin
 *
 * L2 lifecycle ownership:
 *   - WRITE: js-workflow-executor writes bytecode after V8 compilation
 *   - DELETE: Go workflow service deletes when script source changes
 *   - FANOUT invalidation only clears L0 + L1 (local caches)
 */
export class TieredBytecodeCache implements BytecodeCachePort {
	constructor(
		private readonly bytecodeCache: TieredCacheClient,
		private readonly minioClient: MinIOClient,
		private readonly logger: Logger,
	) {}

	/**
	 * Retrieves cached bytecode for a script.
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{workflowId}/bytecode/{nodeId})
	 * @returns Promise resolving to the decompressed bytecode Buffer, or null if not cached
	 */
	async get(cacheKey: string): Promise<Buffer | null> {
		try {
			const result = await this.bytecodeCache.Get(cacheKey);

			if (!result || !result.data || result.data.length === 0) {
				return null;
			}

			const { bytecode } = decompressBytecode(result.data);
			this.logger.debug(`[CACHE:Bytecode] Cache hit for ${cacheKey} (tier: L${result.tier})`);
			return Buffer.from(bytecode);
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			if (errorMessage.includes('cache miss')) {
				return null;
			}
			this.logger.warn(`[CACHE:Bytecode] Failed to get cache for ${cacheKey}: ${errorMessage}`);
			return null;
		}
	}

	/**
	 * Stores compressed bytecode in cache (L0/L1 + L2 MinIO).
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{workflowId}/bytecode/{nodeId})
	 * @param bytecode - The raw bytecode ArrayBuffer to compress and store
	 * @param ttlMs - Time-to-live in milliseconds (default: 3600000 = 1 hour)
	 */
	async set(cacheKey: string, bytecode: ArrayBuffer, ttlMs = 3600000): Promise<void> {
		try {
			const { payload, compressed, originalSize, finalSize } = compressBytecode(bytecode);

			// Store in L0 (RAM) + L1 (Disk)
			this.bytecodeCache.Set(cacheKey, payload, ttlMs);

			// Store in L2 (MinIO) for horizontal scaling
			try {
				await this.minioClient.Put(cacheKey + '.bin', payload);
				this.logger.info(
					`[CACHE:Bytecode] Cached ${cacheKey} in L0/L1/L2 (${compressed ? 'compressed' : 'uncompressed'}: ${originalSize}b → ${finalSize}b)`
				);
			} catch (minioError) {
				const minioErrMsg = minioError instanceof Error ? minioError.message : String(minioError);
				this.logger.warn(`[CACHE:Bytecode] Failed to store in L2 (MinIO) for ${cacheKey}: ${minioErrMsg}`);
			}
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.warn(`[CACHE:Bytecode] Failed to set cache for ${cacheKey}: ${errorMessage}`);
		}
	}

	/**
	 * Invalidates cached bytecode for a specific workflow node (L0 + L1 only).
	 * L2 (MinIO) cleanup is handled by the Go workflow service.
	 *
	 * @param ctx - The cache context with orgId, workflowId, nodeId
	 */
	async invalidate(ctx: BytecodeCacheContext): Promise<void> {
		try {
			const cacheKey = this.buildCacheKey(ctx);

			// Only invalidate L0 (RAM) + L1 (Disk) — Go handles L2 (MinIO)
			this.bytecodeCache.Invalidate(cacheKey);

			this.logger.info(`[CACHE:Bytecode] Invalidated ${cacheKey} from L0/L1`);
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.warn(`[CACHE:Bytecode] Failed to invalidate cache for node ${ctx.nodeId}: ${errorMessage}`);
		}
	}

	/**
	 * Invalidates cached bytecode for specific nodes (L0 + L1 only).
	 * Called on FANOUT with granular nodeIds from Go workflow service.
	 *
	 * @param orgId - Organization ID
	 * @param definitionId - Workflow definition ID
	 * @param nodeIds - Node IDs to invalidate
	 */
	async invalidateNodes(orgId: string, definitionId: string, nodeIds: string[]): Promise<void> {
		for (const nodeId of nodeIds) {
			await this.invalidate({ orgId, workflowId: definitionId, nodeId });
		}
		this.logger.info(
			`[CACHE:Bytecode] Invalidated ${nodeIds.length} node(s) for definition ${definitionId} from L0/L1`
		);
	}

	/**
	 * Invalidates ALL cached bytecode for a workflow definition (L0 + L1 only).
	 * Fallback when nodeIds are not available — relies on TTL expiry.
	 *
	 * @param orgId - Organization ID
	 * @param workflowId - Workflow definition ID
	 */
	async invalidateWorkflow(orgId: string, workflowId: string): Promise<void> {
		this.logger.info(`[CACHE:Bytecode] Workflow ${workflowId} invalidation requested — L0/L1 entries will expire by TTL`);
	}

	/**
	 * Builds a cache key for a workflow node's bytecode.
	 *
	 * Key format: {orgId}/{workflowId}/bytecode/{nodeId}
	 *
	 * @param ctx - The cache context with orgId, workflowId, nodeId
	 * @returns The cache key string
	 */
	buildCacheKey(ctx: BytecodeCacheContext): string {
		return `${ctx.orgId}/${ctx.workflowId}/bytecode/${ctx.nodeId}`;
	}
}
