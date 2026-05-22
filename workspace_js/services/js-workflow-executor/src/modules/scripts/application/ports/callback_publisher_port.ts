import type { WorkflowScriptCallback } from './workflow_script_service_port';

/**
 * CallbackPublisherPort abstracts NATS callback publishing from the service.
 *
 * Implemented by NatsCallbackPublisher (infrastructure layer).
 * Used by WorkflowScriptService to publish execution results to WORKFLOW-RESUME.
 */
export interface CallbackPublisherPort {
	/**
	 * Publishes a callback result to the specified NATS subject.
	 *
	 * @param subject - NATS subject (e.g., workflow.resume.callback.{instanceId})
	 * @param callback - Callback payload with execution result
	 */
	publishCallback(subject: string, callback: WorkflowScriptCallback): Promise<void>;
}
