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
 * STATIC subject for LoRaWAN data: ${env}.mapexos.lorawan.data. Every routing
 * field (orgId, assetUUID) travels in the JSON payload (LorawanUplinkEnvelope)
 * the mapexLNS adapter publishes — never in the subject.
 *
 * Flow:
 * 1. Device sends a LoRaWAN uplink; The Things Stack NS hands the ApplicationUp
 *    to the uplink queue.
 * 2. The mapexLNS dispatcher (the ACL adapter implementing the TTS
 *    ApplicationUplinkQueue port) translates it into a LorawanUplinkEnvelope
 *    and publishes it on ${env}.mapexos.lorawan.data on NATS Core.
 * 3. The js-executor's JetStream consumer captures the message and the handler
 *    reads orgId / assetUUID straight from the payload.
 */
export const LORAWAN_DATA_SUBJECT = subject('lorawan', 'data');

/** Durable consumer name for LoRaWAN data. */
export const LORAWAN_DATA_DURABLE = durable('jsexecutor', 'lorawandata');

/** Event type for DLQ metadata. */
export const LORAWAN_DATA_EVENT_TYPE = 'lorawan.data';
