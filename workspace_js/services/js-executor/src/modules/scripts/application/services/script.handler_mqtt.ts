import type { Message } from '@mapexos/infrastructure';
import type { ScriptServiceInternalDeps, BatchMessageInput, BatchMessageResult } from '@modules/scripts/application/types';

import { processBatch } from './script.handler_batch';

/**
 * Handles a batch of MQTT telemetry messages.
 *
 * Mapex contract: NATS subjects are agnostic, every routing field
 * (orgId, assetUUID) lives in the JSON payload the broker plugin
 * publishes. Same convention the downstream publishers
 * (route.execute, events.raw) follow. The device's own MQTT body
 * arrives base64-encoded inside `payload` because Go marshals []byte
 * to base64 — we decode + JSON.parse when both succeed so the script
 * engine sees the structured event the device intended.
 *
 * @param deps - Service internal dependencies
 * @param messages - Raw NATS messages from MQTT-DATA stream
 * @returns Per-message results for ACK/Nack
 */
export async function handleMqttBatch(
	deps: ScriptServiceInternalDeps,
	messages: Message[],
): Promise<BatchMessageResult[]> {
	deps.logger.info({ count: messages.length }, '[TRACE:MqttBatch] enter handleMqttBatch');

	const parseResults: BatchMessageResult[] = [];
	const validInputs: BatchMessageInput[] = [];

	for (let i = 0; i < messages.length; i++) {
		const msg = messages[i];

		try {
			const rawData = new TextDecoder().decode(msg.data);
			deps.logger.info({ idx: i, subject: msg.subject, rawLen: rawData.length, rawPreview: rawData.slice(0, 300) }, '[TRACE:MqttBatch] raw message');

			const ingress = JSON.parse(rawData) as {
				orgId?: string;
				assetUUID?: string;
				payload?: unknown;
				topic?: string;
				timestamp?: string;
			};

			const orgId = ingress.orgId ?? '';
			const assetUUID = ingress.assetUUID ?? '';
			const decodedDevice = decodeDevicePayload(ingress.payload);

			deps.logger.info(
				{ idx: i, orgId, assetUUID, topic: ingress.topic, deviceShape: typeof decodedDevice, devicePreview: JSON.stringify(decodedDevice).slice(0, 300) },
				'[TRACE:MqttBatch] parsed ingress',
			);

			if (!orgId || !assetUUID) {
				deps.logger.warn({ idx: i, orgId, assetUUID }, '[TRACE:MqttBatch] missing orgId/assetUUID — pushing parse failure');
				parseResults.push({
					index: i,
					success: false,
					error: 'Invalid ingress payload: missing orgId or assetUUID',
					isPermanent: true,
				});
				continue;
			}

			msg.orgId = orgId;
			msg.pathKey = '';

			validInputs.push({
				index: i,
				orgId,
				assetUUID,
				event: decodedDevice,
				sourceType: 'mqtt',
				eventTrackerId: `seq-${msg.streamSequence}`,
			});
		} catch (error) {
			deps.logger.error({ idx: i, err: error instanceof Error ? error.message : String(error) }, '[TRACE:MqttBatch] parse error');
			if (error instanceof SyntaxError) {
				parseResults.push({ index: i, success: false, error: `Invalid JSON: ${error.message}`, isPermanent: true });
			} else {
				parseResults.push({ index: i, success: false, error: error instanceof Error ? error.message : String(error) });
			}
		}
	}

	deps.logger.info({ valid: validInputs.length, parseFails: parseResults.length }, '[TRACE:MqttBatch] before processBatch');

	if (validInputs.length === 0) {
		deps.logger.warn('[TRACE:MqttBatch] no valid inputs — returning parse failures');
		return parseResults;
	}

	const processResults = await processBatch(deps, validInputs);
	deps.logger.info({ processResultsCount: processResults.length, successCount: processResults.filter(r => r.success).length }, '[TRACE:MqttBatch] after processBatch');
	return [...parseResults, ...processResults];
}

/**
 * Turns the broker's wire payload back into the device's original
 * JSON. Go's encoding/json marshals []byte as a base64 string, so
 * IngressMessage.Payload travels as a string here. We decode +
 * JSON.parse when both succeed; everything else passes through so
 * the engine still sees the raw value the device sent.
 */
function decodeDevicePayload(raw: unknown): unknown {
	if (typeof raw !== 'string') return raw;
	try {
		const decoded = Buffer.from(raw, 'base64').toString('utf-8');
		return JSON.parse(decoded);
	} catch {
		return raw;
	}
}
