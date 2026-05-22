import type { PiscinaWorkerInput } from '@modules/engine/infrastructure/worker';
import type { AssetScripts } from '@modules/scripts/domain/types';
import type { ScriptServiceInternalDeps, BatchMessageInput, BatchMessageResult } from '@modules/scripts/application/types';

import { isEmpty } from 'lodash';

import { PUBLIC_ORG_ID } from '@modules/scripts/application/constants';

/** Metadata tracked per input for publishing after engine dispatch */
interface EventMeta {
	inputIndex: number;
	assetUUID: string;
	assetId: string;
	debugEnabled: boolean;
	sourceType: string;
	eventTrackerId: string;
	orgId: string;
	pathKey: string;
	name: string;
	description: string;
	rawPayload: any;
}

/** Result from fetchAssetData — scripts + metadata needed for execution */
interface AssetFetchResult {
	scripts: AssetScripts;
	assetId: string;
	debugEnabled: boolean;
	healthMonitorEnabled: boolean;
	heartbeatMode: 'implicit' | 'explicit';
	assetMetadata: { pathKey: string; name: string; description: string };
}

/**
 * Shared batch processing pipeline.
 *
 * Called by both handleMqttBatch and handleHttpBatch after they normalize
 * their inputs into the common BatchMessageInput contract.
 *
 * Flow: fetch assets → build worker inputs → engine.runBatch() →
 *       publish results → publish heartbeats → flush → return results.
 */
export async function processBatch(
	deps: ScriptServiceInternalDeps,
	inputs: BatchMessageInput[],
): Promise<BatchMessageResult[]> {
	deps.logger.info({ inputCount: inputs.length }, '[TRACE:processBatch] enter');

	const batchTimestamp = new Date().toISOString();
	const results: BatchMessageResult[] = [];

	/** Fetch unique assets in parallel (deduplicated by orgId/assetUUID) */
	const assetCache = new Map<string, AssetFetchResult>();
	const uniqueKeys = [...new Set(inputs.map(i => `${i.orgId}/${i.assetUUID}`))];

	deps.logger.info({ uniqueKeys }, '[TRACE:processBatch] fetching asset data');

	const fetchResults = await Promise.all(
		uniqueKeys.map(async (key) => {
			try {
				const result = await fetchAssetData(deps, key);
				deps.logger.info({ key, debugEnabled: result.debugEnabled, healthMonitorEnabled: result.healthMonitorEnabled, heartbeatMode: result.heartbeatMode, hasTransform: !!result.scripts.transform }, '[TRACE:processBatch] fetchAssetData OK');
				return { key, result, error: null };
			} catch (error) {
				deps.logger.error({ key, err: error instanceof Error ? error.message : String(error) }, '[TRACE:processBatch] fetchAssetData FAILED');
				return { key, result: null, error };
			}
		}),
	);

	for (const { key, result } of fetchResults) {
		if (result) assetCache.set(key, result);
	}

	deps.logger.info({ cachedCount: assetCache.size, uniqueKeysCount: uniqueKeys.length }, '[TRACE:processBatch] asset cache populated');

	/** Build worker inputs + metadata for publishing */
	const workerInputs: PiscinaWorkerInput[] = [];
	const eventMeta: EventMeta[] = [];

	for (const input of inputs) {
		const assetKey = `${input.orgId}/${input.assetUUID}`;
		const cached = assetCache.get(assetKey);

		if (!cached) {
			results.push({ index: input.index, success: false, error: `Asset not found: ${assetKey}` });
			continue;
		}

		const { scripts, assetId, debugEnabled, healthMonitorEnabled, heartbeatMode, assetMetadata } = cached;

		workerInputs.push({
			rawPayload: input.event,
			scripts: {
				decode: scripts.decode || undefined,
				validation: scripts.validation || undefined,
				transform: scripts.transform || undefined,
			},
			templateId: assetId,
		});

		eventMeta.push({
			inputIndex: input.index,
			assetUUID: input.assetUUID,
			assetId,
			debugEnabled,
			sourceType: input.sourceType,
			eventTrackerId: input.eventTrackerId,
			orgId: input.orgId,
			pathKey: assetMetadata.pathKey,
			name: assetMetadata.name,
			description: assetMetadata.description,
			rawPayload: input.event,
		});

		/** Heartbeat is independent of script success — device is alive, data arrived.
		 *  Truth table:
		 *    enabled=false                  → skip (monitoring off)
		 *    enabled=true + implicit/missing → publish (current behavior)
		 *    enabled=true + explicit         → skip; liveness comes from a different
		 *                                      path chosen by the asset's protocol:
		 *                                        • mqtt: NATS broker presence
		 *                                          ($SYS.ACCOUNT.*.CONNECT/DISCONNECT)
		 *                                        • http: POST /api/v1/heartbeat?ds=…
		 *                                          with body { assetUUID } */
		if (healthMonitorEnabled && heartbeatMode === 'implicit') {
			deps.eventPublisher.publishHeartbeat({
				orgId: input.orgId,
				assetUUID: input.assetUUID,
				pathKey: assetMetadata.pathKey,
			});
			deps.metrics?.heartbeatsPublished?.inc();
		} else if (!healthMonitorEnabled) {
			deps.metrics?.heartbeatsSkipped?.inc({ reason: 'disabled' });
		} else {
			deps.metrics?.heartbeatsSkipped?.inc({ reason: 'explicit_mode' });
		}
	}

	deps.logger.info({ workerInputsCount: workerInputs.length, eventMetaCount: eventMeta.length }, '[TRACE:processBatch] before engine.runBatch');

	if (workerInputs.length === 0) {
		deps.logger.warn('[TRACE:processBatch] workerInputs empty — skipping engine.runBatch');
		return results;
	}

	/** Dispatch to engine */
	const engineResults = await deps.scriptEngine.runBatch(workerInputs);
	deps.logger.info({ engineResultsCount: engineResults.length, successCount: engineResults.filter(r => r.success).length, failureCount: engineResults.filter(r => !r.success).length }, '[TRACE:processBatch] engine.runBatch returned');

	/** Publish results + heartbeats */
	for (let i = 0; i < engineResults.length; i++) {
		const result = engineResults[i];
		const meta = eventMeta[i];

		deps.logger.info({ i, success: result.success, debugEnabled: meta.debugEnabled, error: result.error, finalPayloadPreview: JSON.stringify(result.finalPayload).slice(0, 300) }, '[TRACE:processBatch] engine result');

		if (meta.debugEnabled) {
			deps.logger.info({ trackId: meta.eventTrackerId, asset: meta.assetUUID }, '[TRACE:processBatch] publishing raw event');
			deps.eventPublisher.publishRawEvent({
				eventTrackerId: meta.eventTrackerId,
				assetUUID: meta.assetUUID,
				orgId: meta.orgId,
				pathKey: meta.pathKey,
				name: meta.name,
				description: meta.description,
				event: meta.rawPayload,
				sourceType: meta.sourceType,
				timestamp: batchTimestamp,
			});
		}

		if (result.success) {
			deps.logger.info({ trackId: meta.eventTrackerId, asset: meta.assetUUID }, '[TRACE:processBatch] publishing result to route.execute');
			deps.eventPublisher.publishResult({
				assetUUID: meta.assetUUID,
				assetId: meta.assetId,
				orgId: meta.orgId,
				pathKey: meta.pathKey,
				eventTrackerId: meta.eventTrackerId,
				dataSource: {
					id: meta.assetId,
					orgId: meta.orgId,
					pathKey: meta.pathKey,
					name: meta.name,
					description: meta.description,
				},
				event: result.finalPayload,
			});

			if (meta.debugEnabled) {
				deps.eventPublisher.publishExecutionLog({
					eventTrackerId: meta.eventTrackerId,
					assetUUID: meta.assetUUID,
					orgId: meta.orgId,
					pathKey: meta.pathKey,
					name: meta.name,
					description: meta.description,
					execution: { success: true, failedAt: '', totalExecutionTime: result.totalPipelineTime ?? 0, error: '' },
					event: result.finalPayload,
					timestamp: batchTimestamp,
				});
			}

			results.push({ index: meta.inputIndex, success: true });
		} else {
			deps.eventPublisher.publishExecutionLog({
				eventTrackerId: meta.eventTrackerId,
				assetUUID: meta.assetUUID,
				orgId: meta.orgId,
				pathKey: meta.pathKey,
				name: meta.name,
				description: meta.description,
				execution: { success: false, failedAt: result.failedAt ?? '', totalExecutionTime: result.totalPipelineTime ?? 0, error: result.error ?? '' },
				event: meta.rawPayload,
				timestamp: batchTimestamp,
			});

			results.push({
				index: meta.inputIndex,
				success: false,
				error: result.error ?? 'Script error',
				isOOM: result.isOOM,
			});
		}
	}

	/** Flush all publishes in single TCP roundtrip */
	deps.logger.info('[TRACE:processBatch] calling flush');
	await deps.eventPublisher.flush();
	deps.logger.info({ resultsCount: results.length }, '[TRACE:processBatch] flush done, returning results');

	return results;
}

/**
 * Fetches asset and template data from cache ports.
 * Replicates the logic from ScriptService.fetchAssetScripts().
 */
async function fetchAssetData(deps: ScriptServiceInternalDeps, assetCacheKey: string): Promise<AssetFetchResult> {
	const asset = await deps.assetCachePort.get(assetCacheKey);
	if (!asset) {
		throw new Error(`Asset not found in cache: ${assetCacheKey}`);
	}

	if (!asset.assetTemplateId) {
		throw new Error(`Asset ${asset.uuid} has no template assigned`);
	}

	const templateOrgId = asset.assetTemplateOrgId || PUBLIC_ORG_ID;
	const templateCacheKey = `${templateOrgId}/${asset.assetTemplateId}`;

	const template = await deps.templateCachePort.get(templateCacheKey);
	if (!template) {
		throw new Error(`Template not found in cache: ${templateCacheKey}`);
	}

	const scripts: AssetScripts = {
		decode: template.scriptProcessor || '',
		validation: template.scriptValidator,
		transform: template.scriptConversion,
	};

	if (isEmpty(scripts.transform)) {
		throw new Error(`Invalid or missing transform script for template: ${asset.assetTemplateId}`);
	}

	return {
		scripts,
		assetId: asset.id,
		debugEnabled: asset.debugEnabled,
		healthMonitorEnabled: asset.healthMonitor?.enabled === true,
		heartbeatMode: (asset.healthMonitor?.heartbeatMode ?? 'implicit') as 'implicit' | 'explicit',
		assetMetadata: {
			pathKey: asset.pathKey || '',
			name: asset.name || '',
			description: asset.description || '',
		},
	};
}
