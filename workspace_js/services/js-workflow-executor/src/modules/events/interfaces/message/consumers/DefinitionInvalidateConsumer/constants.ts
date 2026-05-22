import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for DefinitionInvalidateConsumer (FANOUT).
 *
 * Stream and subject names resolve at module load from GO_ENV via the local
 * naming helpers. FANOUT is a platform-wide broadcast bus shared by every
 * service that publishes cache invalidation under fanout.>.
 */

/** FANOUT stream name — resolves to e.g. "DEV-MAPEXOS-FANOUT". */
export const FANOUT_STREAM = streamName('FANOUT', '');

/** FANOUT subject for workflow definition invalidation — resolves to e.g.
 * "dev.mapexos.fanout.workflow.definition.invalidate". */
export const DEFINITION_INVALIDATE_SUBJECT = subject('fanout', 'workflow.definition.invalidate');

/** Durable consumer name — resolves to e.g. "dev-jsworkflowexecutor-definition-invalidate-consumer". */
export const DEFINITION_INVALIDATE_DURABLE = durable('jsworkflowexecutor', 'definition-invalidate');
