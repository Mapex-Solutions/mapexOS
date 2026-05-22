import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for TemplateInvalidateConsumer (FANOUT pattern).
 *
 * Stream and subject names resolve at module load from GO_ENV via the local
 * naming helpers. FANOUT is a platform-wide broadcast bus shared by every
 * service that publishes cache invalidation under fanout.>.
 */

/** Stream name for FANOUT events — resolves to e.g. "DEV-MAPEXOS-FANOUT". */
export const FANOUT_STREAM = streamName('FANOUT', '');

/** Subject for template cache invalidation — resolves to e.g. "dev.mapexos.fanout.template.invalidate". */
export const FANOUT_TEMPLATE_SUBJECT = subject('fanout', 'template.invalidate');

/** Durable consumer name for template invalidation. */
export const FANOUT_TEMPLATE_DURABLE = durable('jsexecutor', 'template-invalidate');

/** Event type for logging. */
export const FANOUT_TEMPLATE_EVENT_TYPE = 'mapexos.fanout.template.invalidate';
