import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for MqttDataConsumer.
 *
 * Stream and subject names resolve at module load from GO_ENV via the
 * local naming helpers — the same binary serves multiple environments
 * on a shared NATS cluster.
 */

/** Stream name for MQTT telemetry data — resolves to e.g. "DEV-MAPEXOS-JSEXECUTOR-MQTTDATA". */
export const MQTT_DATA_STREAM = streamName('JSEXECUTOR', 'MQTTDATA');

/**
 * Subject for MQTT data. The topic is intentionally agnostic — every
 * routing field (orgId, assetUUID, clientId, ...) travels in the JSON
 * payload (IngressMessage) the broker plugin publishes. Pattern:
 * ${env}.mapexos.mqtt.data.> — the trailing wildcard absorbs whatever
 * tail tokens the broker uses for stream-side routing without the
 * consumer having to parse them.
 *
 * Flow:
 * 1. Device publishes via MQTT to the broker on events/{assetUUID}/{type}.
 * 2. mapex-broker-mqtt plugin auths the CONNECT, accepts the publish,
 *    and emits IngressMessage on ${env}.mapexos.mqtt.data.* on NATS Core.
 * 3. The js-executor's JetStream consumer captures the message and the
 *    handler reads orgId / assetUUID straight from the payload.
 */
export const MQTT_DATA_SUBJECT = subject('mqtt', 'data') + '.>';

/** Durable consumer name for MQTT data. */
export const MQTT_DATA_DURABLE = durable('jsexecutor', 'mqttdata');

/** Event type for DLQ metadata. */
export const MQTT_DATA_EVENT_TYPE = 'mqtt.data';
