import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for the HTTP-telemetry consumer (http_gateway -> js-executor).
 *
 * Mirrors MqttData / LorawanData: STATIC subject, identity in the payload.
 * Stream and subject names resolve at module load from GO_ENV via the local
 * naming helpers.
 */

/** Stream name for HTTP telemetry data — resolves to e.g. "DEV-MAPEXOS-JSEXECUTOR-HTTPDATA". */
export const HTTP_DATA_STREAM = streamName('JSEXECUTOR', 'HTTPDATA');

/** STATIC subject for HTTP data — resolves to e.g. "dev.mapexos.http.data". */
export const HTTP_DATA_SUBJECT = subject('http', 'data');

/** Durable consumer name for HTTP data. */
export const HTTP_DATA_DURABLE = durable('jsexecutor', 'httpdata');

/** Event type for DLQ metadata. */
export const HTTP_DATA_EVENT_TYPE = 'http.data';
