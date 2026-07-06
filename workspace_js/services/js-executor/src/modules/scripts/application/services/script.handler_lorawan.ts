import type { Message } from '@mapexos/infrastructure';
import type { ScriptServiceInternalDeps, BatchMessageInput, BatchMessageResult } from '@modules/scripts/application/types';
import type { LorawanUplinkEnvelope, LorawanDecodeInput } from '@modules/scripts/application/types';

import { processBatch } from './script.handler_batch';

/**
 * Handles a batch of LoRaWAN uplinks republished by the mapexLNS adapter.
 *
 * Mapex contract: NATS subjects are agnostic, every routing field
 * (orgId, assetUUID) lives in the JSON payload the LNS adapter (the
 * Anti-Corruption Layer over The Things Stack) publishes on
 * ${env}.mapexos.lorawan.data.> — the same convention the broker plugin
 * follows for MQTT. Unlike MQTT, the device sends raw bytes (FRMPayload),
 * not JSON: we hand the decode step the bytes + frame metadata so the
 * device codec (decodeUplink-style) turns them into structured data.
 *
 * @param deps - Service internal dependencies
 * @param messages - Raw NATS messages from LORAWAN-DATA stream
 * @returns Per-message results for ACK/Nack
 */
export async function handleLorawanBatch(
	deps: ScriptServiceInternalDeps,
	messages: Message[],
): Promise<BatchMessageResult[]> {
	deps.logger.info({ count: messages.length }, '[TRACE:LorawanBatch] enter handleLorawanBatch');

	const parseResults: BatchMessageResult[] = [];
	const validInputs: BatchMessageInput[] = [];

	for (let i = 0; i < messages.length; i++) {
		const msg = messages[i];

		try {
			const rawData = new TextDecoder().decode(msg.data);
			deps.logger.info({ idx: i, subject: msg.subject, rawLen: rawData.length, rawPreview: rawData.slice(0, 300) }, '[TRACE:LorawanBatch] raw message');

			const envelope = JSON.parse(rawData) as LorawanUplinkEnvelope;

			const orgId = envelope.orgId ?? '';
			const assetUUID = envelope.assetUUID ?? '';

			if (!orgId || !assetUUID) {
				deps.logger.warn({ idx: i, orgId, assetUUID }, '[TRACE:LorawanBatch] missing orgId/assetUUID — pushing parse failure');
				parseResults.push({
					index: i,
					success: false,
					error: 'Invalid lorawan envelope: missing orgId or assetUUID',
					isPermanent: true,
				});
				continue;
			}

			const decodeInput = buildDecodeInput(envelope);

			deps.logger.info(
				{ idx: i, orgId, assetUUID, fPort: decodeInput.fPort, fCnt: decodeInput.fCnt, byteLen: decodeInput.bytes.length },
				'[TRACE:LorawanBatch] parsed envelope',
			);

			msg.orgId = orgId;
			msg.pathKey = '';

			validInputs.push({
				index: i,
				orgId,
				assetUUID,
				event: decodeInput,
				sourceType: 'lorawan',
				eventTrackerId: `seq-${msg.streamSequence}`,
			});
		} catch (error) {
			deps.logger.error({ idx: i, err: error instanceof Error ? error.message : String(error) }, '[TRACE:LorawanBatch] parse error');
			if (error instanceof SyntaxError) {
				parseResults.push({ index: i, success: false, error: `Invalid JSON: ${error.message}`, isPermanent: true });
			} else {
				parseResults.push({ index: i, success: false, error: error instanceof Error ? error.message : String(error) });
			}
		}
	}

	deps.logger.info({ valid: validInputs.length, parseFails: parseResults.length }, '[TRACE:LorawanBatch] before processBatch');

	if (validInputs.length === 0) {
		deps.logger.warn('[TRACE:LorawanBatch] no valid inputs — returning parse failures');
		return parseResults;
	}

	const processResults = await processBatch(deps, validInputs);
	deps.logger.info({ processResultsCount: processResults.length, successCount: processResults.filter(r => r.success).length }, '[TRACE:LorawanBatch] after processBatch');
	return [...parseResults, ...processResults];
}

/**
 * Builds the codec input the decode step sees as the global `payload`.
 * The FRMPayload travels base64-encoded (Go marshals []byte as base64);
 * we decode it to a plain byte array so the device codec can read
 * `payload.bytes` / `payload.fPort` the TTN/ChirpStack way. A missing or
 * malformed payload yields an empty byte array — the codec decides what
 * that means for the device.
 */
function buildDecodeInput(envelope: LorawanUplinkEnvelope): LorawanDecodeInput {
	return {
		bytes: decodeFrmPayload(envelope.payload),
		fPort: envelope.fPort ?? 0,
		fCnt: envelope.fCnt ?? 0,
		rxInfo: envelope.rxInfo,
	};
}

/**
 * Turns the base64 FRMPayload string into a plain array of byte values.
 * Returns an empty array when the payload is absent or not decodable.
 */
function decodeFrmPayload(raw: unknown): number[] {
	if (typeof raw !== 'string' || raw.length === 0) return [];
	try {
		return Array.from(Buffer.from(raw, 'base64'));
	} catch {
		return [];
	}
}
