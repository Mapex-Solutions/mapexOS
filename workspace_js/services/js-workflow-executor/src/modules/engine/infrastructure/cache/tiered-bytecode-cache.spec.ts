/**
 * TieredBytecodeCache Unit Tests
 *
 * Tests bytecode cache: get/set with compression, invalidation, MinIO integration.
 * Mocks: TieredCacheClient, MinIOClient.
 */

import type { Logger } from '@mapexos/microservices';
import type { TieredCacheClient, MinIOClient } from '@mapexos/infrastructure';

import { TieredBytecodeCache } from './tiered-bytecode-cache';
import { compressBytecode } from '@shared/utils';

// ─── Mock Helpers ────────────────────────────────────────────────────

const createMockLogger = (): Logger => ({
	info: jest.fn(),
	debug: jest.fn(),
	warn: jest.fn(),
	error: jest.fn(),
	trace: jest.fn(),
	fatal: jest.fn(),
	child: jest.fn().mockReturnThis(),
} as unknown as Logger);

const createMockTieredCache = (): TieredCacheClient => ({
	Get: jest.fn(),
	Set: jest.fn(),
	Invalidate: jest.fn(),
} as unknown as TieredCacheClient);

const createMockMinIO = (): MinIOClient => ({
	Get: jest.fn(),
	Put: jest.fn(),
	Delete: jest.fn(),
} as unknown as MinIOClient);

describe('TieredBytecodeCache', () => {
	let cache: TieredBytecodeCache;
	let logger: Logger;
	let tieredCache: TieredCacheClient;
	let minioClient: MinIOClient;

	beforeEach(() => {
		jest.clearAllMocks();
		logger = createMockLogger();
		tieredCache = createMockTieredCache();
		minioClient = createMockMinIO();

		cache = new TieredBytecodeCache(tieredCache, minioClient, logger);
	});

	describe('get', () => {
		it('should return decompressed bytecode on cache hit', async () => {
			const originalBytecode = new ArrayBuffer(1024);
			new Uint8Array(originalBytecode).fill(0x42);
			const { payload } = compressBytecode(originalBytecode);

			(tieredCache.Get as jest.Mock).mockResolvedValue({ data: payload, tier: 0 });

			const result = await cache.get('org-001/wf-001/bytecode/node-001');

			expect(result).toBeInstanceOf(Buffer);
			expect(result!.length).toBe(1024);
		});

		it('should return null on cache miss (null result)', async () => {
			(tieredCache.Get as jest.Mock).mockResolvedValue(null);

			const result = await cache.get('org-001/wf-001/bytecode/node-001');

			expect(result).toBeNull();
		});

		it('should return null on cache miss (empty data)', async () => {
			(tieredCache.Get as jest.Mock).mockResolvedValue({ data: Buffer.alloc(0), tier: 0 });

			const result = await cache.get('org-001/wf-001/bytecode/node-001');

			expect(result).toBeNull();
		});

		it('should return null on cache error', async () => {
			(tieredCache.Get as jest.Mock).mockRejectedValue(new Error('Connection lost'));

			const result = await cache.get('org-001/wf-001/bytecode/node-001');

			expect(result).toBeNull();
			expect(logger.warn).toHaveBeenCalled();
		});

		it('should return null silently on cache miss error message', async () => {
			(tieredCache.Get as jest.Mock).mockRejectedValue(new Error('cache miss'));

			const result = await cache.get('org-001/wf-001/bytecode/node-001');

			expect(result).toBeNull();
			expect(logger.warn).not.toHaveBeenCalled();
		});
	});

	describe('set', () => {
		it('should compress and store in L0/L1 + L2 (MinIO)', async () => {
			const bytecode = new ArrayBuffer(1024);

			await cache.set('org-001/wf-001/bytecode/node-001', bytecode);

			expect(tieredCache.Set).toHaveBeenCalledWith(
				'org-001/wf-001/bytecode/node-001',
				expect.any(Buffer),
				3600000,
			);
			expect(minioClient.Put).toHaveBeenCalledWith(
				'org-001/wf-001/bytecode/node-001.bin',
				expect.any(Buffer),
			);
		});

		it('should use custom TTL when provided', async () => {
			const bytecode = new ArrayBuffer(100);

			await cache.set('key', bytecode, 7200000);

			expect(tieredCache.Set).toHaveBeenCalledWith('key', expect.any(Buffer), 7200000);
		});

		it('should handle MinIO upload failure gracefully', async () => {
			(minioClient.Put as jest.Mock).mockRejectedValue(new Error('MinIO down'));

			const bytecode = new ArrayBuffer(100);
			await cache.set('key', bytecode); // Should not throw

			expect(tieredCache.Set).toHaveBeenCalled(); // L0/L1 still stored
			expect(logger.warn).toHaveBeenCalledWith(expect.stringContaining('MinIO'));
		});

		it('should handle compression error gracefully', async () => {
			// Pass invalid bytecode that causes compression to fail
			await cache.set('key', null as any); // Should not throw

			expect(logger.warn).toHaveBeenCalled();
		});
	});

	describe('invalidate', () => {
		it('should invalidate L0/L1 + L2 (MinIO)', async () => {
			const ctx = { orgId: 'org-001', workflowId: 'wf-001', nodeId: 'node-001' };

			await cache.invalidate(ctx);

			expect(tieredCache.Invalidate).toHaveBeenCalledWith('org-001/wf-001/bytecode/node-001');
			expect(minioClient.Delete).toHaveBeenCalledWith('org-001/wf-001/bytecode/node-001.bin');
		});

		it('should handle MinIO delete failure silently', async () => {
			(minioClient.Delete as jest.Mock).mockRejectedValue(new Error('Not found'));

			const ctx = { orgId: 'org-001', workflowId: 'wf-001', nodeId: 'node-001' };
			await cache.invalidate(ctx); // Should not throw

			expect(tieredCache.Invalidate).toHaveBeenCalled();
		});
	});

	describe('invalidateWorkflow', () => {
		it('should log TTL-based invalidation', async () => {
			await cache.invalidateWorkflow('org-001', 'wf-001');

			expect(logger.info).toHaveBeenCalledWith(
				expect.stringContaining('wf-001')
			);
		});
	});

	describe('buildCacheKey', () => {
		it('should produce correct key format', () => {
			const ctx = { orgId: 'org-001', workflowId: 'wf-001', nodeId: 'node-001' };

			const key = cache.buildCacheKey(ctx);

			expect(key).toBe('org-001/wf-001/bytecode/node-001');
		});
	});
});
