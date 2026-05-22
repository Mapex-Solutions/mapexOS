import type { NatsBus } from '@mapexos/infrastructure';
import type { CallbackPublisherPort } from '../../application/ports/callback_publisher_port';
import type { WorkflowScriptCallback } from '../../application/ports/workflow_script_service_port';

/**
 * NatsCallbackPublisher implements CallbackPublisherPort using NatsBus.
 *
 * Publishes workflow script execution results to WORKFLOW-RESUME stream
 * via the callback subject provided by the Go workflow service.
 */
export class NatsCallbackPublisher implements CallbackPublisherPort {
	constructor(private readonly natsBus: NatsBus) {}

	/**
	 * Publishes callback to the specified NATS subject.
	 *
	 * @param subject - NATS subject for callback delivery
	 * @param callback - Execution result payload
	 */
	async publishCallback(subject: string, callback: WorkflowScriptCallback): Promise<void> {
		await this.natsBus.publish(subject, callback);
	}
}
