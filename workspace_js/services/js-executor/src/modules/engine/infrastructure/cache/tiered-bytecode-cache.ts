import type { BytecodeCachePort, BytecodeCacheContext } from '@modules/engine/application/ports';
import type { Logger } from '@mapexos/microservices';
import type { TieredCacheClient, MinIOClient } from '@mapexos/infrastructure';

import { compressBytecode, decompressBytecode } from '@shared/utils';

/**
 * TieredCache implementation of BytecodeCachePort.
 *
 * This infrastructure service handles storage of compiled V8 bytecode using TieredCache.
 * It compresses bytecode before storage and decompresses on retrieval to reduce
 * memory usage.
 *
 * Cache hierarchy:
 *   - L0 (RAM): Ultra-fast in-memory cache (~50µs) - primary lookup
 *   - L1 (Disk): NVMe/SSD storage (~500µs) - persistence across restarts
 *   - L2 (MinIO): Source of truth for horizontal scaling
 *
 * Key format: {orgId}/{templateId}/{scriptType}.bin
 * Example: mapexos_public/507f1f77bcf86cd799439011/DECODE.bin
 *
 * Benefits of template-based keys:
 *   - Bytecode belongs to TEMPLATES, not assets
 *   - Multiple assets using same template share the bytecode
 *   - Pod A compiles and stores in MinIO
 *   - Pod B, C, D can reuse without recompiling (horizontal scaling)
 *   - Significant CPU savings at scale
 *
 * @remarks
 * This service is part of the Infrastructure layer and implements a Port interface
 * defined in the Application layer, following Hexagonal Architecture.
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
	 * Lookup order: L0 (RAM) → L1 (Disk)
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{templateId}/{scriptType})
	 * @returns Promise resolving to the decompressed bytecode Buffer, or null if not cached
	 */
	async get(cacheKey: string): Promise<Buffer | null> {
		try {
			const result = await this.bytecodeCache.Get(cacheKey);

			if (!result || !result.data || result.data.length === 0) {
				return null;
			}

			// Decompress the cached bytecode
			const { bytecode } = decompressBytecode(result.data);
			this.logger.debug(`[CACHE:Bytecode] Cache hit for ${cacheKey} (tier: L${result.tier})`);
			return Buffer.from(bytecode);
		} catch (error) {
			// ErrCacheMiss is expected for first-time access
			const errorMessage = error instanceof Error ? error.message : String(error);
			if (errorMessage.includes('cache miss')) {
				return null;
			}
			this.logger.warn(`[CACHE:Bytecode] Failed to get cache for ${cacheKey}: ${errorMessage}`);
			return null;
		}
	}

	/**
	 * Stores compressed bytecode in cache.
	 *
	 * Stores in all tiers for fast access, persistence, and horizontal scaling:
	 *   - L0 (RAM): Fast access for current pod
	 *   - L1 (Disk): Persistence across restarts
	 *   - L2 (MinIO): Shared across all pods for horizontal scaling
	 *
	 * @param cacheKey - The cache key (format: {orgId}/{templateId}/{scriptType})
	 * @param bytecode - The raw bytecode ArrayBuffer to compress and store
	 * @param ttlMs - Time-to-live in milliseconds (default: 3600000 = 1 hour)
	 */
	async set(cacheKey: string, bytecode: ArrayBuffer, ttlMs = 3600000): Promise<void> {
		try {
			const { payload, compressed, originalSize, finalSize } = compressBytecode(bytecode);

			// Store in L0 (RAM) + L1 (Disk)
			this.bytecodeCache.Set(cacheKey, payload, ttlMs);

			// Store in L2 (MinIO) for horizontal scaling
			// Other pods will benefit from this compiled bytecode
			try {
				await this.minioClient.Put(cacheKey + '.bin', payload);
				this.logger.info(
					`[CACHE:Bytecode] Cached ${cacheKey} in L0/L1/L2 (${compressed ? 'compressed' : 'uncompressed'}: ${originalSize}b → ${finalSize}b)`
				);
			} catch (minioError) {
				// L2 failure is not critical - L0/L1 still have the data
				const minioErrMsg = minioError instanceof Error ? minioError.message : String(minioError);
				this.logger.warn(`[CACHE:Bytecode] Failed to store in L2 (MinIO) for ${cacheKey}: ${minioErrMsg}`);
				this.logger.info(
					`[CACHE:Bytecode] Cached ${cacheKey} in L0/L1 only (${compressed ? 'compressed' : 'uncompressed'}: ${originalSize}b → ${finalSize}b)`
				);
			}
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.warn(`[CACHE:Bytecode] Failed to set cache for ${cacheKey}: ${errorMessage}`);
		}
	}

	/**
	 * Invalidates cached bytecode for a specific template.
	 *
	 * Removes from all tiers: L0 (RAM), L1 (Disk), and L2 (MinIO).
	 *
	 * @param ctx - The cache context containing templateId and templateOrgId
	 * @param scriptName - Optional specific script to invalidate. If not provided, all scripts for the template are invalidated.
	 */
	async invalidate(ctx: BytecodeCacheContext, scriptName?: string): Promise<void> {
		try {
			if (scriptName) {
				// Invalidate specific script
				const cacheKey = this.buildCacheKey(ctx, scriptName);
				this.bytecodeCache.Invalidate(cacheKey);

				// Also remove from L2 (MinIO)
				try {
					await this.minioClient.Delete(cacheKey + '.bin');
				} catch {
					// Ignore - file may not exist
				}

				this.logger.info(`[CACHE:Bytecode] Invalidated ${cacheKey} from L0/L1/L2`);
			} else {
				// Invalidate all scripts for the template
				const keys = [
					this.buildCacheKey(ctx, 'payloadDecode'),
					this.buildCacheKey(ctx, 'payloadValidation'),
					this.buildCacheKey(ctx, 'payloadTransform'),
				];
				for (const key of keys) {
					this.bytecodeCache.Invalidate(key);

					// Also remove from L2 (MinIO)
					try {
						await this.minioClient.Delete(key + '.bin');
					} catch {
						// Ignore - file may not exist
					}
				}
				this.logger.info(`[CACHE:Bytecode] Invalidated all scripts for template ${ctx.templateId} from L0/L1/L2`);
			}
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.warn(`[CACHE:Bytecode] Failed to invalidate cache for template ${ctx.templateId}: ${errorMessage}`);
		}
	}

	/**
	 * Builds a cache key for a script.
	 *
	 * Key format: {orgId}/{templateId}/{scriptType}
	 * Example: mapexos_public/507f1f77bcf86cd799439011/DECODE
	 *
	 * @param ctx - The cache context containing templateId and templateOrgId
	 * @param scriptName - The script name (payloadDecode, payloadValidation, payloadTransform)
	 * @returns The cache key string
	 */
	buildCacheKey(ctx: BytecodeCacheContext, scriptName: string): string {
		const scriptType = this.normalizeScriptType(scriptName);
		return `${ctx.templateOrgId}/${ctx.templateId}/${scriptType}`;
	}

	/**
	 * Normalizes script name to uppercase type.
	 *
	 * @param scriptName - The script name (payloadDecode, payloadValidation, payloadTransform)
	 * @returns The normalized script type (DECODE, VALIDATION, TRANSFORM)
	 */
	private normalizeScriptType(scriptName: string): string {
		switch (scriptName) {
			case 'payloadDecode':
				return 'DECODE';
			case 'payloadValidation':
				return 'VALIDATION';
			case 'payloadTransform':
				return 'TRANSFORM';
			default:
				return scriptName.toUpperCase();
		}
	}
}
