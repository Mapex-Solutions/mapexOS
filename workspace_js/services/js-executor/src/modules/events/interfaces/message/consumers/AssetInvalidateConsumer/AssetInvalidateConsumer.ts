import type { NatsBus } from '@mapexos/infrastructure';
import type { Logger } from '@mapexos/microservices';
import type { AssetInvalidatePayload } from '@mapexos/schemas';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';

import { FANOUT_STREAM, FANOUT_ASSET_SUBJECT } from './constants';

/**
 * Dependencies for AssetInvalidateConsumer
 */
export interface AssetInvalidateConsumerDeps {
	natsBus: NatsBus;
	logger: Logger;
	assetCachePort: AssetCachePort;
	serviceName: string;
}

/**
 * Initializes the Asset Invalidate FANOUT consumer for TieredCache invalidation.
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
 * 1. Asset service updates MinIO read model
 * 2. Asset service publishes FANOUT invalidation with { orgId, assetUUID }
 * 3. All service instances receive the message
 * 4. Each instance clears L0+L1 using key: {orgId}/{assetUUID}
 * 5. Next request fetches fresh data from L2 → populates L0/L1
 *
 * @param deps - Consumer dependencies
 */
export async function initAssetInvalidateConsumer(deps: AssetInvalidateConsumerDeps): Promise<void> {
	const { natsBus, logger, assetCachePort, serviceName } = deps;

	logger.info('[CONSUMER:AssetInvalidate] Starting FANOUT subscription');

	await natsBus.subscribeFanout({
		stream: FANOUT_STREAM,
		serviceName,
		subject: FANOUT_ASSET_SUBJECT,
		handler: async (data: Uint8Array) => {
			try {
				const payload: AssetInvalidatePayload = JSON.parse(new TextDecoder().decode(data));

				// Validate required fields
				if (!payload.orgId || !payload.assetUUID) {
					logger.warn('[CONSUMER:AssetInvalidate] Invalid payload: missing orgId or assetUUID');
					return;
				}

				const cacheKey = `${payload.orgId}/${payload.assetUUID}`;
				logger.info({ orgId: payload.orgId, assetUUID: payload.assetUUID }, '[CONSUMER:AssetInvalidate] Invalidating cache');
				assetCachePort.invalidate(cacheKey);
			} catch (err) {
				logger.warn({ err }, '[CONSUMER:AssetInvalidate] Error processing message');
			}
		},
	});

	logger.info('[CONSUMER:AssetInvalidate] FANOUT subscription started');
}
