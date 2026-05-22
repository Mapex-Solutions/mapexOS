/**
 * ScriptEngineService Unit Tests
 *
 * Tests the Piscina-based ScriptEngineService (main thread orchestrator).
 * Verifies: initialization, shutdown, dispatch, OOM propagation, metrics, pool stats.
 *
 * Mocks Piscina to test main-thread logic without spawning real worker threads.
 */

import type { Logger } from '@mapexos/microservices';
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

const createMockPiscinaOptions = (): PiscinaOptions => ({
	workers: 2,
	workerPath: '/mock/worker.js',
});

const createMockWorkerConfig = (): PiscinaWorkerConfig => ({
	memoryLimitMb: 32,
	timeoutMs: 5000,
	contextRecycleInterval: 10000,
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

const createWorkerInput = () => ({
	script: 'const result = { output: event.data, statePatch: {} };',
	cacheKey: 'org-001/wf-001/scripts/node-001',
	event: { data: { temperature: 25.3 } },
	state: { processedCount: 0 },
	inputs: { threshold: 80 },
	nodes: {},
});

// ─── Mock Piscina ────────────────────────────────────────────────────

const mockPiscinaRun = jest.fn();
const mockPiscinaDestroy = jest.fn().mockResolvedValue(undefined);

jest.mock('piscina', () => {
	return jest.fn().mockImplementation(() => ({
		run: mockPiscinaRun,
		destroy: mockPiscinaDestroy,
		completed: 42,
		utilization: 0.5,
		histogram: {
			runTime: { average: 1000 },
			waitTime: { average: 100 },
		},
	}));
});

describe('ScriptEngineService', () => {
	let service: ScriptEngineService;
	let logger: Logger;
	let metrics: ScriptEngineMetrics;
	let poolMetrics: PiscinaPoolMetrics;

	beforeEach(() => {
		jest.clearAllMocks();
		logger = createMockLogger();
		metrics = createMockMetrics();
		poolMetrics = createMockPoolMetrics();

		service = new ScriptEngineService(
			logger,
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

		it('should pass workerData config to Piscina', async () => {
			await service.initialize();

			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalledWith(expect.objectContaining({
				workerData: {
					memoryLimitMb: 32,
					timeoutMs: 5000,
					contextRecycleInterval: 10000,
				},
			}));
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

		it('should reset state so re-initialize works', async () => {
			await service.initialize();
			await service.shutdown();
			await service.initialize();

			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalledTimes(2);
		});
	});

	describe('runWorkflowScript', () => {
		it('should dispatch to Piscina and return success result', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: true,
				output: { temperature: 25.3 },
				statePatch: { lastProcessed: 'now' },
				executionTime: 1.5,
			});

			const result = await service.runWorkflowScript(createWorkerInput());

			expect(result.success).toBe(true);
			expect(result.output).toEqual({ temperature: 25.3 });
			expect(result.statePatch).toEqual({ lastProcessed: 'now' });
		});

		it('should pass input directly to piscina.run()', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true });

			const input = createWorkerInput();
			await service.runWorkflowScript(input);

			expect(mockPiscinaRun).toHaveBeenCalledWith(input);
		});

		it('should auto-initialize on first call', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true });

			await service.runWorkflowScript(createWorkerInput());

			const Piscina = require('piscina');
			expect(Piscina).toHaveBeenCalled();
		});

		it('should throw OOMError when worker reports OOM', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'isolate disposed',
				isOOM: true,
			});

			await expect(service.runWorkflowScript(createWorkerInput()))
				.rejects
				.toThrow(OOMError);
		});

		it('should return failure result for script errors (non-OOM)', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'ReferenceError: x is not defined',
				executionTime: 0.5,
			});

			const result = await service.runWorkflowScript(createWorkerInput());

			expect(result.success).toBe(false);
			expect(result.error).toContain('ReferenceError');
		});

		it('should handle piscina.run() rejection (worker crash)', async () => {
			mockPiscinaRun.mockRejectedValue(new Error('Worker thread terminated'));

			const result = await service.runWorkflowScript(createWorkerInput());

			expect(result.success).toBe(false);
			expect(result.error).toContain('Worker thread terminated');
		});

		it('should re-throw OOMError on piscina.run() rejection if already OOMError', async () => {
			mockPiscinaRun.mockRejectedValue(new OOMError('V8 OOM'));

			await expect(service.runWorkflowScript(createWorkerInput()))
				.rejects
				.toThrow(OOMError);
		});
	});

	describe('metrics', () => {
		it('should observe piscina run duration on success', async () => {
			mockPiscinaRun.mockResolvedValue({ success: true });

			await service.runWorkflowScript(createWorkerInput());

			expect(poolMetrics.piscinaRunDuration.observe).toHaveBeenCalled();
			expect(poolMetrics.piscinaCompleted.inc).toHaveBeenCalled();
		});

		it('should increment scriptErrors on failure', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'TypeError: cannot read property',
			});

			await service.runWorkflowScript(createWorkerInput());

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ script_type: 'workflow_code', error_type: 'type' })
			);
		});

		it('should classify syntax errors correctly', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'SyntaxError: Unexpected token',
			});

			await service.runWorkflowScript(createWorkerInput());

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ error_type: 'syntax' })
			);
		});

		it('should classify reference errors correctly', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'x is not defined',
			});

			await service.runWorkflowScript(createWorkerInput());

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ error_type: 'reference' })
			);
		});

		it('should classify timeout errors correctly', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'Script execution timed out',
			});

			await service.runWorkflowScript(createWorkerInput());

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ error_type: 'timeout' })
			);
		});

		it('should classify unknown errors as other', async () => {
			mockPiscinaRun.mockResolvedValue({
				success: false,
				error: 'Something went wrong',
			});

			await service.runWorkflowScript(createWorkerInput());

			expect(metrics.scriptErrors.inc).toHaveBeenCalledWith(
				expect.objectContaining({ error_type: 'other' })
			);
		});

		it('should not observe metrics on piscina crash', async () => {
			mockPiscinaRun.mockRejectedValue(new Error('Worker crash'));

			await service.runWorkflowScript(createWorkerInput());

			expect(poolMetrics.piscinaCompleted.inc).not.toHaveBeenCalled();
		});
	});

	describe('getPoolStats', () => {
		it('should return empty stats before initialization', () => {
			const stats = service.getPoolStats();

			expect(stats.piscina.completed).toBe(0);
			expect(stats.piscina.threads).toBe(0);
			expect(stats.piscina.runTimeAvg).toBe(0);
			expect(stats.piscina.waitTimeAvg).toBe(0);
			expect(stats.piscina.utilization).toBe(0);
		});

		it('should return Piscina stats after initialization', async () => {
			await service.initialize();

			const stats = service.getPoolStats();

			expect(stats.piscina.completed).toBe(42);
			expect(stats.piscina.threads).toBe(2);
			expect(typeof stats.piscina.utilization).toBe('number');
		});

		it('should set piscinaWorkers gauge when poolMetrics provided', async () => {
			await service.initialize();

			service.getPoolStats();

			expect(poolMetrics.piscinaWorkers.set).toHaveBeenCalledWith(2);
		});
	});

	describe('without metrics', () => {
		it('should work without metrics (all optional)', async () => {
			const noMetricsService = new ScriptEngineService(
				logger,
				createMockPiscinaOptions(),
				createMockWorkerConfig(),
			);

			mockPiscinaRun.mockResolvedValue({ success: true });

			const result = await noMetricsService.runWorkflowScript(createWorkerInput());
			expect(result.success).toBe(true);

			await noMetricsService.shutdown();
		});
	});
});
