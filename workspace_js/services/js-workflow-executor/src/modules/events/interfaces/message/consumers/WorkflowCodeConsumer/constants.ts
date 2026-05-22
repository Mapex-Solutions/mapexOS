import { streamName, subject, durable } from '@shared/configuration/naming';

/**
 * Constants for WorkflowCodeConsumer.
 *
 * Stream and subject names resolve at module load from GO_ENV via the local
 * naming helpers — same binary serves multiple environments on a shared cluster.
 */

/** Stream name for workflow code execution requests — resolves to e.g. "DEV-MAPEXOS-JSWORKFLOWEXECUTOR-CODE". */
export const WORKFLOW_JS_CODE_STREAM = streamName('JSWORKFLOWEXECUTOR', 'CODE');

/** Subject for workflow code execution — resolves to e.g. "dev.mapexos.workflow.js.code". */
export const WORKFLOW_JS_CODE_SUBJECT = subject('workflow', 'js.code');

/** Durable consumer name — resolves to e.g. "dev-jsworkflowexecutor-code-consumer". */
export const WORKFLOW_JS_CODE_DURABLE = durable('jsworkflowexecutor', 'code');

/** Event type for DLQ metadata. */
export const WORKFLOW_JS_CODE_EVENT_TYPE = 'mapexos.workflow.js.code';
