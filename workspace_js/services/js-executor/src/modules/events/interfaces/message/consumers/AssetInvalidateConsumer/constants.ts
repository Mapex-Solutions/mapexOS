import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for AssetInvalidateConsumer (FANOUT pattern).
 *
 * Stream and subject names resolve at module load from GO_ENV via the local
 * naming helpers. FANOUT is a platform-wide broadcast bus shared by every
 * service that publishes cache invalidation under fanout.>.
 */

/** Stream name for FANOUT events — resolves to e.g. "DEV-MAPEXOS-FANOUT". */
export const FANOUT_STREAM = streamName('FANOUT', '');

/** Subject for asset cache invalidation — resolves to e.g. "dev.mapexos.fanout.asset.invalidate". */
export const FANOUT_ASSET_SUBJECT = subject('fanout', 'asset.invalidate');

/** Durable consumer name for asset invalidation. */
export const FANOUT_ASSET_DURABLE = durable('jsexecutor', 'asset-invalidate');

/** Event type for logging. */
export const FANOUT_ASSET_EVENT_TYPE = 'mapexos.fanout.asset.invalidate';
