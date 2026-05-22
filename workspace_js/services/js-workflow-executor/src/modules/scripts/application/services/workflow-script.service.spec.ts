/**
 * WorkflowScriptService Unit Tests
 *
 * Tests the script execution orchestrator: source fetching, bytecode caching,
 * engine dispatch, callback publishing, and invalidation.
 *
 * Mocks: ScriptSourceCachePort, BytecodeCachePort, ScriptEngineServicePort, NatsBus.
 */

import type { Logger } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { BytecodeCachePort } from '@modules/engine/application/ports';
import type { WorkflowScriptInput, ScriptSourceCachePort } from '../ports';

import { WorkflowScriptService } from './workflow-script.service';

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

const createMockScriptSourceCache = (): ScriptSourceCachePort => ({
	get: jest.fn(),
	invalidate: jest.fn(),
});

const createMockBytecodeCache = (): BytecodeCachePort => ({
	get: jest.fn().mockResolvedValue(null),
	set: jest.fn().mockResolvedValue(undefined),
	invalidate: jest.fn().mockResolvedValue(undefined),
	invalidateNodes: jest.fn().mockResolvedValue(undefined),
	invalidateWorkflow: jest.fn().mockResolvedValue(undefined),
	buildCacheKey: jest.fn((ctx) => `${ctx.orgId}/${ctx.workflowId}/bytecode/${ctx.nodeId}`),
});

const createMockScriptEngine = (): ScriptEngineServicePort => ({
	initialize: jest.fn().mockResolvedValue(undefined),
	shutdown: jest.fn().mockResolvedValue(undefined),
	getPoolStats: jest.fn().mockReturnValue({ piscina: {} }),
	runWorkflowScript: jest.fn().mockResolvedValue({ success: true, output: { result: 'ok' } }),
});

const createMockNatsBus = (): NatsBus => ({
	publish: jest.fn().mockResolvedValue(undefined),
} as unknown as NatsBus);

const createInput = (overrides?: Partial<WorkflowScriptInput>): WorkflowScriptInput => ({
	orgId: 'org-001',
	pathKey: '000001',
	workflowId: 'wf-001',
	nodeId: 'node-001',
	instanceId: 'inst-001',
	callbackSubject: 'workflow.resume.callback.inst-001',
	eventPayload: { data: { temperature: 25.3 } },
	state: { processedCount: 0 },
	inputs: { threshold: 80 },
	nodeOutputs: { 'trigger-001': { source: 'http' } },
	...overrides,
});

describe('WorkflowScriptService', () => {
	let service: WorkflowScriptService;
	let logger: Logger;
	let scriptSourceCache: ScriptSourceCachePort;
	let bytecodeCache: BytecodeCachePort;
	let scriptEngine: ScriptEngineServicePort;
	let natsBus: NatsBus;

	beforeEach(() => {
		jest.clearAllMocks();
		logger = createMockLogger();
		scriptSourceCache = createMockScriptSourceCache();
		bytecodeCache = createMockBytecodeCache();
		scriptEngine = createMockScriptEngine();
		natsBus = createMockNatsBus();

		const callbackPublisher = { publishCallback: natsBus.publish.bind(natsBus) } as any;
		service = new WorkflowScriptService(
			logger,
			scriptSourceCache,
			bytecodeCache,
			scriptEngine,
			callbackPublisher,
		);
	});

	describe('execute', () => {
		it('should fetch source and bytecode in parallel, dispatch to engine, and publish callback', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue({
				data: 'const result = { output: {}, statePatch: {} };',
				tier: 0,
			});
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({
				success: true,
				output: { temperature: 25.3 },
				statePatch: { lastProcessed: 'now' },
			});

			const input = createInput();
			const result = await service.execute(input);

			expect(result.success).toBe(true);
			expect(scriptSourceCache.get).toHaveBeenCalledWith('org-001/wf-001/scripts/node-001');
			expect(bytecodeCache.get).toHaveBeenCalledWith('org-001/wf-001/bytecode/node-001');
			expect(scriptEngine.runWorkflowScript).toHaveBeenCalledWith(expect.objectContaining({
				script: 'const result = { output: {}, statePatch: {} };',
				cacheKey: 'org-001/wf-001/scripts/node-001',
				event: input.eventPayload,
				state: input.state,
				inputs: input.inputs,
				nodes: input.nodeOutputs,
			}));
			expect(natsBus.publish).toHaveBeenCalledWith(
				'workflow.resume.callback.inst-001',
				expect.objectContaining({
					instanceId: 'inst-001',
					nodeId: 'node-001',
					status: 'success',
					output: { temperature: 25.3 },
					statePatch: { lastProcessed: 'now' },
				}),
			);
		});

		it('should publish error callback when script source is not found', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue(null);

			const input = createInput();
			const result = await service.execute(input);

			expect(result.success).toBe(false);
			expect(result.error).toContain('Script source not found');
			expect(scriptEngine.runWorkflowScript).not.toHaveBeenCalled();
			expect(natsBus.publish).toHaveBeenCalledWith(
				'workflow.resume.callback.inst-001',
				expect.objectContaining({
					status: 'error',
					error: expect.objectContaining({ code: 'SCRIPT_NOT_FOUND' }),
				}),
			);
		});

		it('should publish error callback when script source is empty', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('');

			const result = await service.execute(createInput());

			expect(result.success).toBe(false);
			expect(natsBus.publish).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({
					status: 'error',
					error: expect.objectContaining({ code: 'SCRIPT_NOT_FOUND' }),
				}),
			);
		});

		it('should publish error callback on script execution failure', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('bad code');
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({
				success: false,
				error: 'ReferenceError: x is not defined',
			});

			const result = await service.execute(createInput());

			expect(result.success).toBe(false);
			expect(natsBus.publish).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({
					status: 'error',
					error: expect.objectContaining({ code: 'SCRIPT_ERROR', message: 'ReferenceError: x is not defined' }),
				}),
			);
		});

		it('should pass cached bytecode to engine when available', async () => {
			const fakeBytecode = Buffer.from('fake-bytecode');
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script code');
			(bytecodeCache.get as jest.Mock).mockResolvedValue(fakeBytecode);
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({ success: true });

			await service.execute(createInput());

			expect(scriptEngine.runWorkflowScript).toHaveBeenCalledWith(
				expect.objectContaining({
					cachedBytecode: expect.any(ArrayBuffer),
				}),
			);
		});

		it('should pass undefined cachedBytecode when bytecode cache returns null', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script code');
			(bytecodeCache.get as jest.Mock).mockResolvedValue(null);
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({ success: true });

			await service.execute(createInput());

			expect(scriptEngine.runWorkflowScript).toHaveBeenCalledWith(
				expect.objectContaining({
					cachedBytecode: undefined,
				}),
			);
		});

		it('should store newBytecode in cache when returned by engine (fire-and-forget)', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script');
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({
				success: true,
				output: {},
				newBytecode: new ArrayBuffer(100),
			});

			await service.execute(createInput());

			expect(bytecodeCache.set).toHaveBeenCalledWith(
				'org-001/wf-001/bytecode/node-001',
				expect.any(ArrayBuffer),
			);
		});

		it('should not store bytecode when newBytecode is absent', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script');
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({
				success: true,
				output: {},
			});

			await service.execute(createInput());

			expect(bytecodeCache.set).not.toHaveBeenCalled();
		});

		it('should handle bytecode store failure gracefully (fire-and-forget)', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script');
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({
				success: true,
				newBytecode: new ArrayBuffer(50),
			});
			(bytecodeCache.set as jest.Mock).mockRejectedValue(new Error('MinIO down'));

			// Should not throw
			const result = await service.execute(createInput());
			expect(result.success).toBe(true);
		});

		it('should handle TieredCache fetch error gracefully', async () => {
			(scriptSourceCache.get as jest.Mock).mockRejectedValue(new Error('Cache connection failed'));

			const result = await service.execute(createInput());

			expect(result.success).toBe(false);
			expect(natsBus.publish).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({ status: 'error' }),
			);
		});

		it('should handle bytecode fetch error gracefully (returns null)', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script');
			(bytecodeCache.get as jest.Mock).mockRejectedValue(new Error('Bytecode error'));
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({ success: true });

			const result = await service.execute(createInput());

			// Should still succeed — bytecode is optional
			expect(result.success).toBe(true);
		});

		it('should handle Buffer data from TieredCache', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue({
				data: Buffer.from('const result = { output: "ok", statePatch: {} };'),
				tier: 1,
			});
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({ success: true });

			const result = await service.execute(createInput());

			expect(result.success).toBe(true);
			expect(scriptEngine.runWorkflowScript).toHaveBeenCalledWith(
				expect.objectContaining({
					script: 'const result = { output: "ok", statePatch: {} };',
				}),
			);
		});

		it('should handle NATS publish failure gracefully', async () => {
			(scriptSourceCache.get as jest.Mock).mockResolvedValue('script');
			(scriptEngine.runWorkflowScript as jest.Mock).mockResolvedValue({ success: true, output: {} });
			(natsBus.publish as jest.Mock).mockRejectedValue(new Error('NATS timeout'));

			// Should not throw — publish failure is logged but not fatal
			const result = await service.execute(createInput());
			expect(result.success).toBe(true);
		});
	});

	describe('invalidateWorkflow', () => {
		it('should delegate to bytecodeCache.invalidateWorkflow', async () => {
			await service.invalidateWorkflow('org-001', 'wf-001');

			expect(bytecodeCache.invalidateWorkflow).toHaveBeenCalledWith('org-001', 'wf-001');
		});

		it('should log invalidation', async () => {
			await service.invalidateWorkflow('org-001', 'wf-001');

			expect(logger.info).toHaveBeenCalledWith(
				expect.stringContaining('wf-001')
			);
		});
	});
});
