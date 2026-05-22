/**
 * DefinitionInvalidateConsumer Unit Tests
 *
 * Tests FANOUT consumer: payload parsing, cache invalidation delegation, error handling.
 * Mocks: NatsBus, WorkflowScriptServicePort.
 */

import type { Logger } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { WorkflowScriptServicePort } from '@modules/scripts/application/ports';

import { initDefinitionInvalidateConsumer } from './DefinitionInvalidateConsumer';

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

const createMockScriptService = (): WorkflowScriptServicePort => ({
	execute: jest.fn().mockResolvedValue({ success: true }),
	invalidateNodes: jest.fn().mockResolvedValue(undefined),
	invalidateWorkflow: jest.fn().mockResolvedValue(undefined),
});

describe('DefinitionInvalidateConsumer', () => {
	let logger: Logger;
	let natsBus: NatsBus;
	let scriptService: WorkflowScriptServicePort;

	let fanoutHandler: (data: Uint8Array) => Promise<void>;

	beforeEach(async () => {
		jest.clearAllMocks();
		logger = createMockLogger();
		scriptService = createMockScriptService();

		natsBus = {
			subscribeFanout: jest.fn().mockImplementation(async (opts: any) => {
				fanoutHandler = opts.handler;
			}),
		} as unknown as NatsBus;

		await initDefinitionInvalidateConsumer({
			natsBus,
			logger,
			scriptService,
			serviceName: 'js-workflow-executor',
		});
	});

	it('should subscribe to FANOUT with correct config', () => {
		expect(natsBus.subscribeFanout).toHaveBeenCalledWith(expect.objectContaining({
			stream: 'FANOUT',
			subject: 'fanout.workflow.definition.invalidate',
			serviceName: 'js-workflow-executor',
		}));
	});

	it('should call scriptService.invalidateNodes on payload with nodeIds', async () => {
		const payload = { orgId: 'org-001', definitionId: 'def-001', nodeIds: ['node-1', 'node-2'] };
		const data = new TextEncoder().encode(JSON.stringify(payload));

		await fanoutHandler(data);

		expect(scriptService.invalidateNodes).toHaveBeenCalledWith('org-001', 'def-001', ['node-1', 'node-2']);
		expect(scriptService.invalidateWorkflow).not.toHaveBeenCalled();
	});

	it('should fall back to invalidateWorkflow when nodeIds is empty', async () => {
		const payload = { orgId: 'org-001', definitionId: 'def-001', nodeIds: [] };
		const data = new TextEncoder().encode(JSON.stringify(payload));

		await fanoutHandler(data);

		expect(scriptService.invalidateWorkflow).toHaveBeenCalledWith('org-001', 'def-001');
		expect(scriptService.invalidateNodes).not.toHaveBeenCalled();
	});

	it('should fall back to invalidateWorkflow when nodeIds is absent', async () => {
		const payload = { orgId: 'org-001', definitionId: 'def-001' };
		const data = new TextEncoder().encode(JSON.stringify(payload));

		await fanoutHandler(data);

		expect(scriptService.invalidateWorkflow).toHaveBeenCalledWith('org-001', 'def-001');
		expect(scriptService.invalidateNodes).not.toHaveBeenCalled();
	});

	it('should skip invalidation when orgId is missing', async () => {
		const payload = { definitionId: 'def-001' };
		const data = new TextEncoder().encode(JSON.stringify(payload));

		await fanoutHandler(data);

		expect(scriptService.invalidateNodes).not.toHaveBeenCalled();
		expect(scriptService.invalidateWorkflow).not.toHaveBeenCalled();
		expect(logger.warn).toHaveBeenCalledWith(
			expect.stringContaining('missing orgId or definitionId')
		);
	});

	it('should skip invalidation when definitionId is missing', async () => {
		const payload = { orgId: 'org-001' };
		const data = new TextEncoder().encode(JSON.stringify(payload));

		await fanoutHandler(data);

		expect(scriptService.invalidateNodes).not.toHaveBeenCalled();
		expect(scriptService.invalidateWorkflow).not.toHaveBeenCalled();
	});

	it('should handle JSON parse error gracefully', async () => {
		const data = new TextEncoder().encode('invalid json');

		await fanoutHandler(data); // Should not throw

		expect(scriptService.invalidateNodes).not.toHaveBeenCalled();
		expect(logger.error).toHaveBeenCalledWith(
			expect.stringContaining('Failed to process invalidation')
		);
	});

	it('should handle scriptService.invalidateNodes failure gracefully', async () => {
		(scriptService.invalidateNodes as jest.Mock).mockRejectedValue(new Error('Cache error'));

		const payload = { orgId: 'org-001', definitionId: 'def-001', nodeIds: ['node-1'] };
		const data = new TextEncoder().encode(JSON.stringify(payload));

		await fanoutHandler(data); // Should not throw

		expect(logger.error).toHaveBeenCalledWith(
			expect.stringContaining('Failed to process invalidation')
		);
	});
});
