import type { NatsBus } from '@mapexos/infrastructure';
import type { Logger } from '@mapexos/microservices';
import type { TemplateInvalidatePayload } from '@mapexos/schemas';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';

import { FANOUT_STREAM, FANOUT_TEMPLATE_SUBJECT } from './constants';

/**
 * Dependencies for TemplateInvalidateConsumer
 */
export interface TemplateInvalidateConsumerDeps {
	natsBus: NatsBus;
	logger: Logger;
	templateCachePort: TemplateCachePort;
	serviceName: string;
}

/**
 * Initializes the Template Invalidate FANOUT consumer for TieredCache invalidation.
 *
 * FANOUT Pattern:
 * - Each service instance receives a copy of the message (no queue group)
 * - Used for cache invalidation across all replicas
 * - Ephemeral consumer (not durable) - created fresh on each startup
 *
 * TieredCache Architecture:
 *   L0 (RAM): Hot cache - cleared on invalidation
 *   L1 (Disk): Persistent cache - cleared on invalidation
 *   L2 (MinIO): Source of truth - NOT affected
 *
 * Flow:
 * 1. Template service updates MinIO read model
 * 2. Template service publishes FANOUT invalidation with { orgId, templateId }
 * 3. All service instances receive the message
 * 4. Each instance clears L0+L1 using key: {orgId}/{templateId}
 * 5. Next request fetches fresh data from L2 → populates L0/L1
 *
 * @param deps - Consumer dependencies
 */
export async function initTemplateInvalidateConsumer(deps: TemplateInvalidateConsumerDeps): Promise<void> {
	const { natsBus, logger, templateCachePort, serviceName } = deps;

	logger.info('[CONSUMER:TemplateInvalidate] Starting FANOUT subscription');

	await natsBus.subscribeFanout({
		stream: FANOUT_STREAM,
		serviceName,
		subject: FANOUT_TEMPLATE_SUBJECT,
		handler: async (data: Uint8Array) => {
			try {
				const payload: TemplateInvalidatePayload = JSON.parse(new TextDecoder().decode(data));

				// Validate required fields
				if (!payload.orgId || !payload.templateId) {
					logger.warn('[CONSUMER:TemplateInvalidate] Invalid payload: missing orgId or templateId');
					return;
				}

				const cacheKey = `${payload.orgId}/${payload.templateId}`;
				logger.info({ orgId: payload.orgId, templateId: payload.templateId }, '[CONSUMER:TemplateInvalidate] Invalidating cache');
				templateCachePort.invalidate(cacheKey);
			} catch (err) {
				logger.warn({ err }, '[CONSUMER:TemplateInvalidate] Error processing message');
			}
		},
	});

	logger.info('[CONSUMER:TemplateInvalidate] FANOUT subscription started');
}
