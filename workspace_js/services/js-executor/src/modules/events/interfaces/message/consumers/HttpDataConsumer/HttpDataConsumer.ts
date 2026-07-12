import type { Message } from '@mapexos/infrastructure';
import type { JsExecuteConsumerDeps } from './JsExecuteConsumer.types';

import { DeliverPolicy } from '@mapexos/infrastructure';

import { SERVICE_NAME, SERVICE_TYPE, DEFAULT_RETRY_POLICY, resolveConsumerConfig } from '@shared/constants';
import { JS_EXECUTE_STREAM, JS_EXECUTE_SUBJECT, JS_EXECUTE_DURABLE, JS_EXECUTE_EVENT_TYPE } from './constants';

/**
 * JsExecute consumer — receives HTTP datasource events and delegates to ScriptService.
 *
 * Flow: receive → call service.handleHttpBatch() → ACK/Nack per result.
 *
 * @param deps - Consumer dependencies
 */
export async function initJsExecuteConsumer(deps: JsExecuteConsumerDeps): Promise<void> {
	const { natsBus, logger, scriptService, config } = deps;
	const consumerConfig = resolveConsumerConfig(config);

	logger.info({ config: consumerConfig }, '[CONSUMER:JsExecute] Initializing');

	await natsBus.startConsumer({
		stream: JS_EXECUTE_STREAM,
		subject: JS_EXECUTE_SUBJECT,
		durable: JS_EXECUTE_DURABLE,
		deliverPolicy: DeliverPolicy.All,
		batchSize: consumerConfig.batchSize,
		fetchTimeout: consumerConfig.fetchTimeout,
		maxAckPending: consumerConfig.maxAckPending,
		retryPolicy: DEFAULT_RETRY_POLICY,
		dlqPolicy: {
			serviceName: SERVICE_NAME,
			serviceType: SERVICE_TYPE,
			eventType: JS_EXECUTE_EVENT_TYPE,
		},
		batchMessageHandlerV2: async (messages: Message[]) => {
			deps.batchSize?.observe(messages.length);

			const results = await scriptService.handleHttpBatch(messages);

			for (const result of results) {
				const msg = messages[result.index];

				if (result.success) {
					msg.ack();
					deps.eventsProcessed?.inc({ consumer: 'js_execute', status: 'success' });
				} else if (result.isPermanent) {
					await msg.reject(result.error ?? 'Invalid message');
					deps.eventsProcessed?.inc({ consumer: 'js_execute', status: 'rejected' });
				} else if (result.isOOM) {
					msg.nack(new Error(result.error ?? 'V8 OOM')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'js_execute', status: 'failure' });
				} else {
					msg.nack(new Error(result.error ?? 'Processing error')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'js_execute', status: 'failure' });
				}
			}
		},
	});

	logger.info('[CONSUMER:JsExecute] Initialized');
}
