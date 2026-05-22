/**
 * WorkflowCodeConsumer Unit Tests
 *
 * Tests NATS consumer message handling: parse, dispatch, ack/nack, OOM handling, metrics.
 * Mocks: NatsBus, WorkflowScriptServicePort, ConfigModule.
 */

import type { Logger, ConfigModule } from '@mapexos/microservices';
import type { NatsBus, Message } from '@mapexos/infrastructure';
import type { WorkflowScriptServicePort } from '@modules/scripts/application/ports';
import type { WorkflowCodeConsumerDeps } from './WorkflowCodeConsumer.types';

import { initWorkflowCodeConsumer } from './WorkflowCodeConsumer';
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

const createMockConfig = (): ConfigModule => ({
	get: jest.fn((key: string) => {
		const map: Record<string, any> = {
			cpu_limit: 2,
			piscina_workers: 0,
			nats_consumer_batch_size: 0,
			nats_consumer_fetch_timeout: 0,
			nats_consumer_max_ack_pending: 0,
		};
		return map[key] ?? 0;
	}),
} as unknown as ConfigModule);

const createMockScriptService = (): WorkflowScriptServicePort => ({
	execute: jest.fn().mockResolvedValue({ success: true, output: {} }),
	invalidateNodes: jest.fn().mockResolvedValue(undefined),
	invalidateWorkflow: jest.fn().mockResolvedValue(undefined),
});

const createMockMessage = (data: Record<string, any>): Message => ({
	data: new TextEncoder().encode(JSON.stringify(data)),
	ack: jest.fn(),
	nack: jest.fn(),
	subject: 'workflow.js.code',
	sid: 1,
} as unknown as Message);

const validPayload = {
	orgId: 'org-001',
	workflowId: 'wf-001',
	nodeId: 'node-001',
	instanceId: 'inst-001',
	callbackSubject: 'workflow.resume.callback.inst-001',
	eventPayload: { data: 'test' },
	state: {},
	inputs: {},
	nodeOutputs: {},
};

describe('WorkflowCodeConsumer', () => {
	let logger: Logger;
	let natsBus: NatsBus;
	let scriptService: WorkflowScriptServicePort;
	let config: ConfigModule;
	let executionsTotal: { inc: jest.Mock };
	let executionDuration: { observe: jest.Mock };
	let batchSize: { observe: jest.Mock };

	// Captured handler from natsBus.startConsumer
	let batchHandler: (messages: Message[]) => Promise<void>;

	beforeEach(async () => {
		jest.clearAllMocks();
		logger = createMockLogger();
		config = createMockConfig();
		scriptService = createMockScriptService();
		executionsTotal = { inc: jest.fn() };
		executionDuration = { observe: jest.fn() };
		batchSize = { observe: jest.fn() };

		natsBus = {
			startConsumer: jest.fn().mockImplementation(async (opts: any) => {
				batchHandler = opts.batchMessageHandlerV2;
			}),
		} as unknown as NatsBus;

		const deps: WorkflowCodeConsumerDeps = {
			natsBus,
			logger,
			scriptService,
			config,
			executionsTotal: executionsTotal as any,
			executionDuration: executionDuration as any,
			batchSize: batchSize as any,
		};

		await initWorkflowCodeConsumer(deps);
	});

	it('should register a NATS consumer with correct config', () => {
		expect(natsBus.startConsumer).toHaveBeenCalledWith(expect.objectContaining({
			stream: 'WORKFLOW-JS-CODE',
			subject: 'workflow.js.code',
			durable: 'js-workflow-executor-code',
		}));
	});

	it('should observe batch size', async () => {
		const messages = [createMockMessage(validPayload)];
		await batchHandler(messages);

		expect(batchSize.observe).toHaveBeenCalledWith(1);
	});

	describe('successful processing', () => {
		it('should parse message, call scriptService.execute, and ack', async () => {
			const msg = createMockMessage(validPayload);
			await batchHandler([msg]);

			expect(scriptService.execute).toHaveBeenCalledWith(expect.objectContaining({
				orgId: 'org-001',
				workflowId: 'wf-001',
				nodeId: 'node-001',
				instanceId: 'inst-001',
			}));
			expect(msg.ack).toHaveBeenCalled();
			expect(msg.nack).not.toHaveBeenCalled();
			expect(executionsTotal.inc).toHaveBeenCalledWith({ status: 'success' });
		});

		it('should observe execution duration', async () => {
			await batchHandler([createMockMessage(validPayload)]);

			expect(executionDuration.observe).toHaveBeenCalled();
			const duration = (executionDuration.observe as jest.Mock).mock.calls[0][0];
			expect(typeof duration).toBe('number');
			expect(duration).toBeGreaterThanOrEqual(0);
		});
	});

	describe('invalid messages', () => {
		it('should ack messages with missing required fields', async () => {
			const msg = createMockMessage({ orgId: 'org-001' }); // Missing workflowId, nodeId, instanceId
			await batchHandler([msg]);

			expect(scriptService.execute).not.toHaveBeenCalled();
			expect(msg.ack).toHaveBeenCalled();
			expect(logger.warn).toHaveBeenCalledWith(
				expect.stringContaining('missing required fields')
			);
		});
	});

	describe('OOM handling', () => {
		it('should nack message on OOMError for retry', async () => {
			(scriptService.execute as jest.Mock).mockRejectedValue(new OOMError('V8 OOM'));

			const msg = createMockMessage(validPayload);
			await batchHandler([msg]);

			expect(msg.nack).toHaveBeenCalled();
			expect(msg.ack).not.toHaveBeenCalled();
			expect(executionsTotal.inc).toHaveBeenCalledWith({ status: 'oom' });
		});
	});

	describe('non-OOM error handling', () => {
		it('should ack message on non-OOM error (callback already published)', async () => {
			(scriptService.execute as jest.Mock).mockRejectedValue(new Error('Unexpected error'));

			const msg = createMockMessage(validPayload);
			await batchHandler([msg]);

			expect(msg.ack).toHaveBeenCalled();
			expect(msg.nack).not.toHaveBeenCalled();
			expect(executionsTotal.inc).toHaveBeenCalledWith({ status: 'error' });
		});
	});

	describe('batch processing', () => {
		it('should process multiple messages independently', async () => {
			(scriptService.execute as jest.Mock)
				.mockResolvedValueOnce({ success: true })
				.mockRejectedValueOnce(new OOMError('OOM'))
				.mockResolvedValueOnce({ success: true });

			const msgs = [
				createMockMessage({ ...validPayload, instanceId: 'inst-001' }),
				createMockMessage({ ...validPayload, instanceId: 'inst-002' }),
				createMockMessage({ ...validPayload, instanceId: 'inst-003' }),
			];

			await batchHandler(msgs);

			expect(msgs[0].ack).toHaveBeenCalled();
			expect(msgs[1].nack).toHaveBeenCalled(); // OOM
			expect(msgs[2].ack).toHaveBeenCalled();
		});
	});

	describe('without optional metrics', () => {
		it('should work without metrics (all optional)', async () => {
			const noBus = {
				startConsumer: jest.fn().mockImplementation(async (opts: any) => {
					batchHandler = opts.batchMessageHandlerV2;
				}),
			} as unknown as NatsBus;

			await initWorkflowCodeConsumer({
				natsBus: noBus,
				logger,
				scriptService,
				config,
				// No metrics
			});

			const msg = createMockMessage(validPayload);
			await batchHandler([msg]); // Should not throw

			expect(msg.ack).toHaveBeenCalled();
		});
	});
});
