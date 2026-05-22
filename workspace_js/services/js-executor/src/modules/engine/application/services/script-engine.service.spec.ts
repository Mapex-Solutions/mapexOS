/**
 * ScriptEngineService Unit Tests
 *
 * Tests the Piscina-based ScriptEngineService (main thread orchestrator).
 * Verifies: pipeline dispatch, OOM propagation, metrics, initialization, shutdown.
 *
 * Mocks Piscina to test main-thread logic without spawning real worker threads.
 */

import type { Logger } from '@mapexos/microservices';
import type { BytecodeCachePort } from '@modules/engine/application/ports';
import type { PiscinaWorkerConfig } from '@modules/engine/infrastructure/worker';
import type { PiscinaPoolMetrics, PiscinaOptions, ScriptEngineMetrics } from '@modules/engine/application/types';

import { ScriptEngineService } from './script-engine.service';
import { OOMError } from '@modules/engine/domain/errors';

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

const createMockBytecodeCache = (): BytecodeCachePort => ({
	get: jest.fn().mockResolvedValue(null),
	set: jest.fn().mockResolvedValue(undefined),
	invalidate: jest.fn().mockResolvedValue(undefined),
	buildCacheKey: jest.fn((ctx, name) => `${ctx.templateOrgId}/${ctx.templateId}/${name}`),
});

const createMockPiscinaOptions = (): PiscinaOptions => ({
	workers: 2,
	workerPath: '/mock/worker.js',
});

const createMockWorkerConfig = (): PiscinaWorkerConfig => ({
	memoryLimitMb: 32,
	timeoutMs: 5000,
	contextRecycleInterval: 10000,
	mapexValidatorCode: '',
});

const createMockMetrics = (): ScriptEngineMetrics => ({
	scriptDuration: { observe: jest.fn() } as any,
	scriptErrors: { inc: jest.fn() } as any,
	compileDuration: { observe: jest.fn() } as any,
	bytecodeCache: { inc: jest.fn() } as any,
	scriptRegistry: { inc: jest.fn() } as any,
});

const createMockPoolMetrics = (): PiscinaPoolMetrics => ({
	piscinaCompleted: { inc: jest.fn() } as any,
	piscinaRunDuration: { observe: jest.fn() } as any,
	piscinaWaitDuration: { observe: jest.fn() } as any,
	piscinaWorkers: { set: jest.fn() } as any,
});

// ─── Mock Piscina ────────────────────────────────────────────────────

const mockPiscinaRun = jest.fn();
const mockPiscinaDestroy = jest.fn().mockResolvedValue(undefined);

jest.mock('piscina', () => {
	return jest.fn().mockImplementation(() => ({
		run: mockPiscinaRun,
		destroy: mockPiscinaDestroy,
		completed: 0,
		utilization: 0.5,
		runTime: { average: 1000 },
		waitTime: { average: 100 },
	}));
});

describe('ScriptEngineService', () => {
	let service: ScriptEngineService;
	let logger: Logger;
	let bytecodeCache: BytecodeCachePort;
	let metrics: ScriptEngineMetrics;
	let poolMetrics: PiscinaPoolMetrics;

	beforeEach(() => {
		jest.clearAllMocks();
		logger = createMockLogger();
		bytecodeCache = createMockBytecodeCache();
		metrics = createMockMetrics();
		poolMetrics = createMockPoolMetrics();

		service = new ScriptEngineService(
			logger,
			bytecodeCache,
			createMockPiscinaOptions(),
			createMockWorkerConfig(),
			metrics,
			poolMetrics,
		);
	});

	afterEach(async () => {
		try { await service.shutdown(); } catch { /* ignore */ }
	});

	describe('initialize', () => {
		it('should create Piscina pool on initialize', async () => {
			await service.initialize();

			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalledWith(expect.objectContaining({
				filename: '/mock/worker.js',
				minThreads: 2,
				maxThreads: 2,
			}));
		});

		it('should only initialize once', async () => {
			await service.initialize();
			await service.initialize();

			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalledTimes(1);
		});
	});

	describe('shutdown', () => {
		it('should destroy Piscina pool', async () => {
			await service.initialize();
			await service.shutdown();

			expect(mockPiscinaDestroy).toHaveBeenCalledTimes(1);
		});

		it('should be safe to call shutdown without initialize', async () => {
			await service.shutdown(); // Should not throw
		});
	});

	describe('runScriptPipeline', () => {
		const scripts = { decode: '', validation: '', transform: 'var result = payload;' };
		const payload = { temperature: 25.3 };
		const cacheContext = { templateId: 'tpl-001', templateOrgId: 'org-001' };

		it('should dispatch to Piscina and return success result', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: true,
				finalPayload: { processed: true, temperature: 25.3 },
				totalPipelineTime: 1.5,
			});

			const result = await service.runScriptPipeline(payload, scripts, cacheContext);

			expect(result.success).toBe(true);
			expect(result.finalPayload).toEqual({ processed: true, temperature: 25.3 });
			expect(result.totalPipelineTime).toBe(1.5);
		});

		it('should build correct PiscinaWorkerInput', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true, finalPayload: {} });

			await service.runScriptPipeline(payload, scripts, cacheContext);

			expect(mockPiscinaRun).toHaveBeenCalledWith({
				rawPayload: payload,
				scripts: {
					decode: '',
					validation: '',
					transform: 'var result = payload;',
				},
				templateId: 'tpl-001',
			});
		});

		it('should use "default" templateId when no cacheContext', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true, finalPayload: {} });

			await service.runScriptPipeline(payload, scripts);

			expect(mockPiscinaRun).toHaveBeenCalledWith(
				expect.objectContaining({ templateId: 'default' })
			);
		});

		it('should return failure result for script errors', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				failedAt: 'validation',
				error: 'ReferenceError: x is not defined',
				totalPipelineTime: 0.5,
			});

			const result = await service.runScriptPipeline(payload, scripts, cacheContext);

			expect(result.success).toBe(false);
			expect(result.failedAt).toBe('validation');
			expect(result.error).toContain('ReferenceError');
		});

		it('should throw OOMError when worker reports OOM', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'isolate disposed',
				isOOM: true,
			});

			await expect(service.runScriptPipeline(payload, scripts, cacheContext))
				.rejects
				.toThrow(OOMError);
		});

		it('should handle Piscina.run() rejection (worker crash)', async () => {
			mockPiscinaRun.mockRejectedValue(new Error('Worker thread terminated'));

			const result = await service.runScriptPipeline(payload, scripts, cacheContext);

			expect(result.success).toBe(false);
			expect(result.error).toContain('Worker thread terminated');
		});

		it('should auto-initialize on first runScriptPipeline call', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true, finalPayload: {} });

			// Don't call initialize() explicitly
			const result = await service.runScriptPipeline(payload, scripts);

			expect(result.success).toBe(true);
			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalled();
		});
	});

	describe('Metrics', () => {
		const scripts = { decode: '', validation: '', transform: 'var result = payload;' };
		const payload = { temperature: 25.3 };

		it('should increment piscinaCompleted on success', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true, finalPayload: {} });

			await service.runScriptPipeline(payload, scripts);

			expect(poolMetrics.piscinaCompleted.inc).toHaveBeenCalled();
		});

		it('should observe piscinaRunDuration on success', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true, finalPayload: {} });

			await service.runScriptPipeline(payload, scripts);

			expect(poolMetrics.piscinaRunDuration.observe).toHaveBeenCalled();
		});

		it('should increment scriptErrors on failure', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				failedAt: 'transform',
				error: 'TypeError: cannot read property',
			});

			await service.runScriptPipeline(payload, scripts);

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ script_type: 'transform', error_type: 'type' })
			);
		});

		it('should classify syntax errors correctly', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				failedAt: 'decode',
				error: 'SyntaxError: Unexpected token',
			});

			await service.runScriptPipeline(payload, scripts);

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ error_type: 'syntax' })
			);
		});

		it('should classify reference errors correctly', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				failedAt: 'transform',
				error: 'x is not defined',
			});

			await service.runScriptPipeline(payload, scripts);

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ error_type: 'reference' })
			);
		});
	});

	describe('getPoolStats', () => {
		it('should return empty stats before initialization', () => {
			const stats = service.getPoolStats();

			expect(stats.piscina.completed).toBe(0);
			expect(stats.piscina.threads).toBe(0);
		});

		it('should return Piscina stats after initialization', async () => {
			await service.initialize();

			const stats = service.getPoolStats();

			expect(stats.piscina).toBeDefined();
			expect(stats.piscina.threads).toBe(2);
			expect(typeof stats.piscina.utilization).toBe('number');
		});
	});

	describe('Event Loss Prevention (OOM propagation)', () => {
		const scripts = { decode: '', validation: '', transform: 'var result = payload;' };
		const payload = { data: 'test' };

		it('should propagate OOMError so caller can NACK the message', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'V8 heap out of memory',
				isOOM: true,
			});

			try {
				await service.runScriptPipeline(payload, scripts);
				fail('Should have thrown OOMError');
			} catch (error) {
				expect(error).toBeInstanceOf(OOMError);
				expect((error as OOMError).message).toContain('V8 heap out of memory');
			}
		});

		it('should NOT throw for non-OOM failures (event is consumed, not retried)', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				failedAt: 'transform',
				error: 'Script error: x is not defined',
			});

			// Should NOT throw — returns failure result, event is ack'd (not retried)
			const result = await service.runScriptPipeline(payload, scripts);
			expect(result.success).toBe(false);
			expect(result.error).toBeDefined();
		});
	});
});

// ─── Batch Pipeline Tests ─────────────────────────────────────────────

describe('ScriptEngineService - Batch Pipeline', () => {
	let service: ScriptEngineService;
	let logger: Logger;
	let bytecodeCache: BytecodeCachePort;
	let poolMetrics: PiscinaPoolMetrics;

	const _createOldBatchPoolOptions = (): any => ({
		workers: 3,
		workerPath: '/mock/batch-worker.js',
		eventsPerWorker: 2,
		nats: { url: 'nats://localhost:4222', user: 'test', pass: 'test' },
	});

	const _createBatchEvent = (index: number): any => ({
		rawPayload: { temperature: 20 + index },
		scripts: { decode: '', validation: '', transform: 'var result = payload;' },
		templateId: 'tpl-001',
		assetUUID: `asset-uuid-${index}`,
		assetId: `asset-id-${index}`,
		eventTrackerId: `tracker-${index}`,
		debugEnabled: false,
		sourceType: 'http',
		dataSource: {
			id: `ds-${index}`,
			orgId: 'org-001',
			pathKey: '/root',
			name: 'TestDevice',
			description: 'Test',
		},
	});

	beforeEach(() => {
		jest.clearAllMocks();
		logger = createMockLogger();
		bytecodeCache = createMockBytecodeCache();
		poolMetrics = createMockPoolMetrics();

		service = new ScriptEngineService(
			logger,
			bytecodeCache,
			createMockPiscinaOptions(),
			createMockWorkerConfig(),
			createMockMetrics(),
			poolMetrics,
		);
	});

	afterEach(async () => {
		try { await service.shutdown(); } catch { /* ignore */ }
	});

	describe('initialize with batch pool', () => {
		it('should create both single-event and batch Piscina pools', async () => {
			await service.initialize();

			const Piscina = require('piscina');
			// Called twice: once for single-event pool, once for batch pool
			expect(Piscina).toHaveBeenCalledTimes(2);
		});

		it('should create batch pool with correct options', async () => {
			await service.initialize();

			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalledWith(expect.objectContaining({
				filename: '/mock/batch-worker.js',
				minThreads: 3,
				maxThreads: 3,
			}));
		});
	});

	describe('shutdown with batch pool', () => {
		it('should destroy both pools on shutdown', async () => {
			await service.initialize();
			await service.shutdown();

			// destroy is called on both pool instances
			expect(mockPiscinaDestroy).toHaveBeenCalledTimes(2);
		});
	});

	// TODO: Add runBatch tests (TKT-2026-0024 — replaces runBatchPipeline)
});
