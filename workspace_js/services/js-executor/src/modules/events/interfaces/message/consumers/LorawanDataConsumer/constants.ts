import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for LorawanDataConsumer.
 *
 * Stream and subject names resolve at module load from GO_ENV via the
 * local naming helpers — the same binary serves multiple environments
 * on a shared NATS cluster.
 */

/** Stream name for LoRaWAN telemetry data — resolves to e.g. "DEV-MAPEXOS-JSEXECUTOR-LORAWANDATA". */
export const LORAWAN_DATA_STREAM = streamName('JSEXECUTOR', 'LORAWANDATA');

/**
 * Subject for LoRaWAN data. The topic is intentionally agnostic — every
 * routing field (orgId, assetUUID) travels in the JSON payload
 * (LorawanUplinkEnvelope) the mapexLNS adapter publishes. Pattern:
 * ${env}.mapexos.lorawan.data.> — the trailing wildcard absorbs the
 * {orgId}.{assetUUID} tail the LNS uses for stream-side routing without
 * the consumer having to parse it.
 *
 * Flow:
 * 1. Device sends a LoRaWAN uplink; The Things Stack NS decrypts the
 *    FRMPayload and hands it to the application uplink queue.
 * 2. The mapexLNS dispatcher (the ACL adapter implementing the TTS
 *    ApplicationUplinkQueue port) translates the ApplicationUp into a
 *    LorawanUplinkEnvelope and publishes it on
 *    ${env}.mapexos.lorawan.data.{orgId}.{assetUUID} on NATS Core.
 * 3. The js-executor's JetStream consumer captures the message and the
 *    handler reads orgId / assetUUID straight from the payload.
 */
export const LORAWAN_DATA_SUBJECT = subject('lorawan', 'data') + '.>';

/** Durable consumer name for LoRaWAN data. */
export const LORAWAN_DATA_DURABLE = durable('jsexecutor', 'lorawandata');

/** Event type for DLQ metadata. */
export const LORAWAN_DATA_EVENT_TYPE = 'lorawan.data';
