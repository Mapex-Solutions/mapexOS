import type { DefinitionInvalidatePayload } from '@mapexos/schemas';

import type { DefinitionInvalidateConsumerDeps } from './DefinitionInvalidateConsumer.types';

import { DEFINITION_INVALIDATE_SUBJECT, FANOUT_STREAM } from './constants';

/**
 * Initializes the DefinitionInvalidate FANOUT consumer.
 *
 * Listens on fanout.workflow.definition.invalidate for cache invalidation events.
 * When a workflow definition is updated/deleted, all pods clear their caches
 * so the next execution fetches fresh script source from L2 (MinIO) or fallback HTTP.
 *
 * @param deps - Consumer dependencies
 */
export async function initDefinitionInvalidateConsumer(deps: DefinitionInvalidateConsumerDeps): Promise<void> {
	const { natsBus, logger, scriptService, serviceName } = deps;

	logger.info('[CONSUMER:DefinitionInvalidate] Initializing FANOUT consumer');

	await natsBus.subscribeFanout({
		stream: FANOUT_STREAM,
		subject: DEFINITION_INVALIDATE_SUBJECT,
		serviceName,
		handler: async (data: Uint8Array) => {
			try {
				const payload = JSON.parse(new TextDecoder().decode(data)) as DefinitionInvalidatePayload;

				if (!payload.orgId || !payload.definitionId) {
					logger.warn('[CONSUMER:DefinitionInvalidate] Invalid payload: missing orgId or definitionId');
					return;
				}

				if (payload.nodeIds && payload.nodeIds.length > 0) {
					// Granular invalidation — clear L0+L1 for specific nodes
					await scriptService.invalidateNodes(payload.orgId, payload.definitionId, payload.nodeIds);
					logger.debug(
						`[CONSUMER:DefinitionInvalidate] Invalidated ${payload.nodeIds.length} node(s) for definition ${payload.definitionId} (org: ${payload.orgId})`
					);
				} else {
					// Fallback — no nodeIds available, rely on TTL expiry
					await scriptService.invalidateWorkflow(payload.orgId, payload.definitionId);
					logger.debug(
						`[CONSUMER:DefinitionInvalidate] Workflow-level invalidation for definition ${payload.definitionId} (org: ${payload.orgId})`
					);
				}
			} catch (error) {
				const errorMessage = error instanceof Error ? error.message : String(error);
				logger.error(`[CONSUMER:DefinitionInvalidate] Failed to process invalidation: ${errorMessage}`);
			}
		},
	});

	logger.info('[CONSUMER:DefinitionInvalidate] FANOUT consumer initialized');
}
