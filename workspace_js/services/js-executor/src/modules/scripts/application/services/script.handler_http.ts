import type { Message } from '@mapexos/infrastructure';
import type { ScriptServiceInternalDeps, BatchMessageInput, BatchMessageResult } from '@modules/scripts/application/types';
import type { ScriptProcessorMessage } from '@modules/scripts/application/types';

import { processBatch } from './script.handler_batch';

/**
 * Handles a batch of HTTP datasource messages.
 * Parses JSON payloads, validates required fields, resolves assetUUID,
 * normalizes into BatchMessageInput[], then delegates to processBatch().
 *
 * @param deps - Service internal dependencies
 * @param messages - Raw NATS messages from PROCESSOR-JS-EXECUTE stream
 * @param resolveAssetUUID - Function to resolve assetUUID from message (injected from ScriptService)
 * @returns Per-message results for ACK/Nack
 */
export async function handleHttpBatch(
	deps: ScriptServiceInternalDeps,
	messages: Message[],
	resolveAssetUUID: (message: ScriptProcessorMessage) => Promise<string>,
): Promise<BatchMessageResult[]> {
	const parseResults: BatchMessageResult[] = [];
	const validInputs: BatchMessageInput[] = [];

	for (let i = 0; i < messages.length; i++) {
		const msg = messages[i];

		try {
			const payload: ScriptProcessorMessage = JSON.parse(new TextDecoder().decode(msg.data));

			if (!payload.eventTrackerId) {
				payload.eventTrackerId = `seq-${msg.streamSequence}`;
			}

			msg.orgId = payload.dataSource?.orgId ?? '';
			msg.pathKey = payload.dataSource?.pathKey ?? '';

			if (!payload.event) {
				parseResults.push({ index: i, success: false, error: 'Missing required field: event', isPermanent: true });
				continue;
			}
			if (!payload.dataSource) {
				parseResults.push({ index: i, success: false, error: 'Missing required field: dataSource', isPermanent: true });
				continue;
			}

			const orgId = payload.dataSource.orgId || '';
			if (!orgId) {
				parseResults.push({ index: i, success: false, error: 'Missing required field: dataSource.orgId', isPermanent: true });
				continue;
			}

			const assetUUID = await resolveAssetUUID(payload);

			validInputs.push({
				index: i,
				orgId,
				assetUUID,
				event: payload.event,
				sourceType: payload.sourceType,
				eventTrackerId: payload.eventTrackerId!,
				dataSource: payload.dataSource,
			});
			
		} catch (error) {
			if (error instanceof SyntaxError) {
				parseResults.push({ index: i, success: false, error: `Invalid JSON: ${error.message}`, isPermanent: true });
			} else {
				parseResults.push({ index: i, success: false, error: error instanceof Error ? error.message : String(error) });
			}
		}
	}

	if (validInputs.length === 0) return parseResults;

	const processResults = await processBatch(deps, validInputs);
	return [...parseResults, ...processResults];
}
