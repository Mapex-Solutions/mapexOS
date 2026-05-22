import type { Message } from '@mapexos/infrastructure';
import type { WorkflowCodeConsumerDeps } from './WorkflowCodeConsumer.types';
import type { WorkflowScriptInput } from '@modules/scripts/application/ports';

import { DeliverPolicy } from '@mapexos/infrastructure';

import { SERVICE_NAME, SERVICE_TYPE, DEFAULT_RETRY_POLICY, resolveConsumerConfig } from '@shared/constants';
import {
	WORKFLOW_JS_CODE_STREAM,
	WORKFLOW_JS_CODE_SUBJECT,
	WORKFLOW_JS_CODE_DURABLE,
	WORKFLOW_JS_CODE_EVENT_TYPE,
} from './constants';

import { OOMError } from '@modules/engine/domain/errors';

/**
 * Initializes the WorkflowCode consumer for script execution requests.
 *
 * Listens on WORKFLOW-JS-CODE stream for code node execution requests from
 * the workflow runtime. Each message contains context (event, state, inputs, nodes)
 * and identifies the script to execute via orgId/workflowId/nodeId.
 *
 * @param deps - Consumer dependencies
 */
export async function initWorkflowCodeConsumer(deps: WorkflowCodeConsumerDeps): Promise<void> {
	const { natsBus, logger, scriptService, config } = deps;
	const consumerConfig = resolveConsumerConfig(config);

	logger.info({ config: consumerConfig }, '[CONSUMER:WorkflowCode] Initializing');

	await natsBus.startConsumer({
		stream: WORKFLOW_JS_CODE_STREAM,
		subject: WORKFLOW_JS_CODE_SUBJECT,
		durable: WORKFLOW_JS_CODE_DURABLE,
		deliverPolicy: DeliverPolicy.All,
		batchSize: consumerConfig.batchSize,
		fetchTimeout: consumerConfig.fetchTimeout,
		maxAckPending: consumerConfig.maxAckPending,
		retryPolicy: DEFAULT_RETRY_POLICY,
		dlqPolicy: {
			serviceName: SERVICE_NAME,
			serviceType: SERVICE_TYPE,
			eventType: WORKFLOW_JS_CODE_EVENT_TYPE,
		},
		batchMessageHandlerV2: async (messages: Message[]) => {
			deps.batchSize?.observe(messages.length);

			for (const msg of messages) {
				const startTime = process.hrtime.bigint();

				try {
					const input = JSON.parse(new TextDecoder().decode(msg.data)) as WorkflowScriptInput;

					// Set DLQ context for multi-tenant filtering (MANDATORY)
					msg.orgId = input.orgId ?? '';
					msg.pathKey = input.pathKey ?? '';

					if (!input.orgId || !input.workflowId || !input.nodeId || !input.instanceId) {
						logger.warn('[CONSUMER:WorkflowCode] Invalid message: missing required fields');
						msg.ack();
						continue;
					}

					const result = await scriptService.execute(input);

					if (!result.success) {
						deps.executionsTotal?.inc({ status: 'error' });
						msg.ack(); // ACK to prevent infinite retry — error callback already published
						continue;
					}

					deps.executionsTotal?.inc({ status: 'success' });
					msg.ack();
				} catch (error) {
					if (error instanceof OOMError) {
						// OOM is transient — NACK for retry with backoff
						logger.warn(`[CONSUMER:WorkflowCode] OOM detected, NACKing for retry: ${error.message}`);
						msg.nack(error);
						deps.executionsTotal?.inc({ status: 'oom' });
					} else {
						const err = error instanceof Error ? error : new Error(String(error));
						logger.error(`[CONSUMER:WorkflowCode] Processing failed: ${err.message}`);
						msg.ack(); // ACK to prevent infinite retry — error callback already published
						deps.executionsTotal?.inc({ status: 'error' });
					}
				} finally {
					const elapsed = Number(process.hrtime.bigint() - startTime) / 1e9;
					deps.executionDuration?.observe(elapsed);
				}
			}
		},
	});

	logger.info('[CONSUMER:WorkflowCode] Initialized');
}
