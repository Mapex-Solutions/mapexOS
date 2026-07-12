import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for JsExecuteConsumer.
 *
 * Stream and subject names resolve at module load from GO_ENV via the local
 * naming helpers.
 */

/** Stream name for script execution — resolves to e.g. "DEV-MAPEXOS-JSEXECUTOR-PROCESS". */
export const JS_EXECUTE_STREAM = streamName('JSEXECUTOR', 'PROCESS');

/** Subject for script execution — resolves to e.g. "dev.mapexos.processor.js.execute". */
export const JS_EXECUTE_SUBJECT = subject('processor', 'js.execute');

/** Durable consumer name for script execution. */
export const JS_EXECUTE_DURABLE = durable('jsexecutor', 'process');

/** Event type for DLQ metadata. */
export const JS_EXECUTE_EVENT_TYPE = 'js.execute';
