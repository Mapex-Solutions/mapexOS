import 'reflect-metadata';

import { handleLorawanBatch } from './script.handler_lorawan';
import { processBatch } from './script.handler_batch';

import type { Logger } from '@mapexos/microservices';
import type { Message } from '@mapexos/infrastructure';
import type { BatchMessageInput, BatchMessageResult, ScriptServiceInternalDeps } from '@modules/scripts/application/types';

jest.mock('./script.handler_batch', () => ({
	processBatch: jest.fn(),
}));

const processBatchMock = processBatch as jest.MockedFunction<typeof processBatch>;

/**
 * handleLorawanBatch parses the LorawanUplinkEnvelope published by the mapexLNS
 * adapter and normalizes it into the shared batch contract. These tests pin the
 * contract: routing fields come from the body, the FRMPayload is base64-decoded
 * into bytes for the codec, and the source is tagged 'lorawan'.
 */
describe('script.handler_lorawan — handleLorawanBatch', () => {
	const logger: jest.Mocked<Logger> = {
		info: jest.fn(),
		error: jest.fn(),
		warn: jest.fn(),
		debug: jest.fn(),
	} as any;

	const deps = { logger } as unknown as ScriptServiceInternalDeps;

	const msg = (body: unknown, opts: { streamSequence?: number; subject?: string } = {}): Message => {
		const json = typeof body === 'string' ? body : JSON.stringify(body);
		return {
			data: new TextEncoder().encode(json),
			subject: opts.subject ?? 'dev.mapexos.lorawan.data',
			streamSequence: opts.streamSequence ?? 42,
		} as unknown as Message;
	};

	beforeEach(() => {
		jest.clearAllMocks();
		processBatchMock.mockResolvedValue([]);
	});

	it('maps a valid envelope: base64 payload → bytes, frame metadata, sourceType lorawan', async () => {
		const payload = Buffer.from([0x01, 0x02, 0x03]).toString('base64');
		const rxInfo = [{ rssi: -42, snr: 9.5 }];

		await handleLorawanBatch(deps, [
			msg({ orgId: 'org-1', assetUUID: 'dev-eui-1', payload, fPort: 10, fCnt: 7, rxInfo }, { streamSequence: 99 }),
		]);

		expect(processBatchMock).toHaveBeenCalledTimes(1);
		const inputs = processBatchMock.mock.calls[0][1] as BatchMessageInput[];
		expect(inputs).toHaveLength(1);
		expect(inputs[0]).toEqual({
			index: 0,
			orgId: 'org-1',
			assetUUID: 'dev-eui-1',
			event: { bytes: [1, 2, 3], fPort: 10, fCnt: 7, rxInfo },
			sourceType: 'lorawan',
			eventTrackerId: 'seq-99',
		});
	});

	it('defaults fPort/fCnt to 0 and bytes to [] when payload is absent', async () => {
		await handleLorawanBatch(deps, [msg({ orgId: 'org-1', assetUUID: 'dev-eui-1' })]);

		const inputs = processBatchMock.mock.calls[0][1] as BatchMessageInput[];
		expect(inputs[0].event).toEqual({ bytes: [], fPort: 0, fCnt: 0, rxInfo: undefined });
	});

	it('rejects permanently when orgId is missing and does not call processBatch', async () => {
		const results = await handleLorawanBatch(deps, [msg({ assetUUID: 'dev-eui-1', payload: '' })]);

		expect(processBatchMock).not.toHaveBeenCalled();
		expect(results).toEqual<BatchMessageResult[]>([
			{ index: 0, success: false, error: 'Invalid lorawan envelope: missing orgId or assetUUID', isPermanent: true },
		]);
	});

	it('rejects permanently when assetUUID is missing', async () => {
		const results = await handleLorawanBatch(deps, [msg({ orgId: 'org-1', payload: '' })]);

		expect(processBatchMock).not.toHaveBeenCalled();
		expect(results[0]).toMatchObject({ index: 0, success: false, isPermanent: true });
	});

	it('rejects permanently on invalid JSON (SyntaxError)', async () => {
		const results = await handleLorawanBatch(deps, [msg('{ not json')]);

		expect(processBatchMock).not.toHaveBeenCalled();
		expect(results[0].success).toBe(false);
		expect(results[0].isPermanent).toBe(true);
		expect(results[0].error).toMatch(/Invalid JSON/);
	});

	it('processes valid entries and rejects invalid ones in the same batch', async () => {
		processBatchMock.mockResolvedValue([{ index: 1, success: true }]);
		const payload = Buffer.from([0xff]).toString('base64');

		const results = await handleLorawanBatch(deps, [
			msg({ assetUUID: 'only-asset' }), // invalid: no orgId → index 0
			msg({ orgId: 'org-1', assetUUID: 'dev-eui-1', payload, fPort: 1, fCnt: 2 }), // valid → index 1
		]);

		// One valid input forwarded to processBatch
		const inputs = processBatchMock.mock.calls[0][1] as BatchMessageInput[];
		expect(inputs).toHaveLength(1);
		expect(inputs[0].index).toBe(1);
		expect(inputs[0].event).toEqual({ bytes: [255], fPort: 1, fCnt: 2, rxInfo: undefined });

		// Parse failure for the invalid one is preserved alongside the process result
		expect(results).toEqual(
			expect.arrayContaining([
				expect.objectContaining({ index: 0, success: false, isPermanent: true }),
				{ index: 1, success: true },
			]),
		);
	});
});
